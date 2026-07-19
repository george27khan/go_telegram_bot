package handler

type TimeHandler struct {
	SchedulerEmployeeHandler SchedulerEmployeeHandler
}

func NewTimeHandler(schedulerEmployeeHandler SchedulerEmployeeHandler) *TimeHandler {
	return &TimeHandler{
		SchedulerEmployeeHandler: schedulerEmployeeHandler,
	}
}
