// models/expense.go
package models

import (
	"backend/constants"
	"time"

	"github.com/google/uuid"
)

type Expense struct {
	ID               int64                   `json:"id" gorm:"primaryKey"`
	UUID             uuid.UUID               `json:"uuid" gorm:"type:uuid"`
	UserID           int64                   `json:"user_id"`
	AmountIDR        int64                   `json:"amount_idr" gorm:"column:amount_idr"`
	Description      string                  `json:"description"`
	ReceiptURL       string                  `json:"receipt_url"`
	Status           constants.ExpenseStatus `json:"status" gorm:"type:text"`
	RequiresApproval bool                    `json:"requires_approval"`
	AutoApproved     bool                    `json:"auto_approved"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
	SubmittedAt      time.Time               `json:"submitted_at"`
	ProcessedAt      *time.Time              `json:"processed_at"`

	Approval *Approval `json:"approval" gorm:"foreignKey:ExpenseID"`
}

type Approval struct {
	ID         int64                    `json:"id" gorm:"primaryKey"`
	ExpenseID  int64                    `json:"expense_id"`
	ApproverID *int64                   `json:"approver_id"`
	Status     constants.ApprovalStatus `json:"status" gorm:"type:text"`
	Notes      string                   `json:"notes"`
	CreatedAt  time.Time                `json:"created_at"`
	UpdatedAt  time.Time                `json:"updated_at"`

	Expense Expense `json:"-" gorm:"foreignKey:ExpenseID"`
}
