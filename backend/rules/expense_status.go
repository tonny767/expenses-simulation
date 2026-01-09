package rules

import (
	c "backend/constants"
	"errors"
)

var (
	ErrInvalidStatusTransition = errors.New("invalid expense status transition")
)

func InitialExpenseStatus(requiresApproval bool) c.ExpenseStatus {
	if requiresApproval {
		return c.ExpenseStatusPending
	}
	return c.ExpenseStatusApproved
}

func CanApproveExpense(status c.ExpenseStatus) error {
	if status != c.ExpenseStatusPending {
		return ErrInvalidStatusTransition
	}
	return nil
}

func CanRejectExpense(status c.ExpenseStatus) error {
	if status != c.ExpenseStatusPending {
		return ErrInvalidStatusTransition
	}
	return nil
}
