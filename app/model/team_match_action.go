package model

import "time"

type TeamMatchAction struct {
	Id           int        `gorm:"column:_id;primary_key;AUTO_INCREMENT"`
	TeamMatchId  TeamMatch  `gorm:"column:team_match_id"`
	ActionTypeId ActionType `gorm:"column:action_type_id"`
	Quantity     int        `gorm:"column:quantity"`
	StartTime    time.Time  `gorm:"column:start_time"`
	EndTime      time.Time  `gorm:"column:end_time"`
	ObjectCount  int        `gorm:"column:object_count"`
	TabletUUID   string     `gorm:"column:tablet_uuid;size:500"`
}

func (tma TeamMatchAction) TableName() string {
	return "team_match_action"
}
