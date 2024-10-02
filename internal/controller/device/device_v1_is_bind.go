package device

import (
	"context"
	"eGZ-Board/internal/model"
	"eGZ-Board/internal/model/entity"
	"github.com/gogf/gf/v2/net/ghttp"

	"eGZ-Board/api/device/v1"
)

func (c *ControllerV1) IsBind(ctx context.Context, req *v1.IsBindReq) (res *v1.IsBindRes, err error) {
	// 从 context 中获取 *ghttp.Request
	r := ghttp.RequestFromCtx(ctx)

	// 获取访问者 IP 地址
	clientIP := r.GetClientIp()
	var deviceTable entity.DeviceTable
	model.GetDatabase().Where("device_ID = ?", clientIP).First(&deviceTable)
	res = &v1.IsBindRes{
		Device: deviceTable,
	}
	return
}
