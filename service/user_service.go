package service

import (
	"context"
	"coresamples/common"
	"coresamples/ent"
	"coresamples/ent/contact"
	"coresamples/ent/internaluser"
	"coresamples/ent/loginhistory"
	"coresamples/ent/user"
	"coresamples/ent/userinvitationrecord"
	pb "coresamples/proto"
	"coresamples/publisher"
	"coresamples/util"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
)

type IUserService interface {
	GetUserInfoByRole(ctx context.Context, userID string, userRole string) (*pb.GetUserInfoByRoleResponse, error)
	GetLoginHistory(ctx context.Context, customerID, userID string, startTime, endTime *time.Time, perPage, page int32) (*pb.GetLoginHistoryResponse, error)
	GetInternalUser(ctx context.Context, role string, roleIDs []int32, usernames []string) (*pb.GetInternalUserResponse, error)
	GetInternalUserByRoleID(ctx context.Context, role string, roleID int32) ([]*pb.InternalUser, error)
	GetInternalUserByUsername(ctx context.Context, role string, username string) ([]*pb.InternalUser, error)
	GetInternalUserByID(ctx context.Context, internalUserID int32) (*pb.InternalUser, error)
	TransferSalesCustomer(ctx context.Context, fromSalesID, toSalesID, customerID string) (*pb.TransferSalesCustomerResponse, error)
	GetUserInformation(ctx context.Context, userID string) (*pb.GetUserInformationResponse, error)
	IsEmailUsedAsLoginId(ctx context.Context, email string) (*pb.IsEmailUsedAsLoginIdResponse, error)
	CheckWhetherEmailIsUsedAsLoginId(ctx context.Context, email string, clinicID string) (*pb.CheckWhetherEmailIsUsedAsLoginIdResponse, error)
	GetUser2FAContactInfo(ctx context.Context, token string) (*pb.GetUser2FAContactInfoResponse, error)
	RenewToken(ctx context.Context, jwtToken string) (*pb.RenewTokenResponse, error)
	Send2FAVerificationCode(ctx context.Context, username string, emailAddress, phoneNumber string) (*pb.Send2FAVerificationCodeResponse, error)
	Verify2FAVerificationCode(ctx context.Context, username, verificationCode, emailAddress, phoneNumber string) (*pb.Verify2FAVerificationResponse, error)
	InitialForgetPassword(ctx context.Context, emailAddress string) (*pb.ForgetPasswordResponse, error)
	ForgetPasswordRequest(ctx context.Context, username, requestMethod, requestTarget string) (*pb.ForgetPasswordRequestResponse, error)
	ForgetPassword(ctx context.Context, username, verificationCode, newPassword string) (*pb.ForgetPasswordVerifyResponse, error)
	TurnOff2FASettingPage(ctx context.Context, username, token string) (*pb.TurnOff2FASettingPageResponse, error)
	TurnOn2FASettingPage(ctx context.Context, username, email2faAddress, phone2faNumber, token string) (*pb.TurnOn2FASettingPageResponse, error)
	UpdateUserInvitationRecord(ctx context.Context, customerID int32, invitationLink string) (*pb.UpdateUserInvitationRecordResponse, error)
}

type UserService struct {
	dbClient    *ent.Client
	redisClient *common.RedisClient
	httpClient  *http.Client // Optional HTTP client for testing
}

func NewUserService(dbClient *ent.Client, redisClient *common.RedisClient) IUserService {
	return &UserService{
		dbClient:    dbClient,
		redisClient: redisClient,
	}
}

// NewUserServiceWithHTTPClient creates a new UserService with an injected HTTP client for testing
func NewUserServiceWithHTTPClient(dbClient *ent.Client, redisClient *common.RedisClient, httpClient *http.Client) IUserService {
	return &UserService{
		dbClient:    dbClient,
		redisClient: redisClient,
		httpClient:  httpClient,
	}
}

func GetUserService(dbClient *ent.Client, redisClient *common.RedisClient) IUserService {
	if UserSvc == nil {
		UserSvc = NewUserService(dbClient, redisClient)
	}
	return UserSvc
}

func GetCurrentUserService() IUserService {
	if UserSvc == nil {
		common.Fatal(ErrServiceNotInitialized)
	}
	return UserSvc
}

// convertToProto converts an ent.InternalUser to pb.InternalUser
func convertToProto(user *ent.InternalUser) *pb.InternalUser {
	return &pb.InternalUser{
		InternalUserId:         int32(user.ID),
		InternalUserRole:       user.InternalUserRole,
		InternalUserName:       user.InternalUserName,
		InternalUserFirstname:  user.InternalUserFirstname,
		InternalUserLastname:   user.InternalUserLastname,
		InternalUserMiddlename: user.InternalUserMiddleName,
		InternalUserRoleId:     int32(user.InternalUserRoleID),
		InternalUserType:       user.InternalUserType,
		InternalUserIsFullTime: user.InternalUserIsFullTime,
		InternalUserEmail:      user.InternalUserEmail,
		InternalUserPhone:      user.InternalUserPhone,
		IsActive:               user.IsActive,
		UserId:                 int32(user.UserID),
	}
}

// GetInternalUserByRoleID retrieves internal users based on role and role ID
func (s *UserService) GetInternalUserByRoleID(ctx context.Context, role string, roleID int32) ([]*pb.InternalUser, error) {
	query := s.dbClient.InternalUser.Query().Where(internaluser.IsActive(true))

	if role != "" {
		query = query.Where(internaluser.InternalUserRole(role))
	}

	query = query.Where(internaluser.InternalUserRoleID(int(roleID)))

	users, err := query.All(ctx)
	if err != nil {
		common.Errorf("Error fetching internal users by role ID", err)
		return nil, err
	}

	// Convert to proto message
	protoUsers := make([]*pb.InternalUser, 0, len(users))
	for _, user := range users {
		protoUsers = append(protoUsers, convertToProto(user))
	}

	return protoUsers, nil
}

// GetInternalUserByUsername retrieves internal users based on role and username with Redis caching
func (s *UserService) GetInternalUserByUsername(ctx context.Context, role string, username string) ([]*pb.InternalUser, error) {
	// Try to get from Redis cache first
	cacheKey := fmt.Sprintf("lis::core_service::internal_user_role_%s_username_%s", role, username)
	cachedResult := s.redisClient.Get(ctx, cacheKey)
	cachedData, err := cachedResult.Result()

	if err == nil && cachedData != "" {
		// Cache hit, deserialize and return
		var protoUsers []*pb.InternalUser
		if err := json.Unmarshal([]byte(cachedData), &protoUsers); err == nil {
			return protoUsers, nil
		}
		// If deserialize fails, continue to fetch from DB
	}

	// Cache miss or deserialize error, query the database
	query := s.dbClient.InternalUser.Query().Where(internaluser.IsActive(true))

	if role != "" {
		query = query.Where(internaluser.InternalUserRole(role))
	}

	query = query.Where(internaluser.InternalUserName(username))

	users, err := query.All(ctx)
	if err != nil {
		common.Errorf("Error fetching internal users by username", err)
		return nil, err
	}

	// Convert to proto message
	protoUsers := make([]*pb.InternalUser, 0, len(users))
	for _, user := range users {
		protoUsers = append(protoUsers, convertToProto(user))
	}

	// Cache the result (10 hours = 36000 seconds)
	if cacheData, err := json.Marshal(protoUsers); err == nil {
		s.redisClient.Set(ctx, cacheKey, string(cacheData), 36000*time.Second)
	}

	return protoUsers, nil
}

// GetInternalUserByID retrieves a single internal user by ID
func (s *UserService) GetInternalUserByID(ctx context.Context, internalUserID int32) (*pb.InternalUser, error) {
	// Directly query the database
	user, err := s.dbClient.InternalUser.Get(ctx, int(internalUserID))
	if err != nil {
		common.Errorf(fmt.Sprintf("Error fetching internal user by ID %d", internalUserID), err)
		return nil, err
	}

	// Convert to proto message
	protoUser := convertToProto(user)

	return protoUser, nil
}

// GetUserInfoByRole retrieves user information based on the user's role
func (s *UserService) GetUserInfoByRole(ctx context.Context, userID string, userRole string) (*pb.GetUserInfoByRoleResponse, error) {
	uid, err := strconv.Atoi(userID)
	if err != nil {
		common.Errorf("Invalid user ID format", err)
		return nil, err
	}

	// Query user with related entities
	userObj, err := s.dbClient.User.Query().
		Where(user.ID(uid)).
		WithClinic().
		WithCustomer().
		WithPatient().
		WithInternalUser().
		Only(ctx)

	if err != nil {
		common.Errorf(fmt.Sprintf("Error fetching user with ID %d", uid), err)
		return nil, err
	}

	// Prepare response
	response := &pb.GetUserInfoByRoleResponse{
		UserId:         int32(userObj.ID),
		Username:       userObj.UserName,
		UserGroup:      userObj.UserGroup,
		UserPermission: "", // Fill this if available in your schema
		IsActive:       userObj.IsActive,
	}

	// Populate relations based on user role
	if userObj.Edges.Customer != nil && len(userObj.Edges.Customer) > 0 && (userRole == "" || userRole == "customer") {
		customer := userObj.Edges.Customer[0]
		response.Customer = &pb.UserCustomer{
			CustomerId:         int32(customer.ID),
			CustomerFirstName:  customer.CustomerFirstName,
			CustomerLastName:   customer.CustomerLastName,
			CustomerMiddleName: customer.CustomerMiddleName,
		}
	}

	if userObj.Edges.Clinic != nil && len(userObj.Edges.Clinic) > 0 && (userRole == "" || userRole == "clinic") {
		clinic := userObj.Edges.Clinic[0]
		response.Clinic = &pb.UserClinic{
			ClinicId:   int32(clinic.ID),
			ClinicName: clinic.ClinicName,
		}
	}

	if userObj.Edges.InternalUser != nil && len(userObj.Edges.InternalUser) > 0 && (userRole == "" || userRole == "navigator" || userRole == "internal") {
		internalUser := userObj.Edges.InternalUser[0]
		response.Internal = &pb.UserInternal{
			InternalUserId:         int32(internalUser.ID),
			InternalUserRole:       internalUser.InternalUserRole,
			InternalUserRoleId:     int32(internalUser.InternalUserRoleID),
			InternalUserFirstname:  internalUser.InternalUserFirstname,
			InternalUserLastname:   internalUser.InternalUserLastname,
			InternalUserMiddlename: internalUser.InternalUserMiddleName,
		}
	}

	if userObj.Edges.Patient != nil && len(userObj.Edges.Patient) > 0 && (userRole == "" || userRole == "patient") {
		patient := userObj.Edges.Patient[0]
		response.Patient = &pb.UserPatient{
			PatientId:         int32(patient.ID),
			PatientFirstName:  patient.PatientFirstName,
			PatientLastName:   patient.PatientLastName,
			PatientMiddleName: patient.PatientMiddleName,
		}
	}

	return response, nil
}

// GetLoginHistory retrieves login history for a user
func (s *UserService) GetLoginHistory(ctx context.Context, customerID, userID string, startTime, endTime *time.Time, perPage, page int32) (*pb.GetLoginHistoryResponse, error) {
	// Set default pagination values if not provided
	if perPage <= 0 {
		perPage = 100
	}
	if page <= 0 {
		page = 1
	}

	var username string

	// Find the username based on customer_id or user_id
	if customerID != "" {
		cid, err := strconv.Atoi(customerID)
		if err != nil {
			common.Errorf("Invalid customer ID format", err)
			return nil, err
		}

		// Get customer's user
		customer, err := s.dbClient.Customer.Get(ctx, cid)
		if err != nil {
			common.Errorf(fmt.Sprintf("Error fetching customer with ID %d", cid), err)
			return nil, err
		}

		// Get username from user
		userObj, err := s.dbClient.User.Get(ctx, customer.UserID)
		if err != nil {
			common.Errorf(fmt.Sprintf("Error fetching user with ID %d", customer.UserID), err)
			return nil, err
		}

		username = userObj.UserName
	} else if userID != "" {
		uid, err := strconv.Atoi(userID)
		if err != nil {
			common.Errorf("Invalid user ID format", err)
			return nil, err
		}

		// Get username directly from user
		userObj, err := s.dbClient.User.Get(ctx, uid)
		if err != nil {
			common.Errorf(fmt.Sprintf("Error fetching user with ID %d", uid), err)
			return nil, err
		}

		username = userObj.UserName
	} else {
		return &pb.GetLoginHistoryResponse{
			LoginHistory: []*pb.LoginHistory{},
			TotalCount:   0,
		}, nil
	}

	// Build login history query
	query := s.dbClient.LoginHistory.Query().
		Where(loginhistory.Username(username))

	// Add time range filters if provided
	if startTime != nil {
		query = query.Where(loginhistory.LoginTimeGTE(*startTime))
	}
	if endTime != nil {
		query = query.Where(loginhistory.LoginTimeLTE(*endTime))
	}

	// Get total count
	totalCount, err := query.Count(ctx)
	if err != nil {
		common.Errorf("Error counting login history records", err)
		return nil, err
	}

	// Get paginated results
	skip := (page - 1) * perPage
	limit := perPage

	loginRecords, err := query.
		Order(ent.Desc(loginhistory.FieldLoginTime)).
		Offset(int(skip)).
		Limit(int(limit)).
		All(ctx)

	if err != nil {
		common.Errorf("Error fetching login history records", err)
		return nil, err
	}

	// Convert to protobuf response
	response := &pb.GetLoginHistoryResponse{
		LoginHistory: make([]*pb.LoginHistory, len(loginRecords)),
		TotalCount:   int32(totalCount),
	}

	for i, record := range loginRecords {
		response.LoginHistory[i] = &pb.LoginHistory{
			Id:                int32(record.ID),
			Username:          record.Username,
			LoginTime:         record.LoginTime.Format(time.RFC3339),
			LoginIp:           record.LoginIP,
			LoginSuccessfully: record.LoginSuccessfully,
			FailureReason:     record.FailureReason,
			LoginPortal:       record.LoginPortal,
		}
	}

	return response, nil
}

func (s *UserService) GetInternalUser(ctx context.Context, role string, roleIDs []int32, usernames []string) (*pb.GetInternalUserResponse, error) {
	response := &pb.GetInternalUserResponse{
		Response: []*pb.GetInternalUserResponseMiddleLevel{},
	}

	// Process role IDs - make separate queries for each role ID
	if len(roleIDs) > 0 {
		for _, roleID := range roleIDs {
			protoUsers, err := s.GetInternalUserByRoleID(ctx, role, roleID)
			if err != nil {
				common.Errorf(fmt.Sprintf("Error fetching internal users by role ID %d", roleID), err)
				continue // Skip this roleID but continue with others
			}

			// Add as separate entry in response
			middleLevel := &pb.GetInternalUserResponseMiddleLevel{
				InternalUser: protoUsers,
			}
			response.Response = append(response.Response, middleLevel)
		}
	}

	// Process usernames - make separate queries for each username, with Redis caching
	if len(usernames) > 0 {
		for _, username := range usernames {
			protoUsers, err := s.GetInternalUserByUsername(ctx, role, username)
			if err != nil {
				common.Errorf(fmt.Sprintf("Error fetching internal users by username %s", username), err)
				continue // Skip this username but continue with others
			}

			// Add as separate entry in response
			middleLevel := &pb.GetInternalUserResponseMiddleLevel{
				InternalUser: protoUsers,
			}
			response.Response = append(response.Response, middleLevel)
		}
	}

	// If no specific filters provided, return all active users matching the role
	if len(roleIDs) == 0 && len(usernames) == 0 {
		query := s.dbClient.InternalUser.Query().Where(internaluser.IsActive(true))

		if role != "" {
			query = query.Where(internaluser.InternalUserRole(role))
		}

		users, err := query.All(ctx)
		if err != nil {
			common.Errorf("Error fetching all internal users", err)
			return nil, err
		}

		// Convert to proto message
		protoUsers := make([]*pb.InternalUser, 0, len(users))
		for _, user := range users {
			protoUsers = append(protoUsers, convertToProto(user))
		}

		// Add as single entry
		middleLevel := &pb.GetInternalUserResponseMiddleLevel{
			InternalUser: protoUsers,
		}
		response.Response = append(response.Response, middleLevel)
	}

	return response, nil
}

// TransferSalesCustomer transfers a customer from one sales user to another sales user
func (s *UserService) TransferSalesCustomer(ctx context.Context, fromSalesID, toSalesID, customerID string) (*pb.TransferSalesCustomerResponse, error) {
	// Generate a tracking ID for logging
	trackingID := util.GenerateUUID()

	// Validate and convert the input IDs to integers
	custID, err := strconv.Atoi(customerID)
	if err != nil {
		common.Errorf(fmt.Sprintf("Invalid customer ID format: %s", customerID), err)
		return &pb.TransferSalesCustomerResponse{
			Status: fmt.Sprintf("Failed, Invalid customer ID: %s", customerID),
		}, nil
	}

	fromSID, err := strconv.Atoi(fromSalesID)
	if err != nil {
		common.Errorf(fmt.Sprintf("Invalid from_sales_id format: %s", fromSalesID), err)
		return &pb.TransferSalesCustomerResponse{
			Status: fmt.Sprintf("Failed, Invalid from_sales_id: %s", fromSalesID),
		}, nil
	}

	toSID, err := strconv.Atoi(toSalesID)
	if err != nil {
		common.Errorf(fmt.Sprintf("Invalid to_sales_id format: %s", toSalesID), err)
		return &pb.TransferSalesCustomerResponse{
			Status: fmt.Sprintf("Failed, Invalid to_sales_id: %s", toSalesID),
		}, nil
	}

	// Find the customer
	customer, err := s.dbClient.Customer.Get(ctx, custID)
	if err != nil {
		common.Errorf(fmt.Sprintf("Error finding customer with ID %d", custID), err)
		return &pb.TransferSalesCustomerResponse{
			Status: fmt.Sprintf("Failed, Customer with customer_id %s not found", customerID),
		}, nil
	}

	// Check if the customer is already assigned to the target sales
	if customer.SalesID == toSID {
		return &pb.TransferSalesCustomerResponse{
			Status: fmt.Sprintf("Failed, Customer with customer_id %s is already under sales %s", customerID, toSalesID),
		}, nil
	}

	// Check if the customer is assigned to the source sales
	if customer.SalesID != fromSID {
		return &pb.TransferSalesCustomerResponse{
			Status: fmt.Sprintf("Failed, Customer with customer_id %s is not under sales %s", customerID, fromSalesID),
		}, nil
	}

	// Verify the target sales exists and is active
	targetSales, err := s.dbClient.InternalUser.Get(ctx, toSID)
	if err != nil {
		common.Errorf(fmt.Sprintf("Error finding internal user with ID %d", toSID), err)
		return &pb.TransferSalesCustomerResponse{
			Status: fmt.Sprintf("Failed, Cannot find the internal user %s", toSalesID),
		}, nil
	}

	if !targetSales.IsActive || targetSales.InternalUserRole != "sales" {
		return &pb.TransferSalesCustomerResponse{
			Status: fmt.Sprintf("Failed, Cannot find the internal user %s or is not active", toSalesID),
		}, nil
	}

	// Update the customer's sales ID
	_, err = s.dbClient.Customer.UpdateOneID(custID).
		SetSalesID(toSID).
		Save(ctx)

	if err != nil {
		common.Errorf(fmt.Sprintf("Error updating customer %d sales ID from %d to %d", custID, fromSID, toSID), err)
		return &pb.TransferSalesCustomerResponse{
			Status: fmt.Sprintf("Failed, Error updating customer: %v", err),
		}, nil
	}

	// Log the successful transfer
	common.Infof(fmt.Sprintf("[%s] Transfer Sales Customer: Customer with customer_id %s from sales %s to sales %s",
		trackingID, customerID, fromSalesID, toSalesID))

	return &pb.TransferSalesCustomerResponse{
		Status: fmt.Sprintf("Successfully transfer the customer %s from sales %s to sales %s", customerID, fromSalesID, toSalesID),
	}, nil
}

// GetUserInformation retrieves detailed information about a user based on user_id
func (s *UserService) GetUserInformation(ctx context.Context, userID string) (*pb.GetUserInformationResponse, error) {
	// Generate a tracking ID for logging
	trackingID := util.GenerateUUID()

	// Convert user_id from string to int
	uid, err := strconv.Atoi(userID)
	if err != nil {
		common.Errorf(fmt.Sprintf("[%s] Invalid user ID format: %s", trackingID, userID), err)
		return nil, err
	}

	// Query user with customer relation
	userObj, err := s.dbClient.User.Query().
		Where(user.ID(uid)).
		WithCustomer().
		Only(ctx)

	if err != nil {
		// Check if the error is a "not found" error
		if ent.IsNotFound(err) {
			common.Infof(fmt.Sprintf("[%s] User with ID %d not found", trackingID, uid))
			// Match TypeScript behavior by returning a response with null/zero user_id
			return &pb.GetUserInformationResponse{
				UserId: 0, // In TypeScript this would be null, but in protobuf we use 0
			}, nil
		}

		common.Errorf(fmt.Sprintf("[%s] Error fetching user with ID %d", trackingID, uid), err)
		return nil, err
	}

	// Prepare response
	response := &pb.GetUserInformationResponse{
		UserId:                           int32(userObj.ID),
		Username:                         userObj.UserName,
		EmailUserId:                      userObj.EmailUserID,
		IsTwoFactorAuthenticationEnabled: userObj.IsTwoFactorAuthenticationEnabled,
		UserPermission:                   "", // UserPermission field is deprecated, leave it blank
		IsActive:                         userObj.IsActive,
		ImportedUserWithSaltPassword:     userObj.ImportedUserWithSaltPassword,
	}

	// Include customer information if available
	if userObj.Edges.Customer != nil && len(userObj.Edges.Customer) > 0 {
		customer := userObj.Edges.Customer[0]
		response.Customer = &pb.GetUserInfoCustomer{
			CustomerId:        int32(customer.ID),
			CustomerFirstName: customer.CustomerFirstName,
			CustomerLastName:  customer.CustomerLastName,
		}
	}

	// Log the successful retrieval
	common.Infof(fmt.Sprintf("[%s] User information retrieved for user_id %s", trackingID, userID))

	return response, nil
}

func (s *UserService) IsEmailUsedAsLoginId(ctx context.Context, email string) (*pb.IsEmailUsedAsLoginIdResponse, error) {
	// Generate tracking ID for logging
	trackingID := util.GenerateUUID()

	// Log the request for auditing and debugging purposes
	common.Infof(fmt.Sprintf("[%s] Checking if email is used as login ID: %s", trackingID, email))

	// Query the database - check if user exists with this email as login ID
	user, err := s.dbClient.User.Query().
		Where(user.EmailUserID(email)).
		Only(ctx)

	// Handle the possible outcomes
	if err != nil {
		if ent.IsNotFound(err) {
			// No user found with this email as login ID
			common.Infof(fmt.Sprintf("[%s] Email %s not found as login ID", trackingID, email))
			return &pb.IsEmailUsedAsLoginIdResponse{
				UsedAsEmailLogId: false,
				Message:          "Email is not used as log in email id",
				UserId:           0, // Using 0 as null equivalent
			}, nil
		}
		// Some other error occurred while querying the database
		common.Errorf(fmt.Sprintf("Error checking if email %s is used as login ID: %v", email, err), err)
		return nil, fmt.Errorf("database error while checking email: %w", err)
	}

	// User with this email as login ID exists
	common.Infof(fmt.Sprintf("[%s] Email %s is used as login ID by user_id %d", trackingID, email, user.ID))
	return &pb.IsEmailUsedAsLoginIdResponse{
		UsedAsEmailLogId: true,
		Message:          "Email is already used as log in email id",
		UserId:           int32(user.ID),
	}, nil
}

func (s *UserService) CheckWhetherEmailIsUsedAsLoginId(ctx context.Context, email string, clinicID string) (*pb.CheckWhetherEmailIsUsedAsLoginIdResponse, error) {
	// Generate tracking ID for logging
	trackingID := util.GenerateUUID()

	// Log the request for auditing and debugging purposes
	common.Infof(fmt.Sprintf("[%s] Checking if email %s is used as login ID in clinic %s", trackingID, email, clinicID))

	// Default response is no existing user
	existingUser := false

	// Query the database - check if user exists with this email as login ID
	userObj, err := s.dbClient.User.Query().
		Where(user.EmailUserID(email)).
		WithCustomer(func(q *ent.CustomerQuery) {
			q.WithClinics()
		}).
		WithClinic().
		Only(ctx)

	// Handle the possible outcomes
	if err != nil {
		if !ent.IsNotFound(err) {
			// Some other error occurred while querying the database
			common.Errorf(fmt.Sprintf("Error checking if email %s is used as login ID: %v", email, err), err)
			return nil, fmt.Errorf("database error while checking email: %w", err)
		}

		// No user found with this email as login ID
		common.Infof(fmt.Sprintf("[%s] Email %s not found as login ID", trackingID, email))
		return &pb.CheckWhetherEmailIsUsedAsLoginIdResponse{
			ExistingUser: existingUser,
			Message:      "User Does Not Exist",
		}, nil
	}

	// User with this email as login ID exists
	common.Infof(fmt.Sprintf("[%s] Email %s is used as login ID by user_id %d", trackingID, email, userObj.ID))

	// Check if user has a customer relationship and if that customer is associated with the specified clinic
	if userObj.Edges.Customer != nil && len(userObj.Edges.Customer) > 0 {
		// Convert clinic ID to integer for comparison with database values
		clinicIDInt, err := strconv.Atoi(clinicID)
		if err != nil {
			common.Errorf(fmt.Sprintf("Invalid clinic ID format: %s", clinicID), err)
			return nil, fmt.Errorf("invalid clinic ID format: %w", err)
		}

		// Customer exists, check clinics
		for _, customer := range userObj.Edges.Customer {
			if customer.Edges.Clinics != nil {
				for _, clinic := range customer.Edges.Clinics {
					if clinic.ID == clinicIDInt {
						existingUser = true
						break
					}
				}
			}

			if existingUser {
				break
			}
		}
	}

	// Check if user is directly associated with the clinic
	if !existingUser && userObj.Edges.Clinic != nil && len(userObj.Edges.Clinic) > 0 {
		// Convert clinic ID to integer for comparison with database values
		clinicIDInt, err := strconv.Atoi(clinicID)
		if err != nil {
			common.Errorf(fmt.Sprintf("Invalid clinic ID format: %s", clinicID), err)
			return nil, fmt.Errorf("invalid clinic ID format: %w", err)
		}

		// Check direct clinic association
		for _, clinic := range userObj.Edges.Clinic {
			if clinic.ID == clinicIDInt {
				existingUser = true
				break
			}
		}
	}

	// Construct response based on findings
	if existingUser {
		common.Infof(fmt.Sprintf("[%s] Email %s is associated with clinic %s", trackingID, email, clinicID))
		return &pb.CheckWhetherEmailIsUsedAsLoginIdResponse{
			ExistingUser: true,
			Message:      "User Exists",
		}, nil
	} else {
		common.Infof(fmt.Sprintf("[%s] Email %s is not associated with clinic %s", trackingID, email, clinicID))
		return &pb.CheckWhetherEmailIsUsedAsLoginIdResponse{
			ExistingUser: false,
			Message:      "User Exists but not in this clinic",
		}, nil
	}
}

// GetUser2FAContactInfo retrieves all 2FA contacts for a user based on JWT token
// It returns a list of contacts that are marked as 2FA contacts
func (s *UserService) GetUser2FAContactInfo(ctx context.Context, token string) (*pb.GetUser2FAContactInfoResponse, error) {
	// Generate a tracking ID for logging
	trackingID := util.GenerateUUID()
	common.Infof(fmt.Sprintf("[%s] GetUser2FAContactInfo: Retrieving 2FA contacts for user", trackingID))

	// Parse the JWT token to get user ID
	claims, err := util.ParseJWTToken(token, common.Secrets.JWTSecret)
	if err != nil {
		common.Errorf(fmt.Sprintf("[%s] Failed to parse JWT token: %v", trackingID, err), err)
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	userID := claims.UserId
	if userID <= 0 {
		common.Errorf(fmt.Sprintf("[%s] Invalid user ID in token: %d", trackingID, userID), nil)
		return nil, fmt.Errorf("invalid user ID in token")
	}

	// Query the database for all contacts associated with this user that are marked as 2FA contacts
	contacts, err := s.dbClient.Contact.Query().
		Where(
			contact.UserID(userID),
			contact.Is2faContact(true),
		).
		All(ctx)

	if err != nil {
		common.Errorf(fmt.Sprintf("[%s] Error fetching 2FA contacts for user %d: %v", trackingID, userID, err), err)
		return nil, fmt.Errorf("error fetching 2FA contacts: %w", err)
	}

	// Convert to protobuf response
	response := &pb.GetUser2FAContactInfoResponse{
		Contacts: make([]*pb.User2FAContact, 0, len(contacts)),
	}

	for _, c := range contacts {
		response.Contacts = append(response.Contacts, &pb.User2FAContact{
			ContactId:                 int32(c.ID),
			ContactDescription:        c.ContactDescription,
			ContactDetails:            c.ContactDetails,
			ContactType:               c.ContactType,
			IsPrimaryContact:          c.IsPrimaryContact,
			Is_2FaContact:             c.Is2faContact,
			CustomerId:                int32(c.CustomerID),
			PatientId:                 int32(c.PatientID),
			ClinicId:                  int32(c.ClinicID),
			InternalUserId:            int32(c.InternalUserID),
			ContactLevel:              int32(c.ContactLevel),
			ContactLevelName:          c.ContactLevelName,
			GroupContactId:            int32(c.GroupContactID),
			ApplyToAllGroupMember:     c.ApplyToAllGroupMember,
			HasGroupContact:           c.GroupContactID > 0,
			IsGroupContact:            c.IsGroupContact,
			UseAsDefaultCreateContact: c.UseAsDefaultCreateContact,
			UseGroupContact:           c.UseGroupContact,
		})
	}

	common.Infof(fmt.Sprintf("[%s] GetUser2FAContactInfo: Found %d 2FA contacts for user %d",
		trackingID, len(response.Contacts), userID))
	return response, nil
}

// RenewToken validates and renews a JWT token
// It parses the original token, validates it, and issues a new token with the same claims
func (s *UserService) RenewToken(ctx context.Context, jwtToken string) (*pb.RenewTokenResponse, error) {
	// Generate tracking ID for logging
	trackingID := util.GenerateUUID()
	common.Infof(fmt.Sprintf("[%s] Processing token renewal request", trackingID))

	// Try to parse and validate the token
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(common.Secrets.JWTSecret), nil
	})

	// Handle token validation errors
	if err != nil {
		common.Errorf(fmt.Sprintf("[%s] Failed to parse token: %v", trackingID, err), err)
		return &pb.RenewTokenResponse{
			Code:           400,
			Message:        fmt.Sprintf("Token Renew Failed, error: %v", err),
			Token:          "",
			ExpirationTime: "",
		}, nil
	}

	// Extract claims from the token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Create a new token with the same claims
		newToken := jwt.New(jwt.SigningMethodHS256)
		newClaims := newToken.Claims.(jwt.MapClaims)

		// Copy all claims from the original token
		for key, val := range claims {
			if key != "exp" && key != "iat" { // Skip expiration and issuance time
				newClaims[key] = val
			}
		}

		// Set new issuance and expiration times
		currentTimestamp := time.Now().Unix()

		// Get JWT expiration time from environment variable, default to 2700 seconds (45 minutes)
		// This matches the TypeScript implementation's behavior
		var expirationTimeInSeconds int64 = 2700 // Default 2700 seconds (45 minutes)
		if expTimeStr := os.Getenv("JWT_EXPIRATION_TIME"); expTimeStr != "" {
			if expTime, err := strconv.ParseInt(expTimeStr, 10, 64); err == nil {
				expirationTimeInSeconds = expTime
				common.Infof(fmt.Sprintf("[%s] Using JWT_EXPIRATION_TIME from environment: %d seconds", trackingID, expirationTimeInSeconds))
			} else {
				common.Errorf(fmt.Sprintf("[%s] Invalid JWT_EXPIRATION_TIME value, using default: 2700 seconds", trackingID), err)
			}
		} else {
			common.Infof(fmt.Sprintf("[%s] JWT_EXPIRATION_TIME not set, using default: 2700 seconds", trackingID))
		}

		tokenExpirationTimestamp := currentTimestamp + expirationTimeInSeconds

		newClaims["iat"] = currentTimestamp
		newClaims["exp"] = tokenExpirationTimestamp

		// Sign the new token
		tokenString, err := newToken.SignedString([]byte(common.Secrets.JWTSecret))
		if err != nil {
			common.Errorf(fmt.Sprintf("[%s] Failed to sign new token: %v", trackingID, err), err)
			return &pb.RenewTokenResponse{
				Code:           400,
				Message:        fmt.Sprintf("Token Renew Failed, error: %v", err),
				Token:          "",
				ExpirationTime: "",
			}, nil
		}

		// Create audit log message
		auditLogMessage := common.AuditLogEntry{
			EventID:     trackingID,
			ServiceName: common.ServiceName,
			ServiceType: "backend",
			EventName:   "RenewToken",
			EntityType:  "user",
			EntityID:    fmt.Sprintf("%v", claims["userId"]),
			User:        fmt.Sprintf("%v", claims["userId"]),
			Entrypoint:  "GRPC",
		}
		go func() {
			common.RecordAuditLog(auditLogMessage)
		}()

		// Return successful response
		return &pb.RenewTokenResponse{
			Code:           200,
			Message:        "Token Renewed",
			Token:          tokenString,
			ExpirationTime: fmt.Sprintf("%d", tokenExpirationTimestamp),
		}, nil
	}

	// If we get here, the token was invalid for some reason
	common.Errorf(fmt.Sprintf("[%s] Invalid token claims", trackingID), nil)
	return &pb.RenewTokenResponse{
		Code:           400,
		Message:        "Token Renew Failed, invalid token claims",
		Token:          "",
		ExpirationTime: "",
	}, nil
}

// Send2FAVerificationCode sends a 2FA verification code via email or SMS
func (s *UserService) Send2FAVerificationCode(ctx context.Context, username string, emailAddress, phoneNumber string) (*pb.Send2FAVerificationCodeResponse, error) {
	// Generate a tracking ID for logging
	trackingID := util.GenerateUUID()
	common.Infof(fmt.Sprintf("[%s] Send2FAVerificationCode: Starting for username %s", trackingID, username))

	// Find the user by email or username
	var userObj *ent.User
	var err error

	// Check if username is an email format
	if util.IsValidEmail(username) {
		userObj, err = s.dbClient.User.Query().
			Where(user.EmailUserID(username)).
			Only(ctx)

		if err != nil && !ent.IsNotFound(err) {
			common.Errorf(fmt.Sprintf("[%s] Error querying user by email: %v", trackingID, err), err)
			return nil, err
		}
	}

	// If not found by email, try username
	if userObj == nil {
		userObj, err = s.dbClient.User.Query().
			Where(user.UserName(username)).
			Only(ctx)

		if err != nil {
			if ent.IsNotFound(err) {
				common.Infof(fmt.Sprintf("[%s] User not found: %s", trackingID, username))
				return &pb.Send2FAVerificationCodeResponse{
					Code:    404,
					Message: "User Not Found",
				}, nil
			}
			common.Errorf(fmt.Sprintf("[%s] Error querying user: %v", trackingID, err), err)
			return nil, err
		}
	}

	// Check if 2FA is enabled for the user
	if !userObj.IsTwoFactorAuthenticationEnabled {
		common.Infof(fmt.Sprintf("[%s] 2FA not enabled for user %d", trackingID, userObj.ID))
		return &pb.Send2FAVerificationCodeResponse{
			Code:    400,
			Message: "2FA is Not Enabled",
		}, nil
	}

	// Handle email verification
	if emailAddress != "" {
		// Query for user's 2FA contacts
		contacts, err := s.dbClient.Contact.Query().
			Where(
				contact.UserID(userObj.ID),
				contact.Is2faContact(true),
			).
			All(ctx)

		if err != nil {
			common.Errorf(fmt.Sprintf("[%s] Error querying contacts: %v", trackingID, err), err)
			return nil, err
		}

		// Check if user has 2FA contacts
		if len(contacts) == 0 {
			return &pb.Send2FAVerificationCodeResponse{
				Code:    400,
				Message: "This Account Does Not Have 2FA Contact, Please Contact Support",
			}, nil
		}

		// Find email contact
		var emailContact *ent.Contact
		for _, c := range contacts {
			if c.ContactType == "email" {
				emailContact = c
				break
			}
		}

		if emailContact == nil {
			return &pb.Send2FAVerificationCodeResponse{
				Code:    400,
				Message: "This Account is Not Set with 2FA Email",
			}, nil
		}

		// Verify the email matches
		if emailContact.ContactDetails != emailAddress {
			return &pb.Send2FAVerificationCodeResponse{
				Code:    400,
				Message: "The Input Email Address is Not the 2FA Email of This Account",
			}, nil
		}

		// Generate verification code using OTP library
		verificationCode := util.GenerateOTP(6)

		// Store in Redis with 5 minute expiration
		redisKey := fmt.Sprintf("lis::core_service::user_service:2fa%d_email", userObj.ID)
		err = s.redisClient.Set(ctx, redisKey, verificationCode, 300*time.Second).Err()
		if err != nil {
			common.Errorf(fmt.Sprintf("[%s] Error storing verification code in Redis: %v", trackingID, err), err)
			return nil, fmt.Errorf("error storing verification code: %w", err)
		}

		// Send email via external service
		err = s.send2FAEmailVerification(ctx, emailAddress, verificationCode)
		if err != nil {
			common.Errorf(fmt.Sprintf("[%s] Error sending verification email: %v", trackingID, err), err)
			return nil, fmt.Errorf("error sending verification email: %w", err)
		}

		common.Infof(fmt.Sprintf("[%s] Verification code sent via email to %s", trackingID, emailAddress))
		return &pb.Send2FAVerificationCodeResponse{
			Code:    200,
			Message: "Code sent via email",
		}, nil
	}

	// Handle phone verification
	if phoneNumber != "" {
		// Query for user's 2FA contacts
		contacts, err := s.dbClient.Contact.Query().
			Where(
				contact.UserID(userObj.ID),
				contact.Is2faContact(true),
			).
			All(ctx)

		if err != nil {
			common.Errorf(fmt.Sprintf("[%s] Error querying contacts: %v", trackingID, err), err)
			return nil, err
		}

		// Check if user has 2FA contacts
		if len(contacts) == 0 {
			return &pb.Send2FAVerificationCodeResponse{
				Code:    400,
				Message: "This Account Does Not Have 2FA Contact, Please Contact Support",
			}, nil
		}

		// Find phone contact
		var phoneContact *ent.Contact
		for _, c := range contacts {
			if c.ContactType == "phone" {
				phoneContact = c
				break
			}
		}

		if phoneContact == nil {
			return &pb.Send2FAVerificationCodeResponse{
				Code:    400,
				Message: "This Account is Not Set with 2FA Phone",
			}, nil
		}

		// Verify the phone matches
		if phoneContact.ContactDetails != phoneNumber {
			return &pb.Send2FAVerificationCodeResponse{
				Code:    400,
				Message: "The Input Phone Number is Not the 2FA Phone of This Account",
			}, nil
		}

		// Generate verification code
		verificationCode := util.GenerateOTP(6)

		// Store in Redis with 5 minute expiration
		redisKey := fmt.Sprintf("lis::core_service::user_service:2fa%d_text", userObj.ID)
		err = s.redisClient.Set(ctx, redisKey, verificationCode, 300*time.Second).Err()
		if err != nil {
			common.Errorf(fmt.Sprintf("[%s] Error storing verification code in Redis: %v", trackingID, err), err)
			return nil, fmt.Errorf("error storing verification code: %w", err)
		}

		// Send SMS via Kafka
		smsID := util.GenerateUUID()
		textMessage := map[string]interface{}{
			"SMSID":    smsID,
			"Tag":      "send",
			"From":     "18776760002",
			"To":       phoneNumber,
			"TextBody": fmt.Sprintf("Vibrant America: %s is your security code. This code is valid for 5 minutes.", verificationCode),
			"Delay":    0,
			"Type":     "SMS",
		}

		// Publish to Kafka
		err = s.publishSMSMessage(ctx, smsID, textMessage)
		if err != nil {
			common.Errorf(fmt.Sprintf("[%s] Error sending SMS: %v", trackingID, err), err)
			return nil, fmt.Errorf("error sending SMS: %w", err)
		}

		common.Infof(fmt.Sprintf("[%s] Verification code sent via SMS to %s", trackingID, phoneNumber))
		return &pb.Send2FAVerificationCodeResponse{
			Code:    200,
			Message: "Code sent via text",
		}, nil
	}

	// Neither email nor phone provided
	return &pb.Send2FAVerificationCodeResponse{
		Code:    400,
		Message: "Please provide either email or phone number",
	}, nil
}

// send2FAEmailVerification sends a 2FA verification email via external service
func (s *UserService) send2FAEmailVerification(ctx context.Context, emailAddress string, verificationCode string) error {
	// Generate JWT token for API authentication
	token, err := util.GenerateSystemJWT(common.Secrets.JWTSecret)
	if err != nil {
		return fmt.Errorf("error generating JWT token: %w", err)
	}

	// Prepare request body
	requestBody := map[string]string{
		"email": emailAddress,
		"code":  verificationCode,
	}

	// Send HTTP request to email service
	endpoint := "https://www.vibrant-america.com/lisapi/v1/portal/trans-service/valogin/send2faAuthEmail"
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	// Use injected HTTP client if available, otherwise use util.PostJSON
	if s.httpClient != nil {
		return util.PostJSONWithClient(ctx, s.httpClient, endpoint, requestBody, headers)
	}
	return util.PostJSON(ctx, endpoint, requestBody, headers)
}

// publishSMSMessage publishes an SMS message to Kafka
func (s *UserService) publishSMSMessage(ctx context.Context, messageKey string, message map[string]interface{}) error {
	// Get the Kafka publisher
	publisher := publisher.GetPublisher()

	// Marshall message to JSON
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling SMS message: %w", err)
	}

	// Publish to Kafka
	return publisher.GetWriter().WriteMessages(ctx, kafka.Message{
		Topic: "Notification-SMS",
		Key:   []byte(messageKey),
		Value: messageBytes,
	})
}

// Verify2FAVerificationCode verifies a 2FA verification code sent via email or SMS
func (s *UserService) Verify2FAVerificationCode(ctx context.Context, username, verificationCode, emailAddress, phoneNumber string) (*pb.Verify2FAVerificationResponse, error) {
	// Generate a tracking ID for logging
	trackingID := util.GenerateUUID()
	common.Infof(fmt.Sprintf("[%s] Verify2FAVerificationCode: Starting for username %s", trackingID, username))

	// Find the user by email or username
	var userObj *ent.User
	var err error

	// Check if username is an email format
	if util.IsValidEmail(username) {
		userObj, err = s.dbClient.User.Query().
			Where(user.EmailUserID(username)).
			Only(ctx)

		if err != nil && !ent.IsNotFound(err) {
			common.Errorf(fmt.Sprintf("[%s] Error querying user by email: %v", trackingID, err), err)
			return nil, err
		}
	}

	// If not found by email, try username
	if userObj == nil {
		userObj, err = s.dbClient.User.Query().
			Where(user.UserName(username)).
			Only(ctx)

		if err != nil {
			if ent.IsNotFound(err) {
				common.Infof(fmt.Sprintf("[%s] User not found: %s", trackingID, username))
				return &pb.Verify2FAVerificationResponse{
					Code:    404,
					Message: "User Not Found",
				}, nil
			}
			common.Errorf(fmt.Sprintf("[%s] Error querying user: %v", trackingID, err), err)
			return nil, err
		}
	}

	// Check if 2FA is enabled for the user
	if !userObj.IsTwoFactorAuthenticationEnabled {
		common.Infof(fmt.Sprintf("[%s] 2FA not enabled for user %d", trackingID, userObj.ID))
		return &pb.Verify2FAVerificationResponse{
			Code:    400,
			Message: "2FA is Not Enabled",
		}, nil
	}

	// Determine if email or phone verification
	var redisKey string
	var verificationMethod string

	if emailAddress != "" {
		// Verify email address belongs to the user
		contacts, err := s.dbClient.Contact.Query().
			Where(
				contact.UserID(userObj.ID),
				contact.Is2faContact(true),
				contact.ContactType("email"),
			).
			All(ctx)

		if err != nil {
			common.Errorf(fmt.Sprintf("[%s] Error querying contacts: %v", trackingID, err), err)
			return nil, err
		}

		// Check if email matches any of the user's 2FA email contacts
		var emailMatched bool
		for _, c := range contacts {
			if c.ContactDetails == emailAddress {
				emailMatched = true
				break
			}
		}

		if !emailMatched {
			return &pb.Verify2FAVerificationResponse{
				Code:    400,
				Message: "The provided email does not match any 2FA emails for this account",
			}, nil
		}

		redisKey = fmt.Sprintf("lis::core_service::user_service:2fa%d_email", userObj.ID)
		verificationMethod = "email"
	} else if phoneNumber != "" {
		// Verify phone number belongs to the user
		contacts, err := s.dbClient.Contact.Query().
			Where(
				contact.UserID(userObj.ID),
				contact.Is2faContact(true),
				contact.ContactType("phone"),
			).
			All(ctx)

		if err != nil {
			common.Errorf(fmt.Sprintf("[%s] Error querying contacts: %v", trackingID, err), err)
			return nil, err
		}

		// Check if phone matches any of the user's 2FA phone contacts
		var phoneMatched bool
		for _, c := range contacts {
			if c.ContactDetails == phoneNumber {
				phoneMatched = true
				break
			}
		}

		if !phoneMatched {
			return &pb.Verify2FAVerificationResponse{
				Code:    400,
				Message: "The provided phone number does not match any 2FA phones for this account",
			}, nil
		}

		redisKey = fmt.Sprintf("lis::core_service::user_service:2fa%d_text", userObj.ID)
		verificationMethod = "text"
	} else {
		// Neither email nor phone provided
		return &pb.Verify2FAVerificationResponse{
			Code:    400,
			Message: "Please provide either email or phone number",
		}, nil
	}

	// Get the stored verification code from Redis
	storedCode, err := s.redisClient.Get(ctx, redisKey).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			// Code expired or never sent
			common.Infof(fmt.Sprintf("[%s] No verification code found or expired for user %d", trackingID, userObj.ID))
			return &pb.Verify2FAVerificationResponse{
				Code:    400,
				Message: "Verification code expired or not sent. Please request a new code.",
			}, nil
		}
		// Redis error
		common.Errorf(fmt.Sprintf("[%s] Error retrieving verification code from Redis: %v", trackingID, err), err)
		return nil, fmt.Errorf("error retrieving verification code: %w", err)
	}

	// Compare the codes
	if storedCode != verificationCode {
		common.Infof(fmt.Sprintf("[%s] Invalid verification code for user %d", trackingID, userObj.ID))
		return &pb.Verify2FAVerificationResponse{
			Code:    400,
			Message: "Invalid verification code",
		}, nil
	}

	// Code verified successfully, delete the code from Redis
	err = s.redisClient.Del(ctx, redisKey).Err()
	if err != nil {
		common.Errorf(fmt.Sprintf("[%s] Error deleting verification code from Redis: %v", trackingID, err), err)
		// Continue anyway, this is not a critical error
	}

	// Log the successful verification
	common.Infof(fmt.Sprintf("[%s] 2FA verification successful for user %d via %s", trackingID, userObj.ID, verificationMethod))

	return &pb.Verify2FAVerificationResponse{
		Code:    200,
		Message: "Verification successful",
	}, nil
}

// InitialForgetPassword initiates the forget password process by sending an email
func (s *UserService) InitialForgetPassword(ctx context.Context, emailAddress string) (*pb.ForgetPasswordResponse, error) {
	// Generate a tracking ID for logging
	trackingID := util.GenerateUUID()
	common.Infof(fmt.Sprintf("[%s] InitialForgetPassword: Starting for email %s", trackingID, emailAddress))

	// Defer catch-all to ensure we always return success (security requirement)
	defer func() {
		if r := recover(); r != nil {
			common.Errorf(fmt.Sprintf("[%s] InitialForgetPassword: Recovered from panic: %v", trackingID, r), nil)
		}
	}()

	// OPTIMIZATION: Enhanced input validation and sanitization
	emailAddress = strings.TrimSpace(emailAddress)
	if emailAddress == "" {
		common.Infof(fmt.Sprintf("[%s] Empty email address provided", trackingID))
		return &pb.ForgetPasswordResponse{
			Code:    400,
			Message: "Please Enter a Valid Email Address",
		}, nil
	}
	
	// 1. Input validation - ONLY accept email addresses with regex validation
	if !util.IsValidEmail(emailAddress) {
		common.Infof(fmt.Sprintf("[%s] Invalid email format: %s", trackingID, emailAddress))
		return &pb.ForgetPasswordResponse{
			Code:    400,
			Message: "Please Enter a Valid Email Address",
		}, nil
	}

	// 2. User lookup logic - Find user by email_user_id field only
	userObj, err := s.dbClient.User.Query().
		Where(user.EmailUserID(emailAddress)).
		WithCustomer().
		WithClinic().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			// User not found - return "Contact Support"
			common.Infof(fmt.Sprintf("[%s] User not found for email: %s", trackingID, emailAddress))
			return &pb.ForgetPasswordResponse{
				Code:    400,
				Message: "Contact Support",
			}, nil
		}
		// Database error - still return success message for security
		common.Errorf(fmt.Sprintf("[%s] Database error looking up user: %v", trackingID, err), err)
		return &pb.ForgetPasswordResponse{
			Code:    200,
			Message: "Email Address Sent",
		}, nil
	}

	// 3. Business logic flow - Handle customer vs clinic relationships
	var customerID, clinicID int

	// Check customer relationship
	if userObj.Edges.Customer != nil && len(userObj.Edges.Customer) > 0 {
		customerID = userObj.Edges.Customer[0].ID
	}

	// Check clinic relationship
	if userObj.Edges.Clinic != nil && len(userObj.Edges.Clinic) > 0 {
		clinicID = userObj.Edges.Clinic[0].ID
	}

	// 4. Get 2FA contact information (is_2fa_contact: true)
	contacts, err := s.dbClient.Contact.Query().
		Where(
			contact.UserID(userObj.ID),
			contact.Is2faContact(true),
		).
		All(ctx)

	if err != nil {
		common.Errorf(fmt.Sprintf("[%s] Error getting 2FA contacts: %v", trackingID, err), err)
		// Continue with empty contacts - don't fail the process
		contacts = []*ent.Contact{}
	}

	// Extract 2FA email and phone if available
	var email2FA, phone2FA string
	for _, c := range contacts {
		if c.ContactType == "email" {
			email2FA = c.ContactDetails
		} else if c.ContactType == "phone" {
			phone2FA = c.ContactDetails
		}
	}

	// 5. Generate JWT token using system token generation
	token, err := util.GenerateSystemJWT(common.Secrets.JWTSecret)
	if err != nil {
		common.Errorf(fmt.Sprintf("[%s] Error generating JWT token: %v", trackingID, err), err)
		return &pb.ForgetPasswordResponse{
			Code:    200,
			Message: "Email Address Sent",
		}, nil
	}

	// 6. Build request body exactly matching TypeScript implementation
	requestBody := map[string]interface{}{
		"is_2fa":          len(contacts) > 0,
		"email_2fa":       email2FA,
		"phone_2fa":       phone2FA,
		"email":           emailAddress,
		"token":           token,
		"email_log_in_id": emailAddress,
	}

	// Add customer_id or clinic_id based on relationship
	if customerID > 0 {
		requestBody["customer_id"] = customerID
	}
	if clinicID > 0 {
		requestBody["clinic_id"] = clinicID
	}

	// 7. Email sending - HTTP POST to external API
	endpoint := "https://www.vibrant-america.com/lisapi/v1/portal/trans-service/valogin/sendForgetPasswordEmail"
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	// Send the request (errors are caught but don't affect response)
	err = func() error {
		if s.httpClient != nil {
			return util.PostJSONWithClient(ctx, s.httpClient, endpoint, requestBody, headers)
		}
		return util.PostJSON(ctx, endpoint, requestBody, headers)
	}()

	if err != nil {
		common.Errorf(fmt.Sprintf("[%s] Error sending forget password email: %v", trackingID, err), err)
		// Don't fail - continue to return success for security
	} else {
		common.Infof(fmt.Sprintf("[%s] Forget password email sent successfully", trackingID))
	}

	// 8. Response pattern - Always return success for security
	return &pb.ForgetPasswordResponse{
		Code:    200,
		Message: "Email Address Sent",
	}, nil
}

// ForgetPasswordRequest handles password reset request by sending a verification code
// This method checks if a user exists and has 2FA enabled, then sends a verification code via email or SMS
func (s *UserService) ForgetPasswordRequest(ctx context.Context, username, requestMethod, requestTarget string) (*pb.ForgetPasswordRequestResponse, error) {
	// Generate a tracking ID for logging
	trackingID := util.GenerateUUID()
	common.Infof(fmt.Sprintf("[%s] ForgetPasswordRequest: Starting for username %s", trackingID, username))

	// OPTIMIZATION: Enhanced input validation and sanitization
	username = strings.TrimSpace(username)
	requestMethod = strings.TrimSpace(requestMethod)
	requestTarget = strings.TrimSpace(requestTarget)
	
	if username == "" {
		common.Errorf(fmt.Sprintf("[%s] Empty username provided", trackingID), fmt.Errorf("empty username"))
		return &pb.ForgetPasswordRequestResponse{
			Code:    400,
			Message: "Username is required",
		}, nil
	}
	
	if requestMethod == "" {
		common.Errorf(fmt.Sprintf("[%s] Empty request method provided", trackingID), fmt.Errorf("empty request method"))
		return &pb.ForgetPasswordRequestResponse{
			Code:    400,
			Message: "Request method is required",
		}, nil
	}
	
	if requestTarget == "" {
		common.Errorf(fmt.Sprintf("[%s] Empty request target provided", trackingID), fmt.Errorf("empty request target"))
		return &pb.ForgetPasswordRequestResponse{
			Code:    400,
			Message: "Request target is required",
		}, nil
	}

	// Validate request method
	if requestMethod != "email" && requestMethod != "phone" {
		common.Errorf(fmt.Sprintf("[%s] Invalid request method: %s", trackingID, requestMethod), fmt.Errorf("invalid request method"))
		return &pb.ForgetPasswordRequestResponse{
			Code:    400,
			Message: "Invalid request method. Must be 'email' or 'phone'",
		}, nil
	}

	// Validate request target format based on method
	if requestMethod == "email" && !util.IsValidEmail(requestTarget) {
		common.Errorf(fmt.Sprintf("[%s] Invalid email format: %s", trackingID, requestTarget), fmt.Errorf("invalid email"))
		return &pb.ForgetPasswordRequestResponse{
			Code:    400,
			Message: "Invalid email format",
		}, nil
	}

	// Generate SMS ID for potential SMS sending
	smsID := util.GenerateUUID()

	var userObj *ent.User
	var err error

	// Check if username is an email address
	if util.IsValidEmail(username) {
		userObj, err = s.dbClient.User.Query().
			Where(user.EmailUserID(username)).
			WithCustomer().
			WithInternalUser().
			WithPatient().
			WithClinic().
			Only(ctx)

		if err != nil && !ent.IsNotFound(err) {
			common.Errorf(fmt.Sprintf("[%s] Error querying user by email: %v", trackingID, err), err)
			return nil, err
		}
	}

	// If not found by email, try username
	if userObj == nil {
		userObj, err = s.dbClient.User.Query().
			Where(user.UserName(username)).
			WithCustomer().
			WithInternalUser().
			WithPatient().
			WithClinic().
			Only(ctx)

		if err != nil {
			if ent.IsNotFound(err) {
				common.Infof(fmt.Sprintf("[%s] User not found: %s", trackingID, username))
				return &pb.ForgetPasswordRequestResponse{
					Code:    404,
					Message: "User not found",
				}, nil
			}
			common.Errorf(fmt.Sprintf("[%s] Error querying user: %v", trackingID, err), err)
			return nil, err
		}
	}

	// Check if the user has 2FA turned on
	if !userObj.IsTwoFactorAuthenticationEnabled {
		common.Infof(fmt.Sprintf("[%s] User does not have 2FA enabled: %d", trackingID, userObj.ID))
		return &pb.ForgetPasswordRequestResponse{
			Code:    401,
			Message: "User does not have 2fa enabled",
		}, nil
	}

	// Generate verification code (6 digits)
	verificationCode := util.GenerateOTP(6)

	// Store verification code in Redis with 5 minute expiration
	redisKey := fmt.Sprintf("lis::core_service::user_service:forget_password_verification_user_%s", username)
	err = s.redisClient.Set(ctx, redisKey, verificationCode, 300*time.Second).Err()
	if err != nil {
		common.Errorf(fmt.Sprintf("[%s] Error storing verification code in Redis: %v", trackingID, err), err)
		return nil, fmt.Errorf("error storing verification code: %w", err)
	}

	// Send out the verification code via phone/email
	switch requestMethod {
	case "email":
		// Generate system JWT token for API authentication
		token, err := util.GenerateSystemJWT(common.Secrets.JWTSecret)
		if err != nil {
			common.Errorf(fmt.Sprintf("[%s] Error generating JWT token: %v", trackingID, err), err)
			return nil, fmt.Errorf("error generating JWT token: %w", err)
		}

		// Prepare request body for email service
		requestBody := map[string]string{
			"email": requestTarget,
			"code":  verificationCode,
		}

		// Send email via external service
		endpoint := "https://www.vibrant-america.com/lisapi/v1/portal/trans-service/valogin/send2faResetPassEmail"
		headers := map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", token),
		}

		// Send HTTP request to email service
		err = func() error {
			if s.httpClient != nil {
				return util.PostJSONWithClient(ctx, s.httpClient, endpoint, requestBody, headers)
			}
			return util.PostJSON(ctx, endpoint, requestBody, headers)
		}()

		if err != nil {
			common.Errorf(fmt.Sprintf("[%s] Error sending reset password email: %v", trackingID, err), err)
			return nil, fmt.Errorf("error sending reset password email: %w", err)
		}

		common.Infof(fmt.Sprintf("[%s] Verification code sent via email to %s", trackingID, requestTarget))
		return &pb.ForgetPasswordRequestResponse{
			Code:    200,
			Message: "Code sent via email",
		}, nil

	case "phone":
		// Send SMS via Kafka
		textMessage := map[string]interface{}{
			"SMSID":    smsID,
			"Tag":      "send",
			"From":     "18776760002",
			"To":       requestTarget,
			"TextBody": fmt.Sprintf("Vibrant America: %s is your security code. This code is valid for 5 minutes.", verificationCode),
			"Delay":    0,
			"Type":     "SMS",
		}

		// Publish to Kafka
		err = s.publishSMSMessage(ctx, smsID, textMessage)
		if err != nil {
			common.Errorf(fmt.Sprintf("[%s] Error sending SMS: %v", trackingID, err), err)
			return nil, fmt.Errorf("error sending SMS: %w", err)
		}

		common.Infof(fmt.Sprintf("[%s] Verification code sent via SMS to %s", trackingID, requestTarget))
		return &pb.ForgetPasswordRequestResponse{
			Code:    200,
			Message: "Code sent via text",
		}, nil

	default:
		common.Infof(fmt.Sprintf("[%s] Invalid request method: %s", trackingID, requestMethod))
		return &pb.ForgetPasswordRequestResponse{
			Code:    400,
			Message: "Invalid request method",
		}, nil
	}
}

// ForgetPassword completes the password reset process by verifying the code and updating the password  
func (s *UserService) ForgetPassword(ctx context.Context, username, verificationCode, newPassword string) (*pb.ForgetPasswordVerifyResponse, error) {
	trackingID := util.GenerateUUID()
	common.Infof(fmt.Sprintf("[%s] ForgetPassword: Starting for username %s", trackingID, username))

	// INPUT VALIDATION (OPTIMIZATION: Added proper input validation)
	if strings.TrimSpace(username) == "" {
		common.Errorf(fmt.Sprintf("[%s] Empty username provided", trackingID), fmt.Errorf("empty username"))
		return &pb.ForgetPasswordVerifyResponse{
			Code:    400,
			Message: "Username is required",
		}, nil
	}

	// Sanitize username input (OPTIMIZATION: Added input sanitization)
	username = strings.TrimSpace(username)

	var userInfo *ent.User
	var err error

	// Check if username is an email address
	emailRegex := `^[^\s@]+@[^\s@]+\.[^\s@]+$`
	if matched, _ := regexp.MatchString(emailRegex, username); matched {
		userInfo, err = s.dbClient.User.Query().
			Where(user.EmailUserID(username)).
			WithCustomer().
			WithInternalUser().
			WithPatient().
			WithClinic().
			First(ctx)
		if err != nil && !ent.IsNotFound(err) {
			common.Errorf(fmt.Sprintf("[%s] Error querying user by email: %v", trackingID, err), err)
			return &pb.ForgetPasswordVerifyResponse{
				Code:    500,
				Message: "Internal server error",
			}, nil
		}
	}

	// If not found by email, try by username
	if userInfo == nil {
		userInfo, err = s.dbClient.User.Query().
			Where(user.UserName(username)).
			WithCustomer().
			WithInternalUser().
			WithPatient().
			WithClinic().
			First(ctx)
		if err != nil && !ent.IsNotFound(err) {
			common.Errorf(fmt.Sprintf("[%s] Error querying user by username: %v", trackingID, err), err)
			return &pb.ForgetPasswordVerifyResponse{
				Code:    500,
				Message: "Internal server error",
			}, nil
		}
	}

	// BUG FIX: Proper user existence check (TypeScript bug was accessing userInfo without null check)
	if userInfo == nil {
		common.Infof(fmt.Sprintf("[%s] User not found: %s", trackingID, username))
		return &pb.ForgetPasswordVerifyResponse{
			Code:    404,
			Message: "User not found",
		}, nil
	}

	// Check if 2FA is enabled
	if userInfo.IsTwoFactorAuthenticationEnabled {
		if verificationCode == "" {
			return &pb.ForgetPasswordVerifyResponse{
				Code:    403,
				Message: "Please Enter 2FA Password",
			}, nil
		}

		// Check verification code from Redis
		redisKey := fmt.Sprintf("lis::core_service::user_service:forget_password_verification_user_%s", username)
		storedCode, err := s.redisClient.Get(ctx, redisKey).Result()
		if err != nil {
			if err == redis.Nil {
				common.Infof(fmt.Sprintf("[%s] No verification code found in Redis for user: %s", trackingID, username))
				return &pb.ForgetPasswordVerifyResponse{
					Code:    401,
					Message: fmt.Sprintf("Never received a forget password request from user: %s or code is expired.", username),
				}, nil
			}
			common.Errorf(fmt.Sprintf("[%s] Error retrieving verification code from Redis: %v", trackingID, err), err)
			return &pb.ForgetPasswordVerifyResponse{
				Code:    500,
				Message: "Internal server error",
			}, nil
		}

		// Verify the code
		if verificationCode != storedCode {
			common.Infof(fmt.Sprintf("[%s] Wrong verification code for user: %s", trackingID, username))
			return &pb.ForgetPasswordVerifyResponse{
				Code:    403,
				Message: "Wrong Verification Number",
			}, nil
		}

		// BUG FIX: Fixed TypeScript logic bug (was: new_password || new_password != "" which always true)
		// OPTIMIZATION: Added proper password validation
		if newPassword != "" {
			// OPTIMIZATION: Validate password strength
			if err := s.validatePassword(newPassword); err != nil {
				common.Errorf(fmt.Sprintf("[%s] Password validation failed: %v", trackingID, err), err)
				return &pb.ForgetPasswordVerifyResponse{
					Code:    400,
					Message: fmt.Sprintf("Password validation failed: %s", err.Error()),
				}, nil
			}

			err = s.updateUserPassword(ctx, userInfo, newPassword)
			if err != nil {
				common.Errorf(fmt.Sprintf("[%s] Error updating password: %v", trackingID, err), err)
				return &pb.ForgetPasswordVerifyResponse{
					Code:    500,
					Message: "Error updating password",
				}, nil
			}

			// Delete verification code from Redis
			err = s.redisClient.Del(ctx, redisKey).Err()
			if err != nil {
				common.Errorf(fmt.Sprintf("[%s] Error deleting verification code from Redis: %v", trackingID, err), err)
			}

			common.Infof(fmt.Sprintf("[%s] Password updated successfully for user: %s", trackingID, username))
			return &pb.ForgetPasswordVerifyResponse{
				Code:    201,
				Message: "Password Updated",
			}, nil
		} else {
			return &pb.ForgetPasswordVerifyResponse{
				Code:    200,
				Message: "Verification Passed, Enter New Password to Set the Password",
			}, nil
		}
	} else {
		// BUG FIX: TypeScript allowed empty passwords for non-2FA users, now we validate
		if strings.TrimSpace(newPassword) == "" {
			common.Errorf(fmt.Sprintf("[%s] Empty password provided for non-2FA user", trackingID), fmt.Errorf("empty password"))
			return &pb.ForgetPasswordVerifyResponse{
				Code:    400,
				Message: "Password is required",
			}, nil
		}

		// OPTIMIZATION: Validate password strength for non-2FA users too
		if err := s.validatePassword(newPassword); err != nil {
			common.Errorf(fmt.Sprintf("[%s] Password validation failed for non-2FA user: %v", trackingID, err), err)
			return &pb.ForgetPasswordVerifyResponse{
				Code:    400,
				Message: fmt.Sprintf("Password validation failed: %s", err.Error()),
			}, nil
		}

		err = s.updateUserPassword(ctx, userInfo, newPassword)
		if err != nil {
			common.Errorf(fmt.Sprintf("[%s] Error updating password: %v", trackingID, err), err)
			return &pb.ForgetPasswordVerifyResponse{
				Code:    500,
				Message: "Error updating password",
			}, nil
		}

		common.Infof(fmt.Sprintf("[%s] Password updated successfully for user without 2FA: %s", trackingID, username))
		return &pb.ForgetPasswordVerifyResponse{
			Code:    201,
			Message: "Password Updated",
		}, nil
	}
}

// validatePassword validates password strength according to exact requirements:
// At least 8-20 characters, at least 1 lowercase letter [a-z], at least 1 uppercase letter [A-Z], 
// at least 1 number [0-9], At least 1 special character [*.!@#$%^&(){}[]:;<>,.?/~_+-=|]
func (s *UserService) validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	if len(password) > 20 {
		return fmt.Errorf("password must be at most 20 characters long")
	}
	
	// Check for at least one uppercase letter [A-Z]
	if matched, _ := regexp.MatchString(`[A-Z]`, password); !matched {
		return fmt.Errorf("password must contain at least 1 uppercase letter [A-Z]")
	}
	
	// Check for at least one lowercase letter [a-z]
	if matched, _ := regexp.MatchString(`[a-z]`, password); !matched {
		return fmt.Errorf("password must contain at least 1 lowercase letter [a-z]")
	}
	
	// Check for at least one number [0-9]
	if matched, _ := regexp.MatchString(`[0-9]`, password); !matched {
		return fmt.Errorf("password must contain at least 1 number [0-9]")
	}
	
	// Check for at least one special character [*.!@#$%^&(){}[]:;<>,.?/~_+-=|]
	if matched, _ := regexp.MatchString(`[\*\.!@#\$%\^&\(\)\{\}\[\]:;<>,\.\?/~_+\-=\|]`, password); !matched {
		return fmt.Errorf("password must contain at least 1 special character")
	}
	
	return nil
}

// updateUserPassword is a helper method to update user password with bcrypt hashing
func (s *UserService) updateUserPassword(ctx context.Context, userInfo *ent.User, newPassword string) error {
	// Hash the new password with bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// Update the user's password in the database
	_, err = s.dbClient.User.UpdateOneID(userInfo.ID).
		SetPassword(string(hashedPassword)).
		SetImportedUserWithSaltPassword(false). // Mark salt password as false like TypeScript version
		Save(ctx)
	if err != nil {
		return fmt.Errorf("error updating user password in database: %w", err)
	}

	return nil
}

// TurnOff2FASettingPage disables two-factor authentication for a user
func (s *UserService) TurnOff2FASettingPage(ctx context.Context, username, token string) (*pb.TurnOff2FASettingPageResponse, error) {
	// Generate a tracking ID for logging
	trackingID := util.GenerateUUID()
	common.Infof(fmt.Sprintf("[%s] TurnOff2FASettingPage: Starting for username %s", trackingID, username))

	// Input validation
	username = strings.TrimSpace(username)
	token = strings.TrimSpace(token)

	if username == "" {
		common.Errorf(fmt.Sprintf("[%s] Empty username provided", trackingID), fmt.Errorf("empty username"))
		return &pb.TurnOff2FASettingPageResponse{
			Code:    400,
			Message: "Username is required",
		}, nil
	}

	if token == "" {
		common.Errorf(fmt.Sprintf("[%s] Empty token provided", trackingID), fmt.Errorf("empty token"))
		return &pb.TurnOff2FASettingPageResponse{
			Code:    400,
			Message: "Token is required",
		}, nil
	}

	// Parse the JWT token to validate and get user information
	claims, err := util.ParseJWTToken(token, common.Secrets.JWTSecret)
	if err != nil {
		common.Errorf(fmt.Sprintf("[%s] Failed to parse JWT token: %v", trackingID, err), err)
		return &pb.TurnOff2FASettingPageResponse{
			Code:    401,
			Message: "Invalid token",
		}, nil
	}

	tokenUserID := claims.UserId
	if tokenUserID <= 0 {
		common.Errorf(fmt.Sprintf("[%s] Invalid user ID in token: %d", trackingID, tokenUserID), nil)
		return &pb.TurnOff2FASettingPageResponse{
			Code:    401,
			Message: "Invalid token",
		}, nil
	}

	// Find the user by email or username (matching TypeScript logic)
	var userObj *ent.User

	// Check if username is an email format
	if util.IsValidEmail(username) {
		userObj, err = s.dbClient.User.Query().
			Where(user.EmailUserID(username)).
			Only(ctx)

		if err != nil && !ent.IsNotFound(err) {
			common.Errorf(fmt.Sprintf("[%s] Error querying user by email: %v", trackingID, err), err)
			return nil, err
		}
	}

	// If not found by email, try username
	if userObj == nil {
		userObj, err = s.dbClient.User.Query().
			Where(user.UserName(username)).
			Only(ctx)

		if err != nil {
			if ent.IsNotFound(err) {
				common.Infof(fmt.Sprintf("[%s] User not found: %s", trackingID, username))
				return &pb.TurnOff2FASettingPageResponse{
					Code:    400,
					Message: "2FA Already Enabled",
				}, nil
			}
			common.Errorf(fmt.Sprintf("[%s] Error querying user: %v", trackingID, err), err)
			return nil, err
		}
	}

	// Validate that the token's user ID matches the found user (security check)
	if userObj.ID != tokenUserID {
		common.Errorf(fmt.Sprintf("[%s] Token user ID %d does not match target user ID %d", trackingID, tokenUserID, userObj.ID), nil)
		return &pb.TurnOff2FASettingPageResponse{
			Code:    403,
			Message: "Unauthorized to modify this user's 2FA settings",
		}, nil
	}

	// Check if 2FA is already disabled
	if !userObj.IsTwoFactorAuthenticationEnabled {
		common.Infof(fmt.Sprintf("[%s] 2FA already disabled for user %d", trackingID, userObj.ID))
		return &pb.TurnOff2FASettingPageResponse{
			Code:    400,
			Message: "The 2FA is already off",
		}, nil
	}

	// Update the user to disable 2FA and clear the secret
	_, err = s.dbClient.User.UpdateOneID(userObj.ID).
		SetIsTwoFactorAuthenticationEnabled(false).
		ClearTwoFactorAuthenticationSecret().
		Save(ctx)

	if err != nil {
		common.Errorf(fmt.Sprintf("[%s] Error updating user 2FA settings: %v", trackingID, err), err)
		return nil, fmt.Errorf("error disabling 2FA: %w", err)
	}

	// Log the successful 2FA disable
	common.Infof(fmt.Sprintf("[%s] 2FA disabled successfully for user %d", trackingID, userObj.ID))

	return &pb.TurnOff2FASettingPageResponse{
		Code:    200,
		Message: "2FA Already Disabled",
	}, nil
}

// TurnOn2FASettingPage enables two-factor authentication for a user
func (s *UserService) TurnOn2FASettingPage(ctx context.Context, username, email2faAddress, phone2faNumber, token string) (*pb.TurnOn2FASettingPageResponse, error) {
	// Generate a tracking ID for logging
	trackingID := util.GenerateUUID()
	common.Infof(fmt.Sprintf("[%s] TurnOn2FASettingPage: Starting for username %s", trackingID, username))

	// Input validation
	username = strings.TrimSpace(username)
	email2faAddress = strings.TrimSpace(email2faAddress)
	phone2faNumber = strings.TrimSpace(phone2faNumber)
	token = strings.TrimSpace(token)

	if username == "" {
		common.Errorf(fmt.Sprintf("[%s] Empty username provided", trackingID), fmt.Errorf("empty username"))
		return &pb.TurnOn2FASettingPageResponse{
			Code:    400,
			Message: "Username is required",
		}, nil
	}

	if token == "" {
		common.Errorf(fmt.Sprintf("[%s] Empty token provided", trackingID), fmt.Errorf("empty token"))
		return &pb.TurnOn2FASettingPageResponse{
			Code:    400,
			Message: "Token is required",
		}, nil
	}

	if email2faAddress == "" {
		common.Errorf(fmt.Sprintf("[%s] Empty 2FA email address provided", trackingID), fmt.Errorf("empty 2FA email"))
		return &pb.TurnOn2FASettingPageResponse{
			Code:    400,
			Message: "2FA email address is required",
		}, nil
	}

	if phone2faNumber == "" {
		common.Errorf(fmt.Sprintf("[%s] Empty 2FA phone number provided", trackingID), fmt.Errorf("empty 2FA phone"))
		return &pb.TurnOn2FASettingPageResponse{
			Code:    400,
			Message: "2FA phone number is required",
		}, nil
	}

	// Parse the JWT token to validate and get user information
	claims, err := util.ParseJWTToken(token, common.Secrets.JWTSecret)
	if err != nil {
		common.Errorf(fmt.Sprintf("[%s] Failed to parse JWT token: %v", trackingID, err), err)
		return &pb.TurnOn2FASettingPageResponse{
			Code:    401,
			Message: "Invalid token",
		}, nil
	}

	tokenUserID := claims.UserId
	if tokenUserID <= 0 {
		common.Errorf(fmt.Sprintf("[%s] Invalid user ID in token: %d", trackingID, tokenUserID), nil)
		return &pb.TurnOn2FASettingPageResponse{
			Code:    401,
			Message: "Invalid token",
		}, nil
	}

	// Find the user by email or username (matching TypeScript logic)
	var userObj *ent.User

	// Check if username is an email format
	if util.IsValidEmail(username) {
		userObj, err = s.dbClient.User.Query().
			Where(user.EmailUserID(username)).
			WithCustomer().WithInternalUser().WithPatient().WithClinic().
			Only(ctx)

		if err != nil && !ent.IsNotFound(err) {
			common.Errorf(fmt.Sprintf("[%s] Error querying user by email: %v", trackingID, err), err)
			return nil, err
		}
	}

	// If not found by email, try username
	if userObj == nil {
		userObj, err = s.dbClient.User.Query().
			Where(user.UserName(username)).
			WithCustomer().WithInternalUser().WithPatient().WithClinic().
			Only(ctx)

		if err != nil {
			if ent.IsNotFound(err) {
				common.Infof(fmt.Sprintf("[%s] User not found: %s", trackingID, username))
				return &pb.TurnOn2FASettingPageResponse{
					Code:    400,
					Message: "2FA Already Enabled",
				}, nil
			}
			common.Errorf(fmt.Sprintf("[%s] Error querying user: %v", trackingID, err), err)
			return nil, err
		}
	}

	// Validate that the token's user ID matches the found user (security check)
	if userObj.ID != tokenUserID {
		common.Errorf(fmt.Sprintf("[%s] Token user ID %d does not match target user ID %d", trackingID, tokenUserID, userObj.ID), nil)
		return &pb.TurnOn2FASettingPageResponse{
			Code:    403,
			Message: "Unauthorized to modify this user's 2FA settings",
		}, nil
	}

	// Check if 2FA is already enabled
	if userObj.IsTwoFactorAuthenticationEnabled {
		common.Infof(fmt.Sprintf("[%s] 2FA already enabled for user %d", trackingID, userObj.ID))
		return &pb.TurnOn2FASettingPageResponse{
			Code:    400,
			Message: "2FA Already Enabled",
		}, nil
	}

	// Validate that 2FA email address is not the same as login email
	if userObj.EmailUserID != "" && userObj.EmailUserID == email2faAddress {
		common.Errorf(fmt.Sprintf("[%s] 2FA email same as login email: %s", trackingID, email2faAddress), nil)
		return &pb.TurnOn2FASettingPageResponse{
			Code:    400,
			Message: "2FA Email Address Cannot be the same as Email Login ID",
		}, nil
	}

	// Update or create 2FA email contact
	err = s.updateOrCreate2FAContact(ctx, trackingID, userObj.ID, email2faAddress, "email", "email_2fa")
	if err != nil {
		common.Errorf(fmt.Sprintf("[%s] Error updating 2FA email contact: %v", trackingID, err), err)
		return nil, fmt.Errorf("error updating 2FA email contact: %w", err)
	}

	// Update or create 2FA phone contact
	err = s.updateOrCreate2FAContact(ctx, trackingID, userObj.ID, phone2faNumber, "phone", "phone_2fa")
	if err != nil {
		common.Errorf(fmt.Sprintf("[%s] Error updating 2FA phone contact: %v", trackingID, err), err)
		return nil, fmt.Errorf("error updating 2FA phone contact: %w", err)
	}

	// Generate TOTP secret and OTP Auth URL
	otpAuthURL, secret, err := s.generateTwoFactorAuthenticationSecret(userObj.UserName)
	if err != nil {
		common.Errorf(fmt.Sprintf("[%s] Error generating 2FA secret: %v", trackingID, err), err)
		return nil, fmt.Errorf("error generating 2FA secret: %w", err)
	}

	// Enable 2FA and save the secret to the database
	_, err = s.dbClient.User.UpdateOneID(userObj.ID).
		SetIsTwoFactorAuthenticationEnabled(true).
		SetTwoFactorAuthenticationSecret(secret).
		Save(ctx)

	if err != nil {
		common.Errorf(fmt.Sprintf("[%s] Error enabling 2FA for user: %v", trackingID, err), err)
		return nil, fmt.Errorf("error enabling 2FA: %w", err)
	}

	// Log the successful 2FA enable
	common.Infof(fmt.Sprintf("[%s] 2FA enabled successfully for user %d", trackingID, userObj.ID))

	return &pb.TurnOn2FASettingPageResponse{
		Code:       200,
		Message:    "2FA Enabled",
		OtpauthUrl: otpAuthURL,
	}, nil
}

// Helper function to update or create 2FA contact information
func (s *UserService) updateOrCreate2FAContact(ctx context.Context, trackingID string, userID int, contactDetails, contactType, contactDescription string) error {
	// Query for existing contact
	existingContact, err := s.dbClient.Contact.Query().
		Where(
			contact.UserID(userID),
			contact.ContactDescription(contactDescription),
		).
		First(ctx)

	if err != nil && !ent.IsNotFound(err) {
		return fmt.Errorf("error querying existing contact: %w", err)
	}

	if existingContact != nil {
		// Update existing contact
		common.Infof(fmt.Sprintf("[%s] Updating existing %s contact for user %d", trackingID, contactType, userID))
		_, err = s.dbClient.Contact.UpdateOneID(existingContact.ID).
			SetContactDetails(contactDetails).
			SetIs2faContact(true).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("error updating existing contact: %w", err)
		}
	} else {
		// Create new contact
		common.Infof(fmt.Sprintf("[%s] Creating new %s contact for user %d", trackingID, contactType, userID))
		_, err = s.dbClient.Contact.Create().
			SetContactDetails(contactDetails).
			SetIs2faContact(true).
			SetUserID(userID).
			SetContactType(contactType).
			SetContactDescription(contactDescription).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("error creating new contact: %w", err)
		}
	}

	return nil
}

// Helper function to generate TOTP secret and OTP Auth URL (matching TypeScript implementation)
func (s *UserService) generateTwoFactorAuthenticationSecret(username string) (string, string, error) {
	// Generate a new TOTP secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "LIS", // Matching TypeScript TWO_FACTOR_AUTHENTICATION_APP_NAME config
		AccountName: username,
		SecretSize:  32, // Standard 32-byte secret
	})
	if err != nil {
		return "", "", fmt.Errorf("error generating TOTP key: %w", err)
	}

	// Get the secret string
	secret := key.Secret()

	// Get the OTP Auth URL for QR code generation
	otpAuthURL := key.URL()

	return otpAuthURL, secret, nil
}

func (s *UserService) UpdateUserInvitationRecord(ctx context.Context, customerID int32, invitationLink string) (*pb.UpdateUserInvitationRecordResponse, error) {
	invitationRecord, err := s.dbClient.UserInvitationRecord.Query().
		Where(userinvitationrecord.CustomerID(int(customerID))).
		First(ctx)
	
	if err != nil {
		if ent.IsNotFound(err) {
			return &pb.UpdateUserInvitationRecordResponse{
				Code:         400,
				Message:      "No Invitation Record Found",
				ErrorMessage: "No Invitation Record Found",
			}, nil
		}
		common.Errorf("Error finding invitation record for customer", err)
		return &pb.UpdateUserInvitationRecordResponse{
			Code:         500,
			Message:      "Internal Error",
			ErrorMessage: fmt.Sprintf("Error : %s", err.Error()),
		}, err
	}

	previousInvitationLink := invitationRecord.InvitationLink

	_, err = s.dbClient.UserInvitationRecord.
		UpdateOneID(invitationRecord.ID).
		SetInvitationLink(invitationLink).
		Save(ctx)
	
	if err != nil {
		common.Errorf("Error updating invitation record for customer", err)
		return &pb.UpdateUserInvitationRecordResponse{
			Code:         500,
			Message:      "Internal Error",
			ErrorMessage: fmt.Sprintf("Error : %s", err.Error()),
		}, err
	}

	auditLogMessage := common.AuditLogEntry{
		EventID:                uuid.NewString(),
		ServiceName:            common.ServiceName,
		ServiceType:            "backend",
		EventName:              "UpdateUserInvitationRecord",
		EntityType:             "customer",
		EntityID:               strconv.Itoa(int(customerID)),
		AttributeValuePrior:    previousInvitationLink,
		AttributeValuePost:     invitationLink,
		Entrypoint:             "GRPC",
	}
	go func() {
		common.RecordAuditLog(auditLogMessage)
	}()

	return &pb.UpdateUserInvitationRecordResponse{
		Code:         200,
		Message:      "Invitation Link Updated",
		ErrorMessage: "",
	}, nil
}
