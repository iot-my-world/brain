import os, re, openpyxl
from openpyxl.utils import get_column_letter, column_index_from_string


def first_empty_column(worksheet):
    return get_column_letter(len(list(worksheet.columns)) + 1)


def next_column_letter(column_letter):
    return get_column_letter(column_index_from_string(column_letter) + 1)


stringPattern = r'"(.*?)"'

if __name__ == '__main__':
    readingFilePaths = [readingFile for readingFile in os.listdir('.') if readingFile.endswith('.rdat')]
    for filePath in readingFilePaths:
        with open(filePath) as readingFile:
            for line in readingFile:
                if 'lat' in line and 'lon' in line:
                    print(re.findall(stringPattern, line))
