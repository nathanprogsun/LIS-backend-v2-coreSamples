package service

import (
	"context"
	"coresamples/common"
	"coresamples/dbutils"
	"coresamples/ent"
	"coresamples/ent/betaprogramparticipation"
	"coresamples/ent/clinic"
	"coresamples/ent/customer"
	"coresamples/ent/internaluser"
	"coresamples/model"
	pb "coresamples/proto"
	"coresamples/util"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ICustomerService interface {
	GetCustomer(customerId int, ctx context.Context) (customer *model.FullCustomer, err error)
	ListCustomer(page string, perPage string, tx context.Context) (customers []*model.FullCustomer, hasNextPage bool, total int32, err error)
	GetSalesCustomer(salesName []string, page string, perPage string, ctx context.Context) (customers []*model.FullCustomer, err error)
	GetCustomerSales(customerNames []string, customerIds []string, ctx context.Context) (customerSales []*model.CustomerSales, err error)
	CheckCustomerNPINumber(npiNumber string, clinicId *string, customerId *string, roles []string, ctx context.Context) (result *pb.NPICheckResult, err error)
	GetCustomerByNPINumber(npiNumber string, ctx context.Context) (customers []*ent.Customer, err error)
	ReinviteNPICheck(customerNpiNumber string, customerRoles []string, ctx context.Context) (status string, err error)
	FuzzySearchCustomers(customerSearchInput string, clinicId *string, ctx context.Context) ([]*model.FuzzyClientObject, error)
	FuzzySearchCustomerClinicName(searchInput string, ctx context.Context) ([]*model.CustomerClinicData, error)
	UpdateCustomerOnboardingQuestionnaireStatus(customerId string, ctx context.Context) (customerID int32, status string)
	CheckCustomerOnboardingQuestionnaireStatus(customerId string, ctx context.Context) (*pb.CheckCustomerOnboardingQuestionnaireStatusResponse, error)
	AddCustomerWithNPINumberNative(data *model.AddCustomerWithNPINumberRequest, ctx context.Context) (*model.AddCustomerWithNPINumberResponse, error)
	JoinCustomerToClinic(customerId string, clinic_id string, updatedBy string, roles []string, ctx context.Context) (updateStatus string, errorLog string, code int32)
	RemoveCustomerFromClinic(customerId string, clinicId string, updatedBy string, ctx context.Context) (updateStatus string, errorLog string, code int32)
	FetchCustomerBetaProgramsForClinic(customerID int32, clinicID int32, ctx context.Context) (CustomerBetaPrograms []*model.CustomerBetaPrograms, errorMessage string)

	// internal use
}

type CustomerService struct {
	Service
	rbacService IRBACService
}

func NewCustomerService(dbClient *ent.Client, redisClient *common.RedisClient) ICustomerService {
	s := &CustomerService{
		Service:     InitService(dbClient, redisClient),
		rbacService: GetCurrentRBACService(),
	}
	return s
}

// CheckCustomerNPINumber implements ICustomerService.
func (c *CustomerService) CheckCustomerNPINumber(npiNumber string, clinicId *string, customerId *string, roles []string, ctx context.Context) (*pb.NPICheckResult, error) {
	if npiNumber == "Internal special NPI" || (len(roles) == 1 && roles[0] == "clinicadmin") {
		return &pb.NPICheckResult{NPI_Check: "Valid NPI"}, nil
	}

	npiResp, err := util.NPIOnlineCheck(npiNumber, ctx)
	if err != nil {
		return nil, err
	}

	if isSuccess := util.IsSuccessResponse(npiResp); !isSuccess {
		return &pb.NPICheckResult{NPI_Check: "Invalid NPI"}, nil
	}

	// NPI valid, now check for duplicates in the clinic
	if clinicId != nil && customerId != nil {
		clinicIDInt, err := strconv.Atoi(*clinicId)
		if err != nil {
			return nil, fmt.Errorf("invalid clinic ID: %w", err)
		}

		customerIDInt, err := strconv.Atoi(*customerId)
		if err != nil {
			return nil, fmt.Errorf("invalid customer ID: %w", err)
		}
		clinicEnt, err := c.dbClient.Clinic.
			Query().
			Where(clinic.IDEQ(clinicIDInt)).
			WithCustomers().
			Only(ctx)
		if err != nil {
			return nil, err
		}

		for _, cust := range clinicEnt.Edges.Customers {
			if cust.CustomerNpiNumber == npiNumber && cust.ID != customerIDInt {
				if npiNumber != "Internal special NPI" && npiNumber != "A000000000" {
					return &pb.NPICheckResult{NPI_Check: "Invalid NPI Used in Clinic Already"}, nil
				}
			}
		}
	}

	return &pb.NPICheckResult{NPI_Check: "Valid NPI"}, nil

}

// GetCustomer implements ICustomerService.
func (c *CustomerService) GetCustomer(customerId int, ctx context.Context) (*model.FullCustomer, error) {
	trackingID := uuid.New().String()

	customers, err := dbutils.FetchModelFullCustomers(ctx, c.dbClient, &customerId, 0, 100)
	if err != nil {
		common.ErrorLogger("[tracking:%s] failed to get customer: %v", trackingID, err)
		return nil, err
	}

	return customers[0], err

}

// GetCustomerByNPINumber implements ICustomerService.
func (c *CustomerService) GetCustomerByNPINumber(npiNumber string, ctx context.Context) ([]*ent.Customer, error) {
	customers, err := c.dbClient.Customer.
		Query().
		Where(customer.CustomerNpiNumberEQ(npiNumber)).
		Select(customer.FieldID).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

// GetCustomerSales implements ICustomerService.
func (c *CustomerService) GetCustomerSales(customerNames []string, customerIds []string, ctx context.Context) (sales []*model.CustomerSales, err error) {
	var allSales []*model.CustomerSales

	// Step 1: Handle customer names
	for _, customerName := range customerNames {
		nameParts := util.SplitName(customerName)

		redisKey := dbutils.KeyGetCustomerSalesByCustomerName(customerName)
		redisResult, err := c.redisClient.Get(ctx, redisKey).Result()
		if err != nil {
			common.Error(err)
		}
		if err == nil && redisResult != "" {
			var cachedSales []*model.CustomerSales
			if err := json.Unmarshal([]byte(redisResult), &cachedSales); err == nil {
				allSales = append(allSales, cachedSales...)
				continue
			}
		}

		// Cache miss, query DB
		var customerSales []*model.CustomerSales
		results, err := c.dbClient.Customer.
			Query().
			Where(
				customer.CustomerFirstNameEQ(nameParts.FirstName),
				customer.CustomerLastNameEQ(nameParts.LastName),
			).
			WithSales().
			All(ctx)
		if err != nil {
			common.Error(err)
			continue
		}

		for _, cust := range results {
			customerSales = append(customerSales, toModelCustomerSales(cust))
		}

		// Cache result
		recb, err := json.Marshal(customerSales)
		if err != nil {
			common.Error(err)
			continue
		}
		c.redisClient.SetEX(ctx, redisKey, recb, time.Second*1000)

		// add this customeSales to allCustomerSales
		allSales = append(allSales, customerSales...)
	}

	// Step 2: Handle customer IDs
	for _, customerIDStr := range customerIds {
		customerID, err := strconv.Atoi(customerIDStr)
		if err != nil {
			common.Error(err)
			continue
		}
		redisKey := dbutils.KeyGetCustomerSalesByCustomerId(customerID)
		redisResult, err := c.redisClient.Get(ctx, redisKey).Result()
		if err != nil {
			common.Error(err)
		}
		if err == nil && redisResult != "" {
			var cachedSales []*model.CustomerSales
			if err := json.Unmarshal([]byte(redisResult), &cachedSales); err == nil {
				allSales = append(allSales, cachedSales...)
				continue
			}
		}

		var customerSales []*model.CustomerSales
		results, err := c.dbClient.Customer.
			Query().
			Where(customer.IDEQ(customerID)).
			WithSales().
			All(ctx)
		if err != nil {
			common.Error(err)
			continue
		}

		for _, cust := range results {
			customerSales = append(customerSales, toModelCustomerSales(cust))
		}

		// Cache result
		recb, err := json.Marshal(customerSales)
		if err != nil {
			common.Error(err)
			continue
		}
		c.redisClient.SetEX(ctx, redisKey, recb, time.Second*1000)

		// add this customeSales to allCustomerSales
		allSales = append(allSales, customerSales...)
	}

	return allSales, nil
}

// GetSalesCustomer implements ICustomerService.
func (c *CustomerService) GetSalesCustomer(salesName []string, page string, perPage string, ctx context.Context) ([]*model.FullCustomer, error) {
	var arrayResult []*model.FullCustomer

	pageNum, err := strconv.Atoi(page)
	if err != nil {
		return nil, fmt.Errorf("invalid page: %w", err)
	}

	pageSizeNum, err := strconv.Atoi(perPage)
	if err != nil {
		return nil, fmt.Errorf("invalid perPage: %w", err)
	}

	for _, sales_name := range salesName {
		redisKey := dbutils.KeyGetSalesCustomer(sales_name, page, perPage)
		cachedData, err := c.redisClient.Get(ctx, redisKey).Result()
		if err != nil {
			common.Error(err)
		}
		if err == nil && cachedData != "" {
			var cachedCustomers []*model.FullCustomer
			if err := json.Unmarshal([]byte(cachedData), &cachedCustomers); err == nil {
				arrayResult = append(arrayResult, cachedCustomers...)
				continue
			}
		}

		var currentResult []*model.FullCustomer

		internalUser, err := c.dbClient.InternalUser.
			Query().
			Where(
				internaluser.InternalUserNameEQ(sales_name),
				internaluser.IsActive(true),
			).
			WithCustomers(func(q *ent.CustomerQuery) {
				q.
					WithClinics().
					Offset((pageNum - 1) * pageSizeNum).
					Limit(pageSizeNum)
			}).
			Only(ctx)
		if err != nil {
			common.Error(err)
			continue
		}

		for _, cust := range internalUser.Edges.Customers {
			var clinics []*model.CustomerClinic
			for _, clinic := range cust.Edges.Clinics {
				addresses := dbutils.FetchCustomerClinicAddresses(ctx, c.dbClient, cust.ID, clinic.ID)
				contacts := dbutils.FetchCustomerClinicContacts(ctx, c.dbClient, cust.ID, clinic.ID)
				clinics = append(clinics, &model.CustomerClinic{
					ClinicID:          int32(clinic.ID),
					ClinicName:        clinic.ClinicName,
					UserID:            int32(clinic.UserID),
					IsActive:          clinic.IsActive,
					ClinicAccountID:   int32(clinic.ClinicAccountID),
					CustomerAddresses: addresses,
					CustomerContacts:  contacts,
				})
			}
			currentResult = append(currentResult, &model.FullCustomer{
				CustomerID:                int32(cust.ID),
				UserID:                    int32(cust.UserID),
				CustomerFirstName:         cust.CustomerFirstName,
				CustomerLastName:          cust.CustomerLastName,
				CustomerMiddleName:        cust.CustomerMiddleName,
				CustomerTypeID:            cust.CustomerTypeID,
				CustomerSuffix:            cust.CustomerSuffix,
				CustomerSamplesReceived:   cust.CustomerSamplesReceived,
				CustomerRequestSubmitTime: cust.CustomerRequestSubmitTime.Format(time.RFC3339),
				IsActive:                  cust.IsActive,
				Clinics:                   clinics,
				CustomerNPINumber:         cust.CustomerNpiNumber,
				SalesID:                   int32(cust.SalesID),
				CustomerSignupTime:        cust.CustomerSignupTime.Format(time.RFC3339),
			})
		}

		// Cache the result
		recb, err := json.Marshal(currentResult)
		if err == nil {
			_ = c.redisClient.SetEX(ctx, redisKey, recb, 10*time.Minute).Err()
		}

		arrayResult = append(arrayResult, currentResult...)
	}

	return arrayResult, nil
}

// ListCustomer implements ICustomerService.
func (c *CustomerService) ListCustomer(page string, perPage string, ctx context.Context) (customers []*model.FullCustomer, hasNextPage bool, total int32, err error) {
	trackingID := uuid.New().String()
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		return nil, false, 0, fmt.Errorf("invalid page: %w", err)
	}

	pageSizeNum, err := strconv.Atoi(perPage)
	if err != nil || pageSizeNum <= 0 {
		return nil, false, 0, fmt.Errorf("invalid perPage: %w", err)
	}

	const redisKey = "lis::core_service_v2::customer::total_customer_count"
	redisResult, redisErr := c.redisClient.Get(ctx, redisKey).Result()
	if redisErr == nil && redisResult != "" {
		if err := json.Unmarshal([]byte(redisResult), &total); err != nil {
			common.ErrorLogger("[tracking:%s] failed to unmarshal total_customer_count from redis: %v", trackingID, err)
			return nil, false, 0, err
		}
	} else {
		count, err := c.dbClient.Customer.Query().Count(ctx)
		if err != nil {
			common.ErrorLogger("[tracking:%s] failed to count customers: %v", trackingID, err)
			return nil, false, 0, err
		}
		total = int32(count)

		encoded, _ := json.Marshal(total)
		_ = c.redisClient.SetEX(ctx, redisKey, encoded, 8*time.Hour).Err()
	}

	taken := pageNum * pageSizeNum
	hasNextPage = taken < int(total)

	customers, err = dbutils.FetchModelFullCustomers(ctx, c.dbClient, nil, (pageNum-1)*pageSizeNum, pageSizeNum)
	if err != nil {
		common.ErrorLogger("[tracking:%s] failed to list customers: %v", trackingID, err)
		return nil, hasNextPage, total, err
	}

	return customers, hasNextPage, total, nil
}

// ReinviteNPICheck implements ICustomerService.
func (c *CustomerService) ReinviteNPICheck(customerNpiNumber string, customerRoles []string, ctx context.Context) (status string, err error) {
	if customerNpiNumber == "Internal special NPI" {
		return "success", nil
	}
	var npiRole string
	if customerNpiNumber != "" {
		npiResp, err := util.NPIOnlineCheck(customerNpiNumber, ctx)
		if err != nil {
			return "failed", fmt.Errorf("NPI lookup failed: %w", err)
		}

		if !util.IsSuccessResponse(npiResp) {
			return "failed", fmt.Errorf("invalid NPI")
		}

		enumType := npiResp.Results[0].EnumerationType
		switch enumType {
		case "NPI-1":
			npiRole = "clinicadmin"
		case "NPI-2":
			npiRole = "clinicadminonly"
		default:
			return "failed", fmt.Errorf("unsupported npi type: %s", enumType)
		}
	}
	if npiRole == "clinicadminonly" {
		for _, role := range customerRoles {
			if strings.ToLower(role) == "provider" {
				return "failed", fmt.Errorf("NPI type does not meet requirements for role: provider")
			}
		}
	}

	return "success", nil
}

// FuzzySearchCustomers implements ICustomerService.
func (c *CustomerService) FuzzySearchCustomers(customerSearchInput string, clinicId *string, ctx context.Context) ([]*model.FuzzyClientObject, error) {
	var customers []*ent.Customer
	var err error

	// Fetch customers based on clinicId
	if clinicId != nil {
		clinicIDInt, err := strconv.Atoi(*clinicId)
		if err != nil {
			return nil, fmt.Errorf("invalid clinic ID: %w", err)
		}

		customers, err = c.dbClient.Customer.
			Query().
			Where(customer.HasClinicsWith(clinic.IDEQ(clinicIDInt))).
			All(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch customers for clinic ID %d: %w", clinicIDInt, err)
		}
	} else {
		customers, err = c.dbClient.Customer.
			Query().
			All(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch customers: %w", err)
		}
	}

	if len(customers) == 0 {
		return []*model.FuzzyClientObject{}, nil
	}

	searchInput := strings.TrimSpace(customerSearchInput)
	if searchInput == "" {
		return nil, fmt.Errorf("search input cannot be empty")
	}
	isNumeric := util.IsNumericString(searchInput)
	var matchedCustomers []*model.FuzzyClientObject
	if isNumeric {
		// Numeric search: match customer IDs containing the search input
		for _, cust := range customers {
			if strings.Contains(fmt.Sprintf("%d", cust.ID), searchInput) {
				matchedCustomers = append(matchedCustomers, &model.FuzzyClientObject{
					ClientId:   int64(cust.ID),
					ClientName: util.AssembleFullName(cust.CustomerFirstName, cust.CustomerMiddleName, cust.CustomerLastName),
				})
			}
		}
	} else {
		inputParts := strings.Fields(searchInput)
		for _, cust := range customers {
			fn := strings.ToLower(cust.CustomerFirstName)
			ln := strings.ToLower(cust.CustomerLastName)
			mn := strings.ToLower(cust.CustomerMiddleName)
			switch len(inputParts) {
			case 3:
				if strings.Contains(fn, strings.ToLower(inputParts[0])) &&
					strings.Contains(mn, strings.ToLower(inputParts[1])) &&
					strings.Contains(ln, strings.ToLower(inputParts[2])) {
					matchedCustomers = append(matchedCustomers, &model.FuzzyClientObject{
						ClientId:   int64(cust.ID),
						ClientName: util.AssembleFullName(cust.CustomerFirstName, cust.CustomerMiddleName, cust.CustomerLastName),
					})
				}
			case 2:
				if strings.Contains(fn, strings.ToLower(inputParts[0])) &&
					strings.Contains(ln, strings.ToLower(inputParts[1])) {
					matchedCustomers = append(matchedCustomers, &model.FuzzyClientObject{
						ClientId:   int64(cust.ID),
						ClientName: util.AssembleFullName(cust.CustomerFirstName, cust.CustomerMiddleName, cust.CustomerLastName),
					})
				}
			case 1:
				term := strings.ToLower(inputParts[0])
				if strings.Contains(fn, term) || strings.Contains(mn, term) || strings.Contains(ln, term) {
					matchedCustomers = append(matchedCustomers, &model.FuzzyClientObject{
						ClientId:   int64(cust.ID),
						ClientName: util.AssembleFullName(cust.CustomerFirstName, cust.CustomerMiddleName, cust.CustomerLastName),
					})
				}
			}
		}
	}

	return matchedCustomers, nil
}

func (c *CustomerService) FuzzySearchCustomerClinicName(searchInput string, ctx context.Context) ([]*model.CustomerClinicData, error) {
	trackingID := uuid.New().String()
	const redisKey = "customer_clinic_data"

	var data []*model.CustomerClinicData
	var err error

	redisResult, redisErr := c.redisClient.Get(ctx, redisKey).Result()
	if redisErr == nil && redisResult != "" {
		if err := json.Unmarshal([]byte(redisResult), &data); err != nil {
			common.ErrorLogger("[tracking:%s] failed to unmarshal customer_clinic_data from redis: %v", trackingID, err)
			return nil, err
		}
	} else {
		data, err = dbutils.FetchAndCacheCustomerClinicData(ctx, c.dbClient, c.redisClient)
		if err != nil {
			common.ErrorLogger("[tracking:%s] Fetch and cache error: %v", trackingID, err)
			return nil, fmt.Errorf("fetch and cache error: %w", err)
		}
	}

	searchTerm := strings.ToLower(strings.TrimSpace(searchInput))
	var results []*model.CustomerClinicData

	for _, entry := range data {
		if strings.Contains(strings.ToLower(fmt.Sprint(entry.CustomerId)), searchTerm) ||
			strings.Contains(strings.ToLower(entry.CustomerName), searchTerm) ||
			strings.Contains(strings.ToLower(entry.ClinicName), searchTerm) {
			results = append(results, &model.CustomerClinicData{
				CustomerId:   entry.CustomerId,
				CustomerName: entry.CustomerName,
				ClinicName:   entry.ClinicName,
			})
		}
	}
	return results, nil
}

// UpdateCustomerOnboardingQuestionnaireStatus implements ICustomerService.
func (c *CustomerService) UpdateCustomerOnboardingQuestionnaireStatus(customerId string, ctx context.Context) (customerID int32, status string) {
	// Convert customerId from string to integer
	custIDInt, err := strconv.Atoi(customerId)
	if err != nil {
		// Log and return if the customer ID is not a valid number
		common.ErrorLogger("invalid customer_id: %v", err)
		return 0, "InvalidCustomerID"
	}

	// Query the customer record by ID
	customer, err := c.dbClient.Customer.
		Query().
		Where(customer.IDEQ(custIDInt)).
		Only(ctx)
	if err != nil {
		// If customer not found, return a specific status
		if ent.IsNotFound(err) {
			return 0, "CustomerNotFound"
		}
		// Log other database errors
		common.ErrorLogger("failed to query customer: %v", err)
		return 0, "InternalError"
	}

	// Update the "onboarding_questionnaire_filled_on" field with the current time
	_, err = c.dbClient.Customer.
		UpdateOne(customer).
		SetOnboardingQuestionnaireFilledOn(time.Now()).
		Save(ctx)
	if err != nil {
		// Log and return if update fails
		common.ErrorLogger("failed to update onboarding questionnaire: %v", err)
		return 0, "InternalError"
	}

	// Return success with the customer ID and status
	return int32(customer.ID), "OnboardingQuestionnaireUpdated"
}

// CheckCustomerOnboardingQuestionnaireStatus implements ICustomerService.
func (c *CustomerService) CheckCustomerOnboardingQuestionnaireStatus(customerId string, ctx context.Context) (*pb.CheckCustomerOnboardingQuestionnaireStatusResponse, error) {
	response := &pb.CheckCustomerOnboardingQuestionnaireStatusResponse{}

	custIDInt, err := strconv.Atoi(customerId)
	if err != nil {
		common.ErrorLogger("invalid customer_id: %v", err)
		return nil, fmt.Errorf("invalid customer_id: %w", err)
	}

	// Fetch the customer's basic info (only ID and questionnaire timestamp)
	customer, err := c.dbClient.Customer.
		Query().
		Where(customer.IDEQ(custIDInt)).
		Select(customer.FieldID, customer.FieldOnboardingQuestionnaireFilledOn).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			// If customer doesn't exist, return not found error
			return nil, fmt.Errorf("customer not found: %d", custIDInt)
		}
		common.ErrorLogger("failed to query customer: %v", err)
		return nil, fmt.Errorf("failed to query customer: %w", err)
	}

	// Populate response fields
	response.CustomerId = int32(customer.ID)
	if customer.OnboardingQuestionnaireFilledOn.IsZero() {
		response.IsOnboardingQuestionnaireFilled = false
	} else {
		response.IsOnboardingQuestionnaireFilled = true
		response.OnboardingQuestionnaireFilledOn = customer.OnboardingQuestionnaireFilledOn.Format(time.RFC3339)
	}

	return response, nil
}

// AddCustomerWithNPINumberNative implements ICustomerService.
func (c *CustomerService) AddCustomerWithNPINumberNative(data *model.AddCustomerWithNPINumberRequest, ctx context.Context) (*model.AddCustomerWithNPINumberResponse, error) {

	var npiRole string
	npiResp, err := util.NPIOnlineCheck(data.CustomerNPINumber, ctx)
	if err != nil || !util.IsSuccessResponse(npiResp) {
		return &model.AddCustomerWithNPINumberResponse{
			Status:       "fail",
			CustomerID:   -1,
			ErrorMessage: "Invalid NPI or lookup fail",
		}, err
	}

	enumType := npiResp.Results[0].EnumerationType
	switch enumType {
	case "NPI-1":
		npiRole = "clinicadmin"
	case "NPI-2":
		npiRole = "clinicadminonly"
	default:
		return &model.AddCustomerWithNPINumberResponse{
			Status:       "fail",
			CustomerID:   -1,
			ErrorMessage: "unsupported npi type",
		}, nil
	}

	if npiRole == "clinicadminonly" && util.Contains(data.CustomerRoles, "provider") {
		return &model.AddCustomerWithNPINumberResponse{
			Status:       "fail",
			CustomerID:   -1,
			ErrorMessage: "NPI type does not meet requirements.",
		}, nil
	}

	newCustomerID, err := dbutils.GetNewCustomerID(ctx, c.dbClient)
	clinicID, err := strconv.Atoi(data.ClinicID)
	if err != nil {
		return nil, err
	}
	createdCust, err := c.dbClient.Customer.
		Create().
		SetID(newCustomerID).
		SetCustomerFirstName(data.CustomerFirstName).
		SetCustomerLastName(data.CustomerLastName).
		SetCustomerNpiNumber(data.CustomerNPINumber).
		SetCustomerSuffix(data.CustomerSuffix).
		AddClinicIDs(clinicID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	//address

	addressData := &model.Address{
		AddressType:      "office",
		StreetAddress:    data.CustomerAddressLine1,
		AptPO:            data.CustomerAddressLine2,
		City:             data.CustomerCity,
		State:            data.CustomerState,
		Zipcode:          data.CustomerZipcode,
		Country:          data.CustomerCountry,
		AddressConfirmed: true,
		IsPrimaryAddress: true,
		AddressLevel:     1,
		AddressLevelName: "Customer",
	}
	createdAddressID, err := dbutils.CreateAddress(ctx, c.dbClient, addressData)
	if err != nil {
		return nil, err
	}

	err = dbutils.AddAddressToCustomerClinic(ctx, c.dbClient, &model.CustomerAddressOnClinicsCreation{
		CustomerID:  int32(createdCust.ID),
		ClinicID:    int32(clinicID),
		AddressID:   createdAddressID,
		AddressType: addressData.AddressType,
	})
	if err != nil {
		return nil, err
	}

	//contact
	emailData := &model.Contact{
		ContactDescription: "notification email",
		ContactDetails:     data.CustomerNotificationEmail,
		ContactType:        "email",
		ContactLevel:       1,
		ContactLevelName:   "Customer",
	}
	createdEmailContactID, err := dbutils.CreateContact(ctx, c.dbClient, emailData)
	if err != nil {
		return nil, err
	}

	err = dbutils.AddContactToCustomerClinic(ctx, c.dbClient, &model.CustomerContactOnClinicsCreation{
		CustomerID:  int32(createdCust.ID),
		ClinicID:    int32(clinicID),
		ContactID:   createdEmailContactID,
		ContactType: emailData.ContactType,
	})
	if err != nil {
		return nil, err
	}

	phoneData := &model.Contact{
		ContactDescription: "contact phone",
		ContactDetails:     data.CustomerPhone,
		ContactType:        "phone",
		ContactLevel:       1,
		ContactLevelName:   "Customer",
	}
	createdPhoneContactID, err := dbutils.CreateContact(ctx, c.dbClient, phoneData)
	if err != nil {
		return nil, err
	}

	err = dbutils.AddContactToCustomerClinic(ctx, c.dbClient, &model.CustomerContactOnClinicsCreation{
		CustomerID:  int32(createdCust.ID),
		ClinicID:    int32(clinicID),
		ContactID:   createdPhoneContactID,
		ContactType: phoneData.ContactType,
	})
	if err != nil {
		return nil, err
	}

	//role: for rbac customer must have account (user_id) to be assigned a role. Will add roles during the sign up process

	//invitation link
	if data.CustomerInvitationLink != "" {
		_, _ = c.dbClient.UserInvitationRecord.
			Create().
			SetCustomerID(createdCust.ID).
			SetInvitationLink(data.CustomerInvitationLink).
			Save(ctx)
	}

	return &model.AddCustomerWithNPINumberResponse{
		Status:       "success",
		CustomerID:   int32(createdCust.ID),
		ErrorMessage: "",
	}, nil
}

// JoinCustomerToClinic implements ICustomerService.
func (c *CustomerService) JoinCustomerToClinic(customerId string, clinicId string, updatedBy string, roles []string, ctx context.Context) (updateStatus string, errorLog string, code int32) {

	custID, err := strconv.Atoi(customerId)
	if err != nil {
		common.ErrorLogger("invalid customer_id: %v", err)
		return "failed", "invalid customer_id", 400
	}

	clinicID, err := strconv.Atoi(clinicId)
	if err != nil {
		common.ErrorLogger("invalid clinic_id: %v", err)
		return "failed", "invalid clinic_id", 400
	}

	clinicEnt, err := c.dbClient.Clinic.
		Query().
		Where(clinic.IDEQ(clinicID)).
		WithCustomers(func(q *ent.CustomerQuery) {
			q.Select(customer.FieldID, customer.FieldCustomerNpiNumber, customer.FieldUserID, customer.FieldIsActive)
		}).
		Only(ctx)
	if err != nil {
		common.ErrorLogger("failed to query clinic: %v", err)
		return "failed", "clinic not found", 404
	}

	// Check if customer already in clinic
	for _, cust := range clinicEnt.Edges.Customers {
		if cust.ID == custID {

			c.redisClient.Del(ctx, dbutils.KeyGetCustomerAllClinics(customerId))
			return "success", "Customer is Already in the Clinic", 200
		}
	}

	customerEnt, err := c.dbClient.Customer.Get(ctx, custID)
	if err != nil {
		common.ErrorLogger("customer not found: %v", err)
		return "failed", "customer not found", 404
	}

	for _, cust := range clinicEnt.Edges.Customers {
		if cust.IsActive &&
			cust.CustomerNpiNumber == customerEnt.CustomerNpiNumber &&
			customerEnt.CustomerNpiNumber != "Internal special NPI" &&
			customerEnt.CustomerNpiNumber != "A000000000" {
			fmt.Print("bb")
			return "failed", "NPI Number Duplicate", 500
		}
	}

	_, err = c.dbClient.Customer.
		UpdateOneID(custID).
		AddClinicIDs(clinicID).
		Save(ctx)
	if err != nil {
		common.ErrorLogger("failed to connect customer to clinic: %v", err)
		return "failed", "failed to link customer to clinic", 500
	}

	snapshot := map[string]string{
		"customer_id": customerId,
		"clinic_id":   clinicId,
	}

	auditLogMessage := common.AuditLogEntry{
		EventID:        uuid.NewString(),
		ServiceName:    common.ServiceName,
		ServiceType:    "backend",
		EventName:      "joinCustomerToClinic",
		EntityType:     "clinic_id",
		EntityID:       clinicId,
		User:           updatedBy,
		Entrypoint:     "GRPC",
		EntitySnapshot: util.MustMarshalJSON(snapshot),
	}
	go func() {
		common.RecordAuditLog(auditLogMessage)
	}()

	c.redisClient.Del(ctx, dbutils.KeyGetCustomerAllClinics(customerId))

	//role
	for _, role := range roles {
		switch strings.ToLower(role) {
		case "provider":
			npiResp, err := util.NPIOnlineCheck(customerEnt.CustomerNpiNumber, ctx)
			if err != nil || !util.IsSuccessResponse(npiResp) {
				common.Error(err)
				continue
			}

			if npiResp.Results[0].EnumerationType == "NPI-1" {
				err = c.rbacService.AssignRoleToAccountInDomain(role, int32(customerEnt.UserID), "clinic", int32(clinicID), 0, ctx)
				if err != nil {
					common.ErrorLogger("failed to add provider role: %v", err)
					return "partial success", "customer join clinic but failed to add provider role", 200
				}
			} else {
				return "partial success", "customer join clinic but non NPI-1 type can't have provider role", 200
			}
		case "clinicadmin":
			err = c.rbacService.AssignRoleToAccountInDomain(role, int32(customerEnt.UserID), "clinic", int32(clinicID), 0, ctx)
			if err != nil {
				common.ErrorLogger("failed to add clinicadmin role: %v", err)
				return "partial success", "customer join clinic but failed to add clinicadmin role", 200
			}
		default:
			common.ErrorLogger("unrecognized role: %s", role)
		}
	}

	return "success", "", 200
}

// RemoveCustomerFromClinic implements ICustomerService.
func (c *CustomerService) RemoveCustomerFromClinic(customerId string, clinicId string, updatedBy string, ctx context.Context) (updateStatus string, errorLog string, code int32) {

	// Parse IDs
	custIDInt, err := strconv.Atoi(customerId)
	if err != nil {
		return "failed", "invalid customer_id", 400
	}
	clinicIDInt, err := strconv.Atoi(clinicId)
	if err != nil {
		return "failed", "invalid clinic_id", 400
	}
	clinicEnt, err := c.dbClient.Clinic.
		Query().
		Where(clinic.IDEQ(clinicIDInt)).
		WithCustomers(
			func(q *ent.CustomerQuery) {
				q.Select(customer.FieldID, customer.FieldUserID)
			},
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return "failed", "clinic not found", 404
		}
		common.ErrorLogger("query clinic failed: %v", err)
		return "failed", "internal error", 500
	}
	var foundCustomer *ent.Customer
	for _, cust := range clinicEnt.Edges.Customers {
		if cust.ID == custIDInt {
			fmt.Print(cust.ID)
			foundCustomer = cust
			break
		}
	}
	if foundCustomer == nil {
		return "failed", "no_such_relation", 404
	}

	if foundCustomer.UserID == clinicEnt.UserID {
		return "failed", "Ask Support to Remove the Founder of the Clinic", 500
	}

	snapshot := map[string]string{
		"customer_id": customerId,
		"clinic_id":   clinicId,
	}

	auditLogMessage := common.AuditLogEntry{
		EventID:        uuid.NewString(),
		ServiceName:    common.ServiceName,
		ServiceType:    "backend",
		EventName:      "removeCustomerfromClinic",
		EntityType:     "clinic_id",
		EntityID:       clinicId,
		User:           updatedBy,
		Entrypoint:     "GRPC",
		EntitySnapshot: util.MustMarshalJSON(snapshot),
	}
	go func() {
		common.RecordAuditLog(auditLogMessage)
	}()

	_, err = c.dbClient.Customer.
		UpdateOneID(custIDInt).
		RemoveClinics(clinicEnt).
		Save(ctx)
	if err != nil {
		common.ErrorLogger("failed to disconnect clinic: %v", err)
		return "failed", "internal error", 500
	}

	return "success", "", 200
}

// FetchCustomerBetaProgramsForClinic implements ICustomerService.
func (c *CustomerService) FetchCustomerBetaProgramsForClinic(customerID int32, clinicID int32, ctx context.Context) (CustomerBetaPrograms []*model.CustomerBetaPrograms, errorMessage string) {
	if customerID == 0 && clinicID == 0 {
		return nil, "No Customer ID/Clinic ID Inputted"
	}

	query := c.dbClient.BetaProgramParticipation.Query().WithBetaProgram()
	// Apply filters based on input
	if customerID != 0 && clinicID != 0 {
		query = query.Where(
			betaprogramparticipation.CustomerIDEQ(int(customerID)),
			betaprogramparticipation.ClinicIDEQ(int(clinicID)),
		)
	} else if customerID != 0 {
		query = query.Where(betaprogramparticipation.CustomerIDEQ(int(customerID)))
	} else if clinicID != 0 {
		query = query.Where(betaprogramparticipation.ClinicIDEQ(int(clinicID)))
	}

	// Fetch data
	betaProgramInfo, err := query.All(ctx)
	if err != nil {
		common.ErrorLogger("Failed to fetch beta program participations: %v", err)
		return nil, "Internal Server Error"
	}
	if len(betaProgramInfo) == 0 {
		return []*model.CustomerBetaPrograms{}, ""
	}

	// Grouping by customerID-clinicID
	grouped := make(map[string]*model.CustomerBetaPrograms)
	for _, entry := range betaProgramInfo {
		cID := entry.CustomerID
		clID := entry.ClinicID
		beta := entry.Edges.BetaProgram

		if beta == nil {
			continue
		}

		key := fmt.Sprintf("%d-%d", cID, clID)
		if _, exists := grouped[key]; !exists {
			grouped[key] = &model.CustomerBetaPrograms{
				CustomerID:   int32(cID),
				ClinicID:     int32(clID),
				BetaPrograms: []string{},
			}
		}
		grouped[key].BetaPrograms = append(grouped[key].BetaPrograms, beta.BetaProgramName)
	}
	// Convert map to slice
	var result []*model.CustomerBetaPrograms
	for _, v := range grouped {
		result = append(result, v)
	}

	return result, ""
}

func toSaleDetailWithCustomerV7(cust *ent.Customer) *pb.SaleDetailcWithCustomerV7 {
	var internalUser *pb.SaleDetailcV7
	if cust.Edges.Sales != nil {
		internalUser = &pb.SaleDetailcV7{
			InternalUserRoleId:     int32(cust.Edges.Sales.InternalUserRoleID),
			InternalUserFirstname:  cust.Edges.Sales.InternalUserFirstname,
			InternalUserLastname:   cust.Edges.Sales.InternalUserLastname,
			InternalUserMiddlename: cust.Edges.Sales.InternalUserMiddleName,
			InternalUserEmail:      cust.Edges.Sales.InternalUserEmail,
			InternalUserPhone:      cust.Edges.Sales.InternalUserPhone,
		}
	}

	return &pb.SaleDetailcWithCustomerV7{
		CustomerId:         int32(cust.ID),
		CustomerFirstName:  cust.CustomerFirstName,
		CustomerLastName:   cust.CustomerLastName,
		CustomerMiddleName: cust.CustomerMiddleName,
		InternalUser:       internalUser,
	}
}

func toModelCustomerSales(cust *ent.Customer) *model.CustomerSales {
	var internalUser *model.InternalUser
	if cust.Edges.Sales != nil {
		internalUser = &model.InternalUser{
			InternalUserRoleId:     int32(cust.Edges.Sales.InternalUserRoleID),
			InternalUserFirstname:  cust.Edges.Sales.InternalUserFirstname,
			InternalUserLastname:   cust.Edges.Sales.InternalUserLastname,
			InternalUserMiddlename: cust.Edges.Sales.InternalUserMiddleName,
			InternalUserEmail:      cust.Edges.Sales.InternalUserEmail,
			InternalUserPhone:      cust.Edges.Sales.InternalUserPhone,
		}
	}

	return &model.CustomerSales{
		CustomerId:         int32(cust.ID),
		CustomerFirstName:  cust.CustomerFirstName,
		CustomerLastName:   cust.CustomerLastName,
		CustomerMiddleName: cust.CustomerMiddleName,
		InternalUser:       internalUser,
	}
}

func toFullCustomer(cust *ent.Customer) *pb.FullCustomer {
	return &pb.FullCustomer{
		CustomerId:                int32(cust.ID),
		UserId:                    int32(cust.UserID),
		CustomerFirstName:         cust.CustomerFirstName,
		CustomerLastName:          cust.CustomerLastName,
		CustomerMiddleName:        cust.CustomerMiddleName,
		CustomerTypeId:            cust.CustomerTypeID,
		CustomerSuffix:            cust.CustomerSuffix,
		CustomerSamplesReceived:   cust.CustomerSamplesReceived,
		CustomerRequestSubmitTime: cust.CustomerRequestSubmitTime.Format(time.RFC3339),
		IsActive:                  cust.IsActive,
		Clinics:                   toProtoClinics(cust.Edges.Clinics),
		CustomerNpiNumber:         cust.CustomerNpiNumber,
		SalesId:                   int32(cust.SalesID),
		CustomerSignupTime:        cust.CustomerSignupTime.Format(time.RFC3339),
	}
}

func toProtoAddresses(addresses []*ent.Address) []*pb.Address {
	var result []*pb.Address
	for _, addr := range addresses {
		result = append(result, &pb.Address{
			AddressId:        int32(addr.ID),
			AddressType:      addr.AddressType,
			StreetAddress:    addr.StreetAddress,
			AptPo:            addr.AptPo,
			City:             addr.City,
			State:            addr.State,
			Zipcode:          addr.Zipcode,
			Country:          addr.Country,
			AddressConfirmed: addr.AddressConfirmed,
			IsPrimaryAddress: addr.IsPrimaryAddress,
			CustomerId:       int32(addr.CustomerID),
			PatientId:        int32(addr.PatientID),
			ClinicId:         int32(addr.ClinicID),
			InternalUserId:   int32(addr.InternalUserID),
		})
	}
	return result
}

func toProtoContacts(contacts []*ent.Contact) []*pb.Contact {
	var result []*pb.Contact
	for _, c := range contacts {
		result = append(result, &pb.Contact{
			ContactId:          int32(c.ID),
			ContactDescription: c.ContactDescription,
			ContactDetails:     c.ContactDetails,
			ContactType:        c.ContactType,
			IsPrimaryContact:   c.IsPrimaryContact,
			CustomerId:         int32(c.CustomerID),
			PatientId:          int32(c.PatientID),
			ClinicId:           int32(c.ClinicID),
		})
	}
	return result
}

func toProtoSettings(settings []*ent.Setting) []*pb.Setting {
	var result []*pb.Setting
	for _, s := range settings {
		result = append(result, &pb.Setting{
			SettingId:          int32(s.ID),
			SettingName:        s.SettingName,
			SettingDescription: s.SettingDescription,
			SettingValue:       s.SettingValue,
			SettingType:        s.SettingType,
		})
	}
	return result
}

func toProtoClinics(clinics []*ent.Clinic) []*pb.CustomerClinic {
	var result []*pb.CustomerClinic
	for _, clinic := range clinics {
		result = append(result, &pb.CustomerClinic{
			ClinicId:        int32(clinic.ID),
			ClinicName:      clinic.ClinicName,
			UserId:          int32(clinic.UserID),
			IsActive:        clinic.IsActive,
			ClinicAccountId: int32(clinic.ClinicAccountID),
		})
	}
	return result
}
