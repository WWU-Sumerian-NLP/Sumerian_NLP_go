package IE_Extractor

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

const findParenthesis = `\([^)]*\)|\[[^\]]*\]g`
const findInnerParenthesis = `\([^DU(,)]*\)|\[[^\]]*\]g` //regex gets inner para like 1(disz) except DEL

type RelationExtractorRB struct {
	in               <-chan TaggedTransliterations
	Out              chan RelationData
	done             chan struct{}
	regexRuleList    []string
	re               *regexp.Regexp
	re2              *regexp.Regexp
	RelationDataList []RelationData
}

func NewRelationExtractorRB(in <-chan TaggedTransliterations) *RelationExtractorRB {
	relationExtractor := &RelationExtractorRB{
		in:            in,
		Out:           make(chan RelationData, 10000000),
		regexRuleList: make([]string, 0),
		done:          make(chan struct{}, 1),
	}
	relationExtractor.regexRuleList = []string{`ANIM [O\s?]*PN [O\s?]*DEL`, `ANIM [O\s?]*PN [O\s?]*REC`, `ANIM [O\s?]*FOR [O\s?]*DEL`}
	re, _ := regexp.Compile(findParenthesis)
	re_2, _ := regexp.Compile(findInnerParenthesis)
	relationExtractor.re = re
	relationExtractor.re2 = re_2

	relationExtractor.RelationDataList = readRelationTypesCsv("tests/relation_input.tsv")
	relationExtractor.run()
	return relationExtractor
}

func (r *RelationExtractorRB) WaitUntilDone() {
	r.done <- struct{}{}
}

func (r *RelationExtractorRB) run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		println("relation extraction")
		defer wg.Done()
		for cdliData := range r.in {
			for _, relationData := range r.RelationDataList { //read through relation_input.csv
				relationData.tabletNum = cdliData.TabletNum
				relationData.providence = cdliData.Providence
				relationData.period = cdliData.Period
				relationData.datesReferenced = cdliData.DateReferenced

				extractedRelationTuple := r.extractFromRegexRules(cdliData.taggedTranslit, relationData.regexRules)
				if len(extractedRelationTuple) >= 3 {
					fmt.Printf("extractedRelationTuple: %v\n", extractedRelationTuple)
					relationTuple := relationData.getRelationTuple(strings.Split(relationData.tags, ","), extractedRelationTuple)
					fmt.Printf("relationData.relationTuple: %v\n", relationData.relationTuple)
					relationData.relationTuple = relationTuple
					r.Out <- relationData

				}
			}
		}
		println("DONE")
		close(r.Out)
	}()
	wg.Wait()
}

func (r *RelationExtractorRB) extractFromRegexRules(allTabletLines string, regexRule string) []string {
	//temporarily filter out mu-kux(DU) to remove inner parathesis
	firstPass := strings.ReplaceAll(allTabletLines, "mu-kux(DU)", "mu-kux-DU")

	//Replace inner parathesis except those with characters [DU] (small error for now) 1(disz) gets mapped to 1
	secondPass := r.re2.ReplaceAllString(firstPass, "")

	//regex expression to split by parathesis
	graphemeWithTag := r.re.FindAllString(secondPass, 100)

	tagList := r.createTagList(graphemeWithTag)
	fmt.Printf("tagList: %v\n", tagList)

	finalList := []string{}

	findPat, _ := regexp.Compile(regexRule)
	fmt.Printf("findPat: %v\n", findPat)

	desiredTagSequence := findPat.FindAllString(strings.Join(tagList, " "), 10)
	desiredTagSequence = strings.Split(strings.Join(desiredTagSequence, " "), " ")
	fmt.Printf("desiredTagSequence: %v\n", desiredTagSequence)

	if len(desiredTagSequence) > 1 {
		finalList = r.findRegexMatchFromTagSequence(desiredTagSequence, tagList, graphemeWithTag)
	}
	finalListModified := strings.Split(strings.Join(finalList, " "), " ")
	fmt.Printf("finalListModified: %v\n", finalListModified)
	println("LENGTH:", len(finalListModified))
	return finalListModified

}

func (r *RelationExtractorRB) createTagList(graphemeWithTag []string) []string {
	tagList := []string{}
	for _, tag := range graphemeWithTag {
		if strings.Contains(tag, ",") { //some erros with inner parathesis, which is why this exists, it filters 1(disz)
			newTag := strings.Split(tag, ",")[1]
			tagList = append(tagList, newTag[0:len(newTag)-1]) //Splitting on "," yield [(object,], [tag,)], we don't want parathesis
		} else {
			tagList = append(tagList, "O")
		}
	}
	return tagList
}

/*
	String matching algorithm - Given a desired tag sequence found from our regex expression
	We want to trace the tagList to find this sequence. As we find the seqequence, we start building
	a slice based on corresponding graphemes. If successful, we add this sequence else we reset

	Tail Head Relation
*/

//PN ANIM DEL
func (r *RelationExtractorRB) findRegexMatchFromTagSequence(desiredTagSequence []string, tagList []string, graphemeWithTag []string) []string {
	tupleList := []string{}
	finalList := []string{}

	pos := 0
	for i, tag := range tagList {
		fmt.Printf("pos: %v\n", pos)
		fmt.Printf("tag: %v\n", tag)
		fmt.Printf("desiredTagSequence[pos]: %v\n", desiredTagSequence[pos])

		if tag == desiredTagSequence[pos] {
			pos += 1
			if tag != "O" {
				new_tag := strings.Split(graphemeWithTag[i], ",")[0]
				tupleList = append(tupleList, new_tag[1:])
			}

		} else { //reset sequence because it is not an exact match
			pos = 0
			tupleList = []string{}
		}
		if pos == len(desiredTagSequence) {
			// new_tag := strings.Split(graphemeWithTag[i+1], ",")[0]
			// tupleList = append(tupleList, new_tag[1:]) //TODO: weird issue of not adding last tag
			finalList = append(finalList, strings.Join(tupleList, " "))
			break
		}
	}
	return finalList
}
