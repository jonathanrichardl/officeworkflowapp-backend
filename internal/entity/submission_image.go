package entity

type SubmissionImage struct {
	ID           int
	Image        string
	SubmissionID string
}

func NewImage(ID int, Image string, SubmissionID string) SubmissionImage {
	return SubmissionImage{
		ID:           ID,
		Image:        Image,
		SubmissionID: SubmissionID,
	}

}
