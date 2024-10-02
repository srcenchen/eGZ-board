package count_down

import (
	"context"
	"eGZ-Board/internal/model"
	"eGZ-Board/internal/model/entity"

	"eGZ-Board/api/count_down/v1"
)

func (c *ControllerV1) DelEvent(ctx context.Context, req *v1.DelEventReq) (res *v1.DelEventRes, err error) {
	model.GetDatabase().Delete(&entity.CountDownTable{}, "id = ?", req.Id)
	res = &v1.DelEventRes{}
	return
}
