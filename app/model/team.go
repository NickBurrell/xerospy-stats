package model

type Team struct {
	Id               int    `gorm:"column:_id;primary_key;AUTO_INCREMENT"`
	TeamNumber       string `gorm:"column:team_number;size:25"`
	TBATeamKey       string `gorm:"column:tba_team_key;size:45"`
	LongName         string `gorm:"column:long_name;size:255"`
	Name             string `gorm:"column:name;size:255"`
	LogoFileLocation string `gorm:"column:logo_file_location;size:2000"`
	City             string `gorm:"column:city;size:255"`
	StateCode        string `gorm:"column:state_code;size:45"`
	Country          string `gorm:"column:country;size:255"`
	Motto            string `gorm:"column:motto;size:2000"`
	RookieYear       int    `gorm:"column:rookie_year"`
}

func (t Team) TableName() string {
	return "team"
}
