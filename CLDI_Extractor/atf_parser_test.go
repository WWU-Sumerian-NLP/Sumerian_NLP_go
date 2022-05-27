package CLDI_Extractor

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// /Users/hanselguzman-soto/Desktop/urr3-drehem-KG/Sumerian_KG/Data_Collection/sumerian_tablets/cdli_atf_20220525.txt

type ATFParserTestSuite struct {
	suite.Suite
}

func (suite *ATFParserTestSuite) TestSomething() {
	parser := newATFParser("../../sumerian_tablets/cdli_atf_20220525.txt", "../../output.tsv")
	parser.loadCLDIData()
	parser.parseLines()
	parser.exportToCSV()

}

func TestATFParserTestSuite(t *testing.T) {
	suite.Run(t, new(ATFParserTestSuite))
}
