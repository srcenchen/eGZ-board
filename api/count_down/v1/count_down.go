package v1

import (
	"eGZ-Board/internal/model/entity"
	"github.com/gogf/gf/v2/frame/g"
)

type AddEventReq struct {
	g.Meta `method:"post"`
	Event  string `v:"required"`
	Date   string `v:"required"`
	Type   string `v:"required"`
	During int    `v:"required"`
}

type AddEventRes struct {
	g.Meta `method:"post"`
}

type DelEventReq struct {
	g.Meta `method:"post"`
	Id     int `v:"required"`
}

type DelEventRes struct {
	g.Meta `method:"post"`
}

type GetEventReq struct {
	g.Meta `method:"get"`
}

type GetEventRes struct {
	g.Meta `method:"get"`
	Res    []entity.CountDownTable
}
