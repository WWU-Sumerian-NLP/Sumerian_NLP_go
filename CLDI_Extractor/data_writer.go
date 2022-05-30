package CLDI_Extractor

import (
	"encoding/csv"
	"log"
	"os"
)

type DataWriter struct {
	destPath  string
	in        <-chan CLDIData
	done      chan struct{}
	csvWriter *csv.Writer
}

func newDataWriter(destPath string, in <-chan CLDIData) *DataWriter {
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
	go func() {
		for cldiData := range w.in {
			w.exportToCSV(cldiData)
		}
	}()

	go func() {
		w.done <- struct{}{}
	}()

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
	w.csvWriter.Write([]string{"CLDI", "PUB", "loc", "no", "raw_translit",
		"normalized_translit", "annotations", "transli_entities", "entities"}) //hardcoded
	w.csvWriter.Flush()
}

func (w *DataWriter) exportToCSV(cldiData CLDIData) {
	w.csvWriter.Write([]string{cldiData.CLDI, cldiData.PUB, cldiData.tabletLines["loc"], cldiData.tabletLines["no"],
		cldiData.tabletLines["translit"], cldiData.tabletLines["normalized_translit"],
		cldiData.tabletLines["annotations"], cldiData.tabletLines["transli_entities"],
		cldiData.tabletLines["entities"]})
	w.csvWriter.Flush()
}
