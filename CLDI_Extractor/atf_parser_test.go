package CLDI_Extractor

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ATFParserTestSuite struct {
	suite.Suite
}

func (suite *ATFParserTestSuite) TestATFParser() {
	parser := newATFParser("../../sumerian_tablets/cdli_atf_20220525.txt", "../../cdli_atf_20220525.tsv")
	parser.loadCLDIData()
	parser.parseLines()
	parser.exportToCSV()

}

func TestATFParserTestSuite(t *testing.T) {
	suite.Run(t, new(ATFParserTestSuite))
}
