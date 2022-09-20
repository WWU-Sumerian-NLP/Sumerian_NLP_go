package CDLI_Extractor

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type PipelineTestSuite struct {
	suite.Suite
}

func (suite *PipelineTestSuite) TestPipeline() {
	// runPipeline("test_data/entity_extraction_input.atf", "output/test_data.tsv")
	// runPipeline("../../sumerian_tablets/cdli_atf_20220525.txt", "output/urr3_no_annotations.tsv")
	runPipeline("../../sumerian_tablets/cdli_result_20220525.txt", "output/urr3_annotations.tsv")
	// runPipeline("../../sumerian_tablets/ur3_20110805_public.atf", "output/all_tablets_data.tsv")
	// runPipeline("../../sumerian_tablets/cdliatf_unblocked.atf", "output/unblocked_data.tsv")

}

func TestPipelineTestSuite(t *testing.T) {
	suite.Run(t, new(PipelineTestSuite))
}
