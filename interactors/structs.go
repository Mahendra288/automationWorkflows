package interactors

import "automationWorkflows/enums"

type TriggerEventDetailsStruct struct {
	eventEntity   string
	eventEntityId string
	payload       map[string]any
}
type WorkflowExecReqStruct struct {
	workflowId              string
	sourceId                string
	sourceType              enums.WorkFlowSourceType
	workflowExecReqUniqueId string
	triggerEventDetails     TriggerEventDetailsStruct
}
