package model

type TeamMatch struct {
	Id                       int    `gorm:"column:_id;primary_key;AUTO_INCREMENT"`
	TeamId                   Team   `gorm:"column:team_id"`
	EventId                  Event  `gorm:"column:event_id"`
	Alliance                 string `gorm:"column:alliance;size:45"`
	Position                 int    `gorm:"column:position"`
	ScoutName                string `gorm:"column:scout_name;size:2000"`
	StrategyComments         string `gorm:"column:strategy_comments;size:2000"`
	RobotPerformanceComments string `gorm:"column:robot_performance_comments;size:2000"`
	DriveTeamComments        string `gorm:"column:drive_team_comments;size:2000"`
	FinalThoughts            string `gorm:"column:final_thoughts;size:2000"`
	DriveTeamSkill           int    `gorm:"column:drive_team_skill"`
	PreMatchCooperation      int    `gorm:"column:pre_match_cooperation"`
	InMatchCooperation       int    `gorm:"column:in_match_cooperation"`
	ScoutingDataQuality      int    `gorm:"column:scouting_data_quality"`
	GraciousProfessionalism  int    `gorm:"column:gracious_professionalism"`
}

func (tm TeamMatch) TableName() string {
	return "team_match"
}
