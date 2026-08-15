package schedule

import (
	"context"
	"eGZ-Board/internal/dao"
	"eGZ-Board/utility"

	"eGZ-Board/api/schedule/v1"
)

func (c *ControllerV1) GetSchedule(ctx context.Context, req *v1.GetScheduleReq) (res *v1.GetScheduleRes, err error) {
	classID, err := utility.GetClassID(ctx)
	if err != nil {
		return nil, err
	}
	scheduleTable, err := dao.GetSchedule(classID)
	if err != nil {
		return nil, err
	}
	return &v1.GetScheduleRes{Res: scheduleTable}, nil
}
