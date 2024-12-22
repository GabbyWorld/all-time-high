// internal/models/battle.go
package models

import "time"

type Battle struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	AttackerID  uint      `gorm:"not null;index" json:"attacker_id"`
	Attacker    Agent     `gorm:"foreignKey:AttackerID" json:"-"`
	DefenderID  uint      `gorm:"not null;index" json:"defender_id"`
	Defender    Agent     `gorm:"foreignKey:DefenderID" json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	Outcome     string    `gorm:"type:varchar(20);not null" json:"outcome"`
	Description string    `json:"description"`
}
