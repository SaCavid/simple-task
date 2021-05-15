package main

import (
	"fmt"
	"github.com/SaCavid/simple-task/handlers"
	"github.com/SaCavid/simple-task/models"
	"github.com/SaCavid/simple-task/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	// can be changed in env file. default 8080
	port := os.Getenv("HTTP_SERVER_PORT")

	// initialize server
	srv := handlers.Server{
		TransactionIds: make(map[string]string, 0),
		UserBalances:   make(map[string]models.Balance, 0),
		Repo:           service.NewTaskRepository(os.Getenv("DATABASE_URL")),
	}

	// fetching database information about users and transactions for further use
	err := srv.FetchData()
	if err != nil {
		log.Println(err)
	}

	// goroutine for bulk inserting transaction information to database
	go srv.BulkInsertTransactions()

	// goroutine updating user balances depended on transactions. bulk update.
	go srv.BulkUpdateBalances()

	// -- post processing task
	// Every N minutes 10 latest odd records must be canceled and balance should be corrected by the application.
	// Cancelled records shouldn't be processed twice.
	// can be changed from env file
	// default 5 minutes
	go srv.PostProcessing()

	// starting HTTP route
	e := echo.New()

	// Task url
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "-------\n\nThe main goal of this test task is to develop the application for processing the incoming requests from the 3d-party providers.\nThe application must have an HTTP URL to receive incoming POST requests.\nTo receive the incoming POST requests the application must have an HTTP URL endpoint.\n\nTechnologies: Golang + Postgres.\n\nRequirements:\n1. Processing and saving incoming requests.\n\nImagine that we have a user with the account balance.\n\nExample of the POST request:\nPOST /your_url HTTP/1.1\nSource-Type: client\nContent-Length: 34\nHost: 127.0.0.1\nContent-Type: application/json\n{\"state\": \"win\", \"amount\": \"10.15\", \"transactionId\": \"some generated identificator\"}\n\nHeader “Source-Type” could be in 3 types (game, server, payment). This type probably can be extended in the future.\n\nPossible states (win, lost):\n1. Win requests must increase the user balance\n2. Lost requests must decrease user balance.\nEach request (with the same transaction id) must be processed only once.\n\nThe decision regarding database architecture and table structure is made to you.\n\nYou should know that account balance can't be in a negative value.\nThe application must be competitive ability.\n\n2. Post-processing\nEvery N minutes 10 latest odd records must be canceled and balance should be corrected by the application.\nCancelled records shouldn't be processed twice.\n\n3. The application should be prepared for running via docker containers.\n\nPlease be informed and kindly note that application without description about how to run and test won't be accepted and reviewed. \n\n---------")
	})

	// for registering users
	e.GET("/api/users", srv.FetchUsersForTesting)
	e.POST("/api/register", srv.Register)

	// main route for processing transactions
	e.POST("/api/processing", srv.Handler)
	s := &http.Server{
		ReadTimeout: 5 * time.Second,
	}

	// starting HTTP server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)), s)
}
