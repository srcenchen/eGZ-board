package dao

import (
	"eGZ-Board/internal/model"
	"eGZ-Board/internal/model/entity"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetDevice(deviceID string) (entity.DeviceTable, error) {
	var device entity.DeviceTable
	if err := model.GetDatabase().Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		return device, err
	}
	return device, nil
}

func GetClass(grade, class int) (entity.ClassTable, error) {
	var classTable entity.ClassTable
	if err := model.GetDatabase().Where("grade = ? AND class = ?", grade, class).First(&classTable).Error; err != nil {
		return classTable, err
	}
	return classTable, nil
}

func GetClassByDevice(deviceID string) (entity.ClassTable, error) {
	device, err := GetDevice(deviceID)
	if err != nil {
		return entity.ClassTable{}, err
	}
	var classTable entity.ClassTable
	if err := model.GetDatabase().First(&classTable, device.ClassID).Error; err != nil {
		return classTable, err
	}
	return classTable, nil
}

func UpsertDevice(deviceID string, classID int) (entity.DeviceTable, error) {
	device := entity.DeviceTable{DeviceID: deviceID, ClassID: classID}
	err := model.GetDatabase().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "device_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"class_id"}),
	}).Create(&device).Error
	if err != nil {
		return device, fmt.Errorf("upsert device: %w", err)
	}
	return GetDevice(deviceID)
}

func DeleteDevice(deviceID string) error {
	return model.GetDatabase().Where("device_id = ?", deviceID).Delete(&entity.DeviceTable{}).Error
}

func IsNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}
