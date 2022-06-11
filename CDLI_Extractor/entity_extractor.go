package CDLI_Extractor

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type CDLIEntityExtractor struct {
	in      <-chan CDLIData
	out     chan CDLIData
	done    chan struct{}
	nerList []string
}

func newCDLIEntityExtractor(in <-chan CDLIData) *CDLIEntityExtractor {
	entityExtractor := &CDLIEntityExtractor{
		in:      in,
		out:     make(chan CDLIData, 1000000),
		done:    make(chan struct{}, 1),
		nerList: []string{"city_ner.csv", "months_ner.csv", "royalname_ner.csv", "animals_ner.csv", "governors_ner.csv", "people_ner.csv"},
	}
	entityExtractor.run()
	return entityExtractor
}

func (e *CDLIEntityExtractor) WaitUntilDone() {
	e.done <- struct{}{}
}

func (e *CDLIEntityExtractor) run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		println("entity extracting")
		defer wg.Done()
		for cdliData := range e.in {
			for i, tablet := range cdliData.TabletList {
				tablet.EntitiyLines = make(map[string]string)
				cdliData.TabletList[i].EntitiyLines = e.getFromAnnotations(tablet)
				for _, list := range e.nerList {
					nerMap := e.readNERLists(list)
					cdliData.TabletList[i].EntitiyLines = e.getFromNERLists(tablet, nerMap)
				}
			}
			e.out <- cdliData
		}
		println("DONE")
		close(e.out)
	}()
	wg.Wait()
}

//case 1 - Get entities from seed rules (iti-month, mu-year)
func (e *CDLIEntityExtractor) getFromAnnotations(tableLines TabletLine) map[string]string {

	for line_no, translit := range tableLines.TabletLines {

		if strings.Contains(translit, "iti") {
			tableLines.EntitiyLines[line_no] = strings.ReplaceAll(translit, "iti", "iti[month]")
		}
		if strings.Contains(translit, "mu") {
			tableLines.EntitiyLines[line_no] = strings.ReplaceAll(translit, "mu", "mu(year)")
		}
	}
	return tableLines.EntitiyLines
}

// case 2 - Get from NER_lists
func (e *CDLIEntityExtractor) getFromNERLists(tableLines TabletLine, nerMap map[string]string) map[string]string {
	for line_no, translit := range tableLines.TabletLines {
		//might have to iterate and split string by " "
		for ner := range nerMap {
			if strings.Contains(translit, ner) {
				//todo - Replacement shouldn't happen at the specific instance, instead at the end of the character
				tableLines.EntitiyLines[line_no] = strings.ReplaceAll(translit, ner, ner+"("+nerMap[ner]+")")
			}
		}
	}
	return tableLines.EntitiyLines
}

//generalize function so user just passes a file name
//and then a temporary map is created to parse through
func (e *CDLIEntityExtractor) readNERLists(nerListName string) map[string]string {
	//city
	csvFile, err := os.Open(filepath.Join("../Annotation_lists/NER_lists", nerListName))
	if err != nil {
		log.Fatalf("failed reading file: %s", err)
	}
	csvReader := csv.NewReader(csvFile)
	nerCSV, err := csvReader.ReadAll()
	nerMap := make(map[string]string)
	for _, ner := range nerCSV {
		nerMap[ner[0]] = ner[1]
	}
	csvFile.Close()
	return nerMap
}
