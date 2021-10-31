package entity

import (
	"time"
)

type Submission struct {
	ID             string
	SubmissionTime time.Time
	Images         []SubmissionImage
	Message        string
	TaskID         string
}

func NewSubmission(Message string, Images [][]byte, TaskID string) *Submission {
	id := NewUUID().String()
	submitTime := time.Now()
	var submissionImages []SubmissionImage
	for count, image := range Images {
		submissionImages = append(submissionImages, NewImage(count, image, id))
	}
	return &Submission{
		ID:             id,
		SubmissionTime: submitTime,
		Images:         submissionImages,
		Message:        Message,
		TaskID:         TaskID,
	}
}

func (f *Submission) RefreshTimestamp() {
	f.SubmissionTime = time.Now()
}
