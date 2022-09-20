package IE_Extractor

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type PipelineTestSuite struct {
	suite.Suite
}

func (suite *PipelineTestSuite) TestPipeline() {
	// runPipeline("test_data/entity_extraction_input.atf", "output/test_data.tsv")
	// runPipeline("../CDLI_Extractor/output/urr3_no_annotations.tsv", "output/urr3_ie_no_annotations.tsv")
	runPipeline("../CDLI_Extractor/output/urr3_annotations.tsv", "output/urr3_ie_annotations.tsv")
	// runPipeline("../CDLI_Extractor/output/all_tablets_data.tsv", "output/all_ie_data.tsv")
	// runPipeline("../CDLI_Extractor/output/unblocked_data.tsv", "output/unblocked_ie_data.tsv")

}

func TestPipelineTestSuite(t *testing.T) {
	suite.Run(t, new(PipelineTestSuite))
}
