package entity

import (
	"time"
)

type Status int8

const (
	Unfinished Status = iota
	InReview
	Finished
)

type Task struct {
	AssignerID        string
	ID                string
	Note              string
	RequirementID     int
	UserID            string
	Status            Status
	NumOfPrerequisite uint8
	Prerequisites     []string
	Allowed           bool
	Deadline          time.Time
}

type TaskWithDetails struct {
	ID               string
	AssignedBy       string
	Username         string
	Note             string
	RequirementID    int
	Request          string
	ExpectedOutcome  string
	UserID           string
	Status           Status
	OrderTitle       string
	OrderDescription string
	OrderDeadline    string
}

func NewTask(assignerID string, requirementID int, userID string, note string, prerequisiteTaskID []string, deadline time.Time) *Task {
	taskID := NewUUID().String()
	task := Task{
		AssignerID:        assignerID,
		ID:                taskID,
		RequirementID:     requirementID,
		UserID:            userID,
		Note:              note,
		Allowed:           true,
		Deadline:          deadline,
		NumOfPrerequisite: 0,
		Status:            0,
	}
	if totalPrerequisite := len(prerequisiteTaskID); totalPrerequisite != 0 {
		task.Allowed = false
		task.NumOfPrerequisite = uint8(totalPrerequisite)
		task.Prerequisites = prerequisiteTaskID
	}
	return &task
}

func (t *Task) SetStatus(newStatus Status) {
	t.Status = newStatus
}

func (t *Task) ReducePrerequisite() {
	t.NumOfPrerequisite--
}

func (t *Task) Allow() {
	t.Allowed = true
}
