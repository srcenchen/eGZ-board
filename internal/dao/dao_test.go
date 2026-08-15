package dao

import (
	"eGZ-Board/internal/model"
	"eGZ-Board/internal/model/entity"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDatabase(t *testing.T) *gorm.DB {
	t.Helper()
	database, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := database.AutoMigrate(&entity.DeviceTable{}, &entity.SettingTable{}, &entity.CountDownTable{}); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}
	for _, statement := range []string{
		"CREATE UNIQUE INDEX idx_device_id ON device_tables(device_id)",
		"CREATE UNIQUE INDEX idx_setting_key_class ON setting_tables(key, class)",
	} {
		if err := database.Exec(statement).Error; err != nil {
			t.Fatalf("create test index: %v", err)
		}
	}
	model.SetDatabase(database)
	t.Cleanup(func() {
		sqlDB, err := database.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
		model.SetDatabase(nil)
	})
	return database
}

func TestUpsertDevice(t *testing.T) {
	database := setupTestDatabase(t)

	if _, err := UpsertDevice("192.0.2.1", 10); err != nil {
		t.Fatalf("insert device: %v", err)
	}
	device, err := UpsertDevice("192.0.2.1", 20)
	if err != nil {
		t.Fatalf("update device: %v", err)
	}
	if device.ClassID != 20 {
		t.Fatalf("class ID = %d, want 20", device.ClassID)
	}
	var count int64
	if err := database.Model(&entity.DeviceTable{}).Count(&count).Error; err != nil {
		t.Fatalf("count devices: %v", err)
	}
	if count != 1 {
		t.Fatalf("device count = %d, want 1", count)
	}
}

func TestSetSettingUpsertsByClass(t *testing.T) {
	database := setupTestDatabase(t)

	if err := SetSetting("notice", "class one", 1); err != nil {
		t.Fatalf("set first setting: %v", err)
	}
	if err := SetSetting("notice", "updated", 1); err != nil {
		t.Fatalf("update first setting: %v", err)
	}
	if err := SetSetting("notice", "class two", 2); err != nil {
		t.Fatalf("set second setting: %v", err)
	}
	setting, err := GetSetting("notice", 1)
	if err != nil {
		t.Fatalf("get setting: %v", err)
	}
	if setting.Value != "updated" {
		t.Fatalf("setting value = %q, want updated", setting.Value)
	}
	var count int64
	if err := database.Model(&entity.SettingTable{}).Count(&count).Error; err != nil {
		t.Fatalf("count settings: %v", err)
	}
	if count != 2 {
		t.Fatalf("setting count = %d, want 2", count)
	}
}

func TestDeleteCountDownIsClassScoped(t *testing.T) {
	database := setupTestDatabase(t)
	event := entity.CountDownTable{Event: "exam", Date: "2026-06-07", Type: "count", Class: 1}
	if err := database.Create(&event).Error; err != nil {
		t.Fatalf("create event: %v", err)
	}

	if err := DeleteCountDown(event.Id, 2); err != nil {
		t.Fatalf("delete from wrong class: %v", err)
	}
	if err := database.First(&entity.CountDownTable{}, event.Id).Error; err != nil {
		t.Fatalf("event was deleted by another class: %v", err)
	}
	if err := DeleteCountDown(event.Id, 1); err != nil {
		t.Fatalf("delete from owning class: %v", err)
	}
	if err := database.First(&entity.CountDownTable{}, event.Id).Error; err != gorm.ErrRecordNotFound {
		t.Fatalf("lookup after delete error = %v, want record not found", err)
	}
}
