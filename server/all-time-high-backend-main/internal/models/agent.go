package models

import (
	"time"

	"gorm.io/gorm"
)

type Agent struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	Name               string         `gorm:"type:varchar(100);not null" json:"name"`
	Ticker             string         `gorm:"type:varchar(50);not null" json:"ticker"`
	Prompt             string         `gorm:"type:text;not null" json:"prompt"`
	Description        string         `gorm:"type:text" json:"description"`
	ImageURL           string         `gorm:"type:varchar(255)" json:"image_url"`
	TokenAddress       string         `gorm:"type:varchar(100)" json:"token_address"` // 新增字段
	UserID             uint           `gorm:"not null;index" json:"user_id"`          // 新增字段
	User               User           `gorm:"foreignKey:UserID" json:"-"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
	MarketCap          float64        `gorm:"type:double precision" json:"market_cap"` // 市值
	MarketCapUpdatedAt time.Time      `json:"market_cap_updated_at"`
	HighestPrice       float64        `gorm:"type:double precision" json:"highest_price"`
}
