package repository

import (
	"database/sql"

	"clean/internal/entity"
)

type RequirementsPSQL struct {
	db *sql.DB
}

func NewRequirementsPSQL(db *sql.DB) *RequirementsPSQL {
	return &RequirementsPSQL{
		db: db,
	}
}

func (r *RequirementsPSQL) Create(e *entity.Requirements) (int, error) {
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

func (r *RequirementsPSQL) Get(ID int) (*entity.Requirements, error) {
	stmt, err := r.db.Prepare(`SELECT id, request, expectedoutcome, status FROM requirements where id = $1`)
	if err != nil {
		return nil, err
	}
	var requirements entity.Requirements
	row := stmt.QueryRow(ID)
	err = row.Scan(&requirements.Id, &requirements.Request, &requirements.ExpectedOutcome, &requirements.Status)
	if err != nil {
		return nil, err
	}
	return &requirements, nil
}

func (r *RequirementsPSQL) Update(e *entity.Requirements) error {
	_, err := r.db.Exec("UPDATE requirements SET request = $1,  expectedoutcome = $2, status = $3 where id = $4",
		e.Request, e.ExpectedOutcome, e.Status, e.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *RequirementsPSQL) Search(query string) ([]*entity.Requirements, error) {
	stmt, err := r.db.Prepare(`SELECT id, request, expectedoutcome, status FROM requirements WHERE request like $1`)
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

func (r *RequirementsPSQL) List() ([]*entity.Requirements, error) {
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

func (r *RequirementsPSQL) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM requirements where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *RequirementsPSQL) CustomQuery(query string) (*sql.Rows, error) {
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
