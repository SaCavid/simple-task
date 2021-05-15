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

	port := os.Getenv("HTTP_SERVER_PORT")

	srv := handlers.Server{
		TransactionIds: make(map[string]string, 0),
		UserBalances:   make(map[string]models.Balance, 0),
		Repo:           service.NewTaskRepository(os.Getenv("DATABASE_URL")),
	}

	err := srv.FetchUsers()
	if err != nil {
		log.Println(err)
	}

	go srv.BulkInsertTransactions()
	go srv.BulkUpdateBalances()
	go srv.PostProcessing()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Task!")
	})

	e.GET("/api/users", srv.FetchUsersForTesting)
	e.POST("/api/register", srv.Register)

	e.GET("/api/test", srv.Test)

	e.POST("/api/processing", srv.Handler)
	s := &http.Server{
		ReadTimeout: 5 * time.Second,
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)), s)
}
