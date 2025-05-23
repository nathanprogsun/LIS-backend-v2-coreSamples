package service

import (
	"context"
	"coresamples/common"
	"coresamples/dbutils"
	"coresamples/ent"
	"coresamples/ent/clinic"
	"coresamples/ent/customeraddressonclinics"
	"coresamples/ent/customercontactonclinics"
	"coresamples/ent/enttest"
	"coresamples/ent/userinvitationrecord"
	"coresamples/model"
	"coresamples/publisher"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"coresamples/util"

	"github.com/casbin/casbin/v2"
	casbin_constant "github.com/casbin/casbin/v2/constant"
	entadapter "github.com/casbin/ent-adapter"
	casbin_enttest "github.com/casbin/ent-adapter/ent/enttest"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stvp/tempredis"
)

func setupCustomerTest(t *testing.T) (*CustomerService, *tempredis.Server) {
	dataSource := "file:ent?mode=memory&_fk=1"
	dbClient := enttest.Open(t, "sqlite3", dataSource)

	err := dbClient.Schema.Create(context.Background())
	if err != nil {
		common.Fatalf("failed opening connection to MySQL", err)
	}

	server, err := tempredis.Start(tempredis.Config{
		"port": "0",
	})
	if err != nil {
		common.Fatalf("Failed to start tempredis: %v", err)
	}

	common.InitZapLogger("debug")

	redisClient := redis.NewClient(&redis.Options{
		Network: "unix",
		Addr:    server.Socket(),
	})

	rc := common.NewRedisClient(redisClient, redisClient)

	//
	rs := &RBACService{
		Service: InitService(dbClient, nil),
	}
	c := casbin_enttest.Open(t, "sqlite3", dataSource)
	adapter, err := entadapter.NewAdapterWithClient(c)
	casbinDbClient = c
	if err != nil {
		common.Fatal(err)
	}

	rs.enforcer, err = casbin.NewEnforcer("rbac_model.conf", adapter)
	if err != nil {
		common.Fatal(err)
	}
	rs.enforcer.SetFieldIndex("p", casbin_constant.DomainIndex, 3)
	rs.enforcer.EnableAutoSave(true)
	_, err = setupTestDB(rs)
	if err != nil {
		common.Fatal(err)
	}
	//

	s := &CustomerService{
		Service: InitService(
			dbClient,
			rc),
		rbacService: rs,
	}

	// This will start a mock publisher in the memory
	publisher.InitMockPublisher()
	return s, server
}

func cleanupCustomerTest(svc *CustomerService, server *tempredis.Server) {
	var err error
	if err = server.Kill(); err != nil {
		common.Error(err)
	}
	if svc.dbClient != nil {
		if err = svc.dbClient.Close(); err != nil {
			common.Error(err)
		}
	}
	publisher.GetPublisher().GetWriter().Close()
}

func TestGetCustomer(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)

	ctx := context.Background()
	customer, err := svc.dbClient.Customer.Create().
		SetCustomerFirstName("Test").
		SetCustomerLastName("Customer1").
		SetCustomerTypeID("type_1").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	clinic, err := svc.dbClient.Clinic.Create().
		SetClinicName("Test Clinic").
		AddCustomerIDs(customer.ID).
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	address, err := svc.dbClient.Address.Create().
		SetStreetAddress("123 Main St").
		SetCity("Testville").
		SetState("TS").
		SetZipcode("12345").
		SetIsPrimaryAddress(true).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.dbClient.CustomerAddressOnClinics.Create().
		SetClinicID(clinic.ID).
		SetCustomerID(customer.ID).
		SetAddressID(address.ID).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	contact, err := svc.dbClient.Contact.Create().
		SetContactDescription("phone").
		SetContactDetails("123-456-7890").SetContactType("phone").
		SetIsPrimaryContact(true).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.dbClient.CustomerContactOnClinics.Create().
		SetClinicID(clinic.ID).
		SetCustomerID(customer.ID).
		SetContactID(contact.ID).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	fullCustomer, err := svc.GetCustomer(customer.ID, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if int(fullCustomer.CustomerID) != customer.ID {
		t.Fatalf("expected customer ID %d, got %d", customer.ID, fullCustomer.CustomerID)
	}

	if len(fullCustomer.Clinics) != 1 {
		t.Fatal("customer should have one clinic")
	}
	if int(fullCustomer.Clinics[0].ClinicID) != clinic.ID {
		t.Fatal("unmatched clinic ID")
	}
	if len(fullCustomer.Clinics[0].CustomerAddresses) != 1 {
		t.Fatal("customer clinic should have one address")
	}
	if int(fullCustomer.Clinics[0].CustomerAddresses[0].AddressID) != address.ID {
		t.Fatal("unmatched address ID")
	}
	if len(fullCustomer.Clinics[0].CustomerContacts) != 1 {
		t.Fatal("customer clinic should have one contact")
	}
	if int(fullCustomer.Clinics[0].CustomerContacts[0].ContactID) != contact.ID {
		t.Fatal("unmatched contact ID")
	}
}

func TestListCustomer(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)

	ctx := context.Background()

	var expectedCustomerID int
	var expectedClinicID int
	var expectedAddressID int
	var expectedContactID int

	for i := 0; i < 5; i++ {
		customer, err := svc.dbClient.Customer.Create().
			SetCustomerFirstName(fmt.Sprintf("First%d", i)).
			SetCustomerLastName(fmt.Sprintf("Last%d", i)).
			SetCustomerTypeID("type_test").
			SetIsActive(true).
			Save(ctx)
		if err != nil {
			t.Fatalf("failed to create customer %d: %v", i, err)
		}

		clinic, err := svc.dbClient.Clinic.Create().
			SetClinicName(fmt.Sprintf("Clinic %d", i)).
			AddCustomerIDs(customer.ID).
			SetIsActive(true).
			Save(ctx)
		if err != nil {
			t.Fatalf("failed to create clinic %d: %v", i, err)
		}

		address, err := svc.dbClient.Address.Create().
			SetStreetAddress(fmt.Sprintf("Street %d", i)).
			SetCity("Testville").
			SetState("TS").
			SetZipcode("12345").
			SetIsPrimaryAddress(true).
			Save(ctx)
		if err != nil {
			t.Fatalf("failed to create address %d: %v", i, err)
		}

		_, err = svc.dbClient.CustomerAddressOnClinics.Create().
			SetClinicID(clinic.ID).
			SetCustomerID(customer.ID).
			SetAddressID(address.ID).
			Save(ctx)
		if err != nil {
			t.Fatalf("failed to bind address %d: %v", i, err)
		}

		contact, err := svc.dbClient.Contact.Create().
			SetContactDescription("phone").
			SetContactDetails(fmt.Sprintf("123-456-78%02d", i)).
			SetContactType("phone").
			SetIsPrimaryContact(true).
			Save(ctx)
		if err != nil {
			t.Fatalf("failed to create contact %d: %v", i, err)
		}

		_, err = svc.dbClient.CustomerContactOnClinics.Create().
			SetClinicID(clinic.ID).
			SetCustomerID(customer.ID).
			SetContactID(contact.ID).
			Save(ctx)
		if err != nil {
			t.Fatalf("failed to bind contact %d: %v", i, err)
		}

		if i == 0 {
			expectedCustomerID = customer.ID
			expectedClinicID = clinic.ID
			expectedAddressID = address.ID
			expectedContactID = contact.ID
		}
	}

	_ = svc.redisClient.Del(ctx, "lis::core_service_v2::customer::total_customer_count").Err()

	page := "1"
	perPage := "3"
	customers, hasNextPage, total, err := svc.ListCustomer(page, perPage, ctx)
	if err != nil {
		t.Fatalf("ListCustomer failed: %v", err)
	}

	if len(customers) != 3 {
		t.Fatalf("expected 3 customers, got %d", len(customers))
	}

	if !hasNextPage {
		t.Fatal("expected hasNextPage to be true")
	}

	if total != 5 {
		t.Fatalf("expected total to be 5, got %d", total)
	}

	var found bool
	for _, fc := range customers {
		if int(fc.CustomerID) == expectedCustomerID {
			found = true
			if len(fc.Clinics) != 1 {
				t.Fatal("customer should have one clinic")
			}
			if int(fc.Clinics[0].ClinicID) != expectedClinicID {
				t.Fatal("unmatched clinic ID")
			}
			if len(fc.Clinics[0].CustomerAddresses) != 1 {
				t.Fatal("clinic should have one address")
			}
			if int(fc.Clinics[0].CustomerAddresses[0].AddressID) != expectedAddressID {
				t.Fatal("unmatched address ID")
			}
			if len(fc.Clinics[0].CustomerContacts) != 1 {
				t.Fatal("clinic should have one contact")
			}
			if int(fc.Clinics[0].CustomerContacts[0].ContactID) != expectedContactID {
				t.Fatal("unmatched contact ID")
			}
		}
	}
	if !found {
		t.Fatal("expected customer not found in page results")
	}
}

func TestReinviteNPICheck(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)

	ctx := context.Background()

	status, err := svc.ReinviteNPICheck("Internal special NPI", []string{}, ctx)
	if err != nil {
		t.Fatal(err)
	}

	if status != "success" {
		t.Fatalf("expected status %s, got %s", "success", status)
	}

	statusType1, err := svc.ReinviteNPICheck("1649315763", []string{"provider", "clinicadmin"}, ctx)
	if statusType1 != "success" || err != nil {
		t.Fatal(err)
	}

	statusType2, err := svc.ReinviteNPICheck("1790598225", []string{"provider", "clinicadmin"}, ctx)
	if statusType2 == "success" || err == nil {
		t.Fatal("should fail since NPI type does not meet requirements for role: provider")
	}
}

func TestGetSalesCustomer(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)

	ctx := context.Background()

	// 创建 InternalUser (Sales)
	salesUser, err := svc.dbClient.InternalUser.Create().
		SetInternalUserName("john_sales").
		SetIsActive(true).
		SetInternalUserRole("sales").
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create internal user: %v", err)
	}

	// 创建 Customer 及 Clinic 并关联到 Sales
	customer, err := svc.dbClient.Customer.Create().
		SetCustomerFirstName("Alice").
		SetCustomerLastName("Smith").
		SetCustomerTypeID("type1").
		SetSalesID(salesUser.ID).
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create customer: %v", err)
	}

	clinic, err := svc.dbClient.Clinic.Create().
		SetClinicName("Sunrise Health").
		AddCustomerIDs(customer.ID).
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create clinic: %v", err)
	}

	address, err := svc.dbClient.Address.Create().
		SetStreetAddress("456 Clinic St").
		SetCity("HealthCity").
		SetState("HC").
		SetZipcode("67890").
		SetIsPrimaryAddress(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create address: %v", err)
	}

	_, err = svc.dbClient.CustomerAddressOnClinics.Create().
		SetClinicID(clinic.ID).
		SetCustomerID(customer.ID).
		SetAddressID(address.ID).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to bind address to clinic: %v", err)
	}

	contact, err := svc.dbClient.Contact.Create().
		SetContactDescription("Email").
		SetContactDetails("alice@sunrise.com").
		SetContactType("email").
		SetIsPrimaryContact(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create contact: %v", err)
	}

	_, err = svc.dbClient.CustomerContactOnClinics.Create().
		SetClinicID(clinic.ID).
		SetCustomerID(customer.ID).
		SetContactID(contact.ID).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to bind contact to clinic: %v", err)
	}

	// 测试目标函数
	results, err := svc.GetSalesCustomer([]string{"john_sales"}, "1", "10", ctx)
	if err != nil {
		t.Fatalf("GetSalesCustomer failed: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 customer, got %d", len(results))
	}

	fullCustomer := results[0]
	if int(fullCustomer.CustomerID) != customer.ID {
		t.Fatalf("expected customer ID %d, got %d", customer.ID, fullCustomer.CustomerID)
	}
	if len(fullCustomer.Clinics) != 1 {
		t.Fatalf("expected 1 clinic, got %d", len(fullCustomer.Clinics))
	}
	c := fullCustomer.Clinics[0]
	if int(c.ClinicID) != clinic.ID {
		t.Fatalf("expected clinic ID %d, got %d", clinic.ID, c.ClinicID)
	}
	if len(c.CustomerAddresses) != 1 {
		t.Fatal("clinic should have 1 address")
	}
	if int(c.CustomerAddresses[0].AddressID) != address.ID {
		t.Fatalf("expected address ID %d, got %d", address.ID, c.CustomerAddresses[0].AddressID)
	}
	if len(c.CustomerContacts) != 1 {
		t.Fatal("clinic should have 1 contact")
	}
	if int(c.CustomerContacts[0].ContactID) != contact.ID {
		t.Fatalf("expected contact ID %d, got %d", contact.ID, c.CustomerContacts[0].ContactID)
	}
}

func TestGetCustomerSales(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)

	ctx := context.Background()

	// 创建 InternalUser (Sales)
	sales, err := svc.dbClient.InternalUser.Create().
		SetInternalUserName("susan_sales").
		SetInternalUserFirstname("Susan").
		SetInternalUserLastname("Lee").
		SetIsActive(true).
		SetInternalUserRole("sales").
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create internal user: %v", err)
	}

	// 创建 Customer 并绑定 Sales
	customer, err := svc.dbClient.Customer.Create().
		SetCustomerFirstName("Bob").
		SetCustomerLastName("Brown").
		SetCustomerTypeID("type1").
		SetIsActive(true).
		SetSalesID(sales.ID).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create customer: %v", err)
	}

	nameResults, err := svc.GetCustomerSales([]string{"Bob Brown"}, nil, ctx)
	if err != nil {
		t.Fatalf("GetCustomerSales by name failed: %v", err)
	}
	if len(nameResults) != 1 {
		t.Fatalf("expected 1 result, got %d", len(nameResults))
	}
	salesData := nameResults[0].InternalUser
	if salesData == nil {
		t.Fatal("expected InternalUser to be not nil")
	}
	if salesData.InternalUserFirstname != "Susan" || salesData.InternalUserLastname != "Lee" {
		t.Fatalf("expected InternalUser name Susan Lee, got %s %s", salesData.InternalUserFirstname, salesData.InternalUserLastname)
	}

	// 测试通过 customer ID 查询
	idResults, err := svc.GetCustomerSales(nil, []string{strconv.Itoa(customer.ID)}, ctx)
	if err != nil {
		t.Fatalf("GetCustomerSales by ID failed: %v", err)
	}
	if len(idResults) != 1 {
		t.Fatalf("expected 1 result by ID, got %d", len(idResults))
	}
	salesDataByID := idResults[0].InternalUser
	if salesDataByID == nil {
		t.Fatal("expected InternalUser to be not nil in ID query")
	}
	if salesDataByID.InternalUserFirstname != "Susan" || salesDataByID.InternalUserLastname != "Lee" {
		t.Fatalf("expected InternalUser name Susan Lee in ID query, got %s %s", salesDataByID.InternalUserFirstname, salesDataByID.InternalUserLastname)
	}
}

func TestGetCustomerByNPINumber(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)

	ctx := context.Background()
	npi := "9999999999"

	// 插入两个具有相同 NPI 的客户
	customer1, err := svc.dbClient.Customer.Create().
		SetCustomerFirstName("John").
		SetCustomerLastName("Doe").
		SetCustomerNpiNumber(npi).
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create customer1: %v", err)
	}

	customer2, err := svc.dbClient.Customer.Create().
		SetCustomerFirstName("Jane").
		SetCustomerLastName("Smith").
		SetCustomerNpiNumber(npi).
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create customer2: %v", err)
	}

	// 插入一个不同 NPI 的客户
	_, err = svc.dbClient.Customer.Create().
		SetCustomerFirstName("Foo").
		SetCustomerLastName("Bar").
		SetCustomerNpiNumber("1234567890").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create customer3: %v", err)
	}

	// 查询
	customers, err := svc.GetCustomerByNPINumber(npi, ctx)
	if err != nil {
		t.Fatalf("GetCustomerByNPINumber failed: %v", err)
	}
	if len(customers) != 2 {
		t.Fatalf("expected 2 customers with NPI %s, got %d", npi, len(customers))
	}
	// 验证ID是否匹配
	found1, found2 := false, false
	for _, c := range customers {
		if c.ID == customer1.ID {
			found1 = true
		}
		if c.ID == customer2.ID {
			found2 = true
		}
	}
	if !found1 || !found2 {
		t.Fatal("returned customers do not match expected IDs")
	}
}

func TestFuzzySearchCustomers(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)

	ctx := context.Background()

	// 插入 Customer 数据
	customer1, _ := svc.dbClient.Customer.Create().
		SetCustomerFirstName("Alice").
		SetCustomerMiddleName("Marie").
		SetCustomerLastName("Johnson").
		SetIsActive(true).
		Save(ctx)

	customer2, _ := svc.dbClient.Customer.Create().
		SetCustomerFirstName("Bob").
		SetCustomerMiddleName("Lee").
		SetCustomerLastName("Smith").
		SetIsActive(true).
		Save(ctx)

	customer3, _ := svc.dbClient.Customer.Create().
		SetCustomerFirstName("Charlie").
		SetCustomerMiddleName("Alpha").
		SetCustomerLastName("Beta").
		SetIsActive(true).
		Save(ctx)

	// 插入 Clinic 并绑定 customer3
	clinic, _ := svc.dbClient.Clinic.Create().
		SetClinicName("TestClinic").
		SetIsActive(true).
		AddCustomerIDs(customer3.ID).
		Save(ctx)
	clinicIDStr := strconv.Itoa(clinic.ID)

	tests := []struct {
		name      string
		search    string
		clinicId  *string
		expectIDs []int
		expectErr bool
	}{
		{
			name:      "Fuzzy by one word name match",
			search:    "bob",
			expectIDs: []int{customer2.ID},
		},
		{
			name:      "Fuzzy by two words name match",
			search:    "alice johnson",
			expectIDs: []int{customer1.ID},
		},
		{
			name:      "Fuzzy by three words full match",
			search:    "alice marie johnson",
			expectIDs: []int{customer1.ID},
		},
		{
			name:      "Numeric ID match",
			search:    strconv.Itoa(customer2.ID),
			expectIDs: []int{customer2.ID},
		},
		{
			name:      "Scoped to clinicId",
			search:    "charlie",
			clinicId:  &clinicIDStr,
			expectIDs: []int{customer3.ID},
		},
		{
			name:      "Empty input error",
			search:    "   ",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := svc.FuzzySearchCustomers(tt.search, tt.clinicId, ctx)
			if tt.expectErr {
				if err == nil {
					t.Fatal("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var resultIDs []int
			for _, r := range results {
				resultIDs = append(resultIDs, int(r.ClientId))
			}

			if len(resultIDs) != len(tt.expectIDs) {
				t.Fatalf("expected %v, got %v", tt.expectIDs, resultIDs)
			}
			for _, id := range tt.expectIDs {
				found := false
				for _, rid := range resultIDs {
					if id == rid {
						found = true
					}
				}
				if !found {
					t.Fatalf("expected to find ID %d, but not found in result", id)
				}
			}
		})
	}
}

func TestFuzzySearchCustomerClinicName(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)

	ctx := context.Background()

	// 清空 redis 缓存，模拟 miss
	_ = svc.redisClient.Del(ctx, "customer_clinic_data").Err()

	// mock dbutils.FetchAndCacheCustomerClinicData
	originalFetch := dbutils.FetchAndCacheCustomerClinicData
	defer func() {
		dbutils.FetchAndCacheCustomerClinicData = originalFetch
	}()

	dbutils.FetchAndCacheCustomerClinicData = func(ctx context.Context, dbClient *ent.Client, redisClient *common.RedisClient) ([]*model.CustomerClinicData, error) {
		return []*model.CustomerClinicData{
			{CustomerId: 101, CustomerName: "Alice Johnson", ClinicName: "Sunrise Clinic"},
			{CustomerId: 102, CustomerName: "Bob Smith", ClinicName: "HealthCare Center"},
			{CustomerId: 103, CustomerName: "Charlie Doe", ClinicName: "Wellness Lab"},
		}, nil
	}

	tests := []struct {
		name      string
		search    string
		expectIDs []int
	}{
		{
			name:      "Search by numeric CustomerId",
			search:    "101",
			expectIDs: []int{101},
		},
		{
			name:      "Search by CustomerName part",
			search:    "alice",
			expectIDs: []int{101},
		},
		{
			name:      "Search by ClinicName part",
			search:    "care",
			expectIDs: []int{102},
		},
		{
			name:      "Search by unmatched input",
			search:    "zzz",
			expectIDs: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := svc.FuzzySearchCustomerClinicName(tt.search, ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			var gotIDs []int
			for _, r := range results {
				gotIDs = append(gotIDs, int(r.CustomerId))
			}
			if !util.EqualIntSlices(gotIDs, tt.expectIDs) {
				t.Fatalf("expected IDs %v, got %v", tt.expectIDs, gotIDs)
			}
		})
	}
}

func TestUpdateCustomerOnboardingQuestionnaireStatus(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)

	ctx := context.Background()

	// 🧪 Case 1: Invalid customer ID
	id, status := svc.UpdateCustomerOnboardingQuestionnaireStatus("abc", ctx)
	if id != 0 || status != "InvalidCustomerID" {
		t.Fatalf("expected (0, InvalidCustomerID), got (%d, %s)", id, status)
	}

	// 🧪 Case 2: Customer not found
	id, status = svc.UpdateCustomerOnboardingQuestionnaireStatus("999999", ctx)
	if id != 0 || status != "CustomerNotFound" {
		t.Fatalf("expected (0, CustomerNotFound), got (%d, %s)", id, status)
	}

	// 🧪 Case 3: Success path
	customer, err := svc.dbClient.Customer.Create().
		SetCustomerFirstName("Test").
		SetCustomerLastName("User").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create customer: %v", err)
	}

	idStr := strconv.Itoa(customer.ID)
	retID, status := svc.UpdateCustomerOnboardingQuestionnaireStatus(idStr, ctx)
	if retID != int32(customer.ID) || status != "OnboardingQuestionnaireUpdated" {
		t.Fatalf("expected (%d, OnboardingQuestionnaireUpdated), got (%d, %s)", customer.ID, retID, status)
	}

	// 🧪 Verify field was updated (optional)
	updated, err := svc.dbClient.Customer.Get(ctx, customer.ID)
	if err != nil {
		t.Fatalf("failed to fetch updated customer: %v", err)
	}
	if updated.OnboardingQuestionnaireFilledOn.IsZero() {
		t.Fatal("expected OnboardingQuestionnaireFilledOn to be set")
	}
}

func TestCheckCustomerOnboardingQuestionnaireStatus(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)

	ctx := context.Background()

	// 🧪 Case 1: Invalid customer ID
	_, err := svc.CheckCustomerOnboardingQuestionnaireStatus("abc", ctx)
	if err == nil || !strings.Contains(err.Error(), "invalid customer_id") {
		t.Fatalf("expected error for invalid ID, got: %v", err)
	}

	// 🧪 Case 2: Customer not found
	_, err = svc.CheckCustomerOnboardingQuestionnaireStatus("999999", ctx)
	if err == nil || !strings.Contains(err.Error(), "customer not found") {
		t.Fatalf("expected not found error, got: %v", err)
	}

	// 🧪 Case 3: Customer exists but has not filled questionnaire
	c1, err := svc.dbClient.Customer.Create().
		SetCustomerFirstName("Test").
		SetCustomerLastName("User").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create customer: %v", err)
	}
	resp, err := svc.CheckCustomerOnboardingQuestionnaireStatus(strconv.Itoa(c1.ID), ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.CustomerId != int32(c1.ID) || resp.IsOnboardingQuestionnaireFilled {
		t.Fatalf("expected IsOnboardingQuestionnaireFilled=false for %d", c1.ID)
	}

	// 🧪 Case 4: Customer exists and has filled questionnaire
	c2, err := svc.dbClient.Customer.Create().
		SetCustomerFirstName("Filled").
		SetCustomerLastName("Form").
		SetIsActive(true).
		SetOnboardingQuestionnaireFilledOn(time.Now()).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create customer with filled questionnaire: %v", err)
	}
	resp, err = svc.CheckCustomerOnboardingQuestionnaireStatus(strconv.Itoa(c2.ID), ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.IsOnboardingQuestionnaireFilled {
		t.Fatalf("expected questionnaire filled to be true for customer %d", c2.ID)
	}
	if resp.OnboardingQuestionnaireFilledOn == "" {
		t.Fatalf("expected non-empty filled-on timestamp")
	}
}

func TestAddCustomerWithNPINumberNative(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)

	ctx := context.Background()

	clinic, _ := svc.dbClient.Clinic.Create().
		SetClinicName("TestClinic").
		SetIsActive(true).
		Save(ctx)
	clinicIDStr := strconv.Itoa(clinic.ID)

	input := &model.AddCustomerWithNPINumberRequest{
		CustomerFirstName:         "John",
		CustomerLastName:          "Doe",
		CustomerNPINumber:         "1295351534",
		CustomerSuffix:            "Jr.",
		ClinicID:                  clinicIDStr,
		CustomerAddressLine1:      "123 Main St",
		CustomerAddressLine2:      "Apt 4",
		CustomerCity:              "Anytown",
		CustomerState:             "CA",
		CustomerZipcode:           "12345",
		CustomerCountry:           "USA",
		CustomerNotificationEmail: "john.doe@example.com",
		CustomerPhone:             "555-1234",
		CustomerRoles:             []string{"clinicadmin", "provider"},
		CustomerInvitationLink:    "http://example.com/invite",
	}

	resp, err := svc.AddCustomerWithNPINumberNative(input, ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "success" || resp.CustomerID == -1 || resp.ErrorMessage != "" {
		t.Fatalf("unexpected unsuccessful status")
	}

	createdCust, err := svc.dbClient.Customer.Get(ctx, int(resp.CustomerID))
	if err != nil {
		t.Fatalf("failed to fetch created customer: %v", err)
	}

	// Verify that the customer is linked to the correct clinic
	clinics, err := createdCust.QueryClinics().All(ctx)
	if err != nil {
		t.Fatalf("failed to query clinics for customer: %v", err)
	}
	if len(clinics) != 1 || clinics[0].ID != clinic.ID {
		t.Errorf("expected customer to be linked to clinic ID %d, got %+v", clinic.ID, clinics)
	}

	// Verify that exactly one address is linked to the customer and clinic
	addrs, err := svc.dbClient.CustomerAddressOnClinics.
		Query().
		Where(
			customeraddressonclinics.CustomerIDEQ(createdCust.ID),
			customeraddressonclinics.ClinicIDEQ(clinic.ID),
		).
		All(ctx)
	if err != nil {
		t.Fatalf("failed to query address links: %v", err)
	}
	if len(addrs) != 1 {
		t.Errorf("expected 1 address linked to customer and clinic, got %d", len(addrs))
	}
	addr, err := svc.dbClient.Address.Get(ctx, addrs[0].AddressID)
	if err != nil {
		t.Fatalf("failed to fetch address by ID: %v", err)
	}
	if addr.StreetAddress != input.CustomerAddressLine1 {
		t.Errorf("address mismatch: expected %s, got %s", input.CustomerAddressLine1, addr.StreetAddress)
	}

	// Verify that two contacts (email and phone) are linked to the customer and clinic
	contacts, err := svc.dbClient.CustomerContactOnClinics.
		Query().
		Where(
			customercontactonclinics.CustomerIDEQ(createdCust.ID),
			customercontactonclinics.ClinicIDEQ(clinic.ID),
		).
		All(ctx)
	if err != nil {
		t.Fatalf("failed to query contact links: %v", err)
	}
	if len(contacts) != 2 {
		t.Errorf("expected 2 contacts (email and phone), got %d", len(contacts))
	}
	for _, link := range contacts {
		contact, err := svc.dbClient.Contact.Get(ctx, link.ContactID)
		if err != nil {
			t.Fatalf("failed to fetch contact by ID: %v", err)
		}
		if contact.ContactType == "email" && contact.ContactDetails != input.CustomerNotificationEmail {
			t.Errorf("email contact mismatch: got %s, expected %s", contact.ContactDetails, input.CustomerNotificationEmail)
		}
		if contact.ContactType == "phone" && contact.ContactDetails != input.CustomerPhone {
			t.Errorf("phone contact mismatch: got %s, expected %s", contact.ContactDetails, input.CustomerPhone)
		}
	}

	invitationRecord, err := svc.dbClient.UserInvitationRecord.
		Query().
		Where(
			userinvitationrecord.CustomerIDEQ(createdCust.ID),
		).Only(ctx)

	if err != nil {
		t.Fatalf("failed to query invitation record: %v", err)
	}
	if invitationRecord.InvitationLink != input.CustomerInvitationLink {
		t.Errorf("invitation link mismatch: got %s, expected %s", invitationRecord.InvitationLink, input.CustomerInvitationLink)
	}

	tests := []struct {
		name           string
		input          *model.AddCustomerWithNPINumberRequest
		wantStatus     string
		wantCustomerID int32
	}{
		{
			name: "Invalid NPI",
			input: &model.AddCustomerWithNPINumberRequest{
				CustomerFirstName: "Jane",
				CustomerLastName:  "Smith",
				CustomerNPINumber: "invalid-npi",
				ClinicID:          clinicIDStr,
			},
			wantStatus:     "fail",
			wantCustomerID: -1,
		},
		{
			name: "NPI-2 with Provider Role",
			input: &model.AddCustomerWithNPINumberRequest{
				CustomerFirstName: "Bob",
				CustomerLastName:  "Johnson",
				CustomerNPINumber: "1790598225",
				ClinicID:          clinicIDStr,
				CustomerRoles:     []string{"provider"},
			},
			wantStatus:     "fail",
			wantCustomerID: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.AddCustomerWithNPINumberNative(tt.input, ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			require.NotNil(t, resp)
			assert.Equal(t, tt.wantStatus, resp.Status)

			if tt.wantCustomerID == -1 {
				assert.Equal(t, int32(-1), resp.CustomerID)
				assert.NotEmpty(t, resp.ErrorMessage)
			} else {
				assert.NotEqual(t, int32(-1), resp.CustomerID)
				assert.Equal(t, "", resp.ErrorMessage)
			}
		})
	}
}

func TestRemoveCustomerFromClinic(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)
	ctx := context.Background()

	// Case 1: invalid customerId
	status, msg, code := svc.RemoveCustomerFromClinic("abc", "123", "tester", ctx)
	if status != "failed" || msg != "invalid customer_id" || code != 400 {
		t.Errorf("expected invalid customer_id error, got %s %s %d", status, msg, code)
	}

	// Case 2: invalid clinicId
	status, msg, code = svc.RemoveCustomerFromClinic("1", "xyz", "tester", ctx)
	if status != "failed" || msg != "invalid clinic_id" || code != 400 {
		t.Errorf("expected invalid clinic_id error, got %s %s %d", status, msg, code)
	}

	// Case 3: clinic not found
	status, msg, code = svc.RemoveCustomerFromClinic("1", "9999", "tester", ctx)
	if status != "failed" || msg != "clinic not found" || code != 404 {
		t.Errorf("expected clinic not found error, got %s %s %d", status, msg, code)
	}

	createdUser1, err := svc.dbClient.User.Create().SetUserName("user1").SetPassword("test").Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	createdUser2, err := svc.dbClient.User.Create().SetUserName("user2").SetPassword("test").Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Setup: create clinic and customers
	createdClinic, err := svc.dbClient.Clinic.Create().
		SetClinicName("Test Clinic").
		SetUserID(createdUser1.ID).
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Case 4: customer not in clinic
	createdCust1, err := svc.dbClient.Customer.Create().
		SetCustomerFirstName("Orphan").
		SetUserID(createdUser2.ID).
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	status, msg, code = svc.RemoveCustomerFromClinic(strconv.Itoa(createdCust1.ID), strconv.Itoa(createdClinic.ID), "tester", ctx)
	if status != "failed" || msg != "no_such_relation" || code != 404 {
		t.Errorf("expected no_such_relation error, got %s %s %d", status, msg, code)
	}

	// Case 5: customer is founder (same userID as clinic)
	createdCust2, _ := svc.dbClient.Customer.Create().
		SetCustomerFirstName("Founder").
		SetUserID(createdUser1.ID). // same as clinic.UserID
		AddClinics(createdClinic).
		SetIsActive(true).
		Save(ctx)

	status, msg, code = svc.RemoveCustomerFromClinic(strconv.Itoa(createdCust2.ID), strconv.Itoa(createdClinic.ID), "tester", ctx)
	if status != "failed" || !strings.Contains(msg, "Founder") || code != 500 {
		t.Errorf("expected founder restriction error, got %s %s %d", status, msg, code)
	}

	// Case 6: normal customer removal
	createdCust3, _ := svc.dbClient.Customer.Create().
		SetCustomerFirstName("Regular").
		SetUserID(createdUser2.ID).
		AddClinics(createdClinic).
		SetIsActive(true).
		Save(ctx)

	status, msg, code = svc.RemoveCustomerFromClinic(strconv.Itoa(createdCust3.ID), strconv.Itoa(createdClinic.ID), "tester", ctx)
	if status != "success" || msg != "" || code != 200 {
		t.Errorf("expected success, got %s %s %d", status, msg, code)
	}

	// Verify disconnection
	updatedClinic, _ := svc.dbClient.Clinic.Query().Where(clinic.IDEQ(createdClinic.ID)).WithCustomers().Only(ctx)
	for _, c := range updatedClinic.Edges.Customers {
		if c.ID == createdCust3.ID {
			t.Errorf("customer was not actually removed from clinic")
		}
	}
}

func TestJoinCustomerToClinic_InvalidCustomerID(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)

	status, msg, code := svc.JoinCustomerToClinic("abc", "1", "tester", nil, context.Background())
	if code != 400 || msg != "invalid customer_id" {
		t.Errorf("expected invalid customer_id, got: %s %s %d", status, msg, code)
	}
}

func TestJoinCustomerToClinic_InvalidClinicID(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)

	status, msg, code := svc.JoinCustomerToClinic("1", "xyz", "tester", nil, context.Background())
	if code != 400 || msg != "invalid clinic_id" {
		t.Errorf("expected invalid clinic_id, got: %s %s %d", status, msg, code)
	}
}

func TestJoinCustomerToClinic_ClinicNotFound(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)

	status, msg, code := svc.JoinCustomerToClinic("1", "9999", "tester", nil, context.Background())
	if code != 404 {
		t.Errorf("expected clinic not found, got: %s %s %d", status, msg, code)
	}
}

func TestJoinCustomerToClinic_CustomerNotFound(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)
	ctx := context.Background()

	clinic, _ := svc.dbClient.Clinic.Create().SetClinicName("TestClinic").SetIsActive(true).Save(ctx)

	status, msg, code := svc.JoinCustomerToClinic("9999", strconv.Itoa(clinic.ID), "tester", nil, ctx)
	if code != 404 || !strings.Contains(msg, "customer") {
		t.Errorf("expected customer not found, got: %s %s %d", status, msg, code)
	}
}

func TestJoinCustomerToClinic_CustomerAlreadyInClinic(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)
	ctx := context.Background()

	clinic, _ := svc.dbClient.Clinic.Create().SetClinicName("TestClinic").SetIsActive(true).Save(ctx)

	customer, _ := svc.dbClient.Customer.Create().
		SetCustomerFirstName("Joiner").
		SetCustomerNpiNumber("NPI123").
		SetIsActive(true).
		AddClinics(clinic).
		Save(ctx)

	status, msg, code := svc.JoinCustomerToClinic(strconv.Itoa(customer.ID), strconv.Itoa(clinic.ID), "tester", nil, ctx)
	if code != 200 || !strings.Contains(msg, "Already in the Clinic") {
		t.Errorf("expected already in clinic, got: %s %s %d", status, msg, code)
	}
}

func TestJoinCustomerToClinic_NPIDuplicate(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)
	ctx := context.Background()

	clinic, _ := svc.dbClient.Clinic.Create().SetClinicName("TestClinic").SetIsActive(true).Save(ctx)

	_ = svc.dbClient.Customer.Create().
		SetCustomerFirstName("Dup").
		SetCustomerNpiNumber("DUPNPI").
		SetIsActive(true).
		AddClinics(clinic).
		SaveX(ctx)

	cust, _ := svc.dbClient.Customer.Create().
		SetCustomerFirstName("New").
		SetCustomerNpiNumber("DUPNPI").
		SetIsActive(true).
		Save(ctx)

	status, msg, code := svc.JoinCustomerToClinic(strconv.Itoa(cust.ID), strconv.Itoa(clinic.ID), "tester", nil, ctx)
	if code != 500 || !strings.Contains(msg, "NPI Number Duplicate") {
		t.Errorf("expected NPI duplicate failure, got: %s %s %d", status, msg, code)
	}
}

func TestJoinCustomerToClinic_Success(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)
	ctx := context.Background()

	err := svc.rbacService.CreateRole("provider", "external", 0, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = svc.rbacService.CreateRole("clinicadmin", "external", 0, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}

	user, _ := svc.dbClient.User.Create().SetUserName("userX").SetPassword("pwd").Save(ctx)
	clinic, _ := svc.dbClient.Clinic.Create().SetClinicName("TestClinic").SetIsActive(true).Save(ctx)

	customer, _ := svc.dbClient.Customer.Create().
		SetCustomerFirstName("Unique").
		SetCustomerNpiNumber("1295351534").
		SetUserID(user.ID).
		SetIsActive(true).
		Save(ctx)

	status, msg, code := svc.JoinCustomerToClinic(strconv.Itoa(customer.ID), strconv.Itoa(clinic.ID), "tester", []string{"provider", "clinicadmin"}, ctx)
	if code != 200 || status != "success" {
		t.Errorf("expected success, got: %s %s %d", status, msg, code)
	}

	roles, err := svc.rbacService.GetAccountRolesInDomain(int32(user.ID), "clinic", int32(clinic.ID), ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) < 2 {
		t.Fatalf("expected at least one role, got none")
	}
	expectedRoles := map[string]bool{
		"provider":    false,
		"clinicadmin": false,
	}

	for _, role := range roles {
		if _, ok := expectedRoles[role.Name]; ok {
			expectedRoles[role.Name] = true
		}
	}

	for roleName, found := range expectedRoles {
		if !found {
			t.Fatalf("expected role '%s' not found in domain", roleName)
		}
	}
}

func TestJoinCustomerToClinic_NPI2ProviderFail(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)
	ctx := context.Background()

	err := svc.rbacService.CreateRole("provider", "external", 0, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = svc.rbacService.CreateRole("clinicadmin", "external", 0, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}

	user, _ := svc.dbClient.User.Create().SetUserName("userX").SetPassword("pwd").Save(ctx)
	clinic, _ := svc.dbClient.Clinic.Create().SetClinicName("TestClinic").SetIsActive(true).Save(ctx)

	customer, _ := svc.dbClient.Customer.Create().
		SetCustomerFirstName("Unique").
		SetCustomerNpiNumber("1134835283").
		SetUserID(user.ID).
		SetIsActive(true).
		Save(ctx)

	status, msg, code := svc.JoinCustomerToClinic(strconv.Itoa(customer.ID), strconv.Itoa(clinic.ID), "tester", []string{"provider", "clinicadmin"}, ctx)
	if status != "partial success" || !strings.Contains(msg, "non NPI-1 type") || code != 200 {
		t.Errorf("expected partial success due to NPI-2 restriction, got: %s %s %d", status, msg, code)
	}

	roles, err := svc.rbacService.GetAccountRolesInDomain(int32(user.ID), "clinic", int32(clinic.ID), ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	foundProvider := false
	for _, r := range roles {
		if r.Name == "provider" {
			foundProvider = true
		}
	}
	if foundProvider {
		t.Error("expected provider role not to be assigned for NPI-2")
	}
}

func TestFetchCustomerBetaProgramsForClinic(t *testing.T) {
	svc, redisServer := setupCustomerTest(t)
	defer cleanupCustomerTest(svc, redisServer)
	ctx := context.Background()

	// Create customer
	customer, err := svc.dbClient.Customer.Create().
		SetCustomerFirstName("Beta").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create customer: %v", err)
	}

	// Create clinic
	clinic, err := svc.dbClient.Clinic.Create().
		SetClinicName("BetaClinic").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create clinic: %v", err)
	}

	// Create beta programs
	beta1, err := svc.dbClient.BetaProgram.Create().
		SetBetaProgramName("Program A").
		SetBetaProgramDescription("Program A").
		SetUpdatedTime(time.Now()).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create beta program A: %v", err)
	}

	beta2, err := svc.dbClient.BetaProgram.Create().
		SetBetaProgramName("Program B").
		SetBetaProgramDescription("Program B").
		SetUpdatedTime(time.Now()).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create beta program B: %v", err)
	}

	// Create participations
	_, err = svc.dbClient.BetaProgramParticipation.Create().
		SetCustomerID(customer.ID).
		SetClinicID(clinic.ID).
		SetBetaProgramID(beta1.ID).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create beta participation for beta1: %v", err)
	}

	_, err = svc.dbClient.BetaProgramParticipation.Create().
		SetCustomerID(customer.ID).
		SetClinicID(clinic.ID).
		SetBetaProgramID(beta2.ID).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create beta participation for beta2: %v", err)
	}

	// Case 1: both IDs are 0
	res, msg := svc.FetchCustomerBetaProgramsForClinic(0, 0, ctx)
	if msg != "No Customer ID/Clinic ID Inputted" {
		t.Errorf("expected validation error, got: %s", msg)
	}
	if res != nil {
		t.Errorf("expected nil result when no IDs passed")
	}

	// Case 2: customerID + clinicID
	res, msg = svc.FetchCustomerBetaProgramsForClinic(int32(customer.ID), int32(clinic.ID), ctx)
	if msg != "" {
		t.Errorf("unexpected error: %s", msg)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if len(res[0].BetaPrograms) != 2 {
		t.Errorf("expected 2 programs, got %v", res[0].BetaPrograms)
	}

	// Case 3: only customerID
	res, msg = svc.FetchCustomerBetaProgramsForClinic(int32(customer.ID), 0, ctx)
	if msg != "" || len(res) != 1 {
		t.Errorf("expected match by customerID only, got msg: %s, res: %v", msg, res)
	}

	// Case 4: only clinicID
	res, msg = svc.FetchCustomerBetaProgramsForClinic(0, int32(clinic.ID), ctx)
	if msg != "" || len(res) != 1 {
		t.Errorf("expected match by clinicID only, got msg: %s, res: %v", msg, res)
	}

	// Case 5: no match
	res, msg = svc.FetchCustomerBetaProgramsForClinic(9999, 9999, ctx)
	if msg != "" {
		t.Errorf("unexpected error msg for no match: %s", msg)
	}
	if len(res) != 0 {
		t.Errorf("expected empty result for unmatched input, got %v", res)
	}
}
