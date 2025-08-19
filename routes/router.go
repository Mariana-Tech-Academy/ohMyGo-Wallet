package routes

import (
	"vaqua/handler"
	"vaqua/middleware"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func SetupRouter(userHandler *handler.UserHandler, transferRequestHandler *handler.TransferHandler, transactionHandler *handler.TransactionHandler, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		sqlDB, err := db.DB()
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

	//  Authenticated user routes updated
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())

	authorized.POST("/logout", userHandler.LogoutUser)
	authorized.POST("/transfer", transferRequestHandler.CreateTransfer)
	authorized.PATCH("/profile", userHandler.UpdateUserProfile)
	authorized.GET("/user/id/me", userHandler.GetUserByID)
	authorized.GET("/user/email/me", userHandler.GetUserByEmail)

	authorized.GET("/dashboard/expenses", transactionHandler.GetExpensesHandler)
	authorized.GET("/transactions/:id", transactionHandler.GetTransactionByIDHandler)
	return r

}
