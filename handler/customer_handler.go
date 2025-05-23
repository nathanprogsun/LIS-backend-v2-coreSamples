package handler

import (
	"context"
	"coresamples/model"
	pb "coresamples/proto"
	"coresamples/service"
	"coresamples/util"
)

type CustomerHandler struct {
	CustomerService service.ICustomerService
}

func (h *CustomerHandler) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest, resp *pb.Customer) error {
	return nil
}

func (h *CustomerHandler) ListCustomers(ctx context.Context, req *pb.CustomerPage, resp *pb.ListCustomersResponse) error {
	customers, hasNextPage, total, err := h.CustomerService.ListCustomer(req.Page, req.PerPage, ctx)
	if err != nil {
		return err
	}
	err = util.Swap(customers, &resp.Customers)
	if err != nil {
		return err
	}
	err = util.Swap(hasNextPage, &resp.HasNextPage)
	if err != nil {
		return err
	}
	err = util.Swap(total, &resp.Total)
	if err != nil {
		return err
	}
	return nil
}

func (h *CustomerHandler) GetCustomer(ctx context.Context, req *pb.CustomerID, resp *pb.FullCustomer) error {
	customer, err := h.CustomerService.GetCustomer(int(req.CustomerId), ctx)
	if err != nil {
		return err
	}
	err = util.Swap(customer, resp)
	if err != nil {
		return err
	}
	return nil
}

func (h *CustomerHandler) GetSalesCustomers(ctx context.Context, req *pb.SalesInfo, resp *pb.FullCustomerList) error {
	result, err := h.CustomerService.GetSalesCustomer(req.SalesName, req.Page, req.PerPage, ctx)
	if err != nil {
		return err
	}
	err = util.Swap(result, &resp.Customers)
	if err != nil {
		return err
	}

	return nil
}

func (h *CustomerHandler) GetClinicSalesSamples(ctx context.Context, req *pb.GetSampleDataByPracticeAndSalesRequest, resp *pb.SampleDataByPracticeAndSalesResponse) error {
	return nil
}

func (h *CustomerHandler) GetCustomerSetting(ctx context.Context, req *pb.GetCustomerSettingRequest, resp *pb.GetCustomerSettingResponse) error {
	return nil
}

func (h *CustomerHandler) GetCustomerSales(ctx context.Context, req *pb.GetCustomerSalesRequest, resp *pb.ListSalesCustomerResponseV7) error {
	customerSales, err := h.CustomerService.GetCustomerSales(req.CustomerNames, req.CustomerIds, ctx)
	if err != nil {
		return err
	}
	var saleDetail []*pb.SaleDetailcWithCustomerV7
	for _, item := range customerSales {
		saleDetail = append(saleDetail, &pb.SaleDetailcWithCustomerV7{
			CustomerId:         item.CustomerId,
			CustomerFirstName:  item.CustomerFirstName,
			CustomerLastName:   item.CustomerLastName,
			CustomerMiddleName: item.CustomerMiddleName,
			InternalUser: &pb.SaleDetailcV7{
				InternalUserRoleId:     item.InternalUser.InternalUserRoleId,
				InternalUserFirstname:  item.InternalUser.InternalUserFirstname,
				InternalUserLastname:   item.InternalUser.InternalUserLastname,
				InternalUserMiddlename: item.InternalUser.InternalUserMiddlename,
				InternalUserEmail:      item.InternalUser.InternalUserEmail,
				InternalUserPhone:      item.InternalUser.InternalUserPhone,
			},
		})
	}

	resp.Sales = []*pb.SalesCustomerV7{
		{CustomerSales: saleDetail},
	}

	return nil
}

func (h *CustomerHandler) UpdateCustomerSetting(ctx context.Context, req *pb.UpdateCustomerSettingRequest, resp *pb.UpdateCustomerSettingResponse) error {
	return nil
}

func (h *CustomerHandler) UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest, resp *pb.UpdateCustomerResponse) error {
	return nil
}

func (h *CustomerHandler) CreatePatientInternalNotes(ctx context.Context, req *pb.CreatePatientInternalNotesRequest, resp *pb.CreatePatientInternalNotesResponse) error {
	return nil
}

func (h *CustomerHandler) ModifyPatientInternalNotes(ctx context.Context, req *pb.ModifyPatientInternalNotesRequest, resp *pb.ModifyPatientInternalNotesResponse) error {
	return nil
}

func (h *CustomerHandler) DeletePatientInternalNotes(ctx context.Context, req *pb.DeletePatientInternalNotesRequest, resp *pb.DeletePatientInternalNotesResponse) error {
	return nil
}

// Version 0.7.3.7
func (h *CustomerHandler) IsNewCustomer(ctx context.Context, req *pb.CustomerID, resp *pb.GetIsNewCustomerResponse) error {
	return nil
}

// Version 0.7.3.9
func (h *CustomerHandler) UpdateCustomerNPI(ctx context.Context, req *pb.UpdateCustomerNPIRequest, resp *pb.Customer) error {
	return nil
}

func (h *CustomerHandler) UpdateCustomerSettingFull(ctx context.Context, req *pb.UpdateCustomerSettingFullRequest, resp *pb.UpdateCustomerSettingResponse) error {
	return nil
}

func (h *CustomerHandler) EditCustomerSettingProperties(ctx context.Context, req *pb.EditCustomerSettingPropertiesequest, resp *pb.UpdateCustomerSettingResponse) error {
	return nil
}

func (h *CustomerHandler) EditCustomerProfileOnSettingPage(ctx context.Context, req *pb.EditCustomerProfileOnSettingPageRequest, resp *pb.EditCustomerProfileOnSettingPageResponse) error {
	return nil
}

func (h *CustomerHandler) RemoveCustomerFromClinic(ctx context.Context, req *pb.RemoveCustomerFromClinicRequest, resp *pb.EditCustomerProfileOnSettingPageResponse) error {
	var updatedBy string = "Unknown"

	// comment out metadata log temporarily

	// var xRequestID string = uuid.NewString()
	// var serviceName string = "Unknown Caller"

	// md, ok := metadata.FromIncomingContext(ctx)
	// if ok {
	// 	if val := md.Get("user"); len(val) > 0 {
	// 		updatedBy = val[0]
	// 	}
	// 	if val := md.Get("x-request-id"); len(val) > 0 {
	// 		xRequestID = val[0]
	// 	}
	// 	if val := md.Get("service-name"); len(val) > 0 {
	// 		serviceName = val[0]
	// 	}
	// }
	// common.Infof("[xRequestID:%s] RemoveCustomerFromClinic called by: %s, input: %+v", xRequestID, serviceName, req)

	status, errorLog, code := h.CustomerService.RemoveCustomerFromClinic(req.CustomerId, req.ClinicId, updatedBy, ctx)

	// Set response
	resp.UpdateStatus = status
	resp.ErrorLog = errorLog
	resp.Code = code

	return nil
}

func (h *CustomerHandler) JoinCustomerToClinic(ctx context.Context, req *pb.JoinCustomerFromClinicRequest, resp *pb.EditCustomerProfileOnSettingPageResponse) error {
	var roles []string
	var updatedBy string = "Unknown"

	// Prefer req.Roles if it's provided
	if len(req.Roles) > 0 {
		roles = req.Roles
	} else if req.Role != "" {
		roles = []string{req.Role}
	}
	status, errorLog, code := h.CustomerService.JoinCustomerToClinic(req.CustomerId, req.ClinicId, updatedBy, roles, ctx)

	resp.UpdateStatus = status
	resp.ErrorLog = errorLog
	resp.Code = code
	return nil
}

func (h *CustomerHandler) GetCustomerByIDs(ctx context.Context, req *pb.GetCustomerByIDsRequest, resp *pb.GetCustomerByIDsResponse) error {
	return nil
}

// Version 0.7.4
func (h *CustomerHandler) CheckCustomerNPINumber(ctx context.Context, req *pb.CheckCustomerNPINumberRequest, resp *pb.NPICheckResult) error {
	result, err := h.CustomerService.CheckCustomerNPINumber(req.NpiNumber, req.ClinicId, req.CustomerId, req.Roles, ctx)
	if err != nil {
		return err
	}
	err = util.Swap(result, resp)
	if err != nil {
		return err
	}

	return nil
}

func (h *CustomerHandler) GetCustomer2FAContact(ctx context.Context, req *pb.GetCustomer2FAContactRequest, resp *pb.GetCustomer2FAContactResponse) error {
	return nil
}

func (h *CustomerHandler) SignUpCustomer(ctx context.Context, req *pb.CustomerSignUpRequest, resp *pb.SignUpResponse) error {
	return nil
}

func (h *CustomerHandler) SearchClientsByName(ctx context.Context, req *pb.SearchcliensNameRequest, resp *pb.SearchclientsInfoResponse) error {
	return nil
}

func (h *CustomerHandler) ListCustomerAllClinics(ctx context.Context, req *pb.ListCustomerAllClinicsRequest, resp *pb.ListCustomerAllClinicsResponse) error {
	return nil
}

func (h *CustomerHandler) CheckClientAttributes(ctx context.Context, req *pb.CheckClientAttributesRequest, resp *pb.CheckClientAttributesResponse) error {
	return nil
}

func (h *CustomerHandler) NewEditCustomerProfileOnSettingPage(ctx context.Context, req *pb.NewEditCustomerProfileOnSettingPageRequest, resp *pb.NewEditCustomerProfileOnSettingPageResponse) error {
	return nil
}

func (h *CustomerHandler) AddCustomerWithNPINumber(ctx context.Context, req *pb.AddCustomerWithNPINumberRequest, resp *pb.AddCustomerWithNPINumberResponse) error {
	return nil
}

// Version 2
func (h *CustomerHandler) V2_EditCustomerProfileOnSettingPage(ctx context.Context, req *pb.V2_EditCustomerProfileOnSettingPageRequest, resp *pb.V2_EditCustomerProfileOnSettingPageResponse) error {
	return nil
}

// VP-4964 OnboardingQuestionnair Check
func (h *CustomerHandler) CheckCustomerOnboardingQuestionnaireStatus(ctx context.Context, req *pb.CheckCustomerOnboardingQuestionnaireStatusRequest, resp *pb.CheckCustomerOnboardingQuestionnaireStatusResponse) error {
	response, err := h.CustomerService.CheckCustomerOnboardingQuestionnaireStatus(req.CustomerId, ctx)
	if err != nil {
		return err
	}
	err = util.Swap(response, resp)
	if err != nil {
		return err
	}

	return nil
}

func (h *CustomerHandler) UpdateCustomerOnboardingQuestionnaireStatus(ctx context.Context, req *pb.UpdateCustomerOnboardingQuestionnaireStatusRequest, resp *pb.UpdateCustomerOnboardingQuestionnaireStatusResponse) error {
	customerID, status := h.CustomerService.UpdateCustomerOnboardingQuestionnaireStatus(req.CustomerId, ctx)
	resp.CustomerId = customerID
	resp.Status = status
	return nil
}

func (h *CustomerHandler) AddCustomerWithNPINumberNative(ctx context.Context, req *pb.AddCustomerWithNPINumberRequest, resp *pb.AddCustomerWithNPINumberResponse) error {

	data := &model.AddCustomerWithNPINumberRequest{
		CustomerFirstName:         req.CustomerFirstName,
		CustomerLastName:          req.CustomerLastName,
		CustomerNPINumber:         req.CustomerNpiNumber,
		CustomerLoginEmail:        req.CustomerLoginEmail,
		CustomerNotificationEmail: req.CustomerNotificationEmail,
		CustomerPhone:             req.CustomerPhone,
		CustomerAddressLine1:      req.CustomerAddressLine_1,
		CustomerAddressLine2:      req.CustomerAddressLine_2,
		CustomerCity:              req.CustomerCity,
		CustomerState:             req.CustomerState,
		CustomerZipcode:           req.CustomerZipcode,
		CustomerCountry:           req.CustomerCountry,
		ClinicID:                  req.ClinicId,
		// InvitedFromCustomer:       req.InvitedFromCustomer,
		CustomerInvitationLink: req.CustomerInvitationLink,
		CustomerSuffix:         req.CustomerSuffix,
		CustomerRoles:          req.CustomerRoles,
	}

	response, err := h.CustomerService.AddCustomerWithNPINumberNative(data, ctx)
	if err != nil {
		return err
	}
	err = util.Swap(response, resp)
	if err != nil {
		return err
	}
	return nil
}

func (h *CustomerHandler) SignUpCustomerV2(ctx context.Context, req *pb.CustomerSignUpRequest, resp *pb.SignUpResponse) error {
	return nil
}

func (h *CustomerHandler) FuzzySearchCustomers(ctx context.Context, req *pb.FuzzySearchCustomersRequest, resp *pb.SearchclientsInfoResponse) error {
	clients, err := h.CustomerService.FuzzySearchCustomers(req.CustomerSearchInput, &req.ClinicId, ctx)
	if err != nil {
		return err
	}
	err = util.Swap(clients, &resp.Clients)
	if err != nil {
		return err
	}
	return nil
}

func (h *CustomerHandler) GetCustomerByNPINumber(ctx context.Context, req *pb.NPINumber, resp *pb.GetCustomerByNPINumberResponse) error {
	customers, err := h.CustomerService.GetCustomerByNPINumber(req.NpiNumber, ctx)
	if err != nil {
		return err
	}
	var result []*pb.GetCustomerByNPINumberCustomerIDResponse
	for _, cust := range customers {
		result = append(result, &pb.GetCustomerByNPINumberCustomerIDResponse{
			CustomerId: int32(cust.ID),
		})
	}

	err = util.Swap(result, &resp.Result)
	if err != nil {
		return err
	}

	return nil
}

func (h *CustomerHandler) ReinviteNPICheck(ctx context.Context, req *pb.ReinviteNPICheckRequest, resp *pb.ReinviteNPICheckResponse) error {
	status, err := h.CustomerService.ReinviteNPICheck(req.CustomerNpiNumber, req.CustomerRoles, ctx)
	resp.Status = status

	if err != nil {
		resp.ErrorMessage = err.Error()
	}
	return nil
}

func (h *CustomerHandler) FuzzySearchCustomerClinicName(ctx context.Context, req *pb.FuzzySearchRequest, resp *pb.FuzzySearchResponse) error {
	customerClinic, err := h.CustomerService.FuzzySearchCustomerClinicName(req.CustomerSearchInput, ctx)
	if err != nil {
		return err
	}

	err = util.Swap(customerClinic, &resp.Results)
	if err != nil {
		return err
	}
	return nil
}

func (h *CustomerHandler) GetStatementData(ctx context.Context, req *pb.GetStatementRequest, resp *pb.GetStatementResponse) error {
	return nil
}

func (h *CustomerHandler) FetchCustomerBetaProgramsForClinic(ctx context.Context, req *pb.FetchCustomerBetaProgramsForClinicInput, resp *pb.FetchCustomerBetaProgramsForClinicResponse) error {
	customerBetaPrograms, errMessage := h.CustomerService.FetchCustomerBetaProgramsForClinic(req.CustomerId, req.ClinicId, ctx)
	resp.ErrorMessage = errMessage
	err := util.Swap(customerBetaPrograms, &resp.Result)
	if err != nil {
		return err
	}

	return nil
}
