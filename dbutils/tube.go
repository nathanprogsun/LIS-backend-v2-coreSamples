package dbutils

import (
	"context"
	"coresamples/ent"
	"coresamples/ent/tube"
)

func GetTubeByTubeID(tubeId string, client *ent.Client, ctx context.Context) (*ent.Tube, error) {
	return client.Tube.Query().Where(tube.TubeIDEQ(tubeId)).WithTubeType().Only(ctx)
}
