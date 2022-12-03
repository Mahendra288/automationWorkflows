package interactorTests

import (
	"automationWorkflows/customErrors"
	"automationWorkflows/enums"
	"automationWorkflows/interactors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type executeWorkflowStorageInterfacesMock struct {
	mock.Mock
}

func (o *executeWorkflowStorageInterfacesMock) GetWorkflowExecStatusFromRequestId(execReqUniqueId string) (string, error) {
	args := o.Called(execReqUniqueId)
	return enums.WorkFlowExecLogStatusEnum().Done, args.Error(1)
}

func (o *executeWorkflowStorageInterfacesMock) GetLatestPublishedWorkflowExecConfigId(workflowId string) (string, error) {
	args := o.Called(workflowId)
	return "1223", args.Error(1)
}

func TestAlreadyWorkflowExecutionConfigDone(t *testing.T) {
	// Arrange
	execReqUniqueId := "12341"
	reqDetails := interactors.WorkflowExecReqStruct{
		WorkflowId: "1", SourceId: "Source1", WorkflowExecReqUniqueId: execReqUniqueId}

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
