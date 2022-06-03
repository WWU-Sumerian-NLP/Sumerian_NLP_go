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
	entityExtractor := &CLDIEntityExtractor{in: parser.out, out: entitiesParsed, done: make(chan struct{}, 1)}
	parser.loadCLDIData()
	parser.run()
	entityExtractor.run()

	entityExtractor.WaitUntilDone()
	parser.WaitUntilDone()

	for entities := range entitiesParsed {
		fmt.Printf("entities: %v\n", entities.TabletList[0].EntitiyLines)
	}
}

func TestEntityExtractionTestSuite(t *testing.T) {
	suite.Run(t, new(CLDIEntityExtractorTest))
}
