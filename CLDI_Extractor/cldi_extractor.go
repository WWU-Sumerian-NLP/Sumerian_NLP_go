package CLDI_Extractor

type CLDIExtractor struct {
	atfParser *ATFParser
}

func newCLDIExtractor() *CLDIExtractor {
	parser := newATFParser("../../sumerian_tablets/cdli_atf_20220525.txt", "output.tsv")
	return &CLDIExtractor{atfParser: parser}
}
