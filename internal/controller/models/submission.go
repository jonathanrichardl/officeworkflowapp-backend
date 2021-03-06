package models

import "order-validation-v2/internal/entity"

type Submission struct {
	ID             string  `json:"submission_id"`
	SubmissionTime string  `json:"submission_time"`
	TaskID         string  `json:"task_id,omitempty"`
	Images         []Image `json:"images"`
	Message        string  `json:"message"`
}

type Image struct {
	ID    int    `json:"image_id,omitempty"`
	Image string `json:"image"`
}

type SubmissionUpdate struct {
	ID      string
	Images  []Image
	Message string
}

func BuildSubmissionPayload(submissions []*entity.Submission) []Submission {
	var submissionJSON []Submission
	for _, s := range submissions {
		var images []Image
		for _, imageData := range s.Images {
			image := Image{
				ID:    imageData.ID,
				Image: imageData.Image,
			}
			images = append(images, image)

		}
		submission := Submission{
			ID:             s.ID,
			SubmissionTime: s.SubmissionTime.Format("02-Jan-2006 15:04:05"),
			TaskID:         s.TaskID,
			Images:         images,
			Message:        s.Message,
		}
		submissionJSON = append(submissionJSON, submission)
	}
	return submissionJSON

}

func DecodeSubmissionPayload(submission Submission) []string {
	var imageEntity []string
	for _, imageData := range submission.Images {
		imageEntity = append(imageEntity, imageData.Image)
	}
	return imageEntity

}
