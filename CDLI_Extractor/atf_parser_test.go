package CDLI_Extractor

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ATFParserTestSuite struct {
	suite.Suite
}

func (suite *ATFParserTestSuite) TestATFParser() {
	parsedLines := make(chan CDLIData, 10)
	parser := &ATFParser{path: "test_data/atf_test_input.atf", Out: parsedLines, done: make(chan struct{}, 1)}
	parser.loadCDLIData()
	parser.run()
	parser.WaitUntilDone()

	givenTabletNum := ""
	givenPub := ""
	for item := range parsedLines {
		givenTabletNum = item.TabletNum
		givenPub = item.PUB

	}
	expectedTabletNum := "P142761"
	expectPub := "AAICAB 1/1, pl. 048, 1911-488"

	suite.Assert().Equal(givenTabletNum, expectedTabletNum)
	suite.Assert().Equal(givenPub, expectPub)

}

func TestATFParserTestSuite(t *testing.T) {
	suite.Run(t, new(ATFParserTestSuite))
}
