// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package schedule

import (
	"context"

	"eGZ-Board/api/schedule/v1"
)

type IScheduleV1 interface {
	GetSchedule(ctx context.Context, req *v1.GetScheduleReq) (res *v1.GetScheduleRes, err error)
	UploadSchedule(ctx context.Context, req *v1.UploadScheduleReq) (res *v1.UploadScheduleRes, err error)
}
