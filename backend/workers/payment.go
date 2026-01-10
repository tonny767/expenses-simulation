package workers

import (
	"context"
	"fmt"
	"log"
	"time"

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
						db.DB.Save(expense.Approval)
					}
					return
				}

				retryWait := RetryDelay * time.Duration(attempt)
				time.Sleep(retryWait)
				continue
			}

			if err := db.DB.Save(updatedExpense).Error; err != nil {
				fmt.Printf("Failed to update expense %d: %v", expense.ID, err)
				return
			}

			if updatedApproval != nil {
				if err := db.DB.Save(updatedApproval).Error; err != nil {
					fmt.Printf("Failed to update approval for expense %d: %v", expense.ID, err)
					return
				}
			}
		}
	}()
}
