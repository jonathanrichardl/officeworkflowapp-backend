package repository

import (
	"database/sql"
	"order-validation-v2/internal/entity"
)

type TaskMySQL struct {
	db *sql.DB
}

func NewTaskMySQL(db *sql.DB) *TaskMySQL {
	return &TaskMySQL{
		db: db,
	}
}

func (r *TaskMySQL) Create(e *entity.Task) (string, error) {
	stmt, err := r.db.Prepare(`
		INSERT INTO tasks (ID, user_id, requirement_id, status, submission_id) 
		values($1,$2,$3,$4)`)

	if err != nil {
		return e.ID, err
	}
	_, err = stmt.Exec(
		e.ID,
		e.UserID,
		e.RequirementID,
		e.Status,
		e.SubmissionID,
	)
	if err != nil {
		return e.ID, err
	}
	return e.ID, nil

}

func (r *TaskMySQL) Get(userID string) ([]*entity.Task, error) {
	stmt, err := r.db.Prepare(`SELECT tasks.id, requirements.request, requirements.expected_outcome,  
								orders.title, orders.description, orders.deadline,tasks.fulfillment_status
								FROM tasks INNER JOIN requirements ON tasks.requirement_id=requirements.id 
								INNER JOIN orders ON requirements.order_id = orders.id 
								where user_id = $1 `)
	if err != nil {
		return nil, err
	}
	var tasks []*entity.Task
	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var t entity.Task
		err = rows.Scan(&t.ID, &t.Request, &t.ExpectedOutcome, &t.OrderTitle, &t.OrderDescription, &t.OrderDeadline,
			&t.Status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	return tasks, nil

}
func (r *TaskMySQL) Update(e *entity.Task) error {
	_, err := r.db.Exec("UPDATE task SET user_id = $1 where id = $2",
		e.UserID, e.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskMySQL) List() ([]*entity.Task, error) {
	stmt, err := r.db.Prepare(`SELECT tasks.id, requirements.request, requirements.expected_outcome,  
								orders.title, orders.description, orders.deadline,
								tasks.fulfillment_status FROM tasks INNER JOIN requirements ON tasks.requirement_id=requirements.id 
								INNER JOIN orders ON requirements.order_id = orders.id `)
	if err != nil {
		return nil, err
	}
	var tasks []*entity.Task
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var t entity.Task
		err = rows.Scan(&t.ID, &t.Request, &t.ExpectedOutcome, &t.OrderTitle, &t.OrderDescription, &t.OrderDeadline,
			&t.Status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	return tasks, nil

}

func (r *TaskMySQL) Delete(TaskID string) error {
	_, err := r.db.Exec("DELETE FROM task where requirement_id = $1", TaskID)
	if err != nil {
		return err
	}
	return nil
}
