package models

import (
	"time"

	"gorm.io/gorm"
)

type Agent struct {
	ID        string `gorm:"primaryKey"`
	Hostname  string
	Username  string
	IP        string
	LastSeen  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
