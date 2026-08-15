package count_down

import (
	"context"
	"eGZ-Board/internal/dao"
	"eGZ-Board/utility"

	"eGZ-Board/api/count_down/v1"
)

func (c *ControllerV1) DelEvent(ctx context.Context, req *v1.DelEventReq) (res *v1.DelEventRes, err error) {
	classID, err := utility.GetClassID(ctx)
	if err != nil {
		return nil, err
	}
	if err := dao.DeleteCountDown(req.Id, classID); err != nil {
		return nil, err
	}
	return &v1.DelEventRes{}, nil
}
