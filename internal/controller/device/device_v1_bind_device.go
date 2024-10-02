package device

import (
	"context"
	"eGZ-Board/internal/model"
	"eGZ-Board/internal/model/entity"
	"github.com/gogf/gf/v2/net/ghttp"

	"eGZ-Board/api/device/v1"
)

func (c *ControllerV1) BindDevice(ctx context.Context, req *v1.BindDeviceReq) (res *v1.BindDeviceRes, err error) {
	// 获取ClassID
	var classTable entity.ClassTable
	// 获取访问者 IP 地址
	r := ghttp.RequestFromCtx(ctx)
	clientIP := r.GetClientIp()
	model.GetDatabase().Where("grade = ? and class=?", req.Grade, req.Class).First(&classTable)
	model.GetDatabase().Create(&entity.DeviceTable{DeviceID: clientIP, ClassID: classTable.Id})
	return
}
