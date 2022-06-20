package CDLI_Extractor

import "sync"

type RelationExtractorRB struct {
	in   <-chan CDLIData
	out  chan CDLIData
	done chan struct{}
}

func newRelationExtractorRB(in <-chan CDLIData) *RelationExtractorRB {
	relationExtractor := &RelationExtractorRB{
		in:   in,
		out:  make(chan CDLIData, 1000000),
		done: make(chan struct{}, 1),
	}
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
		println("rule-based relation extraction")
		defer wg.Done()
		for cdliData := range r.in {
			for i, tablet := range cdliData.TabletSections {
				tablet.RelationLines = make(map[int]string)
				cdliData.TabletSections[i].RelationLines = r.getFromRules(tablet)
			}
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

*/
func (r *RelationExtractorRB) getFromRules(tabletLines TabletSection) map[int]string {

	/*
	   mu du (delivery) rule

	   ANIM
	   Person 1
	   mu-kux(DU)
	   Person 2 i3-dab5
	*/

	/*
	   ba-zi (disposition) rule

	   ANIM
	   Disposition
	   ki Person-ta ba-zi
	*/

	/*
	   t3-dab5 (receieved rule) rule

	   ANIM
	   ki Person1-ta
	   Person2 i3-dab5
	*/

	/*
		sz ba-ti
	*/

	return tabletLines.RelationLines
}
