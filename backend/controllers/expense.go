package controllers

import (
	"net/http"

	"backend/actions"
	"backend/constants"
	"backend/db"
	"backend/rules"
	"backend/services"

	"github.com/gin-gonic/gin"

	"backend/models"
)

// ---- Get All Expenses ----
func GetExpenses(c *gin.Context) {
	var expenses []models.Expense

	if err := db.DB.Preload("Approval").Order("submitted_at DESC").Find(&expenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}
	c.JSON(http.StatusOK, expenses)
}

func GetUserExpenses(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var expenses []models.Expense
	if err := db.DB.Preload("Approval").Where("user_id = ?", userID).Order("submitted_at DESC").Find(&expenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user expenses"})
		return
	}

	c.JSON(http.StatusOK, expenses)
}

func GetExpense(c *gin.Context) {
	id := c.Param("id")

	var expense models.Expense
	if err := db.DB.Preload("Approval").First(&expense, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	c.JSON(http.StatusOK, expense)
}

func CreateExpense(c *gin.Context) {
	var input actions.SubmitExpenseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expense, approval, err := actions.SubmitExpense(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Begin a single transaction
	tx := db.DB.Begin()

	if err := tx.Create(&expense).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save expense"})
		return
	}

	approval.ExpenseID = expense.ID
	if err := tx.Create(&approval).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save approval"})
		return
	}

	// If auto-approved, process payment and update the same transaction
	if expense.AutoApproved {
		updatedExpense, updatedApproval, err := services.NewPaymentService().ProcessPayment(c.Request.Context(), expense, approval)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := tx.Save(updatedExpense).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense after payment"})
			return
		}

		if err := tx.Save(updatedApproval).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update approval after payment"})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"expense":  expense,
		"approval": approval,
	})
}

// ---- Approve Expense ----
func ApproveExpense(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		ApproverID *int64 `json:"approver_id"`
		Notes      string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var expense models.Expense
	if err := db.DB.Preload("Approval").First(&expense, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	if err := rules.CanApproveExpense(expense.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update expense & approval
	expense.Approval.ApproverID = input.ApproverID
	expense.Approval.Notes = input.Notes

	tx := db.DB.Begin()
	updatedExpense, updatedApproval, err := services.NewPaymentService().ProcessPayment(c.Request.Context(), &expense, expense.Approval)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Save(&updatedExpense).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense"})
		return
	}
	if err := tx.Save(&updatedApproval).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update approval"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, expense)
}

func RejectExpense(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		ApproverID *int64 `json:"approver_id"`
		Notes      string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var expense models.Expense
	if err := db.DB.Preload("Approval").First(&expense, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	if err := rules.CanRejectExpense(expense.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update expense & approval
	expense.Status = constants.ExpenseStatusRejected
	expense.Approval.Status = constants.ApprovalStatusRejected
	expense.Approval.ApproverID = input.ApproverID
	expense.Approval.Notes = input.Notes

	tx := db.DB.Begin()
	if err := tx.Save(&expense).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense"})
		return
	}
	if err := tx.Save(&expense.Approval).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update approval"})
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, expense)
}
