package entity

import (
	"time"

	"github.com/google/uuid"
)

type Orders struct {
	ID          uuid.UUID
	Title       string
	Description string
	Deadline    time.Time
}

func NewOrder(title string, description string, deadline time.Time) *Orders {
	o := &Orders{
		ID:          NewUUID(),
		Title:       title,
		Description: description,
		Deadline:    deadline,
	}
	return o
}
