package models

type Requirements struct {
	Id              int    `json:"id,omitempty"`
	Request         string `json:"request"`
	ExpectedOutcome string `json:"outcome"`
	Status          bool   `json:"status"`
	OrderID         string `json:"order_id,omitempty"`
}

type RequirementPatch struct {
	Patches []Patch `json:"patch"`
}

type Patch struct {
	Id              int    `json:"id"`
	ExpectedOutcome string `json:"outcome"`
}

type ProgressForm struct {
	Fufillments []Fufillment `json:"fufillments"`
}

type Fufillment struct {
	Requirementid int    `json:"reqid"`
	Outcome       string `json:"outcome"`
}

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
