package CDLI_Extractor

//runPipeline will run entire pipeline
func runPipeline(path, destPath string) {
	atfParser := NewATFParser(path)
	transliterationCleaner := newTransliterationCleaner(false, atfParser.out)
	atfNormalizer := newATFNormalizer(false, transliterationCleaner.out)
	entityExtractor := newCDLIEntityExtractor(atfNormalizer.out)
	dataWriter := newDataWriter(destPath, entityExtractor.out)

	go func() {
		println("running pipeline")
		dataWriter.WaitUntilDone()
		entityExtractor.WaitUntilDone()
		atfNormalizer.WaitUntilDone()
		transliterationCleaner.WaitUntilDone()
		atfParser.WaitUntilDone()
	}()
}

func runCDLIParserPipeline(path, destPath string) {
	atfParser := NewATFParser(path)
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
