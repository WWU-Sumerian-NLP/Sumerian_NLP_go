package CDLI_Extractor

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"sync"
)

type CDLIDataWriter struct {
	destPath  string
	in        <-chan CDLIData
	done      chan struct{}
	csvWriter *csv.Writer
}

func newCDLIDataWriter(destPath string, in <-chan CDLIData) *CDLIDataWriter {
	dataWriter := &CDLIDataWriter{
		destPath: destPath,
		in:       in,
		done:     make(chan struct{}, 1),
	}
	dataWriter.makeWriter()
	dataWriter.run()
	return dataWriter
}

func (w *CDLIDataWriter) run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for cdliData := range w.in {
			w.exportToCSV(cdliData)
		}
	}()
	wg.Wait()
}

func (w *CDLIDataWriter) WaitUntilDone() {
	w.done <- struct{}{}
}

func (w *CDLIDataWriter) makeWriter() {
	csvFile, err := os.Create(w.destPath)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvWriter := csv.NewWriter(csvFile)
	csvWriter.Comma = '\t'
	w.csvWriter = csvWriter
	w.csvWriter.Write([]string{"tablet num", "PUB", "Providence", "Period", "Dates Referenced", "loc", "no", "raw_translit", "annotations",
		"normalized_translit"}) //hardcoded
	w.csvWriter.Flush()
}

func (w *CDLIDataWriter) exportToCSV(cdliData CDLIData) {
	cldiNo := cdliData.TabletNum
	cldiPub := cdliData.PUB
	cdliProv := cdliData.Provenience
	cdliPeriod := cdliData.Period
	cdliDates := cdliData.DatesReferenced

	for _, tablet := range cdliData.TabletSections {
		for i, lineNo := range tablet.LineNumbers {
			w.csvWriter.Write([]string{cldiNo, cldiPub, cdliProv, cdliPeriod, cdliDates, tablet.TabletLocation, strconv.Itoa(lineNo),
				tablet.TabletLines[i], tablet.Annotations[i], tablet.NormalizedLines[i]})
		}
		w.csvWriter.Flush()
	}
	w.csvWriter.Flush()

}
