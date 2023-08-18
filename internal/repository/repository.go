package repository

import (
	"crtexBalance/internal/models"
	"database/sql"
	"time"
)

type Repository struct {
	Control
}

type Control interface {
	UpdateBalanceTx(tx *sql.Tx, userId int, amount int) error
	GetUser(userId int) (*models.User, error)
	GetUserForUpdate(tx *sql.Tx, userId int) (*models.User, error)
	InsertUserTx(tx *sql.Tx, userId int, amount int) error
	InsertLogTx(tx *sql.Tx, userId int, date time.Time, amount int, description string) error
	GetService(serviceId int) (string, error)
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Control: NewControlPostgres(db),
	}
}
