package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"

	_ "backend/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func AuthRoutes(r *gin.RouterGroup) {

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health-check", controllers.HealthCheck)

	auth := r.Group("/auth")
	{
		auth.POST("/login", controllers.Login)
	}

	protected := r.Group("/", middleware.JWTAuthMiddleware())

	manager := protected.Group("/manager", middleware.RequireRole("manager"))

	manager.GET("/dashboard", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome Manager"})
	})

	managerExpenses := manager.Group("/expenses")
	{
		managerExpenses.GET("", controllers.GetExpenses)
		managerExpenses.GET("/:id", controllers.GetExpense)
		managerExpenses.PUT("/:id/approve", controllers.ApproveExpense)
		managerExpenses.PUT("/:id/reject", controllers.RejectExpense)
	}

	user := protected.Group("/user", middleware.RequireRole("user"))

	user.GET("/dashboard", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome user"})
	})

	// user expense routes
	userExpenses := user.Group("/expenses")
	{
		userExpenses.GET("", controllers.GetUserExpenses)
		userExpenses.GET("/:id", controllers.GetExpense)
		userExpenses.POST("", controllers.CreateExpense)
	}
}
