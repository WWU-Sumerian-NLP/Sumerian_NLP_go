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

	providence      string
	datesReferenced string
}

func newRelationData(relationType string, regexRules string, subject string, object string) *RelationData {
	return &RelationData{relationType: relationType, regexRules: regexRules, subjectTag: subject, objectTag: object}
}

// Relation
// person:delivers_animal
// subject - PN
// object - ANIM
// Expected: RelationTuple - (person:delivers_animals, person_word, animal_word)

// extractedRegexTuple - (PN, ANIM DEL) ---> (DEL, PN, ANIM)

func (r *RelationData) getRelationTuple(extractedRegexTuple []string, transliterationTuple []string) [3]string {
	relationTuple := [3]string{}

	//make function
	translitMap := make(map[string]int, 0)
	for i, translit := range transliterationTuple {
		println("HERE", translit, i)
		translitMap[translit] = i
	}
	fmt.Printf("transliterationTuple: %v\n", transliterationTuple)
	relationTuple[0] = r.relationType
	fmt.Printf("extractedTuple: %v\n", extractedRegexTuple)
	for i, tag := range extractedRegexTuple {
		tag = strings.TrimSpace(tag)
		r.subjectTag = strings.TrimSpace(r.subjectTag) //fix later
		r.objectTag = strings.TrimSpace(r.objectTag)
		if tag == r.subjectTag {
			relationTuple[1] = transliterationTuple[i]
		} else if tag == r.objectTag {
			relationTuple[2] = transliterationTuple[i]
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
