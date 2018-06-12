package model

type Match struct {
	Id                int    `gorm:"column:_id;primary_key;AUTO_INCREMENT"`
	EventId           Event  `gorm:"column:event_id"`
	TBAMatchKey       string `gorm:"column:tba_match_key;size:45"`
	CompLevel         string `gorm:"column:comp_level;size:45"`
	SetNumber         int    `gorm:"column:set_number"`
	MatchNumber       int    `gorm:"column:match_number"`
	Status            string `gorm:"column:status;size:45"`
	Red1TeamId        int    `gorm:"column:red_1_team_id"`
	Red2TeamId        int    `gorm:"column:red_2_team_id"`
	Red3TeamId        int    `gorm:"column:red_3_team_id"`
	RedAutoScore      int    `gorm:"column:red_auto_score"`
	RedTeleopScore    int    `gorm:"column:red_teleop_score"`
	RedTotalScore     int    `gorm:"column:red_total_score"`
	RedQp             int    `gorm:"column:red_qp"`
	RedFoulPoints     int    `gorm:"column:red_foul_points"`
	Blue1TeamId       int    `gorm:"column:blue_1_team_id"`
	Blue2TeamId       int    `gorm:"column:blue_2_team_id"`
	Blue3TeamId       int    `gorm:"column:blue_3_team_id"`
	BlueAutoScore     int    `gorm:"column:blue_auto_score"`
	BlueTeleopScore   int    `gorm:"column:blue_teleop_score"`
	BlueTotalScore    int    `gorm:"column:blue_total_score"`
	BlueQp            int    `gorm:"column:blue_qp"`
	BlueFoulPoints    int    `gorm:"column:blue_foul_points"`
	Winner            string `gorm:"column:winner;size:25"`
	DriveTeamComments string `gorm:"column:drive_team_comments;size:2000"`
}

func (m Match) TableName() string {
	return "match"
}
