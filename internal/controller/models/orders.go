package models

import "order-validation-v2/internal/entity"

type Orders struct {
	ID           string         `json:"id,omitempty"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	Deadline     string         `json:"deadline"`
	Requirements []Requirements `json:"requirements"`
}

type Tasks struct {
	Id              int    `json:"id"`
	Request         string `json:"task"`
	ExpectedOutcome string `json:"outcome"`
	Status          bool   `json:"status"`
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

func BuildTasks(R []*entity.Requirements) []*Tasks {
	var tasks []*Tasks
	for _, r := range R {
		task := Tasks{
			Id:              r.Id,
			ExpectedOutcome: r.ExpectedOutcome,
			Request:         r.Request,
			Status:          r.Status,
		}
		tasks = append(tasks, &task)

	}
	return tasks
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
