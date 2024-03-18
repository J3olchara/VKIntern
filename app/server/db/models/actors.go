package models

import (
	"github.com/J3olchara/VKIntern/app/server/db"
	"time"
)

type Actor struct {
	ID       uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name     string `json:"name" gorm:"type:varchar(150);not null"`
	Gender   bool   `json:"gender" gorm:"not null"`
	Birthday Date   `json:"birthday" gorm:"type:date;not null"`
	Films    []Film `gorm:"many2many:actor_films;"`
}

func NewActor(name string, gender bool, Birthday time.Time, films []Film) *Actor {
	return &Actor{
		Name:   name,
		Gender: gender,
		Birthday: Date(time.Date(Birthday.Year(), Birthday.Month(), Birthday.Day(),
			0, 0, 0, 0, Birthday.Location())),
		Films: films,
	}
}

func (a *Actor) Valid() bool {
	return a.ValidName()
}

func (a *Actor) ValidName() bool {
	return len(a.Name) >= 1 && len(a.Name) <= 150
}

func (a *Actor) Save() bool {
	if !a.Valid() {
		return false
	}
	res := db.Conn.Omit("id").Save(a)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected == 1
}

func (a *Actor) Delete() bool {
	res := db.Conn.Delete(a, a.ID)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected == 1
}

func (a *Actor) Create() bool {
	if !a.Valid() {
		return false
	}
	res := db.Conn.Omit("id").Create(a)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected == 1
}
