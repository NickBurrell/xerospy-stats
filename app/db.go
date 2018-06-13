// +build mysql !sqlite

package app

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func (s *Server) InitDatabase() {
	s.Database, _ = gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)%s?charset=utf8",
			s.ServerConfig.DatabaseSettings.Username,
			s.ServerConfig.DatabaseSettings.Password,
			s.ServerConfig.DatabaseSettings.Address,
			s.ServerConfig.DatabaseSettings.DatabaseName))
}
