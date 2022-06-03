package CLDI_Extractor

import (
	"regexp"
	"strings"
)

const comment_regex = `((|-|)\(\$.+\$\)(|-|))`

type ATFNormalizer struct {
	rawTransliteration string
	normalizeDefective bool
	in                 <-chan CLDIData
	out                chan CLDIData
	done               chan struct{}
}

func newATFNormalizer(normalizeDefective bool, in <-chan CLDIData) *ATFNormalizer {
	atfNormalizer := &ATFNormalizer{
		normalizeDefective: false,
		in:                 in,
		out:                make(chan CLDIData, 1000),
		done:               make(chan struct{}, 1),
	}
	atfNormalizer.run()
	return atfNormalizer
}

func (n *ATFNormalizer) run() {
	go func() {
		for CLDIData := range n.in {
			for i, tabletLine := range CLDIData.TabletList {
				tabletLine.NormalizedLines = make(map[string]string)
				CLDIData.TabletList[i].NormalizedLines = n.parseRawTransliteration(tabletLine)
			}
			n.out <- CLDIData
		}
		close(n.out)
		n.done <- struct{}{}
	}()

}

func (n *ATFNormalizer) WaitUntilDone() {
	n.done <- struct{}{}
}

func (n *ATFNormalizer) parseRawTransliteration(tabletLine TabletLine) map[string]string {
	for line_no, translit := range tabletLine.TabletLines {
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
		tabletLine.NormalizedLines[line_no] = strings.TrimSpace(allGraphemes)
	}
	return tabletLine.NormalizedLines
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
