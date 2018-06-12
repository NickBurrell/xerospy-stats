package model

type ActionType struct {
	Id             int     `gorm:"column:_id;primary_key;AUTO_INCREMENT"`
	Name           string  `gorm:"column:name;size:255"`
	Description    string  `gorm:"column:description;size:2000"`
	MatchPhase     string  `gorm:"column:match_phase;size:45"`
	Points         float32 `gorm:"column:points"`
	OpponentPoints int     `gorm:"column:opponent_points"`
	QualPoints     int     `gorm:"column:qual_points"`
	FoulPoints     int     `gorm:"column:foul_points"`
	CoopFlag       string  `gorm:"column:coop_flag;size:1"`
	Category       string  `gorm:"column:category;size:255"`
}

func (at ActionType) TableName() string {
	return "action_type"
}
