package CDLI_Extractor

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

//Lines map to a line_no, transliterations, normalizations and annotations
type CDLIData struct {
	TabletNum      string
	PUB            string
	TabletSections []TabletSection
}

type TabletSection struct {
	TabletLocation  string
	TabletLines     map[int]string //map line_no to transliterations
	NormalizedLines map[int]string
	EntitiyLines    map[int]string
	RelationLines   map[int]string
	Annotations     map[int]string //map line_no to annotations
	maxLine         int
}

type ATFParser struct {
	path         string
	data         []string
	currCLDIData CDLIData
	out          chan CDLIData
	done         chan struct{}
	re           regexp.Regexp
}

func newATFParser(path string) *ATFParser {
	atfParser := &ATFParser{
		path: path,
		out:  make(chan CDLIData, 100000),
		done: make(chan struct{}, 1),
	}
	atfParser.loadCDLIData()
	atfParser.re = *regexp.MustCompile("[0-9]+")
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
			p.parseLines(line)
		}
		println("DONE")
		p.out <- p.currCLDIData
	}()
	wg.Wait()
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
		// Send tablet downstream and init new tablet object
		if p.currCLDIData.PUB != "" {
			p.out <- p.currCLDIData
			p.currCLDIData = CDLIData{}
		}
		line = strings.ReplaceAll(line, "&", "")
		line = strings.TrimSpace(line)
		data := strings.SplitN(line, " = ", 2)
		p.currCLDIData.TabletNum = data[0]
		p.currCLDIData.PUB = data[1]

		// Create new tablet section
	} else if line != "" && strings.Contains(string(line[0:1]), "@") && NotTabletSubsection(line) {
		newTableLine := TabletSection{}
		newTableLine.TabletLocation = strings.TrimSpace(line)
		newTableLine.TabletLines = make(map[int]string)

		p.currCLDIData.TabletSections = append(p.currCLDIData.TabletSections, newTableLine)
		p.currCLDIData.TabletSections[len(p.currCLDIData.TabletSections)-1].Annotations = make(map[int]string)

		// Get tablet annotations
	} else if strings.Contains(line, "#tr.en") {
		// You can translate tr.en entries
		line = strings.Replace(line, "#tr.en", "", 1)
		line = strings.Replace(line, ":", "", 1)
		line := strings.TrimSpace(line)

		currTabletSections := p.currCLDIData.TabletSections
		if len(currTabletSections) > 0 {
			line_no := currTabletSections[len(currTabletSections)-1].maxLine
			currTabletSections[len(currTabletSections)-1].Annotations[line_no] = line
		}

		// Get tablet lines
	} else if !strings.Contains(line, "$") && strings.Contains(line, ". ") && !strings.Contains(string(line[0:1]), "#") {
		_, err := strconv.Atoi(line[0:1])
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		data := strings.SplitN(line, ". ", 2)
		line_no := findLineNumber(&p.re, data[0])
		translit := strings.TrimSpace(data[1])

		currTabletSections := p.currCLDIData.TabletSections
		if len(currTabletSections) > 0 {
			currTabletSections[len(currTabletSections)-1].TabletLines[line_no] = translit

			//update line number
			if line_no > currTabletSections[len(currTabletSections)-1].maxLine {
				currTabletSections[len(currTabletSections)-1].maxLine = line_no
			}
		}
	}
}

// @ sign may denote a tablet secton, but in cases below they do not. Rather they are subsections
func NotTabletSubsection(line string) bool {
	return !strings.Contains(line, "object") && !strings.Contains(line, "tablet") &&
		!strings.Contains(line, "envelope") && !strings.Contains(line, "bulla")
}

//Line number comes in the form of strings like (1.), ('1.) which needs to be converted to an int
func findLineNumber(re *regexp.Regexp, line string) int {
	intString := re.FindString(line)
	lineNumber, err := strconv.Atoi(intString)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return lineNumber

}
