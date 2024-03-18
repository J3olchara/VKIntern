package models

import (
	"github.com/J3olchara/VKIntern/app/server/db"
	"log"
)

type User struct {
	Username string `gorm:"unique"`
	Password string
	Staff    bool
}

func (u *User) Create() {
	db.Conn.Create(u)
	db.Conn.First(u, "username = ?", u.Username)
}

func (u *User) Delete() {
	res := db.Conn.Where("username = ?", u.Username).Delete(u)
	log.Println(res.Error)
}

//func HashPassword(password string) (string, error) {
//	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
//	return string(bytes), err
//}
//
//func CheckPasswordHash(password, hash string) bool {
//	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
//	return err == nil
//}
//
//func Authenticate(username string, password string) bool {
//	var user User
//	db.Conn.Where("username = ?", username).First(&user)
//	return CheckPasswordHash(password, user.Password)
//}
// для удобства убрал чтобы можно было юзера создавать прям в базе и закидывать в заголовки

func Authenticate(username string, password string) bool {
	var user User
	db.Conn.Where("username = ?", username).First(&user)
	return user.Password == password
}
