package handler

import (
	"context"
	"coresamples/ent"
	pb "coresamples/proto"
	"coresamples/service"
	"coresamples/util"
	"strconv"
)

type OrderHandler struct {
	OrderService service.IOrderService
}

func (oh *OrderHandler) GetOrder(ctx context.Context, req *pb.OrderID, resp *pb.Order) error {
	order, err := oh.OrderService.GetOrder(int(req.OrderID), ctx)
	if err != nil {
		return err
	}
	err = util.Swap(order, resp)
	if err != nil {
		return err
	}
	err = util.Swap(order.Edges.OrderFlags, &resp.OrderFlags)
	if err != nil {
		return err
	}

	samples := []*ent.Sample{order.Edges.Sample}
	err = util.Swap(samples, &resp.Samples)
	if err != nil {
		return err
	}
	return nil
}

func (oh *OrderHandler) CancelOrder(ctx context.Context, req *pb.OrderID, resp *pb.CancelOrderResponse) error {
	order, err := oh.OrderService.CancelOrder(int(req.OrderID), ctx)
	if err != nil {
		return err
	}
	err = util.Swap(order, resp)
	if err != nil {
		return err
	}
	return nil
}

func (oh *OrderHandler) RestoreCanceledOrder(ctx context.Context, req *pb.OrderID, resp *pb.RestoreOrderResponse) error {
	order, err := oh.OrderService.RestoreCanceledOrder(int(req.OrderID), ctx)
	if err != nil {
		return err
	}
	err = util.Swap(order, resp)
	if err != nil {
		return err
	}
	return nil
}

func (oh *OrderHandler) AddOrderFlag(ctx context.Context, req *pb.AddOrderFlagRequest, resp *pb.AddOrderFlagResponse) error {
	flag, err := oh.OrderService.AddOrderFlag(req, ctx)
	if err != nil {

		return nil
	}
	err = util.Swap(flag, resp)
	if err != nil {
		return err
	}
	return nil
}

func (oh *OrderHandler) FlagOrder(ctx context.Context, req *pb.FlagOrderRequest, resp *pb.FlagOrderResponse) error {
	orderId, err := strconv.Atoi(req.OrderId)
	if err != nil {
		return nil
	}
	order, err := oh.OrderService.FlagOrder(orderId, req.OrderFlagNames, ctx)
	err = util.Swap(order, resp)
	if err != nil {
		return err
	}
	return nil
}

func (oh *OrderHandler) UnflagOrder(ctx context.Context, req *pb.UnflagOrderRequest, resp *pb.UnflagOrderResponse) error {
	orderId, err := strconv.Atoi(req.OrderId)
	if err != nil {
		return nil
	}
	order, err := oh.OrderService.UnflagOrder(orderId, req.OrderFlagNames, ctx)
	err = util.Swap(order, resp)
	if err != nil {
		return err
	}
	return nil
}

func (oh *OrderHandler) ListOrderFlagTypes(ctx context.Context, req *pb.ListOrderFlagTypesRequest, resp *pb.ListOrderFlagTypesResponse) error {
	flags, err := oh.OrderService.ListOrderFlagTypes(ctx)
	if err != nil {
		return nil
	}
	err = util.Swap(flags, &resp.OrderFlags)
	if err != nil {
		return err
	}
	return nil
}

func (oh *OrderHandler) ChangeOrderStatus(ctx context.Context, req *pb.ChangeOrderStatusRequest, resp *pb.Order) error {
	orderId, err := strconv.Atoi(req.OrderId)
	if err != nil {
		return nil
	}
	order, err := oh.OrderService.ChangeOrderStatus(orderId, req.Status, ctx)
	if err != nil {
		return nil
	}
	err = util.Swap(order, resp)
	if err != nil {
		return err
	}
	return nil
}

func (oh *OrderHandler) GetOrderStatusForDisplay(ctx context.Context, req *pb.GetOrderStatusDisplay, resp *pb.GetOrderStatusDisplayResponse) error {
	of, pf, err := oh.OrderService.GetOrderStatusForDisplay(ctx)
	if err != nil {
		return err
	}
	err = util.Swap(of, &resp.OrderFlag)
	if err != nil {
		return err
	}

	err = util.Swap(pf, &resp.PatientFlag)
	if err != nil {
		return err
	}
	return nil
}

func (oh *OrderHandler) RerunSampleTests(ctx context.Context, req *pb.RerunSampleTestsRequest, resp *pb.RerunSampleTestsResponse) error {
	//TODO: implement this
	return nil
}

func (oh *OrderHandler) RestoreOrderStatus(ctx context.Context, req *pb.RestoreOrderStatusRequest, resp *pb.RestoreOrderStatusResponse) error {
	//TODO: implement this
	return nil
}
