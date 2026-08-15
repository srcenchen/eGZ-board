package device

import (
	"context"
	"eGZ-Board/internal/dao"
	"github.com/gogf/gf/v2/net/ghttp"

	"eGZ-Board/api/device/v1"
)

func (c *ControllerV1) BindDevice(ctx context.Context, req *v1.BindDeviceReq) (res *v1.BindDeviceRes, err error) {
	classTable, err := dao.GetClass(req.Grade, req.Class)
	if err != nil {
		return nil, err
	}
	r := ghttp.RequestFromCtx(ctx)
	clientIP := r.GetClientIp()
	if _, err = dao.UpsertDevice(clientIP, classTable.Id); err != nil {
		return nil, err
	}
	return &v1.BindDeviceRes{Msg: "Success"}, nil
}
