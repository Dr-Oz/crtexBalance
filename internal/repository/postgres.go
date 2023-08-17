package repository

import (
	"crtexBalance/internal/models"
	"database/sql"
	"time"
)

type ControlPosgres struct {
	DB *sql.DB
}

func NewControlPostgres(db *sql.DB) *ControlPosgres {
	return &ControlPosgres{DB: db}
}

func (m *ControlPosgres) GetUser(userId int) (*models.User, error) {
	var balance int
	var id int
	rows, err := m.DB.Query("SELECT id, balance FROM users WHERE id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&id, &balance)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return &models.User{Id: id, Balance: balance}, err
}

func (m *ControlPosgres) GetUserForUpdate(tx *sql.Tx, userId int) (*models.User, error) {
	var balance int
	var id int

	stmt, err := tx.Prepare(`SELECT id, balance FROM users WHERE id = $1 FOR UPDATE;`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&id, &balance)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return &models.User{Id: id, Balance: balance}, err
}

func (m *ControlPosgres) UpdateBalanceTx(tx *sql.Tx, userId int, amount int) error {

	stmt, err := tx.Prepare(`UPDATE users SET balance = $1 WHERE id = $2;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(amount, userId); err != nil {
		return err
	}
	return err
}

func (m *ControlPosgres) InsertUserTx(tx *sql.Tx, userId int, amount int) error {

	stmt, err := tx.Prepare(`INSERT INTO users (id, balance) VALUES ($1, $2);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(userId, amount); err != nil {
		return err
	}
	return err
}

func (m *ControlPosgres) InsertLogTx(tx *sql.Tx, userId int, date time.Time, amount int, description string) error {

	stmt, err := tx.Prepare(`INSERT INTO logs (user_id, date, amount, description) VALUES ($1, $2, $3, $4);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(userId, date, amount, description); err != nil {
		return err
	}
	return err
}

func (m *ControlPosgres) GetService(serviceId int) (string, error) {
	var title string

	rows, err := m.DB.Query("SELECT title FROM services WHERE id = $1", serviceId)
	if err != nil {
		return title, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&title)
		if err != nil {
			return title, err
		}
	}

	return title, err
}
