package CDLI_Extractor

import (
	"encoding/csv"
	"log"
	"os"
	"sync"
)

type DataWriter struct {
	destPath  string
	in        <-chan CDLIData
	done      chan struct{}
	csvWriter *csv.Writer
}

func newDataWriter(destPath string, in <-chan CDLIData) *DataWriter {
	dataWriter := &DataWriter{
		destPath: destPath,
		in:       in,
		done:     make(chan struct{}, 1),
	}
	dataWriter.makeWriter()
	dataWriter.run()
	return dataWriter
}

func (w *DataWriter) run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for cdliData := range w.in {
			w.exportToCSV(cdliData)
		}
	}()
	wg.Wait()

	// go func() {
	// 	w.done <- struct{}{}
	// }()

}

func (w *DataWriter) WaitUntilDone() {
	w.done <- struct{}{}
}

func (w *DataWriter) makeWriter() {
	csvFile, err := os.Create(w.destPath)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvWriter := csv.NewWriter(csvFile)
	csvWriter.Comma = '\t'
	w.csvWriter = csvWriter
	w.csvWriter.Write([]string{"tablet num", "PUB", "loc", "no", "raw_translit", "annotations",
		"normalized_translit", "transli_entities", "entities"}) //hardcoded
	w.csvWriter.Flush()
}

func (w *DataWriter) exportToCSV(cdliData CDLIData) {
	cldiNo := cdliData.TabletNum
	cldiPub := cdliData.PUB
	for _, tablet := range cdliData.TabletList {
		for lineNo, translit := range tablet.TabletLines {
			w.csvWriter.Write([]string{cldiNo, cldiPub, tablet.TabletLocation, lineNo,
				translit, tablet.Annotation[lineNo], tablet.NormalizedLines[lineNo], tablet.EntitiyLines[lineNo]})
		}
		w.csvWriter.Flush()
	}
	w.csvWriter.Flush()

}
