package CLDI_Extractor

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type PipelineTestSuite struct {
	suite.Suite
}

func (suite *PipelineTestSuite) TestPipeline() {
	runPipeline("../../sumerian_tablets/cdli_atf_20220525.txt", "pipeline.tsv")

}

func TestPipelineTestSuite(t *testing.T) {
	suite.Run(t, new(PipelineTestSuite))
}
