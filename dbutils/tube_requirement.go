package dbutils

import (
	"context"
	"coresamples/ent"
	"coresamples/ent/sample"
	"coresamples/ent/tuberequirement"
)

func FindTubeRequirement(sampleId int, tubeType string, dbClient *ent.Client, ctx context.Context) (*ent.TubeRequirement, error) {
	return dbClient.TubeRequirement.Query().Where(
		tuberequirement.And(
			tuberequirement.HasSampleWith(sample.ID(sampleId)),
			tuberequirement.TubeTypeEQ(tubeType),
		)).
		First(ctx)
}

func CreateTubeRequirement(sampleId int, tubeType string, requiredCnt int32, requiredBy string, dbClient *ent.Client, ctx context.Context) error {
	return dbClient.TubeRequirement.Create().
		SetSampleID(sampleId).
		SetTubeType(tubeType).
		SetRequiredCount(int(requiredCnt)).
		SetRequiredBy(requiredBy).
		Exec(ctx)
}

func GetSampleRequiredTubes(sampleId int, client *ent.Client, ctx context.Context) ([]*ent.TubeRequirement, error) {
	return client.TubeRequirement.Query().Where(tuberequirement.SampleID(sampleId)).All(ctx)
}
