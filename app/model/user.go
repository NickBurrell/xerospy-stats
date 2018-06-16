package model

import (
	"regexp"
	"time"
)

type User struct {
	CreatedAt time.Time
	Username  string `gorm:"not null;unique"`
	Email     string `gorm:"not null;unique"`
	Password  string `gorm:"not null"`
	APIKey    string `gorm:"not null;size:36;unique"`
}

func ValidateEmail(email string) bool {

	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	return re.MatchString(email)

}
