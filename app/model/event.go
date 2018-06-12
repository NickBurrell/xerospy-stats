package model

type Event struct {
	Id            int    `gorm:"column:_id;primary_key;AUTO_INCREMENT"`
	TBAEventKey   string `gorm:"column:tba_event_key;size:45"`
	Name          string `gorm:"column:name;size:255"`
	ShortName     string `gorm:"column:short_name;size:255"`
	EventType     string `gorm:"column:event_type;size:255"`
	EventDistrict string `gorm:"column:event_district;size:255"`
	Year          int    `gorm:"column:year"`
	Week          int    `gorm:"column:week"`
	Location      string `gorm:"column:location;size:255"`
	TBAEventCode  string `gorm:"column:tba_event_code;size:45"`
}

func (e Event) TableName() string {
	return "event"
}
