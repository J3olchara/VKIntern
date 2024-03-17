package models

type Film struct {
	ID          uint    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name        string  `json:"name" gorm:"type:varchar(150);check: name <> ''"`
	Description string  `json:"description" gorm:"type:varchar(1000)"`
	Date        Date    `json:"date" gorm:"type:date"`
	Rating      uint8   `json:"rating" gorm:"type:smallint;check: rating >= 0 and rating <= 10"`
	Actors      []Actor `gorm:"many2many:actor_films;"`
}
