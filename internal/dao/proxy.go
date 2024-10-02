package dao

import (
	"eGZ-Board/internal/model"
	"eGZ-Board/internal/model/entity"
)

// SetWeather Set
func SetWeather(key string, value string) {
	var weatherTable entity.WeatherTable
	if model.GetDatabase().Where("key = ?", key).First(&weatherTable).RowsAffected == 1 {
		weatherTable.Value = value
		model.GetDatabase().Save(&weatherTable)
	} else {
		model.GetDatabase().Create(&entity.WeatherTable{Key: key, Value: value})
	}
}

func GetWeather(key string) entity.WeatherTable {
	var weatherTable entity.WeatherTable
	model.GetDatabase().Where("key = ?", key).First(&weatherTable)
	return weatherTable
}
