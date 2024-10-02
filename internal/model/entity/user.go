package entity

// UserTable /*
type UserTable struct {
	Id          int    `gorm:"primary_key;auto_increment;not null"`
	Username    string `gorm:"type:varchar(255);not null"`
	Password    string `gorm:"type:varchar(255);not null"`
	PhoneNumber string `gorm:"type:varchar(255);not null"`
	Role        string `gorm:"type:varchar(255);not null"`
}

type ClassTable struct {
	Id    int `gorm:"primary_key;auto_increment;not null"`
	Grade int `gorm:"not null"`
	Class int `gorm:"not null"`
}

type DeviceTable struct {
	Id       int    `gorm:"primary_key;auto_increment;not null"`
	ClassID  int    `gorm:"not null"`
	DeviceID string `gorm:"not null"`
}
