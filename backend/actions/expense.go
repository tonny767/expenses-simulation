package actions

import (
	"backend/constants"
	"backend/models"
	"backend/rules"
	"time"

	"github.com/google/uuid"
)

type SubmitExpenseInput struct {
	UserID      int64  `json:"user_id"`
	AmountIDR   int64  `json:"amount_idr"`
	Description string `json:"description"`
	ReceiptURL  string `json:"receipt_url"`
}

func SubmitExpense(input SubmitExpenseInput) (*models.Expense, *models.Approval, error) {
	if err := rules.ValidateExpense(input.AmountIDR, input.Description); err != nil {
		return nil, nil, err
	}

	now := time.Now().UTC()
	requiresApproval := rules.RequiresManagerApproval(input.AmountIDR)

	expense := &models.Expense{
		UUID:             uuid.New(),
		UserID:           input.UserID,
		AmountIDR:        input.AmountIDR,
		Description:      input.Description,
		ReceiptURL:       input.ReceiptURL,
		SubmittedAt:      now,
		Status:           rules.InitialExpenseStatus(requiresApproval),
		RequiresApproval: requiresApproval,
		AutoApproved:     !requiresApproval,
	}

	approval := &models.Approval{
		ApproverID: nil,
		Notes:      "",
		Status:     constants.ApprovalStatusPending,
		CreatedAt:  now,
	}

	return expense, approval, nil
}

type ApproveExpenseInput struct {
	Expense    *models.Expense
	ApproverID *int64
	Notes      string
}

func ApproveExpense(input ApproveExpenseInput) (*models.Expense, *models.Approval, error) {
	if err := rules.CanApproveExpense(input.Expense.Status); err != nil {
		return nil, nil, err
	}

	if err := rules.CanProceed(input.Expense); err != nil {
		return nil, nil, err
	}

	toStatus := constants.ExpenseStatusApproved // manually for now
	if err := rules.CanTransition(input.Expense.Status, toStatus); err != nil {
		return nil, nil, err
	}

	input.Expense.Status = toStatus

	input.Expense.Approval.Status = constants.ApprovalStatusApproved
	input.Expense.Approval.ApproverID = input.ApproverID
	input.Expense.Approval.Notes = input.Notes

	return input.Expense, input.Expense.Approval, nil
}

type RejectExpenseInput struct {
	Expense    *models.Expense
	ApproverID *int64
	Notes      string
}

func RejectExpense(input RejectExpenseInput) (*models.Expense, *models.Approval, error) {
	if err := rules.CanRejectExpense(input.Expense.Status); err != nil {
		return nil, nil, err
	}
	if err := rules.CanProceed(input.Expense); err != nil {
		return nil, nil, err
	}

	toStatus := constants.ExpenseStatusRejected // manually for now
	if err := rules.CanTransition(input.Expense.Status, toStatus); err != nil {
		return nil, nil, err
	}

	input.Expense.Status = toStatus

	input.Expense.Approval.Status = constants.ApprovalStatusRejected
	input.Expense.Approval.ApproverID = input.ApproverID
	input.Expense.Approval.Notes = input.Notes

	return input.Expense, input.Expense.Approval, nil
}
