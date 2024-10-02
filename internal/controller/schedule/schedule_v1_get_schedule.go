package schedule

import (
	"context"
	"eGZ-Board/internal/model"
	"eGZ-Board/internal/model/entity"
	"eGZ-Board/utility"

	"eGZ-Board/api/schedule/v1"
)

func (c *ControllerV1) GetSchedule(ctx context.Context, req *v1.GetScheduleReq) (res *v1.GetScheduleRes, err error) {
	var scheduleTable entity.ScheduleTable
	model.GetDatabase().Where("class = ?", utility.GetClassID(ctx)).First(&scheduleTable)
	res = &v1.GetScheduleRes{
		Res: scheduleTable,
	}
	return
}
