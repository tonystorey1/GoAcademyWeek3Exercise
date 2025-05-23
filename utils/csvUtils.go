package utils

import (
	"Basic_CLI_Application/consts"
	"encoding/csv"
	"os"
)

// OpenOrCreateFile opens the supplied CSV file; panics if it can't
func OpenOrCreateFile() *os.File {
	f, err := os.OpenFile(consts.FileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	return f
}

func CreateCSVWriter(filename string) (*csv.Writer, *os.File, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, nil, err
	}
	writer := csv.NewWriter(f)
	writer.Comma = ','
	writer.UseCRLF = true
	return writer, f, nil
}
