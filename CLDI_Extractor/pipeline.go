package CLDI_Extractor

import (
	"time"
)

func runPipeline(path, destPath string) {
	atfParser := newATFParser(path, "")
	atfNormalizer := newATFNormalizer(false, atfParser.out)
	entityExtractor := newCLDIEntityExtractor(atfNormalizer.out)
	dataWriter := newDataWriter(destPath, entityExtractor.out)

	//does nothing for now
	go func() {
		println("finishing up")
		atfParser.WaitUntilDone()
		atfNormalizer.WaitUntilDone()
		dataWriter.WaitUntilDone()
		entityExtractor.WaitUntilDone()

	}()
	time.Sleep(time.Second * 5)
}
