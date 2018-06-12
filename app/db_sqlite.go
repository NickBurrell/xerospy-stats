// +build !mysql,sqlite

package app

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/zero-frost/xerospy-stats/app/model"
	"os"
)

var db *gorm.DB

func InitDatabase(user string, password string, address string, db_name string) {

	if _, err := os.Stat("./db/"); err != nil {
		os.Mkdir("./db", 0755)
	}

	db_conn, _ := gorm.Open("sqlite3", "./db/"+db_name+".db")

	db_conn.AutoMigrate(&model.User{})
	db_conn.AutoMigrate(&model.Team{})
	db_conn.AutoMigrate(&model.Event{})
	db_conn.AutoMigrate(&model.Match{})
	db_conn.AutoMigrate(&model.ActionType{})
	db_conn.AutoMigrate(&model.TeamMatch{})
	db_conn.AutoMigrate(&model.TeamMatchAction{})

	db = db_conn

}

func GetDatabase() *gorm.DB {
	return db
}
