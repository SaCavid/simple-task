package handlers

import (
	"bytes"
	"github.com/SaCavid/simple-task/service"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
)

var (
	Users []string
)

func user(id int, wg *sync.WaitGroup) {

	defer wg.Done()
	payload := []byte(`{"name":"test user","age":30}`)
	request, _ := http.NewRequest("POST", "", bytes.NewBuffer(payload))

	log.Println(request.Response)

}

func BenchmarkServer_Handler(b *testing.B) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Print("No .env file found")
	}

	Database, err := service.Crea(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	err = Database.Find(&Users).Error
	if err != nil {
		log.Println(err)
	}

	var wg sync.WaitGroup
	log.Println("Users registered:", len(Users))
	for i := 1; i <= 1; i++ {
		wg.Add(1)
		go user(i, &wg)
	}

	wg.Wait()
}
