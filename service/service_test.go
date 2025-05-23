package service

import (
	"coresamples/ent"
	"coresamples/ent/testlist"
	"coresamples/external"
	"github.com/shopspring/decimal"
	"testing"
)

func setupGetRequiredTubeVolumeService() (*GetRequiredTubeVolumeService, *InternalRequiredTubeVolumeResponse) {
	s := &GetRequiredTubeVolumeService{}
	s.dbClient = nil
	s.redisClient = nil
	s.groupInfo = external.TestGroupInfo{}

	resp := &InternalRequiredTubeVolumeResponse{}
	resp.init([]int{})
	return s, resp
}

func TestCollectTubeVolume(t *testing.T) {

	tests := []*ent.TestList{
		{
			ID:             4885,
			TubeType:       testlist.TubeTypeBLOOD_FINGERPRICK,
			VolumeRequired: 1,
			TestInstrument: testlist.TestInstrumentTSP,
		},
	}
	s, resp := setupGetRequiredTubeVolumeService()
	var volumeForVGTestsOneMoreEDTATube []decimal.Decimal
	tubesForProduct := map[external.SpecialProduct]testlist.TubeType{}
	groupMap := map[string]*TubeVolEntry{}
	for _, test := range tests {
		s.collectTubeVolume(test, resp, volumeForVGTestsOneMoreEDTATube, tubesForProduct, groupMap)
	}
	if !resp.internalVolumeRequired[testlist.TubeTypeBLOOD_FINGERPRICK].Equal(decimal.NewFromInt(1)) {
		t.Fatalf("volume required should be 1, get %s", resp.internalVolumeRequired[testlist.TubeTypeBLOOD_FINGERPRICK].String())
	}
}

func TestCalcVolumeForGroups(t *testing.T) {
	groupMap := map[string]*TubeVolEntry{
		"a": {volume: decimal.NewFromFloat(0.25), tubeType: testlist.TubeTypeSST},
		"b": {volume: decimal.NewFromFloat(0.25), tubeType: testlist.TubeTypeSST},
	}
	s, resp := setupGetRequiredTubeVolumeService()
	s.calcVolumeForAllGroups(groupMap, resp)
	vol, _ := decimal.NewFromString("0.6")
	if !resp.internalVolumeRequired[testlist.TubeTypeSST].Equal(vol) {
		t.Fatalf("volume should be 0.6, get %s", resp.internalVolumeRequired[testlist.TubeTypeSST].String())
	}
}

func TestGetVolumeForTube(t *testing.T) {
	s, resp := setupGetRequiredTubeVolumeService()
	s.getVolumeForTube(resp, testlist.TubeTypeSST, decimal.NewFromInt(500), resp.internalVolumeRequired)
	if resp.internalNumberOfTubes[testlist.TubeTypeSST] != 1 {
		t.Fatalf("number of tubes should be 1, get %d", resp.internalNumberOfTubes[testlist.TubeTypeSST])
	}
	if !resp.internalVolumeRequired[testlist.TubeTypeSST].Equal(decimal.NewFromInt(1500)) {
		t.Fatalf("volume required should be 1500m get %s", resp.internalVolumeRequired[testlist.TubeTypeSST].String())
	}
}
