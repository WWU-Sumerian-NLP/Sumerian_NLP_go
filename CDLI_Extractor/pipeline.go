package CDLI_Extractor

func runPipeline(path, destPath string) {
	atfParser := newATFParser(path)
	atfNormalizer := newATFNormalizer(false, atfParser.out)
	entityExtractor := newCDLIEntityExtractor(atfNormalizer.out)
	dataWriter := newDataWriter(destPath, entityExtractor.out)

	// // time.Sleep(2 * time.Second)
	go func() {
		println("running pipeline")
		dataWriter.WaitUntilDone()
		entityExtractor.WaitUntilDone()
		atfNormalizer.WaitUntilDone()
		atfParser.WaitUntilDone()

	}()
}
