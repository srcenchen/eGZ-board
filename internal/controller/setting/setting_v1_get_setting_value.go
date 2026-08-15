package setting

import (
	"context"
	"eGZ-Board/internal/dao"
	"eGZ-Board/utility"

	"eGZ-Board/api/setting/v1"
)

func (c *ControllerV1) GetSettingValue(ctx context.Context, req *v1.GetSettingValueReq) (res *v1.GetSettingValueRes, err error) {
	classID, err := utility.GetClassID(ctx)
	if err != nil {
		return nil, err
	}
	setting, err := dao.GetSetting(req.Key, classID)
	if err != nil {
		return nil, err
	}
	return &v1.GetSettingValueRes{Res: setting}, nil
}
