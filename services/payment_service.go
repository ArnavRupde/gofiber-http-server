package services

import (
	"context"
	"errors"
	"fmt"

	"example.com/server/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TransferDetails struct {
	FromUserId int             `json:"from_user_id"`
	ToUserId   int             `json:"to_user_id"`
	Amount     decimal.Decimal `json:"amount"`
}

type PaymentService struct {
	db *gorm.DB
}

func NewPaymentService(db *gorm.DB) *PaymentService {
	return &PaymentService{db: db}
}

func (s *PaymentService) TransferFunds(
	ctx context.Context,
	fromUserId int,
	toUserId int,
	amount decimal.Decimal,
) error {
	tx := s.db.Begin() // Start a new transaction

	fromUser := &models.User{Id: fromUserId}
	toUser := &models.User{Id: toUserId}

	err := tx.Where("id = ?", fromUserId).First(fromUser).Error
	if err != nil {
		tx.Rollback()
		return errors.New(
			fmt.Sprintf(
				"Error while fetching user with id: %d , Error: %v",
				fromUserId,
				err,
			),
		)
	}

	err = tx.Where("id = ?", toUserId).First(toUser).Error
	if err != nil {
		tx.Rollback()
		return errors.New(
			fmt.Sprintf(
				"Error while fetching user with id: %d , Error: %v",
				toUserId,
				err,
			),
		)
	}

	err = tx.Model(&models.Wallet{}).Where("user_id = ?", fromUserId).Update("balance", gorm.Expr("balance - ?", amount)).Error
	if err != nil {
		tx.Rollback()
		return errors.New(
			fmt.Sprintf(
				"Error while deducting funds from user with id: %d , Error: %v",
				fromUserId,
				err,
			),
		)
	}

	err = tx.Model(&models.Wallet{}).Where("user_id = ?", toUserId).Update("balance", gorm.Expr("balance + ?", amount)).Error
	if err != nil {
		tx.Rollback()
		return errors.New(
			fmt.Sprintf(
				"Error while adding funds to user with id: %d , Error: %v",
				toUserId,
				err,
			),
		)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return errors.New(
			fmt.Sprintf(
				"Error while committing transaction, Error: %v",
				err,
			),
		)
	}

	return nil
}
