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
		nerList: []string{"city_ner.csv", "months_ner.csv", "royalname_ner.csv", "governors_ner.csv", "people_ner.csv", "animals_ner.csv", "foreigners_ner.csv"},
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
			for i, tablet := range cdliData.TabletSections {
				tablet.EntitiyLines = make(map[int]string)
				cdliData.TabletSections[i].EntitiyLines = e.getFromNERLists(tablet)
				cdliData.TabletSections[i].EntitiyLines = e.getFromAnnotations(tablet)
			}
			e.out <- cdliData
		}
		println("DONE")
		close(e.out)
	}()
	wg.Wait()
}

//case 1 - Get entities from seed rules (iti-month, mu-year)
func (e *CDLIEntityExtractor) getFromAnnotations(tableLines TabletSection) map[int]string {

	for line_no, translit := range tableLines.TabletLines {
		//regex expression tags the next word after iti
		if strings.Contains(translit, "iti ") {
			listString := strings.SplitAfter(translit, "iti")
			listString[1] = " (" + listString[1] + ", " + "MN" + ")"
			tableLines.EntitiyLines[line_no] = strings.Join(listString, " ")
		}
		if strings.Contains(translit, "mu ") {
			listString := strings.SplitAfterN(translit, "mu", 1)
			listString = strings.Split(strings.Join(listString, " "), " ")
			listString[1] = " (" + listString[1] + ", " + "YR" + ")"
			tableLines.EntitiyLines[line_no] = strings.Join(listString, " ")
		}
	}
	return tableLines.EntitiyLines
}

// case 2 - Get from NER_lists
func (e *CDLIEntityExtractor) getFromNERLists(tableLines TabletSection) map[int]string {
	for line_no, translit := range tableLines.TabletLines {
		new_translit := strings.Split(translit, " ")
		for i, grapheme := range strings.Split(translit, " ") {
			for _, list := range e.nerList { //fix
				nerMap := e.readNERLists(list)
				if ner, ok := nerMap[grapheme]; ok {
					new_translit[i] = "(" + grapheme + ", " + ner + ")"
				}
			}
		}
		tableLines.EntitiyLines[line_no] = strings.Join(new_translit, " ")
	}

	return tableLines.EntitiyLines
}

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
