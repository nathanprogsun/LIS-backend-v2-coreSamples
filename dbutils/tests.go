package dbutils

import (
	"context"
	"coresamples/ent"
	"coresamples/ent/test"
	"coresamples/ent/testdetail"
	pb "coresamples/proto"
	"entgo.io/ent/dialect/sql"
)

func GetTestsByTestIds(testIds []int, client *ent.Client, ctx context.Context) ([]*ent.Test, error) {
	return client.Test.Query().Where(func(s *sql.Selector) {
		s.Where(sql.InInts(test.FieldID, testIds...))
	}).WithTestDetails().All(ctx)
}

func GetTestsWithFields(testIds []int, testDetailNames []string, client *ent.Client, ctx context.Context) ([]*ent.Test, error) {
	return client.Test.Query().Where(func(s *sql.Selector) {
		s.Where(sql.InInts(test.FieldID, testIds...))
	}).WithTestDetails(func(q *ent.TestDetailQuery) {
		q.Where(testdetail.TestDetailNameIn(testDetailNames...))
	}).All(ctx)
}

func GetAllTests(client *ent.Client, ctx context.Context) ([]*ent.Test, error) {
	return client.Test.Query().WithTestDetails().All(ctx)
}

func CreateTest(test *pb.CreateTestRequest, client *ent.Client, ctx context.Context) (*ent.Test, error) {
	return client.Test.Create().
		SetIsActive(test.IsActive).
		SetTestName(test.TestName).
		SetTestCode(test.TestCode).
		SetDisplayName(test.DisplayName).
		SetTestDescription(test.TestDescription).
		SetAssayName(test.AssayName).
		Save(ctx)
}

func GetTestIdsByTestCode(code string, client *ent.Client, ctx context.Context) ([]int, error) {
	return client.Test.Query().Where(test.TestCodeEQ(code)).Select(test.FieldID).Ints(ctx)
}

func GetTestByTestId(testId int, client *ent.Client, ctx context.Context) (*ent.Test, error) {
	return client.Test.Get(ctx, testId)
}

func GetTestByTestIdWithDetailName(testId int, testDetailName string, client *ent.Client, ctx context.Context) (*ent.Test, error) {
	return client.Test.Query().
		Where(test.IDEQ(testId)).
		WithTestDetails(
			func(q *ent.TestDetailQuery) {
				q.Where(testdetail.TestDetailNameEQ(testDetailName))
			}).
		Only(ctx)
}

func GetTestDetailsByDetailValueAndDetailName(testDetailValue string, testDetailName string, client *ent.Client, ctx context.Context) ([]*ent.TestDetail, error) {
	return client.TestDetail.Query().
		Where(testdetail.And(
			testdetail.TestDetailsValueEQ(testDetailValue),
			testdetail.TestDetailNameEQ(testDetailName),
		)).
		All(ctx)
}
