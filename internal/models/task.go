package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID        uint   `gorm:"primaryKey"`
	AgentID   string `gorm:"index"` // e.g., "archlinux-Nandu"
	Command   string // The shell command: "whoami", "ls -la", etc.
	Status    string // "pending", "completed", "failed"
	Result    string // The terminal output returned by Hermes
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
