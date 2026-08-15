package count_down

import (
	"context"
	"eGZ-Board/internal/dao"
	"eGZ-Board/utility"

	"eGZ-Board/api/count_down/v1"
)

func (c *ControllerV1) GetEvent(ctx context.Context, req *v1.GetEventReq) (res *v1.GetEventRes, err error) {
	classID, err := utility.GetClassID(ctx)
	if err != nil {
		return nil, err
	}
	events, err := dao.GetCountDown(classID)
	if err != nil {
		return nil, err
	}
	return &v1.GetEventRes{Res: events}, nil
}
