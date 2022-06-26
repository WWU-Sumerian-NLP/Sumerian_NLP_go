package CDLI_Extractor

//runPipeline will run entire pipeline
func runPipeline(path, destPath string) {
	atfParser := newATFParser(path)
	atfNormalizer := newATFNormalizer(false, atfParser.out)
	entityExtractor := newCDLIEntityExtractor(atfNormalizer.out)
	RelationExtractorRB := newRelationExtractorRB(entityExtractor.out)
	dataWriter := newDataWriter(destPath, RelationExtractorRB.out)

	go func() {
		println("running pipeline")
		dataWriter.WaitUntilDone()
		RelationExtractorRB.WaitUntilDone()
		entityExtractor.WaitUntilDone()
		atfNormalizer.WaitUntilDone()
		atfParser.WaitUntilDone()
	}()
}

func runCDLIParserPipeline(path, destPath string) {
	atfParser := newATFParser(path)
	dataWriter := newDataWriter(destPath, atfParser.out)

	go func() {
		println("running CDLI Parser")
		dataWriter.WaitUntilDone()
		atfParser.WaitUntilDone()
	}()
}

func runEntityPipeline(path, destPath string) {
	in := readCDLIData(path)
	println("test")
	entityExtractor := newCDLIEntityExtractor(in)
	dataWriter := newDataWriter(destPath, entityExtractor.out)

	go func() {
		println("running entity extraction")
		dataWriter.WaitUntilDone()
		entityExtractor.WaitUntilDone()
	}()
}
