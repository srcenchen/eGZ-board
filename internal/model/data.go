package model

import (
	"eGZ-Board/internal/model/entity"
	"fmt"
	"strings"
	"sync"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	dbMu sync.RWMutex
)

func GetDatabase() *gorm.DB {
	dbMu.RLock()
	defer dbMu.RUnlock()
	return db
}

// SetDatabase replaces the active database. It is primarily useful for isolated tests.
func SetDatabase(database *gorm.DB) {
	dbMu.Lock()
	db = database
	dbMu.Unlock()
}

func InitData(path string) error {
	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	sqlDB, err := database.DB()
	if err != nil {
		return fmt.Errorf("get database connection: %w", err)
	}
	// SQLite PRAGMAs are connection-local; one connection keeps them consistent.
	sqlDB.SetMaxOpenConns(1)
	SetDatabase(database)

	for _, pragma := range []string{
		"PRAGMA busy_timeout = 5000",
		"PRAGMA foreign_keys = ON",
		"PRAGMA journal_mode = WAL",
	} {
		if err := database.Exec(pragma).Error; err != nil {
			return fmt.Errorf("configure sqlite (%s): %w", pragma, err)
		}
	}

	if err := database.AutoMigrate(&entity.SettingTable{}, &entity.UserTable{}, &entity.CountDownTable{},
		&entity.ScheduleTable{}, &entity.WeatherTable{}, &entity.ClassTable{}, &entity.DeviceTable{}); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}
	if err := deduplicateAndIndex(database); err != nil {
		return err
	}

	var userTable entity.UserTable
	result := database.Where("id = ?", 1).First(&userTable)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return fmt.Errorf("check seed data: %w", result.Error)
	}
	if result.Error == gorm.ErrRecordNotFound {
		if err := database.Create(&entity.UserTable{Username: "admin", Password: "8dbc6dfc58f9f7cf07eff4bef62c0468", Role: "SuperAdmin"}).Error; err != nil {
			return fmt.Errorf("seed administrator: %w", err)
		}
		for grade := 1; grade <= 3; grade++ {
			for class := 1; class <= 24; class++ {
				classTable := entity.ClassTable{Grade: grade, Class: class}
				if err := database.Create(&classTable).Error; err != nil {
					return fmt.Errorf("seed class %d-%d: %w", grade, class, err)
				}
				settings := []entity.SettingTable{
					{Key: "slogan_time", Value: "60", Class: classTable.Id},
					{Key: "wall_url", Value: "https://bing.img.run/1920x1080.php", Class: classTable.Id},
					{Key: "notice", Value: "<h2 style=\"text-align: center;\">本周任务</h2><p>注意：字体大小不允许超过H2</p>", Class: classTable.Id},
					{Key: "notice_switch", Value: "false", Class: classTable.Id},
				}
				if err := database.Create(&settings).Error; err != nil {
					return fmt.Errorf("seed settings for class %d: %w", classTable.Id, err)
				}
				if err := database.Create(&entity.CountDownTable{Event: "高考", Date: "2025-6-7", Class: classTable.Id, Type: "count"}).Error; err != nil {
					return fmt.Errorf("seed countdown for class %d: %w", classTable.Id, err)
				}
			}
		}
	}
	return nil
}

func deduplicateAndIndex(database *gorm.DB) error {
	return database.Transaction(func(tx *gorm.DB) error {
		// Merge duplicate class IDs before enforcing the natural class identity.
		var duplicates []entity.ClassTable
		if err := tx.Raw(`SELECT c.id, c.grade, c.class FROM class_tables c
			JOIN (SELECT grade, class, MIN(id) keep_id FROM class_tables GROUP BY grade, class HAVING COUNT(*) > 1) d
			ON c.grade = d.grade AND c.class = d.class WHERE c.id <> d.keep_id`).Scan(&duplicates).Error; err != nil {
			return fmt.Errorf("find duplicate classes: %w", err)
		}
		for _, duplicate := range duplicates {
			var keepID int
			if err := tx.Raw("SELECT MIN(id) FROM class_tables WHERE grade = ? AND class = ?", duplicate.Grade, duplicate.Class).Scan(&keepID).Error; err != nil {
				return fmt.Errorf("resolve duplicate class: %w", err)
			}
			for _, update := range []string{
				"UPDATE device_tables SET class_id = ? WHERE class_id = ?",
				"UPDATE setting_tables SET class = ? WHERE class = ?",
				"UPDATE count_down_tables SET class = ? WHERE class = ?",
				"UPDATE schedule_tables SET class = ? WHERE class = ?",
			} {
				if err := tx.Exec(update, keepID, duplicate.Id).Error; err != nil {
					return fmt.Errorf("merge duplicate class references: %w", err)
				}
			}
			if err := tx.Delete(&entity.ClassTable{}, duplicate.Id).Error; err != nil {
				return fmt.Errorf("remove duplicate class: %w", err)
			}
		}

		for _, statement := range []string{
			"DELETE FROM device_tables WHERE id NOT IN (SELECT MAX(id) FROM device_tables GROUP BY device_id)",
			"DELETE FROM setting_tables WHERE id NOT IN (SELECT MAX(id) FROM setting_tables GROUP BY key, class)",
			"DELETE FROM schedule_tables WHERE id NOT IN (SELECT MAX(id) FROM schedule_tables GROUP BY class)",
			"DELETE FROM weather_tables WHERE id NOT IN (SELECT MAX(id) FROM weather_tables GROUP BY key)",
		} {
			if err := tx.Exec(statement).Error; err != nil {
				return fmt.Errorf("deduplicate database: %w", err)
			}
		}
		indexes := []string{
			"CREATE UNIQUE INDEX IF NOT EXISTS idx_class_grade_class ON class_tables(grade, class)",
			"CREATE UNIQUE INDEX IF NOT EXISTS idx_device_id ON device_tables(device_id)",
			"CREATE UNIQUE INDEX IF NOT EXISTS idx_setting_key_class ON setting_tables(key, class)",
			"CREATE UNIQUE INDEX IF NOT EXISTS idx_schedule_class ON schedule_tables(class)",
			"CREATE UNIQUE INDEX IF NOT EXISTS idx_weather_key ON weather_tables(key)",
		}
		for _, statement := range indexes {
			if err := tx.Exec(statement).Error; err != nil {
				return fmt.Errorf("create unique index %s: %w", strings.Fields(statement)[6], err)
			}
		}
		return nil
	})
}
