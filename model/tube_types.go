package model

var tubeTypeNames = [...]string{"SST", "Urine", "NA", "PLASMA", "TES", "ESR", "EDTA", "STOOL", "Swab", "NP Swab", "Parasites",
	"EDTA or Swab", "METAL_FREE_URINE", "SODIUM_CITRATE_PLASMA", "URINE_M", "URINE_A", "URINE_E",
	"URINE_N", "SALIVA_M", "SALIVA_A", "SALIVA_E", "SALIVA_N", "UNPRESERVED_STOOL",
	"SALIVA", "COVID19_NPSWAB", "COVID19_FINGERPRICK", "COVID19_EDTA", "BLOOD_FINGERPRICK", "DNA fingerprick",
	"COVID19_SST", "PLASMA_EDTA", "PLASMA_EDTA_PLATELET_FREE", "FROZEN_SERUM", "PLASMA_CITRATE_PLATELET_POOR", "PLASMA_EDTA_TRASYLOL",
	"BLOOD_MICROTUBE", "DNA_FINGERPRICK", "URINE_UTI", "URINE_ANALYSIS"}

func GetTubeTypeNamedValues() []string {
	tubeTypeNamedValues := make([]string, len(tubeTypeNames)*2)
	for idx, name := range tubeTypeNames {
		tubeTypeNamedValues[idx*2] = name
		if name == "NA" {
			tubeTypeNamedValues[idx*2+1] = "N/A"
		} else {
			tubeTypeNamedValues[idx*2+1] = name
		}
	}
	return tubeTypeNamedValues
}
