package model

var testInstrumentsNames = [...]string{
	"WELLNESS", "Roche Calculation", "Diazyme", "Outsource", "NA", "Phadia_250", "Roche", "Sed rate", "Sysmex", "TSP", "VG", "Xevo TQS", "Beckman", "Siemens", "Xevo TQ-XS", "Perkin ICP-MS", "Xevo TQD", "Urinalysis Clintek Novus", "Urinalysis UF",
}

func GetTestInstrumentNamedValues() []string {
	testInstrumentsNamedValues := make([]string, len(testInstrumentsNames)*2)
	for idx, name := range testInstrumentsNames {
		testInstrumentsNamedValues[idx*2] = name
		if name == "NA" {
			testInstrumentsNamedValues[idx*2+1] = "N/A"
		} else if name == "Roche Calculation" {
			testInstrumentsNamedValues[idx*2+1] = "Calculation"
		} else if name == "Outsource" {
			testInstrumentsNamedValues[idx*2+1] = "LC"
		} else {
			testInstrumentsNamedValues[idx*2+1] = name
		}
	}
	return testInstrumentsNamedValues
}
