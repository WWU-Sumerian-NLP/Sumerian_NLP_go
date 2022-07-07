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
	in            <-chan TaggedTransliterations
	out           chan TaggedTransliterations
	done          chan struct{}
	regexRuleList []string
	re            *regexp.Regexp
	re2           *regexp.Regexp
}

func newRelationExtractorRB(in <-chan TaggedTransliterations) *RelationExtractorRB {
	relationExtractor := &RelationExtractorRB{
		in:            in,
		out:           make(chan TaggedTransliterations, 1000000),
		regexRuleList: make([]string, 0),
		done:          make(chan struct{}, 1),
	}
	relationExtractor.regexRuleList = []string{`ANIM PN [O\s?]*DEL`, `ANIM [O\s?]* PN REC`}
	re, _ := regexp.Compile(findParenthesis)
	re_2, _ := regexp.Compile(findInnerParenthesis)
	relationExtractor.re = re
	relationExtractor.re2 = re_2

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
			for _, regexRule := range r.regexRuleList {
				fmt.Printf("regexRule: %v\n", regexRule)
				r.extractFromRegexRules(cdliData.taggedTranslit, regexRule)
			}
			r.out <- cdliData
		}
		println("DONE")
		close(r.out)
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

	desiredTagSequence := findPat.FindAllString(strings.Join(tagList, " "), 10)
	desiredTagSequence = strings.Split(strings.Join(desiredTagSequence, " "), " ")
	fmt.Printf("desiredTagSequence: %v\n", desiredTagSequence)

	if len(desiredTagSequence) > 1 {
		finalList = r.findRegexMatchFromTagSequence(desiredTagSequence, tagList, graphemeWithTag)
	}
	fmt.Printf("finalList: %v\n", finalList)
	return finalList

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
*/
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

		} else {
			pos = 0
			tupleList = []string{}
		}

		if pos == len(desiredTagSequence)-1 {
			new_tag := strings.Split(graphemeWithTag[i+1], ",")[0]
			tupleList = append(tupleList, new_tag[1:]) //weird issue of not adding last tag
			finalList = append(finalList, strings.Join(tupleList, " "))
			break
		}
	}
	return finalList
}
