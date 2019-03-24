package workbook

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	workbookException "gitlab.com/iotTracker/brain/workbook/exception"
)

/*
	Returns a header map for the given excel file sheet
*/
func ColumnHeaderMap(xlsxFile *excelize.File, sheet string, topRowIdx int) (map[string]string, error) {
	columnHeaderMap := make(map[string]string)
	rows, err := xlsxFile.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	if len(rows)-1 < topRowIdx {
		return nil, workbookException.NotEnoughRowsInSheet{
			Reasons: []string{
				fmt.Sprintf("only %d rows", len(rows)),
				fmt.Sprintf("should be %d rows", topRowIdx+1),
			}}
	}
	for colIdx, colCell := range rows[topRowIdx] {
		columnHeaderMap[colCell], err = excelize.ColumnNumberToName(colIdx)
	}

	return columnHeaderMap, nil
}
