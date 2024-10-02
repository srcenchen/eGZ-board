package dao

import (
	"eGZ-Board/internal/model"
	"eGZ-Board/internal/model/entity"
)

// GetSetting 获取设置
func GetSetting(key string, class int) entity.SettingTable {
	var settingTable entity.SettingTable
	model.GetDatabase().Where("key = ? and class = ?", key, class).First(&settingTable)
	return settingTable
}

// SetSetting Set设置
func SetSetting(key string, value string, class int) {
	var settingTable entity.SettingTable
	model.GetDatabase().Where("key = ? and class = ?", key, class).First(&settingTable)
	settingTable.Value = value
	model.GetDatabase().Save(&settingTable)
}
