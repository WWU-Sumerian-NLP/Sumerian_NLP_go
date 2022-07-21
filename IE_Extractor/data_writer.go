package IE_Extractor

import (
	"encoding/csv"
	"log"
	"os"
	"sync"
)

type DataWriter struct {
	destPath  string
	in        <-chan RelationData
	done      chan struct{}
	csvWriter *csv.Writer
}

func newDataWriter(destPath string, in <-chan RelationData) *DataWriter {
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
		for relationData := range w.in {
			w.exportToCSV(relationData)
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
	w.csvWriter.Write([]string{"tablet_num", "relation_type", "subject", "object"}) //hardcoded
	w.csvWriter.Flush()
}

func (w *DataWriter) exportToCSV(relationData RelationData) {
	tabletNum := relationData.tabletNum
	relationType := relationData.relationType
	relationSubject := relationData.relationTuple[1]
	relationObject := relationData.relationTuple[2]

	w.csvWriter.Write([]string{tabletNum, relationType, relationSubject, relationObject})

	w.csvWriter.Flush()

}
