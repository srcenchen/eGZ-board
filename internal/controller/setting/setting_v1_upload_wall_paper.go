package setting

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"path"
	"strings"
	"time"

	"eGZ-Board/api/setting/v1"
)

func (c *ControllerV1) UploadWallPaper(ctx context.Context, req *v1.UploadWallPaperReq) (res *v1.UploadWallPaperRes, err error) {
	_ = req
	_ = ctx
	file := req.File                               // 获取文件
	file.Filename = strings.ToLower(file.Filename) // 降为小写
	timeStamp := gconv.String(time.Now().Unix())   // 获取时间戳
	randomFileName := grand.Letters(16)            // 生成16位随机字符串
	file.Filename = timeStamp + randomFileName + path.Ext(file.Filename)
	fileName, err := file.Save("./resource/upload")
	res = &v1.UploadWallPaperRes{
		FileName: fileName,
	}
	return
}
