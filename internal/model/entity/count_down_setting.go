package entity

// CountDownTable /*
type CountDownTable struct {
	Id     int    `gorm:"primary_key;auto_increment;not null"`
	Event  string `gorm:"type:varchar(255);not null"`
	Date   string `gorm:"type:varchar(255);not null"`
	Type   string `gorm:"type:varchar(255);not null"`
	During int    `gorm:"default: 0"`
	Class  int    `gorm:"not null"`
}
