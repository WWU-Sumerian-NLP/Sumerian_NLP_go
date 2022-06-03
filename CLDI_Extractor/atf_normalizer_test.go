package CLDI_Extractor

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

// /Users/hanselguzman-soto/Desktop/urr3-drehem-KG/Sumerian_KG/Data_Collection/sumerian_tablets/cdli_atf_20220525.txt

type ATFNormalizerTestSuite struct {
	suite.Suite
}

func (suite *ATFNormalizerTestSuite) TestATFNormalizer() {
	transliterationString := "[mu ha-ar]-szi#{ki} u3 [ki-masz{ki}] ba-hul"
	print(transliterationString)
	// test2 := "4(disz) gu4 amar ga 5(disz) ab2 [amar ga]"
}

func (suite *ATFNormalizerTestSuite) TestATFStandardizer() {
	transliterationString := "[mu ha-ar]-szi#{ki} u3 [ki-masz{ki}] ba-hul"
	expectedString := "[mu ha-ar]-ci#{ki} u3 [ki-mac{ki}] ba-hul"

	normalizedString := ""
	for _, grapheme := range strings.Split(transliterationString, " ") {
		normalizedGrapheme := standardizeGrapheme(grapheme)
		normalizedString += normalizedGrapheme + " "
	}
	normalizedString = strings.TrimSpace(normalizedString)
	suite.Assert().Equal(expectedString, normalizedString)

}

func TestATFNormalizerTestSuite(t *testing.T) {
	suite.Run(t, new(ATFNormalizerTestSuite))
}
