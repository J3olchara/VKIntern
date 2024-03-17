package models

type Actor struct {
	ID       uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name     string `json:"name" gorm:"type:varchar(150);not null"`
	Gender   bool   `json:"gender" gorm:"not null"`
	Birthday Date   `json:"birthday" gorm:"type:date;not null"`
	Films    []Film `gorm:"many2many:actor_films;"`
}
