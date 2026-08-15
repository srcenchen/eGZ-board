package main

import (
	"eGZ-Board/internal/model"
	_ "eGZ-Board/internal/packed"
	"eGZ-Board/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"log"
	"os"

	"github.com/gogf/gf/v2/os/gctx"

	"eGZ-Board/internal/cmd"
)

func main() {
	// 设置默认配置文件目录
	_ = g.Cfg().GetAdapter().(*gcfg.AdapterFile).AddPath("manifest/conf")
	_ = g.Cfg().GetAdapter().(*gcfg.AdapterFile).AddPath("conf")
	if err := initPath(); err != nil {
		log.Fatalf("initialize resource paths: %v", err)
	}
	if err := model.InitData("./resource/database/data.db"); err != nil {
		log.Fatalf("initialize database: %v", err)
	}
	// logic.SettingInit() // 初始化设置
	go initProxyData()
	if err := service.Timer(); err != nil {
		log.Fatalf("initialize weather timer: %v", err)
	}
	cmd.Main.Run(gctx.GetInitCtx())
}

// createPath /**
func createPath(path string) error {
	return os.MkdirAll(path, 0755)
}

// initPath /**
func initPath() error {
	for _, path := range []string{"./resource/database", "./resource/upload", "./resource/upload/images", "./resource/upload/videos"} {
		if err := createPath(path); err != nil {
			return err
		}
	}
	return nil
}

func initProxyData() {
	if err := service.WeatherProxy(); err != nil {
		log.Printf("initial weather refresh failed: %v", err)
	}
}
