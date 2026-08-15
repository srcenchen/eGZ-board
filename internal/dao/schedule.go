package dao

import (
	"eGZ-Board/internal/model"
	"eGZ-Board/internal/model/entity"

	"gorm.io/gorm/clause"
)

func GetSchedule(classID int) (entity.ScheduleTable, error) {
	var schedule entity.ScheduleTable
	if err := model.GetDatabase().Where("class = ?", classID).First(&schedule).Error; err != nil {
		return schedule, err
	}
	return schedule, nil
}

func UpsertSchedule(classID int, value string) error {
	schedule := entity.ScheduleTable{Class: classID, Schedule: value}
	return model.GetDatabase().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "class"}},
		DoUpdates: clause.AssignmentColumns([]string{"schedule"}),
	}).Create(&schedule).Error
}

func DeleteSchedule(classID int) error {
	return model.GetDatabase().Where("class = ?", classID).Delete(&entity.ScheduleTable{}).Error
}
