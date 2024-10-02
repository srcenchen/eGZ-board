package v1

import (
	"eGZ-Board/internal/model/entity"
	"github.com/gogf/gf/v2/frame/g"
)

type GetWeatherReq struct {
	g.Meta `tags:"Hello" method:"get"`
	Key    string `v:"required"`
}

type GetWeatherRes struct {
	g.Meta `tags:"Hello" method:"get"`
	Data   entity.WeatherTable
}
