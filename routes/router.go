package routes

import (
	"vaqua/handler"
	"vaqua/middleware"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func SetupRouter(userHandler *handler.UserHandler, transferRequestHandler *handler.TransferRequestHandler, transactionHandler *handler.TransactionHandler, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	auth := r.Group("/")

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		sqlDB, err:= db.DB()
		if err != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}

		if err := sqlDB.Ping(); err != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "error": "cannot connect to database"})
			return
		}

		c.JSON(200, gin.H{"status": "healthy", "db": "connected to database"})
	})

	// public routes
	r.POST("/signup", userHandler.SignUpNewUserAcct)
	r.POST("/login", userHandler.LoginUser)
	// r.GET("/user/:id", userHandler.GetUserByID)
	// r.GET("/user", userHandler.GetUserByEmail)



	//  Authenticated user routes
	auth.Use(middleware.AuthMiddleware())
	{
		auth.PATCH("/profile", userHandler.UpdateUserProfile)
		auth.GET("/user/id/me", userHandler.GetUserByID)
		auth.GET("/user/email/me", userHandler.GetUserByEmail)
	}

	r.POST("/logout", userHandler.LogoutUser)

	auth.GET("/dashboard/expenses", transactionHandler.GetExpensesHandler)
	auth.GET("/transactions/:id", transactionHandler.GetTransactionByIDHandler)
	return r

}
