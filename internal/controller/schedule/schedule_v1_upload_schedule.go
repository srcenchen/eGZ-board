package schedule

import (
	"context"
	"eGZ-Board/api/schedule/v1"
	"eGZ-Board/internal/dao"
	"eGZ-Board/utility"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func (c *ControllerV1) UploadSchedule(ctx context.Context, req *v1.UploadScheduleReq) (res *v1.UploadScheduleRes, err error) {
	originalName := filepath.Base(req.File.Filename)
	if strings.ToLower(filepath.Ext(originalName)) != ".xlsx" {
		return nil, fmt.Errorf("schedule file must be .xlsx")
	}

	fileName, err := req.File.Save("./resource/upload", true)
	if err != nil {
		return nil, fmt.Errorf("save schedule upload: %w", err)
	}
	filePath := filepath.Join("./resource/upload", fileName)
	defer func() {
		if removeErr := os.Remove(filePath); err == nil && removeErr != nil && !os.IsNotExist(removeErr) {
			err = fmt.Errorf("remove temporary schedule: %w", removeErr)
			res = nil
		}
	}()

	data, err := parseSchedule(filePath)
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("encode schedule: %w", err)
	}
	classID, err := utility.GetClassID(ctx)
	if err != nil {
		return nil, err
	}
	if err := dao.UpsertSchedule(classID, string(jsonData)); err != nil {
		return nil, fmt.Errorf("store schedule: %w", err)
	}
	return &v1.UploadScheduleRes{FileName: originalName}, nil
}

func parseSchedule(filePath string) (data map[string][]string, err error) {
	workbook, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("open schedule workbook: %w", err)
	}
	defer func() {
		if closeErr := workbook.Close(); err == nil && closeErr != nil {
			data = nil
			err = fmt.Errorf("close schedule workbook: %w", closeErr)
		}
	}()

	sheetName := workbook.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("schedule workbook has no worksheet")
	}
	rows, err := workbook.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("read schedule worksheet: %w", err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("schedule worksheet has no data rows")
	}

	data = make(map[string][]string)
	for _, row := range rows[1:] {
		for column := 1; column < len(row); column++ {
			key := strconv.Itoa(column)
			data[key] = append(data[key], row[column])
		}
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("schedule worksheet has no course columns")
	}
	return data, nil
}
