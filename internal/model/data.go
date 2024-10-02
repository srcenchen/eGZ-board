package model

import (
	"eGZ-Board/internal/model/entity"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"sync"
)

var once sync.Once
var db *gorm.DB

func GetDatabase() *gorm.DB {
	once.Do(func() {
		db, _ = gorm.Open(sqlite.Open("./resource/database/data.db"), &gorm.Config{})
	})
	return db
}

// InitData /**
func InitData() {
	// 自动迁移模式
	if GetDatabase().AutoMigrate(&entity.SettingTable{}, &entity.UserTable{}, &entity.CountDownTable{},
		&entity.ScheduleTable{}, &entity.WeatherTable{}, &entity.ClassTable{}, entity.DeviceTable{}) != nil {
		panic("failed to migrate database")
	}
	// 判断是否已经初始化过
	var userTable entity.UserTable
	if GetDatabase().Where("id = ?", 1).First(&userTable).RowsAffected == 0 { // 没初始化过
		// 初始化SuperAdmin
		GetDatabase().Create(&entity.UserTable{Username: "admin", Password: "8dbc6dfc58f9f7cf07eff4bef62c0468", Role: "SuperAdmin"})
		// 创建班级
		for grade := 1; grade <= 3; grade++ {
			// 遍历每个年级的1班到24班
			for class := 1; class <= 24; class++ {
				// 输出班级名称
				// fmt.Printf("高%d %d班\n", grade, class)
				classTable := entity.ClassTable{Grade: grade, Class: class}
				GetDatabase().Create(&classTable)
				GetDatabase().Create(&entity.SettingTable{Key: "slogan_time", Value: "60", Class: classTable.Id})
				GetDatabase().Create(&entity.SettingTable{Key: "wall_url", Value: "https://bing.img.run/1920x1080.php", Class: classTable.Id})
				GetDatabase().Create(&entity.SettingTable{Key: "notice", Value: "<h2 style=\"text-align: center;\">本周任务</h2><p>注意：字体大小不允许超过H2</p>", Class: classTable.Id})
				GetDatabase().Create(&entity.SettingTable{Key: "notice_switch", Value: "false", Class: classTable.Id})
				GetDatabase().Create(&entity.CountDownTable{Event: "高考", Date: "2025-6-6", Class: classTable.Id, Type: "count"})
			}
		}
	}
}
