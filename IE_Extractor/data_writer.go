package IE_Extractor

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
	"sync"
)

type DataWriter struct {
	destPath  string
	in        <-chan TaggedTransliterations
	done      chan struct{}
	csvWriter *csv.Writer
}

func newDataWriter(destPath string, in <-chan TaggedTransliterations) *DataWriter {
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
	w.csvWriter.Write([]string{"tablet num", "tagged_translit", "delivery_rel", "recieved_rel"}) //hardcoded
	w.csvWriter.Flush()
}

func (w *DataWriter) exportToCSV(cdliData TaggedTransliterations) {
	cldiNo := cdliData.TabletNum
	cdliTaggedTranslit := cdliData.taggedTranslit
	delieveryTuples := strings.Join(cdliData.deliveryRelations, " ")
	recievedTuples := strings.Join(cdliData.recievedRelations, " ")
	w.csvWriter.Write([]string{cldiNo, cdliTaggedTranslit, delieveryTuples, recievedTuples})

	w.csvWriter.Flush()

}
