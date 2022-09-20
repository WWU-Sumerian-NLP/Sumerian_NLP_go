package IE_Extractor

//runPipeline will run entire pipeline
func runPipeline(path, destPath string) {
	cdliParser := NewCDLIParser(path)
	RelationExtractorRB := NewRelationExtractorRB(cdliParser.Out)
	dataWriter := NewDataWriter(destPath, RelationExtractorRB.Out)

	go func() {
		println("running pipeline")
		dataWriter.WaitUntilDone()
		RelationExtractorRB.WaitUntilDone()
		cdliParser.WaitUntilDone()
	}()
}
