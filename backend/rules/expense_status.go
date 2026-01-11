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

func CanTransition(fromStatus c.ExpenseStatus, toStatus c.ExpenseStatus) error {
	validTransitions := map[c.ExpenseStatus][]c.ExpenseStatus{
		c.ExpenseStatusPending:   {c.ExpenseStatusApproved, c.ExpenseStatusRejected},
		c.ExpenseStatusApproved:  {c.ExpenseStatusCompleted},
		c.ExpenseStatusRejected:  {},
		c.ExpenseStatusCompleted: {},
	}

	allowedStatuses, exists := validTransitions[fromStatus]
	if !exists {
		return ErrInvalidStatusTransition
	}

	for _, status := range allowedStatuses {
		if status == toStatus {
			return nil
		}
	}

	return ErrInvalidStatusTransition
}
