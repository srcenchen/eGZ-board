package entity

// WeatherTable /*
type WeatherTable struct {
	Id    int    `gorm:"primary_key;auto_increment;not null"`
	Key   string `gorm:"type:varchar(255);not null"`
	Value string `gorm:"type:varchar(255);not null"`
}
