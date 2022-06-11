package CDLI_Extractor

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

//Lines map to a line_no, transliterations, normalizations and annotations
type CDLIData struct {
	TabletNum  string
	PUB        string
	TabletList []TabletLine
}

type TabletLine struct {
	TabletLocation  string
	TabletLines     map[string]string //map line_no to transliterations
	NormalizedLines map[string]string
	EntitiyLines    map[string]string
	Annotation      map[string]string //map line_no to annotation
	maxLine         int
}

type ATFParser struct {
	path         string
	data         []string
	currCLDIData CDLIData
	out          chan CDLIData
	done         chan struct{}
	t            int
}

func newATFParser(path string) *ATFParser {
	atfParser := &ATFParser{
		path: path,
		out:  make(chan CDLIData, 100000),
		done: make(chan struct{}, 1),
	}
	atfParser.loadCDLIData()
	atfParser.run()
	return atfParser
}

func (p *ATFParser) run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		println("atf parsing")
		defer wg.Done()
		defer close(p.out)
		for _, line := range p.data {
			// println(line)
			p.parseLines(line)
		}
		println("DONE")
		p.out <- p.currCLDIData
	}()

	// go func() {
	// 	println("WAIT")
	wg.Wait()
	// }()

	// defer p.WaitUntilDone()
	// go func() {
	// 	p.done <- struct{}{}
	// }()
}

func (p *ATFParser) WaitUntilDone() {
	p.done <- struct{}{}
}

func (p *ATFParser) loadCDLIData() {
	f, err := os.Open(p.path)
	if err != nil {
		log.Fatalf("failed reading file: %s", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		p.data = append(p.data, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("issue reading this file: %s", err)
	}
}

/*
Parse Tablet Line-By-Line.
	- &P indicates a new tablet, we initalize p.currCLDIData for a new tablet
	- Store CLDI, PUB, annotation and transliterations data to this object
*/

func (p *ATFParser) parseLines(line string) {
	line = strings.TrimSpace(line)
	if line != "" && strings.Contains(line, "&P") && strings.Contains(line, " =") {
		if p.currCLDIData.PUB == "" {
			line = strings.ReplaceAll(line, "&", "")
			line = strings.TrimSpace(line)
			data := strings.SplitN(line, " = ", 2)
			p.currCLDIData.TabletNum = data[0]
			p.currCLDIData.PUB = data[1]

		} else {
			p.out <- p.currCLDIData
			p.currCLDIData = CDLIData{}
			p.t++

		}
		line = strings.ReplaceAll(line, "&", "")
		line = strings.TrimSpace(line)
		data := strings.SplitN(line, " = ", 2)
		p.currCLDIData.TabletNum = data[0]
		p.currCLDIData.PUB = data[1]

	} else if line != "" && strings.Contains(string(line[0:1]), "@") && !strings.Contains(line, "object") && !strings.Contains(line, "tablet") && !strings.Contains(line, "envelope") && !strings.Contains(line, "bulla") {
		newTableLine := TabletLine{}
		newTableLine.TabletLocation = strings.TrimSpace(line)
		newTableLine.TabletLines = make(map[string]string)
		p.currCLDIData.TabletList = append(p.currCLDIData.TabletList, newTableLine)
		p.currCLDIData.TabletList[len(p.currCLDIData.TabletList)-1].Annotation = make(map[string]string)

	} else if strings.Contains(line, "#tr.en") {
		// You can translate tr.en entries
		line = strings.Replace(line, "#tr.en", "", 1)
		line = strings.Replace(line, ":", "", 1)
		line := strings.TrimSpace(line)

		// TabletLine.
		if len(p.currCLDIData.TabletList) > 0 {
			line_no := strconv.Itoa(p.currCLDIData.TabletList[len(p.currCLDIData.TabletList)-1].maxLine)
			p.currCLDIData.TabletList[len(p.currCLDIData.TabletList)-1].Annotation[line_no] = line
		}

	} else if !strings.Contains(line, "$") && strings.Contains(line, ". ") && !strings.Contains(string(line[0:1]), "#") {
		_, err := strconv.Atoi(line[0:1])
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		data := strings.SplitN(line, ". ", 2)
		line_no := strings.TrimSpace(data[0])
		translit := strings.TrimSpace(data[1])

		if len(p.currCLDIData.TabletList) > 0 {
			p.currCLDIData.TabletList[len(p.currCLDIData.TabletList)-1].TabletLines[line_no] = translit

			//update line number
			line_int, _ := strconv.Atoi(line_no)
			if line_int > p.currCLDIData.TabletList[len(p.currCLDIData.TabletList)-1].maxLine {
				p.currCLDIData.TabletList[len(p.currCLDIData.TabletList)-1].maxLine = line_int
			}
		}
	}
}
