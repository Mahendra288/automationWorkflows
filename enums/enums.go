package enums

type TriggerEventEntity struct {
	Lead     string
	Activity string
	Task     string
}

func TriggerEventEntityEnum() *TriggerEventEntity {
	return &TriggerEventEntity{
		Lead:     "LEAD",
		Activity: "ACTIVITY",
		Task:     "TASK",
	}
}

type WorkFlowSourceType struct {
	AsyncEvents string
	IbEvents    string
}

func WorkFlowSourceTypeEnum() *WorkFlowSourceType {
	return &WorkFlowSourceType{
		AsyncEvents: "ASYNC_EVENTS",
		IbEvents:    "IB_EVENTS",
	}
}

type WorkFlowExecLogStatus struct {
	Created         string
	Processing      string
	Failed          string
	PartiallyFailed string
	Done            string
}

func WorkFlowExecLogStatusEnum() *WorkFlowExecLogStatus {
	return &WorkFlowExecLogStatus{
		Created:         "CREATED",
		Processing:      "PROCESSING",
		Failed:          "FAILED",
		PartiallyFailed: "PARTIALLY_FAILED",
		Done:            "DONE",
	}
}
