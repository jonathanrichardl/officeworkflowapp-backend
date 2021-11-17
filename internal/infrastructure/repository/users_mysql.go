package repository

import (
	"database/sql"

	"order-validation-v2/internal/entity"
)

type UserMySQL struct {
	db *sql.DB
}

//NewBookMySQL create new repository
func NewUserMySQL(db *sql.DB) *UserMySQL {
	return &UserMySQL{
		db: db,
	}
}

func (r *UserMySQL) Create(u *entity.User) (string, error) {

	stmt, err := r.db.Prepare(`
		INSERT INTO users (id, username, email, pswd, user_role) 
		values(?, ?,?, sha2(?,256), ?)`)
	if err != nil {
		return u.ID, err
	}
	_, err = stmt.Exec(
		u.ID,
		u.Username,
		u.Email,
		u.Password,
		u.UserRole,
	)
	if err != nil {
		return u.ID, err
	}
	err = stmt.Close()
	if err != nil {
		return u.ID, err
	}
	return u.ID, nil
}

func (r *UserMySQL) GetbyUsername(username string) (*entity.User, error) {
	stmt, err := r.db.Prepare(`SELECT id, username, email, pswd, user_role from users where username = ?`)
	if err != nil {
		return nil, err
	}
	var user entity.User
	row := stmt.QueryRow(username)
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.UserRole)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserMySQL) GetbyID(ID string) (*entity.User, error) {
	stmt, err := r.db.Prepare(`SELECT id, username, email, user_role from users where ID = ?`)
	if err != nil {
		return nil, err
	}
	var user entity.User
	row := stmt.QueryRow(ID)
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.UserRole)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserMySQL) Update(u *entity.User) error {
	_, err := r.db.Exec("UPDATE users SET pswd = sha2(?,256),  username = ?, email = ? , user_role = ? where id = ?",
		u.Password, u.Username, u.Email, u.UserRole, u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserMySQL) Search(query string) ([]*entity.User, error) {
	stmt, err := r.db.Prepare(`SELECT id, username, email , user_role FROM users WHERE username like ?`)
	if err != nil {
		return nil, err
	}
	var users []*entity.User
	rows, err := stmt.Query("%" + query + "%")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var u entity.User
		err = rows.Scan(&u.ID, &u.Username, &u.Email, &u.UserRole)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

func (r *UserMySQL) List() ([]*entity.User, error) {
	stmt, err := r.db.Prepare(`SELECT ID, username, email, user_role FROM users`)
	if err != nil {
		return nil, err
	}
	var users []*entity.User
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var u entity.User
		err = rows.Scan(&u.ID,
			&u.Email, &u.Username, &u.UserRole)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

func (r *UserMySQL) Delete(ID string) error {
	_, err := r.db.Exec("DELETE FROM users where id = ?", ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserMySQL) CustomQuery(query string) (*sql.Rows, error) {
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
