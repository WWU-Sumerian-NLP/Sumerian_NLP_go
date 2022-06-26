package CDLI_Extractor

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

type TSVData struct {
	tabletNum    string
	PUB          string
	loc          string
	no           string
	raw_translit string
}

func readCDLIData(path string) chan CDLIData {
	csvFile, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed reading file: %s", err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = '\t'

	in := make(chan CDLIData, 1000000)
	var cdliDataList []TSVData

	for {
		var data TSVData
		tsvData, err := csvReader.Read()
		fmt.Printf("tsvData: %v\n", tsvData)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("issue here")
			log.Fatal(err)
		}
		data.tabletNum = tsvData[0]
		data.PUB = tsvData[1]
		data.loc = tsvData[2]
		data.no = tsvData[3]
		data.raw_translit = tsvData[4]
		cdliDataList = append(cdliDataList, data)
	}

	// for _, cdliData := range cdliDataList {
	// 	fmt.Printf("cdliData: %v\n", cdliData)
	// 	in <- cdliData
	// }
	// close(in)
	return in
}
