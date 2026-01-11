package controllers

import (
	"backend/db"
	"backend/helpers"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAuditLog godoc
// @Summary Get audit logs
// @Description Get paginated list of audit logs for expenses status changes (manager only)
// @Tags Manager
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param expense_id query int false "Filter by expense ID"
// @Success 200 {object} object{data=[]models.ExpenseAuditLog,meta=PaginationMeta}
// @Failure 401 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /manager/expense-logs [get]
func GetExpenseAuditLog(c *gin.Context) {
	var auditLogs []models.ExpenseAuditLog
	var total int64

	page, limit, offset := helpers.GetPagination(c)
	expenseID := c.Query("expense_id")

	query := db.DB.Model(&models.ExpenseAuditLog{})

	if expenseID != "" {
		query = query.Where("expense_id = ?", expenseID)
	}

	// count first
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count audit logs"})
		return
	}

	// fetch paginated data
	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&auditLogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audit logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": auditLogs,
		"meta": PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	})
}
