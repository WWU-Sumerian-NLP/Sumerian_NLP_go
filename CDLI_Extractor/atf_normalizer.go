package CDLI_Extractor

import (
	"regexp"
	"strings"
	"sync"
)

const comment_regex = `((|-|)\(\$.+\$\)(|-|))`

type ATFNormalizer struct {
	rawTransliteration string
	normalizeDefective bool
	in                 <-chan CDLIData
	out                chan CDLIData
	done               chan struct{}
}

func newATFNormalizer(normalizeDefective bool, in <-chan CDLIData) *ATFNormalizer {
	atfNormalizer := &ATFNormalizer{
		normalizeDefective: false,
		in:                 in,
		out:                make(chan CDLIData, 100000),
		done:               make(chan struct{}, 1),
	}
	atfNormalizer.run()
	return atfNormalizer
}

func (n *ATFNormalizer) run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		println("normalizing")
		for CDLIData := range n.in {
			for i, TabletSection := range CDLIData.TabletSections {
				TabletSection.NormalizedLines = make(map[int]string)
				CDLIData.TabletSections[i].NormalizedLines = n.parseRawTransliteration(TabletSection)
			}
			n.out <- CDLIData
		}
		println("done")
		close(n.out)
	}()
	wg.Wait()

}

func (n *ATFNormalizer) WaitUntilDone() {
	n.done <- struct{}{}
}

func (n *ATFNormalizer) parseRawTransliteration(TabletSection TabletSection) map[int]string {
	for line_no, translit := range TabletSection.TabletLines {
		//remove comments - Apply regex expression
		if strings.Contains(translit, "($") {
			r, _ := regexp.Compile(comment_regex)
			r.ReplaceAllString(translit, "")
		}
		allGraphemes := ""
		for _, grapheme := range strings.Split(translit, " ") {
			grapheme = n.parseTransliterationGraphemes(grapheme)
			allGraphemes += " " + grapheme
		}
		TabletSection.NormalizedLines[line_no] = strings.TrimSpace(allGraphemes)
	}
	return TabletSection.NormalizedLines
}

func (n *ATFNormalizer) parseTransliterationGraphemes(grapheme string) string {
	grapheme = strings.ReplaceAll(grapheme, " ", "")
	grapheme = n.replaceBracesNSlashes()
	grapheme = standardizeGrapheme(grapheme)
	// sign_list := parseSigns(grapheme)
	// // fmt.Printf("sign_list: %v\n", sign_list)
	// for i := 0; i < len(sign_list); i++ {
	// 	if strings.Contains(sign_list[i], "det") {
	// 		sign_list[i] = n.getUnicodeIndex(sign_list)
	// 	}
	// }

	return grapheme
}

func (n *ATFNormalizer) replaceBracesNSlashes() string {
	bracesList := map[string]struct{}{
		"(": {}, ")": {},
		"{": {}, "}": {},
	}
	slashList := map[string]struct{}{
		"\\": {}, "/": {},
	}

	newGrapheme := ""
	for _, test := range n.rawTransliteration {
		if newGrapheme != "" {
			if _, ok := bracesList[string(test)]; ok {
				if _, ok := slashList[string(newGrapheme[len(newGrapheme)-1])]; ok {
					newGrapheme += newGrapheme[:len(newGrapheme)-1] + string(test)
				}
			}
		} else {
			newGrapheme += string(test)
		}
	}
	return newGrapheme

}

func (n *ATFNormalizer) getUnicodeIndex(signDictionary []string) string {
	return ""
}
