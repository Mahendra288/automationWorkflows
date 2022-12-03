package interactors

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
	err := interactor.validateForPrevExecution(reqDetails.workflowExecReqUniqueId)
	if err != nil {
		return err
	}

	latestExecConfigId, err := interactor.getLatestPublishedWorkflowExecConfigId(reqDetails.workflowId)
	if err != nil {
		return err
	}
	_, err = interactor.executeWorkflowNodes(latestExecConfigId, reqDetails)
	if err != nil {
		return err
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
	latestExecConfigId string, reqDetails WorkflowExecReqStruct,
) ([]string, error) {
	var execFailedNodeIds []string
	leadDetails, leadId, err := interactor.getLatestLeadDetails(
		reqDetails.triggerEventDetails.eventEntity,
		reqDetails.triggerEventDetails.eventEntityId,
	)
	if err != nil {
		return execFailedNodeIds, err
	}
	err = interactor.executeNode(leadDetails, leadId)
	if err != nil {
		return execFailedNodeIds, err
	}
	return execFailedNodeIds, err
}

func (interactor ExecuteWorkflowInteractor) getLatestLeadDetails(
	entityType string, entityId string) (leadDetails map[string]string, leadId string, err error) {
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
	leadDetails map[string]string, leadId string) error {
	return nil
}
