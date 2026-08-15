package proxy

import (
	"context"
	"eGZ-Board/internal/dao"

	"eGZ-Board/api/proxy/v1"
)

func (c *ControllerV1) GetWeather(ctx context.Context, req *v1.GetWeatherReq) (res *v1.GetWeatherRes, err error) {
	weather, err := dao.GetWeather(req.Key)
	if err != nil {
		return nil, err
	}
	return &v1.GetWeatherRes{Data: weather}, nil
}
