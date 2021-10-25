package entity

type Requirements struct {
	Id              int
	Request         string
	ExpectedOutcome string
	OrderID         string
	UserID          string
	Status          bool
}

func NewRequirement(request string, expectedOutcome string, orderID string, userID *string) *Requirements {
	r := &Requirements{
		Request:         request,
		ExpectedOutcome: expectedOutcome,
		OrderID:         orderID,
	}
	if userID != nil {
		r.UserID = *userID
	}
	return r

}

func (r *Requirements) SetStatus(newStatus bool) {
	r.Status = newStatus

}

func (r *Requirements) AssignUser(UserID string) {
	r.UserID = UserID
}
