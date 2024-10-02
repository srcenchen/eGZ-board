package v1

import (
	"eGZ-Board/internal/model/entity"
	"github.com/gogf/gf/v2/frame/g"
)

type GetSettingValueReq struct {
	g.Meta `tags:"Hello" method:"get"`
	Key    string `v:"required"`
}

type GetSettingValueRes struct {
	g.Meta `tags:"Hello" method:"get"`
	Res    entity.SettingTable
}

type SetSettingValueReq struct {
	g.Meta `tags:"Hello" method:"post"`
	Value  string `v:"required"`
	Key    string `v:"required"`
}

type SetSettingValueRes struct {
	g.Meta `tags:"Hello" method:"post"`
	Res    string
}
