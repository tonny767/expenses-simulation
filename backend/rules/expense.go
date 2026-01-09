package rules

import (
	c "backend/constants"
	"backend/models"
	"errors"
)

var (
	ErrAmountTooSmall            = errors.New("amount is below minimum expense limit")
	ErrAmountTooLarge            = errors.New("amount exceeds maximum expense limit")
	ErrInvalidAmount             = errors.New("amount must be a positive integer")
	ErrEmptyDesc                 = errors.New("description is required")
	ErrExpenseAlreadyFinalized   = errors.New("expense has already been approved and finalized")
	ErrExpenseNotPendingApproval = errors.New("expense is not pending for approval")
)

func ValidateExpense(amount int64, description string) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	if amount < c.MinExpenseAmount {
		return ErrAmountTooSmall
	}

	if amount > c.MaxExpenseAmount {
		return ErrAmountTooLarge
	}

	if description == "" {
		return ErrEmptyDesc
	}

	return nil
}

func RequiresManagerApproval(amount int64) bool {
	return amount >= c.ApprovalThreshold
}

func CanProceed(exp *models.Expense) error {
	switch exp.Status {
	case c.ExpenseStatusPending:
		return nil

	case c.ExpenseStatusApproved,
		c.ExpenseStatusRejected,
		c.ExpenseStatusAutoApproved:
		return ErrExpenseAlreadyFinalized

	default:
		return ErrExpenseNotPendingApproval
	}
}
