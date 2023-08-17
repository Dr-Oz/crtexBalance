package service

import (
	c "crtexBalance/internal/config"
	"crtexBalance/internal/models"
	"crtexBalance/internal/repository"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const layout string = "2006-01-02"

type ControlService struct {
	repo repository.Control
	conf *c.Config
	db   *sql.DB
}

func NewControlService(repo repository.Control, conf *c.Config, db *sql.DB) *ControlService {
	return &ControlService{
		repo: repo,
		conf: conf,
		db:   db,
	}
}

func (c *ControlService) GetBalance(userId int) (*models.User, error) {
	var user *models.User
	var err error

	if user, err = c.repo.GetUser(userId); err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("пользователь не найден")
	}
	return user, err
}

func (c *ControlService) ReplenishmentBalance(replenishment *models.Replenishment) error {
	var tx *sql.Tx
	var err error
	var user *models.User

	date, _ := time.Parse(layout, replenishment.Date)
	if date.IsZero() {
		date = time.Now()
	}

	tx, err = c.db.Begin()
	if err != nil {
		return err
	}

	user, err = c.repo.GetUserForUpdate(tx, replenishment.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if user != nil {
		if err = c.repo.UpdateBalanceTx(tx, replenishment.UserID, user.Balance+replenishment.Amount); err != nil {
			tx.Rollback()
			return err
		}
		if err = c.repo.InsertLogTx(tx, replenishment.UserID, date, replenishment.Amount, "Пополнение баланса"); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if err = c.repo.InsertUserTx(tx, replenishment.UserID, replenishment.Amount); err != nil {
			tx.Rollback()
			return err
		}
		if err = c.repo.InsertLogTx(tx, replenishment.UserID, date, replenishment.Amount, "Пополнение баланса"); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (c *ControlService) Transfer(money *models.Money) error {
	var tx *sql.Tx
	var err error
	var fromUser, toUser *models.User

	date, _ := time.Parse(layout, money.Date)
	if date.IsZero() {
		date = time.Now()
	}

	tx, err = c.db.Begin()
	if err != nil {
		return err
	}

	if fromUser, err = c.repo.GetUserForUpdate(tx, money.FromUserID); err != nil {
		tx.Rollback()
		return err
	}
	if fromUser == nil {
		tx.Rollback()
		return errors.New("пользователь не найден")
	}
	if toUser, err = c.repo.GetUserForUpdate(tx, money.ToUserID); err != nil {
		tx.Rollback()
		return err
	}
	if toUser == nil {
		tx.Rollback()
		return errors.New("пользователь не найден")
	}

	if fromUser.Balance-money.Amount < 0 {
		tx.Rollback()
		return errors.New("недостаточно средств")
	}

	if err = c.repo.UpdateBalanceTx(tx, fromUser.Id, fromUser.Balance-money.Amount); err != nil {
		tx.Rollback()
		return err
	}
	if err = c.repo.InsertLogTx(tx, money.FromUserID, date, money.Amount, fmt.Sprintf("Перевод средств пользователю %d", money.ToUserID)); err != nil {
		tx.Rollback()
		return err
	}

	if err = c.repo.UpdateBalanceTx(tx, toUser.Id, toUser.Balance+money.Amount); err != nil {
		tx.Rollback()
		return err
	}
	if err = c.repo.InsertLogTx(tx, money.ToUserID, date, money.Amount, fmt.Sprintf("Перевод средств от пользователя %d", money.FromUserID)); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
