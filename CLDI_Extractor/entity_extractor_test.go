package CLDI_Extractor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CLDIEntityExtractorTest struct {
	suite.Suite
}

func (suite *CLDIEntityExtractorTest) TestEntityExtraction() {
	parsedLines := make(chan CLDIData, 10)
	entitiesParsed := make(chan CLDIData, 10)
	parser := &ATFParser{path: "test_data/entity_extraction_input.atf", out: parsedLines, done: make(chan struct{}, 1)}
	parser.run()
	parser.WaitUntilDone()
	entityExtractor := &CLDIEntityExtractor{in: parsedLines, out: entitiesParsed, done: make(chan struct{}, 1)}
	entityExtractor.run()
	entityExtractor.WaitUntilDone()

	// close(outputChan)

	for test := range entitiesParsed {
		fmt.Printf("test: %v\n", test)
	}

}

func TestEntityExtractionTestSuite(t *testing.T) {
	suite.Run(t, new(CLDIEntityExtractorTest))
}
