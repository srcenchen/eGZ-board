package schedule

import (
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestParseSchedule(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "schedule.xlsx")
	workbook := excelize.NewFile()
	sheet := workbook.GetSheetName(0)
	rows := [][]interface{}{
		{"period", "Monday", "Tuesday"},
		{"1", "Math", "English"},
		{"2", "Physics", "History"},
	}
	for rowIndex, row := range rows {
		for columnIndex, value := range row {
			cell, err := excelize.CoordinatesToCellName(columnIndex+1, rowIndex+1)
			if err != nil {
				t.Fatalf("create cell name: %v", err)
			}
			if err := workbook.SetCellValue(sheet, cell, value); err != nil {
				t.Fatalf("set cell: %v", err)
			}
		}
	}
	if err := workbook.SaveAs(filePath); err != nil {
		t.Fatalf("save workbook: %v", err)
	}
	if err := workbook.Close(); err != nil {
		t.Fatalf("close workbook: %v", err)
	}

	data, err := parseSchedule(filePath)
	if err != nil {
		t.Fatalf("parse schedule: %v", err)
	}
	if got := data["1"]; len(got) != 2 || got[0] != "Math" || got[1] != "Physics" {
		t.Fatalf("Monday data = %#v", got)
	}
	if got := data["2"]; len(got) != 2 || got[0] != "English" || got[1] != "History" {
		t.Fatalf("Tuesday data = %#v", got)
	}
}

func TestParseScheduleRejectsEmptyWorkbook(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "empty.xlsx")
	workbook := excelize.NewFile()
	if err := workbook.SaveAs(filePath); err != nil {
		t.Fatalf("save workbook: %v", err)
	}
	if err := workbook.Close(); err != nil {
		t.Fatalf("close workbook: %v", err)
	}
	if _, err := parseSchedule(filePath); err == nil {
		t.Fatal("parseSchedule returned nil error for empty workbook")
	}
}
