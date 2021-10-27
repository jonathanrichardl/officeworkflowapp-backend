package repository

import (
	"database/sql"

	"order-validation-v2/internal/entity"
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
		INSERT INTO requirements (request, expectedoutcome, orderid, status, userid) 
		values($1,$2,$3,'0', $4)`)
	if err != nil {
		return -1, err
	}
	_, err = stmt.Exec(
		e.Request,
		e.ExpectedOutcome,
		e.OrderID,
		e.UserID,
	)
	if err != nil {
		return -1, err
	}
	err = stmt.Close()
	if err != nil {
		return -1, err
	}
	row := r.db.QueryRow("SELECT CURRVAL(pg_get_serial_sequence('requirements','id'));")
	var id int
	row.Scan(&id)
	return id, nil
}

func (r *RequirementsPSQL) Get(ID int) (*entity.Requirements, error) {
	stmt, err := r.db.Prepare(`SELECT * FROM requirements where id = $1`)
	if err != nil {
		return nil, err
	}
	var q entity.Requirements
	row := stmt.QueryRow(ID)
	err = row.Scan(&q.Id, &q.Request, &q.ExpectedOutcome, &q.OrderID, &q.UserID, &q.Status)
	if err != nil {
		return nil, err
	}
	return &q, nil
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
		err = rows.Scan(&q.Id, &q.Request, &q.ExpectedOutcome, &q.OrderID, &q.UserID, &q.Status)
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
		err = rows.Scan(&q.Id, &q.Request, &q.ExpectedOutcome, &q.OrderID, &q.UserID, &q.Status)
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

func (r *RequirementsPSQL) GetByUserID(userID string) ([]*entity.Requirements, error) {
	stmt, err := r.db.Prepare("SELECT * FROM requirements where userid = $1")
	if err != nil {
		return nil, err
	}
	var requirements []*entity.Requirements
	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var q entity.Requirements
		err = rows.Scan(&q.Id, &q.Request, &q.ExpectedOutcome, &q.OrderID, &q.UserID, &q.Status)
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, &q)
	}
	if len(requirements) == 0 {
		return nil, nil
	}
	return requirements, nil
}

func (r *RequirementsPSQL) GetByOrderID(orderID string) ([]*entity.Requirements, error) {
	stmt, err := r.db.Prepare("SELECT * FROM requirements where orderid = $1")
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
		err = rows.Scan(&q.Id, &q.Request, &q.ExpectedOutcome, &q.UserID, &q.OrderID, &q.Status)
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, &q)
	}

	if len(requirements) == 0 {
		return nil, nil
	}
	return requirements, nil
}
