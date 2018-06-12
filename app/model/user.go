package model

import (
	"github.com/jinzhu/gorm"
	"regexp"
)

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string `gorm:"not null"`
}

func ValidateEmail(email string) bool {

	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`]")

	return re.MatchString(email)

}
