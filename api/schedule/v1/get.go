package v1

import (
	"eGZ-Board/internal/model/entity"
	"github.com/gogf/gf/v2/frame/g"
)

type GetScheduleReq struct {
	g.Meta `method:"get"`
}

type GetScheduleRes struct {
	g.Meta `method:"get"`
	Res    entity.ScheduleTable
}
