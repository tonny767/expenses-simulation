package workers

import (
	"context"
	"fmt"
	"log"
	"time"

	"backend/actions"
	"backend/constants"
	"backend/db"
	"backend/models"
	"backend/services"
)

type PaymentWorker struct {
	paymentService services.PaymentService
}

func NewPaymentWorker() *PaymentWorker {
	return &PaymentWorker{
		paymentService: services.NewPaymentService(),
	}
}

const (
	MaxRetries = 3
	RetryDelay = 5 * time.Second
)

// ProcessExpensePaymentAsync processes payment in background
func (w *PaymentWorker) ProcessExpensePaymentAsync(expenseID int64) {
	go func() {
		// Simulate processing delay
		time.Sleep(3 * time.Second)

		ctx := context.Background()

		log.Printf("Starting payment processing for expense %d", expenseID)

		tx := db.DB.Begin()

		// Retry logic, 3 times max, can be improved by making an api to retry failed payments manually
		for attempt := 1; attempt <= MaxRetries; attempt++ {
			// Fetch fresh expense data
			var expense models.Expense
			if err := db.DB.Preload("Approval").First(&expense, expenseID).Error; err != nil {
				log.Printf("Failed to fetch expense %d: %v", expenseID, err)
				return
			}

			if expense.Status == constants.ExpenseStatusCompleted {
				log.Printf("Expense %d is completed", expenseID)
				return
			}

			// part of log
			originalStatus := expense.Status

			updatedExpense, updatedApproval, err := w.paymentService.ProcessPayment(ctx, &expense, expense.Approval)

			if err != nil {
				if attempt >= MaxRetries {
					if expense.Approval != nil {
						failureNote := fmt.Sprintf("Payment failed after %d attempts", MaxRetries)
						if expense.Approval.Notes != "" {
							expense.Approval.Notes += " - " + failureNote
						} else {
							expense.Approval.Notes = failureNote
						}
						if err := tx.Save(expense.Approval).Error; err != nil {
							tx.Rollback()
							log.Printf("Failed to update approval %d: %v", expense.Approval.ID, err)
							return
						}
					}
					return
				}

				retryWait := RetryDelay * time.Duration(attempt)
				time.Sleep(retryWait)
				continue
			}

			audit, err := actions.ExpenseAuditLog(actions.ExpenseAuditLogInput{
				ExpenseID:  expense.ID,
				ActorID:    nil,
				FromStatus: originalStatus,
				ToStatus:   updatedExpense.Status,
				Reason:     "Expense completed by payment processing",
			})
			if err != nil {
				log.Printf("Failed to create audit log for expense %d: %v", expense.ID, err)
				return
			}

			if err := tx.Save(updatedExpense).Error; err != nil {
				tx.Rollback()
				fmt.Printf("Failed to update expense %d: %v", expense.ID, err)
				return
			}

			if updatedApproval != nil {
				if err := tx.Save(updatedApproval).Error; err != nil {
					tx.Rollback()
					fmt.Printf("Failed to update approval for expense %d: %v", expense.ID, err)
					return
				}
			}

			if err := tx.Create(audit).Error; err != nil {
				tx.Rollback()
				fmt.Printf("Failed to create audit log for expense %d: %v", expense.ID, err)
				return
			}

			tx.Commit()
		}
	}()
}
