package entity

type Requirements struct {
	Id              int
	Request         string
	ExpectedOutcome string
	Status          bool
	OrderID         ID
}

func NewRequirement(request string, expectedOutcome string, orderID ID) *Requirements {
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
