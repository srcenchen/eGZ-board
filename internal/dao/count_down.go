package dao

import (
	"eGZ-Board/internal/model"
	"eGZ-Board/internal/model/entity"
)

// GetCountDown 获取倒计时
func GetCountDown(class int) []entity.CountDownTable {
	var countDownTable []entity.CountDownTable
	model.GetDatabase().Where("class = ?", class).Find(&countDownTable)
	return countDownTable
}

// AddCountDown 添加倒计时
func AddCountDown(event string, date string, countType string, during int, class int) {
	model.GetDatabase().Create(&entity.CountDownTable{Event: event, Date: date, Type: countType, During: during, Class: class})
}
