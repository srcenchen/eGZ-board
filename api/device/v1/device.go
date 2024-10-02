package v1

import (
	"eGZ-Board/internal/model/entity"
	"github.com/gogf/gf/v2/frame/g"
)

type IsBindReq struct {
	g.Meta `method:"get"`
}

type IsBindRes struct {
	g.Meta `method:"get"`
	Device entity.DeviceTable `json:"isBind"`
}

type BindDeviceReq struct {
	g.Meta `method:"post"`
	Grade  int    `json:"grade"`
	Class  string `json:"class"`
}

type BindDeviceRes struct {
	g.Meta `method:"post"`
	Msg    string `json:"msg"`
}

type GetClassReq struct {
	g.Meta `method:"get"`
}

type GetClassRes struct {
	g.Meta `method:"get"`
	Grade  int `json:"grade"`
	Class  int `json:"class"`
}

type UnBindReq struct {
	g.Meta `method:"get"`
}

type UnBindRes struct {
	g.Meta `method:"get"`
}
