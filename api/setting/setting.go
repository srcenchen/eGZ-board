// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package setting

import (
	"context"

	"eGZ-Board/api/setting/v1"
)

type ISettingV1 interface {
	GetSettingValue(ctx context.Context, req *v1.GetSettingValueReq) (res *v1.GetSettingValueRes, err error)
	SetSettingValue(ctx context.Context, req *v1.SetSettingValueReq) (res *v1.SetSettingValueRes, err error)
	UploadWallPaper(ctx context.Context, req *v1.UploadWallPaperReq) (res *v1.UploadWallPaperRes, err error)
}
