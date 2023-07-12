package CDLI_Extractor

import (
	"regexp"
	"strings"
	"sync"
)

// TransliterationCleaner is a type that provides methods for cleaning the transliteration data within CDLIData types.
// It reads from an input channel of CDLIData, performs cleaning operations, and sends the cleaned data to an output channel.
//
//the bool is to replace characters
// cleaning always happens to conform to our pipeline
// clean = conform to our pipeline, replace = alter the meaning or representation or remove entirely
type TransliterationCleaner struct {
	// path        string
	in          <-chan CDLIData
	dropTablets bool
	Out         chan CDLIData
	done        chan struct{}
}

// NewTransliterationCleaner creates a new instance of a TransliterationCleaner. The cleaner immediately starts reading from
// its input channel and cleaning data in a separate goroutine. The cleaned data is sent to its output channel.
func NewTransliterationCleaner(dropTablets bool, in <-chan CDLIData) *TransliterationCleaner {
	transliterationCleaner := &TransliterationCleaner{
		in:          in,
		Out:         make(chan CDLIData, 10000000),
		done:        make(chan struct{}, 1),
		dropTablets: dropTablets,
	}
	transliterationCleaner.run()
	return transliterationCleaner
}

// run() manages the main cleaning loop for the cleaner. It reads tablets from the input channel, cleans their transliterations,
// and sends the cleaned tablets to the output channel. It finishes when all input data has been read and cleaned.
func (c *TransliterationCleaner) run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		println("cleaning tablets transliteration")
		defer wg.Done()
		defer close(c.Out)
		for cdliTablet := range c.in { //we should analyze tablet section by section + line by line
			for i, tabletSection := range cdliTablet.TabletSections {
				cdliTablet.TabletSections[i].isBroken = false //temporary
				cdliTablet.TabletSections[i] = c.cleanTablets(tabletSection)
				if c.dropTablets && cdliTablet.TabletSections[i].isBroken {
					continue //if drop tablet line is onn and section is borken then skip tablet
				}
			}
			//else pass downstream
			c.Out <- cdliTablet
		}
		println("DONE")
	}()
	wg.Wait()
}

// WaitUntilDone allows external callers to wait until the cleaner has finished processing all of its input data.
func (c *TransliterationCleaner) WaitUntilDone() {
	c.done <- struct{}{}
}

// cleanTablets performs all cleaning operations on a given TabletSection. The cleaned TabletSection is returned.
func (c *TransliterationCleaner) cleanTablets(tabletSection TabletSection) TabletSection {
	c.cleanNumericSymbols(false, tabletSection)
	c.cleanDamagedMetaCharacters(false, tabletSection)
	c.replaceAnnotatorsCorrections(false, tabletSection)
	// c.replaceCommasInNames(false, tabletSection)
	return tabletSection

}

// cleanNumericSymbols() cleans numeric symbols from the transliteration lines of a given TabletSection. The cleaned TabletSection is returned.
//
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

// cleanDamagedMetaCharacters() cleans damaged meta characters from the transliteration lines of a given TabletSection. The cleaned TabletSection is returned.
//
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
	}
	return &tabletSection
}

// replaceAnnotatorsCorrections() removes annotators corrections from the transliteration lines of a given TabletSection. The cleaned TabletSection is returned.
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

// replaceCommasInNames replaces commas in names within the transliteration lines of a given TabletSection. The cleaned TabletSection is returned.
// func (c *TransliterationCleaner) replaceCommasInNames(flag bool, tabletSection TabletSection) *TabletSection {
//  ...
// }
//This is a specific case(s,e-lu-usz-da-gan), but due to using .csv files there is a name
//which contains a comma which should be replaced with an underscore(_)
// func (c *TransliterationCleaner) replaceCommasInNames(flag bool, tabletSection TabletSection) *TabletSection {
// 	return &tabletSection
// }
