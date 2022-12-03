package customErrors

type AlreadyExecutedWorkflowError struct {
}

func (err *AlreadyExecutedWorkflowError) Error() string {
	return "Workflow has already been executed with given unique request id."
}
