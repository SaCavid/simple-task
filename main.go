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
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Task!")
	})

	// for registering users.
	e.GET("/api/users", srv.FetchUsersForTesting)
	e.POST("/api/register", srv.Register)

	// main route for processing transactions
	e.POST("/api/processing", srv.Handler)
	s := &http.Server{
		ReadTimeout: 5 * time.Second,
	}

	// Starting HTTP server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)), s)
}
