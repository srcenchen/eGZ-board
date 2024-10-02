package entity

// ScheduleTable /*
type ScheduleTable struct {
	Id       int    `gorm:"primary_key;auto_increment;not null"`
	Schedule string `gorm:"type:varchar(255);not null"`
	Class    int    `gorm:"not null"`
}
