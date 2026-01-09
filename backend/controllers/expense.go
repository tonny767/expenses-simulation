package controllers

import (
	"net/http"

	"backend/actions"
	"backend/constants"
	"backend/db"
	"backend/helpers"
	"backend/rules"
	"backend/workers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/models"
)

func GetExpenses(c *gin.Context) {
	var expenses []models.Expense
	var total int64

	page, limit, offset := helpers.GetPagination(c)
	status := c.Query("status")

	query := db.DB.Model(&models.Expense{}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("Approval")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// count first
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count expenses"})
		return
	}

	// fetch paginated data
	if err := query.
		Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&expenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": expenses,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func GetUserExpenses(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var expenses []models.Expense
	var total int64

	page, limit, offset := helpers.GetPagination(c)
	status := c.Query("status")

	query := db.DB.Model(&models.Expense{}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("Approval").
		Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count user expenses"})
		return
	}

	if err := query.
		Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&expenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user expenses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": expenses,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func GetExpense(c *gin.Context) {
	id := c.Param("id")

	var expense models.Expense
	if err := db.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).
		Preload("Approval").First(&expense, "id = ?", id).Error; err != nil {
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

	tx.Commit()

	if expense.AutoApproved {
		worker := workers.NewPaymentWorker()
		worker.ProcessExpensePaymentAsync(expense.ID)
	}

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

	tx := db.DB.Begin()
	expense.Status = constants.ExpenseStatusApproved

	expense.Approval.ApproverID = input.ApproverID
	expense.Approval.Notes = input.Notes

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

	worker := workers.NewPaymentWorker()
	worker.ProcessExpensePaymentAsync(expense.ID)

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
