package handlers

import (
	"bytes"
	"fmt"
	"github.com/SaCavid/simple-task/models"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
)

var (
	Users []models.User
)

func user(id int, wg *sync.WaitGroup) {
	log.SetFlags(log.Lshortfile)
	defer wg.Done()
	payload := []byte(`{"state": "win", "amount": "10.15", "transactionId": "some generated identificator"}`)
	resp, err := http.NewRequest("POST", "", bytes.NewBuffer(payload))
	if err != nil {
		log.Println(id, err)
	}

	defer resp.Body.Close()
	for true {

		bs := make([]byte, 1014)
		n, err := resp.Body.Read(bs)
		fmt.Println(string(bs[:n]))

		if n == 0 || err != nil {
			break
		}
	}
}

func BenchmarkServer_Handler(b *testing.B) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Print("No .env file found")
	}

	Database, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	err = Database.Find(&Users).Error
	if err != nil {
		log.Println(err)
	}

	var wg sync.WaitGroup

	log.Println("Users registered:", len(Users))
	for i := 1; i <= len(Users); i++ {
		wg.Add(1)
		go user(i, &wg)
	}

	wg.Wait()
}
