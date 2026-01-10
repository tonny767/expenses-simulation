package actions

import (
	"backend/actions"
	"backend/constants"
	"backend/models"
	"backend/services"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ---------- Helpers ----------
func setupDB(t *testing.T) *gorm.DB {
	// using sqlite in-memory for tests
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	if err := db.AutoMigrate(&models.Expense{}, &models.Approval{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func ptrInt64(v int64) *int64 { return &v }

func TestSubmitExpense_AutoApproved(t *testing.T) {
	db := setupDB(t)

	input := actions.SubmitExpenseInput{
		UserID:      1,
		AmountIDR:   constants.MinExpenseAmount, //10k -- auto approved
		Description: "Auto-approved expense",
		ReceiptURL:  "https://via.placeholder.com",
	}

	expense, approval, err := actions.SubmitExpense(input)

	db.Create(&expense)

	approval.ExpenseID = expense.ID

	db.Create(&approval)

	// mock payment service
	ctx := context.Background()
	paymentService := services.NewPaymentServiceWithBaseURL("https://1620e98f-7759-431c-a2aa-f449d591150b.mock.pstmn.io")
	updatedExpense, updatedApproval, err := paymentService.ProcessPayment(ctx, expense, approval)

	// Wait a bit for goroutine
	time.Sleep(4 * time.Second)

	assert.NoError(t, err)
	assert.NotNil(t, updatedExpense)
	assert.NotNil(t, updatedApproval)
	assert.True(t, updatedExpense.AutoApproved)
	assert.False(t, updatedExpense.RequiresApproval)
	assert.Equal(t, constants.ExpenseStatusCompleted, updatedExpense.Status)
	fmt.Println("Test for auto-approved succeeded")

}

func TestSubmitExpense_PendingApproval(t *testing.T) {
	setupDB(t)

	input := actions.SubmitExpenseInput{
		UserID:      2,
		AmountIDR:   constants.ApprovalThreshold + 10000, // above threshold
		Description: "Pending approval expense",
		ReceiptURL:  "https://via.placeholder.com",
	}

	expense, approval, err := actions.SubmitExpense(input)
	assert.NoError(t, err)
	assert.NotNil(t, expense)
	assert.NotNil(t, approval)
	assert.False(t, expense.AutoApproved)
	assert.True(t, expense.RequiresApproval)
	assert.Equal(t, constants.ExpenseStatusPending, expense.Status)
	fmt.Println("Test for pending approval succeeded")
}

func TestSubmitExpense_InvalidAmount(t *testing.T) {
	setupDB(t)

	input := actions.SubmitExpenseInput{
		UserID:      3,
		AmountIDR:   5000, // invalid amount
		Description: "Invalid expense",
		ReceiptURL:  "https://via.placeholder.com",
	}

	expense, approval, err := actions.SubmitExpense(input)
	assert.Error(t, err)
	assert.Nil(t, expense)
	assert.Nil(t, approval)
	fmt.Println("Test for invalid amount succeeded")
}

func TestApproveExpense_Success(t *testing.T) {
	setupDB(t)

	expense, approval, _ := actions.SubmitExpense(actions.SubmitExpenseInput{
		UserID:      4,
		AmountIDR:   constants.ApprovalThreshold + 10000, // requires approval
		Description: "Approval test",
		ReceiptURL:  "https://via.placeholder.com",
	})
	expense.Approval = approval

	approvedExpense, approvedApproval, err := actions.ApproveExpense(actions.ApproveExpenseInput{
		Expense:    expense,
		ApproverID: ptrInt64(99),
		Notes:      "Approved",
	})

	// mock payment service
	ctx := context.Background()
	paymentService := services.NewPaymentServiceWithBaseURL("https://1620e98f-7759-431c-a2aa-f449d591150b.mock.pstmn.io")
	updatedExpense, updatedApproval, err := paymentService.ProcessPayment(ctx, approvedExpense, approvedApproval)

	assert.NoError(t, err)
	assert.Equal(t, constants.ExpenseStatusCompleted, updatedExpense.Status)
	assert.Equal(t, constants.ApprovalStatusApproved, updatedApproval.Status)
	assert.Equal(t, int64(99), *updatedApproval.ApproverID)
	assert.Equal(t, "Approved", updatedApproval.Notes)
	fmt.Println("Test for approve expense succeeded")
}

func TestRejectExpense_Success(t *testing.T) {
	setupDB(t)

	// expected status should be pending
	expense, approval, _ := actions.SubmitExpense(actions.SubmitExpenseInput{
		UserID:      5,
		AmountIDR:   constants.ApprovalThreshold + 20000, // requires approval
		Description: "Rejection test",
		ReceiptURL:  "https://via.placeholder.com",
	})
	expense.Approval = approval

	rejectedExpense, rejectedApproval, err := actions.RejectExpense(actions.RejectExpenseInput{
		Expense:    expense,
		ApproverID: ptrInt64(101),
		Notes:      "Not allowed",
	})

	assert.NoError(t, err)
	assert.Equal(t, constants.ExpenseStatusRejected, rejectedExpense.Status)
	assert.Equal(t, constants.ApprovalStatusRejected, rejectedApproval.Status)
	assert.Equal(t, int64(101), *rejectedApproval.ApproverID)
	assert.Equal(t, "Not allowed", rejectedApproval.Notes)
	fmt.Println("Test for reject expense succeeded")
}
