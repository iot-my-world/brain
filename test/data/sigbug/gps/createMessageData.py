from openpyxl import load_workbook
from numpy import float32

if __name__ == '__main__':
    workbook = load_workbook('data.xlsx')
    for worksheet in workbook.worksheets:
        for row in list(worksheet.rows)[1:]:
            latHexBytes = ['00' if b == 0 else hex(b).strip('0x') for b in float32(float(row[0].value)).tobytes('C')]
            latHexValue = ''.join(b if len(b) == 2 else '0' + b for b in latHexBytes)
            lonHexBytes = ['00' if b == 0 else hex(b).strip('0x') for b in float32(float(row[1].value)).tobytes('C')]
            lonHexValue = ''.join(b if len(b) == 2 else '0' + b for b in lonHexBytes)
            row[3].value = '00%s%s' % (latHexValue, lonHexValue)
    workbook.save('data.xlsx')
    print('done')
