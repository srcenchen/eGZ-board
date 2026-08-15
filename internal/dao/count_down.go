package dao

import (
	"eGZ-Board/internal/model"
	"eGZ-Board/internal/model/entity"
)

// GetCountDown 获取倒计时
func GetCountDown(class int) ([]entity.CountDownTable, error) {
	var countDownTable []entity.CountDownTable
	err := model.GetDatabase().Where("class = ?", class).Find(&countDownTable).Error
	return countDownTable, err
}

// AddCountDown 添加倒计时
func AddCountDown(event string, date string, countType string, during int, class int) error {
	return model.GetDatabase().Create(&entity.CountDownTable{Event: event, Date: date, Type: countType, During: during, Class: class}).Error
}

func DeleteCountDown(id, class int) error {
	return model.GetDatabase().Where("id = ? AND class = ?", id, class).Delete(&entity.CountDownTable{}).Error
}
