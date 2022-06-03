package CLDI_Extractor

import "time"

func runPipeline(path, destPath string) {
	atfParser := newATFParser(path)
	atfNormalizer := newATFNormalizer(false, atfParser.out)
	entityExtractor := newCLDIEntityExtractor(atfNormalizer.out)
	dataWriter := newDataWriter(destPath, entityExtractor.out)

	go func() {
		println("running pipeline")
		dataWriter.WaitUntilDone()
		entityExtractor.WaitUntilDone()
		atfNormalizer.WaitUntilDone()
		atfParser.WaitUntilDone()
	}()
	time.Sleep(time.Second * 20) //todo - synchronize channels better
}
