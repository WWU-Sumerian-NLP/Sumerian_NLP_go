package IE_Extractor

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

const findParenthesis = `\([^)]*\)|\[[^\]]*\]g`

// const findParenthesis = `\(*,[^)]*\)|\[[^\]]*\]g` //regex gets ,O where O is a tag
// const findParenthesis = `\(*[^,)]*\)|\[[^\]]*\]g` //regex gets

const findInnerParenthesis = `\([^DU(,)]*\)|\[[^\]]*\]g` //regex gets inner para like 1(disz) except DEL

type RelationExtractorRB struct {
	in   <-chan TaggedTransliterations
	out  chan TaggedTransliterations
	done chan struct{}
	re   *regexp.Regexp
	re2  *regexp.Regexp
}

func newRelationExtractorRB(in <-chan TaggedTransliterations) *RelationExtractorRB {
	relationExtractor := &RelationExtractorRB{
		in:   in,
		out:  make(chan TaggedTransliterations, 1000000),
		done: make(chan struct{}, 1),
	}
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
			cdliData.deliveryRelations = r.getFromRules(cdliData.taggedTranslit)
			cdliData.recievedRelations = r.getFromRecRules(cdliData.taggedTranslit)
			r.out <- cdliData
		}
		println("DONE")
		close(r.out)
	}()
	wg.Wait()
}

/*
Evaluate an entire tablet. Parse (sumerian_graphme, NER_label) - function to extract NER_labels = everything must have a label

If we find this sequence, lets extract the entire relation

ANIM
PN
DEL
PN REC

Assuming PN1 != PN2

Iterate through (, )
1. Check to see if there is a comma (grapheme, O)
2. If there is a comma, then split by commma [grapheme, O]
3. Put all tags in a list [O, O, ANIM, O, PN, DEL, PN, REC]


*/
func (r *RelationExtractorRB) getFromRules(allTabletLines string) []string {
	// fmt.Printf("allTabletLines: %v\n", allTabletLines)
	//temporarily filter out strings like  1(disz)
	tabletLines := strings.ReplaceAll(allTabletLines, "mu-kux(DU)", "mu-ku-DU") //temporary
	r.re2.ReplaceAllString(tabletLines, "")
	inTag := r.re.FindAllString(tabletLines, 100)
	fmt.Printf("inTag: %v\n", inTag)
	// fmt.Printf("inTag: %v\n", inTag)
	tagList := []string{}
	for _, tag := range inTag { //filters out items like 1(disz)
		if strings.Contains(tag, ",") {
			new_tag := strings.Split(tag, ",")[1]
			tagList = append(tagList, new_tag[0:len(new_tag)-1])
		} else {
			tagList = append(tagList, "O")
		}
	}
	fmt.Printf("tagList: %v\n", tagList)
	/*
	   mu du (delivery) rule

	   ANIM
	   Person 1 PN
	   mu-kux(DU) DEL
	   Person 2 PN i3-dab5 REC
	*/
	finalList := []string{}
	tupleList := []string{}

	//ANIM PN DEL
	regexAbove := `ANIM PN [O\s?]*DEL`
	findPat, _ := regexp.Compile(regexAbove)

	//PN REC ANIM

	// test := findPat.FindAllString(strings.Join(tagList, " "), 100)
	// fmt.Printf("test: %v\n", test)
	//Person delivered animal

	// (ANIM, 0), (ANIM, 1), (PN, 2), (O, 3), (DEL, 4)

	// println("HERE")
	desiredSequence := findPat.FindAllString(strings.Join(tagList, " "), 10)
	desiredSequence = strings.Split(strings.Join(desiredSequence, " "), " ")

	fmt.Printf("NEWdesiredSequence: %v\n", desiredSequence)
	fmt.Printf("len(desiredSequence): %v\n", len(desiredSequence))

	if len(desiredSequence) > 1 {
		pos := 0
		for i, tag := range tagList {
			fmt.Printf("pos: %v\n", pos)
			fmt.Printf("tag: %v\n", tag)
			fmt.Printf("desiredSequence[pos]: %v\n", desiredSequence[pos])

			if tag == desiredSequence[pos] {
				pos += 1
				if tag != "O" {
					new_tag := strings.Split(inTag[i], ",")[0]
					tupleList = append(tupleList, new_tag[1:])
				}

			} else {
				pos = 0
				tupleList = []string{}
			}

			if pos == len(desiredSequence)-1 {
				new_tag := strings.Split(inTag[i+1], ",")[0]
				tupleList = append(tupleList, new_tag[1:]) //weird issue of not adding ending
				finalList = append(finalList, strings.Join(tupleList, " "))
				break
			}
		}
	}
	fmt.Printf("tupleList: %v\n", tupleList)
	fmt.Printf("finalList: %v\n", finalList)
	return finalList

	/*
		ba-zi (disposition) rule
		ANIM DIS ki PN-ta REC

		ANIM
		Disposition
		ki Person-ta ba-zi
	*/

	/*
		sz ba-ti
	*/
}

func (r *RelationExtractorRB) getFromRecRules(allTabletLines string) []string {
	// fmt.Printf("allTabletLines: %v\n", allTabletLines)
	//temporarily filter out strings like  1(disz)
	tabletLines := strings.ReplaceAll(allTabletLines, "mu-kux(DU)", "mu-ku-DU") //temporary
	r.re2.ReplaceAllString(tabletLines, "")
	inTag := r.re.FindAllString(tabletLines, 100)
	fmt.Printf("inTag: %v\n", inTag)
	// fmt.Printf("inTag: %v\n", inTag)
	tagList := []string{}
	for _, tag := range inTag { //filters out items like 1(disz)
		if strings.Contains(tag, ",") {
			new_tag := strings.Split(tag, ",")[1]
			tagList = append(tagList, new_tag[0:len(new_tag)-1])
		} else {
			tagList = append(tagList, "O")
		}
	}
	fmt.Printf("tagList: %v\n", tagList)
	/*
	   mu du (delivery) rule

	   ANIM
	   Person 1 PN
	   mu-kux(DU) DEL
	   Person 2 PN i3-dab5 REC
	*/
	finalList := []string{}
	tupleList := []string{}

	regexAbove := `ANIM [O\s?]* PN REC`
	findPat, _ := regexp.Compile(regexAbove)

	//PN REC ANIM

	// test := findPat.FindAllString(strings.Join(tagList, " "), 100)
	// fmt.Printf("test: %v\n", test)
	//Person delivered animal

	// (ANIM, 0), (ANIM, 1), (PN, 2), (O, 3), (DEL, 4)

	// println("HERE")
	desiredSequence := findPat.FindAllString(strings.Join(tagList, " "), 10)
	desiredSequence = strings.Split(strings.Join(desiredSequence, " "), " ")

	fmt.Printf("NEWdesiredSequence: %v\n", desiredSequence)
	fmt.Printf("len(desiredSequence): %v\n", len(desiredSequence))

	if len(desiredSequence) > 1 {
		pos := 0
		for i, tag := range tagList {
			fmt.Printf("pos: %v\n", pos)
			fmt.Printf("tag: %v\n", tag)
			fmt.Printf("desiredSequence[pos]: %v\n", desiredSequence[pos])

			if tag == desiredSequence[pos] {
				pos += 1
				if tag != "O" {
					new_tag := strings.Split(inTag[i], ",")[0]
					tupleList = append(tupleList, new_tag[1:])
				}

			} else {
				pos = 0
				tupleList = []string{}
			}

			if pos == len(desiredSequence)-1 {
				new_tag := strings.Split(inTag[i+1], ",")[0]
				tupleList = append(tupleList, new_tag[1:]) //weird issue of not adding ending
				finalList = append(finalList, strings.Join(tupleList, " "))
				break
			}
		}
	}
	fmt.Printf("tupleList: %v\n", tupleList)
	fmt.Printf("finalList: %v\n", finalList)
	return finalList

	/*
		ba-zi (disposition) rule
		ANIM DIS ki PN-ta REC

		ANIM
		Disposition
		ki Person-ta ba-zi
	*/

	/*
		sz ba-ti
	*/
}
