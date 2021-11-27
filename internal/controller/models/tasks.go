package models

import "order-validation-v2/internal/entity"

type ReviewForm struct {
	Message   string   `json:"message"`
	Approved  bool     `json:"approved"`
	ForwardTo []string `json:"forwarded"`
}

type TaskWithDetail struct {
	Id               string        `json:"id"`
	User             string        `json:"assigned_user,omitempty"`
	Request          string        `json:"task"`
	ExpectedOutcome  string        `json:"outcome"`
	Status           entity.Status `json:"status,omitempty"`
	OrderTitle       string        `json:"order_title"`
	OrderDescription string        `json:"order_description"`
	OrderDeadline    string        `json:"order_deadline"`
}

type NewTask struct {
	Num           string   `json:"num,omitempty"`
	RequirementID int      `json:"requirement_id"`
	Note          string   `json:"note"`
	UserID        string   `json:"user_id"`
	Prerequisite  []string `json:"prerequisite"`
	Deadline      string   `json:"deadline"`
}

type BulkAddedTasks struct {
	Tasks []NewTask `json:"tasks"`
}

func BuildTasks(T []*entity.TaskWithDetails) []*TaskWithDetail {
	var tasks []*TaskWithDetail
	for _, t := range T {
		task := TaskWithDetail{
			Id:               t.ID,
			User:             t.Username,
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
