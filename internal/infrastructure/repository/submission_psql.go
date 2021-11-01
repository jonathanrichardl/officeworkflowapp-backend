package repository

import (
	"database/sql"
	"fmt"
	"order-validation-v2/internal/entity"
)

type SubmissionPSQL struct {
	db *sql.DB
}

func NewSubmissionPSQL(db *sql.DB) *SubmissionPSQL {
	return &SubmissionPSQL{
		db: db,
	}
}

func (r *SubmissionPSQL) Create(e *entity.Submission) (string, error) {
	statement, err := r.db.Prepare(`
		INSERT INTO submissions (id, submit_time, message, task_id) 
		values($1,$2,$3,$4)`)

	if err != nil {
		return e.ID, err
	}

	imageStatement, err := r.db.Prepare(`
		INSERT INTO image_submissions (id, submission_id, image) 
		values($1,$2,$3)`)

	if err != nil {
		return e.ID, err
	}
	_, err = statement.Exec(
		e.ID,
		e.SubmissionTime,
		e.Message,
		e.TaskID,
	)
	fmt.Println("OK")
	if err != nil {
		return e.ID, err
	}
	err = statement.Close()
	if err != nil {
		return e.ID, err
	}

	for _, image := range e.Images {
		_, err = imageStatement.Exec(
			image.ID,
			e.ID,
			image.Image,
		)
		if err != nil {
			return e.ID, err
		}

	}
	err = imageStatement.Close()
	if err != nil {
		return e.ID, err
	}
	return e.ID, nil
}

func (r *SubmissionPSQL) Get(id string) (*entity.Submission, error) {
	statement, err := r.db.Prepare(`SELECT id, submit_time, message FROM submissions where id = $1`)
	if err != nil {
		return nil, err
	}
	imageStatement, err := r.db.Prepare(`SELECT id, image FROM image_submissions where submission_id = $1`)
	if err != nil {
		return nil, err
	}
	var s entity.Submission
	submission := statement.QueryRow(id)
	if err != nil {
		return nil, err
	}

	err = submission.Scan(&s.ID, &s.SubmissionTime, &s.Message)
	if err != nil {
		return nil, err
	}

	images, err := imageStatement.Query(id)
	if err != nil {
		return nil, err
	}

	for images.Next() {
		var i entity.SubmissionImage
		err = images.Scan(&i.ID, &i.Image)
		if err != nil {
			return nil, err
		}
		s.Images = append(s.Images, i)

	}
	return &s, nil
}

func (r *SubmissionPSQL) Update(e *entity.Submission) error {
	_, err := r.db.Exec("UPDATE submissions SET message = $1, submit_time =$2 where id = $3", e.Message, e.SubmissionTime, e.ID)
	if err != nil {
		return err
	}
	for _, image := range e.Images {
		_, err := r.db.Exec("UPDATE image_submission SET IMAGE = $1 where id = $2 ", image.Image, image.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SubmissionPSQL) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM submissions where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
