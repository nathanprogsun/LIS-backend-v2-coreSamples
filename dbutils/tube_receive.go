package dbutils

import (
	"context"
	"coresamples/ent"
	"time"
)

func CreateTubeReceive(sampleId int32,
	tubeType string,
	collectionTime time.Time,
	receivedCount int32,
	receivedBy string,
	receivedTime time.Time,
	isRedraw bool,
	dbClient *ent.Client,
	ctx context.Context) (*ent.TubeReceive, error) {
	return dbClient.TubeReceive.Create().
		SetSampleID(int(sampleId)).
		SetTubeType(tubeType).
		SetCollectionTime(collectionTime).
		SetReceivedTime(receivedTime).
		SetReceivedCount(int(receivedCount)).
		SetReceivedBy(receivedBy).
		SetIsRedraw(isRedraw).Save(ctx)
}
