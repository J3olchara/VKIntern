package db

import (
	"fmt"
	"github.com/J3olchara/VKIntern/app/server/db/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
	"sync"
)

type DB struct {
	*gorm.DB
}

var (
	once sync.Once

	Conn *DB
)

func NewConnection() *DB {
	once.Do(func() {
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		dbname := os.Getenv("POSTGRES_DB")
		host := os.Getenv("POSTGRES_HOST")
		port := os.Getenv("POSTGRES_PORT")

		connectionData := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, host, port)
		db, err := gorm.Open("postgres", connectionData)
		if err != nil {
			log.Fatal(err)
		}

		Conn = &DB{db}
	})
	return Conn
}

func (db *DB) Prepare() {
	db.AutoMigrate(&models.Actor{}, &models.Film{}, &models.FilmActor{})
}
