package main

import (
	"fmt"
	"vaqua/config"
	"vaqua/db"
	"vaqua/handler"
	"vaqua/repository"
	"vaqua/routes"
	"vaqua/services"

	"log"
	"vaqua/redis"
)
//THIS IS A TEST
func main() {
	// load up variables
	config.LoadEnv()

	// connect to database
	db.InitDb()

    // connect to Redis
	if err := redis.ConnectRedis(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis successfully!")

	// initialise the repo
	userRepo := &repository.UserRepo{}
	transferRequestRepo := &repository.TransferRepo{}
	transactionRepo := &repository.TransactionRepo{}

	// initialise the service
	userService := &services.UserService{Repo: userRepo}
	transferService := &services.TransferServices{Repo: *transferRequestRepo}
	transactionService := &services.TransactionService{Repo: transactionRepo}

	// initialise the handler
	userHandler := &handler.UserHandler{Service: userService}
	transferRequestHandler := &handler.TransferHandler{Service: transferService}
	transactionHandler := &handler.TransactionHandler{Service: transactionService}

	// define routes
	db := db.InitDb()
	r := routes.SetupRouter(userHandler, transferRequestHandler, transactionHandler, db)

	// start the server
	fmt.Println("server is running on localhost:8080...")
	r.Run(":8080")

}
