package dtos

import (
	"fmt"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

const (
	DevicesFilename = "Devices"
)

type ExportFile struct {
	Excel    *excelize.File
	FileName string
}

func NewExportFile(filename string) (*ExportFile, error) {
	ef := &ExportFile{
		FileName: newFileName(filename),
		Excel:    excelize.NewFile(),
	}
	return ef, nil
}

func newFileName(name string) string {
	date := time.Now().Format("2006-01-02")
	unix := strconv.FormatInt(time.Now().Unix(), 10)
	return fmt.Sprintf("%s_%s_%s.xlsx", name, date, unix)
}

func (f *ExportFile) GetCenterStyle() int {
	style, _ := f.Excel.NewStyle(`{"alignment":{"horizontal": "center","vertical": "center"}}`)
	return style
}

func (f *ExportFile) GetRequiredStyle() int {
	style, _ := f.Excel.NewStyle(`{"alignment":{"horizontal": "center","vertical": "center"},"font":{"color": "#ea4335"}}`)
	return style
}
