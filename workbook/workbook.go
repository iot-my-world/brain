package workbook

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	workbookException "github.com/iot-my-world/brain/workbook/exception"
)

type Workbook struct {
	File              *excelize.File
	SheetHeaderRowMap map[string]int
	SheetHeaderMaps   map[string]map[string]string
}

func New(
	pathToWorkBook string,
	sheetHeaderRowMap map[string]int,
) (*Workbook, error) {
	if sheetHeaderRowMap == nil {
		sheetHeaderRowMap = make(map[string]int)
	}

	// open the workbook
	file, err := excelize.OpenFile(pathToWorkBook)
	if err != nil {
		return nil, workbookException.OpeningFile{Reasons: []string{err.Error()}}
	}

	// build header map for each sheet
	sheetHeaderMaps := make(map[string]map[string]string)
	for _, sheetName := range file.GetSheetMap() {
		sheetHeaderMaps[sheetName], err = ColumnHeaderMap(file, sheetName, sheetHeaderRowMap[sheetName])
		if err != nil {
			return nil, err
		}
	}

	return &Workbook{
		File:              file,
		SheetHeaderRowMap: sheetHeaderRowMap,
		SheetHeaderMaps:   sheetHeaderMaps,
	}, nil
}

func (w *Workbook) DataRows(sheetName string) {

}

func (w *Workbook) SheetAsSliceMap(sheetName string) ([]map[string]string, error) {
	noSheets := len(w.GetSheetNames())
	sheetCount := 0
	for _, sheetInBookName := range w.File.GetSheetMap() {
		if sheetInBookName == sheetName {
			break
		}
		sheetCount++
		if sheetCount == noSheets {
			return nil, workbookException.SheetDoesNotExist{SheetName: sheetName}
		}
	}

	sheetSliceMap := make([]map[string]string, 0)
	sheetHeaderRowIdx := w.SheetHeaderRowMap[sheetName]
	rowsWithHeader, err := w.File.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	rows := rowsWithHeader[sheetHeaderRowIdx+1:]
	for rowIdx := range rows {
		if err != nil {
			return nil, err
		}
		rowMap := make(map[string]string)
		for header, column := range w.SheetHeaderMaps[sheetName] {
			cellRef := fmt.Sprintf("%s%d", column, rowIdx+sheetHeaderRowIdx+2)
			rowMap[header], err = w.File.GetCellValue(sheetName, cellRef)
		}
		sheetSliceMap = append(sheetSliceMap, rowMap)
	}

	return sheetSliceMap, nil
}

func (w *Workbook) GetSheetNames() []string {
	sheetNames := make([]string, 0)
	for _, sheetName := range w.File.GetSheetMap() {
		sheetNames = append(sheetNames, sheetName)
	}

	return sheetNames
}
