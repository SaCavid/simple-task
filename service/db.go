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
	taskRepo, err := CreateDbConnectionSensors(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	return &TaskRepository{Db: taskRepo}
}

func CreateDbConnectionSensors(connectionUri string) (*gorm.DB, error) {

	db, err := gorm.Open("postgres", connectionUri)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Data{})

	return db, nil
}
