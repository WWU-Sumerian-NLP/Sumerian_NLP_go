package CLDI_Extractor

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
)

//Lines map to a line_no, transliterations and normalizations
type CLDIData struct {
	CLDI        string
	PUB         string
	annotations string
	lines       map[string]string
}

type ATFParser struct {
	path     string
	destPath string
	data     []string
	CLDIList []CLDIData
}

func newATFParser(path string, destPath string) *ATFParser {
	return &ATFParser{path: path, destPath: destPath}
}

func (p *ATFParser) loadCLDIData() {
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
4 cases

1 - Skip empty lines
2 - Tablet information denotated by "&P"
3 - Tablet annotations/translation denotated by "#tr.en"
4 - Tablet line and writings ex: "1 _mu_ qa-ba-ra-a{ki}"
*/

func (p *ATFParser) parseLines() {
	cldiData := &CLDIData{}
	for _, line := range p.data {
		p.CLDIList = append(p.CLDIList, *cldiData)
		line = strings.TrimSpace(line)
		if line != "" && strings.Contains(line, "&P") && strings.Contains(line, " =") {
			line = strings.ReplaceAll(line, "&", "")
			line = strings.TrimSpace(line)
			data := strings.SplitN(line, " = ", 2)
			cldiData.CLDI = data[0]
			cldiData.PUB = data[1]
		} else if strings.Contains(line, "#tr.en") {
			// You can translate tr.en entries
			line = strings.Replace(line, "#tr.en", "", 1)
			line = strings.Replace(line, ":", "", 1)
			cldiData.annotations = line
		} else if strings.Contains(line, ". ") {
			_, err := strconv.Atoi(line[0:1])
			if err != nil {
				continue
			}
			data := strings.SplitN(line, ". ", 2)
			cldiData.lines = make(map[string]string)
			cldiData.lines["no"] = strings.TrimSpace(data[0])
			cldiData.lines["translit"] = strings.TrimSpace(data[1])
			cldiData.lines["normalized_translit"] = strings.TrimSpace(data[1])
		} else {
			continue
		}

	}
}

func (p *ATFParser) exportToCSV() {
	csvFile, err := os.Create(p.destPath)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvWriter := csv.NewWriter(csvFile)
	csvWriter.Comma = '\t'
	csvWriter.Write([]string{"CLDI", "PUB", "no", "translit", "annotations"}) //hardcoded
	for _, CLDIData := range p.CLDIList {
		csvWriter.Write([]string{CLDIData.CLDI, CLDIData.PUB, CLDIData.lines["no"], CLDIData.lines["translit"], CLDIData.annotations})
	}
	csvFile.Close()
}
