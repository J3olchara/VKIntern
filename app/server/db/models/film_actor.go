package models

type FilmActor struct {
	FilmID  uint `json:"film_id" gorm:"primary_key"`
	ActorID uint `json:"actor_id" gorm:"primary_key"`
}
