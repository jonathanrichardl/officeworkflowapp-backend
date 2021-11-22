package entity

const (
	NotAssigned Status = iota
	Assigned
	AssignedAndFinished
)

type Requirements struct {
	Id              int
	Request         string
	ExpectedOutcome string
	OrderID         string
	Status          Status
	TaskID          string
}

func NewRequirement(request string, expectedOutcome string, orderID string) *Requirements {
	r := &Requirements{
		Request:         request,
		ExpectedOutcome: expectedOutcome,
		OrderID:         orderID,
		Status:          0,
	}
	return r

}

func (r *Requirements) Assign(taskID string) {
	r.TaskID = taskID
	r.setStatus(1)

}
func (r *Requirements) setStatus(newStatus Status) {
	r.Status = newStatus
}
