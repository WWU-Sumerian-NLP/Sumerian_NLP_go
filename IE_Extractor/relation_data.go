package IE_Extractor

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type RelationData struct {
	tabletNum     string
	relationType  string
	regexRules    string
	subjectTag    string
	objectTag     string
	tags          string    //how they would be extracted (ANIM PN DEL)
	relationTuple [3]string //3-tuple (relationType, subjectString, objectString)
}

func newRelationData(relationType string, regexRules string, subject string, object string) *RelationData {
	return &RelationData{relationType: relationType, regexRules: regexRules, subjectTag: subject, objectTag: object}
}

// Relation
// person:delivers_animal
// subject - PN
// object - ANIM
// Expected: RelationTuple - (person:delivers_animals, person_word, animal_word)

// ExtractedTuple - (PN, ANIM DEL)
// TransliterationTuple - (person_word, animal_word, mu-kux(DU))
//Refine this algorithm to be more in general
//we should loop better
func (r *RelationData) getRelationTuple(extractedTuple []string, transliterationTuple []string) [3]string {
	relationTuple := [3]string{}

	relationTuple[0] = r.relationType
	fmt.Printf("extractedTuple: %v\n", extractedTuple)
	for _, tag := range extractedTuple {
		tag = strings.TrimSpace(tag)
		r.subjectTag = strings.TrimSpace(r.subjectTag) //fix later
		r.objectTag = strings.TrimSpace(r.objectTag)
		if tag == r.subjectTag {
			relationTuple[1] = transliterationTuple[1]
		} else if tag == r.objectTag {
			relationTuple[2] = transliterationTuple[0]
		}
	}
	fmt.Printf("relationTuple: %v\n", relationTuple)
	return relationTuple
}

func readRelationTypesCsv(pathToRelationTypesCSV string) []RelationData {
	csvFile, err := os.Open(pathToRelationTypesCSV)
	if err != nil {
		log.Fatalf("failed reading file: %s", err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = '\t'

	var relationDataList []RelationData
	i := 0
	for {
		var data RelationData
		csvData, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("issue here")
			log.Fatal(err)
		}
		data.relationType = csvData[0]
		data.subjectTag = csvData[1]
		data.objectTag = csvData[2]
		data.regexRules = csvData[3]
		data.tags = csvData[4]
		if i != 0 {
			relationDataList = append(relationDataList, data)
		}
		i += 1
	}
	return relationDataList
}
