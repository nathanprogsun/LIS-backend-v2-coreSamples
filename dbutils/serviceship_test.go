package dbutils

import (
	"context"
	"coresamples/common"
	"coresamples/ent"
	"coresamples/ent/enttest"
	"coresamples/ent/serviceship"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"testing"
	"time"
)

const (
	BasicMembership = "BasicMembership"
	Membership      = "membership"
)

func setupMembershipDB(t *testing.T) *ent.Client {
	dataSource := "file:ent?mode=memory&cache=shared&_fk=1"
	dbClient := enttest.Open(t, "sqlite3", dataSource, enttest.WithOptions(ent.Log(t.Log)))

	err := dbClient.Schema.Create(context.Background())

	if err != nil {
		common.Fatal(err)
	}
	return dbClient
}

func TestCreateServiceship(t *testing.T) {
	client := setupMembershipDB(t)
	defer client.Close()

	ctx := context.Background()
	err := CreateServiceship(
		BasicMembership,
		"membership",
		[]float32{20, 30},
		[]int32{1, 6},
		[]string{"monthly", "monthly"},
		client,
		ctx)
	if err != nil {
		t.Fatal(err)
	}

	m, err := client.Serviceship.Query().Where(serviceship.Tag(BasicMembership)).WithServiceshipBillingPlan().Only(ctx)
	if err != nil {
		t.Fatal(err)
	}

	plans, err := m.Edges.ServiceshipBillingPlanOrErr()
	if err != nil {
		t.Fatal(err)
	}

	if len(plans) != 2 {
		t.Fatal("should have 2 plans")
	}
	for _, plan := range plans {
		pm, err := client.ServiceshipBillingPlan.QueryServiceship(plan).Only(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if pm.Tag != m.Tag || pm.Type != m.Type {
			t.Fatal("serviceship does not match")
		}
	}
}

func TestAddBillingPlan(t *testing.T) {
	client := setupMembershipDB(t)
	defer client.Close()

	ctx := context.Background()
	err := CreateServiceship(
		BasicMembership,
		"membership",
		[]float32{20, 30},
		[]int32{1, 6},
		[]string{"monthly", "monthly"},
		client,
		ctx)
	if err != nil {
		t.Fatal(err)
	}

	m, err := GetServiceshipByTag(BasicMembership, "membership", client, ctx)
	if err != nil {
		t.Fatal(err)
	}
	plans, err := GetLatestBillingPlanSet(m.ID, client, ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(plans) != 2 {
		fmt.Println(plans)
		t.Fatal("should be 2 billing plans")
	}

	err = AddBillingPlan(m.ID, 40, 6, "monthly", client, ctx)
	if err == nil {
		t.Fatal("should not be able to add duplicate billing cycle")
	}

	err = AddBillingPlan(m.ID, 40, 12, "monthly", client, ctx)
	if err != nil {
		t.Fatal("should be able to create billing cycle")
	}

	plans, err = GetLatestBillingPlanSet(m.ID, client, ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(plans) != 3 {
		t.Fatal("should be 3 billing plans")
	}

	now := time.Now()
	err = CreateBillingPlanSet(m.ID, []float32{10, 20, 30, 40}, []int32{1, 2, 3, 4}, []string{"monthly", "monthly", "monthly", "monthly"}, now, client, ctx)

	if err != nil {
		t.Fatal(err)
	}

	plans, err = GetLatestBillingPlanSet(m.ID, client, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(plans) != 4 {
		t.Fatal("should be 4 billing plans")
	}

	for _, plan := range plans {
		if plan.EffectiveTime.Equal(now) {
			continue
		}
		t.Fatal("effective time does not match")
	}
}

func TestCreateSubscription(t *testing.T) {
	client := setupMembershipDB(t)
	defer client.Close()
	ctx := context.Background()

	err := CreateServiceship(
		BasicMembership,
		Membership,
		[]float32{20, 30},
		[]int32{1, 6}, []string{"monthly", "monthly"},
		client,
		ctx)
	if err != nil {
		t.Fatal(err)
	}

	ms1, err := GetServiceshipByTag(BasicMembership, "membership", client, ctx)
	mbs1, err := GetLatestBillingPlanSet(ms1.ID, client, ctx)
	if err != nil {
		t.Fatal(err)
	}

	s, err := CreateSubscription(0, "clinic", "a", "a@gmail.com", mbs1[0].ID, false, client, ctx)
	if err != nil {
		t.Fatal(err)
	}
	b, err := GetBillingPlanBySubscriptionID(s.ID, client, ctx)
	if err != nil {
		t.Fatal(err)
	}
	bf, _ := json.Marshal(b)
	common.Infof(string(bf))
	if b.ID != mbs1[0].ID {
		t.Fatal("wrong billing plan")
	}

	_, err = CreateSubscription(0, "clinic", "a", "a@gmail.com", mbs1[0].ID, false, client, ctx)
	if err == nil {
		t.Fatal("should not be able to add membership that's still valid")
	}

	err = CancelSubscription(0, "clinic", ms1.ID, client, ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = CreateSubscription(0, "clinic", "a", "a@gmail.com", mbs1[0].ID, false, client, ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCancelSubscription(t *testing.T) {
	client := setupMembershipDB(t)
	defer client.Close()
	ctx := context.Background()

	err := CreateServiceship(
		BasicMembership,
		Membership,
		[]float32{20, 30},
		[]int32{1, 6}, []string{"monthly", "monthly"},
		client,
		ctx)
	if err != nil {
		t.Fatal(err)
	}

	ms1, err := GetServiceshipByTag(BasicMembership, "membership", client, ctx)
	mbs1, err := GetLatestBillingPlanSet(ms1.ID, client, ctx)
	if err != nil {
		t.Fatal(err)
	}

	_, err = CreateSubscription(0, "clinic", "a", "a@gmail.com", mbs1[0].ID, false, client, ctx)
	if err != nil {
		t.Fatal(err)
	}
	subs, err := GetAccountSubscriptions(0, "clinic", false, client, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(subs) != 1 {
		t.Fatal(errors.New("account 0 should have 1 subscription"))
	}

	sb, _ := json.Marshal(subs[0])
	common.Infof(string(sb))

	err = CancelSubscription(0, "clinic", ms1.ID, client, ctx)
	if err != nil {
		t.Fatal(err)
	}

	subs, err = GetAccountSubscriptions(0, "clinic", false, client, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(subs) != 0 {
		t.Fatal(errors.New("account 0 should have no subscription"))
	}
}

func TestUpdateSubscriptionBillingPlan(t *testing.T) {
	client := setupMembershipDB(t)
	defer client.Close()
	ctx := context.Background()

	err := CreateServiceship(
		BasicMembership,
		Membership,
		[]float32{20, 30},
		[]int32{1, 6}, []string{"monthly", "monthly"},
		client,
		ctx)
	if err != nil {
		t.Fatal(err)
	}

	ms1, err := GetServiceshipByTag(BasicMembership, "membership", client, ctx)
	mbs1, err := GetLatestBillingPlanSet(ms1.ID, client, ctx)
	if err != nil {
		t.Fatal(err)
	}

	_, err = CreateSubscription(0, "clinic", "a", "a@gmail.com", mbs1[0].ID, false, client, ctx)
	if err != nil {
		t.Fatal(err)
	}
	subs, err := GetAccountSubscriptions(0, "clinic", false, client, ctx)
	bp, err := subs[0].QueryServiceshipBillingPlan().Only(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if bp.Fee != 20 {
		t.Fatal("wrong billing plan")
	}

	err = UpdateSubscriptionBillingPlan(subs[0].ID, mbs1[1].ID, client, ctx)
	if err != nil {
		t.Fatal(err)
	}

	bp, err = subs[0].QueryServiceshipBillingPlan().Only(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if bp.Fee != 30 {
		t.Fatal("wrong billing plan")
	}
}

func TestAddSubscriptionEndTime(t *testing.T) {
	client := setupMembershipDB(t)
	defer client.Close()
	ctx := context.Background()

	err := CreateServiceship(
		BasicMembership,
		Membership,
		[]float32{20, 30},
		[]int32{1, 6}, []string{"monthly", "monthly"},
		client,
		ctx)
	if err != nil {
		t.Fatal(err)
	}

	ms1, err := GetServiceshipByTag(BasicMembership, "membership", client, ctx)
	mbs1, err := GetLatestBillingPlanSet(ms1.ID, client, ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = CreateSubscription(0, "clinic", "a", "a@gmail.com", mbs1[0].ID, false, client, ctx)
	if err != nil {
		t.Fatal(err)
	}
	subs, _ := GetAccountSubscriptions(0, "clinic", false, client, ctx)
	endtime := subs[0].EndTime
	err = AddSubscriptionEndTime(12, subs[0].ID, client, ctx)
	if err != nil {
		t.Fatal(err)
	}
	subs, _ = GetAccountSubscriptions(0, "clinic", false, client, ctx)
	if !endtime.AddDate(1, 0, 0).Equal(subs[0].EndTime) {
		t.Fatal("wrong end time")
	}
}

func TestGetSubscriptionsByEndTime(t *testing.T) {
	client := setupMembershipDB(t)
	defer client.Close()
	ctx := context.Background()

	err := CreateServiceship(
		BasicMembership,
		Membership,
		[]float32{20, 30},
		[]int32{1, 6}, []string{"monthly", "monthly"},
		client,
		ctx)
	if err != nil {
		t.Fatal(err)
	}

	ms1, err := GetServiceshipByTag(BasicMembership, "membership", client, ctx)
	mbs1, err := GetLatestBillingPlanSet(ms1.ID, client, ctx)
	if err != nil {
		t.Fatal(err)
	}
	sub, _ := CreateSubscription(0, "clinic", "a", "a@gmail.com", mbs1[0].ID, false, client, ctx)
	bf, _ := json.Marshal(sub)
	fmt.Println(string(bf))
	expire := sub.EndTime.AddDate(0, 0, 1)
	subs, err := GetSubscriptionsByEndTime(expire, client, ctx)
	if len(subs) != 1 {
		t.Fatal("should find 1 subscription")
	}

	expire = sub.EndTime.AddDate(0, 0, 2)
	subs, err = GetSubscriptionsByEndTime(expire, client, ctx)
	if len(subs) != 0 {
		t.Fatal("should find 0 subscription")
	}

	expire = sub.EndTime
	subs, err = GetSubscriptionsByEndTime(expire, client, ctx)
	if len(subs) != 0 {
		t.Fatal("should find 0 subscription")
	}
}

//func TestMaybeMatchNewestBillingPlan(t *testing.T) {
//	client := setupMembershipDB(t)
//	defer client.Close()
//	ctx := context.Background()
//
//	err := CreateServiceship(
//		BasicMembership,
//		"",
//		false,
//		0.02,
//		0.01,
//		[]int32{20, 30},
//		[]int32{1, 6}, []string{"monthly", "monthly"},
//		client,
//		ctx)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	ms1, err := GetServiceshipByTag(BasicMembership, client, ctx)
//	mbs1, err := GetLatestBillingPlanSet(ms1.ID, client, ctx)
//	CreateSubscription(0, "a", "a@gmail.com",  true, mbs1[0].ID, client, ctx)
//	subs, _ := GetAccountSubscriptions(0, false,  client, ctx)
//	bp, err := MaybeMatchNewestBillingPlan(subs[0], client, ctx)
//	if bp.MonthlyFee != 20 {
//		t.Fatal("wrong billing plan")
//	}
//
//	now := time.Now()
//
//	CreateBillingPlanSet(ms1.ID, []int32{10, 40}, []int32{1, 12}, now.AddDate(1, 0, 0), client, ctx)
//	bp, err = MaybeMatchNewestBillingPlan(subs[0], client, ctx)
//	if bp.MonthlyFee != 20 {
//		t.Fatal("wrong billing plan")
//	}
//
//	CreateBillingPlanSet(ms1.ID, []int32{10, 40}, []int32{1, 12}, now, client, ctx)
//	bp, err = MaybeMatchNewestBillingPlan(subs[0], client, ctx)
//	if bp.MonthlyFee != 10 {
//		t.Fatal("wrong billing plan3")
//	}
//
//	bp1, _ := GetLatestEffectiveBillingPlanSet(ms1.ID, client, ctx)
//	bp2, _ := GetLatestBillingPlanSet(ms1.ID, client, ctx)
//	if !bp1[0].EffectiveTime.Equal(now) || !bp2[0].EffectiveTime.Equal(now.AddDate(1, 0, 0)) {
//		t.Fatal("find wrong billing plans")
//	}
//}
