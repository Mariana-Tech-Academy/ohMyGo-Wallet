package main

import (
	"fmt"
	"os"
	"vaqua/config"
	"vaqua/db"
	"vaqua/handler"
	"vaqua/models"
	"vaqua/repository"
	"vaqua/routes"
	"vaqua/services"

	"log"
	"vaqua/redis"
)

func main() {
	// load up variables
	config.LoadEnv()

	// connect to database
	db := db.InitDb()
	db.AutoMigrate(&models.User{})

    // connect to Redis
	if err := redis.ConnectRedis(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis successfully!")
	
	// initialise the repo
	userRepo := &repository.UserRepo{}
	transferRequestRepo := &repository.TransferRequestRepo{}
	transactionRepo := &repository.TransactionRepo{}
	// initialise the service
	userService := &services.UserService{Repo: userRepo}
	transferRequestService := &services.TransferRequestService{Repo: transferRequestRepo}
	transactionService := &services.TransactionService{Repo: transactionRepo}

	// initialise the handler
	userHandler := &handler.UserHandler{Service: userService}
	transferRequestHandler := &handler.TransferRequestHandler{Service: transferRequestService}
	transactionHandler := &handler.TransactionHandler{Service: transactionService}

	// setup routes
	r := routes.SetupRouter(userHandler, transferRequestHandler, transactionHandler, db)

	// start the server
	fmt.Println("server is running on localhost:8080...")
	fmt.Println("TEST_VAR:", os.Getenv("TEST_VAR"))
	r.Run(":8080")

}
