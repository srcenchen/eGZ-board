package utility

import (
	"eGZ-Board/internal/dao"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func GetClassID(ctx g.Ctx) (int, error) {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return 0, fmt.Errorf("request context is unavailable")
	}
	clientIP := r.GetClientIp()
	device, err := dao.GetDevice(clientIP)
	if err != nil {
		if dao.IsNotFound(err) {
			return 0, fmt.Errorf("device %s is not bound to a class", clientIP)
		}
		return 0, fmt.Errorf("resolve device class: %w", err)
	}
	if device.ClassID <= 0 {
		return 0, fmt.Errorf("device %s has an invalid class binding", clientIP)
	}
	return device.ClassID, nil
}
