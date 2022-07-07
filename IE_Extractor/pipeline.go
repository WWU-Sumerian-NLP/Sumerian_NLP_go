package IE_Extractor

//runPipeline will run entire pipeline
func runPipeline(path, destPath string) {
	cdliParser := newCDLIParser(path)
	RelationExtractorRB := newRelationExtractorRB(cdliParser.out)
	dataWriter := newDataWriter(destPath, RelationExtractorRB.out)

	go func() {
		println("running pipeline")
		dataWriter.WaitUntilDone()
		RelationExtractorRB.WaitUntilDone()
		cdliParser.WaitUntilDone()
	}()
}
