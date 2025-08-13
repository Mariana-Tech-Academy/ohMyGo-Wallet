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
)

func main() {
	// load up variables
	config.LoadEnv()

	// connect to database
	db := db.InitDb()
	db.AutoMigrate(&models.User{})

	// initialise the repo
	userRepo := &repository.UserRepo{DB: db}
	transferRequestRepo := &repository.TransferRequestRepo{DB: db}
	transactionRepo := &repository.TransactionRepo{DB: db}
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
