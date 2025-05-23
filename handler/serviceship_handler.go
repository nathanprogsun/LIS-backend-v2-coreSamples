package handler

import (
	"context"
	pb "coresamples/proto"
	"coresamples/service"
	"coresamples/util"
)

type MembershipHandler struct {
	MembershipService service.IServiceshipService
}

func (mh *MembershipHandler) SubscriptionAllowed(ctx context.Context, request *pb.ClinicName, response *pb.CheckPermissionResponse) error {
	response.Message = MsgSuccess
	response.Granted = mh.MembershipService.SubscriptionAllowed(request.GetClinicName(), ctx)
	return nil
}

func (mh *MembershipHandler) Subscribe(ctx context.Context, request *pb.SubscribeRequest, response *pb.SimpleResponse) error {
	err := mh.MembershipService.Subscribe(request.AccountId,
		request.AccountType,
		request.SubscriberName,
		request.Email,
		int(request.BillingPlanId),
		request.KeepActive,
		request.Platform,
		request.PaymentToken,
		request.CustomerToken,
		request.PaymentMethodId,
		ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return err
}

func (mh *MembershipHandler) GetAccountSubscriptions(ctx context.Context, request *pb.GetAccountSubscriptionsRequest, response *pb.AccountSubscriptionsInfo) error {
	subs, err := mh.MembershipService.GetAccountSubscriptions(request.AccountId, request.AccountType, request.IncludeOutdated, ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	for _, sub := range subs {
		response.Subscriptions = append(response.Subscriptions, &pb.AccountSubscription{
			Id:             int32(sub.ID),
			AccountId:      sub.AccountID,
			AccountType:    string(sub.AccountType),
			SubscriberName: sub.SubscriberName,
			Email:          sub.Email,
			StartTime:      sub.StartTime.Unix(),
			EndTime:        sub.EndTime.Unix(),
		})
	}
	response.Message = MsgSuccess
	return nil
}

func (mh *MembershipHandler) GetChargingSubscription(ctx context.Context, request *pb.GetChargingSubscriptionRequest, response *pb.ChargingSubscriptionInfo) error {
	csub, err := mh.MembershipService.GetChargingSubscription(int(request.SubscriptionId), ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	err = util.Swap(csub, response.ChargingSubscription)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return nil
}

func (mh *MembershipHandler) AddBillingPlan(ctx context.Context, request *pb.AddBillingPlanRequest, response *pb.SimpleResponse) error {
	err := mh.MembershipService.AddBillingPlan(int(request.ServiceshipId), request.Fee, request.BillingCycle, request.Interval, ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return nil
}

func (mh *MembershipHandler) CreateBillingPlanSet(ctx context.Context, request *pb.CreateBillingPlanSetRequest, response *pb.SimpleResponse) error {
	err := mh.MembershipService.CreateBillingPlanSet(int(request.ServiceshipId),
		request.Fee, request.BillingCycles, request.Intervals, request.Effective, ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return nil
}

func (mh *MembershipHandler) GetLatestBillingPlanSet(ctx context.Context, request *pb.ServiceshipID, response *pb.MembershipBillingPlansInfo) error {
	plans, err := mh.MembershipService.GetLatestBillingPlanSet(int(request.GetServiceshipId()), ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	for _, plan := range plans {
		response.BillingPlans = append(response.BillingPlans, &pb.MembershipBillingPlan{
			Id:            int32(plan.ID),
			Fee:           plan.Fee,
			BillingCycle:  plan.BillingCycle,
			EffectiveTime: plan.EffectiveTime.Unix(),
		})
	}
	response.Message = MsgSuccess
	return nil
}

func (mh *MembershipHandler) GetBillingPlanByID(ctx context.Context, request *pb.BillingPlanID, response *pb.MembershipBillingPlanInfo) error {
	plan, err := mh.MembershipService.GetBillingPlanByID(int(request.Id), ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.BillingPlan = &pb.MembershipBillingPlan{
		Id:            int32(plan.ID),
		Fee:           plan.Fee,
		BillingCycle:  plan.BillingCycle,
		EffectiveTime: plan.EffectiveTime.Unix(),
	}
	response.Message = MsgSuccess
	return nil
}

func (mh *MembershipHandler) GetBillingPlanBySubscriptionID(ctx context.Context, request *pb.SubscriptionID, response *pb.MembershipBillingPlanInfo) error {
	plan, err := mh.MembershipService.GetBillingPlanByID(int(request.Id), ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.BillingPlan = &pb.MembershipBillingPlan{
		Id:            int32(plan.ID),
		Fee:           plan.Fee,
		BillingCycle:  plan.BillingCycle,
		EffectiveTime: plan.EffectiveTime.Unix(),
	}
	response.Message = MsgSuccess
	return nil
}

func (mh *MembershipHandler) GetPaymentMethods(ctx context.Context, request *pb.AccountID, response *pb.PaymentMethodsInfo) error {
	methods, err := mh.MembershipService.GetPaymentMethods(request.GetId(), request.Type, ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	err = util.Swap(methods, &response.PaymentMethods)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return nil
}

func (mh *MembershipHandler) CreatePaymentMethod(ctx context.Context, request *pb.CreatePaymentMethodRequest, response *pb.CreatePaymentMethodResponse) error {
	methodID, err := mh.MembershipService.CreatePaymentMethod(request.AccountId, request.AccountType, request.Platform, request.PaymentToken, request.CustomerToken, request.CardType, request.ExpiryDate, request.LastFour, ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.PaymentMethodId = methodID
	response.Message = MsgSuccess
	return nil
}

func (mh *MembershipHandler) UpdateSubscriptionPaymentMethod(ctx context.Context, request *pb.UpdateSubscriptionPaymentMethodRequest, response *pb.SimpleResponse) error {
	err := mh.MembershipService.UpdateSubscriptionPaymentMethod(int(request.SubscriptionId), request.PaymentMethodId, ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return nil
}

func (mh *MembershipHandler) DeletePaymentMethod(ctx context.Context, request *pb.DeletePaymentMethodRequest, response *pb.SimpleResponse) error {
	err := mh.MembershipService.DeletePaymentMethod(request.PaymentMethodId, request.AccountId, request.AccountType, ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return nil
}

func (mh *MembershipHandler) PauseAutoRenew(ctx context.Context, request *pb.PauseAutoRenewRequest, response *pb.SimpleResponse) error {
	err := mh.MembershipService.PauseAutoRenew(int(request.SubscriptionId), ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return nil
}

func (mh *MembershipHandler) ResumeAutoRenew(ctx context.Context, request *pb.ResumeAutoRenewRequest, response *pb.SimpleResponse) error {
	err := mh.MembershipService.ResumeAutoRenew(int(request.SubscriptionId), request.Platform, request.PaymentToken, request.CustomerToken, request.PaymentMethodId, ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return nil
}

func (mh *MembershipHandler) UpdateSubscriptionBillingPlan(ctx context.Context, request *pb.UpdateSubscriptionBillingPlanRequest, response *pb.SimpleResponse) error {
	err := mh.MembershipService.UpdateSubscriptionBillingPlan(int(request.SubscriptionId), int(request.BillingPlanId), ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return nil
}

func (mh *MembershipHandler) UpdateSubscriptionEmail(ctx context.Context, request *pb.UpdateSubscriptionEmailRequest, response *pb.SimpleResponse) error {
	err := mh.MembershipService.UpdateSubscriptionEmail(request.Email, int(request.SubscriptionId), ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return nil
}

func (mh *MembershipHandler) UpdateAllSubscriptionEmail(ctx context.Context, request *pb.UpdateAllSubscriptionEmailRequest, response *pb.SimpleResponse) error {
	err := mh.MembershipService.UpdateAllSubscriptionEmail(request.Email, request.AccountId, request.AccountType, ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return nil
}

func (mh *MembershipHandler) GetSubscriptionTransactionInfos(ctx context.Context, request *pb.GetSubscriptionTransactionInfosRequest, response *pb.GetSubscriptionTransactionInfosResponse) error {
	txns, err := mh.MembershipService.GetSubscriptionTransactionInfos(int(request.SubscriptionId), ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	err = util.Swap(txns, response.Transactions)
	if err != nil {
		response.Message = err.Error()
		return err
	}

	response.Message = MsgSuccess
	return nil
}
