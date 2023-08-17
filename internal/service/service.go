package service

import (
	c "crtexBalance/internal/config"
	"crtexBalance/internal/models"
	"crtexBalance/internal/repository"
	"database/sql"
)

type Control interface {
	ReplenishmentBalance(replenishment *models.Replenishment) error
	Transfer(money *models.Money) error
	GetBalance(userId int) (*models.User, error)
}

type Service struct {
	Control
}

func NewService(repos *repository.Repository, conf *c.Config, db *sql.DB) *Service {
	return &Service{
		Control: NewControlService(repos.Control, conf, db),
	}
}
