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

/*CDLIData is a structure that encapsulates information about a tablet from CDLI website.

It includes various identifiers, details of the tablet's provenance and period, any referenced dates,
relation tuples, and an array of sections that make up the tablet.

Each line maps to a line_no, transliterations, normalizations and annotations.*/
type CDLIData struct {
	TabletNum       string
	PUB             string
	Provenience     string
	Period          string
	DatesReferenced string
	RelationTuples  string
	TabletSections  []TabletSection
}

// TabletSection represents a distinct section of a tablet.
// It contains information about the location of the section on the tablet, the line numbers in the section,
// the content of each line (including normalized versions of the lines), entity lines, and any associated annotations.
// It also includes a flag to indicate if the section is broken.
type TabletSection struct {
	TabletLocation  string
	LineNumbers     []int //reference each line in case of dropped lines (errors or intelligble)
	TabletLines     []string
	NormalizedLines []string
	EntitiyLines    []string
	Annotations     map[int]string
	isBroken        bool
}

// ATFParser represents a parser for ATF (Assyriological Text Format) data.
// It includes the path to the file being parsed, the raw data, the current CDLIData being constructed,
// output and done channels for signaling completion of the parsing process, and a regex expression.
type ATFParser struct {
	path         string
	data         []string
	currCLDIData CDLIData
	Out          chan CDLIData
	done         chan struct{}
	re           regexp.Regexp
}

// NewATFParser() initializes a new ATFParser with the provided path.
// It then loads the data from the file at the path, compiles the regex used in parsing,
// and starts the parsing process.
func NewATFParser(path string) *ATFParser {
	atfParser := &ATFParser{
		path: path,
		Out:  make(chan CDLIData, 10000000),
		done: make(chan struct{}, 1),
	}
	atfParser.loadCDLIData()
	atfParser.re = *regexp.MustCompile("[0-9]+")
	atfParser.run()
	return atfParser
}

// run() starts the parsing process. It spawns a goroutine that goes through each line in the data,
// parses it, and sends the completed CDLIData to the output channel.
func (p *ATFParser) run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		println("atf parsing")
		defer wg.Done()
		defer close(p.Out)
		for _, line := range p.data {
			p.parseLines(line)
		}
		println("DONE")
		// fmt.Printf("p.currCLDIData: %v\n", p.currCLDIData)
		p.Out <- p.currCLDIData
	}()
	wg.Wait()
}

// WaitUntilDone() sends a signal to the done channel of the parser, indicating the parsing process is complete.
func (p *ATFParser) WaitUntilDone() {
	p.done <- struct{}{}
}

// loadCDLIData() opens the file at the parser's path, reads it line by line, and appends each line to the parser's data.
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

//Parse Tablet Line-By-Line.
//	- &P indicates a new tablet, we initalize p.currCLDIData for a new tablet
//	- Store CLDI, PUB, annotation and transliterations data to this object
//	- @ indicates a new tablet section, we initalize a new TabletSection object
//	- #tr.en indicates a transliteration annotation, we store this annotation to the current TabletSection object
//	- . indicates a new line, we store this line to the current TabletSection object
//	- $ indicates a new line, we store this line to the current TabletSection object
//	- x indicates a new line, we store this line to the current TabletSection object
//	- If we encounter a new tablet, we send the current tablet downstream and initalize a new tablet object
//	- If we encounter a new tablet section, we initalize a new TabletSection object
//	- If we encounter a new line, we store this line to the current TabletSection object
//	- If we encounter a new annotation, we store this annotation to the current TabletSection object
//
//
// parseLines() processes a given line from the file. It recognizes and extracts various types of information
// from the line based on its contents, such as the tablet's primary publication, provenience, period,
// dates referenced, and tablet number. If a new tablet or tablet section is found, the current one is sent
// to the output channel and a new one is initialized. If a line containing tablet content or an annotation is found,
// it is added to the current tablet section.
func (p *ATFParser) parseLines(line string) {
	line = strings.TrimSpace(line)
	// fmt.Printf("line: %v\n", line)

	if line != "" && strings.Contains(line, "Primary publication") {
		// Send tablet downstream and init new tablet object
		if p.currCLDIData.PUB != "" {
			p.Out <- p.currCLDIData
			fmt.Printf("p.currCLDIData: %v\n", p.currCLDIData)
			p.currCLDIData = CDLIData{}
		}
	} else if line != "" && strings.Contains(line, "Provenience") {
		line = strings.TrimSpace(line)
		data := strings.Split(line, ":")[1]
		p.currCLDIData.Provenience = data
	} else if line != "" && strings.Contains(line, "Period") {
		line = strings.TrimSpace(line)
		data := strings.Split(line, ":")[1]
		p.currCLDIData.Period = data
	} else if line != "" && strings.Contains(line, "Dates referenced") {
		line = strings.TrimSpace(line)
		data := strings.Split(line, ":")[1]
		p.currCLDIData.DatesReferenced = data

	} else if line != "" && strings.Contains(line, "&P") && strings.Contains(line, " =") {
		// TODO: Send tablet downstream and init new tablet object with alternative format
		if p.currCLDIData.PUB != "" {
			p.Out <- p.currCLDIData
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
		newTableLine.TabletLines = make([]string, 0)

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
			line_no := len(currTabletSections[len(currTabletSections)-1].TabletLines) - 1
			currTabletSections[len(currTabletSections)-1].Annotations[line_no] = line
		}

		// Get tablet lines
	} else if !strings.Contains(line, "$") && strings.Contains(line, ". ") && !strings.Contains(string(line[0:1]), "#") && !strings.Contains(string(line[0:1]), "x") {
		_, err := strconv.Atoi(line[0:1])
		if err != nil {
			fmt.Printf("err: %v for %s\n", err, line)
		} else { // fix later

			data := strings.SplitN(line, ". ", 2)
			line_no := findLineNumber(&p.re, data[0])
			translit := strings.TrimSpace(data[1])

			currTabletSections := p.currCLDIData.TabletSections
			if len(currTabletSections) > 0 {
				currTabletSections[len(currTabletSections)-1].LineNumbers = append(currTabletSections[len(currTabletSections)-1].LineNumbers, line_no)
				currTabletSections[len(currTabletSections)-1].TabletLines = append(currTabletSections[len(currTabletSections)-1].TabletLines, translit)
			}
		}

	}
}

// NotTabletSubsection() checks if the given line does not represent a subsection of a tablet. It determines this by checking
// if the line does not contain certain keywords such as "object", "tablet", "envelope", or "bulla". It returns true
// if none of these keywords are present, and false otherwise.
//
// @ sign may denote a tablet secton, but in cases below they do not. Rather they are subsections
func NotTabletSubsection(line string) bool {
	return !strings.Contains(line, "object") && !strings.Contains(line, "tablet") &&
		!strings.Contains(line, "envelope") && !strings.Contains(line, "bulla")
}

// findLineNumber() extracts the line number from the given string. The line number is assumed to be an integer embedded within the string.
// If the line number cannot be extracted or converted to an integer, an error message is printed and a zero value is returned.
//
//Line number comes in the form of strings like (1.), ('1.) which needs to be converted to an int
func findLineNumber(re *regexp.Regexp, line string) int {
	intString := re.FindString(line)
	lineNumber, err := strconv.Atoi(intString)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return lineNumber
}
