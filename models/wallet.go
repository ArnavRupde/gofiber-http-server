package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Wallet struct {
	ID         int `json:"id" gorm:"primaryKey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	WalletUuid string          `json:"wallet_uuid"`                        // Define field as UUID and create an unique index
	Balance    decimal.Decimal `json:"balance" gorm:"type:decimal(18,2);"` // Define field as decimal with precision 18 and scale 2
	UserId     int
	User       User `gorm:"foreignKey:UserId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Define foreign key relationship
}

// BeforeCreate is a GORM hook that is triggered before a record is created
func (wallet *Wallet) BeforeCreate(tx *gorm.DB) error {
	wallet.WalletUuid = uuid.NewString()
	return nil
}
