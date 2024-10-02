package schedule

import (
	"context"
	"eGZ-Board/api/schedule/v1"
	"eGZ-Board/internal/model"
	"eGZ-Board/internal/model/entity"
	"eGZ-Board/utility"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func (c *ControllerV1) UploadSchedule(ctx context.Context, req *v1.UploadScheduleReq) (res *v1.UploadScheduleRes, err error) {
	_ = req
	_ = ctx
	file := req.File                               // 获取文件
	file.Filename = strings.ToLower(file.Filename) // 降为小写
	timeStamp := gconv.String(time.Now().Unix())   // 获取时间戳
	randomFileName := grand.Letters(16)            // 生成16位随机字符串
	file.Filename = timeStamp + randomFileName + path.Ext(file.Filename)
	fileName, err := file.Save("./resource/upload")
	// 打开 Excel 文件
	f, err := excelize.OpenFile("./resource/upload/" + fileName)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}

	// 获取第一个工作表的名称
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		log.Fatalf("Sheet name is empty")
	}

	// 获取工作表的所有列数据
	cols, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatalf("Failed to get columns: %v", err)
	}

	// 定义数据结构
	data := map[string][]string{}

	// 跳过标题行，从第二行开始读取数据
	for i, col := range cols {
		if i == 0 {
			// 跳过标题行
			continue
		}

		for j := 1; j < len(col); j++ {
			data[strconv.Itoa(j)] = append(data[strconv.Itoa(j)], col[j])
		}
	}

	// 将数据转换为 JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// 打印 JSON 数据
	fmt.Println(string(jsonData))
	// 判断是否存在课程表
	var scheduleTable entity.ScheduleTable
	if model.GetDatabase().Where("class = ?", utility.GetClassID(ctx)).First(&scheduleTable).RowsAffected == 0 {
		model.GetDatabase().Create(&entity.ScheduleTable{Schedule: string(jsonData), Class: utility.GetClassID(ctx)})
	} else {
		model.GetDatabase().Where("class = ?", utility.GetClassID(ctx)).First(&scheduleTable)
		scheduleTable.Schedule = string(jsonData)
		model.GetDatabase().Save(&scheduleTable)
	}
	err = os.Remove("./resource/upload/" + fileName)
	if err != nil {
		return nil, err
	}
	res = &v1.UploadScheduleRes{
		FileName: fileName,
	}
	return
}
