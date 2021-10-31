package models

import "order-validation-v2/internal/entity"

type Tasks struct {
	Id               string        `json:"id"`
	Request          string        `json:"task"`
	ExpectedOutcome  string        `json:"outcome"`
	Status           entity.Status `json:"status"`
	OrderTitle       string        `json:"order_title"`
	OrderDescription string        `json:"order_description"`
	OrderDeadline    string        `json:"order_deadline"`
}

func BuildTasks(T []*entity.Task) []*Tasks {
	var tasks []*Tasks
	for _, t := range T {
		task := Tasks{
			Id:               t.ID,
			ExpectedOutcome:  t.ExpectedOutcome,
			Request:          t.Request,
			Status:           t.Status,
			OrderTitle:       t.OrderTitle,
			OrderDescription: t.OrderDescription,
			OrderDeadline:    t.OrderDeadline,
		}
		tasks = append(tasks, &task)

	}
	return tasks
}
