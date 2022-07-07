package IE_Extractor

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type TaggedTransliterations struct {
	TabletNum         string
	taggedTranslit    string //the entire tablets content
	deliveryRelations []string
	recievedRelations []string
}
type CDLIParser struct {
	path               string
	data               []TSVData
	currTabletTranslit string
	out                chan TaggedTransliterations
	done               chan struct{}
}

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

func newCDLIParser(path string) *CDLIParser {
	cdliParser := &CDLIParser{
		path: path,
		out:  make(chan TaggedTransliterations, 100000),
		done: make(chan struct{}, 1),
	}
	cdliParser.readCDLIData()
	cdliParser.run()
	return cdliParser
}

func (p *CDLIParser) run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		println("CDLI parsing")
		defer wg.Done()
		defer close(p.out)

		prevTabletNum := ""
		for _, tablet := range p.data {
			// fmt.Printf("tablet: %v\n", tablet)
			p.collapseTabletTranslit(tablet, prevTabletNum)
			prevTabletNum = tablet.tabletNum
		}
		println("DONE")
		// p.out <- p.currCLDIData
	}()
	wg.Wait()
}

func (p *CDLIParser) WaitUntilDone() {
	p.done <- struct{}{}
}

func (p *CDLIParser) collapseTabletTranslit(tablet TSVData, prevTabletNum string) {
	if prevTabletNum != tablet.tabletNum {
		taggedTranslit := &TaggedTransliterations{TabletNum: prevTabletNum, taggedTranslit: p.currTabletTranslit}
		// fmt.Printf("taggedTranslit: %v\n", taggedTranslit)
		p.currTabletTranslit = ""
		p.out <- *taggedTranslit
	} else {
		p.currTabletTranslit += tablet.transli_entites
	}
}

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
		fmt.Printf("tsvData: %v\n", tsvData)
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
		// cdliDataList = append(cdliDataList, data)
	}

}
