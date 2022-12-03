package interactorTests

import (
	"automationWorkflows/customErrors"
	"automationWorkflows/enums"
	"automationWorkflows/interactors"
	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type executeWorkflowStorageInterfacesMock struct {
	mock.Mock
}

func (o *executeWorkflowStorageInterfacesMock) GetWorkflowExecStatusFromRequestId(execReqUniqueId string) (string, error) {
	args := o.Called(execReqUniqueId)
	return args.Get(0).(string), args.Error(1)
}

func (o *executeWorkflowStorageInterfacesMock) GetLatestPublishedWorkflowExecConfigId(workflowId string) (string, error) {
	args := o.Called(workflowId)
	return args.Get(0).(string), args.Error(1)
}

func TestAlreadyWorkflowExecutionConfigIsInProgress(t *testing.T) {
	// Arrange
	execReqUniqueId := "12341"
	workflowId := "workflow1"
	reqDetails := interactors.WorkflowExecReqStruct{
		WorkflowId: workflowId, SourceId: "Source1", WorkflowExecReqUniqueId: execReqUniqueId}

	// Mocking storage method
	StorageInterfaceMockObj := new(executeWorkflowStorageInterfacesMock)
	StorageInterfaceMockObj.On(
		"GetWorkflowExecStatusFromRequestId", execReqUniqueId,
	).Return(enums.WorkFlowExecLogStatusEnum().Created, nil)

	interactor := interactors.ExecuteWorkflowInteractor{
		Storage: StorageInterfaceMockObj,
	}

	// Act
	err := interactor.ExecuteWorkflow(reqDetails)

	// Assert
	assert.Equal(t, &customErrors.AlreadyExecutedWorkflowError{}, err)
}

func TestValidDetailsExecutesWorkflow(t *testing.T) {
	// Arrange
	execReqUniqueId := "12341"
	workflowId := "workflow1"
	latestExecConfigId := "LatestWorkflowExecutionConfigId1"

	payload := map[string]string{}
	triggerEventDetails := interactors.TriggerEventDetailsStruct{
		EventEntity:   enums.TriggerEventEntityEnum().Lead,
		EventEntityId: "1",
		Payload:       payload,
	}
	reqDetails := interactors.WorkflowExecReqStruct{
		WorkflowId: workflowId, SourceId: "Source1", WorkflowExecReqUniqueId: execReqUniqueId,
		TriggerEventDetails: triggerEventDetails}
	//
	//initialLog := interactors.InitialWorkflowExecLogStruct{
	//	ExecLogId:     uuid.NewString(),
	//	ExecRequestId: reqDetails.WorkflowExecReqUniqueId,
	//	SourceId:      reqDetails.SourceId,
	//	SourceType:    reqDetails.SourceType,
	//	ExecConfigId:  latestExecConfigId,
	//	Payload:       string(payload),
	//}
	patches := gomonkey.ApplyFunc(
		interactors.WorkflowExecLogInteractor{}.CreateInitialWorkflowExecLog, func(execLogStruct interactors.InitialWorkflowExecLogStruct) error {
			return nil
		})

	defer patches.Reset()

	// Mocking storage method
	StorageInterfaceMockObj := new(executeWorkflowStorageInterfacesMock)
	StorageInterfaceMockObj.On(
		"GetWorkflowExecStatusFromRequestId", execReqUniqueId,
	).Return("", nil)
	StorageInterfaceMockObj.On(
		"GetLatestPublishedWorkflowExecConfigId", workflowId,
	).Return(latestExecConfigId, nil)
	interactor := interactors.ExecuteWorkflowInteractor{
		Storage: StorageInterfaceMockObj,
	}

	// Act
	err := interactor.ExecuteWorkflow(reqDetails)

	// Assert
	assert.Equal(t, err, nil)
}
