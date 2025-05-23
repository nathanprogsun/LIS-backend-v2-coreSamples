package dbutils

import (
	"context"
	"coresamples/ent"
	"coresamples/ent/testlist"
	"entgo.io/ent/dialect/sql"
	"github.com/opentracing/opentracing-go"
)

func GetTestsByIds(ids []int, client *ent.Client, ctx context.Context) ([]*ent.TestList, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GetTestsByIds")
	defer span.Finish()
	tests, err := client.TestList.Query().Where(func(s *sql.Selector) {
		s.Where(sql.InInts(testlist.FieldID, ids...))
	}).All(ctx)

	return tests, err
}

func GetTestsByBloodType(blood bool, client *ent.Client, ctx context.Context) ([]*ent.TestList, error) {
	return client.TestList.Query().Where(testlist.BloodTypeEQ(blood)).All(ctx)
}
