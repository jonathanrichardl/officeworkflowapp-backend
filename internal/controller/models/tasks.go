package models

import "order-validation-v2/internal/entity"

type ReviewForm struct {
	Message   string   `json:"message"`
	Approved  bool     `json:"approved"`
	ForwardTo []string `json:"forwarded"`
}

type TaskWithDetail struct {
	Id               string        `json:"id"`
	Note             string        `json:"note"`
	User             string        `json:"assigned_user,omitempty"`
	Username         string        `json:"assigned_username,omitempty"`
	Request          string        `json:"task"`
	ExpectedOutcome  string        `json:"outcome,omitempty"`
	Prerequisites    []string      `json:"prerequisites,omitempty"`
	Status           entity.Status `json:"status"`
	TaskDeadline     string        `json:"deadline,omitempty"`
	OrderTitle       string        `json:"order_title,omitempty"`
	OrderDescription string        `json:"order_description,omitempty"`
	OrderDeadline    string        `json:"order_deadline,omitempty"`
	Feedbacks        []Feedback    `json:"feedbacks"`
}

type Feedback struct {
	UserName string `json:"from"`
	Message  string `json:"message"`
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
			Note:             t.Note,
			User:             t.UserID,
			Username:         t.Username,
			ExpectedOutcome:  t.ExpectedOutcome,
			Request:          t.Request,
			TaskDeadline:     t.Deadline.Format("2/Jan/2006 15:04:05"),
			Status:           t.Status,
			OrderTitle:       t.OrderTitle,
			OrderDescription: t.OrderDescription,
			OrderDeadline:    t.OrderDeadline.Format("2/Jan/2006 15:04:05"),
			Prerequisites:    t.Prerequisites,
		}
		var feedbacks []Feedback
		for _, review := range t.Messages {
			var feedback Feedback
			feedback.UserName = review.Username
			feedback.Message = review.Message
		}
		task.Feedbacks = feedbacks
		tasks = append(tasks, &task)

	}
	return tasks
}
