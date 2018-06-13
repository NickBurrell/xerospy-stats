// +build !mysql,sqlite

package app

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/zero-frost/xerospy-stats/app/model"
	"os"
)

func (s *Server) InitDatabase() {
	if _, err := os.Stat("./db/"); err != nil {
		os.Mkdir("./db", 0755)
	}

	s.Database, _ = gorm.Open("sqlite3", "./db/"+s.ServerConfig.DatabaseSettings.DatabaseName+".db")

	s.Database.AutoMigrate(&model.User{})
	s.Database.AutoMigrate(&model.Team{})
	s.Database.AutoMigrate(&model.Event{})
	s.Database.AutoMigrate(&model.Match{})
	s.Database.AutoMigrate(&model.ActionType{})
	s.Database.AutoMigrate(&model.TeamMatch{})
	s.Database.AutoMigrate(&model.TeamMatchAction{})

}
