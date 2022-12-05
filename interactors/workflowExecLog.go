package interactors

type WorkflowExecLogInteractor struct {
}

func (interactor WorkflowExecLogInteractor) CreateInitialWorkflowExecLog(
	execLogStruct InitialWorkflowExecLogStruct) error {
	return nil
}

func (interactor WorkflowExecLogInteractor) updateFailedStatusForExecLog(execLogId string, err error) {
	// todo: figure out how to get traceback here from error object
}

func (interactor WorkflowExecLogInteractor) updatePartiallyFailedStatusForExecLog(
	execLogId string, traceback map[string][]string) {

}

func (interactor WorkflowExecLogInteractor) updateSuccessStatusForExecLog(execLogId string) {

}
