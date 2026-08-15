package setting

import (
	"context"
	"eGZ-Board/internal/dao"
	"eGZ-Board/utility"

	"eGZ-Board/api/setting/v1"
)

func (c *ControllerV1) SetSettingValue(ctx context.Context, req *v1.SetSettingValueReq) (res *v1.SetSettingValueRes, err error) {
	classID, err := utility.GetClassID(ctx)
	if err != nil {
		return nil, err
	}
	if err := dao.SetSetting(req.Key, req.Value, classID); err != nil {
		return nil, err
	}
	res = &v1.SetSettingValueRes{
		Res: "Success",
	}
	return
}
