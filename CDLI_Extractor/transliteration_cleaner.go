package CDLI_Extractor

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

//the bool is to replace characters
// cleaning always happens to conform to our pipeline
// clean = conform to our pipeline, replace = alter the meaning or representation or remove entirely
type TransliterationCleaner struct {
	// path        string
	in          <-chan CDLIData
	dropTablets bool
	out         chan CDLIData
	done        chan struct{}
}

func newTransliterationCleaner(dropTablets bool, in <-chan CDLIData) *TransliterationCleaner {
	transliterationCleaner := &TransliterationCleaner{
		in:          in,
		out:         make(chan CDLIData, 100000),
		done:        make(chan struct{}, 1),
		dropTablets: dropTablets,
	}
	transliterationCleaner.run()
	return transliterationCleaner
}

func (c *TransliterationCleaner) run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		println("cleaning tablets transliteration")
		defer wg.Done()
		defer close(c.out)
		for cdliTablet := range c.in { //we should analyze tablet section by section + line by line
			for i, tabletSection := range cdliTablet.TabletSections {
				cdliTablet.TabletSections[i].isBroken = false //temporary
				cdliTablet.TabletSections[i] = c.cleanTablets(tabletSection)
				if c.dropTablets && cdliTablet.TabletSections[i].isBroken {
					continue //if drop tablet line is onn and section is borken then skip tablet
				}
			}
			//else pass downstream
			c.out <- cdliTablet
		}
		println("DONE")
	}()
	wg.Wait()
}

func (c *TransliterationCleaner) WaitUntilDone() {
	c.done <- struct{}{}
}

func (c *TransliterationCleaner) cleanTablets(tabletSection TabletSection) TabletSection {
	c.cleanNumericSymbols(false, tabletSection)
	c.cleanDamagedMetaCharacters(false, tabletSection)
	c.replaceAnnotatorsCorrections(false, tabletSection)
	// c.replaceCommasInNames(false, tabletSection)
	return tabletSection

}

//1(disz) should map to 1
//check if before a parthesis there is a number
func (c *TransliterationCleaner) cleanNumericSymbols(flag bool, tabletSection TabletSection) *TabletSection {
	const regexReplaceNumeric = `-?(?:\d)\([^)]*\)|\[[^\]]*\]g`
	re, _ := regexp.Compile(regexReplaceNumeric)

	for i, line := range tabletSection.TabletLines {
		re.ReplaceAllString(line, "")
		tabletSection.TabletLines[i] = line
	}
	return &tabletSection
}

//regexCleanArrows: <ki> --> ki
//regexCleanBrackets: [...] --> ... Note: Each dot represents a potential word. This is intentional
//regexReplaceBrackets: [...] --> ""
func (c *TransliterationCleaner) cleanDamagedMetaCharacters(flag bool, tabletSection TabletSection) *TabletSection {
	const regexCleanArrows = `([<]|[>])+`
	re, _ := regexp.Compile(regexCleanArrows)

	const regexCleanBrackets = `(\[|\])`
	re2, _ := regexp.Compile(regexCleanBrackets)

	for i, line := range tabletSection.TabletLines {
		line = re.ReplaceAllString(line, "")
		line = re2.ReplaceAllString(line, "")
		tabletSection.TabletLines[i] = line
		fmt.Printf("line: %v\n", line)
	}
	return &tabletSection
}

func (c *TransliterationCleaner) replaceAnnotatorsCorrections(flag bool, tabletSection TabletSection) *TabletSection {
	for i, line := range tabletSection.TabletLines {
		//remove #
		line = strings.ReplaceAll(line, "#", "")

		//remove ?
		line = strings.ReplaceAll(line, "?", "")

		//remove !
		line = strings.ReplaceAll(line, "!", "")
		tabletSection.TabletLines[i] = line
	}
	return &tabletSection
}

//This is a specific case(s,e-lu-usz-da-gan), but due to using .csv files there is a name
//which contains a comma which should be replaced with an underscore(_)
// func (c *TransliterationCleaner) replaceCommasInNames(flag bool, tabletSection TabletSection) *TabletSection {
// 	return &tabletSection
// }
