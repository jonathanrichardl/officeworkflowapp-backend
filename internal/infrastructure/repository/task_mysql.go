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

func (r *TaskMySQL) Create(t *entity.Task) (string, error) {
	stmt, err := r.db.Prepare(`
		INSERT INTO tasks (ID, user_id, requirement_id, note, fulfillment_status, allowed, deadline, num_of_prerequisite) 
		values(?,?,?,?,?,?,?,?)`)

	if err != nil {
		return t.ID, err
	}

	_, err = stmt.Exec(
		t.ID,
		t.UserID,
		t.RequirementID,
		t.Note,
		t.Status,
		t.Allowed,
		t.Deadline,
		t.NumOfPrerequisite,
	)
	if err != nil {
		return t.ID, err
	}
	return t.ID, nil

}

func (r *TaskMySQL) RemovePrerequisite(taskID string) ([]*entity.Task, error) {
	stmt, err := r.db.Prepare(`SELECT tasks.id, tasks.allowed, tasks.user_id, tasks.fulfillment_status, tasks.num_of_prerequisite, tasks.deadline
								FROM prerequisite INNER JOIN tasks on tasks.id = prerequisite.task_id
								 WHERE prerequisite = ?"`)

	if err != nil {
		return nil, err
	}
	var affectedTasks []*entity.Task
	rows, err := stmt.Query(taskID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t entity.Task
		err = rows.Scan(&t.ID, &t.Allowed, &t.UserID, &t.Status, &t.NumOfPrerequisite, &t.Deadline)
		if err != nil {
			return nil, err
		}
		affectedTasks = append(affectedTasks, &t)
	}
	_, err = r.db.Exec("DELETE FROM prerequisite WHERE prerequisite=?", taskID)
	if err != nil {
		return nil, err
	}
	return affectedTasks, nil
}

func (r *TaskMySQL) Get(id string) (*entity.Task, error) {
	stmt, err := r.db.Prepare(`SELECT id, requirement_id, allowed, user_id, fulfillment_status, num_of_prerequisite, deadline 
								from tasks where id = ?`)
	var task entity.Task
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(id)
	if row == nil {
		return nil, err
	}
	err = row.Scan(&task.ID, &task.RequirementID, &task.Allowed, &task.UserID,
		&task.Status, &task.NumOfPrerequisite, &task.Deadline)
	if err != nil {
		return nil, err
	}
	return &task, nil

}

func (r *TaskMySQL) GetbyUserID(userID string) ([]*entity.TaskWithDetails, error) {
	stmt, err := r.db.Prepare(`SELECT tasks.id, requirements.request, requirements.expected_outcome,  
								orders.title, orders.description, orders.deadline,tasks.fulfillment_status
								FROM tasks INNER JOIN requirements ON tasks.requirement_id=requirements.id 
								INNER JOIN orders ON requirements.order_id = orders.id 
								where user_id = ? and tasks.allowed = true`)
	if err != nil {
		return nil, err
	}
	var tasks []*entity.TaskWithDetails
	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var t entity.TaskWithDetails
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
	_, err := r.db.Exec(`UPDATE tasks SET user_id = ?, fulfillment_status = ?, deadline = ?, num_of_prerequisite = ?,
						 allowed = ?, where id = ?`,
		e.UserID, e.Status, e.Deadline, e.NumOfPrerequisite, e.Allowed, e.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskMySQL) List() ([]*entity.TaskWithDetails, error) {
	stmt, err := r.db.Prepare(`SELECT tasks.id, users.username, requirements.request, requirements.expected_outcome,  
								orders.title, orders.description, orders.deadline, tasks.fulfillment_status 
								FROM tasks INNER JOIN requirements ON tasks.requirement_id=requirements.id 
								INNER JOIN users on users.id = tasks.user_id
								INNER JOIN orders ON requirements.order_id = orders.id `)
	if err != nil {
		return nil, err
	}
	var tasks []*entity.TaskWithDetails
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var t entity.TaskWithDetails
		err = rows.Scan(&t.ID, &t.Username, &t.Request, &t.ExpectedOutcome, &t.OrderTitle, &t.OrderDescription, &t.OrderDeadline,
			&t.Status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	if len(tasks) == 0 {
		return nil, nil
	}

	return tasks, nil

}

func (r *TaskMySQL) Delete(TaskID string) error {
	_, err := r.db.Exec("DELETE FROM task where requirement_id = ?", TaskID)
	if err != nil {
		return err
	}
	return nil
}
