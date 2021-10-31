package entity

type SubmissionImage struct {
	ID           int
	Image        []byte
	SubmissionID string
}

func NewImage(ID int, Image []byte, SubmissionID string) SubmissionImage {
	return SubmissionImage{
		ID:           ID,
		Image:        Image,
		SubmissionID: SubmissionID,
	}

}
