package models

import (
	"backend/constants"
	"time"
)

type User struct {
	ID           int64              `json:"id" gorm:"primaryKey"`
	Email        string             `json:"email" gorm:"uniqueIndex"`
	Name         string             `json:"name"`
	Role         constants.UserRole `json:"role"` // "user" or "manager"
	PasswordHash string             `json:"-"`    // Never send to frontend
	CreatedAt    time.Time          `json:"created_at"`

	Expenses []Expense `json:"expenses,omitempty" gorm:"foreignKey:UserID"`
}
