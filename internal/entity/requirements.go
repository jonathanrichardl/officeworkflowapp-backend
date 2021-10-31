package entity

type Requirements struct {
	Id              int
	Request         string
	ExpectedOutcome string
	OrderID         string
	Status          bool
}

func NewRequirement(request string, expectedOutcome string, orderID string) *Requirements {
	r := &Requirements{
		Request:         request,
		ExpectedOutcome: expectedOutcome,
		OrderID:         orderID,
	}
	return r

}

func (r *Requirements) SetStatus(newStatus bool) {
	r.Status = newStatus

}
