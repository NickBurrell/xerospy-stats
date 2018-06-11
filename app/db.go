package app

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func InitDatabase(user string, password string, address string, db_name string) {
	db_conn, _ := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)%s?charset=utf8", user, password, address, db_name))
	db = db_conn
}

func GetDatabase() *gorm.DB {
	return db
}
