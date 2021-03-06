package repository

import (
	"database/sql"

	"order-validation-v2/internal/entity"
)

type OrdersMySQL struct {
	db *sql.DB
}

func NewOrdersMySQL(db *sql.DB) *OrdersMySQL {
	return &OrdersMySQL{
		db: db,
	}
}

func (r *OrdersMySQL) Create(e *entity.Orders) (string, error) {
	stmt, err := r.db.Prepare(`
		INSERT INTO orders (id, title, description, deadline) 
		values(?,?,?,?)`)
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

func (r *OrdersMySQL) Get(id string) (*entity.Orders, error) {
	stmt, err := r.db.Prepare(`SELECT * FROM orders where id = ?`)
	if err != nil {
		return nil, err
	}
	var o entity.Orders
	row := stmt.QueryRow(id)
	if err != nil {
		return nil, err
	}

	err = row.Scan(&o.ID, &o.Title, &o.Description, &o.Deadline)
	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (r *OrdersMySQL) Update(e *entity.Orders) error {
	_, err := r.db.Exec("UPDATE orders SET title = ?, description = ?, deadline = ? where id = '?'",
		e.Title, e.Description, e.Deadline, e.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrdersMySQL) Search(query string) ([]*entity.Orders, error) {
	stmt, err := r.db.Prepare(`SELECT * FROM orders WHERE title like ?`)
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

func (r *OrdersMySQL) List() ([]*entity.Orders, error) {
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

func (r *OrdersMySQL) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM orders where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrdersMySQL) CustomQuery(query string) (*sql.Rows, error) {
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
