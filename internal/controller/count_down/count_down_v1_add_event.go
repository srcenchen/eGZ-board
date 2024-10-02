package count_down

import (
	"context"
	"eGZ-Board/internal/dao"
	"eGZ-Board/utility"

	"eGZ-Board/api/count_down/v1"
)

func (c *ControllerV1) AddEvent(ctx context.Context, req *v1.AddEventReq) (res *v1.AddEventRes, err error) {
	dao.AddCountDown(req.Event, req.Date, req.Type, req.During, utility.GetClassID(ctx))
	res = &v1.AddEventRes{}
	return
}
