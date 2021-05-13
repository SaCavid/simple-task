package service

import (
	"github.com/SaCavid/simple-task/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
	"time"
)

type TaskRepository struct {
	Db *gorm.DB
}

func NewTaskRepository() *TaskRepository {

	// docker-compose sometimes starts processing container faster than expected.
	// timeout for not to get error. docker-compose depends on configuration didnt helps. can be adjusted.
	time.Sleep(10 * time.Second)
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

	db.AutoMigrate(&models.Data{}, &models.User{})

	// while development can be triggered to drop database tables
	b := os.Getenv("DROP_TABLES")
	if b == "true" {
		log.Println("Dropping tables data and users")
		db.DropTableIfExists(&models.Data{}, &models.User{})
	}
	return db, nil
}
