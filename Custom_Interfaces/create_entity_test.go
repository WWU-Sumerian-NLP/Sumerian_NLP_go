package custom_interfaces

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CustomEntityTest struct {
	suite.Suite
}

func (suite *CustomEntityTest) TestCustomEntity() {
	//Create new entity
	entityName := "Shovel"
	entityTag := "TOOLS"
	pathToList := ""
	testEntity := newCustomEntity(entityName, entityTag, pathToList)

	fmt.Printf("testEntity: %v\n", testEntity)

	// Parse ATF test files and extract entities

	// parsedLines := make(chan CDLIData, 10)
	// entitiesParsed := make(chan CDLIData, 10)
	// parser := &ATFParser{path: "test_data/entity_extraction_input.atf", out: parsedLines, done: make(chan struct{}, 1)}
	// entityExtractor := &CDLIEntityExtractor{in: parser.out, out: entitiesParsed, done: make(chan struct{}, 1)}
	// parser.loadCDLIData()
	// parser.run()
	// entityExtractor.run()

	// entityExtractor.WaitUntilDone()
	// parser.WaitUntilDone()

	// for entities := range entitiesParsed {
	// 	fmt.Printf("entities: %v\n", entities.TabletSections[0].EntitiyLines)
	// }

	//Check to see if new entities are extracted
}

func TestCustomEntityTestSuite(t *testing.T) {
	suite.Run(t, new(CustomEntityTest))
}
