package repository

import (
	"database/sql"

	"order-validation-v2/internal/entity"
)

type RequirementsMySQL struct {
	db *sql.DB
}

func NewRequirementsMySQL(db *sql.DB) *RequirementsMySQL {
	return &RequirementsMySQL{
		db: db,
	}
}

func (r *RequirementsMySQL) Create(e *entity.Requirements) (int, error) {
	stmt, err := r.db.Prepare(`
		INSERT INTO requirements (request, expected_outcome, order_id, status) 
		values(?,?,?,0)`)
	if err != nil {
		return -1, err
	}
	_, err = stmt.Exec(
		e.Request,
		e.ExpectedOutcome,
		e.OrderID,
	)
	if err != nil {
		return -1, err
	}
	err = stmt.Close()
	if err != nil {
		return -1, err
	}
	id := r.db.QueryRow("SELECT last_insert_id();")
	var createdID int
	id.Scan(&createdID)

	return createdID, nil
}

func (r *RequirementsMySQL) Get(ID int) (*entity.Requirements, error) {
	stmt, err := r.db.Prepare(`SELECT * FROM requirements where id = ?`)
	if err != nil {
		return nil, err
	}
	var requirements entity.Requirements
	row := stmt.QueryRow(ID)
	err = row.Scan(&requirements.Id, &requirements.Request, &requirements.ExpectedOutcome,
		&requirements.Status, &requirements.OrderID)
	if err != nil {
		return nil, err
	}
	return &requirements, nil
}

func (r *RequirementsMySQL) Update(e *entity.Requirements) error {
	_, err := r.db.Exec("UPDATE requirements SET request = ?,  expected_outcome = ?, status = ?, where id = ?",
		e.Request, e.ExpectedOutcome, e.Status, e.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *RequirementsMySQL) Search(query string) ([]*entity.Requirements, error) {
	stmt, err := r.db.Prepare(`SELECT id, request, expected_outcome, status FROM requirements WHERE request like ?`)
	if err != nil {
		return nil, err
	}
	var requirements []*entity.Requirements
	rows, err := stmt.Query("%" + query + "%")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var q entity.Requirements
		err = rows.Scan(&q.Id, &q.Request, &q.ExpectedOutcome, &q.Status, &q.OrderID)
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, &q)
	}

	return requirements, nil
}

func (r *RequirementsMySQL) List() ([]*entity.Requirements, error) {
	stmt, err := r.db.Prepare(`SELECT * FROM requirements`)
	if err != nil {
		return nil, err
	}
	var requirements []*entity.Requirements
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var q entity.Requirements
		err = rows.Scan(&q.Id, &q.Request, &q.ExpectedOutcome, &q.Status, &q.OrderID)
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, &q)
	}
	return requirements, nil
}

func (r *RequirementsMySQL) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM requirements where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *RequirementsMySQL) GetByOrderID(orderID string) ([]*entity.Requirements, error) {
	stmt, err := r.db.Prepare("SELECT * FROM requirements where order_id = ?")
	if err != nil {
		return nil, err
	}
	var requirements []*entity.Requirements
	rows, err := stmt.Query(orderID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var q entity.Requirements
		err = rows.Scan(&q.Id, &q.Request, &q.ExpectedOutcome, &q.OrderID, &q.Status)
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, &q)
	}
	return requirements, nil
}
