CLDI_Extractors is a pipeline composed of the following 


1. ATF Parser - Parses individual CLDI Data Tablets and stores information in CLDIData object which contains tablet_num, publication, tablet_location (@observer, @reverse, etc.) and line_no corresponding to transliterations 


2. ATF Normalizer - Normalizes transliterations to a common format. This includes standardizing graphemes, sign_lists and dates. Normalized transliterations are stored under 'normalized_translit' in the CLDI Objects.

3. Entity Extractor - Extracts entity from transliterations

4. Data Writer - Writes out CLDIData object to specified output path