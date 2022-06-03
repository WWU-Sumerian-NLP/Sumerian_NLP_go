package CLDI_Extractor

import (
	"regexp"
	"strings"
)

func standardizeGrapheme(grapheme string) string {
	standardizeMap := map[string]string{
		"š": "c", "ŋ": "j", "₀": "0", "₁": "1", "₂": "2",
		"₃": "3", "₄": "4", "₅": "5", "₆": "6", "₇": "7",
		"₈": "8", "₉": "9", "+": "-", "Š": "C", "Ŋ": "J",
		"·": "''", "°": "''", "sz": "c", "SZ": "C",
		"Sz": "C", "ʾ": "'", "’": "'",
	}
	for key := range standardizeMap {
		if strings.Contains(grapheme, key) {
			grapheme = strings.ReplaceAll(grapheme, key, standardizeMap[key])
		}

	}
	times, _ := regexp.Compile(`(?P<a>[\w])x(?P<b>[\w])`)
	if times.MatchString(grapheme) {
		grapheme = times.ReplaceAllString(`\g<a>×\g<b>`, grapheme)
	}
	return grapheme
}

func parseSigns(grapheme string) []string {
	sign_list := []string{}
	grapheme = restyleDeterminates(grapheme)

	return sign_list

}

func restyleDeterminates(grapheme string) string {
	re_det, _ := regexp.Compile(`((?P<o_brc>\{)(?P<det>.*?)(?P<c_brc>\}))`)
	all_matches := re_det.FindAllString(grapheme, 5)
	print(all_matches)
	return grapheme
}
