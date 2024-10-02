// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package proxy

import (
	"context"

	"eGZ-Board/api/proxy/v1"
)

type IProxyV1 interface {
	GetWeather(ctx context.Context, req *v1.GetWeatherReq) (res *v1.GetWeatherRes, err error)
}
