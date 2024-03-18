package db

import (
	"database/sql"
	"fmt"
	"github.com/J3olchara/VKIntern/app/server/support"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

		connectionDataSql := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=disable", user, password, host, port)
		connectionData := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, host, port)
		dbSql, err := sql.Open("postgres", connectionDataSql)
		support.FatalErr(err)
		_, err = dbSql.Exec("CREATE DATABASE " + dbname)
		support.WarningErr(err)
		err = dbSql.Close()
		support.FatalErr(err)
		db, err := gorm.Open("postgres", connectionData)
		support.FatalErr(err)

		Conn = &DB{db}
	})
	return Conn
}
