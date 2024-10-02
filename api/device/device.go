// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package device

import (
	"context"

	"eGZ-Board/api/device/v1"
)

type IDeviceV1 interface {
	IsBind(ctx context.Context, req *v1.IsBindReq) (res *v1.IsBindRes, err error)
	BindDevice(ctx context.Context, req *v1.BindDeviceReq) (res *v1.BindDeviceRes, err error)
	GetClass(ctx context.Context, req *v1.GetClassReq) (res *v1.GetClassRes, err error)
	UnBind(ctx context.Context, req *v1.UnBindReq) (res *v1.UnBindRes, err error)
}
