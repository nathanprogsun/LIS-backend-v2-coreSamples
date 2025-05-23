package dbutils

import (
	"context"
	"coresamples/ent"
)

func GetAllPatientFlags(client *ent.Client, ctx context.Context) ([]*ent.PatientFlag, error) {
	return client.PatientFlag.Query().All(ctx)
}
