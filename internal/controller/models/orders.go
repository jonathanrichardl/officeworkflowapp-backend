package models

import "order-validation-v2/internal/entity"

type Orders struct {
	ID           string         `json:"id,omitempty"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	Deadline     string         `json:"deadline"`
	Requirements []Requirements `json:"requirements"`
}

type OrderPatch struct {
	ID          string  `json:"id"`
	Title       *string `json:"new_title"`
	Description *string `json:"new_description"`
	Deadline    *string `json:"new_deadline"`
}

func BuildPayload(O []*entity.Orders) []*Orders {
	var response []*Orders
	for _, o := range O {
		r := Orders{
			ID:          o.ID,
			Description: o.Description,
			Title:       o.Title,
			Deadline:    o.Deadline.Format("2 Jan 2006"),
		}

		response = append(response, &r)

	}
	return response

}

func (o *Orders) AddRequirements(R []*entity.Requirements) {
	var requirements []Requirements
	for _, r := range R {
		requirement := Requirements{
			Id:              r.Id,
			OrderID:         r.OrderID,
			ExpectedOutcome: r.ExpectedOutcome,
			Request:         r.Request,
			Status:          r.Status,
		}
		requirements = append(requirements, requirement)

	}
	o.Requirements = requirements

}
