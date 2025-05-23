package service

import (
	"context"
	"coresamples/common"
	"coresamples/dbutils"
	"coresamples/ent"
	"coresamples/ent/testlist"
	"coresamples/external"
	"coresamples/util"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/shopspring/decimal"
)

// global constants
var decimalZero = decimal.NewFromInt(0)
var volumeForEachTube = decimal.NewFromInt(4000)
var volumeForEachTubeSerum = decimal.NewFromInt(2000)
var volumeForEachTubeEDTA = decimal.NewFromInt(3000)
var volumeForDead = decimal.NewFromInt(1000)
var volumeForEDTAGenotypeException = decimal.NewFromInt(2000)
var tubesWithNoVolumeRequirement = map[testlist.TubeType]bool{
	testlist.TubeTypeEDTAOrSwab: true, testlist.TubeTypeNPSwab: true, testlist.TubeTypeParasites: true,
	testlist.TubeTypeESR: true, testlist.TubeTypeUrine: true, testlist.TubeTypeTES: true, testlist.TubeTypeSTOOL: true,
	testlist.TubeTypeUNPRESERVED_STOOL: true, testlist.TubeTypeMETAL_FREE_URINE: true, testlist.TubeTypePLASMA: true,
	testlist.TubeTypeURINE_M: true, testlist.TubeTypeURINE_A: true, testlist.TubeTypeURINE_E: true, testlist.TubeTypeURINE_N: true,
	testlist.TubeTypeSALIVA_M: true, testlist.TubeTypeSALIVA_A: true, testlist.TubeTypeSALIVA_E: true, testlist.TubeTypeSALIVA_N: true,
	testlist.TubeTypeSwab: true, testlist.TubeTypeCOVID19_FINGERPRICK: true, testlist.TubeTypeCOVID19_NPSWAB: true,
	testlist.TubeTypeCOVID19_SST: true, testlist.TubeTypeSALIVA: true, testlist.TubeTypeNA: true,
}

type InternalRequiredTubeVolumeResponse struct {
	VolumeRequired                        map[string]string                `json:"volume_required,omitempty"`
	NumberOfTubes                         map[string]int32                 `json:"number_of_tubes,omitempty"`
	NumberOfDBSBloodTubes                 map[string]int32                 `json:"number_of_DBS_blood_tubes,omitempty"`
	TubeOrder                             map[string]int32                 `json:"tube_order,omitempty"`
	TubeInfo                              map[string]*ent.TubeInstructions `json:"tube_information,omitempty"`
	ActualVolumeRequired                  map[string]string                `json:"actual_volume_required,omitempty"`
	ActualNumberOfTubes                   map[string]int32                 `json:"actual_number_of_tubes,omitempty"`
	isDbsPossible                         bool
	internalNumberOfTubes                 map[testlist.TubeType]int32
	internalActualNumberOfTubes           map[testlist.TubeType]int32
	internalNumberOfDBSBloodTubes         map[testlist.TubeType]int32
	internalVolumeRequired                map[testlist.TubeType]decimal.Decimal
	internalActualVolumeRequired          map[testlist.TubeType]decimal.Decimal
	isMicroNutrientsDBSOrdered            bool
	isMicroNutrientsDBSAloneOrdered       bool
	isMicroNutrientsDBSTheOnlyTestOrdered bool
}

type GetRequiredTubeVolumeService struct {
	Service
	groupInfo    external.TestGroupInfo
	tubeTypeInfo map[testlist.TubeType]*ent.TubeInstructions
	ticker       *time.Ticker
}

type TubeVolEntry struct {
	tubeType testlist.TubeType
	volume   decimal.Decimal
}

func (r *InternalRequiredTubeVolumeResponse) init(testIds []int) {
	r.VolumeRequired = make(map[string]string)
	r.NumberOfDBSBloodTubes = make(map[string]int32)
	r.NumberOfTubes = make(map[string]int32)
	r.ActualNumberOfTubes = make(map[string]int32)
	r.ActualVolumeRequired = make(map[string]string)
	r.internalNumberOfTubes = make(map[testlist.TubeType]int32)
	r.TubeInfo = make(map[string]*ent.TubeInstructions)
	r.TubeOrder = make(map[string]int32)
	r.internalActualNumberOfTubes = make(map[testlist.TubeType]int32)
	r.internalNumberOfDBSBloodTubes = make(map[testlist.TubeType]int32)
	r.internalVolumeRequired = make(map[testlist.TubeType]decimal.Decimal)
	r.internalActualVolumeRequired = make(map[testlist.TubeType]decimal.Decimal)
	r.isDbsPossible = testIds == nil || len(testIds) == 0
	r.isMicroNutrientsDBSOrdered = false
	r.isMicroNutrientsDBSAloneOrdered = true
	r.isMicroNutrientsDBSTheOnlyTestOrdered = true
}

func (r *InternalRequiredTubeVolumeResponse) convertVolume() {
	for name, volume := range r.internalVolumeRequired {
		if _, ok := r.internalActualVolumeRequired[name]; !ok {
			r.internalActualVolumeRequired[name] = volume
		}
		r.VolumeRequired[name.String()] = volume.String()
	}
	for name, volume := range r.internalActualVolumeRequired {
		r.ActualVolumeRequired[name.String()] = volume.String()
	}
	for name, number := range r.internalNumberOfTubes {
		if _, ok := r.internalActualNumberOfTubes[name]; !ok {
			r.internalActualNumberOfTubes[name] = number
		}
		r.NumberOfTubes[name.String()] = number
	}
	for name, number := range r.internalActualNumberOfTubes {
		r.ActualNumberOfTubes[name.String()] = number
	}
	for name, number := range r.internalNumberOfDBSBloodTubes {
		r.NumberOfDBSBloodTubes[name.String()] = number
	}
}

func (r *InternalRequiredTubeVolumeResponse) isTwoEDTANeeded(tubeType testlist.TubeType) bool {
	return tubeType == testlist.TubeTypeBLOOD_FINGERPRICK && r.isMicroNutrientsDBSOrdered && !r.isMicroNutrientsDBSAloneOrdered
}

func (s *GetRequiredTubeVolumeService) Init(dbClient *ent.Client, redisClient *common.RedisClient) {
	s.Service = InitService(dbClient, redisClient)
	s.groupInfo = external.TestGroupInfo{}
	s.groupInfo.Init()
	s.tubeTypeInfo = make(map[testlist.TubeType]*ent.TubeInstructions)
	s.ticker = time.NewTicker(30 * time.Minute)
}

func (s *GetRequiredTubeVolumeService) GetRequiredTubeVolume(testIds []int32, ctx context.Context) (*InternalRequiredTubeVolumeResponse, error) {
	select {
	case <-s.ticker.C:
		s.groupInfo.ReSync(common.Env.ConsulClient, common.Env.ConsulPrefix)
	default:
		break
	}
	if len(s.groupInfo.GroupTestMappings) == 0 || len(s.groupInfo.SpecialGroupTests) == 0 || len(s.groupInfo.Groups) == 0 {
		s.groupInfo.ReSync(common.Env.ConsulClient, common.Env.ConsulPrefix)
		if len(s.groupInfo.GroupTestMappings) == 0 || len(s.groupInfo.SpecialGroupTests) == 0 || len(s.groupInfo.Groups) == 0 {
			return nil, fmt.Errorf("Unable to find group info")
		}
	}
	span, sctx := opentracing.StartSpanFromContext(ctx, "GetRequiredTubeVolume")
	defer span.Finish()
	var volumeForVGTestsOneMoreEDTATube []decimal.Decimal
	tubesForProduct := map[external.SpecialProduct]testlist.TubeType{}
	groupMap := map[string]*TubeVolEntry{}
	ids := make([]int, len(testIds))
	response := &InternalRequiredTubeVolumeResponse{}
	testMap := map[int]*ent.TestList{}
	otherTubeRequired := map[testlist.TubeType]bool{}

	// convert int32 to int
	for idx, id := range testIds {
		ids[idx] = int(id)
	}
	// get all corresponding tests
	tests, err := dbutils.GetTestsByIds(ids, s.dbClient, sctx)
	if err != nil {
		common.Errorf("error querying test_list", err)
		return nil, err
	}

	for _, test := range tests {
		testMap[test.ID] = test
	}

	response.init(ids)

	s.calcTubeVolumeException(testMap, response)
	for _, test := range testMap {
		s.collectTubeVolume(test, response, volumeForVGTestsOneMoreEDTATube, tubesForProduct, groupMap)
	}
	// calculate volume required for all groups
	s.calcVolumeForAllGroups(groupMap, response)
	for tubeType, volume := range response.internalVolumeRequired {
		if tubeType != testlist.TubeTypeBLOOD_FINGERPRICK && tubeType != testlist.TubeTypeDNAFingerprick {
			s.getVolumeForTube(response, tubeType, volume, response.internalVolumeRequired)
			if tubeType == testlist.TubeTypeFROZEN_SERUM || tubeType == testlist.TubeTypePLASMA_EDTA {
				// Max number of tubes is limited to 2
				response.internalNumberOfTubes[tubeType] = util.Min(response.internalNumberOfTubes[tubeType], 2)
			}
		} else {
			response.internalNumberOfTubes[tubeType] = int32(volume.RoundCeil(0).IntPart())
		}
	}

	for _, test := range testMap {
		tubesWithNoVolumeReq := tubesWithNoVolumeRequirement[test.TubeType] && !s.isCalculationTest(test.TestInstrument)
		if tubesWithNoVolumeReq || test.TubeType == testlist.TubeTypeBLOOD_FINGERPRICK || test.TubeType == testlist.TubeTypeDNAFingerprick {
			otherTubeRequired[test.TubeType] = true
		}
	}

	for tubeType := range otherTubeRequired {
		if tubeType == testlist.TubeTypeEDTAOrSwab || tubeType == testlist.TubeTypeBLOOD_FINGERPRICK || tubeType == testlist.TubeTypeDNAFingerprick {
			response.isDbsPossible = true
			continue
		}

		if tubeType == testlist.TubeTypeBLOOD_FINGERPRICK || tubeType == testlist.TubeTypeDNAFingerprick {
			response.internalNumberOfTubes[tubeType] = int32(response.internalVolumeRequired[tubeType].Ceil().IntPart())
			var tubeCnt int32 = 1
			if response.isTwoEDTANeeded(tubeType) {
				tubeCnt = 2
			}

			response.internalNumberOfTubes[tubeType] = tubeCnt
			continue
		}

		if tubeType == testlist.TubeTypeSTOOL {
			response.internalNumberOfTubes[testlist.TubeTypeUNPRESERVED_STOOL] = 1
		}

		if tubeType == testlist.TubeTypeUNPRESERVED_STOOL {
			response.internalNumberOfTubes[testlist.TubeTypeSTOOL] = 1
		}
		response.internalNumberOfTubes[tubeType] = 1
	}

	// At least 2 sst tubes needed
	if response.internalNumberOfTubes[testlist.TubeTypeSST] == 1 {
		response.internalNumberOfTubes[testlist.TubeTypeSST] = 2
		response.internalVolumeRequired[testlist.TubeTypeSST] = volumeForEachTubeSerum.Mul(decimal.NewFromInt(2))
	}

	// covert volume from decimal to string to preserve precision
	response.convertVolume()
	spanGather, _ := opentracing.StartSpanFromContext(sctx, "GatherTestInstructions")
	for t, _ := range response.internalNumberOfTubes {
		info, found := s.tubeTypeInfo[t]
		if !found {
			info, err = dbutils.GetTubeInfoByEnum(t, s.dbClient, ctx)
			if err != nil {
				common.Error(err)
				continue
			}
			s.tubeTypeInfo[t] = info
		}

		response.TubeInfo[t.String()] = info
		response.TubeOrder[t.String()] = int32(info.SortOrder)
	}
	spanGather.Finish()
	return response, nil
}

func (s *GetRequiredTubeVolumeService) calcVolumeForAllGroups(groupMap map[string]*TubeVolEntry, response *InternalRequiredTubeVolumeResponse) {
	for _, entry := range groupMap {
		if _, ok := response.internalVolumeRequired[entry.tubeType]; !ok {
			response.internalVolumeRequired[entry.tubeType] = entry.volume.Round(1)
			response.internalActualVolumeRequired[entry.tubeType] = entry.volume.Round(1)
		} else {
			response.internalVolumeRequired[entry.tubeType] = response.internalVolumeRequired[entry.tubeType].Add(entry.volume.Round(1))
			response.internalActualVolumeRequired[entry.tubeType] = response.internalActualVolumeRequired[entry.tubeType].Add(entry.volume.Round(1))
		}
	}
}

func (s *GetRequiredTubeVolumeService) collectTubeVolume(test *ent.TestList, resp *InternalRequiredTubeVolumeResponse,
	volumeForVGTestsOneMoreEDTATube []decimal.Decimal, tubesForProduct map[external.SpecialProduct]testlist.TubeType, groupMap map[string]*TubeVolEntry) {
	testGroup := test.DIGroupName
	volume := decimal.NewFromFloat(test.VolumeRequired)
	tubeType := test.TubeType
	testInstrument := test.TestInstrument
	if volume.Equal(decimalZero) || s.isCalculationTest(testInstrument) || s.requireNoVolume(tubeType) {
		return
	}
	if len(volumeForVGTestsOneMoreEDTATube) == 0 && s.testGroupRequireOneMoreEDTATube(testGroup) {
		volumeForVGTestsOneMoreEDTATube = append(volumeForVGTestsOneMoreEDTATube, volumeForEDTAGenotypeException)
	}

	resp.isMicroNutrientsDBSTheOnlyTestOrdered = s.isMicroNutrientsDBSTheOnlyTestOrdered(testGroup)

	if s.isGeneticsOrdered(testGroup) {
		groupMap[testGroup] = &TubeVolEntry{tubeType: tubeType, volume: volume}
		if util.StringEqualIgnoreCase(testGroup, "MicroNutrients DBS") {
			resp.isMicroNutrientsDBSOrdered = true
		} else if tubeType == testlist.TubeTypeBLOOD_FINGERPRICK {
			resp.isMicroNutrientsDBSAloneOrdered = false
		}
	} else {
		resp.isMicroNutrientsDBSAloneOrdered = tubeType == testlist.TubeTypeBLOOD_FINGERPRICK
		if !s.checkIfSpecialProductName(tubesForProduct, test, tubeType, volume, resp) {
			if _, ok := resp.internalVolumeRequired[tubeType]; !ok {
				resp.internalVolumeRequired[tubeType] = volume
			} else {
				resp.internalVolumeRequired[tubeType] = resp.internalVolumeRequired[tubeType].Add(volume)
			}
		}
	}
}

func (s *GetRequiredTubeVolumeService) checkIfSpecialProductName(tubesForProduct map[external.SpecialProduct]testlist.TubeType, test *ent.TestList,
	tubeType testlist.TubeType, volume decimal.Decimal, resp *InternalRequiredTubeVolumeResponse) bool {
	for specialProduct, testList := range s.groupInfo.SpecialGroupTests {
		if _, ok := testList[test.ID]; ok {
			if _, ok := tubesForProduct[specialProduct]; !ok {
				if _, ok := resp.internalVolumeRequired[tubeType]; !ok {
					resp.internalVolumeRequired[tubeType] = volume
				} else {
					resp.internalVolumeRequired[tubeType] = resp.internalVolumeRequired[tubeType].Add(volume)
				}
				tubesForProduct[specialProduct] = tubeType
			}
			return true
		}
	}
	return false
}

func (s *GetRequiredTubeVolumeService) calcTubeVolumeException(testMap map[int]*ent.TestList, resp *InternalRequiredTubeVolumeResponse) {
	// hormone
	hormoneTests := s.groupInfo.GetTestsUnderGroup(453)
	hasHormoneTests := true
	for _, id := range hormoneTests {
		if _, ok := testMap[id]; !ok {
			hasHormoneTests = false
			break
		}
	}
	if hasHormoneTests {
		// Heavy Metals - Urine, Environmental toxin, Organic Acids, PFAS Chemicals, Mycotoxins, Oxidative stress profile
		// when ordered with hormone zoomer, will add and only add a tube of metal free urine
		addonPackages := []int{259, 298, 307, 426, 434, 449}
		packages := []int{}
		for _, pkg := range addonPackages {
			hasPkg := true
			tests := s.groupInfo.GetTestsUnderGroup(pkg)
			for _, id := range tests {
				if _, ok := testMap[id]; !ok {
					hasPkg = false
					break
				}
			}
			if hasPkg {
				packages = append(packages, pkg)
			}
		}

		for _, pkg := range packages {
			tests := s.groupInfo.GetTestsUnderGroup(pkg)
			for _, id := range tests {
				delete(testMap, id)
			}
		}

		for _, id := range hormoneTests {
			delete(testMap, id)
		}
		resp.internalNumberOfTubes[testlist.TubeTypeURINE_A] = 1
		resp.internalNumberOfTubes[testlist.TubeTypeURINE_M] = 1
		resp.internalNumberOfTubes[testlist.TubeTypeURINE_N] = 1
		resp.internalNumberOfTubes[testlist.TubeTypeURINE_E] = 1
		if len(packages) > 0 {
			resp.internalNumberOfTubes[testlist.TubeTypeMETAL_FREE_URINE] = 1
		}
	}
	// UTI zoomer panel
	utiZoomerTests := s.groupInfo.GetTestsUnderGroup(452)
	hasZoomerTests := true

	for _, id := range utiZoomerTests {
		if _, ok := testMap[id]; !ok {
			hasZoomerTests = false
			break
		}
	}
	if hasZoomerTests {
		for _, id := range utiZoomerTests {
			delete(testMap, id)
		}
		resp.internalNumberOfTubes[testlist.TubeTypeURINE_UTI] = 1
		resp.internalVolumeRequired[testlist.TubeTypeURINE_UTI] = decimal.NewFromInt(1000)

		resp.internalNumberOfTubes[testlist.TubeTypeURINE_ANALYSIS] = 1
		resp.internalVolumeRequired[testlist.TubeTypeURINE_ANALYSIS] = decimal.NewFromInt(1000)
	}

	nutriProTests := s.groupInfo.GetTestsUnderGroup(362)
	// micronutrient tests
	microNutrientsTests := s.groupInfo.GetTestsUnderGroup(348)
	hasMicroNutrient := true
	hasNutriPro := true

	for _, id := range microNutrientsTests {
		if _, ok := testMap[id]; !ok {
			hasMicroNutrient = false
			break
		}
	}
	for _, id := range nutriProTests {
		if _, ok := testMap[id]; !ok {
			hasNutriPro = false
			break
		}
	}
	if hasMicroNutrient {
		for _, id := range microNutrientsTests {
			delete(testMap, id)
		}
		if !hasNutriPro {
			resp.internalNumberOfTubes[testlist.TubeTypeEDTA] = 3
			resp.internalNumberOfTubes[testlist.TubeTypeTES] = 1
			resp.internalNumberOfTubes[testlist.TubeTypeSST] = 2
			resp.internalVolumeRequired[testlist.TubeTypeEDTA] = decimal.NewFromInt(8000)
		}
	}
	// allergy tests
	hasAllergyTests := false
	allergyTests := s.groupInfo.GetTestsUnderGroups(441, 203)

	for id := range allergyTests {
		if _, ok := testMap[id]; ok {
			hasAllergyTests = true
			break
		}
	}

	if hasAllergyTests {
		for id := range allergyTests {
			delete(testMap, id)
		}
		resp.internalNumberOfTubes[testlist.TubeTypeSST] = 1
		resp.internalVolumeRequired[testlist.TubeTypeSST] = decimal.NewFromInt(1000)
	}

	// fatty acids tests
	hasFattyAcidsTests := true
	FattyAcidsTests := s.groupInfo.GetTestsUnderGroup(158)

	for _, id := range FattyAcidsTests {
		if _, ok := testMap[id]; !ok {
			hasFattyAcidsTests = false
			break
		}
	}

	if hasFattyAcidsTests {
		for _, id := range FattyAcidsTests {
			delete(testMap, id)
		}
		resp.internalNumberOfTubes[testlist.TubeTypeEDTA] = 1
	}
}

func (s *GetRequiredTubeVolumeService) getVolumeForTube(resp *InternalRequiredTubeVolumeResponse, tubeType testlist.TubeType, volume decimal.Decimal, volumeRequired map[testlist.TubeType]decimal.Decimal) {
	eachVolume := volumeForEachTube
	if tubeType == testlist.TubeTypeEDTA {
		eachVolume = volumeForEachTubeEDTA
	} else if tubeType == testlist.TubeTypeSST {
		eachVolume = volumeForEachTubeSerum
	}
	volumeWithDeadVol := volume.Add(volumeForDead)
	if volumeRequired != nil {
		volumeRequired[tubeType] = volumeWithDeadVol
	}
	resp.internalNumberOfTubes[tubeType] = int32(volumeWithDeadVol.Div(eachVolume).RoundUp(0).IntPart())
	resp.internalActualNumberOfTubes[tubeType] = resp.internalNumberOfTubes[tubeType]
}

func (s *GetRequiredTubeVolumeService) requireNoVolume(tubeType testlist.TubeType) bool {
	for t := range tubesWithNoVolumeRequirement {
		if tubeType == t {
			return true
		}
	}
	return false
}

func (s *GetRequiredTubeVolumeService) isCalculationTest(instrument testlist.TestInstrument) bool {
	return instrument == testlist.TestInstrumentRocheCalculation
}

func (s *GetRequiredTubeVolumeService) testGroupRequireOneMoreEDTATube(testGroup string) bool {
	return testGroup == "CELIAC_GEN" || testGroup == "APOE" || testGroup == "FACTOR" ||
		testGroup == "MTHFR" || testGroup == "Nutriprofile" || testGroup == "Methylation Panel" || testGroup == "CARDIAX"
}

func (s *GetRequiredTubeVolumeService) isMicroNutrientsDBSTheOnlyTestOrdered(testGroup string) bool {
	return testGroup == "WHOLE_BLOOD_FATTY_ACIDS" || util.StringEqualIgnoreCase(testGroup, "MicroNutrients DBS")
}

func (s *GetRequiredTubeVolumeService) isGeneticsOrdered(testGroup string) bool {
	return len(testGroup) != 0 && !util.StringEqualIgnoreCase(testGroup, "NULL") && !util.StringEqualIgnoreCase(testGroup, "NH") &&
		!util.StringEqualIgnoreCase(testGroup, "CIRS") && !util.StringEqualIgnoreCase(testGroup, "CYB")
}
