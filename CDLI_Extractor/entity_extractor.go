/*
Package CDLI_Extractor is used to extract entities from the CDLI (Cuneiform Digital Library Initiative) data.
It uses Named Entity Recognition (NER) lists to label entities in the input data.
*/

package CDLI_Extractor

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// CDLIEntityExtractor is the main struct for this package.
// It holds the input and output channels for the data, a list of NER lists, and maps for storing and handling the NER data.
type CDLIEntityExtractor struct {
	in         <-chan CDLIData
	Out        chan CDLIData
	done       chan struct{}
	nerList    []string
	nerMap     map[string]map[string]string
	TempNERMap map[string]string
}

// NewCDLIEntityExtractor is the constructor function for the CDLIEntityExtractor.
// It initializes a new instance of the CDLIEntityExtractor with a given input channel.
func NewCDLIEntityExtractor(in <-chan CDLIData) *CDLIEntityExtractor {
	entityExtractor := &CDLIEntityExtractor{
		in:   in,
		Out:  make(chan CDLIData, 100000000),
		done: make(chan struct{}, 1),
		nerList: []string{"city_ner.csv", "months_ner.csv", "royalname_ner.csv", "governors_ner.csv", "people_ner.csv", "animals_ner.csv", "foreigners_ner.csv",
			"agricultural_locus_ner.csv", "ancestral_clan_line_ner.csv", "celestial_ner.csv", "city_quarter_ner.csv", "divine_ner.csv", "ethnos_ner.csv", "field_ner.csv",
			"geographical_ner.csv", "object_ner.csv", "temple_ner.csv", "watercourse_ner.csv"},
		nerMap:     make(map[string]map[string]string, 0),
		TempNERMap: make(map[string]string),
	}
	entityExtractor.readNERLists()
	entityExtractor.run()
	return entityExtractor
}

// WaitUntilDone sends a signal that the entity extraction is done.
func (e *CDLIEntityExtractor) WaitUntilDone() {
	e.done <- struct{}{}
}

// run starts the entity extraction process. It traverses each CDLI data from the input channel and labels the entities found.
func (e *CDLIEntityExtractor) run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		println("entity extracting")
		defer wg.Done()
		for cdliData := range e.in {
			for i, tablet := range cdliData.TabletSections {
				tablet.EntitiyLines = make([]string, len(tablet.LineNumbers))
				cdliData.TabletSections[i].EntitiyLines = e.labelAllGraphemes(tablet)
			}
			e.Out <- cdliData
		}
		println("DONE")
		close(e.Out)
	}()
	wg.Wait()
}

/*

Function that traverses a tablet section grapheme by grapheme

For each grapheme, it will check to see if
1. It exists in one of the NER lists
2. The previous word was iti or mu
3. Default

Then it will simply just do
(grapheme, O) where 0 = default

Then, I'll need a way to parse all this info (parse by parenthesis)

*/

// labelAllGraphemes labels the entities in all graphemes of a tablet section.
func (e *CDLIEntityExtractor) labelAllGraphemes(tabletLines TabletSection) []string {
	for line_no, translit := range tabletLines.TabletLines {
		grapheme_list := strings.Split(translit, " ")
		for i, grapheme := range grapheme_list {
			// grapheme = e.getFromNERLists(grapheme) //first get from annotation lists
			grapheme = e.getFromTempNERList(grapheme)
			grapheme = e.labelRelation(grapheme)
			if !strings.Contains(grapheme, ",") { //second, based on context (n-1)
				if i > 0 && grapheme_list[i-1] == "iti" {
					grapheme = "(" + grapheme + "," + "MN" + ")"
				} else if i > 0 && grapheme_list[i-1] == "mu" {
					grapheme = "(" + grapheme + "," + "YR" + ")"
				} else {
					println(grapheme)
					grapheme = "(" + grapheme + "," + "O" + ")"
				}
			}
			tabletLines.EntitiyLines[line_no] += grapheme + " "
		}
	}
	return tabletLines.EntitiyLines
}

// getFromTempNERList retrieves the entity label for a grapheme from the temporary NER list.
func (e *CDLIEntityExtractor) getFromTempNERList(grapheme string) string {
	new_grapheme := grapheme
	if ner, ok := e.TempNERMap[grapheme]; ok {
		new_grapheme = "(" + grapheme + "," + ner + ")"
	}
	return new_grapheme
}

// getFromNERLists retrieves the entity label for a grapheme from the NER lists.
//Get from NER_lists
func (e *CDLIEntityExtractor) getFromNERLists(grapheme string) string {
	new_grapheme := grapheme
	for fileName := range e.nerMap {
		nerMap := e.nerMap[fileName]
		if ner, ok := nerMap[grapheme]; ok {
			new_grapheme = "(" + grapheme + "," + ner + ")"
		}
	}

	return new_grapheme
}

// labelRelation labels a grapheme based on a specific list of relations.
//TODO - Read a list?
func (e *CDLIEntityExtractor) labelRelation(grapheme string) string {
	if grapheme == "mu-kux(DU)" {
		grapheme = "(" + grapheme + "," + "DEL" + ")"
	} else if grapheme == "i3-dab5" {
		grapheme = "(" + grapheme + "," + "REC" + ")"
	} else if grapheme == "ba-zi" {
		grapheme = "(" + grapheme + "," + "DIS" + ")"
	} else if grapheme == "ba-ti" { //this is wrong, should be sz ba-ti
		grapheme = "(" + grapheme + "," + "REC" + ")"
	} else if grapheme == "dumu" {
		grapheme = "(" + grapheme + "," + "SON" + ")"
	} else if grapheme == "ab" {
		grapheme = "(" + grapheme + "," + "FATHER" + ")"
	} else if grapheme == "dam" {
		grapheme = "(" + grapheme + "," + "WIFE" + ")"
	}
	return grapheme

}

// readNERLists reads the NER lists from CSV files and stores them in a map.
//Read a list of NER lists [ner.csv, ner2.csv] and store a map of filenames to a map of entity relations
func (e *CDLIEntityExtractor) readNERLists() {
	fileNameToNerMap := make(map[string]map[string]string, len(e.nerList))
	for _, nerListName := range e.nerList {
		csvFile, err := os.Open(filepath.Join("../Annotation_lists/NER_lists", nerListName))
		if err != nil {
			log.Fatalf("failed reading file: %s", err)
		}
		csvReader := csv.NewReader(csvFile)
		nerCSV, err := csvReader.ReadAll()
		if err != nil {
			log.Fatalf("error: %s failed parsing file: %s", nerListName, err)
		}
		nerMap := make(map[string]string)
		for _, ner := range nerCSV {
			nerMap[ner[0]] = ner[1]
		}
		fileNameToNerMap[nerListName] = nerMap
		csvFile.Close()
	}
	e.nerMap = fileNameToNerMap
}
