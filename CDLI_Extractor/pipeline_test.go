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
	runPipeline("../../sumerian_tablets/cdli_atf_20220525.txt", "output/new_pipeline.tsv")
	// runCDLIParserPipeline("../../sumerian_tablets/cdli_atf_20220525.txt", "output/parsed_cdli.tsv")
	// runEntityPipeline("output/new_pipeline.tsv", "output/new_entity_pipeline.tsv")

}

func TestPipelineTestSuite(t *testing.T) {
	suite.Run(t, new(PipelineTestSuite))
}
