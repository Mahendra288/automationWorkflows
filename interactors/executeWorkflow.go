package interactors

import (
	"encoding/json"
	"github.com/google/uuid"
)
import (
	"automationWorkflows/adapters"
	"automationWorkflows/customErrors"
	"automationWorkflows/enums"
)

// todo: naming style, whether to keep interfaces as suffix or not
type executeWorkflowStorageInterfaces interface {
	GetWorkflowExecStatusFromRequestId(execReqUniqueId string) (enums.WorkFlowExecLogStatus, error)
	GetLatestPublishedWorkflowExecConfigId(workflowId string) (string, error)
}

type ExecuteWorkflowInteractor struct {
	Storage executeWorkflowStorageInterfaces
}

func (interactor ExecuteWorkflowInteractor) ExecuteWorkflow(reqDetails WorkflowExecReqStruct) error {
	err := interactor.validateForPrevExecution(reqDetails.WorkflowExecReqUniqueId)
	if err != nil {
		return err
	}

	latestExecConfigId, err := interactor.getLatestPublishedWorkflowExecConfigId(reqDetails.WorkflowId)
	if err != nil {
		return err
	}
	initialExecLog, logInteractor, err := interactor.createInitialWorkflowExecLog(latestExecConfigId, reqDetails)
	if err != nil {
		return err
	}

	execFailedNodeIds, err := interactor.executeWorkflowNodes(latestExecConfigId, reqDetails, initialExecLog)
	if err != nil {
		logInteractor.updateFailedStatusForExecLog(initialExecLog.ExecLogId, err)
		return err
	}

	hasAnyNodesFailed := len(execFailedNodeIds) > 0
	if hasAnyNodesFailed {
		traceback := map[string][]string{"failed_node_ids": execFailedNodeIds}
		logInteractor.updatePartiallyFailedStatusForExecLog(initialExecLog.ExecLogId, traceback)
	} else {
		logInteractor.updateSuccessStatusForExecLog(initialExecLog.ExecLogId)
	}

	return nil
}

func (interactor ExecuteWorkflowInteractor) validateForPrevExecution(execReqUniqueId string) error {
	execLogStatus, err := interactor.Storage.GetWorkflowExecStatusFromRequestId(execReqUniqueId)
	if err != nil {
		return err
	}
	if execLogStatus != (enums.WorkFlowExecLogStatus{}) {
		return &customErrors.AlreadyExecutedWorkflowError{}
	}
	return nil
}

func (interactor ExecuteWorkflowInteractor) getLatestPublishedWorkflowExecConfigId(workflowId string) (string, error) {
	latestExecConfigId, err := interactor.Storage.GetLatestPublishedWorkflowExecConfigId(workflowId)
	if err != nil {
		return latestExecConfigId, err
	}
	return latestExecConfigId, nil
}

func (interactor ExecuteWorkflowInteractor) executeWorkflowNodes(
	latestExecConfigId string,
	reqDetails WorkflowExecReqStruct,
	initialExecLog InitialWorkflowExecLogStruct,
) ([]string, error) {
	var execFailedNodeIds []string

	leadDetails, leadId, err := interactor.getLatestLeadDetails(
		reqDetails.TriggerEventDetails.EventEntity,
		reqDetails.TriggerEventDetails.EventEntityId,
	)
	if err != nil {
		return execFailedNodeIds, err
	}

	leadDetails = interactor.mergePayloadWithLeadDetailsForLeadTrigger(
		reqDetails.TriggerEventDetails.EventEntity,
		reqDetails.TriggerEventDetails.Payload,
		leadDetails)

	err = interactor.executeNode(leadDetails, leadId)
	if err != nil {
		return execFailedNodeIds, err
	}
	return execFailedNodeIds, err
}

func (interactor ExecuteWorkflowInteractor) getLatestLeadDetails(
	entityType string, entityId string) (leadDetails map[string]any, leadId string, err error) {
	if entityType == enums.TriggerEventEntityEnum().Lead {
		leadId = entityId
	} else if entityType == enums.TriggerEventEntityEnum().Activity {
		leadId, err = adapters.SalesCrmService{}.GetLeadIdForActivityId(entityId)
		if err != nil {
			return leadDetails, leadId, err
		}
	} else if entityType == enums.TriggerEventEntityEnum().Task {
		leadId, err = adapters.SalesCrmService{}.GetLeadIdForTaskId(entityId)
		if err != nil {
			return leadDetails, leadId, err
		}
	}
	return leadDetails, leadId, err
}

func (interactor ExecuteWorkflowInteractor) executeNode(
	leadDetails map[string]any, leadId string) error {
	return nil
}

func (interactor ExecuteWorkflowInteractor) createInitialWorkflowExecLog(
	latestExecConfigId string, reqDetails WorkflowExecReqStruct) (
	InitialWorkflowExecLogStruct, WorkflowExecLogInteractor, error,
) {

	reqDetailsMap, _ := json.Marshal(reqDetails)
	logInteractor := WorkflowExecLogInteractor{}

	initialLog := InitialWorkflowExecLogStruct{
		ExecLogId:     uuid.NewString(),
		ExecRequestId: reqDetails.WorkflowExecReqUniqueId,
		SourceId:      reqDetails.SourceId,
		SourceType:    reqDetails.SourceType,
		ExecConfigId:  latestExecConfigId,
		Payload:       string(reqDetailsMap),
	}
	err := logInteractor.CreateInitialWorkflowExecLog(initialLog)
	if err != nil {
		return InitialWorkflowExecLogStruct{}, logInteractor, err
	}
	return initialLog, logInteractor, err
}

func (interactor ExecuteWorkflowInteractor) mergePayloadWithLeadDetailsForLeadTrigger(
	triggerEventEntity string,
	triggerEventPayload map[string]any,
	leadDetails map[string]any) map[string]any {
	isLeadTrigger := triggerEventEntity == enums.TriggerEventEntityEnum().Lead
	if !isLeadTrigger {
		return leadDetails
	}
	for leadFieldId, response := range triggerEventPayload {
		leadDetails[leadFieldId] = response
	}
	return leadDetails
}
