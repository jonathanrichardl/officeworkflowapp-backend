package entity

type Status int8

const (
	Unfinished Status = iota
	InReview
	Finished
)

type Task struct {
	ID               string
	RequirementID    int
	Request          string
	ExpectedOutcome  string
	UserID           string
	Status           Status
	OrderTitle       string
	OrderDescription string
	OrderDeadline    string
	SubmissionID     string
	SubmissionTime   string
}

func NewTask(requirementID int, userID string) *Task {
	return &Task{
		RequirementID: requirementID,
		UserID:        userID,
		Status:        0,
	}
}

func (t *Task) Submit(submissionID string) {
	t.SubmissionID = submissionID
}
