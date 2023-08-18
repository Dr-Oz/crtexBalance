package service

import (
	c "crtexBalance/internal/config"
	"crtexBalance/internal/models"
	"crtexBalance/internal/repository"
	"crtexBalance/mq"
	"database/sql"
)

type Control interface {
	ReplenishmentBalance(replenishment *models.Replenishment) error
	Transfer(money *models.Money) error
	GetBalance(userId int) (*models.User, error)
	PublishReplenishment(replenishment *models.Replenishment) error
}

type Service struct {
	Control
}

func NewService(repos *repository.Repository, conf *c.Config, db *sql.DB, rmq *mq.RabbitMQ) *Service {
	return &Service{
		Control: NewControlService(repos.Control, conf, db, rmq),
	}
}
