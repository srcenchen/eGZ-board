package device

import (
	"context"
	"eGZ-Board/internal/model"
	"eGZ-Board/internal/model/entity"
	"github.com/gogf/gf/v2/net/ghttp"

	"eGZ-Board/api/device/v1"
)

func (c *ControllerV1) UnBind(ctx context.Context, req *v1.UnBindReq) (res *v1.UnBindRes, err error) {
	// 从 context 中获取 *ghttp.Request
	r := ghttp.RequestFromCtx(ctx)

	// 获取访问者 IP 地址
	clientIP := r.GetClientIp()
	model.GetDatabase().Delete(&entity.DeviceTable{}, "device_ID = ?", clientIP)
	return
}
