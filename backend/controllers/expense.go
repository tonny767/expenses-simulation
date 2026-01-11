package controllers

import (
	"net/http"
	"time"

	"backend/actions"
	"backend/db"
	"backend/helpers"
	"backend/rules"
	"backend/workers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/models"
)

type PaginationMeta struct {
	Page  int   `json:"page" example:"1"`
	Limit int   `json:"limit" example:"10"`
	Total int64 `json:"total" example:"42"`
}

type ExpensesListResponse struct {
	Data []models.Expense `json:"data"`
	Meta PaginationMeta   `json:"meta"`
}
type MessageResponse struct {
	Message string `json:"message" example:"operation successful"`
}

type StatusExpenseRequest struct {
	ApproverID *int64 `json:"approver_id" example:"99"`
	Notes      string `json:"notes" example:"Approved"`
}

type HealthCheckResponse struct {
	Status   string    `json:"status" example:"ok"`
	Database string    `json:"database" example:"up"`
	Time     time.Time `json:"time" example:"2023-10-05T14:48:00Z"`
}

// HealthCheck godoc
// @Summary Health check
// @Description Check the health status of the application and database connectivity
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} HealthCheckResponse
// @Failure 503 {object} HealthCheckResponse
// @Router /health-check [get]
func HealthCheck(c *gin.Context) {
	sqlDB, err := db.DB.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, HealthCheckResponse{
			Status:   "down",
			Database: "error",
			Time:     time.Now().UTC(),
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, HealthCheckResponse{
			Status:   "down",
			Database: "down",
			Time:     time.Now().UTC(),
		})
		return
	}

	c.JSON(http.StatusOK, HealthCheckResponse{
		Status:   "ok",
		Database: "up",
		Time:     time.Now().UTC(),
	})
}

// GetExpenses godoc
// @Summary Get all expenses (manager only)
// @Description Get paginated list of all expenses (manager only)
// @Tags ManagerExpenses
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param status query string false "Filter by expense status"
// @Success 200 {object} ExpensesListResponse
// @Failure 401 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /manager/expenses [get]
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

	c.JSON(http.StatusOK, ExpensesListResponse{
		Data: expenses,
		Meta: PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	})
}

// GetUserExpenses godoc
// @Summary Get user's expenses
// @Description Get paginated list of expenses for the authenticated user
// @Tags Expenses
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param status query string false "Filter by expense status"
// @Success 200 {object} object{models.Expense}
// @Failure 401 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /expenses [get]
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

	c.JSON(http.StatusOK, ExpensesListResponse{
		Data: expenses,
		Meta: PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	})
}

// GetExpense godoc
// @Summary Get expense by ID
// @Description Get detailed information of a single expense
// @Tags Expenses
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param id path int true "Expense ID"
// @Success 200 {object} models.Expense
// @Failure 401 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Router /expenses/{id} [get]
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

// CreateExpense godoc
// @Summary Create a new expense
// @Description Submit a new expense for reimbursement
// @Tags Expenses
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param request body actions.SubmitExpenseInput true "Expense payload"
// @Success 201 {object} MessageResponse
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /expenses [post]
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

	c.JSON(http.StatusCreated, MessageResponse{
		Message: "Expense created successfully",
	})
}

// ApproveExpense godoc
// @Summary Approve an expense
// @Description Approve a pending expense (manager only)
// @Tags Manager
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param id path int true "Expense ID"
// @Param request body StatusExpenseRequest true "Approval payload"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /manager/expenses/{id}/approve [put]
func ApproveExpense(c *gin.Context) {
	id := c.Param("id")
	var input StatusExpenseRequest

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

	// part of log
	originalStatus := expense.Status

	updatedExpense, updatedApproval, err := actions.ApproveExpense(actions.ApproveExpenseInput{
		Expense:    &expense,
		ApproverID: input.ApproverID,
		Notes:      input.Notes,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	audit, err := actions.ExpenseAuditLog(actions.ExpenseAuditLogInput{
		ExpenseID:  updatedExpense.ID,
		ActorID:    input.ApproverID,
		FromStatus: originalStatus,
		ToStatus:   updatedExpense.Status,
		Reason:     "Expense approved",
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := db.DB.Begin()

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

	if err := tx.Create(&audit).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create audit log"})
		return
	}

	tx.Commit()

	worker := workers.NewPaymentWorker()
	worker.ProcessExpensePaymentAsync(updatedExpense.ID)

	c.JSON(http.StatusOK, MessageResponse{
		Message: "Expense has been approved",
	})
}

// RejectExpense godoc
// @Summary Reject an expense
// @Description Reject a pending expense (manager only)
// @Tags Manager
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param id path int true "Expense ID"
// @Param request body StatusExpenseRequest true "Rejection payload"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} httputil.HTTPError
// @Failure 401 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /manager/expenses/{id}/reject [put]
func RejectExpense(c *gin.Context) {
	id := c.Param("id")

	var input StatusExpenseRequest

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

	// part of log
	originalStatus := expense.Status

	// Update expense & approval
	updatedExpense, updatedApproval, err := actions.RejectExpense(actions.RejectExpenseInput{
		Expense:    &expense,
		ApproverID: input.ApproverID,
		Notes:      input.Notes,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	audit, err := actions.ExpenseAuditLog(actions.ExpenseAuditLogInput{
		ExpenseID:  updatedExpense.ID,
		ActorID:    input.ApproverID,
		FromStatus: originalStatus,
		ToStatus:   updatedExpense.Status,
		Reason:     "Expense rejected",
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := db.DB.Begin()
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

	if err := tx.Create(&audit).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create audit log"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, MessageResponse{
		Message: "Expense has been rejected",
	})
}
