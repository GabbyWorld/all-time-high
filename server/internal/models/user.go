// internal/models/user.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	WalletAddress string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"wallet_address"`
	Username      string         `gorm:"type:varchar(50)" json:"username,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
