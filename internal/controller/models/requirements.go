package models

import "order-validation-v2/internal/entity"

type Requirements struct {
	Id              int           `json:"id,omitempty"`
	Request         string        `json:"request"`
	ExpectedOutcome string        `json:"outcome"`
	Status          entity.Status `json:"status"`
	OrderID         string        `json:"order_id,omitempty"`
}

type RequirementPatch struct {
	Patches []Patch `json:"patch"`
}

type Patch struct {
	Id              int     `json:"id"`
	ExpectedOutcome *string `json:"new_outcome"`
	Request         *string `json:"new_request"`
}
