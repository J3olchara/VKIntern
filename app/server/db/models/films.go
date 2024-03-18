package models

import (
	"github.com/J3olchara/VKIntern/app/server/db"
	"time"
)

type Film struct {
	ID          uint    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name        string  `json:"name" gorm:"type:varchar(150);check: name <> ''"`
	Description string  `json:"description" gorm:"type:varchar(1000)"`
	Date        Date    `json:"date" gorm:"type:date"`
	Rating      uint8   `json:"rating" gorm:"type:smallint;check: rating >= 0 and rating <= 10"`
	Actors      []Actor `gorm:"many2many:actor_films;"`
}

func NewFilm(name, description string, date time.Time, rating uint8, actors []Actor) *Film {
	return &Film{
		Name:        name,
		Description: description,
		Date: Date(time.Date(date.Year(), date.Month(), date.Day(),
			0, 0, 0, 0, date.Location())),
		Rating: rating,
		Actors: actors,
	}
}

func (f *Film) Valid() bool {
	return f.ValidName() && f.ValidRating() && f.ValidDescription()
}

func (f *Film) ValidName() bool {
	return len(f.Name) >= 1 && len(f.Name) <= 150
}

func (f *Film) ValidDescription() bool {
	return len(f.Description) <= 1000
}

func (f *Film) ValidRating() bool {
	return f.Rating <= 10
}

func (f *Film) Save() bool {
	if !f.Valid() {
		return false
	}
	res := db.Conn.Omit("id").Save(f)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected == 1
}

func (f *Film) Delete() bool {
	res := db.Conn.Delete(f, f.ID)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected == 1
}

func (f *Film) Create() bool {
	if !f.Valid() {
		return false
	}
	res := db.Conn.Omit("id").Create(f)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected == 1
}

func SearchFilms(search Search) []Film {
	var films []Film
	if search.Ordering == "" {
		search.Ordering = "desc"
	}
	if search.Field == "" {
		search.Field = "rating"
	}
	res := db.Conn.
		Preload("Actors").
		Joins("LEFT JOIN film_actors ON films.id = film_actors.film_id").
		Joins("LEFT JOIN actors ON film_actors.actor_id = actors.id").
		Where("LOWER(films.name) LIKE LOWER(?)", "%"+search.Search+"%")
	if search.ActorName != "" {
		res = res.Where("LOWER(actors.name) LIKE LOWER(?)", "%"+search.ActorName+"%")
	}
	res = res.
		Order(search.Field + " " + search.Ordering).
		Find(&films)
	if res.Error != nil {
		panic(res.Error)
	}
	return films
}
