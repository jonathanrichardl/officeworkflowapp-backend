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
		INSERT INTO requirements (request, expectedoutcome, order_id, status) 
		values($1,$2,$3,'0')`)
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
	// rows, _ := r.db.Query("SELECT last_insert_id();")
	rows, _ := r.db.Query("SELECT CURRVAL(pg_get_serial_sequence('requirements','id'));")
	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return -1, err
		}
	}
	return id, nil
}

func (r *RequirementsMySQL) Get(ID int) (*entity.Requirements, error) {
	stmt, err := r.db.Prepare(`SELECT id, request, expectedoutcome, status FROM requirements where id = $1`)
	if err != nil {
		return nil, err
	}
	var requirements entity.Requirements
	rows, err := stmt.Query(ID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&requirements.Id, &requirements.Request, &requirements.ExpectedOutcome, &requirements.Status)
		if err != nil {
			return nil, err
		}
	}
	return &requirements, nil
}

//Update a book
func (r *RequirementsMySQL) Update(e *entity.Requirements) error {
	_, err := r.db.Exec("UPDATE requirements SET request = $1, order_id = $2, expectedoutcome = $3, status = $4 where id = $5",
		e.Request, e.OrderID, e.ExpectedOutcome, e.Status, e.Id)
	if err != nil {
		return err
	}
	return nil
}

//Search books
func (r *RequirementsMySQL) Search(query string) ([]*entity.Requirements, error) {
	stmt, err := r.db.Prepare(`SELECT id, request, expectedoutcome, status FROM requirements WHERE request like '$1'`)
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
	_, err := r.db.Exec("DELETE FROM requirements where id = $1", id)
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
