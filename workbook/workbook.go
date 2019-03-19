package workbook

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	workbookException "gitlab.com/iotTracker/brain/workbook/exception"
)

type Workbook struct {
	File            *excelize.File
	SheetHeaderMaps map[string]map[string]string
}

func New(
	pathToWorkBook string,
	sheetFirstRowMap map[string]int,
) (*Workbook, error) {
	if sheetFirstRowMap == nil {
		sheetFirstRowMap = make(map[string]int)
	}

	// open the workbook
	file, err := excelize.OpenFile(pathToWorkBook)
	if err != nil {
		return nil, workbookException.OpeningFile{Reasons: []string{err.Error()}}
	}

	// build header map for each sheet
	sheetHeaderMaps := make(map[string]map[string]string)
	for _, sheetName := range file.GetSheetMap() {
		sheetHeaderMaps[sheetName], err = ColumnHeaderMap(file, sheetName, sheetFirstRowMap[sheetName])
		if err != nil {
			return nil, err
		}
	}

	return &Workbook{
		File:            file,
		SheetHeaderMaps: sheetHeaderMaps,
	}, nil
}
