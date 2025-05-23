// Code generated by ent, DO NOT EDIT.

package testlist

import (
	"fmt"

	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the testlist type in the database.
	Label = "test_list"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "test_id"
	// FieldTestInstrument holds the string denoting the test_instrument field in the database.
	FieldTestInstrument = "test_instrument"
	// FieldTubeType holds the string denoting the tube_type field in the database.
	FieldTubeType = "tube_type"
	// FieldDIGroupName holds the string denoting the di_group_name field in the database.
	FieldDIGroupName = "DI_group_name"
	// FieldVolumeRequired holds the string denoting the volume_required field in the database.
	FieldVolumeRequired = "volume_required"
	// FieldBloodType holds the string denoting the blood_type field in the database.
	FieldBloodType = "blood_type"
	// Table holds the table name of the testlist in the database.
	Table = "test_list"
)

// Columns holds all SQL columns for testlist fields.
var Columns = []string{
	FieldID,
	FieldTestInstrument,
	FieldTubeType,
	FieldDIGroupName,
	FieldVolumeRequired,
	FieldBloodType,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DIGroupNameValidator is a validator for the "DI_group_name" field. It is called by the builders before save.
	DIGroupNameValidator func(string) error
)

// TestInstrument defines the type for the "test_instrument" enum field.
type TestInstrument string

// TestInstrumentNA is the default value of the TestInstrument enum.
const DefaultTestInstrument = TestInstrumentNA

// TestInstrument values.
const (
	TestInstrumentWELLNESS               TestInstrument = "WELLNESS"
	TestInstrumentRocheCalculation       TestInstrument = "Calculation"
	TestInstrumentDiazyme                TestInstrument = "Diazyme"
	TestInstrumentOutsource              TestInstrument = "LC"
	TestInstrumentNA                     TestInstrument = "N/A"
	TestInstrumentPhadia_250             TestInstrument = "Phadia_250"
	TestInstrumentRoche                  TestInstrument = "Roche"
	TestInstrumentSedRate                TestInstrument = "Sed rate"
	TestInstrumentSysmex                 TestInstrument = "Sysmex"
	TestInstrumentTSP                    TestInstrument = "TSP"
	TestInstrumentVG                     TestInstrument = "VG"
	TestInstrumentXevoTQS                TestInstrument = "Xevo TQS"
	TestInstrumentBeckman                TestInstrument = "Beckman"
	TestInstrumentSiemens                TestInstrument = "Siemens"
	TestInstrumentXevoTQXS               TestInstrument = "Xevo TQ-XS"
	TestInstrumentPerkinICPMS            TestInstrument = "Perkin ICP-MS"
	TestInstrumentXevoTQD                TestInstrument = "Xevo TQD"
	TestInstrumentUrinalysisClintekNovus TestInstrument = "Urinalysis Clintek Novus"
	TestInstrumentUrinalysisUF           TestInstrument = "Urinalysis UF"
)

func (ti TestInstrument) String() string {
	return string(ti)
}

// TestInstrumentValidator is a validator for the "test_instrument" field enum values. It is called by the builders before save.
func TestInstrumentValidator(ti TestInstrument) error {
	switch ti {
	case TestInstrumentWELLNESS, TestInstrumentRocheCalculation, TestInstrumentDiazyme, TestInstrumentOutsource, TestInstrumentNA, TestInstrumentPhadia_250, TestInstrumentRoche, TestInstrumentSedRate, TestInstrumentSysmex, TestInstrumentTSP, TestInstrumentVG, TestInstrumentXevoTQS, TestInstrumentBeckman, TestInstrumentSiemens, TestInstrumentXevoTQXS, TestInstrumentPerkinICPMS, TestInstrumentXevoTQD, TestInstrumentUrinalysisClintekNovus, TestInstrumentUrinalysisUF:
		return nil
	default:
		return fmt.Errorf("testlist: invalid enum value for test_instrument field: %q", ti)
	}
}

// TubeType defines the type for the "tube_type" enum field.
type TubeType string

// TubeType values.
const (
	TubeTypeSST                          TubeType = "SST"
	TubeTypeUrine                        TubeType = "Urine"
	TubeTypeNA                           TubeType = "N/A"
	TubeTypePLASMA                       TubeType = "PLASMA"
	TubeTypeTES                          TubeType = "TES"
	TubeTypeESR                          TubeType = "ESR"
	TubeTypeEDTA                         TubeType = "EDTA"
	TubeTypeSTOOL                        TubeType = "STOOL"
	TubeTypeSwab                         TubeType = "Swab"
	TubeTypeNPSwab                       TubeType = "NP Swab"
	TubeTypeParasites                    TubeType = "Parasites"
	TubeTypeEDTAOrSwab                   TubeType = "EDTA or Swab"
	TubeTypeMETAL_FREE_URINE             TubeType = "METAL_FREE_URINE"
	TubeTypeSODIUM_CITRATE_PLASMA        TubeType = "SODIUM_CITRATE_PLASMA"
	TubeTypeURINE_M                      TubeType = "URINE_M"
	TubeTypeURINE_A                      TubeType = "URINE_A"
	TubeTypeURINE_E                      TubeType = "URINE_E"
	TubeTypeURINE_N                      TubeType = "URINE_N"
	TubeTypeSALIVA_M                     TubeType = "SALIVA_M"
	TubeTypeSALIVA_A                     TubeType = "SALIVA_A"
	TubeTypeSALIVA_E                     TubeType = "SALIVA_E"
	TubeTypeSALIVA_N                     TubeType = "SALIVA_N"
	TubeTypeUNPRESERVED_STOOL            TubeType = "UNPRESERVED_STOOL"
	TubeTypeSALIVA                       TubeType = "SALIVA"
	TubeTypeCOVID19_NPSWAB               TubeType = "COVID19_NPSWAB"
	TubeTypeCOVID19_FINGERPRICK          TubeType = "COVID19_FINGERPRICK"
	TubeTypeCOVID19_EDTA                 TubeType = "COVID19_EDTA"
	TubeTypeBLOOD_FINGERPRICK            TubeType = "BLOOD_FINGERPRICK"
	TubeTypeDNAFingerprick               TubeType = "DNA fingerprick"
	TubeTypeCOVID19_SST                  TubeType = "COVID19_SST"
	TubeTypePLASMA_EDTA                  TubeType = "PLASMA_EDTA"
	TubeTypePLASMA_EDTA_PLATELET_FREE    TubeType = "PLASMA_EDTA_PLATELET_FREE"
	TubeTypeFROZEN_SERUM                 TubeType = "FROZEN_SERUM"
	TubeTypePLASMA_CITRATE_PLATELET_POOR TubeType = "PLASMA_CITRATE_PLATELET_POOR"
	TubeTypePLASMA_EDTA_TRASYLOL         TubeType = "PLASMA_EDTA_TRASYLOL"
	TubeTypeBLOOD_MICROTUBE              TubeType = "BLOOD_MICROTUBE"
	TubeTypeDNA_FINGERPRICK              TubeType = "DNA_FINGERPRICK"
	TubeTypeURINE_UTI                    TubeType = "URINE_UTI"
	TubeTypeURINE_ANALYSIS               TubeType = "URINE_ANALYSIS"
)

func (tt TubeType) String() string {
	return string(tt)
}

// TubeTypeValidator is a validator for the "tube_type" field enum values. It is called by the builders before save.
func TubeTypeValidator(tt TubeType) error {
	switch tt {
	case TubeTypeSST, TubeTypeUrine, TubeTypeNA, TubeTypePLASMA, TubeTypeTES, TubeTypeESR, TubeTypeEDTA, TubeTypeSTOOL, TubeTypeSwab, TubeTypeNPSwab, TubeTypeParasites, TubeTypeEDTAOrSwab, TubeTypeMETAL_FREE_URINE, TubeTypeSODIUM_CITRATE_PLASMA, TubeTypeURINE_M, TubeTypeURINE_A, TubeTypeURINE_E, TubeTypeURINE_N, TubeTypeSALIVA_M, TubeTypeSALIVA_A, TubeTypeSALIVA_E, TubeTypeSALIVA_N, TubeTypeUNPRESERVED_STOOL, TubeTypeSALIVA, TubeTypeCOVID19_NPSWAB, TubeTypeCOVID19_FINGERPRICK, TubeTypeCOVID19_EDTA, TubeTypeBLOOD_FINGERPRICK, TubeTypeDNAFingerprick, TubeTypeCOVID19_SST, TubeTypePLASMA_EDTA, TubeTypePLASMA_EDTA_PLATELET_FREE, TubeTypeFROZEN_SERUM, TubeTypePLASMA_CITRATE_PLATELET_POOR, TubeTypePLASMA_EDTA_TRASYLOL, TubeTypeBLOOD_MICROTUBE, TubeTypeDNA_FINGERPRICK, TubeTypeURINE_UTI, TubeTypeURINE_ANALYSIS:
		return nil
	default:
		return fmt.Errorf("testlist: invalid enum value for tube_type field: %q", tt)
	}
}

// OrderOption defines the ordering options for the TestList queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByTestInstrument orders the results by the test_instrument field.
func ByTestInstrument(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTestInstrument, opts...).ToFunc()
}

// ByTubeType orders the results by the tube_type field.
func ByTubeType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTubeType, opts...).ToFunc()
}

// ByDIGroupName orders the results by the DI_group_name field.
func ByDIGroupName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDIGroupName, opts...).ToFunc()
}

// ByVolumeRequired orders the results by the volume_required field.
func ByVolumeRequired(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldVolumeRequired, opts...).ToFunc()
}

// ByBloodType orders the results by the blood_type field.
func ByBloodType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBloodType, opts...).ToFunc()
}
