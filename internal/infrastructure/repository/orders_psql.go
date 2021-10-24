package repository

import (
	"database/sql"

	"order-validation-v2/internal/entity"
)

type OrdersPSQL struct {
	db *sql.DB
}

func NewOrdersPSQL(db *sql.DB) *OrdersPSQL {
	return &OrdersPSQL{
		db: db,
	}
}

func (r *OrdersPSQL) Create(e *entity.Orders) (string, error) {
	stmt, err := r.db.Prepare(`
		INSERT INTO orders (id, title, description, deadline) 
		values($1,$2,$3,$4)`)
	if err != nil {
		return e.ID, err
	}
	_, err = stmt.Exec(
		e.ID,
		e.Title,
		e.Description,
		e.Deadline,
	)
	if err != nil {
		return e.ID, err
	}
	err = stmt.Close()
	if err != nil {
		return e.ID, err
	}
	return e.ID, nil
}

func (r *OrdersPSQL) Get(id string) (*entity.Orders, error) {
	stmt, err := r.db.Prepare(`SELECT * FROM orders where id = $1`)
	if err != nil {
		return nil, err
	}
	var b entity.Orders
	row := stmt.QueryRow(id)
	err = row.Scan(&b.ID, &b.Title, &b.Description, &b.Deadline)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *OrdersPSQL) Update(e *entity.Orders) error {
	_, err := r.db.Exec("UPDATE orders SET title = $1, description = $2, deadline = $3 where id = $4",
		e.Title, e.Description, e.Deadline, e.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrdersPSQL) Search(query string) ([]*entity.Orders, error) {
	stmt, err := r.db.Prepare(`SELECT * FROM orders WHERE title like $1`)
	if err != nil {
		return nil, err
	}
	var orders []*entity.Orders
	rows, err := stmt.Query("%" + query + "%")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var o entity.Orders
		err = rows.Scan(&o.ID, &o.Title, &o.Description, &o.Deadline)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}

	return orders, nil
}

//List books
func (r *OrdersPSQL) List() ([]*entity.Orders, error) {
	stmt, err := r.db.Prepare(`SELECT * FROM orders`)
	if err != nil {
		return nil, err
	}
	var orders []*entity.Orders
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var o entity.Orders
		err = rows.Scan(&o.ID, &o.Title, &o.Description, &o.Deadline)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}
	return orders, nil
}

//Delete a book
func (r *OrdersPSQL) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM orders where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrdersPSQL) CustomQuery(query string) (*sql.Rows, error) {
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
