package CDLI_Extractor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CLDIEntityExtractorTest struct {
	suite.Suite
}

func (suite *CLDIEntityExtractorTest) TestEntityExtraction() {
	parsedLines := make(chan CDLIData, 10)
	entitiesParsed := make(chan CDLIData, 10)
	parser := &ATFParser{path: "test_data/entity_extraction_input.atf", Out: parsedLines, done: make(chan struct{}, 1)}
	entityExtractor := &CDLIEntityExtractor{in: parser.Out, Out: entitiesParsed, done: make(chan struct{}, 1)}
	parser.loadCDLIData()
	parser.run()
	entityExtractor.run()

	entityExtractor.WaitUntilDone()
	parser.WaitUntilDone()

	for entities := range entitiesParsed {
		fmt.Printf("entities: %v\n", entities.TabletSections[0].EntitiyLines)
	}
}

func TestEntityExtractionTestSuite(t *testing.T) {
	suite.Run(t, new(CLDIEntityExtractorTest))
}
