package main

import (
	"eGZ-Board/internal/model"
	_ "eGZ-Board/internal/packed"
	"eGZ-Board/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"os"

	"github.com/gogf/gf/v2/os/gctx"

	"eGZ-Board/internal/cmd"
)

func main() {
	// 设置默认配置文件目录
	_ = g.Cfg().GetAdapter().(*gcfg.AdapterFile).AddPath("manifest/conf")
	_ = g.Cfg().GetAdapter().(*gcfg.AdapterFile).AddPath("conf")
	initPath()       // 初始化文件夹
	model.InitData() // 初始化数据库
	// logic.SettingInit() // 初始化设置
	go initProxyData()
	service.Timer() // 代理计时器
	cmd.Main.Run(gctx.GetInitCtx())
}

// createPath /**
func createPath(path string) {
	// 检测是否存在 path 文件夹 如果不存在则创建
	if _, err := os.Stat(path); os.IsNotExist(err) {
		errs := os.Mkdir(path, os.ModePerm)
		if errs != nil {
			return
		}
	}
}

// initPath /**
func initPath() {
	// 创建文件夹
	createPath("./resource")
	createPath("./resource/database")
	createPath("./resource/upload")
	createPath("./resource/upload/images")
	createPath("./resource/upload/videos")
}

func initProxyData() {
	service.WeatherProxy()
}
