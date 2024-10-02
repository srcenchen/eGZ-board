package utility

import (
	"eGZ-Board/internal/model"
	"eGZ-Board/internal/model/entity"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func GetClassID(ctx g.Ctx) int {
	// 从 context 中获取 *ghttp.Request
	r := ghttp.RequestFromCtx(ctx)

	// 获取访问者 IP 地址
	clientIP := r.GetClientIp()
	var deviceTable entity.DeviceTable
	model.GetDatabase().Where("device_ID = ?", clientIP).First(&deviceTable)
	return deviceTable.ClassID
}
