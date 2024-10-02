package entity

// SettingTable /*
type SettingTable struct {
	Id    int    `gorm:"primary_key;auto_increment;not null"`
	Key   string `gorm:"type:varchar(255);not null"`
	Value string `gorm:"type:varchar(255);not null"`
	Class int    `gorm:"not null"`
}
