package dao

import (
	"eGZ-Board/internal/model"
	"eGZ-Board/internal/model/entity"

	"gorm.io/gorm/clause"
)

// SetWeather Set
func SetWeather(key string, value string) error {
	weather := entity.WeatherTable{Key: key, Value: value}
	return model.GetDatabase().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&weather).Error
}

func GetWeather(key string) (entity.WeatherTable, error) {
	var weatherTable entity.WeatherTable
	err := model.GetDatabase().Where("key = ?", key).First(&weatherTable).Error
	return weatherTable, err
}
