package dbutils

import (
	"context"
	"coresamples/ent"
	"coresamples/ent/sampletype"
	"coresamples/ent/testlist"
	"coresamples/ent/tubeinstructions"
	"coresamples/ent/tubetype"
)

func GetTubeInfoByEnum(tubeType testlist.TubeType, client *ent.Client, ctx context.Context) (*ent.TubeInstructions, error) {
	return client.TubeInstructions.
		Query().
		Where(tubeinstructions.TubeNameEnumEQ(tubeType.String())).
		Only(ctx)
}

func GetTubeTypeInfoByTubeTypeEnum(tubeType string, client *ent.Client, ctx context.Context) ([]*ent.TubeType, error) {
	return client.TubeType.Query().Where(tubetype.TubeTypeEnumEQ(tubeType)).WithSampleTypes().WithTests().All(ctx)
}

func GetTubeTypeInfoBySampleTypeCode(sampleTypeCode string, client *ent.Client, ctx context.Context) (*ent.SampleType, error) {
	return client.SampleType.Query().Where(sampletype.SampleTypeCodeEQ(sampleTypeCode)).WithTubeTypes().Only(ctx)
}

func GetTubeTypeBySampleTypeEnum(sampleTypeEnum string, client *ent.Client, ctx context.Context) (*ent.SampleType, error) {
	return client.SampleType.Query().Where(sampletype.SampleTypeEnumEQ(sampleTypeEnum)).WithTubeTypes().First(ctx)
}
