package workbook

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	workbookException "gitlab.com/iotTracker/brain/workbook/exception"
)

type Workbook struct {
	File             *excelize.File
	SheetFirstRowMap map[string]int
	SheetHeaderMaps  map[string]map[string]string
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
		File:             file,
		SheetFirstRowMap: sheetFirstRowMap,
		SheetHeaderMaps:  sheetHeaderMaps,
	}, nil
}

func (w *Workbook) DataRows(sheetName string) {

}

func (w *Workbook) SheetAsSliceMap(sheetName string) ([]map[string]string, error) {
	for sheetIdx, sheetInBookName := range w.File.GetSheetMap() {
		if sheetInBookName == sheetName {
			break
		}
		if sheetIdx == len(w.File.GetSheetMap())-1 {
			return nil, workbookException.SheetDoesNotExist{SheetName: sheetName}
		}
	}

	sheetSliceMap := make([]map[string]string, 0)
	sheetFirstRowIdx := w.SheetFirstRowMap[sheetName]
	for rowIdx := range w.File.GetRows(sheetName)[sheetFirstRowIdx:] {
		rowMap := make(map[string]string)
		for header, column := range w.SheetHeaderMaps[sheetName] {
			rowMap[header] = w.File.GetCellValue(sheetName, fmt.Sprintf("%s%d", column, rowIdx+1))
		}
		sheetSliceMap = append(sheetSliceMap, rowMap)
	}

	return sheetSliceMap, nil
}
