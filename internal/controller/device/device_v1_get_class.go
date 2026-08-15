package device

import (
	"context"
	"eGZ-Board/internal/dao"
	"github.com/gogf/gf/v2/net/ghttp"

	"eGZ-Board/api/device/v1"
)

func (c *ControllerV1) GetClass(ctx context.Context, req *v1.GetClassReq) (res *v1.GetClassRes, err error) {
	// 从 context 中获取 *ghttp.Request
	r := ghttp.RequestFromCtx(ctx)

	// 获取访问者 IP 地址
	clientIP := r.GetClientIp()
	classTable, err := dao.GetClassByDevice(clientIP)
	if err != nil {
		return nil, err
	}
	res = &v1.GetClassRes{
		Grade: classTable.Grade,
		Class: classTable.Class,
	}
	return
}
