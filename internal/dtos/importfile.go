package dtos

import (
	"io"

	"github.com/xuri/excelize/v2"
)

type ImportFile struct {
	Excel *excelize.File
}

func NewImportFile(f io.Reader) (*ImportFile, error) {
	file, err := excelize.OpenReader(f)
	if err != nil {
		return nil, err
	}
	return &ImportFile{
		Excel: file,
	}, nil
}

type DeviceAddResponse struct {
	List       []DeviceAddResult `json:"list"`
	ProcessNum int               `json:"processNum"`
	SuccessNum int               `json:"successNum"`
	FailNum    int               `json:"failNum"`
}

type DeviceAddResult struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
