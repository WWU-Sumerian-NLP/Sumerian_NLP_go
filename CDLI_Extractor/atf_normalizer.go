package CDLI_Extractor

import (
	"regexp"
	"strings"
	"sync"
)

const comment_regex = `((|-|)\(\$.+\$\)(|-|))`

// ATFNormalizer is a type that provides methods for normalizing the transliteration data within CDLIData types.
// It reads from an input channel of CDLIData, performs normalization operations, and sends the normalized data to an output channel.
type ATFNormalizer struct {
	rawTransliteration []string
	normalizeDefective bool
	in                 <-chan CDLIData
	out                chan CDLIData
	done               chan struct{}
}

// newATFNormalizer creates a new instance of an ATFNormalizer. The normalizer immediately starts reading from
// its input channel and normalizing data in a separate goroutine. The normalized data is sent to its output channel.
func newATFNormalizer(normalizeDefective bool, in <-chan CDLIData) *ATFNormalizer {
	atfNormalizer := &ATFNormalizer{
		normalizeDefective: false,
		in:                 in,
		out:                make(chan CDLIData, 10000000),
		done:               make(chan struct{}, 1),
	}
	atfNormalizer.run()
	return atfNormalizer
}

// run manages the main normalization loop for the normalizer. It reads tablets from the input channel, normalizes their transliterations,
// and sends the normalized tablets to the output channel. It finishes when all input data has been read and normalized.
func (n *ATFNormalizer) run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		println("normalizing")
		for CDLIData := range n.in {
			for i, TabletSection := range CDLIData.TabletSections {
				TabletSection.NormalizedLines = make([]string, len(TabletSection.LineNumbers))
				n.rawTransliteration = CDLIData.TabletSections[i].TabletLines
				CDLIData.TabletSections[i].NormalizedLines = n.parseRawTransliteration(TabletSection)
			}
			n.out <- CDLIData
		}
		println("done")
		close(n.out)
	}()
	wg.Wait()

}

// WaitUntilDone allows external callers to wait until the normalizer has finished processing all of its input data.
func (n *ATFNormalizer) WaitUntilDone() {
	n.done <- struct{}{}
}

// parseRawTransliteration performs normalization operations on the transliteration lines of a given TabletSection.
// The normalized transliteration lines are returned.
func (n *ATFNormalizer) parseRawTransliteration(TabletSection TabletSection) []string {
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

// parseTransliterationGraphemes performs normalization operations on a given grapheme within a transliteration line.
// The normalized grapheme is returned.
func (n *ATFNormalizer) parseTransliterationGraphemes(grapheme string) string {
	grapheme = strings.ReplaceAll(grapheme, " ", "")
	// grapheme = n.replaceBracesNSlashes()

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

// replaceBracesNSlashes replaces braces and slashes within a grapheme with a suitable replacement.
// It returns the grapheme with replaced characters.
//broken
func (n *ATFNormalizer) replaceBracesNSlashes() string {
	bracesList := map[string]struct{}{
		"(": {}, ")": {},
		"{": {}, "}": {},
	}
	slashList := map[string]struct{}{
		"\\": {}, "/": {},
	}

	newGrapheme := ""
	for _, test := range strings.Split(strings.Join(n.rawTransliteration, " "), " ") {
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
	return string(newGrapheme)

}

// getUnicodeIndex returns the Unicode index for a given signDictionary. The returned string is the Unicode index.
func (n *ATFNormalizer) getUnicodeIndex(signDictionary []string) string {
	return ""
}
