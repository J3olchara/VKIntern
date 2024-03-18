package main

import (
	"github.com/J3olchara/VKIntern/app/server/db"
	"github.com/J3olchara/VKIntern/app/server/db/models"
	"github.com/J3olchara/VKIntern/app/server/server"
	"log"
)

func main() {
	db.NewConnection()
	defer func(Conn *db.DB) {
		err := Conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db.Conn)
	db.Conn.AutoMigrate(&models.Actor{}, &models.Film{}, &models.FilmActor{}, &models.User{})

	server.StartServer()
}
