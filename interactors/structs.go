package interactors

import "automationWorkflows/enums"

type TriggerEventDetailsStruct struct {
	EventEntity   string
	EventEntityId string
	Payload       map[string]string
}
type WorkflowExecReqStruct struct {
	WorkflowId              string
	SourceId                string
	SourceType              enums.WorkFlowSourceType
	WorkflowExecReqUniqueId string
	TriggerEventDetails     TriggerEventDetailsStruct
}

type InitialWorkflowExecLogStruct struct {
	ExecLogId     string
	ExecRequestId string
	SourceId      string
	SourceType    enums.WorkFlowSourceType
	ExecConfigId  string
	Payload       string
}
