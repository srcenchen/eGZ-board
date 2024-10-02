package setting

import (
	"context"
	"eGZ-Board/internal/dao"
	"eGZ-Board/utility"

	"eGZ-Board/api/setting/v1"
)

func (c *ControllerV1) GetSettingValue(ctx context.Context, req *v1.GetSettingValueReq) (res *v1.GetSettingValueRes, err error) {
	res = &v1.GetSettingValueRes{Res: dao.GetSetting(req.Key, utility.GetClassID(ctx))}
	return
}
