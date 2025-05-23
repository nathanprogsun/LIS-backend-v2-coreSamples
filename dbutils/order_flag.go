package dbutils

import (
	"context"
	"coresamples/ent"
	"coresamples/ent/orderflag"
	pb "coresamples/proto"
)

// a cache of order flag name to order flags
var orderFlagsByName map[string]*ent.OrderFlag

func GetOrderFlagByName(orderFlagName string, client *ent.Client, ctx context.Context) (*ent.OrderFlag, error) {
	if orderFlagsByName == nil {
		orderFlagsByName = make(map[string]*ent.OrderFlag)
	}
	if flag, exist := orderFlagsByName[orderFlagName]; exist {
		return flag, nil
	}
	flag, err := client.OrderFlag.Query().Where(orderflag.OrderFlagNameEQ(orderFlagName)).Only(ctx)
	if err != nil {
		return nil, err
	}
	orderFlagsByName[orderFlagName] = flag
	return flag, nil
}

func CreateOrderFlag(info *pb.AddOrderFlagRequest, client *ent.Client, ctx context.Context) (*ent.OrderFlag, error) {
	return client.OrderFlag.Create().
		SetOrderFlagName(info.OrderFlagName).
		SetOrderFlagDescription(info.OrderFlagDescription).
		SetOrderFlagAllowDuplicatesUnderSameCategory(info.OrderFlagAllowDuplicatesUnderSameCategory).
		SetOrderFlagCategory(info.OrderFlagCategory).
		SetOrderFlagedBy(info.OrderFlagedBy).
		Save(ctx)
}

func GetAllOrderFlags(client *ent.Client, ctx context.Context) ([]*ent.OrderFlag, error) {
	return client.OrderFlag.Query().All(ctx)
}
