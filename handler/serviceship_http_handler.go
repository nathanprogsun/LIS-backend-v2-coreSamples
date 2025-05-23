package handler

import (
	pb "coresamples/proto"
	"coresamples/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ServiceshipHTTPHandler struct {
	service service.IServiceshipService
}

type UpdateSubscriptionRequest struct {
	SubscriptionID  int    `json:"subscription_id,omitempty"`
	AccountID       int64  `json:"account_id,omitempty"`
	AccountType     string `json:"account_type,omitempty"`
	PaymentMethodID int64  `json:"payment_method_id,omitempty"`
	Email           string `json:"email,omitempty"`
	UpdateAllEmail  bool   `json:"update_all_email,omitempty"`
	BillingPlanID   int    `json:"billing_plan_id,omitempty"`
}

type CreateServiceRequest struct {
	Tag           string    `json:"tag,omitempty"`
	Type          string    `json:"type,omitempty"`
	Fees          []float32 `json:"fees,omitempty"`
	BillingCycles []int32   `json:"billing_cycles,omitempty"`
	Intervals     []string  `json:"intervals,omitempty"`
}

func NewServiceshipHTTPHandler(service service.IServiceshipService) *ServiceshipHTTPHandler {
	return &ServiceshipHTTPHandler{
		service: service,
	}
}

func (s *ServiceshipHTTPHandler) CreateServiceship(c *gin.Context) {
	var err error
	req := &CreateServiceRequest{}
	if err = c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	err = s.service.CreateServiceship(req.Tag, req.Type, req.Fees, req.BillingCycles, req.Intervals, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (s *ServiceshipHTTPHandler) SubscriptionAllowed(c *gin.Context) {
	clinicName := c.Param("clinic_name")
	resp := &pb.CheckPermissionResponse{
		Granted: s.service.SubscriptionAllowed(clinicName, c),
	}
	c.JSON(http.StatusOK, resp)
}

func (s *ServiceshipHTTPHandler) GetServiceships(c *gin.Context) {
	mems, err := s.service.GetServiceships(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error getting memberships": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mems)
}

func (s *ServiceshipHTTPHandler) GetServiceshipsByType(c *gin.Context) {
	svcType := c.Query("type")
	svcships, err := s.service.GetServiceshipsByType(svcType, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, svcships)
}

func (s *ServiceshipHTTPHandler) GetServiceshipByID(c *gin.Context) {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	svc, _, err := s.service.GetServiceshipConfig(id, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error getting membership": err.Error()})
		return
	}
	c.JSON(http.StatusOK, svc)
}

func (s *ServiceshipHTTPHandler) Subscribe(c *gin.Context) {
	var err error
	subreq := &pb.SubscribeRequest{}
	if err = c.ShouldBindJSON(subreq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	if subreq.SubscriberName == "" || subreq.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": "account info should not be empty"})
		return
	}
	err = s.service.Subscribe(subreq.AccountId,
		subreq.AccountType,
		subreq.SubscriberName,
		subreq.Email,
		int(subreq.BillingPlanId),
		subreq.KeepActive,
		subreq.Platform,
		subreq.PaymentToken,
		subreq.CustomerToken,
		subreq.PaymentMethodId, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (s *ServiceshipHTTPHandler) GetAccountSubscriptions(c *gin.Context) {
	var includeOutdated bool
	if include := c.DefaultQuery("include_outdated", "0"); include == "0" {
		includeOutdated = false
	} else if include == "1" {
		includeOutdated = true
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": "invalid include_outdated flag"})
		return
	}
	accountID, err := getAccountIDFromParam(c, c.Query("account_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	accountType := c.Query("account_type")
	subs, err := s.service.GetAccountSubscriptions(accountID, accountType, includeOutdated, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subs)
}

func (s *ServiceshipHTTPHandler) GetChargingSubscription(c *gin.Context) {
	subscriptionID, err := strconv.Atoi(c.Query("subscription_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	info, err := s.service.GetChargingSubscription(subscriptionID, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

func (s *ServiceshipHTTPHandler) CreateBillingPlanSet(c *gin.Context) {
	req := &pb.CreateBillingPlanSetRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	err := s.service.CreateBillingPlanSet(int(req.ServiceshipId), req.Fee, req.BillingCycles, req.Intervals, req.Effective, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (s *ServiceshipHTTPHandler) AddBillingPlan(c *gin.Context) {
	req := &pb.AddBillingPlanRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	err := s.service.AddBillingPlan(int(req.ServiceshipId), req.Fee, req.BillingCycle, req.Interval, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (s *ServiceshipHTTPHandler) GetLatestBillingPlanSet(c *gin.Context) {
	memID, err := strconv.Atoi(c.Query("membership_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	plans, err := s.service.GetLatestBillingPlanSet(memID, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plans)
}

func (s *ServiceshipHTTPHandler) GetPaymentMethods(c *gin.Context) {
	accountID, err := strconv.ParseInt(c.Query("account_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	accountType := c.Query("account_type")
	methods, err := s.service.GetPaymentMethods(accountID, accountType, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, methods)
}

func (s *ServiceshipHTTPHandler) CreatePaymentMethod(c *gin.Context) {
	var err error
	req := &pb.CreatePaymentMethodRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	paymentID, err := s.service.CreatePaymentMethod(req.AccountId, req.AccountType, req.Platform, req.PaymentToken, req.CustomerToken, req.CardType, req.ExpiryDate, req.LastFour, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	resp := struct {
		PaymentMethodID int64 `json:"payment_method_id,omitempty"`
	}{PaymentMethodID: paymentID}
	c.JSON(http.StatusOK, resp)
}

func (s *ServiceshipHTTPHandler) DeletePaymentMethod(c *gin.Context) {
	accountID, err := getAccountIDFromParam(c, c.Query("account_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	accountType := c.Query("account_type")
	paymentMethodID, err := strconv.ParseInt(c.Query("payment_method_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	err = s.service.DeletePaymentMethod(paymentMethodID, accountID, accountType, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (s *ServiceshipHTTPHandler) UpdateSubscription(c *gin.Context) {
	var err error
	req := &UpdateSubscriptionRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	// these updates require subscription id
	if (req.Email != "" && !req.UpdateAllEmail) || req.PaymentMethodID != 0 || req.BillingPlanID != 0 {
		if req.SubscriptionID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errorMessage": "subscription id should not be empty"})
			return
		}
	}
	if req.UpdateAllEmail {
		if req.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"errorMessage": "email should not be empty"})
			return
		}
		err = s.service.UpdateAllSubscriptionEmail(req.Email, req.AccountID, req.AccountType, c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
			return
		}
	} else if req.Email != "" {
		err = s.service.UpdateSubscriptionEmail(req.Email, req.SubscriptionID, c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
			return
		}
	}

	if req.BillingPlanID != 0 {
		err = s.service.UpdateSubscriptionBillingPlan(req.SubscriptionID, req.BillingPlanID, c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
			return
		}
	}

	if req.PaymentMethodID != 0 {
		err = s.service.UpdateSubscriptionPaymentMethod(req.SubscriptionID, req.PaymentMethodID, c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, nil)
}

func (s *ServiceshipHTTPHandler) PauseAutoRenew(c *gin.Context) {
	subscriptionID, err := strconv.Atoi(c.Query("subscription_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	err = s.service.PauseAutoRenew(subscriptionID, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (s *ServiceshipHTTPHandler) ResumeAutoRenew(c *gin.Context) {
	var err error
	req := &pb.ResumeAutoRenewRequest{}
	if err = c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	err = s.service.ResumeAutoRenew(int(req.SubscriptionId), req.Platform, req.PaymentToken, req.CustomerToken, req.PaymentMethodId, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (s *ServiceshipHTTPHandler) GetSubscriptionTransactions(c *gin.Context) {
	subscriptionID, err := strconv.Atoi(c.Query("subscription_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	info, err := s.service.GetSubscriptionTransactionInfos(subscriptionID, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

//func getClinicIDFromToken(c *gin.Context) (int64, error) {
//	accountID := c.GetString("account_id")
//	return strconv.ParseInt(accountID, 10, 64)
//}
//
//func getClinicIDFromTokenOrParam(c *gin.Context, id string) (int64, error) {
//	var accountID int64
//	var err error
//	if accountID, err = strconv.ParseInt(id, 10, 64); err != nil {
//		if accountID, err = getClinicIDFromToken(c); err != nil {
//			return 0, err
//		} else {
//			return accountID, nil
//		}
//	} else {
//		return accountID, nil
//	}
//}

func getAccountIDFromParam(c *gin.Context, id string) (int64, error) {
	var accountID int64
	var err error
	if accountID, err = strconv.ParseInt(id, 10, 64); err != nil {
		return 0, err
	} else {
		return accountID, nil
	}
}
