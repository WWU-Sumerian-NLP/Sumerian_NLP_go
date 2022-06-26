package Data_Collection

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

type CLDITSV struct {
	tabletNum       string
	publication     string
	lineNumber      string
	transliteration string
	annotation      string
}

func readTSV(path string) []CLDITSV {
	csvFile, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed reading file: %s", err)
	}
	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = '\t'

	var dataList []CLDITSV
	for {
		var data CLDITSV

		tsvData, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			println("issue")
		}
		data.tabletNum = tsvData[0]
		data.publication = tsvData[1]
		data.lineNumber = tsvData[2]
		data.transliteration = tsvData[3]
		data.annotation = tsvData[4]
		dataList = append(dataList, data)
	}
	csvFile.Close()
	return dataList
}
