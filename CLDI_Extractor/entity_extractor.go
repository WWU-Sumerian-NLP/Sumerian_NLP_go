package CLDI_Extractor

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type CLDIEntityExtractor struct {
	in      <-chan CLDIData
	out     chan CLDIData
	done    chan struct{}
	nerList []string
}

func newCLDIEntityExtractor(in <-chan CLDIData) *CLDIEntityExtractor {
	entityExtractor := &CLDIEntityExtractor{
		in:      in,
		out:     make(chan CLDIData, 1000),
		done:    make(chan struct{}, 1),
		nerList: []string{"city_ner.csv", "months_ner.csv", "royalname_ner.csv", "animals_ner.csv", "governors_ner.csv", "people_ner.csv"},
	}
	entityExtractor.run()
	return entityExtractor
}

//generalize function so user just passes a file name
//and then a temporary map is created to parse through
func (e *CLDIEntityExtractor) readNERLists(nerListName string) map[string]string {
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
	return nerMap
}

func (e *CLDIEntityExtractor) WaitUntilDone() {
	e.done <- struct{}{}
}

func (e *CLDIEntityExtractor) run() {
	go func() {
		for cldiData := range e.in {
			for i, tablet := range cldiData.TabletList {
				tablet.EntitiyLines = make(map[string]string)
				cldiData.TabletList[i].EntitiyLines = e.getFromAnnotations(tablet)
				for _, list := range e.nerList {
					nerMap := e.readNERLists(list)
					cldiData.TabletList[i].EntitiyLines = e.getFromNERLists(tablet, nerMap)
				}
			}
			e.out <- cldiData
		}
		close(e.out)
	}()

	go func() {
		e.done <- struct{}{}
	}()
}

//case 1 - Get entities from seed rules (iti-month, mu-year)
func (e *CLDIEntityExtractor) getFromAnnotations(tableLines TabletLine) map[string]string {

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
func (e *CLDIEntityExtractor) getFromNERLists(tableLines TabletLine, nerMap map[string]string) map[string]string {
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
