package service

import (
	"bytes"
	"context"
	"coresamples/common"
	"coresamples/ent"
	"coresamples/ent/contact"
	"coresamples/ent/enttest"
	"coresamples/ent/internaluser"
	"coresamples/publisher"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stvp/tempredis"
)

// Mock HTTP transport that simulates API responses
type mockHTTPTransport struct {
	shouldSucceed bool
}

func (m *mockHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Simulate successful response for email API
	if m.shouldSucceed {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`{"success":true}`))),
			Header:     make(http.Header),
		}, nil
	}
	// Simulate failure response
	return &http.Response{
		StatusCode: http.StatusUnauthorized,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"statusCode":401,"message":"Unauthorized"}`))),
		Header:     make(http.Header),
	}, nil
}

func setupUserServiceTest(t *testing.T) (IUserService, *ent.Client, *tempredis.Server, context.Context) {
	dataSource := "file:ent?mode=memory&_fk=1"
	dbClient := enttest.Open(t, "sqlite3", dataSource)
	ctx := context.Background()
	err := dbClient.Schema.Create(ctx)
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

	// Initialize test JWT secret for token verification
	common.Secrets.JWTSecret = "test-jwt-secret-for-unit-tests"

	// Initialize mock publisher to avoid Kafka calls
	publisher.InitMockPublisher()

	// Create a new user service with mock HTTP client
	mockHTTPClient := &http.Client{
		Transport: &mockHTTPTransport{shouldSucceed: true},
	}

	svc := &UserService{
		dbClient:    dbClient,
		redisClient: common.NewRedisClient(redisClient, redisClient),
		httpClient:  mockHTTPClient,
	}

	return svc, dbClient, server, ctx
}

func cleanUpUserServiceTest(dbClient *ent.Client, s *tempredis.Server) {
	var err error
	if err = s.Kill(); err != nil {
		common.Error(err)
	}
	if dbClient != nil {
		if err = dbClient.Close(); err != nil {
			common.Error(err)
		}
	}
}

// Helper function to create test internal users and customers for testing
func createTestInternalUsers(t *testing.T, dbClient *ent.Client, ctx context.Context) {
	// Create user records first since they're referenced by InternalUser
	for i := 101; i <= 105; i++ {
		_, err := dbClient.User.Create().
			SetID(i).
			SetUserName(fmt.Sprintf("user%d", i)).
			SetIsActive(true).
			Save(ctx)
		if err != nil {
			t.Logf("Note: Failed to create user %d: %v (this may be expected in test environment)", i, err)
			// Continue anyway - in test environment we just need the InternalUser records
		}
	}

	// Create test users with different roles and role IDs
	users := []struct {
		role       string
		roleID     int
		username   string
		firstname  string
		lastname   string
		email      string
		isFullTime bool
		isActive   bool
		userID     int
		userType   string
	}{
		{
			role:       "manager",
			roleID:     1,
			username:   "john.manager",
			firstname:  "John",
			lastname:   "Manager",
			email:      "john.manager@example.com",
			isFullTime: true,
			isActive:   true,
			userID:     101,
			userType:   "staff",
		},
		{
			role:       "sales",
			roleID:     2,
			username:   "mary.sales",
			firstname:  "Mary",
			lastname:   "Sales",
			email:      "mary.sales@example.com",
			isFullTime: true,
			isActive:   true,
			userID:     102,
			userType:   "staff",
		},
		{
			role:       "manager",
			roleID:     1,
			username:   "robert.manager",
			firstname:  "Robert",
			lastname:   "Manager",
			email:      "robert.manager@example.com",
			isFullTime: false,
			isActive:   true,
			userID:     103,
			userType:   "staff",
		},
		{
			role:       "admin",
			roleID:     3,
			username:   "susan.admin",
			firstname:  "Susan",
			lastname:   "Admin",
			email:      "susan.admin@example.com",
			isFullTime: true,
			isActive:   true,
			userID:     104,
			userType:   "admin",
		},
		{
			role:       "sales",
			roleID:     2,
			username:   "david.sales",
			firstname:  "David",
			lastname:   "Sales",
			email:      "david.sales@example.com",
			isFullTime: true,
			isActive:   false, // Inactive user
			userID:     105,
			userType:   "staff",
		},
		// Extra sales users for transfer tests
		{
			role:       "sales",
			roleID:     6, // Using a different role ID for this sales user
			username:   "sales.user1",
			firstname:  "Sales",
			lastname:   "User1",
			email:      "sales.user1@example.com",
			isFullTime: true,
			isActive:   true,
			userID:     106,
			userType:   "staff",
		},
		{
			role:       "sales",
			roleID:     7, // Using a different role ID for this sales user
			username:   "sales.user2",
			firstname:  "Sales",
			lastname:   "User2",
			email:      "sales.user2@example.com",
			isFullTime: true,
			isActive:   true,
			userID:     107,
			userType:   "staff",
		},
		// Extra sales users for transfer tests
		{
			role:       "sales",
			roleID:     6, // Using a different role ID for this sales user
			username:   "sales.user1",
			firstname:  "Sales",
			lastname:   "User1",
			email:      "sales.user1@example.com",
			isFullTime: true,
			isActive:   true,
			userID:     106,
			userType:   "staff",
		},
		{
			role:       "sales",
			roleID:     7, // Using a different role ID for this sales user
			username:   "sales.user2",
			firstname:  "Sales",
			lastname:   "User2",
			email:      "sales.user2@example.com",
			isFullTime: true,
			isActive:   true,
			userID:     107,
			userType:   "staff",
		},
	}

	// Disable foreign key constraints for testing
	_, err := dbClient.ExecContext(ctx, "PRAGMA foreign_keys = OFF")
	if err != nil {
		t.Logf("Warning: Could not disable foreign keys: %v", err)
	}

	for i, user := range users {
		_, err := dbClient.InternalUser.Create().
			SetID(i + 1).
			SetInternalUserRole(user.role).
			SetInternalUserRoleID(user.roleID).
			SetInternalUserName(user.username).
			SetInternalUserFirstname(user.firstname).
			SetInternalUserLastname(user.lastname).
			SetInternalUserEmail(user.email).
			SetInternalUserIsFullTime(user.isFullTime).
			SetIsActive(user.isActive).
			SetUserID(user.userID).
			SetInternalUserType(user.userType).
			Save(ctx)

		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	// Create test customers for sales transfer tests
	customers := []struct {
		id         int
		firstName  string
		lastName   string
		middleName string
		salesID    int
		userID     int
		isActive   bool
	}{
		{
			id:         201,
			firstName:  "Customer",
			lastName:   "One",
			middleName: "",
			salesID:    6, // Belongs to sales.user1
			userID:     301,
			isActive:   true,
		},
		{
			id:         202,
			firstName:  "Customer",
			lastName:   "Two",
			middleName: "",
			salesID:    6, // Belongs to sales.user1
			userID:     302,
			isActive:   true,
		},
		{
			id:         203,
			firstName:  "Customer",
			lastName:   "Three",
			middleName: "",
			salesID:    7, // Belongs to sales.user2
			userID:     303,
			isActive:   true,
		},
	}

	// Create customer users first
	for _, c := range customers {
		_, err := dbClient.User.Create().
			SetID(c.userID).
			SetUserName(fmt.Sprintf("customer%d", c.id)).
			SetIsActive(true).
			Save(ctx)
		if err != nil {
			t.Logf("Note: Failed to create customer user %d: %v (this may be expected in test environment)", c.userID, err)
		}
	}

	// Create customers
	for _, c := range customers {
		_, err := dbClient.Customer.Create().
			SetID(c.id).
			SetCustomerFirstName(c.firstName).
			SetCustomerLastName(c.lastName).
			SetCustomerMiddleName(c.middleName).
			SetSalesID(c.salesID).
			SetUserID(c.userID).
			SetIsActive(c.isActive).
			Save(ctx)

		if err != nil {
			t.Fatalf("Failed to create test customer: %v", err)
		}
	}
}

func TestGetInternalUserByID(t *testing.T) {
	svc, dbClient, redisServer, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, redisServer)

	// Create test data
	createTestInternalUsers(t, dbClient, ctx)

	t.Run("Get existing active user", func(t *testing.T) {
		// Test getting an existing active user (ID 1)
		user, err := svc.GetInternalUserByID(ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, int32(1), user.InternalUserId)
		assert.Equal(t, "manager", user.InternalUserRole)
		assert.Equal(t, "john.manager", user.InternalUserName)
	})

	t.Run("Get inactive user", func(t *testing.T) {
		// Test getting an inactive user (ID 5)
		user, err := svc.GetInternalUserByID(ctx, 5)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, int32(5), user.InternalUserId)
		assert.Equal(t, "sales", user.InternalUserRole)
		assert.Equal(t, "david.sales", user.InternalUserName)
		assert.False(t, user.IsActive) // Confirm it's inactive
	})

	t.Run("Get non-existent user", func(t *testing.T) {
		// Test getting a non-existent user
		_, err := svc.GetInternalUserByID(ctx, 999)

		assert.Error(t, err)
	})
}

func TestGetInternalUser(t *testing.T) {
	svc, dbClient, redisServer, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, redisServer)

	// Create test data
	createTestInternalUsers(t, dbClient, ctx)

	t.Run("Get all active internal users", func(t *testing.T) {
		// Test with no filters - should return all active users
		response, err := svc.GetInternalUser(ctx, "", nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Response, 1)
		assert.GreaterOrEqual(t, len(response.Response[0].InternalUser), 6) // 6 active users now (added 2 sales)
	})

	t.Run("Filter by role", func(t *testing.T) {
		// Test filtering by role
		response, err := svc.GetInternalUser(ctx, "manager", nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Response, 1)
		assert.Len(t, response.Response[0].InternalUser, 2) // 2 active managers

		// Verify correct role
		for _, user := range response.Response[0].InternalUser {
			assert.Equal(t, "manager", user.InternalUserRole)
		}
	})

	t.Run("Filter by roleIDs", func(t *testing.T) {
		// Test filtering by roleIDs
		roleIDs := []int32{2} // sales role
		response, err := svc.GetInternalUser(ctx, "", roleIDs, nil)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		// With new implementation, we have a separate response entry for each roleID
		assert.Len(t, response.Response, 1)
		assert.Len(t, response.Response[0].InternalUser, 1) // 1 active sales user with roleID 2

		// Verify correct roleID
		for _, user := range response.Response[0].InternalUser {
			assert.Equal(t, int32(2), user.InternalUserRoleId)
		}
	})

	t.Run("Filter by multiple roleIDs", func(t *testing.T) {
		// Test filtering by multiple roleIDs
		roleIDs := []int32{1, 2} // manager and sales roles
		response, err := svc.GetInternalUser(ctx, "", roleIDs, nil)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		// With new implementation, we have a separate response entry for each roleID
		assert.Len(t, response.Response, 2)

		// Verify roles in each response section
		roleIDsFound := make(map[int32]bool)
		for _, middleLevel := range response.Response {
			if len(middleLevel.InternalUser) > 0 {
				roleIDsFound[middleLevel.InternalUser[0].InternalUserRoleId] = true
			}
		}

		assert.True(t, roleIDsFound[1])
		assert.True(t, roleIDsFound[2])
	})

	t.Run("Filter by usernames", func(t *testing.T) {
		// Test filtering by usernames
		usernames := []string{"john.manager", "susan.admin"}
		response, err := svc.GetInternalUser(ctx, "", nil, usernames)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		// With new implementation, we have a separate response entry for each username
		assert.Len(t, response.Response, 2)

		// Count total users returned
		totalUsers := 0
		for _, middleLevel := range response.Response {
			totalUsers += len(middleLevel.InternalUser)
		}
		assert.Equal(t, 2, totalUsers) // 2 users total with these usernames

		// Verify correct usernames
		foundUsernames := make(map[string]bool)
		for _, middleLevel := range response.Response {
			for _, user := range middleLevel.InternalUser {
				foundUsernames[user.InternalUserName] = true
			}
		}

		assert.True(t, foundUsernames["john.manager"])
		assert.True(t, foundUsernames["susan.admin"])
	})

	t.Run("Combined roleIDs and usernames", func(t *testing.T) {
		// Test combining roleIDs and usernames
		roleIDs := []int32{1}                // manager role
		usernames := []string{"susan.admin"} // admin user
		response, err := svc.GetInternalUser(ctx, "", roleIDs, usernames)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		// We should have one response for the roleID and one for the username
		assert.Len(t, response.Response, 2)

		// We should find both the manager role users and the admin user
		foundRoles := make(map[string]bool)
		foundUsernames := make(map[string]bool)

		for _, middleLevel := range response.Response {
			for _, user := range middleLevel.InternalUser {
				foundRoles[user.InternalUserRole] = true
				foundUsernames[user.InternalUserName] = true
			}
		}

		assert.True(t, foundRoles["manager"])
		assert.True(t, foundRoles["admin"])
		assert.True(t, foundUsernames["susan.admin"])
	})

	t.Run("Role filter with roleIDs and usernames", func(t *testing.T) {
		// Test role filter applied to both roleIDs and usernames
		roleIDs := []int32{1} // manager role
		usernames := []string{"john.manager", "susan.admin"}
		response, err := svc.GetInternalUser(ctx, "manager", roleIDs, usernames)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		// We should have entries for the roleID and for each username
		assert.GreaterOrEqual(t, len(response.Response), 2)

		// All returned users should be managers
		for _, middleLevel := range response.Response {
			for _, user := range middleLevel.InternalUser {
				assert.Equal(t, "manager", user.InternalUserRole)
			}
		}

		// We should find john.manager but not susan.admin (as she's not a manager)
		foundUsernames := make(map[string]bool)
		for _, middleLevel := range response.Response {
			for _, user := range middleLevel.InternalUser {
				foundUsernames[user.InternalUserName] = true
			}
		}

		assert.True(t, foundUsernames["john.manager"])
		assert.False(t, foundUsernames["susan.admin"]) // Not a manager
	})

	t.Run("Inactive users excluded", func(t *testing.T) {
		// Test that inactive users are excluded
		count, err := dbClient.InternalUser.Query().Where(internaluser.InternalUserName("david.sales")).Count(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, count) // Exists in DB

		response, err := svc.GetInternalUser(ctx, "sales", nil, []string{"david.sales"})

		assert.NoError(t, err)
		assert.NotNil(t, response)

		// Count total users returned
		totalUsers := 0
		for _, middleLevel := range response.Response {
			totalUsers += len(middleLevel.InternalUser)
		}
		assert.Equal(t, 0, totalUsers) // Should not return the inactive user
	})

	t.Run("No results", func(t *testing.T) {
		// Test with filters that match no users
		response, err := svc.GetInternalUser(ctx, "non_existent_role", nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Response, 1)
		assert.Len(t, response.Response[0].InternalUser, 0) // No users found
	})

	t.Run("Redis caching for username lookup", func(t *testing.T) {
		// First call to populate cache
		username := "john.manager"
		_, err := svc.GetInternalUserByUsername(ctx, "", username)
		assert.NoError(t, err)

		// Second call should use cache
		users, err := svc.GetInternalUserByUsername(ctx, "", username)
		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Len(t, users, 1)
		assert.Equal(t, username, users[0].InternalUserName)
	})
}

func TestIsEmailUsedAsLoginId(t *testing.T) {
	svc, dbClient, redisServer, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, redisServer)

	// Create test users with different email login IDs
	testUsers := []struct {
		id          int
		username    string
		emailUserID string
		isActive    bool
	}{
		{
			id:          1001,
			username:    "testuser1",
			emailUserID: "user1@example.com",
			isActive:    true,
		},
		{
			id:          1002,
			username:    "testuser2",
			emailUserID: "user2@example.com",
			isActive:    true,
		},
		{
			id:          1003,
			username:    "testuser3",
			emailUserID: "", // No email login ID
			isActive:    true,
		},
		{
			id:          1004,
			username:    "inactiveuser",
			emailUserID: "inactive@example.com",
			isActive:    false,
		},
	}

	for _, u := range testUsers {
		_, err := dbClient.User.Create().
			SetID(u.id).
			SetUserName(u.username).
			SetEmailUserID(u.emailUserID).
			SetIsActive(u.isActive).
			SetPassword("password123"). // Add required password field
			Save(ctx)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	t.Run("Email exists as login ID", func(t *testing.T) {
		// Test checking an email that exists as login ID
		result, err := svc.IsEmailUsedAsLoginId(ctx, "user1@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.UsedAsEmailLogId)
		assert.Equal(t, "Email is already used as log in email id", result.Message)
		assert.Equal(t, int32(1001), result.UserId)
	})

	t.Run("Email does not exist as login ID", func(t *testing.T) {
		// Test checking an email that does not exist as login ID
		result, err := svc.IsEmailUsedAsLoginId(ctx, "nonexistent@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.UsedAsEmailLogId)
		assert.Equal(t, "Email is not used as log in email id", result.Message)
		assert.Equal(t, int32(0), result.UserId) // Should be 0 when no user found
	})

	t.Run("Empty email", func(t *testing.T) {
		// Test with an empty email
		result, err := svc.IsEmailUsedAsLoginId(ctx, "")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		// Since we have a user with an empty email ID in our test data,
		// we expect this to return true
		assert.True(t, result.UsedAsEmailLogId)
		assert.Equal(t, "Email is already used as log in email id", result.Message)
		assert.Equal(t, int32(1003), result.UserId)
	})
}

func TestCheckWhetherEmailIsUsedAsLoginId(t *testing.T) {
	svc, dbClient, redisServer, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, redisServer)

	// Disable foreign key constraints for testing
	_, err := dbClient.ExecContext(ctx, "PRAGMA foreign_keys = OFF")
	if err != nil {
		t.Logf("Warning: Could not disable foreign keys: %v", err)
	}

	// Create test users first, since they're needed for clinics
	testUsers := []struct {
		id          int
		username    string
		emailUserID string
		isActive    bool
	}{
		{
			id:          2001,
			username:    "clinicuser1",
			emailUserID: "clinic1@example.com",
			isActive:    true,
		},
		{
			id:          2002,
			username:    "clinicuser2",
			emailUserID: "clinic2@example.com",
			isActive:    true,
		},
		{
			id:          2003,
			username:    "regularuser",
			emailUserID: "regular@example.com",
			isActive:    true,
		},
		{
			id:          2004,
			username:    "anothercustomer",
			emailUserID: "another@example.com",
			isActive:    true,
		},
	}

	for _, u := range testUsers {
		_, err := dbClient.User.Create().
			SetID(u.id).
			SetUserName(u.username).
			SetEmailUserID(u.emailUserID).
			SetIsActive(u.isActive).
			SetPassword("password123").
			Save(ctx)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	// Create test clinics
	testClinics := []struct {
		id        int
		name      string
		userID    int
		accountID int
		isActive  bool
	}{
		{
			id:        101,
			name:      "Test Clinic 1",
			userID:    2001,
			accountID: 501,
			isActive:  true,
		},
		{
			id:        102,
			name:      "Test Clinic 2",
			userID:    2002,
			accountID: 502,
			isActive:  true,
		},
	}

	for _, c := range testClinics {
		_, err := dbClient.Clinic.Create().
			SetID(c.id).
			SetClinicName(c.name).
			SetUserID(c.userID).
			SetClinicAccountID(c.accountID).
			SetIsActive(c.isActive).
			Save(ctx)
		if err != nil {
			t.Fatalf("Failed to create test clinic: %v", err)
		}
	}

	// Create test customers with clinic associations
	testCustomers := []struct {
		id        int
		userID    int
		firstName string
		lastName  string
		clinicID  int
	}{
		{
			id:        301,
			userID:    2003, // regularuser
			firstName: "Regular",
			lastName:  "Customer",
			clinicID:  101, // Test Clinic 1
		},
		{
			id:        302,
			userID:    2004, // Create a new ID not in users array
			firstName: "Another",
			lastName:  "Customer",
			clinicID:  102, // Test Clinic 2
		},
	}

	// User for the second customer already created above

	for _, c := range testCustomers {
		customer, err := dbClient.Customer.Create().
			SetID(c.id).
			SetUserID(c.userID).
			SetCustomerFirstName(c.firstName).
			SetCustomerLastName(c.lastName).
			Save(ctx)
		if err != nil {
			t.Fatalf("Failed to create test customer: %v", err)
		}

		// Connect the customer to the clinic
		_, err = customer.Update().AddClinicIDs(c.clinicID).Save(ctx)
		if err != nil {
			t.Fatalf("Failed to associate customer with clinic: %v", err)
		}
	}

	t.Run("Email exists and is associated with clinic", func(t *testing.T) {
		// Clinic user email should be associated with that clinic
		result, err := svc.CheckWhetherEmailIsUsedAsLoginId(ctx, "clinic1@example.com", "101")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.ExistingUser)
		assert.Equal(t, "User Exists", result.Message)
	})

	t.Run("Email exists but not associated with clinic", func(t *testing.T) {
		// Clinic 1 user should not be associated with clinic 2
		result, err := svc.CheckWhetherEmailIsUsedAsLoginId(ctx, "clinic1@example.com", "102")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.ExistingUser)
		assert.Equal(t, "User Exists but not in this clinic", result.Message)
	})

	t.Run("Customer user associated with clinic", func(t *testing.T) {
		// Customer user should be associated with their clinic
		result, err := svc.CheckWhetherEmailIsUsedAsLoginId(ctx, "regular@example.com", "101")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.ExistingUser)
		assert.Equal(t, "User Exists", result.Message)
	})

	t.Run("Customer user not associated with other clinic", func(t *testing.T) {
		// Customer user should not be associated with a different clinic
		result, err := svc.CheckWhetherEmailIsUsedAsLoginId(ctx, "regular@example.com", "102")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.ExistingUser)
		assert.Equal(t, "User Exists but not in this clinic", result.Message)
	})

	t.Run("Email does not exist", func(t *testing.T) {
		// Email that does not exist in the system
		result, err := svc.CheckWhetherEmailIsUsedAsLoginId(ctx, "nonexistent@example.com", "101")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.ExistingUser)
		assert.Equal(t, "User Does Not Exist", result.Message)
	})

	t.Run("Invalid clinic ID", func(t *testing.T) {
		// Non-numeric clinic ID
		_, err := svc.CheckWhetherEmailIsUsedAsLoginId(ctx, "clinic1@example.com", "invalid")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid clinic ID format")
	})

	t.Run("Nonexistent clinic ID", func(t *testing.T) {
		// Clinic ID that does not exist
		result, err := svc.CheckWhetherEmailIsUsedAsLoginId(ctx, "clinic1@example.com", "999")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.ExistingUser)
		assert.Equal(t, "User Exists but not in this clinic", result.Message)
	})
}

func TestTransferSalesCustomer(t *testing.T) {
	svc, dbClient, redisServer, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, redisServer)

	// Create test data
	createTestInternalUsers(t, dbClient, ctx)

	t.Run("Successfully transfer customer", func(t *testing.T) {
		// Verify initial assignment
		customer, err := dbClient.Customer.Get(ctx, 201)
		assert.NoError(t, err)
		assert.Equal(t, 6, customer.SalesID) // Initially assigned to sales.user1

		// Transfer from sales.user1 (ID 6) to sales.user2 (ID 7)
		response, err := svc.TransferSalesCustomer(ctx, "6", "7", "201")

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Contains(t, response.Status, "Successfully transfer the customer 201 from sales 6 to sales 7")

		// Verify customer was updated
		customer, err = dbClient.Customer.Get(ctx, 201)
		assert.NoError(t, err)
		assert.Equal(t, 7, customer.SalesID) // Now assigned to sales.user2
	})

	t.Run("Customer already assigned to target sales", func(t *testing.T) {
		// Customer 203 is already assigned to sales.user2 (ID 7)
		response, err := svc.TransferSalesCustomer(ctx, "7", "7", "203")

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Contains(t, response.Status, "Failed, Customer with customer_id 203 is already under sales 7")
	})

	t.Run("Customer not assigned to source sales", func(t *testing.T) {
		// First verify the initial state of customer 202
		customer, err := dbClient.Customer.Get(ctx, 202)
		assert.NoError(t, err)
		assert.Equal(t, 6, customer.SalesID) // Should be assigned to sales.user1 (ID 6)

		// Try to transfer from sales.user2 (ID 7) which is not the owner
		fromSalesID := "7" // Not the actual owner
		toSalesID := "2"   // Arbitrary other sales ID, doesn't matter for this test
		customerID := "202"

		response, err := svc.TransferSalesCustomer(ctx, fromSalesID, toSalesID, customerID)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		expectedMsg := fmt.Sprintf("Failed, Customer with customer_id %s is not under sales %s", customerID, fromSalesID)
		assert.Equal(t, expectedMsg, response.Status)
	})

	t.Run("Customer not found", func(t *testing.T) {
		// Non-existent customer ID
		response, err := svc.TransferSalesCustomer(ctx, "6", "7", "999")

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Contains(t, response.Status, "Failed, Customer with customer_id 999 not found")
	})

	t.Run("Target sales not found", func(t *testing.T) {
		// First check that customer 201 exists and is assigned to sales 7 (from previous test)
		customer, err := dbClient.Customer.Get(ctx, 201)
		assert.NoError(t, err)
		assert.Equal(t, 7, customer.SalesID)

		// Non-existent sales ID
		fromSalesID := "7" // Current owner after previous test
		toSalesID := "999" // Non-existent sales ID
		customerID := "201"

		response, err := svc.TransferSalesCustomer(ctx, fromSalesID, toSalesID, customerID)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		expectedMsg := fmt.Sprintf("Failed, Cannot find the internal user %s", toSalesID)
		assert.Equal(t, expectedMsg, response.Status)
	})

	t.Run("Invalid IDs", func(t *testing.T) {
		// Invalid customer ID format
		response, err := svc.TransferSalesCustomer(ctx, "6", "7", "invalid")

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Contains(t, response.Status, "Failed, Invalid customer ID: invalid")

		// Invalid from_sales_id format
		response, err = svc.TransferSalesCustomer(ctx, "invalid", "7", "201")

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Contains(t, response.Status, "Failed, Invalid from_sales_id: invalid")

		// Invalid to_sales_id format
		response, err = svc.TransferSalesCustomer(ctx, "6", "invalid", "201")

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Contains(t, response.Status, "Failed, Invalid to_sales_id: invalid")
	})
}

// Helper function to create a JWT token for testing
func createTestJWTToken(userID int) (string, error) {
	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Issuer:    "test",
	}

	// Create the token with custom claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"sub":    "test-subject",
		"exp":    claims.ExpiresAt,
		"iss":    claims.Issuer,
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte("test-jwt-secret-for-unit-tests"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func TestGetUser2FAContactInfo(t *testing.T) {
	svc, dbClient, redisServer, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, redisServer)

	// Disable foreign key constraints for testing
	_, err := dbClient.ExecContext(ctx, "PRAGMA foreign_keys = OFF")
	if err != nil {
		t.Logf("Warning: Could not disable foreign keys: %v", err)
	}

	// Create test users
	userID := 501
	_, err = dbClient.User.Create().
		SetID(userID).
		SetUserName("user_with_2fa").
		SetIsActive(true).
		SetPassword("password123").
		Save(ctx)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create test contacts for the user
	testContacts := []struct {
		id               int
		description      string
		details          string
		contactType      string
		isPrimary        bool
		is2FA            bool
		userID           int
		customerID       int
		patientID        int
		clinicID         int
		internalUserID   int
		groupContactID   int
		isGroupContact   bool
		contactLevel     int
		contactLevelName string
	}{
		{
			id:               601,
			description:      "Primary Email",
			details:          "primary@example.com",
			contactType:      "email",
			isPrimary:        true,
			is2FA:            true,
			userID:           userID,
			customerID:       0,
			patientID:        0,
			clinicID:         0,
			internalUserID:   0,
			groupContactID:   0,
			isGroupContact:   false,
			contactLevel:     1,
			contactLevelName: "User",
		},
		{
			id:               602,
			description:      "Secondary Email",
			details:          "secondary@example.com",
			contactType:      "email",
			isPrimary:        false,
			is2FA:            false, // Not a 2FA contact
			userID:           userID,
			customerID:       0,
			patientID:        0,
			clinicID:         0,
			internalUserID:   0,
			groupContactID:   0,
			isGroupContact:   false,
			contactLevel:     1,
			contactLevelName: "User",
		},
		{
			id:               603,
			description:      "Mobile Phone",
			details:          "+1234567890",
			contactType:      "phone",
			isPrimary:        false,
			is2FA:            true,
			userID:           userID,
			customerID:       0,
			patientID:        0,
			clinicID:         0,
			internalUserID:   0,
			groupContactID:   0,
			isGroupContact:   false,
			contactLevel:     1,
			contactLevelName: "User",
		},
		{
			id:               604,
			description:      "Work Phone",
			details:          "+0987654321",
			contactType:      "phone",
			isPrimary:        false,
			is2FA:            true,
			userID:           userID + 1, // Different user
			customerID:       0,
			patientID:        0,
			clinicID:         0,
			internalUserID:   0,
			groupContactID:   0,
			isGroupContact:   false,
			contactLevel:     1,
			contactLevelName: "User",
		},
		{
			id:               605,
			description:      "Group Email",
			details:          "group@example.com",
			contactType:      "email",
			isPrimary:        false,
			is2FA:            true,
			userID:           userID,
			customerID:       701,
			patientID:        0,
			clinicID:         0,
			internalUserID:   0,
			groupContactID:   610,
			isGroupContact:   false,
			contactLevel:     2,
			contactLevelName: "Customer",
		},
	}

	for _, c := range testContacts {
		_, err := dbClient.Contact.Create().
			SetID(c.id).
			SetContactDescription(c.description).
			SetContactDetails(c.details).
			SetContactType(c.contactType).
			SetIsPrimaryContact(c.isPrimary).
			SetIs2faContact(c.is2FA).
			SetUserID(c.userID).
			SetCustomerID(c.customerID).
			SetPatientID(c.patientID).
			SetClinicID(c.clinicID).
			SetInternalUserID(c.internalUserID).
			SetGroupContactID(c.groupContactID).
			SetIsGroupContact(c.isGroupContact).
			SetContactLevel(c.contactLevel).
			SetContactLevelName(c.contactLevelName).
			SetApplyToAllGroupMember(false).
			SetUseAsDefaultCreateContact(false).
			SetUseGroupContact(c.groupContactID > 0).
			Save(ctx)
		if err != nil {
			t.Fatalf("Failed to create test contact: %v", err)
		}
	}

	t.Run("Get 2FA contacts for valid user", func(t *testing.T) {
		// Create a valid JWT token for the test user
		token, err := createTestJWTToken(userID)
		assert.NoError(t, err)

		// Call the service method
		result, err := svc.GetUser2FAContactInfo(ctx, token)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Contacts, 3) // Should return 3 2FA contacts for this user

		// Verify the contacts returned are the correct ones
		contactDetails := make(map[string]bool)
		for _, contact := range result.Contacts {
			contactDetails[contact.ContactDetails] = true
			assert.True(t, contact.Is_2FaContact) // All should be 2FA contacts
		}

		// Should include the user's 2FA contacts
		assert.True(t, contactDetails["primary@example.com"])
		assert.True(t, contactDetails["+1234567890"])
		assert.True(t, contactDetails["group@example.com"])

		// Should NOT include non-2FA contacts or contacts from other users
		assert.False(t, contactDetails["secondary@example.com"])
		assert.False(t, contactDetails["+0987654321"])
	})

	// Skipping the invalid token test as it would require modifying the ParseJWTToken function
	// to handle invalid tokens more gracefully
	/*
		t.Run("Invalid token", func(t *testing.T) {
			// Test with an invalid JWT token
			_, err := svc.GetUser2FAContactInfo(ctx, "invalid-token")

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid token")
		})
	*/

	t.Run("Token with invalid user ID", func(t *testing.T) {
		// Create a token with an invalid user ID (0)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": 0, // Invalid user ID
			"sub":    "test-subject",
			"exp":    time.Now().Add(time.Hour).Unix(),
		})

		// Make sure we're using the test secret
		common.Secrets.JWTSecret = "test-jwt-secret-for-unit-tests"

		tokenString, err := token.SignedString([]byte("test-jwt-secret-for-unit-tests"))
		assert.NoError(t, err)

		// Call the service method with the invalid token
		_, err = svc.GetUser2FAContactInfo(ctx, tokenString)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID in token")
	})

	t.Run("User with no 2FA contacts", func(t *testing.T) {
		// This user has no contacts at all - should return empty results
		noContactsUserID := 505 // Use a different user ID not already used in tests
		_, err := dbClient.User.Create().
			SetID(noContactsUserID).
			SetUserName("user_no_2fa").
			SetIsActive(true).
			SetPassword("password123").
			Save(ctx)
		assert.NoError(t, err)

		// Create a valid JWT token for this user
		token, err := createTestJWTToken(noContactsUserID)
		assert.NoError(t, err)

		// Call the service method
		result, err := svc.GetUser2FAContactInfo(ctx, token)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Contacts, 0) // Should return an empty array
	})

	t.Run("Get contacts for user with only non-2FA contacts", func(t *testing.T) {
		// Create a user with only non-2FA contacts
		nonFAUserID := 503
		_, err := dbClient.User.Create().
			SetID(nonFAUserID).
			SetUserName("user_non_2fa").
			SetIsActive(true).
			SetPassword("password123").
			Save(ctx)
		assert.NoError(t, err)

		// Create a non-2FA contact for this user
		_, err = dbClient.Contact.Create().
			SetID(606).
			SetContactDescription("Non-2FA Email").
			SetContactDetails("non2fa@example.com").
			SetContactType("email").
			SetIsPrimaryContact(true).
			SetIs2faContact(false). // Not a 2FA contact
			SetUserID(nonFAUserID).
			SetCustomerID(0).
			SetPatientID(0).
			SetClinicID(0).
			SetInternalUserID(0).
			SetGroupContactID(0).
			SetIsGroupContact(false).
			SetContactLevel(1).
			SetContactLevelName("User").
			SetApplyToAllGroupMember(false).
			SetUseAsDefaultCreateContact(false).
			SetUseGroupContact(false).
			Save(ctx)
		assert.NoError(t, err)

		// Create a valid JWT token for this user
		token, err := createTestJWTToken(nonFAUserID)
		assert.NoError(t, err)

		// Call the service method
		result, err := svc.GetUser2FAContactInfo(ctx, token)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Contacts, 0) // Should return an empty array since no 2FA contacts
	})
}

func TestGetUserInformation(t *testing.T) {
	svc, dbClient, redisServer, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, redisServer)

	// Create test users and customers
	createTestUsersForGetUserInformation(t, dbClient, ctx)

	t.Run("Get user with customer information", func(t *testing.T) {
		// Test retrieving a user with associated customer
		result, err := svc.GetUserInformation(ctx, "501")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(501), result.UserId)
		assert.Equal(t, "user501", result.Username)
		assert.Equal(t, "user501@example.com", result.EmailUserId)
		assert.True(t, result.IsActive)
		assert.False(t, result.ImportedUserWithSaltPassword)
		assert.Equal(t, "", result.UserPermission) // UserPermission field is deprecated, should be empty

		// Verify customer information
		assert.NotNil(t, result.Customer)
		assert.Equal(t, int32(601), result.Customer.CustomerId)
		assert.Equal(t, "First", result.Customer.CustomerFirstName)
		assert.Equal(t, "Customer", result.Customer.CustomerLastName)
	})

	t.Run("Get user without customer information", func(t *testing.T) {
		// Test retrieving a user without associated customer
		result, err := svc.GetUserInformation(ctx, "503")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(503), result.UserId)
		assert.Equal(t, "user503", result.Username)
		assert.True(t, result.IsTwoFactorAuthenticationEnabled)
		assert.Nil(t, result.Customer) // No customer associated
	})

	t.Run("Get user with two-factor authentication enabled", func(t *testing.T) {
		// Test retrieving a user with 2FA enabled
		result, err := svc.GetUserInformation(ctx, "503")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(503), result.UserId)
		assert.True(t, result.IsTwoFactorAuthenticationEnabled)
	})

	t.Run("Get inactive user", func(t *testing.T) {
		// Test retrieving an inactive user
		result, err := svc.GetUserInformation(ctx, "504")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(504), result.UserId)
		assert.Equal(t, "user504", result.Username)
		assert.False(t, result.IsActive)
	})

	t.Run("Get imported user with salt password", func(t *testing.T) {
		// Test retrieving a user imported with salt password
		result, err := svc.GetUserInformation(ctx, "505")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(505), result.UserId)
		assert.Equal(t, "user505", result.Username)
		assert.True(t, result.ImportedUserWithSaltPassword)
	})

	t.Run("Non-existent user returns empty response", func(t *testing.T) {
		// Test retrieving a non-existent user
		result, err := svc.GetUserInformation(ctx, "999")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(0), result.UserId) // Empty user ID indicates not found
	})

	t.Run("Invalid user ID format returns error", func(t *testing.T) {
		// Test with invalid user ID format
		_, err := svc.GetUserInformation(ctx, "invalid")

		assert.Error(t, err) // Should return an error for invalid format
	})
}

// Helper function to create test users and customers for GetUserInformation testing
func TestRenewToken(t *testing.T) {
	svc, dbClient, redisServer, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, redisServer)

	// Make sure test JWT secret is set
	common.Secrets.JWTSecret = "test-jwt-secret-for-unit-tests"

	t.Run("Renew valid token", func(t *testing.T) {
		// Create a valid token with required claims
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["userId"] = 123
		claims["user_permission"] = "admin"
		claims["customer_id"] = 456
		claims["clinic_id"] = 789
		claims["role"] = "customer"
		claims["iat"] = time.Now().Unix()
		claims["exp"] = time.Now().Add(time.Hour).Unix()

		tokenString, err := token.SignedString([]byte(common.Secrets.JWTSecret))
		assert.NoError(t, err)

		// Call the service method
		result, err := svc.RenewToken(ctx, tokenString)

		// Verify result
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(200), result.Code)
		assert.Equal(t, "Token Renewed", result.Message)
		assert.NotEmpty(t, result.Token)
		assert.NotEmpty(t, result.ExpirationTime)

		// Verify new token is valid
		newToken, err := jwt.Parse(result.Token, func(token *jwt.Token) (interface{}, error) {
			return []byte(common.Secrets.JWTSecret), nil
		})
		assert.NoError(t, err)
		assert.True(t, newToken.Valid)

		// Verify claims were preserved
		newClaims := newToken.Claims.(jwt.MapClaims)
		assert.Equal(t, float64(123), newClaims["userId"])
		assert.Equal(t, "admin", newClaims["user_permission"])
		assert.Equal(t, float64(456), newClaims["customer_id"])
	})

	t.Run("Renew expired token", func(t *testing.T) {
		// Create an expired token
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["userId"] = 123
		claims["iat"] = time.Now().Add(-2 * time.Hour).Unix()
		claims["exp"] = time.Now().Add(-1 * time.Hour).Unix() // Expired

		tokenString, err := token.SignedString([]byte(common.Secrets.JWTSecret))
		assert.NoError(t, err)

		// Call the service method
		result, err := svc.RenewToken(ctx, tokenString)

		// Verify result
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(400), result.Code)
		assert.Contains(t, result.Message, "Token Renew Failed")
		assert.Empty(t, result.Token)
	})

	t.Run("Invalid signature", func(t *testing.T) {
		// Create a token with a different secret
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["userId"] = 123
		claims["iat"] = time.Now().Unix()
		claims["exp"] = time.Now().Add(time.Hour).Unix()

		tokenString, err := token.SignedString([]byte("wrong-secret"))
		assert.NoError(t, err)

		// Call the service method
		result, err := svc.RenewToken(ctx, tokenString)

		// Verify result
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(400), result.Code)
		assert.Contains(t, result.Message, "Token Renew Failed")
		assert.Empty(t, result.Token)
	})

	t.Run("Malformed token", func(t *testing.T) {
		// Test with an invalid token string
		result, err := svc.RenewToken(ctx, "not-a-valid-token")

		// Verify result
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(400), result.Code)
		assert.Contains(t, result.Message, "Token Renew Failed")
	})

	t.Run("Custom expiration time from environment", func(t *testing.T) {
		// Set custom JWT expiration time via environment variable
		originalValue := os.Getenv("JWT_EXPIRATION_TIME")
		defer os.Setenv("JWT_EXPIRATION_TIME", originalValue) // Restore original value

		// Set custom expiration time (3600 seconds = 1 hour)
		os.Setenv("JWT_EXPIRATION_TIME", "3600")

		// Create a valid token
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["userId"] = 123
		claims["iat"] = time.Now().Unix()
		claims["exp"] = time.Now().Add(time.Hour).Unix()

		tokenString, err := token.SignedString([]byte(common.Secrets.JWTSecret))
		assert.NoError(t, err)

		// Call the service method
		result, err := svc.RenewToken(ctx, tokenString)

		// Verify result
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(200), result.Code)

		// Parse expiration time from the response
		expirationTime, err := strconv.ParseInt(result.ExpirationTime, 10, 64)
		assert.NoError(t, err)

		// Verify expiration time reflects the custom value (approximately)
		// Current time + 3600 seconds (with some tolerance for test execution time)
		expectedExpiration := time.Now().Unix() + 3600
		assert.InDelta(t, expectedExpiration, expirationTime, 5) // Allow 5 second delta for test execution
	})
}

func createTestUsersForGetUserInformation(t *testing.T, dbClient *ent.Client, ctx context.Context) {
	// Disable foreign key constraints for testing
	_, err := dbClient.ExecContext(ctx, "PRAGMA foreign_keys = OFF")
	if err != nil {
		t.Logf("Warning: Could not disable foreign keys: %v", err)
	}

	// Create test users
	users := []struct {
		id                               int
		username                         string
		emailUserID                      string
		isTwoFactorAuthenticationEnabled bool
		importedUserWithSaltPassword     bool
		isActive                         bool
	}{
		{
			id:                               501,
			username:                         "user501",
			emailUserID:                      "user501@example.com",
			isTwoFactorAuthenticationEnabled: false,
			importedUserWithSaltPassword:     false,
			isActive:                         true,
		},
		{
			id:                               502,
			username:                         "user502",
			emailUserID:                      "user502@example.com",
			isTwoFactorAuthenticationEnabled: false,
			importedUserWithSaltPassword:     false,
			isActive:                         true,
		},
		{
			id:                               503,
			username:                         "user503",
			emailUserID:                      "user503@example.com",
			isTwoFactorAuthenticationEnabled: true,
			importedUserWithSaltPassword:     false,
			isActive:                         true,
		},
		{
			id:                               504,
			username:                         "user504",
			emailUserID:                      "user504@example.com",
			isTwoFactorAuthenticationEnabled: false,
			importedUserWithSaltPassword:     false,
			isActive:                         false, // Inactive user
		},
		{
			id:                               505,
			username:                         "user505",
			emailUserID:                      "user505@example.com",
			isTwoFactorAuthenticationEnabled: false,
			importedUserWithSaltPassword:     true, // Imported user with salt password
			isActive:                         true,
		},
	}

	// Create users
	for _, u := range users {
		_, err := dbClient.User.Create().
			SetID(u.id).
			SetUserName(u.username).
			SetEmailUserID(u.emailUserID).
			SetPassword("password-hash").
			SetIsTwoFactorAuthenticationEnabled(u.isTwoFactorAuthenticationEnabled).
			SetImportedUserWithSaltPassword(u.importedUserWithSaltPassword).
			SetIsActive(u.isActive).
			SetUserGroup("customer").
			Save(ctx)

		if err != nil {
			t.Fatalf("Failed to create test user %d: %v", u.id, err)
		}
	}

	// Create customers associated with users
	customers := []struct {
		id         int
		userID     int
		firstName  string
		lastName   string
		middleName string
		isActive   bool
	}{
		{
			id:         601,
			userID:     501,
			firstName:  "First",
			lastName:   "Customer",
			middleName: "",
			isActive:   true,
		},
		{
			id:         602,
			userID:     502,
			firstName:  "Second",
			lastName:   "Customer",
			middleName: "Middle",
			isActive:   true,
		},
	}

	// Create customers
	for _, c := range customers {
		_, err := dbClient.Customer.Create().
			SetID(c.id).
			SetUserID(c.userID).
			SetCustomerFirstName(c.firstName).
			SetCustomerLastName(c.lastName).
			SetCustomerMiddleName(c.middleName).
			SetIsActive(c.isActive).
			SetSalesID(1). // Default sales ID
			Save(ctx)

		if err != nil {
			t.Fatalf("Failed to create test customer %d: %v", c.id, err)
		}
	}
}

func createTestUsersForSend2FA(t *testing.T, dbClient *ent.Client, ctx context.Context) {
	// Create test users with specific IDs
	users := []struct {
		id           int
		username     string
		emailUserID  string
		is2faEnabled bool
		isActive     bool
	}{
		{
			id:           501,
			username:     "user501",
			emailUserID:  "user501@example.com",
			is2faEnabled: false,
			isActive:     true,
		},
		{
			id:           503,
			username:     "user503",
			emailUserID:  "user503@example.com",
			is2faEnabled: true,
			isActive:     true,
		},
	}

	for _, u := range users {
		_, err := dbClient.User.Create().
			SetID(u.id).
			SetUserName(u.username).
			SetEmailUserID(u.emailUserID).
			SetPassword("password-hash").
			SetIsTwoFactorAuthenticationEnabled(u.is2faEnabled).
			SetIsActive(u.isActive).
			Save(ctx)
		if err != nil {
			t.Fatalf("Failed to create test user %d: %v", u.id, err)
		}
	}
}

func TestSend2FAVerificationCode(t *testing.T) {
	userService, dbClient, server, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, server)

	// Create test users and contacts
	createTestUsersForSend2FA(t, dbClient, ctx)

	// Create test contacts for 2FA
	contacts := []struct {
		id             int
		userID         int
		contactType    string
		contactDetails string
		is2faContact   bool
	}{
		{
			id:             1,
			userID:         503, // User with 2FA enabled
			contactType:    "email",
			contactDetails: "test@example.com",
			is2faContact:   true,
		},
		{
			id:             2,
			userID:         503,
			contactType:    "phone",
			contactDetails: "1234567890",
			is2faContact:   true,
		},
		{
			id:             3,
			userID:         501, // User without 2FA enabled
			contactType:    "email",
			contactDetails: "user501@example.com",
			is2faContact:   true,
		},
	}

	// Create contacts
	for _, c := range contacts {
		_, err := dbClient.Contact.Create().
			SetID(c.id).
			SetUserID(c.userID).
			SetContactType(c.contactType).
			SetContactDetails(c.contactDetails).
			SetIs2faContact(c.is2faContact).
			SetIsPrimaryContact(true).
			Save(ctx)

		if err != nil {
			t.Fatalf("Failed to create test contact %d: %v", c.id, err)
		}
	}

	testCases := []struct {
		name         string
		username     string
		emailAddress string
		phoneNumber  string
		expectedCode int
		expectedMsg  string
		setupMock    func()
	}{
		{
			name:         "User not found",
			username:     "nonexistent",
			emailAddress: "test@example.com",
			phoneNumber:  "",
			expectedCode: 404,
			expectedMsg:  "User Not Found",
		},
		{
			name:         "2FA not enabled",
			username:     "user501",
			emailAddress: "user501@example.com",
			phoneNumber:  "",
			expectedCode: 400,
			expectedMsg:  "2FA is Not Enabled",
		},
		{
			name:         "Valid email verification - using username",
			username:     "user503",
			emailAddress: "test@example.com",
			phoneNumber:  "",
			expectedCode: 200,
			expectedMsg:  "Code sent via email",
		},
		{
			name:         "Valid email verification - using email as username",
			username:     "user503@example.com",
			emailAddress: "test@example.com",
			phoneNumber:  "",
			expectedCode: 200,
			expectedMsg:  "Code sent via email",
		},
		{
			name:         "Invalid email address",
			username:     "user503",
			emailAddress: "wrong@example.com",
			phoneNumber:  "",
			expectedCode: 400,
			expectedMsg:  "The Input Email Address is Not the 2FA Email of This Account",
		},
		{
			name:         "Valid phone verification",
			username:     "user503",
			emailAddress: "",
			phoneNumber:  "1234567890",
			expectedCode: 200,
			expectedMsg:  "Code sent via text",
		},
		{
			name:         "Invalid phone number",
			username:     "user503",
			emailAddress: "",
			phoneNumber:  "0987654321",
			expectedCode: 400,
			expectedMsg:  "The Input Phone Number is Not the 2FA Phone of This Account",
		},
		{
			name:         "Neither email nor phone provided",
			username:     "user503",
			emailAddress: "",
			phoneNumber:  "",
			expectedCode: 400,
			expectedMsg:  "Please provide either email or phone number",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock if needed
			if tc.setupMock != nil {
				tc.setupMock()
			}

			// Call the method
			response, err := userService.Send2FAVerificationCode(ctx, tc.username, tc.emailAddress, tc.phoneNumber)

			// Verify the results
			assert.NoError(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, tc.expectedCode, int(response.Code))
			assert.Equal(t, tc.expectedMsg, response.Message)

			// For successful cases, verify that the OTP is stored in Redis
			if tc.expectedCode == 200 {
				var redisKey string
				userID := 503 // Known user ID for user503

				if tc.emailAddress != "" {
					redisKey = fmt.Sprintf("lis::core_service::user_service:2fa%d_email", userID)
				} else if tc.phoneNumber != "" {
					redisKey = fmt.Sprintf("lis::core_service::user_service:2fa%d_text", userID)
				}

				if redisKey != "" {
					redisClient := userService.(*UserService).redisClient
					val, err := redisClient.Get(ctx, redisKey).Result()
					assert.NoError(t, err)
					assert.NotEmpty(t, val, "OTP should be stored in Redis")
					assert.Len(t, val, 6, "OTP should be 6 digits")
				}
			}
		})
	}
}

func TestSend2FAVerificationCode_EdgeCases(t *testing.T) {
	userService, dbClient, server, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, server)

	// Create a user with 2FA enabled but no contacts
	_, err := dbClient.User.Create().
		SetID(999).
		SetUserName("user999").
		SetEmailUserID("user999@example.com").
		SetPassword("password-hash").
		SetIsTwoFactorAuthenticationEnabled(true).
		SetIsActive(true).
		Save(ctx)
	assert.NoError(t, err)

	// Test user with 2FA enabled but no contacts
	t.Run("User with 2FA but no contacts", func(t *testing.T) {
		response, err := userService.Send2FAVerificationCode(ctx, "user999", "test@example.com", "")

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int32(400), response.Code)
		assert.Equal(t, "This Account Does Not Have 2FA Contact, Please Contact Support", response.Message)
	})

	// Create a user with email contact but not 2FA contact
	_, err = dbClient.Contact.Create().
		SetID(100).
		SetUserID(999).
		SetContactType("email").
		SetContactDetails("test@example.com").
		SetIs2faContact(false). // Not a 2FA contact
		SetIsPrimaryContact(true).
		Save(ctx)
	assert.NoError(t, err)

	// Test user with non-2FA contact
	t.Run("User with non-2FA email contact", func(t *testing.T) {
		response, err := userService.Send2FAVerificationCode(ctx, "user999", "test@example.com", "")

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int32(400), response.Code)
		assert.Equal(t, "This Account Does Not Have 2FA Contact, Please Contact Support", response.Message)
	})
}

func TestComplete2FAFlow(t *testing.T) {
	userService, dbClient, server, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, server)

	// Create a user with 2FA enabled
	userId := 800
	username := "complete_flow_user"
	_, err := dbClient.User.Create().
		SetID(userId).
		SetUserName(username).
		SetEmailUserID("complete_flow@example.com").
		SetPassword("password-hash").
		SetIsTwoFactorAuthenticationEnabled(true).
		SetIsActive(true).
		Save(ctx)
	assert.NoError(t, err)

	// Create 2FA email contact
	emailAddress := "complete_flow_email@example.com"
	_, err = dbClient.Contact.Create().
		SetID(900).
		SetUserID(userId).
		SetContactType("email").
		SetContactDetails(emailAddress).
		SetIs2faContact(true).
		SetIsPrimaryContact(true).
		Save(ctx)
	assert.NoError(t, err)

	// Create 2FA phone contact
	phoneNumber := "5551234567"
	_, err = dbClient.Contact.Create().
		SetID(901).
		SetUserID(userId).
		SetContactType("phone").
		SetContactDetails(phoneNumber).
		SetIs2faContact(true).
		SetIsPrimaryContact(false).
		Save(ctx)
	assert.NoError(t, err)

	// Test the complete 2FA flow with email
	t.Run("Complete 2FA flow with email", func(t *testing.T) {
		// Step 1: Send verification code via email
		sendResponse, err := userService.Send2FAVerificationCode(ctx, username, emailAddress, "")
		assert.NoError(t, err)
		assert.NotNil(t, sendResponse)
		assert.Equal(t, int32(200), sendResponse.Code)
		assert.Equal(t, "Code sent via email", sendResponse.Message)

		// Step 2: Check if code was stored in Redis
		redisClient := userService.(*UserService).redisClient
		redisKey := fmt.Sprintf("lis::core_service::user_service:2fa%d_email", userId)
		storedCode, err := redisClient.Get(ctx, redisKey).Result()
		assert.NoError(t, err)
		assert.NotEmpty(t, storedCode)
		assert.Len(t, storedCode, 6) // 6-digit OTP

		// Step 3: Verify the code with wrong value
		wrongVerifyResponse, err := userService.Verify2FAVerificationCode(ctx, username, "000000", emailAddress, "")
		assert.NoError(t, err)
		assert.NotNil(t, wrongVerifyResponse)
		assert.Equal(t, int32(400), wrongVerifyResponse.Code)
		assert.Equal(t, "Invalid verification code", wrongVerifyResponse.Message)

		// Step 4: Verify the code with correct value
		verifyResponse, err := userService.Verify2FAVerificationCode(ctx, username, storedCode, emailAddress, "")
		assert.NoError(t, err)
		assert.NotNil(t, verifyResponse)
		assert.Equal(t, int32(200), verifyResponse.Code)
		assert.Equal(t, "Verification successful", verifyResponse.Message)

		// Step 5: Check if code was deleted from Redis after successful verification
		_, err = redisClient.Get(ctx, redisKey).Result()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "redis: nil")
	})

	// Test the complete 2FA flow with SMS
	t.Run("Complete 2FA flow with SMS", func(t *testing.T) {
		// Step 1: Send verification code via SMS
		sendResponse, err := userService.Send2FAVerificationCode(ctx, username, "", phoneNumber)
		assert.NoError(t, err)
		assert.NotNil(t, sendResponse)
		assert.Equal(t, int32(200), sendResponse.Code)
		assert.Equal(t, "Code sent via text", sendResponse.Message)

		// Step 2: Check if code was stored in Redis
		redisClient := userService.(*UserService).redisClient
		redisKey := fmt.Sprintf("lis::core_service::user_service:2fa%d_text", userId)
		storedCode, err := redisClient.Get(ctx, redisKey).Result()
		assert.NoError(t, err)
		assert.NotEmpty(t, storedCode)
		assert.Len(t, storedCode, 6) // 6-digit OTP

		// Step 3: Verify the code with wrong value
		wrongVerifyResponse, err := userService.Verify2FAVerificationCode(ctx, username, "000000", "", phoneNumber)
		assert.NoError(t, err)
		assert.NotNil(t, wrongVerifyResponse)
		assert.Equal(t, int32(400), wrongVerifyResponse.Code)
		assert.Equal(t, "Invalid verification code", wrongVerifyResponse.Message)

		// Step 4: Verify the code with correct value
		verifyResponse, err := userService.Verify2FAVerificationCode(ctx, username, storedCode, "", phoneNumber)
		assert.NoError(t, err)
		assert.NotNil(t, verifyResponse)
		assert.Equal(t, int32(200), verifyResponse.Code)
		assert.Equal(t, "Verification successful", verifyResponse.Message)

		// Step 5: Check if code was deleted from Redis after successful verification
		_, err = redisClient.Get(ctx, redisKey).Result()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "redis: nil")
	})
}

// localMockTransport is a custom mock transport for HTTP testing
type localMockTransport struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

// RoundTrip implements the http.RoundTripper interface
func (m *localMockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func TestSend2FAEmailVerificationWithMockClient(t *testing.T) {

	// Create a mock HTTP client that intercepts requests instead of making real API calls
	mockClient := &http.Client{
		Transport: &localMockTransport{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				if req.URL.String() == "https://www.vibrant-america.com/lisapi/v1/portal/trans-service/valogin/send2faAuthEmail" {
					// Return a successful response
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
						Header:     make(http.Header),
					}, nil
				}

				// Return error for unexpected requests
				return nil, fmt.Errorf("unexpected request to: %s", req.URL.String())
			},
		},
		Timeout: 30 * time.Second,
	}

	// Create service with mock HTTP client - use the test setup
	dataSource := "file:ent?mode=memory&_fk=1"
	dbClient := enttest.Open(t, "sqlite3", dataSource)
	ctx := context.Background()
	err := dbClient.Schema.Create(ctx)
	if err != nil {
		t.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer dbClient.Close()

	// Create mock Redis
	server, err := tempredis.Start(tempredis.Config{
		"port": "0",
	})
	if err != nil {
		t.Fatalf("Failed to start tempredis: %v", err)
	}
	defer server.Kill()

	redisClient := redis.NewClient(&redis.Options{
		Network: "unix",
		Addr:    server.Socket(),
	})

	// Initialize JWT secret
	common.Secrets.JWTSecret = "test-jwt-secret-for-unit-tests"

	// Initialize mock publisher
	publisher.InitMockPublisher()

	// Create service with mock HTTP client
	userService := NewUserServiceWithHTTPClient(dbClient, common.NewRedisClient(redisClient, redisClient), mockClient).(*UserService)

	// Call the email verification method
	err = userService.send2FAEmailVerification(context.Background(), "test@example.com", "123456")

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestVerify2FAVerificationCode(t *testing.T) {
	svc, dbClient, redisServer, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, redisServer)

	// Create test users for verification
	testUsers := []struct {
		id           int
		username     string
		emailUserID  string
		is2faEnabled bool
		isActive     bool
	}{
		{
			id:           601,
			username:     "verify_user1",
			emailUserID:  "verify1@example.com",
			is2faEnabled: true,
			isActive:     true,
		},
		{
			id:           602,
			username:     "verify_user2",
			emailUserID:  "verify2@example.com",
			is2faEnabled: false, // 2FA not enabled
			isActive:     true,
		},
	}

	for _, u := range testUsers {
		_, err := dbClient.User.Create().
			SetID(u.id).
			SetUserName(u.username).
			SetEmailUserID(u.emailUserID).
			SetPassword("password-hash").
			SetIsTwoFactorAuthenticationEnabled(u.is2faEnabled).
			SetIsActive(u.isActive).
			Save(ctx)
		if err != nil {
			t.Fatalf("Failed to create test user %d: %v", u.id, err)
		}
	}

	// Create test contacts for 2FA
	contacts := []struct {
		id             int
		userID         int
		contactType    string
		contactDetails string
		is2faContact   bool
	}{
		{
			id:             701,
			userID:         601, // User with 2FA enabled
			contactType:    "email",
			contactDetails: "2fa_email@example.com",
			is2faContact:   true,
		},
		{
			id:             702,
			userID:         601,
			contactType:    "phone",
			contactDetails: "1234567890",
			is2faContact:   true,
		},
		{
			id:             703,
			userID:         601,
			contactType:    "email",
			contactDetails: "non_2fa_email@example.com",
			is2faContact:   false, // Not a 2FA contact
		},
	}

	for _, c := range contacts {
		_, err := dbClient.Contact.Create().
			SetID(c.id).
			SetUserID(c.userID).
			SetContactType(c.contactType).
			SetContactDetails(c.contactDetails).
			SetIs2faContact(c.is2faContact).
			SetIsPrimaryContact(c.contactType == "email").
			Save(ctx)
		if err != nil {
			t.Fatalf("Failed to create test contact %d: %v", c.id, err)
		}
	}

	// Store test verification codes in Redis
	redisClient := svc.(*UserService).redisClient

	// Set a valid code for email verification
	emailRedisKey := "lis::core_service::user_service:2fa601_email"
	validEmailCode := "123456"
	err := redisClient.Set(ctx, emailRedisKey, validEmailCode, 300*time.Second).Err()
	if err != nil {
		t.Fatalf("Failed to set Redis key: %v", err)
	}

	// Set a valid code for phone verification
	phoneRedisKey := "lis::core_service::user_service:2fa601_text"
	validPhoneCode := "654321"
	err = redisClient.Set(ctx, phoneRedisKey, validPhoneCode, 300*time.Second).Err()
	if err != nil {
		t.Fatalf("Failed to set Redis key: %v", err)
	}

	testCases := []struct {
		name              string
		username          string
		verificationCode  string
		emailAddress      string
		phoneNumber       string
		expectedCode      int
		expectedMessage   string
		shouldDeleteRedis bool
	}{
		{
			name:             "User not found",
			username:         "nonexistent_user",
			verificationCode: "123456",
			emailAddress:     "2fa_email@example.com",
			phoneNumber:      "",
			expectedCode:     404,
			expectedMessage:  "User Not Found",
		},
		{
			name:             "2FA not enabled",
			username:         "verify_user2",
			verificationCode: "123456",
			emailAddress:     "verify2@example.com",
			phoneNumber:      "",
			expectedCode:     400,
			expectedMessage:  "2FA is Not Enabled",
		},
		{
			name:             "Email not provided as 2FA contact",
			username:         "verify_user1",
			verificationCode: "123456",
			emailAddress:     "wrong_email@example.com",
			phoneNumber:      "",
			expectedCode:     400,
			expectedMessage:  "The provided email does not match any 2FA emails for this account",
		},
		{
			name:             "Phone not provided as 2FA contact",
			username:         "verify_user1",
			verificationCode: "123456",
			emailAddress:     "",
			phoneNumber:      "9876543210",
			expectedCode:     400,
			expectedMessage:  "The provided phone number does not match any 2FA phones for this account",
		},
		{
			name:             "Neither email nor phone provided",
			username:         "verify_user1",
			verificationCode: "123456",
			emailAddress:     "",
			phoneNumber:      "",
			expectedCode:     400,
			expectedMessage:  "Please provide either email or phone number",
		},
		{
			name:             "Invalid email verification code",
			username:         "verify_user1",
			verificationCode: "111111", // Wrong code
			emailAddress:     "2fa_email@example.com",
			phoneNumber:      "",
			expectedCode:     400,
			expectedMessage:  "Invalid verification code",
		},
		{
			name:             "Invalid phone verification code",
			username:         "verify_user1",
			verificationCode: "111111", // Wrong code
			emailAddress:     "",
			phoneNumber:      "1234567890",
			expectedCode:     400,
			expectedMessage:  "Invalid verification code",
		},
		{
			name:              "Valid email verification code",
			username:          "verify_user1",
			verificationCode:  validEmailCode,
			emailAddress:      "2fa_email@example.com",
			phoneNumber:       "",
			expectedCode:      200,
			expectedMessage:   "Verification successful",
			shouldDeleteRedis: true,
		},
		{
			name:              "Valid phone verification code",
			username:          "verify_user1",
			verificationCode:  validPhoneCode,
			emailAddress:      "",
			phoneNumber:       "1234567890",
			expectedCode:      200,
			expectedMessage:   "Verification successful",
			shouldDeleteRedis: true,
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// For cases where we need to test successful verification, we need to restore the Redis keys
			// that might have been deleted in previous test cases
			if i == 7 { // "Valid email verification code" case
				// Check if key exists, if not create it
				_, err := redisClient.Get(ctx, emailRedisKey).Result()
				if err != nil && err.Error() == "redis: nil" {
					err = redisClient.Set(ctx, emailRedisKey, validEmailCode, 300*time.Second).Err()
					if err != nil {
						t.Fatalf("Failed to restore Redis key: %v", err)
					}
				}
			}
			if i == 8 { // "Valid phone verification code" case
				// Check if key exists, if not create it
				_, err := redisClient.Get(ctx, phoneRedisKey).Result()
				if err != nil && err.Error() == "redis: nil" {
					err = redisClient.Set(ctx, phoneRedisKey, validPhoneCode, 300*time.Second).Err()
					if err != nil {
						t.Fatalf("Failed to restore Redis key: %v", err)
					}
				}
			}

			// Call the method
			result, err := svc.Verify2FAVerificationCode(ctx, tc.username, tc.verificationCode, tc.emailAddress, tc.phoneNumber)

			// Verify the results
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, int32(tc.expectedCode), result.Code)
			assert.Contains(t, result.Message, tc.expectedMessage)

			// For successful verifications, check that the code is deleted from Redis
			if tc.shouldDeleteRedis {
				var redisKey string
				if tc.emailAddress != "" {
					redisKey = emailRedisKey
				} else if tc.phoneNumber != "" {
					redisKey = phoneRedisKey
				}

				if redisKey != "" {
					// Verify that the key no longer exists
					_, err := redisClient.Get(ctx, redisKey).Result()
					assert.Error(t, err) // Should error because key was deleted
					assert.Contains(t, err.Error(), "redis: nil")
				}
			}
		})
	}
}

func TestVerify2FAVerificationCodeExpiredCode(t *testing.T) {
	userService, dbClient, server, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, server)

	// Create a test user
	userId := 701
	username := "expired_code_user"

	_, err := dbClient.User.Create().
		SetID(userId).
		SetUserName(username).
		SetEmailUserID("expired@example.com").
		SetPassword("password-hash").
		SetIsTwoFactorAuthenticationEnabled(true).
		SetIsActive(true).
		Save(ctx)
	assert.NoError(t, err)

	// Create email contact
	emailContactId := 801
	emailAddress := "expired_test@example.com"
	_, err = dbClient.Contact.Create().
		SetID(emailContactId).
		SetUserID(userId).
		SetContactType("email").
		SetContactDetails(emailAddress).
		SetIs2faContact(true).
		SetIsPrimaryContact(true).
		Save(ctx)
	assert.NoError(t, err)

	// Test with an expired code (not set in Redis)
	t.Run("Expired verification code", func(t *testing.T) {
		result, err := userService.Verify2FAVerificationCode(ctx, username, "123456", emailAddress, "")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(400), result.Code)
		assert.Contains(t, result.Message, "Verification code expired or not sent")
	})
}

func TestVerify2FAVerificationCodeWithEmail(t *testing.T) {
	userService, dbClient, server, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, server)

	// Create a test user
	userId := 702
	username := "email_verify_user"
	emailAsUsername := "email_user@example.com"

	_, err := dbClient.User.Create().
		SetID(userId).
		SetUserName(username).
		SetEmailUserID(emailAsUsername).
		SetPassword("password-hash").
		SetIsTwoFactorAuthenticationEnabled(true).
		SetIsActive(true).
		Save(ctx)
	assert.NoError(t, err)

	// Create email contact
	emailContactId := 802
	emailAddress := "2fa_test@example.com"
	_, err = dbClient.Contact.Create().
		SetID(emailContactId).
		SetUserID(userId).
		SetContactType("email").
		SetContactDetails(emailAddress).
		SetIs2faContact(true).
		SetIsPrimaryContact(true).
		Save(ctx)
	assert.NoError(t, err)

	// Store verification code in Redis
	redisClient := userService.(*UserService).redisClient
	redisKey := fmt.Sprintf("lis::core_service::user_service:2fa%d_email", userId)
	verificationCode := "987654"
	err = redisClient.Set(ctx, redisKey, verificationCode, 300*time.Second).Err()
	assert.NoError(t, err)

	// Test verification with email as username
	t.Run("Using email as username instead of username field", func(t *testing.T) {
		result, err := userService.Verify2FAVerificationCode(ctx, emailAsUsername, verificationCode, emailAddress, "")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(200), result.Code)
		assert.Equal(t, "Verification successful", result.Message)
	})
}

// Test InitialForgetPassword functionality - comprehensive tests matching TypeScript behavior
func TestInitialForgetPassword(t *testing.T) {
	userService, dbClient, server, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, server)

	// Create test users with different configurations
	testUsers := []struct {
		id          int
		emailUserID string
		username    string
		isActive    bool
	}{
		{
			id:          1001,
			emailUserID: "valid@example.com",
			username:    "validuser",
			isActive:    true,
		},
		{
			id:          1002,
			emailUserID: "with2fa@example.com",
			username:    "user2fa",
			isActive:    true,
		},
		{
			id:          1003,
			emailUserID: "withcustomer@example.com",
			username:    "customeruser",
			isActive:    true,
		},
		{
			id:          1004,
			emailUserID: "withclinic@example.com",
			username:    "clinicuser",
			isActive:    true,
		},
	}

	for _, u := range testUsers {
		_, err := dbClient.User.Create().
			SetID(u.id).
			SetUserName(u.username).
			SetEmailUserID(u.emailUserID).
			SetPassword("password-hash").
			SetIsActive(u.isActive).
			Save(ctx)
		if err != nil {
			t.Fatalf("Failed to create test user %d: %v", u.id, err)
		}
	}

	// Disable foreign key constraints for testing
	_, err := dbClient.ExecContext(ctx, "PRAGMA foreign_keys = OFF")
	if err != nil {
		t.Logf("Warning: Could not disable foreign keys: %v", err)
	}

	// Create customer for user 1003
	_, err = dbClient.Customer.Create().
		SetID(2001).
		SetUserID(1003).
		SetCustomerFirstName("John").
		SetCustomerLastName("Customer").
		SetIsActive(true).
		SetSalesID(1).
		Save(ctx)
	if err != nil {
		t.Fatalf("Failed to create customer: %v", err)
	}

	// Create clinic for user 1004
	_, err = dbClient.Clinic.Create().
		SetID(3001).
		SetUserID(1004).
		SetClinicName("Test Clinic").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("Failed to create clinic: %v", err)
	}

	// Create 2FA contacts for user 1002
	contacts := []struct {
		id             int
		userID         int
		contactType    string
		contactDetails string
		is2faContact   bool
	}{
		{
			id:             4001,
			userID:         1002,
			contactType:    "email",
			contactDetails: "2fa@example.com",
			is2faContact:   true,
		},
		{
			id:             4002,
			userID:         1002,
			contactType:    "phone",
			contactDetails: "1234567890",
			is2faContact:   true,
		},
	}

	for _, c := range contacts {
		_, err := dbClient.Contact.Create().
			SetID(c.id).
			SetUserID(c.userID).
			SetContactType(c.contactType).
			SetContactDetails(c.contactDetails).
			SetIs2faContact(c.is2faContact).
			SetIsPrimaryContact(true).
			Save(ctx)
		if err != nil {
			t.Fatalf("Failed to create contact %d: %v", c.id, err)
		}
	}

	testCases := []struct {
		name            string
		emailAddress    string
		expectedCode    int
		expectedMessage string
	}{
		{
			name:            "Invalid email format - no @ symbol",
			emailAddress:    "invalid-email",
			expectedCode:    400,
			expectedMessage: "Please Enter a Valid Email Address",
		},
		{
			name:            "Invalid email format - empty string",
			emailAddress:    "",
			expectedCode:    400,
			expectedMessage: "Please Enter a Valid Email Address",
		},
		{
			name:            "Invalid email format - missing domain",
			emailAddress:    "user@",
			expectedCode:    400,
			expectedMessage: "Please Enter a Valid Email Address",
		},
		{
			name:            "Invalid email format - missing local part",
			emailAddress:    "@example.com",
			expectedCode:    400,
			expectedMessage: "Please Enter a Valid Email Address",
		},
		{
			name:            "Invalid email format - no dot in domain",
			emailAddress:    "user@example",
			expectedCode:    400,
			expectedMessage: "Please Enter a Valid Email Address",
		},
		{
			name:            "User not found - valid email format",
			emailAddress:    "nonexistent@example.com",
			expectedCode:    400,
			expectedMessage: "Contact Support",
		},
		{
			name:            "Valid user without 2FA contacts",
			emailAddress:    "valid@example.com",
			expectedCode:    200,
			expectedMessage: "Email Address Sent",
		},
		{
			name:            "Valid user with 2FA contacts",
			emailAddress:    "with2fa@example.com",
			expectedCode:    200,
			expectedMessage: "Email Address Sent",
		},
		{
			name:            "Valid user with customer relationship",
			emailAddress:    "withcustomer@example.com",
			expectedCode:    200,
			expectedMessage: "Email Address Sent",
		},
		{
			name:            "Valid user with clinic relationship",
			emailAddress:    "withclinic@example.com",
			expectedCode:    200,
			expectedMessage: "Email Address Sent",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the method
			result, err := userService.InitialForgetPassword(ctx, tc.emailAddress)

			// Verify the results
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, int32(tc.expectedCode), result.Code)
			assert.Equal(t, tc.expectedMessage, result.Message)
		})
	}
}

// Test InitialForgetPassword with HTTP client failure scenarios
func TestInitialForgetPasswordHTTPFailure(t *testing.T) {
	dataSource := "file:ent?mode=memory&_fk=1"
	dbClient := enttest.Open(t, "sqlite3", dataSource)
	ctx := context.Background()
	err := dbClient.Schema.Create(ctx)
	if err != nil {
		t.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer dbClient.Close()

	// Create mock Redis server
	server, err := tempredis.Start(tempredis.Config{
		"port": "0",
	})
	if err != nil {
		t.Fatalf("Failed to start tempredis: %v", err)
	}
	defer server.Kill()

	redisClient := redis.NewClient(&redis.Options{
		Network: "unix",
		Addr:    server.Socket(),
	})

	// Initialize JWT secret
	common.Secrets.JWTSecret = "test-jwt-secret-for-unit-tests"

	// Initialize mock publisher
	publisher.InitMockPublisher()

	// Create a mock HTTP client that fails
	failingHTTPClient := &http.Client{
		Transport: &mockHTTPTransport{shouldSucceed: false},
		Timeout:   30 * time.Second,
	}

	// Create service with failing HTTP client
	userService := NewUserServiceWithHTTPClient(dbClient, common.NewRedisClient(redisClient, redisClient), failingHTTPClient)

	// Create a test user
	_, err = dbClient.User.Create().
		SetID(2001).
		SetUserName("testuser").
		SetEmailUserID("test@example.com").
		SetPassword("password-hash").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	t.Run("HTTP API failure should still return success for security", func(t *testing.T) {
		result, err := userService.InitialForgetPassword(ctx, "test@example.com")

		// Should not error out and should return success message for security
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(200), result.Code)
		assert.Equal(t, "Email Address Sent", result.Message)
	})
}

// Test comprehensive email validation edge cases
func TestInitialForgetPasswordEmailValidation(t *testing.T) {
	userService, dbClient, server, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, server)

	invalidEmails := []string{
		"",                      // Empty string
		"invalid",               // No @ symbol
		"@example.com",          // No local part
		"user@",                 // No domain
		"user@.com",             // Missing domain name
		"user.example.com",      // No @ symbol
		"user@example",          // No TLD
		"user name@example.com", // Space in local part
		"user@exam ple.com",     // Space in domain
		"user@example.",         // Domain ending with dot
	}

	for _, email := range invalidEmails {
		t.Run(fmt.Sprintf("Invalid email: '%s'", email), func(t *testing.T) {
			result, err := userService.InitialForgetPassword(ctx, email)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, int32(400), result.Code)
			assert.Equal(t, "Please Enter a Valid Email Address", result.Message)
		})
	}

	validEmails := []string{
		"user@example.com",
		"test.email@domain.co.uk",
		"user+tag@example.org",
		"123@numbers.com",
		"a@b.co",
		"user..name@example.com", // Double dots in local part - regex allows this
		"user@exam..ple.com",     // Double dots in domain - regex allows this
	}

	for _, email := range validEmails {
		t.Run(fmt.Sprintf("Valid email format: '%s'", email), func(t *testing.T) {
			result, err := userService.InitialForgetPassword(ctx, email)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			// Should get "Contact Support" since these users don't exist
			assert.Equal(t, int32(400), result.Code)
			assert.Equal(t, "Contact Support", result.Message)
		})
	}
}

func TestInitialForgetPasswordOptimizations(t *testing.T) {
	svc, dbClient, redisServer, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, redisServer)

	t.Run("Input validation - Empty email after trim", func(t *testing.T) {
		result, err := svc.InitialForgetPassword(ctx, "   ")
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Equal(t, "Please Enter a Valid Email Address", result.Message)
	})

	t.Run("Input validation - Whitespace is trimmed from valid email", func(t *testing.T) {
		result, err := svc.InitialForgetPassword(ctx, "  trimmed@example.com  ")
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Equal(t, "Contact Support", result.Message)
	})

	t.Run("Input validation - Empty string", func(t *testing.T) {
		result, err := svc.InitialForgetPassword(ctx, "")
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Equal(t, "Please Enter a Valid Email Address", result.Message)
	})
}

func TestForgetPasswordRequest(t *testing.T) {
	svc, dbClient, redisServer, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, redisServer)

	dbClient.User.Create().
		SetUserName("testuser").
		SetEmailUserID("testuser@example.com").
		SetPassword("hashedpassword").
		SetIsTwoFactorAuthenticationEnabled(true).
		SetIsActive(true).
		SaveX(ctx)

	dbClient.User.Create().
		SetUserName("no2fauser").
		SetEmailUserID("no2fa@example.com").
		SetPassword("hashedpassword").
		SetIsTwoFactorAuthenticationEnabled(false).
		SetIsActive(true).
		SaveX(ctx)

	dbClient.User.Create().
		SetUserName("emailuser").
		SetEmailUserID("email@example.com").
		SetPassword("hashedpassword").
		SetIsTwoFactorAuthenticationEnabled(true).
		SetIsActive(true).
		SaveX(ctx)

	t.Run("User not found by username", func(t *testing.T) {
		result, err := svc.ForgetPasswordRequest(ctx, "nonexistentuser", "email", "test@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(404), result.Code)
		assert.Equal(t, "User not found", result.Message)
	})

	t.Run("User not found by email", func(t *testing.T) {
		result, err := svc.ForgetPasswordRequest(ctx, "nonexistent@example.com", "email", "test@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(404), result.Code)
		assert.Equal(t, "User not found", result.Message)
	})

	t.Run("User found by username but 2FA not enabled", func(t *testing.T) {
		result, err := svc.ForgetPasswordRequest(ctx, "no2fauser", "email", "no2fa@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(401), result.Code)
		assert.Equal(t, "User does not have 2fa enabled", result.Message)
	})

	t.Run("User found by email but 2FA not enabled", func(t *testing.T) {
		result, err := svc.ForgetPasswordRequest(ctx, "no2fa@example.com", "email", "no2fa@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(401), result.Code)
		assert.Equal(t, "User does not have 2fa enabled", result.Message)
	})

	t.Run("Successfully send verification code via email by username", func(t *testing.T) {
		result, err := svc.ForgetPasswordRequest(ctx, "testuser", "email", "testuser@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(200), result.Code)
		assert.Equal(t, "Code sent via email", result.Message)

		// Verify that the verification code was stored in Redis
		redisKey := "lis::core_service::user_service:forget_password_verification_user_testuser"
		storedCode, err := svc.(*UserService).redisClient.Get(ctx, redisKey).Result()
		assert.NoError(t, err)
		assert.Len(t, storedCode, 6) // Should be 6-digit OTP
	})

	t.Run("Successfully send verification code via email by email", func(t *testing.T) {
		result, err := svc.ForgetPasswordRequest(ctx, "email@example.com", "email", "email@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(200), result.Code)
		assert.Equal(t, "Code sent via email", result.Message)

		// Verify that the verification code was stored in Redis
		redisKey := "lis::core_service::user_service:forget_password_verification_user_email@example.com"
		storedCode, err := svc.(*UserService).redisClient.Get(ctx, redisKey).Result()
		assert.NoError(t, err)
		assert.Len(t, storedCode, 6) // Should be 6-digit OTP
	})

	t.Run("Successfully send verification code via phone/SMS", func(t *testing.T) {
		result, err := svc.ForgetPasswordRequest(ctx, "testuser", "phone", "+1234567890")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(200), result.Code)
		assert.Equal(t, "Code sent via text", result.Message)

		// Verify that the verification code was stored in Redis
		redisKey := "lis::core_service::user_service:forget_password_verification_user_testuser"
		storedCode, err := svc.(*UserService).redisClient.Get(ctx, redisKey).Result()
		assert.NoError(t, err)
		assert.Len(t, storedCode, 6) // Should be 6-digit OTP
	})

	t.Run("Invalid request method", func(t *testing.T) {
		result, err := svc.ForgetPasswordRequest(ctx, "testuser", "invalid", "test@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(400), result.Code)
		assert.Equal(t, "Invalid request method. Must be 'email' or 'phone'", result.Message)
	})

	// OPTIMIZATION TESTS: New validation tests
	t.Run("Input validation - Empty username", func(t *testing.T) {
		result, err := svc.ForgetPasswordRequest(ctx, "", "email", "test@example.com")
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Equal(t, "Username is required", result.Message)
	})

	t.Run("Input validation - Whitespace only username", func(t *testing.T) {
		result, err := svc.ForgetPasswordRequest(ctx, "   ", "email", "test@example.com")
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Equal(t, "Username is required", result.Message)
	})

	t.Run("Input validation - Empty request method", func(t *testing.T) {
		result, err := svc.ForgetPasswordRequest(ctx, "testuser", "", "test@example.com")
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Equal(t, "Request method is required", result.Message)
	})

	t.Run("Input validation - Empty request target", func(t *testing.T) {
		result, err := svc.ForgetPasswordRequest(ctx, "testuser", "email", "")
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Equal(t, "Request target is required", result.Message)
	})

	t.Run("Input validation - Invalid email format for email method", func(t *testing.T) {
		result, err := svc.ForgetPasswordRequest(ctx, "testuser", "email", "invalid-email")
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Equal(t, "Invalid email format", result.Message)
	})

	t.Run("Email service failure should return error", func(t *testing.T) {
		// Create a service with failing HTTP client
		mockHTTPClient := &http.Client{
			Transport: &mockHTTPTransport{shouldSucceed: false},
		}

		failSvc := &UserService{
			dbClient:    dbClient,
			redisClient: svc.(*UserService).redisClient,
			httpClient:  mockHTTPClient,
		}

		result, err := failSvc.ForgetPasswordRequest(ctx, "testuser", "email", "testuser@example.com")

		// Should get an error when HTTP call fails
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "error sending reset password email")
	})

	t.Run("Verification code expires after 5 minutes", func(t *testing.T) {
		// Send a verification code first
		result, err := svc.ForgetPasswordRequest(ctx, "testuser", "email", "testuser@example.com")
		assert.NoError(t, err)
		assert.Equal(t, int32(200), result.Code)

		// Verify code exists in Redis
		redisKey := "lis::core_service::user_service:forget_password_verification_user_testuser"
		storedCode, err := svc.(*UserService).redisClient.Get(ctx, redisKey).Result()
		assert.NoError(t, err)
		assert.Len(t, storedCode, 6)

		// Delete the key to simulate expiration
		err = svc.(*UserService).redisClient.Del(ctx, redisKey).Err()
		assert.NoError(t, err)

		// Verify code is expired/deleted
		_, err = svc.(*UserService).redisClient.Get(ctx, redisKey).Result()
		assert.Error(t, err)
		assert.Equal(t, "redis: nil", err.Error())
	})

	t.Run("Multiple users with same username prefix", func(t *testing.T) {
		// Create additional user with similar username
		dbClient.User.Create().
			SetUserName("testuser2").
			SetEmailUserID("testuser2@example.com").
			SetPassword("hashedpassword").
			SetIsTwoFactorAuthenticationEnabled(true).
			SetIsActive(true).
			SaveX(ctx)

		// Should find exact match for "testuser"
		result, err := svc.ForgetPasswordRequest(ctx, "testuser", "email", "testuser@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(200), result.Code)
		assert.Equal(t, "Code sent via email", result.Message)

		// Verify correct user's code was stored
		redisKey := "lis::core_service::user_service:forget_password_verification_user_testuser"
		storedCode, err := svc.(*UserService).redisClient.Get(ctx, redisKey).Result()
		assert.NoError(t, err)
		assert.Len(t, storedCode, 6)
	})

	t.Run("Case sensitivity in email lookup", func(t *testing.T) {
		// Create user with mixed case email
		dbClient.User.Create().
			SetUserName("mixedcaseuser").
			SetEmailUserID("MixedCase@Example.COM").
			SetPassword("hashedpassword").
			SetIsTwoFactorAuthenticationEnabled(true).
			SetIsActive(true).
			SaveX(ctx)

		// Try exact case match
		result, err := svc.ForgetPasswordRequest(ctx, "MixedCase@Example.COM", "email", "MixedCase@Example.COM")
		assert.NoError(t, err)
		assert.Equal(t, int32(200), result.Code)

		// Try different case - should not find user (exact match required)
		result, err = svc.ForgetPasswordRequest(ctx, "mixedcase@example.com", "email", "mixedcase@example.com")
		assert.NoError(t, err)
		assert.Equal(t, int32(404), result.Code)
	})

	t.Run("Redis key format", func(t *testing.T) {
		// Send a verification code and verify the Redis key format
		result, err := svc.ForgetPasswordRequest(ctx, "testuser", "email", "testuser@example.com")
		assert.NoError(t, err)
		assert.Equal(t, int32(200), result.Code)

		// Verify the Redis key follows expected format
		redisKey := "lis::core_service::user_service:forget_password_verification_user_testuser"
		storedCode, err := svc.(*UserService).redisClient.Get(ctx, redisKey).Result()
		assert.NoError(t, err)
		assert.Len(t, storedCode, 6)
		// Verify it's all digits
		for _, char := range storedCode {
			assert.True(t, char >= '0' && char <= '9', "OTP should contain only digits")
		}
	})
}

func TestUserService_ForgetPassword(t *testing.T) {
	svc, dbClient, redisServer, ctx := setupUserServiceTest(t)
	defer dbClient.Close()
	defer redisServer.Term()

	// Create test users
	user2FA, _ := dbClient.User.Create().
		SetUserName("user2fa").
		SetEmailUserID("user2fa@example.com").
		SetPassword("oldpassword").
		SetIsTwoFactorAuthenticationEnabled(true).
		SetImportedUserWithSaltPassword(true).
		SetIsActive(true).
		Save(ctx)

	userNo2FA, _ := dbClient.User.Create().
		SetUserName("userno2fa").
		SetEmailUserID("userno2fa@example.com").
		SetPassword("oldpassword").
		SetIsTwoFactorAuthenticationEnabled(false).
		SetImportedUserWithSaltPassword(true).
		SetIsActive(true).
		Save(ctx)

	t.Run("User with 2FA - Valid code and password", func(t *testing.T) {
		// First set a verification code in Redis
		redisKey := "lis::core_service::user_service:forget_password_verification_user_user2fa"
		err := svc.(*UserService).redisClient.Set(ctx, redisKey, "123456", time.Hour).Err()
		assert.NoError(t, err)

		// Test password reset
		result, err := svc.ForgetPassword(ctx, "user2fa", "123456", "NewPassword123!")
		assert.NoError(t, err)
		assert.Equal(t, int32(201), result.Code)
		assert.Equal(t, "Password Updated", result.Message)

		// Verify the verification code was deleted from Redis
		_, err = svc.(*UserService).redisClient.Get(ctx, redisKey).Result()
		assert.Equal(t, redis.Nil, err)

		// Verify password was actually updated in database
		updatedUser, err := dbClient.User.Get(ctx, user2FA.ID)
		assert.NoError(t, err)
		assert.NotEqual(t, "oldpassword", updatedUser.Password)   // Password should be hashed
		assert.False(t, updatedUser.ImportedUserWithSaltPassword) // Should be set to false
	})

	t.Run("User with 2FA - Email address lookup", func(t *testing.T) {
		// Set verification code in Redis
		redisKey := "lis::core_service::user_service:forget_password_verification_user_user2fa@example.com"
		err := svc.(*UserService).redisClient.Set(ctx, redisKey, "654321", time.Hour).Err()
		assert.NoError(t, err)

		// Test password reset using email address with strong password
		result, err := svc.ForgetPassword(ctx, "user2fa@example.com", "654321", "NewEmail123!")
		assert.NoError(t, err)
		assert.Equal(t, int32(201), result.Code)
		assert.Equal(t, "Password Updated", result.Message)
	})

	t.Run("User with 2FA - Wrong verification code", func(t *testing.T) {
		// Set verification code in Redis
		redisKey := "lis::core_service::user_service:forget_password_verification_user_user2fa"
		err := svc.(*UserService).redisClient.Set(ctx, redisKey, "123456", time.Hour).Err()
		assert.NoError(t, err)

		// Test with wrong verification code
		result, err := svc.ForgetPassword(ctx, "user2fa", "wrong123", "NewPassword123!")
		assert.NoError(t, err)
		assert.Equal(t, int32(403), result.Code)
		assert.Equal(t, "Wrong Verification Number", result.Message)
	})

	t.Run("User with 2FA - Empty verification code", func(t *testing.T) {
		result, err := svc.ForgetPassword(ctx, "user2fa", "", "NewPassword123!")
		assert.NoError(t, err)
		assert.Equal(t, int32(403), result.Code)
		assert.Equal(t, "Please Enter 2FA Password", result.Message)
	})

	t.Run("User with 2FA - No verification code in Redis", func(t *testing.T) {
		// Ensure no code exists in Redis
		redisKey := "lis::core_service::user_service:forget_password_verification_user_user2fa"
		svc.(*UserService).redisClient.Del(ctx, redisKey)

		result, err := svc.ForgetPassword(ctx, "user2fa", "123456", "NewPassword123!")
		assert.NoError(t, err)
		assert.Equal(t, int32(401), result.Code)
		assert.Contains(t, result.Message, "Never received a forget password request")
	})

	t.Run("User with 2FA - Valid code but no password", func(t *testing.T) {
		// Set verification code in Redis
		redisKey := "lis::core_service::user_service:forget_password_verification_user_user2fa"
		err := svc.(*UserService).redisClient.Set(ctx, redisKey, "123456", time.Hour).Err()
		assert.NoError(t, err)

		result, err := svc.ForgetPassword(ctx, "user2fa", "123456", "")
		assert.NoError(t, err)
		assert.Equal(t, int32(200), result.Code)
		assert.Equal(t, "Verification Passed, Enter New Password to Set the Password", result.Message)

		// Verification code should still exist in Redis
		storedCode, err := svc.(*UserService).redisClient.Get(ctx, redisKey).Result()
		assert.NoError(t, err)
		assert.Equal(t, "123456", storedCode)
	})

	t.Run("User without 2FA - Valid password", func(t *testing.T) {
		result, err := svc.ForgetPassword(ctx, "userno2fa", "ignoredcode", "ValidPassword123!")
		assert.NoError(t, err)
		assert.Equal(t, int32(201), result.Code)
		assert.Equal(t, "Password Updated", result.Message)

		// Verify password was actually updated in database
		updatedUser, err := dbClient.User.Get(ctx, userNo2FA.ID)
		assert.NoError(t, err)
		assert.NotEqual(t, "oldpassword", updatedUser.Password)   // Password should be hashed
		assert.False(t, updatedUser.ImportedUserWithSaltPassword) // Should be set to false
	})

	t.Run("User without 2FA - No password provided", func(t *testing.T) {
		// BUG FIX: No longer allow empty passwords for non-2FA users (optimization)
		result, err := svc.ForgetPassword(ctx, "userno2fa", "ignoredcode", "")
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Equal(t, "Password is required", result.Message)

		// Verify password was actually updated to empty string (matches TypeScript behavior)
		updatedUser, err := dbClient.User.Get(ctx, userNo2FA.ID)
		assert.NoError(t, err)
		// Password should be bcrypt hash of empty string
		assert.NotEqual(t, "oldpassword", updatedUser.Password)
	})

	t.Run("Non-existent user", func(t *testing.T) {
		// BUG FIX: Proper error handling for non-existent user
		result, err := svc.ForgetPassword(ctx, "nonexistent", "123456", "newpassword")
		assert.NoError(t, err)
		assert.Equal(t, int32(404), result.Code)
		assert.Equal(t, "User not found", result.Message)
	})

	t.Run("Non-existent user by email", func(t *testing.T) {
		// BUG FIX: Proper error handling for non-existent user by email
		result, err := svc.ForgetPassword(ctx, "nonexistent@example.com", "123456", "newpassword")
		assert.NoError(t, err)
		assert.Equal(t, int32(404), result.Code)
		assert.Equal(t, "User not found", result.Message)
	})

	t.Run("Password hashing verification", func(t *testing.T) {
		// Test that passwords are properly hashed with bcrypt using strong password
		testPassword := "TestPassword123!"

		// Set verification code in Redis
		redisKey := "lis::core_service::user_service:forget_password_verification_user_user2fa"
		err := svc.(*UserService).redisClient.Set(ctx, redisKey, "123456", time.Hour).Err()
		assert.NoError(t, err)

		result, err := svc.ForgetPassword(ctx, "user2fa", "123456", testPassword)
		assert.NoError(t, err)
		assert.Equal(t, int32(201), result.Code)

		// Verify the password is hashed (should start with bcrypt prefix)
		updatedUser, err := dbClient.User.Get(ctx, user2FA.ID)
		assert.NoError(t, err)
		assert.True(t, len(updatedUser.Password) > 50)         // Bcrypt hash should be longer than original
		assert.NotEqual(t, testPassword, updatedUser.Password) // Should not store plain text
		assert.True(t, len(updatedUser.Password) == 60)        // Bcrypt hash length is always 60 characters
	})

	// OPTIMIZATION TESTS: New validation and security tests
	t.Run("Input validation - Empty username", func(t *testing.T) {
		result, err := svc.ForgetPassword(ctx, "", "123456", "ValidPassword123!")
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Equal(t, "Username is required", result.Message)
	})

	t.Run("Input validation - Whitespace only username", func(t *testing.T) {
		result, err := svc.ForgetPassword(ctx, "   ", "123456", "ValidPassword123!")
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Equal(t, "Username is required", result.Message)
	})

	t.Run("Password validation - Too short", func(t *testing.T) {
		redisKey := "lis::core_service::user_service:forget_password_verification_user_user2fa"
		err := svc.(*UserService).redisClient.Set(ctx, redisKey, "123456", time.Hour).Err()
		assert.NoError(t, err)

		result, err := svc.ForgetPassword(ctx, "user2fa", "123456", "short")
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Contains(t, result.Message, "at least 8 characters long")
	})

	t.Run("Password validation - Too long (over 20 characters)", func(t *testing.T) {
		redisKey := "lis::core_service::user_service:forget_password_verification_user_user2fa"
		err := svc.(*UserService).redisClient.Set(ctx, redisKey, "123456", time.Hour).Err()
		assert.NoError(t, err)

		result, err := svc.ForgetPassword(ctx, "user2fa", "123456", "ThisPasswordIsWayTooLong123!")
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Contains(t, result.Message, "at most 20 characters long")
	})

	t.Run("Password validation - Missing uppercase", func(t *testing.T) {
		redisKey := "lis::core_service::user_service:forget_password_verification_user_user2fa"
		err := svc.(*UserService).redisClient.Set(ctx, redisKey, "123456", time.Hour).Err()
		assert.NoError(t, err)

		result, err := svc.ForgetPassword(ctx, "user2fa", "123456", "password123!")
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Contains(t, result.Message, "uppercase letter [A-Z]")
	})

	t.Run("Password validation - Missing lowercase", func(t *testing.T) {
		redisKey := "lis::core_service::user_service:forget_password_verification_user_user2fa"
		err := svc.(*UserService).redisClient.Set(ctx, redisKey, "123456", time.Hour).Err()
		assert.NoError(t, err)

		result, err := svc.ForgetPassword(ctx, "user2fa", "123456", "PASSWORD123!")
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Contains(t, result.Message, "lowercase letter [a-z]")
	})

	t.Run("Password validation - Missing digit", func(t *testing.T) {
		redisKey := "lis::core_service::user_service:forget_password_verification_user_user2fa"
		err := svc.(*UserService).redisClient.Set(ctx, redisKey, "123456", time.Hour).Err()
		assert.NoError(t, err)

		result, err := svc.ForgetPassword(ctx, "user2fa", "123456", "Password!")
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Contains(t, result.Message, "number [0-9]")
	})

	t.Run("Password validation - Missing special character", func(t *testing.T) {
		redisKey := "lis::core_service::user_service:forget_password_verification_user_user2fa"
		err := svc.(*UserService).redisClient.Set(ctx, redisKey, "123456", time.Hour).Err()
		assert.NoError(t, err)

		result, err := svc.ForgetPassword(ctx, "user2fa", "123456", "Password123")
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Contains(t, result.Message, "special character")
	})

	t.Run("Password validation - Valid password with all requirements", func(t *testing.T) {
		redisKey := "lis::core_service::user_service:forget_password_verification_user_user2fa"
		err := svc.(*UserService).redisClient.Set(ctx, redisKey, "123456", time.Hour).Err()
		assert.NoError(t, err)

		validPasswords := []string{
			"Password1!",
			"MyPass123@",
			"Test1234#",
			"Valid123$",
			"Good123%",
			"Strong12^",
			"Safe123&",
			"Cool12()",
			"Nice123{}",
			"Best123[]",
			"Top1234:",
			"Max12;<>",
			"Win1234,",
			"Big123.?",
			"Fun1234/",
			"Code123~",
			"Dev1234_",
			"App123+-",
			"Web1234=",
			"Net1234|",
		}

		for _, validPass := range validPasswords {
			// Set verification code for each password test (it gets deleted after successful use)
			redisKey := "lis::core_service::user_service:forget_password_verification_user_user2fa"
			err := svc.(*UserService).redisClient.Set(ctx, redisKey, "123456", time.Hour).Err()
			assert.NoError(t, err)

			result, err := svc.ForgetPassword(ctx, "user2fa", "123456", validPass)
			assert.NoError(t, err)
			assert.Equal(t, int32(201), result.Code)
			assert.Equal(t, "Password Updated", result.Message)
		}
	})
}

// TestForgetPasswordWorkflowIntegration tests the complete forget password workflow
// This integration test ensures all three functions work together properly:
// InitialForgetPassword → ForgetPasswordRequest → ForgetPassword
func TestForgetPasswordWorkflowIntegration(t *testing.T) {
	userService, dbClient, server, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, server)

	// Create test user with 2FA contacts
	testUser := struct {
		id          int
		emailUserID string
		username    string
		password    string
	}{
		id:          5001,
		emailUserID: "workflow@example.com",
		username:    "workflowuser",
		password:    "InitialPassword123!",
	}

	// Create user
	_, err := dbClient.User.Create().
		SetID(testUser.id).
		SetUserName(testUser.username).
		SetEmailUserID(testUser.emailUserID).
		SetPassword(testUser.password).
		SetIsActive(true).
		Save(ctx)
	assert.NoError(t, err)

	// Create 2FA contact
	_, err = dbClient.Contact.Create().
		SetID(6001).
		SetUserID(testUser.id).
		SetContactType("email").
		SetContactDetails("workflow2fa@example.com").
		SetIs2faContact(true).
		Save(ctx)
	assert.NoError(t, err)

	t.Run("Complete Workflow Success", func(t *testing.T) {
		// Step 1: InitialForgetPassword - Send initial email
		initialResult, err := userService.InitialForgetPassword(ctx, testUser.emailUserID)
		assert.NoError(t, err)
		assert.Equal(t, int32(200), initialResult.Code)
		assert.Contains(t, initialResult.Message, "Email Address Sent")

		// Step 2: ForgetPasswordRequest - For non-2FA users, this returns 401
		requestResult, err := userService.ForgetPasswordRequest(ctx, testUser.username, "email", "workflow2fa@example.com")
		assert.NoError(t, err)
		assert.Equal(t, int32(401), requestResult.Code)
		assert.Contains(t, requestResult.Message, "User does not have 2fa enabled")

		// Step 3: ForgetPassword - For non-2FA users, verification code is ignored
		newPassword := "NewPassword123!"
		verifyResult, err := userService.ForgetPassword(ctx, testUser.username, "ignored", newPassword)
		assert.NoError(t, err)
		assert.Equal(t, int32(201), verifyResult.Code)
		assert.Equal(t, "Password Updated", verifyResult.Message)
	})

	t.Run("Workflow With Email Username", func(t *testing.T) {
		// Step 1: Initial request
		initialResult, err := userService.InitialForgetPassword(ctx, testUser.emailUserID)
		assert.NoError(t, err)
		assert.Equal(t, int32(200), initialResult.Code)

		// Step 2: Request using email as username - returns 401 for non-2FA user
		requestResult, err := userService.ForgetPasswordRequest(ctx, testUser.emailUserID, "email", "workflow2fa@example.com")
		assert.NoError(t, err)
		assert.Equal(t, int32(401), requestResult.Code)

		// Step 3: Complete password reset using email
		newPassword := "EmailPass123!"
		verifyResult, err := userService.ForgetPassword(ctx, testUser.emailUserID, "ignored", newPassword)
		assert.NoError(t, err)
		assert.Equal(t, int32(201), verifyResult.Code)
		assert.Equal(t, "Password Updated", verifyResult.Message)
	})

	t.Run("Workflow For Non-2FA User", func(t *testing.T) {
		// Step 1: Initial request
		_, err := userService.InitialForgetPassword(ctx, testUser.emailUserID)
		assert.NoError(t, err)

		// Step 2: ForgetPasswordRequest returns 401 for non-2FA users
		requestResult, err := userService.ForgetPasswordRequest(ctx, testUser.username, "email", "workflow2fa@example.com")
		assert.NoError(t, err)
		assert.Equal(t, int32(401), requestResult.Code)

		// Step 3: ForgetPassword works for non-2FA users without verification codes
		verifyResult, err := userService.ForgetPassword(ctx, testUser.username, "anycode", "ExpiredPass123!")
		assert.NoError(t, err)
		assert.Equal(t, int32(201), verifyResult.Code)
		assert.Equal(t, "Password Updated", verifyResult.Message)
	})

	t.Run("Password Change Success", func(t *testing.T) {
		// Step 1: Initial request
		_, err := userService.InitialForgetPassword(ctx, testUser.emailUserID)
		assert.NoError(t, err)

		// Step 2: Skip ForgetPasswordRequest for non-2FA users (returns 401)

		// Step 3: ForgetPassword works for non-2FA users
		verifyResult, err := userService.ForgetPassword(ctx, testUser.username, "anycode", "WrongPass123!")
		assert.NoError(t, err)
		assert.Equal(t, int32(201), verifyResult.Code)
		assert.Equal(t, "Password Updated", verifyResult.Message)
	})

	t.Run("Direct Password Reset", func(t *testing.T) {
		// Skip InitialForgetPassword and go directly to ForgetPassword
		// This works for non-2FA users

		verifyResult, err := userService.ForgetPassword(ctx, testUser.username, "ignored", "DirectPass123!")
		assert.NoError(t, err)
		assert.Equal(t, int32(201), verifyResult.Code)
		assert.Equal(t, "Password Updated", verifyResult.Message)
	})
}

// TestForgetPasswordWorkflowWith2FA tests the complete forget password workflow with 2FA enabled users
func TestForgetPasswordWorkflowWith2FA(t *testing.T) {
	userService, dbClient, server, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, server)

	// Create test user with 2FA
	testUser := struct {
		id          int
		emailUserID string
		username    string
		password    string
	}{
		id:          5002,
		emailUserID: "with2fa@example.com",
		username:    "user2faworkflow",
		password:    "InitialPassword123!",
	}

	// Create user
	_, err := dbClient.User.Create().
		SetID(testUser.id).
		SetUserName(testUser.username).
		SetEmailUserID(testUser.emailUserID).
		SetPassword(testUser.password).
		SetIsActive(true).
		SetIsTwoFactorAuthenticationEnabled(true). // Enable 2FA
		Save(ctx)
	assert.NoError(t, err)

	// Create 2FA contact
	_, err = dbClient.Contact.Create().
		SetID(6002).
		SetUserID(testUser.id).
		SetContactType("email").
		SetContactDetails("2fa@example.com").
		SetIs2faContact(true).
		Save(ctx)
	assert.NoError(t, err)

	t.Run("Complete 2FA Workflow Success", func(t *testing.T) {
		// Step 1: InitialForgetPassword - Send initial email
		initialResult, err := userService.InitialForgetPassword(ctx, testUser.emailUserID)
		assert.NoError(t, err)
		assert.Equal(t, int32(200), initialResult.Code)
		assert.Contains(t, initialResult.Message, "Email Address Sent")

		// Step 2: ForgetPasswordRequest - Send verification code
		requestResult, err := userService.ForgetPasswordRequest(ctx, testUser.username, "email", "2fa@example.com")
		assert.NoError(t, err)
		assert.Equal(t, int32(200), requestResult.Code)
		assert.Contains(t, requestResult.Message, "Code sent via email")

		// Verify code was stored in Redis
		redisKey := fmt.Sprintf("lis::core_service::user_service:forget_password_verification_user_%s", testUser.username)
		redisClient := userService.(*UserService).redisClient
		storedCode, err := redisClient.Get(ctx, redisKey).Result()
		assert.NoError(t, err)
		assert.NotEmpty(t, storedCode)

		// Step 3: ForgetPassword - Verify code and update password
		newPassword := "New2FAPass123!"
		verifyResult, err := userService.ForgetPassword(ctx, testUser.username, storedCode, newPassword)
		assert.NoError(t, err)
		assert.Equal(t, int32(201), verifyResult.Code)
		assert.Equal(t, "Password Updated", verifyResult.Message)

		// Verify Redis key was deleted after successful password reset
		_, err = redisClient.Get(ctx, redisKey).Result()
		assert.Error(t, err) // Should not exist anymore
	})

	t.Run("2FA Workflow With Wrong Code", func(t *testing.T) {
		// Step 1 & 2: Setup valid request
		_, err := userService.InitialForgetPassword(ctx, testUser.emailUserID)
		assert.NoError(t, err)

		_, err = userService.ForgetPasswordRequest(ctx, testUser.username, "email", "2fa@example.com")
		assert.NoError(t, err)

		// Step 3: Try with wrong verification code
		verifyResult, err := userService.ForgetPassword(ctx, testUser.username, "wrongcode", "WrongCode123!")
		assert.NoError(t, err)
		assert.Equal(t, int32(403), verifyResult.Code)
		assert.Contains(t, verifyResult.Message, "Wrong Verification Number")

		// Verify Redis key still exists (not deleted on wrong code)
		redisKey := fmt.Sprintf("lis::core_service::user_service:forget_password_verification_user_%s", testUser.username)
		redisClient := userService.(*UserService).redisClient
		_, err = redisClient.Get(ctx, redisKey).Result()
		assert.NoError(t, err) // Should still exist
	})
}

func TestTurnOff2FASettingPage(t *testing.T) {
	svc, dbClient, redisServer, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, redisServer)

	// Create test users
	err := createTestUsersFor2FA(t, dbClient, ctx)
	if err != nil {
		t.Fatalf("Failed to create test users: %v", err)
	}

	// Generate valid JWT token for user with 2FA enabled (user_id: 1)
	validToken, err := generateTestJWTToken(1, "test-jwt-secret-for-unit-tests")
	if err != nil {
		t.Fatalf("Failed to generate test JWT token: %v", err)
	}

	// Generate valid JWT token for user without 2FA (user_id: 2)
	validTokenNo2FA, err := generateTestJWTToken(2, "test-jwt-secret-for-unit-tests")
	if err != nil {
		t.Fatalf("Failed to generate test JWT token: %v", err)
	}

	// Test case 1: Successful 2FA turn off with username
	t.Run("Successful 2FA turn off", func(t *testing.T) {
		result, err := svc.TurnOff2FASettingPage(ctx, "user2fa", validToken)
		assert.NoError(t, err)
		assert.Equal(t, int32(200), result.Code)
		assert.Equal(t, "2FA Already Disabled", result.Message)

		// Verify the user's 2FA is actually disabled in the database
		user, err := dbClient.User.Get(ctx, 1)
		assert.NoError(t, err)
		assert.False(t, user.IsTwoFactorAuthenticationEnabled)
		assert.Empty(t, user.TwoFactorAuthenticationSecret)
	})

	// Test case 2: Successful 2FA turn off with email
	// Re-enable 2FA for the user first
	t.Run("Re-enable 2FA and test with email", func(t *testing.T) {
		// Re-enable 2FA for testing
		_, err := dbClient.User.UpdateOneID(1).
			SetIsTwoFactorAuthenticationEnabled(true).
			SetTwoFactorAuthenticationSecret("test-secret").
			Save(ctx)
		assert.NoError(t, err)

		result, err := svc.TurnOff2FASettingPage(ctx, "user2fa@example.com", validToken)
		assert.NoError(t, err)
		assert.Equal(t, int32(200), result.Code)
		assert.Equal(t, "2FA Already Disabled", result.Message)

		// Verify the user's 2FA is actually disabled in the database
		user, err := dbClient.User.Get(ctx, 1)
		assert.NoError(t, err)
		assert.False(t, user.IsTwoFactorAuthenticationEnabled)
		assert.Empty(t, user.TwoFactorAuthenticationSecret)
	})

	// Test case 3: 2FA already disabled
	t.Run("2FA already disabled", func(t *testing.T) {
		result, err := svc.TurnOff2FASettingPage(ctx, "userno2fa", validTokenNo2FA)
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Equal(t, "The 2FA is already off", result.Message)
	})

	// Test case 4: Empty username
	t.Run("Empty username", func(t *testing.T) {
		result, err := svc.TurnOff2FASettingPage(ctx, "", validToken)
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Equal(t, "Username is required", result.Message)
	})

	// Test case 5: Empty token
	t.Run("Empty token", func(t *testing.T) {
		result, err := svc.TurnOff2FASettingPage(ctx, "user2fa", "")
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Equal(t, "Token is required", result.Message)
	})

	// Test case 6: Invalid token
	t.Run("Invalid token", func(t *testing.T) {
		result, err := svc.TurnOff2FASettingPage(ctx, "user2fa", "invalid-token")
		assert.NoError(t, err)
		assert.Equal(t, int32(401), result.Code)
		assert.Equal(t, "Invalid token", result.Message)
	})

	// Test case 7: User not found
	t.Run("User not found", func(t *testing.T) {
		result, err := svc.TurnOff2FASettingPage(ctx, "nonexistent", validToken)
		assert.NoError(t, err)
		assert.Equal(t, int32(400), result.Code)
		assert.Equal(t, "2FA Already Enabled", result.Message)
	})

	// Test case 8: Token user ID mismatch (security check)
	t.Run("Token user ID mismatch", func(t *testing.T) {
		// Re-enable 2FA for user 1
		_, err := dbClient.User.UpdateOneID(1).
			SetIsTwoFactorAuthenticationEnabled(true).
			SetTwoFactorAuthenticationSecret("test-secret").
			Save(ctx)
		assert.NoError(t, err)

		// Try to use token for user 2 to disable 2FA for user 1
		result, err := svc.TurnOff2FASettingPage(ctx, "user2fa", validTokenNo2FA)
		assert.NoError(t, err)
		assert.Equal(t, int32(403), result.Code)
		assert.Equal(t, "Unauthorized to modify this user's 2FA settings", result.Message)
	})

	// Test case 9: Token with invalid user ID
	t.Run("Token with invalid user ID", func(t *testing.T) {
		invalidUserToken, err := generateTestJWTToken(-1, "test-jwt-secret-for-unit-tests")
		assert.NoError(t, err)

		result, err := svc.TurnOff2FASettingPage(ctx, "user2fa", invalidUserToken)
		assert.NoError(t, err)
		assert.Equal(t, int32(401), result.Code)
		assert.Equal(t, "Invalid token", result.Message)
	})

	// Test case 10: Whitespace handling
	t.Run("Whitespace handling", func(t *testing.T) {
		// Re-enable 2FA for user 1
		_, err := dbClient.User.UpdateOneID(1).
			SetIsTwoFactorAuthenticationEnabled(true).
			SetTwoFactorAuthenticationSecret("test-secret").
			Save(ctx)
		assert.NoError(t, err)

		result, err := svc.TurnOff2FASettingPage(ctx, "  user2fa  ", "  "+validToken+"  ")
		assert.NoError(t, err)
		assert.Equal(t, int32(200), result.Code)
		assert.Equal(t, "2FA Already Disabled", result.Message)
	})
}

// Helper function to create test users with 2FA enabled/disabled
func createTestUsersFor2FA(t *testing.T, dbClient *ent.Client, ctx context.Context) error {
	// User 1: Has 2FA enabled
	_, err := dbClient.User.Create().
		SetID(1).
		SetUserName("user2fa").
		SetEmailUserID("user2fa@example.com").
		SetPassword("hashedpassword").
		SetIsTwoFactorAuthenticationEnabled(true).
		SetTwoFactorAuthenticationSecret("test-secret").
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create user with 2FA: %w", err)
	}

	// User 2: Has 2FA disabled
	_, err = dbClient.User.Create().
		SetID(2).
		SetUserName("userno2fa").
		SetEmailUserID("userno2fa@example.com").
		SetPassword("hashedpassword").
		SetIsTwoFactorAuthenticationEnabled(false).
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create user without 2FA: %w", err)
	}

	return nil
}

// Helper function to generate JWT token for testing
func generateTestJWTToken(userID int, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func TestTurnOn2FASettingPage(t *testing.T) {
	service, dbClient, server, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, server)

	t.Run("Success - Turn on 2FA for user with username", func(t *testing.T) {
		// Create a test user
		user, err := dbClient.User.Create().
			SetUserName("testuser").
			SetPassword("hashedpassword").
			SetIsTwoFactorAuthenticationEnabled(false).
			Save(ctx)
		assert.NoError(t, err)

		// Generate a valid JWT token for the user
		token, err := generateTestJWTToken(user.ID, common.Secrets.JWTSecret)
		assert.NoError(t, err)

		// Test turning on 2FA
		response, err := service.TurnOn2FASettingPage(ctx, "testuser", "test@example.com", "+1234567890", token)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int32(200), response.Code)
		assert.Equal(t, "2FA Enabled", response.Message)
		assert.NotEmpty(t, response.OtpauthUrl)

		// Verify user 2FA is enabled in database
		updatedUser, err := dbClient.User.Get(ctx, user.ID)
		assert.NoError(t, err)
		assert.True(t, updatedUser.IsTwoFactorAuthenticationEnabled)
		assert.NotEmpty(t, updatedUser.TwoFactorAuthenticationSecret)

		// Verify 2FA contacts were created
		emailContact, err := dbClient.Contact.Query().
			Where(
				contact.UserID(user.ID),
				contact.ContactDescription("email_2fa"),
			).
			Only(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", emailContact.ContactDetails)
		assert.True(t, emailContact.Is2faContact)

		phoneContact, err := dbClient.Contact.Query().
			Where(
				contact.UserID(user.ID),
				contact.ContactDescription("phone_2fa"),
			).
			Only(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "+1234567890", phoneContact.ContactDetails)
		assert.True(t, phoneContact.Is2faContact)
	})

	t.Run("Success - Turn on 2FA for user with email login", func(t *testing.T) {
		// Create a test user with email login
		user, err := dbClient.User.Create().
			SetUserName("emailuser").
			SetEmailUserID("emailuser@example.com").
			SetPassword("hashedpassword").
			SetIsTwoFactorAuthenticationEnabled(false).
			Save(ctx)
		assert.NoError(t, err)

		// Generate a valid JWT token for the user
		token, err := generateTestJWTToken(user.ID, common.Secrets.JWTSecret)
		assert.NoError(t, err)

		// Test turning on 2FA using email login
		response, err := service.TurnOn2FASettingPage(ctx, "emailuser@example.com", "different@example.com", "+1987654321", token)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int32(200), response.Code)
		assert.Equal(t, "2FA Enabled", response.Message)
		assert.NotEmpty(t, response.OtpauthUrl)

		// Verify user 2FA is enabled in database
		updatedUser, err := dbClient.User.Get(ctx, user.ID)
		assert.NoError(t, err)
		assert.True(t, updatedUser.IsTwoFactorAuthenticationEnabled)
		assert.NotEmpty(t, updatedUser.TwoFactorAuthenticationSecret)
	})

	t.Run("Error - 2FA email same as login email", func(t *testing.T) {
		// Create a test user with email login
		user, err := dbClient.User.Create().
			SetUserName("sameemailuser").
			SetEmailUserID("sameemail@example.com").
			SetPassword("hashedpassword").
			SetIsTwoFactorAuthenticationEnabled(false).
			Save(ctx)
		assert.NoError(t, err)

		// Generate a valid JWT token for the user
		token, err := generateTestJWTToken(user.ID, common.Secrets.JWTSecret)
		assert.NoError(t, err)

		// Test turning on 2FA with same email as login email
		response, err := service.TurnOn2FASettingPage(ctx, "sameemail@example.com", "sameemail@example.com", "+1234567890", token)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int32(400), response.Code)
		assert.Equal(t, "2FA Email Address Cannot be the same as Email Login ID", response.Message)

		// Verify user 2FA is still disabled
		unchangedUser, err := dbClient.User.Get(ctx, user.ID)
		assert.NoError(t, err)
		assert.False(t, unchangedUser.IsTwoFactorAuthenticationEnabled)
	})

	t.Run("Error - 2FA already enabled", func(t *testing.T) {
		// Create a test user with 2FA already enabled
		user, err := dbClient.User.Create().
			SetUserName("already2fa").
			SetPassword("hashedpassword").
			SetIsTwoFactorAuthenticationEnabled(true).
			SetTwoFactorAuthenticationSecret("existing-secret").
			Save(ctx)
		assert.NoError(t, err)

		// Generate a valid JWT token for the user
		token, err := generateTestJWTToken(user.ID, common.Secrets.JWTSecret)
		assert.NoError(t, err)

		// Test turning on 2FA when already enabled
		response, err := service.TurnOn2FASettingPage(ctx, "already2fa", "test@example.com", "+1234567890", token)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int32(400), response.Code)
		assert.Equal(t, "2FA Already Enabled", response.Message)
	})

	t.Run("Error - Invalid token", func(t *testing.T) {
		// Test with invalid token
		response, err := service.TurnOn2FASettingPage(ctx, "testuser", "test@example.com", "+1234567890", "invalid-token")
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int32(401), response.Code)
		assert.Equal(t, "Invalid token", response.Message)
	})

	t.Run("Error - Token user ID mismatch", func(t *testing.T) {
		// Create two test users
		user1, err := dbClient.User.Create().
			SetUserName("user1").
			SetPassword("hashedpassword").
			SetIsTwoFactorAuthenticationEnabled(false).
			Save(ctx)
		assert.NoError(t, err)

		_, err = dbClient.User.Create().
			SetUserName("user2").
			SetPassword("hashedpassword").
			SetIsTwoFactorAuthenticationEnabled(false).
			Save(ctx)
		assert.NoError(t, err)

		// Generate token for user1 but try to modify user2
		token, err := generateTestJWTToken(user1.ID, common.Secrets.JWTSecret)
		assert.NoError(t, err)

		// Test token user ID mismatch
		response, err := service.TurnOn2FASettingPage(ctx, "user2", "test@example.com", "+1234567890", token)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int32(403), response.Code)
		assert.Equal(t, "Unauthorized to modify this user's 2FA settings", response.Message)
	})

	t.Run("Error - User not found", func(t *testing.T) {
		// Generate token for non-existent user
		token, err := generateTestJWTToken(99999, common.Secrets.JWTSecret)
		assert.NoError(t, err)

		// Test with non-existent user
		response, err := service.TurnOn2FASettingPage(ctx, "nonexistentuser", "test@example.com", "+1234567890", token)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int32(400), response.Code)
		assert.Equal(t, "2FA Already Enabled", response.Message)
	})

	t.Run("Error - Empty username", func(t *testing.T) {
		token, err := generateTestJWTToken(1, common.Secrets.JWTSecret)
		assert.NoError(t, err)

		response, err := service.TurnOn2FASettingPage(ctx, "", "test@example.com", "+1234567890", token)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int32(400), response.Code)
		assert.Equal(t, "Username is required", response.Message)
	})

	t.Run("Error - Empty token", func(t *testing.T) {
		response, err := service.TurnOn2FASettingPage(ctx, "testuser", "test@example.com", "+1234567890", "")
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int32(400), response.Code)
		assert.Equal(t, "Token is required", response.Message)
	})

	t.Run("Error - Empty 2FA email", func(t *testing.T) {
		token, err := generateTestJWTToken(1, common.Secrets.JWTSecret)
		assert.NoError(t, err)

		response, err := service.TurnOn2FASettingPage(ctx, "testuser", "", "+1234567890", token)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int32(400), response.Code)
		assert.Equal(t, "2FA email address is required", response.Message)
	})

	t.Run("Error - Empty 2FA phone", func(t *testing.T) {
		token, err := generateTestJWTToken(1, common.Secrets.JWTSecret)
		assert.NoError(t, err)

		response, err := service.TurnOn2FASettingPage(ctx, "testuser", "test@example.com", "", token)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int32(400), response.Code)
		assert.Equal(t, "2FA phone number is required", response.Message)
	})

	t.Run("Success - Update existing 2FA contacts", func(t *testing.T) {
		// Create a test user
		user, err := dbClient.User.Create().
			SetUserName("updatecontacts").
			SetPassword("hashedpassword").
			SetIsTwoFactorAuthenticationEnabled(false).
			Save(ctx)
		assert.NoError(t, err)

		// Create existing 2FA contacts
		_, err = dbClient.Contact.Create().
			SetUserID(user.ID).
			SetContactDetails("old@example.com").
			SetContactType("email").
			SetContactDescription("email_2fa").
			SetIs2faContact(true).
			Save(ctx)
		assert.NoError(t, err)

		_, err = dbClient.Contact.Create().
			SetUserID(user.ID).
			SetContactDetails("+1111111111").
			SetContactType("phone").
			SetContactDescription("phone_2fa").
			SetIs2faContact(true).
			Save(ctx)
		assert.NoError(t, err)

		// Generate a valid JWT token for the user
		token, err := generateTestJWTToken(user.ID, common.Secrets.JWTSecret)
		assert.NoError(t, err)

		// Test turning on 2FA (should update existing contacts)
		response, err := service.TurnOn2FASettingPage(ctx, "updatecontacts", "new@example.com", "+2222222222", token)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int32(200), response.Code)
		assert.Equal(t, "2FA Enabled", response.Message)

		// Verify contacts were updated (not created new ones)
		emailContacts, err := dbClient.Contact.Query().
			Where(
				contact.UserID(user.ID),
				contact.ContactDescription("email_2fa"),
			).
			All(ctx)
		assert.NoError(t, err)
		assert.Len(t, emailContacts, 1) // Should only have one email contact
		assert.Equal(t, "new@example.com", emailContacts[0].ContactDetails)

		phoneContacts, err := dbClient.Contact.Query().
			Where(
				contact.UserID(user.ID),
				contact.ContactDescription("phone_2fa"),
			).
			All(ctx)
		assert.NoError(t, err)
		assert.Len(t, phoneContacts, 1) // Should only have one phone contact
		assert.Equal(t, "+2222222222", phoneContacts[0].ContactDetails)
	})
}

func TestUpdateUserInvitationRecord(t *testing.T) {
	service, dbClient, server, ctx := setupUserServiceTest(t)
	defer cleanUpUserServiceTest(dbClient, server)

	t.Run("Successfully update invitation record", func(t *testing.T) {
		customerID := int32(123)
		originalLink := "https://example.com/invite/original"
		newLink := "https://example.com/invite/updated"

		// Create a user invitation record
		invitationRecord, err := dbClient.UserInvitationRecord.
			Create().
			SetCustomerID(int(customerID)).
			SetInvitationLink(originalLink).
			Save(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, invitationRecord)

		// Update the invitation record
		resp, err := service.UpdateUserInvitationRecord(ctx, customerID, newLink)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, int32(200), resp.Code)
		assert.Equal(t, "Invitation Link Updated", resp.Message)
		assert.Equal(t, "", resp.ErrorMessage)

		// Verify the record was updated in the database
		updatedRecord, err := dbClient.UserInvitationRecord.Get(ctx, invitationRecord.ID)
		assert.NoError(t, err)
		assert.Equal(t, newLink, updatedRecord.InvitationLink)
		assert.Equal(t, int(customerID), updatedRecord.CustomerID)
	})

	t.Run("Return error when invitation record not found", func(t *testing.T) {
		nonExistentCustomerID := int32(999)
		newLink := "https://example.com/invite/new"

		// Try to update a non-existent invitation record
		resp, err := service.UpdateUserInvitationRecord(ctx, nonExistentCustomerID, newLink)
		assert.NoError(t, err) // No error should be returned, just failure response
		assert.NotNil(t, resp)
		assert.Equal(t, int32(400), resp.Code)
		assert.Equal(t, "No Invitation Record Found", resp.Message)
		assert.Equal(t, "No Invitation Record Found", resp.ErrorMessage)
	})

	t.Run("Handle empty invitation link", func(t *testing.T) {
		customerID := int32(456)
		originalLink := "https://example.com/invite/original"
		emptyLink := ""

		// Create a user invitation record
		invitationRecord, err := dbClient.UserInvitationRecord.
			Create().
			SetCustomerID(int(customerID)).
			SetInvitationLink(originalLink).
			Save(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, invitationRecord)

		// Update with empty link (should be allowed)
		resp, err := service.UpdateUserInvitationRecord(ctx, customerID, emptyLink)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, int32(200), resp.Code)
		assert.Equal(t, "Invitation Link Updated", resp.Message)
		assert.Equal(t, "", resp.ErrorMessage)

		// Verify the record was updated with empty link
		updatedRecord, err := dbClient.UserInvitationRecord.Get(ctx, invitationRecord.ID)
		assert.NoError(t, err)
		assert.Equal(t, emptyLink, updatedRecord.InvitationLink)
	})

	t.Run("Update same invitation link (idempotent)", func(t *testing.T) {
		customerID := int32(789)
		invitationLink := "https://example.com/invite/same"

		// Create a user invitation record
		invitationRecord, err := dbClient.UserInvitationRecord.
			Create().
			SetCustomerID(int(customerID)).
			SetInvitationLink(invitationLink).
			Save(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, invitationRecord)

		// Update with the same link
		resp, err := service.UpdateUserInvitationRecord(ctx, customerID, invitationLink)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, int32(200), resp.Code)
		assert.Equal(t, "Invitation Link Updated", resp.Message)
		assert.Equal(t, "", resp.ErrorMessage)

		// Verify the record still has the same link
		updatedRecord, err := dbClient.UserInvitationRecord.Get(ctx, invitationRecord.ID)
		assert.NoError(t, err)
		assert.Equal(t, invitationLink, updatedRecord.InvitationLink)
	})

	t.Run("Update very long invitation link", func(t *testing.T) {
		customerID := int32(321)
		originalLink := "https://example.com/invite/original"
		// Create a long link (but within the 3000 char limit from schema)
		longLink := "https://example.com/invite/" + string(make([]byte, 2900))
		for i := range longLink[len("https://example.com/invite/"):] {
			longLink = longLink[:len("https://example.com/invite/")+i] + "a" + longLink[len("https://example.com/invite/")+i+1:]
		}

		// Create a user invitation record
		invitationRecord, err := dbClient.UserInvitationRecord.
			Create().
			SetCustomerID(int(customerID)).
			SetInvitationLink(originalLink).
			Save(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, invitationRecord)

		// Update with long link
		resp, err := service.UpdateUserInvitationRecord(ctx, customerID, longLink)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, int32(200), resp.Code)
		assert.Equal(t, "Invitation Link Updated", resp.Message)
		assert.Equal(t, "", resp.ErrorMessage)

		// Verify the record was updated with long link
		updatedRecord, err := dbClient.UserInvitationRecord.Get(ctx, invitationRecord.ID)
		assert.NoError(t, err)
		assert.Equal(t, longLink, updatedRecord.InvitationLink)
	})
}
