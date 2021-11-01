package repository

import (
	"database/sql"
	"order-validation-v2/internal/entity"
)

type TaskPSQL struct {
	db *sql.DB
}

func NewTaskPSQL(db *sql.DB) *TaskPSQL {
	return &TaskPSQL{
		db: db,
	}
}

func (r *TaskPSQL) Create(e *entity.Task) (string, error) {
	stmt, err := r.db.Prepare(`
		INSERT INTO tasks (ID, user_id, requirement_id, fulfillment_status) 
		values($1,$2,$3,$4)`)

	if err != nil {
		return e.ID, err
	}
	_, err = stmt.Exec(
		e.ID,
		e.UserID,
		e.RequirementID,
		e.Status,
	)
	if err != nil {
		return e.ID, err
	}
	return e.ID, nil

}

func (r *TaskPSQL) Get(id string) (*entity.Task, error) {
	stmt, err := r.db.Prepare(`SELECT id, requirement_id, user_id, fulfillment_status from tasks where id = $1`)
	var task entity.Task
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(id)
	err = row.Scan(&task.ID, &task.RequirementID, &task.UserID, &task.Status)
	if err != nil {
		return nil, err
	}
	return &task, nil

}

func (r *TaskPSQL) GetbyUserID(userID string) ([]*entity.Task, error) {
	stmt, err := r.db.Prepare(`SELECT tasks.id, requirements.request, requirements.expected_outcome,  
								orders.title, orders.description, orders.deadline,tasks.fulfillment_status
								FROM tasks INNER JOIN requirements ON tasks.requirement_id=requirements.id 
								INNER JOIN orders ON requirements.order_id = orders.id 
								where user_id = $1`)
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
func (r *TaskPSQL) Update(e *entity.Task) error {
	_, err := r.db.Exec("UPDATE tasks SET user_id = $1, fulfillment_status = $2 where id = $3",
		e.UserID, e.Status, e.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskPSQL) List() ([]*entity.Task, error) {
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

func (r *TaskPSQL) Delete(TaskID string) error {
	_, err := r.db.Exec("DELETE FROM task where requirement_id = $1", TaskID)
	if err != nil {
		return err
	}
	return nil
}
