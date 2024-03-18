package models

import "github.com/J3olchara/VKIntern/app/server/db"

type FilmActor struct {
	FilmID  uint `json:"film_id" gorm:"primary_key"`
	ActorID uint `json:"actor_id" gorm:"primary_key"`
}

func CreateRelations(film *Film, actors []Actor) {
	for _, actor := range actors {
		res := db.Conn.
			Set("gorm:insert_option", "ON CONFLICT DO NOTHING").
			Create(&FilmActor{
				FilmID:  film.ID,
				ActorID: actor.ID,
			})
		if res.Error != nil {
			panic(res.Error)
		}
	}
}

func DeleteRelations(film *Film) {
	db.Conn.Where("film.id = ?", film.ID).Delete(&FilmActor{})
}

func UpdateRelations(film *Film, actors []Actor) {
	DeleteRelations(film)
	CreateRelations(film, actors)
}
