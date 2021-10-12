package entity

import (
	"time"
)

type Orders struct {
	ID          string
	Title       string
	Description string
	Deadline    time.Time
}

func NewOrder(title string, description string, deadline time.Time) *Orders {
	o := &Orders{
		ID:          NewUUID().String(),
		Title:       title,
		Description: description,
		Deadline:    deadline,
	}
	return o
}
