package CDLI_Extractor

import (
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/jszwec/csvutil"
)

type Strings []string

func readCDLIData(path string) chan CDLIData {
	csvFile, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed reading file: %s", err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = '\t'

	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		log.Fatal(err)
	}
	in := make(chan CDLIData, 1000000)
	var cdliDataList []CDLIData

	for {
		var c CDLIData
		if err := dec.Decode(&c); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		cdliDataList = append(cdliDataList, c)
	}

	for _, cdliData := range cdliDataList {
		in <- cdliData
	}
	close(in)
	return in
}
