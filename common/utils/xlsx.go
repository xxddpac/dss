package utils

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"os"
)

func WriteToXlsx(filePath string, header []string, members [][]string) error {
	file, err := xlsxFile(filePath, header)
	if err != nil {
		return err
	}
	xlsxSheet := file.Sheets[0]
	for _, member := range members {
		insertRow(xlsxSheet, member)
	}
	return file.Save(filePath)
}

func xlsxFile(filePath string, headers []string) (*xlsx.File, error) {
	var file *xlsx.File
	if Exists(filePath) {
		_ = os.Remove(filePath)
	}
	f, _ := os.Create(filePath)
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}(f)
	file = xlsx.NewFile()
	sheet, err := file.AddSheet("scan_result")
	if err != nil {
		return nil, fmt.Errorf("failed add xlsx add sheet, err:%s", err.Error())
	}
	if headers != nil {
		insertRow(sheet, headers)
	}
	if err = file.Save(filePath); err != nil {
		return nil, fmt.Errorf("failed save xlsx file, err:%s", err.Error())
	}
	return file, nil
}

func insertRow(sheet *xlsx.Sheet, rowData []string) {
	row := sheet.AddRow()
	for _, v := range rowData {
		cell := row.AddCell()
		cell.Value = v
	}
}
