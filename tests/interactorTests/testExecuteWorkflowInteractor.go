package interactorTests

import (
	"automationWorkflows/enums"
	"automationWorkflows/interactors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type TestExecuteWorkflowInteractor struct {
}

type executeWorkflowStorageInterfacesMock struct {
	mock.Mock
}

func (o *executeWorkflowStorageInterfacesMock) GetWorkflowExecStatusFromRequestId(execReqUniqueId string) (enums.WorkFlowExecLogStatus, error) {
	args := o.Called(execReqUniqueId)
	return args.Get(0).(enums.WorkFlowExecLogStatus), args.Error(1)
}

func (o *executeWorkflowStorageInterfacesMock) GetLatestPublishedWorkflowExecConfigId(workflowId string) (string, error) {
	args := o.Called(workflowId)
	return "1223", args.Error(1)
}

func (testInteractor TestExecuteWorkflowInteractor) TestAlreadyWorkflowExecutionConfigDone(t *testing.T) {
	// Arrange
	reqDetails := interactors.WorkflowExecReqStruct{}

	// Mocking storage method
	StorageInterfaceMockObj := new(executeWorkflowStorageInterfacesMock)
	StorageInterfaceMockObj.On(
		"GetWorkflowExecStatusFromRequestId", reqDetails,
	).Return(enums.WorkFlowExecLogStatusEnum().Created, nil)

	interactor := interactors.ExecuteWorkflowInteractor{
		Storage: StorageInterfaceMockObj,
	}

	// Act
	err := interactor.ExecuteWorkflow(reqDetails)

	// Assert
	assert.Equal(t, expectedResponse, response)
}
