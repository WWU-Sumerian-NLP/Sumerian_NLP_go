package IE_Extractor

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// TaggedTransliterations is a type that holds tagged transliteration data for a tablet.
type TaggedTransliterations struct {
	TabletNum      string
	taggedTranslit string //the entire tablets content
	Providence     string
	Period         string
	DateReferenced string
}

// CDLIParser is a type that parses TSV data and outputs the parsed data in TaggedTransliterations format.
type CDLIParser struct {
	path               string
	data               []TSVData
	currTabletTranslit string
	Out                chan TaggedTransliterations
	done               chan struct{}
}

// TSVData represents one row of data from the TSV file, with each column represented as a field in the struct.
type TSVData struct {
	tabletNum           string
	PUB                 string
	Providence          string
	Period              string
	DatesReferenced     string
	loc                 string
	no                  string
	raw_translit        string
	annotations         string
	normalized_translit string
	transli_entites     string
}

// NewCDLIParser constructs a new CDLIParser. The parser immediately starts reading and parsing the data at the given path.
func NewCDLIParser(path string) *CDLIParser {
	cdliParser := &CDLIParser{
		path: path,
		Out:  make(chan TaggedTransliterations, 10000000),
		done: make(chan struct{}, 1),
	}
	cdliParser.readCDLIData()
	cdliParser.run()
	return cdliParser
}

// run manages the main parsing loop for the parser. It reads rows from the data slice, parses the data, and sends the parsed data to the output channel.
func (p *CDLIParser) run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		println("CDLI parsing")
		defer wg.Done()
		defer close(p.Out)

		prevTabletNum := ""
		for _, tablet := range p.data {
			fmt.Printf("tablet: %v\n", tablet)
			p.collapseTabletTranslit(tablet, prevTabletNum)
			prevTabletNum = tablet.tabletNum
		}
		println("DONE")
		// p.out <- p.currCLDIData
	}()
	wg.Wait()
}

// WaitUntilDone allows external callers to wait until the parser has finished processing all of its data.
func (p *CDLIParser) WaitUntilDone() {
	p.done <- struct{}{}
}

// collapseTabletTranslit concatenates the transliterations of tablets that share a tablet number.
// It outputs a TaggedTransliterations value for each distinct tablet number.
func (p *CDLIParser) collapseTabletTranslit(tablet TSVData, prevTabletNum string) {
	if prevTabletNum != tablet.tabletNum {
		taggedTranslit := &TaggedTransliterations{TabletNum: prevTabletNum, taggedTranslit: p.currTabletTranslit, Providence: tablet.Providence, Period: tablet.Period, DateReferenced: tablet.DatesReferenced}
		p.currTabletTranslit = ""
		p.Out <- *taggedTranslit
	} else {
		p.currTabletTranslit += tablet.transli_entites
	}
}

// readCDLIData reads TSV data from the file at the path specified in the CDLIParser.
func (p *CDLIParser) readCDLIData() {
	csvFile, err := os.Open(p.path)
	if err != nil {
		log.Fatalf("failed reading file: %s", err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = '\t'

	// var cdliDataList []TSVData
	for {
		var data TSVData
		tsvData, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("issue here")
			log.Fatal(err)
		}
		data.tabletNum = tsvData[0]
		data.PUB = tsvData[1]
		data.Providence = tsvData[2]
		data.Period = tsvData[3]
		data.DatesReferenced = tsvData[4]
		data.loc = tsvData[5]
		data.no = tsvData[6]
		data.raw_translit = tsvData[7]
		data.annotations = tsvData[8]
		data.normalized_translit = tsvData[9]
		data.transli_entites = tsvData[10]

		p.data = append(p.data, data)
	}

}
