package main

import (
	"fmt"
	"github.com/SaCavid/simple-task/handlers"
	"github.com/SaCavid/simple-task/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	port := os.Getenv("PORT")
	address := os.Getenv("ADDRESS")
	endPoint := os.Getenv("TASK_END_POINT")

	if port == "" {
		port = "8080"
	}
	if address == "" {
		address = "127.0.0.1"
	}

	srv := handlers.Server{
		TransactionIds: make(map[string]string, 0),
		Repo:           service.NewTaskRepository(),
	}

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST(endPoint, srv.Handler)

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", address, port)))
}
