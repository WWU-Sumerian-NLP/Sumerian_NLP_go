package CDLI_Extractor

// runPipeline will run entire pipeline
func runPipeline(path, destPath string) {
	atfParser := NewATFParser(path, false)
	transliterationCleaner := NewTransliterationCleaner(false, atfParser.Out)
	atfNormalizer := newATFNormalizer(false, transliterationCleaner.Out)
	entityExtractor := NewCDLIEntityExtractor(atfNormalizer.out)
	dataWriter := NewDataWriter(destPath, entityExtractor.Out)

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
	dataWriter := NewDataWriter(destPath, atfParser.Out)

	go func() {
		println("running CDLI Parser")
		dataWriter.WaitUntilDone()
		atfParser.WaitUntilDone()
	}()
}

func runEntityPipeline(path, destPath string) {
	in := readCDLIData(path)
	println("test")
	entityExtractor := NewCDLIEntityExtractor(in)
	dataWriter := NewDataWriter(destPath, entityExtractor.Out)

	go func() {
		println("running entity extraction")
		dataWriter.WaitUntilDone()
		entityExtractor.WaitUntilDone()
	}()
}
