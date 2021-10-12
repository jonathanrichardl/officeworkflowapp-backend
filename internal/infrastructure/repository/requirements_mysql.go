package repository

import (
	"database/sql"

	"clean/internal/entity"
)

type RequirementsMySQL struct {
	db *sql.DB
}

//NewBookMySQL create new repository
func NewRequirementsMySQL(db *sql.DB) *RequirementsMySQL {
	return &RequirementsMySQL{
		db: db,
	}
}

//Create a book
func (r *RequirementsMySQL) Create(e *entity.Requirements) (int, error) {
	stmt, err := r.db.Prepare(`
		INSERT INTO requirements (request, outcome, order_id, status) 
		values(?,?,?,'0')`)
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
	rows, _ := r.db.Query("SELECT last_insert_id();")
	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return -1, err
		}
	}
	return id, nil
}

func (r *RequirementsMySQL) Get(orderID entity.ID) (*[]entity.Requirements, error) {
	stmt, err := r.db.Prepare(`SELECT id, request,outcome, status FROM requirements where order_id = ?`)
	if err != nil {
		return nil, err
	}
	var requirements []entity.Requirements
	rows, err := stmt.Query(orderID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var q entity.Requirements
		err = rows.Scan(&q.Id, &q.Request, &q.ExpectedOutcome, &q.Status)
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, q)
	}
	return &requirements, nil
}

//Update a book
func (r *RequirementsMySQL) Update(e *entity.Requirements) error {
	_, err := r.db.Exec("UPDATE requirements SET request = ?, order_id = ?, outcome = ?, status = ? where id = ?",
		e.Request, e.OrderID, e.ExpectedOutcome, e.Status, e.Id)
	if err != nil {
		return err
	}
	return nil
}

//Search books
func (r *RequirementsMySQL) Search(query string) ([]*entity.Requirements, error) {
	stmt, err := r.db.Prepare(`SELECT id, request,outcome, status FROM requirements WHERE request like ?`)
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
		err = rows.Scan(&q.Id, &q.Request, &q.ExpectedOutcome, &q.Status)
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, &q)
	}

	return requirements, nil
}

//List books
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

//Delete a book
func (r *RequirementsMySQL) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM requirements where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *RequirementsMySQL) CustomQuery(query string) (*sql.Rows, error) {
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
