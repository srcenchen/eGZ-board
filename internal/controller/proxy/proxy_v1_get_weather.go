package proxy

import (
	"context"
	"eGZ-Board/internal/dao"

	"eGZ-Board/api/proxy/v1"
)

func (c *ControllerV1) GetWeather(ctx context.Context, req *v1.GetWeatherReq) (res *v1.GetWeatherRes, err error) {
	res = &v1.GetWeatherRes{
		Data: dao.GetWeather(req.Key),
	}
	return
}
