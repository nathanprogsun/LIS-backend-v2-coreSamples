package handler

import (
	"context"
	"coresamples/common"
	pb "coresamples/proto"
	"coresamples/service"
	"coresamples/util"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type UserHandler struct {
	UserService service.IUserService
}

func (h *UserHandler) LogIn(ctx context.Context, req *pb.LogInRequest, resp *pb.LogInResponse) error {
	return nil
}

func (h *UserHandler) GetUserInfoByRole(ctx context.Context, req *pb.GetUserInfoByRoleRequest, resp *pb.GetUserInfoByRoleResponse) error {
	// Extract parameters from the request
	userID := req.UserId
	userRole := req.UserRole

	// Call the service method
	result, err := h.UserService.GetUserInfoByRole(ctx, userID, userRole)
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("nil result returned from service")
	}

	// Copy the response
	*resp = *result
	return nil
}

func (h *UserHandler) UpdatePasswordByOldPassword(ctx context.Context, req *pb.UpdatePasswordByOldPasswordRequest, resp *pb.UpdateUserPasswordResponse) error {
	return nil
}

func (h *UserHandler) TransferLISTokenToPortal(ctx context.Context, req *pb.TransferLISTokenToPortalRequest, resp *pb.TransferLISTokenToPortalResponse) error {
	return nil
}

func (h *UserHandler) Send2FATokenRequest(ctx context.Context, req *pb.Send2FATokenRequestMessage, resp *pb.Send2FATokenResponse) error {
	return nil
}

// ForgetPasswordRequest handles password reset request by sending a verification code
// @Summary Handle password reset request
// @Description Handles password reset request by sending a verification code via email or SMS to users with 2FA enabled
// @Tags users,authentication,password-reset
// @Accept grpc
// @Produce grpc
// @Param username query string true "Username or email for password reset"
// @Param request_method query string true "Method to send code (email or phone)"
// @Param request_target query string true "Target email or phone number"
// @Success 200 {object} pb.ForgetPasswordRequestResponse
// @Router /user/forgetPasswordRequest [grpc]
func (h *UserHandler) ForgetPasswordRequest(ctx context.Context, req *pb.ForgetPasswordRequestRequest, resp *pb.ForgetPasswordRequestResponse) error {
	username := req.Username
	requestMethod := req.RequestMethod
	requestTarget := req.RequestTarget
	
	result, err := h.UserService.ForgetPasswordRequest(ctx, username, requestMethod, requestTarget)
	if err != nil {
		return err
	}
	
	err = util.Swap(result, resp)
	if err != nil {
		return err
	}
	
	return nil
}

// ForgetPassword completes password reset with verification code
// @Summary Complete password reset with verification code
// @Description Verifies the code sent via ForgetPasswordRequest and updates the user's password
// @Tags users,authentication,password-reset
// @Accept grpc
// @Produce grpc
// @Param username query string true "Username or email for password reset"
// @Param verification_code query string true "Verification code received via email/SMS"
// @Param new_password query string true "New password to set"
// @Success 200 {object} pb.ForgetPasswordVerifyResponse
// @Router /user/forgetPassword [grpc]
func (h *UserHandler) ForgetPassword(ctx context.Context, req *pb.ForgetPasswordRequestMessage, resp *pb.ForgetPasswordVerifyResponse) error {
	username := req.Username
	verificationCode := req.VerificationCode
	newPassword := req.NewPassword
	
	result, err := h.UserService.ForgetPassword(ctx, username, verificationCode, newPassword)
	if err != nil {
		return err
	}
	
	err = util.Swap(result, resp)
	if err != nil {
		return err
	}
	
	return nil
}

func (h *UserHandler) TransferCustomerClinic(ctx context.Context, req *pb.TransferCustomerClinicRequest, resp *pb.TransferCustomerClinicResponse) error {
	return nil
}

// Set Up Sign In Email

func (h *UserHandler) SetUpEmailRequest(ctx context.Context, req *pb.SetUpLoginEmailRequestRequest, resp *pb.SetUpLoginEmailRequestResponse) error {
	return nil
}

func (h *UserHandler) VerifySetUpUserEmailLogIn(ctx context.Context, req *pb.VerifySetUpUserEmailLogInRequest, resp *pb.VerifySetUpUserEmailLogInResponse) error {
	return nil
}

// IsEmailUsedAsLoginId checks if an email is already used as a login ID in the database
// @Summary Check if email is used as login ID
// @Description Checks if the provided email is already used as a login ID by any user
// @Tags users,authentication
// @Accept grpc
// @Produce grpc
// @Param email query string true "Email to check"
// @Success 200 {object} pb.IsEmailUsedAsLoginIdResponse
// @Router /user/isEmailUsedAsLoginId [grpc]
func (h *UserHandler) IsEmailUsedAsLoginId(ctx context.Context, req *pb.EmailRequest, resp *pb.IsEmailUsedAsLoginIdResponse) error {
	// Extract email from the request
	email := req.Email

	// Call the service method to check if the email is used as login ID
	result, err := h.UserService.IsEmailUsedAsLoginId(ctx, email)

	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("nil result returned from service")
	}

	// Copy the response from the service result
	resp.UsedAsEmailLogId = result.UsedAsEmailLogId
	resp.Message = result.Message
	resp.UserId = result.UserId

	return nil
}

// InitialForgetPassword initiates the forget password process by sending a reset email
// @Summary Initiate forget password process
// @Description Initiates the forget password process by validating the email and sending a password reset email to the user
// @Tags users,authentication,password-reset
// @Accept grpc
// @Produce grpc
// @Param email_address query string true "Email address to send password reset link"
// @Success 200 {object} pb.ForgetPasswordResponse
// @Router /user/initialForgetPassword [grpc]
func (h *UserHandler) InitialForgetPassword(ctx context.Context, req *pb.InitialForgetPasswordRequest, resp *pb.ForgetPasswordResponse) error {
	emailAddress := req.EmailAddress
	
	result, err := h.UserService.InitialForgetPassword(ctx, emailAddress)
	if err != nil {
		return err
	}
	
	err = util.Swap(result, resp)
	if err != nil {
		return nil
	}
	
	return nil
}

// Setting Page Enable 2FA
func (h *UserHandler) SendVerify2FASetUpContactInfo(ctx context.Context, req *pb.SendVerify2FASetUpContactInfoRequest, resp *pb.SendVerify2FASetUpContactInfoResponse) error {
	return nil
}

func (h *UserHandler) Verify2FASetUpContactInfo(ctx context.Context, req *pb.Verify2FASetUpContactRequest, resp *pb.Verify2FASetUpContactResponse) error {
	return nil
}

// TurnOn2FASettingPage enables two-factor authentication for a user
// @Summary Turn on 2FA for user
// @Description Enables two-factor authentication for a user with email and phone contact information
// @Tags users,authentication,2fa
// @Accept grpc
// @Produce grpc
// @Param username query string true "Username or email used for login"
// @Param email_2fa_address query string true "Email address for 2FA verification"
// @Param phone_2fa_number query string true "Phone number for 2FA verification"
// @Param token query string true "JWT token for authentication"
// @Success 200 {object} pb.TurnOn2FASettingPageResponse
// @Router /user/turnOn2FASettingPage [grpc]
func (h *UserHandler) TurnOn2FASettingPage(ctx context.Context, req *pb.TurnOn2FASettingPageRequest, resp *pb.TurnOn2FASettingPageResponse) error {
	result, err := h.UserService.TurnOn2FASettingPage(ctx, req.Username, req.Email_2FaAddress, req.Phone_2FaNumber, req.Token)
	if err != nil {
		return err
	}
	err = util.Swap(result, resp)
	if err != nil {
		return err
	}
	return nil
}

// Send2FAVerificationCode sends a verification code for 2FA via email or SMS
// @Summary Send 2FA verification code
// @Description Sends a verification code for two-factor authentication via email or SMS
// @Tags users,authentication,2fa
// @Accept grpc
// @Produce grpc
// @Param username query string true "Username or email used for login"
// @Param email_address query string false "Email address to send verification code"
// @Param phone_number query string false "Phone number to send verification code"
// @Success 200 {object} pb.Send2FAVerificationCodeResponse
// @Router /user/send2FAVerificationCode [grpc]
func (h *UserHandler) Send2FAVerificationCode(ctx context.Context, req *pb.Send2FAVerificationCodeRequest, resp *pb.Send2FAVerificationCodeResponse) error {
	// Extract parameters from the request
	username := req.Username
	emailAddress := req.EmailAddress
	phoneNumber := req.PhoneNumber
	
	// Call the service method
	result, err := h.UserService.Send2FAVerificationCode(ctx, username, emailAddress, phoneNumber)
	if err != nil {
		return err
	}
	
	// Use util.Swap to transfer data from result to response
	err = util.Swap(result, resp)
	if err != nil {
		return err
	}
	
	return nil
}

// Verify2FAVerificationCode verifies a 2FA verification code
// @Summary Verify 2FA verification code
// @Description Verifies a 2FA verification code sent via email or SMS
// @Tags users,authentication,2fa
// @Accept grpc
// @Produce grpc
// @Param username query string true "Username or email used for login"
// @Param verification_code query string true "2FA verification code to verify"
// @Param email_address query string false "Email address used for verification"
// @Param phone_number query string false "Phone number used for verification"
// @Success 200 {object} pb.Verify2FAVerificationResponse
// @Router /user/verify2FAVerificationCode [grpc]
func (h *UserHandler) Verify2FAVerificationCode(ctx context.Context, req *pb.Verify2FAVerificationCodeRequest, resp *pb.Verify2FAVerificationResponse) error {
	// Extract parameters from the request
	username := req.Username
	verificationCode := req.VerificationCode
	emailAddress := req.EmailAddress
	phoneNumber := req.PhoneNumber
	
	// Call the service method
	result, err := h.UserService.Verify2FAVerificationCode(ctx, username, verificationCode, emailAddress, phoneNumber)
	if err != nil {
		return err
	}
	
	// Use util.Swap to transfer data from result to response
	err = util.Swap(result, resp)
	if err != nil {
		return err
	}
	
	return nil
}

// TurnOff2FASettingPage disables two-factor authentication for a user
// @Summary Turn off 2FA for user
// @Description Disables two-factor authentication for a user by setting isTwoFactorAuthenticationEnabled to false and clearing the authentication secret
// @Tags users,authentication,2fa
// @Accept grpc
// @Produce grpc
// @Param username query string true "Username or email used for login"
// @Param token query string true "JWT token for authentication"
// @Success 200 {object} pb.TurnOff2FASettingPageResponse
// @Router /user/turnOff2FASettingPage [grpc]
func (h *UserHandler) TurnOff2FASettingPage(ctx context.Context, req *pb.TurnOff2FASettingPageRequest, resp *pb.TurnOff2FASettingPageResponse) error {
	// Extract parameters from the request
	username := req.Username
	token := req.Token
	
	// Call the service method
	result, err := h.UserService.TurnOff2FASettingPage(ctx, username, token)
	if err != nil {
		return err
	}
	
	// Use util.Swap to transfer data from result to response
	err = util.Swap(result, resp)
	if err != nil {
		return err
	}
	
	return nil
}

// GetUserInformation retrieves detailed information about a user based on user_id
// @Summary Get user information
// @Description Retrieves detailed user information including customer details based on user ID
// @Tags users
// @Accept grpc
// @Produce grpc
// @Param request body pb.GetUserInformationRequest true "User ID to look up"
// @Success 200 {object} pb.GetUserInformationResponse "User information with customer details if available"
// @Router /user/getUserInformation [grpc]
func (h *UserHandler) GetUserInformation(ctx context.Context, req *pb.GetUserInformationRequest, resp *pb.GetUserInformationResponse) error {
	// Get user information from service
	result, err := h.UserService.GetUserInformation(ctx, req.UserId)
	if err != nil {
		return err
	}

	// Use util.Swap to transfer data from result to response
	err = util.Swap(result, resp)
	if err != nil {
		return err
	}
	return nil
}

// GetUser2FAContactInfo retrieves 2FA contacts for the user based on JWT token
// @Summary Get 2FA contact information for a user
// @Description Retrieves all contacts marked as 2FA contacts for the user identified by the JWT token
// @Tags users,2fa,authentication
// @Accept grpc
// @Produce grpc
// @Param token query string true "JWT token for user authentication"
// @Success 200 {object} pb.GetUser2FAContactInfoResponse
// @Router /user/getUser2FAContactInfo [grpc]
func (h *UserHandler) GetUser2FAContactInfo(ctx context.Context, req *pb.GetUser2FAContactInfoRequest, resp *pb.GetUser2FAContactInfoResponse) error {
	token := req.Token

	result, err := h.UserService.GetUser2FAContactInfo(ctx, token)
	if err != nil {
		return err
	}

	err = util.Swap(result, resp)
	if err != nil {
		return err
	}

	return nil
}

// Change Email
func (h *UserHandler) ChangeUserEmailLogInID(ctx context.Context, req *pb.UserChangeEmailRequest, resp *pb.UserChangeEmailResponse) error {
	return nil
}

func (h *UserHandler) AdminUserLoginSearch(ctx context.Context, req *pb.AdminLoginSearchRequest, resp *pb.AdminLoginSearchResponse) error {
	return nil
}

func (h *UserHandler) AdminLogin(ctx context.Context, req *pb.AdminLoginRequest, resp *pb.AdminLoginResponse) error {
	return nil
}

// RenewToken validates and renews a JWT token
// @Summary Renew JWT token
// @Description Validates an existing JWT token and issues a new token with extended expiration. The default JWT_EXPIRATION_TIME is 2700 seconds.
// @Tags authentication,tokens
// @Accept grpc
// @Produce grpc
// @Param request body pb.RenewTokenRequest true "JWT token to renew"
// @Success 200 {object} pb.RenewTokenResponse "Renewed token with expiration time"
// @Failure 400 {object} pb.RenewTokenResponse "Error message with no token"
// @Router /user/renewToken [grpc]
func (h *UserHandler) RenewToken(ctx context.Context, req *pb.RenewTokenRequest, resp *pb.RenewTokenResponse) error {
	// Extract JWT token from request
	jwtToken := req.JwtToken

	// Call the service method
	result, err := h.UserService.RenewToken(ctx, jwtToken)
	if err != nil {
		return err
	}

	// Use util.Swap to transfer data from result to response
	err = util.Swap(result, resp)
	if err != nil {
		return err
	}

	return nil
}

func (h *UserHandler) LISLogin(ctx context.Context, req *pb.LogInRequest, resp *pb.LogInResponse) error {
	return nil
}

// Version 0.7.5

func (h *UserHandler) CreateUserLogInForInvitedCustomer(ctx context.Context, req *pb.CreateUserLogForInvitedCustomerRequest, resp *pb.CreateUserLogForInvitedCustomerResponse) error {
	return nil
}

// CheckWhetherEmailIsUsedAsLoginId checks if an email is already used as a login ID by a user in the specified clinic
// @Summary Check if email is used as login ID in a specific clinic
// @Description Checks if the provided email is already used as a login ID by any user in the specified clinic
// @Tags users,authentication
// @Accept grpc
// @Produce grpc
// @Param email query string true "Email to check"
// @Param clinic_id query string true "Clinic ID to check against"
// @Success 200 {object} pb.CheckWhetherEmailIsUsedAsLoginIdResponse
// @Router /user/checkWhetherEmailIsUsedAsLoginId [grpc]
func (h *UserHandler) CheckWhetherEmailIsUsedAsLoginId(ctx context.Context, req *pb.CheckWhetherEmailIsUsedAsLoginIdRequest, resp *pb.CheckWhetherEmailIsUsedAsLoginIdResponse) error {
	// Extract parameters from the request
	email := req.Email
	clinicID := req.ClinicId

	// Call the service method
	result, err := h.UserService.CheckWhetherEmailIsUsedAsLoginId(ctx, email, clinicID)
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("nil result returned from service")
	}

	// Copy the response from the service result
	resp.ExistingUser = result.ExistingUser
	resp.Message = result.Message

	return nil
}

func (h *UserHandler) CannySSO(ctx context.Context, req *pb.CannySSORequest, resp *pb.CannySSOResponse) error {
	return nil
}

func (h *UserHandler) ForceChangeLoginEmailInternal(ctx context.Context, req *pb.ForceChangeLoginEmailInternalRequest, resp *pb.ForceChangeLoginEmailInternalResponse) error {
	return nil
}

// UpdateUserInvitationRecord updates the invitation link for a user invitation record
// @Summary Update user invitation record
// @Description Updates the invitation link for a customer's user invitation record
// @Tags users,invitation
// @Accept grpc
// @Produce grpc
// @Param customer_id query int32 true "Customer ID"
// @Param invitation_link query string true "New invitation link"
// @Success 200 {object} pb.UpdateUserInvitationRecordResponse
// @Router /user/updateUserInvitationRecord [grpc]
func (h *UserHandler) UpdateUserInvitationRecord(ctx context.Context, req *pb.UpdateUserInvitationRecordRequest, resp *pb.UpdateUserInvitationRecordResponse) error {
	trackingID := uuid.NewString()
	xRequestID := trackingID
	serviceName := "Unknown Caller"

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if val := md.Get("x-request-id"); len(val) > 0 {
			xRequestID = val[0]
		}
		if val := md.Get("service-name"); len(val) > 0 {
			serviceName = val[0]
		}
	}

	common.Infof("[%s] UpdateUserInvitationRecord called by: %s, input: %+v", xRequestID, serviceName, req)

	result, err := h.UserService.UpdateUserInvitationRecord(ctx, req.CustomerId, req.InvitationLink)
	if err != nil {
		return err
	}
	
	err = util.Swap(result, resp)
	if err != nil {
		return err
	}
	return nil
}

func (h *UserHandler) GetLoginHistory(ctx context.Context, req *pb.GetLoginHistoryRequest, resp *pb.GetLoginHistoryResponse) error {
	// Extract parameters from the request
	customerID := req.CustomerId
	userID := req.UserId

	// Parse start time and end time if provided
	var startTime, endTime *time.Time

	if req.StartTime != "" {
		parsedStartTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
		if err != nil {
			return err
		}
		startTime = &parsedStartTime
	}

	if req.EndTime != "" {
		parsedEndTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
		if err != nil {
			return err
		}
		endTime = &parsedEndTime
	}

	// Call the service method
	result, err := h.UserService.GetLoginHistory(ctx, customerID, userID, startTime, endTime, req.PerPage, req.Page)
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("nil result returned from service")
	}

	// Copy the response
	*resp = *result
	return nil
}
