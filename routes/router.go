package routes

import (
	"vaqua/handler"
	"vaqua/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(
	userHandler *handler.UserHandler,
	transferRequestHandler *handler.TransferHandler,
	transactionHandler *handler.TransactionHandler,
	db *gorm.DB,
) *gin.Engine {
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

	// Public routes
	r.POST("/signup", userHandler.SignUpNewUserAcct)
	r.POST("/login", userHandler.LoginUser)



	//authorised

	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())

	//user routes

	userRoutes := authorized.Group("/user")
	{
	userRoutes.POST("/logout", userHandler.LogoutUser)
	userRoutes.PATCH("/profile", userHandler.UpdateUserProfile)
	userRoutes.GET("/id/me", userHandler.GetUserByID)
	userRoutes.GET("email/me", userHandler.GetUserByEmail)
	}

	//dashboard routes

	dashboardRoutes := authorized.Group("/dashboard")
	{
		dashboardRoutes.GET("/income", transactionHandler.GetUserIncome)
		dashboardRoutes.GET("/expenses", transactionHandler.GetUserExpenses)
		dashboardRoutes.GET("/balance", transactionHandler.GetBalance)
		dashboardRoutes.GET("/transactions", transactionHandler.GetAllTransactions)
		dashboardRoutes.GET("/transaction/:id", transactionHandler.GetTransaction)
		dashboardRoutes.POST("/transfer", transferRequestHandler.CreateTransfer)
		

	}
	
	return r
}
