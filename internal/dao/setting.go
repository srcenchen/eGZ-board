package dao

import (
	"eGZ-Board/internal/model"
	"eGZ-Board/internal/model/entity"

	"gorm.io/gorm/clause"
)

// GetSetting 获取设置
func GetSetting(key string, class int) (entity.SettingTable, error) {
	var settingTable entity.SettingTable
	err := model.GetDatabase().Where("key = ? and class = ?", key, class).First(&settingTable).Error
	return settingTable, err
}

// SetSetting Set设置
func SetSetting(key string, value string, class int) error {
	setting := entity.SettingTable{Key: key, Value: value, Class: class}
	return model.GetDatabase().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}, {Name: "class"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&setting).Error
}
