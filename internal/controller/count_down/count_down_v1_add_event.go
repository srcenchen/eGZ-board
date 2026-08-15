package count_down

import (
	"context"
	"eGZ-Board/internal/dao"
	"eGZ-Board/utility"

	"eGZ-Board/api/count_down/v1"
)

func (c *ControllerV1) AddEvent(ctx context.Context, req *v1.AddEventReq) (res *v1.AddEventRes, err error) {
	classID, err := utility.GetClassID(ctx)
	if err != nil {
		return nil, err
	}
	if err := dao.AddCountDown(req.Event, req.Date, req.Type, req.During, classID); err != nil {
		return nil, err
	}
	res = &v1.AddEventRes{}
	return
}
