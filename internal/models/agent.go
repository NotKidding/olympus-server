package models

import (
	"time"

	"gorm.io/gorm"
)

type Agent struct {
	ID       string `gorm:"primaryKey"`
	Hostname string
	Username string
	IP       string
	// NEW FIELDS
	OSVersion string
	Arch      string

	LastSeen  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
