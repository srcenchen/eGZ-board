package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// UploadScheduleReq 上传课程表文件请求
type UploadScheduleReq struct {
	g.Meta `method:"POST" summary:"上传课程表文件" tags:"课程表资源"`
	File   *ghttp.UploadFile `json:"file" v:"required#文件不能为空"`
}

// UploadScheduleRes 上传课程表文件响应
type UploadScheduleRes struct {
	FileName string `json:"fileName" example:"test.mp3"`
}
