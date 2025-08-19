package routes

import (
	"vaqua/handler"
	"vaqua/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(
	userHandler *handler.UserHandler,
	transferHandler *handler.TransferHandler,
	transactionHandler *handler.TransactionHandler,
	db *gorm.DB,
) *gin.Engine {
	r := gin.Default()

	// Health check
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

	// Public auth
	r.POST("/signup", userHandler.SignUpNewUserAcct)
	r.POST("/login", userHandler.LoginUser)

	// Authenticated routes
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())

	// User, profile and transfer
	auth.POST("/logout", userHandler.LogoutUser)
	auth.PATCH("/profile", userHandler.UpdateUserProfile)
	auth.POST("/transfer", transferHandler.CreateTransfer)
	auth.GET("/user/id/me", userHandler.GetUserByID)
	auth.GET("/user/email/me", userHandler.GetUserByEmail)

	// Transactions 
	auth.POST("/transactions", transactionHandler.CreateTransaction)  // create
	auth.GET("/transactions", transactionHandler.GetAllTransactions)  // list (paginated)

	

	return r
}
