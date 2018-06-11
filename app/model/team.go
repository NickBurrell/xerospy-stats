package model

import (
	_ "github.com/jinzhu/gorm"
)

type Team struct {
	Id               int    `gorm:"column:_id"`
	TeamNumber       string `gorm:"column:team_number"`
	TBATeamKey       string `gorm:"column:tba_team_key"`
	LongName         string `gorm:"column:long_name"`
	Name             string `gorm:"column:name"`
	LogoFileLocation string `gorm:"column:logo_file_location"`
	City             string `gorm:"column:city"`
	StateCode        string `gorm:"column:state_code"`
	Country          string `gorm:"column:country"`
	Motto            string `gorm:"column:motto"`
	RookieYear       int    `gorm:"column:rookie_year"`
}
