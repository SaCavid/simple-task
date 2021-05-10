package service

import (
	"github.com/SaCavid/simple-task/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
)

type TaskRepository struct {
	Db *gorm.DB
}

func NewTaskRepository() *TaskRepository {
	taskRepo, err := CreateDbConnection(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	return &TaskRepository{Db: taskRepo}
}

func CreateDbConnection(connectionUri string) (*gorm.DB, error) {

	db, err := gorm.Open("postgres", connectionUri)
	if err != nil {
		return nil, err
	}

	b := os.Getenv("DROP_TABLES")
	if b == "true" {
		log.Println("Dropping tables data and users")
		db.DropTableIfExists(&models.Data{}, &models.User{})
	}

	db.AutoMigrate(&models.Data{}, &models.User{})

	db.DropTableIfExists()
	return db, nil
}
