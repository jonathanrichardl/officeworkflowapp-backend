package models

type Requirements struct {
	Id              int    `json:"id"`
	Request         string `json:"request"`
	ExpectedOutcome string `json:"outcome"`
	Status          bool   `json:"status"`
	OrderID         string `json:"order_id"`
}

type ProgressForm struct {
	Fufillments []Fufillment `json:"fufillments"`
}

type Fufillment struct {
	Requirementid int    `json:"reqid"`
	Outcome       string `json:"outcome"`
}
