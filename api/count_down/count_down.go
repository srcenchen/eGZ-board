// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package count_down

import (
	"context"

	"eGZ-Board/api/count_down/v1"
)

type ICountDownV1 interface {
	AddEvent(ctx context.Context, req *v1.AddEventReq) (res *v1.AddEventRes, err error)
	DelEvent(ctx context.Context, req *v1.DelEventReq) (res *v1.DelEventRes, err error)
	GetEvent(ctx context.Context, req *v1.GetEventReq) (res *v1.GetEventRes, err error)
}
