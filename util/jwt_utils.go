package util

import (
	"coresamples/common"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

// ExtendedTokenClaim represents all the fields used in JWT tokens across the application
type ExtendedTokenClaim struct {
	UserId            int      `json:"userId,omitempty"`
	UserPermission    string   `json:"user_permission,omitempty"`
	CustomerId        int      `json:"customer_id,omitempty"`
	ClinicId          int      `json:"clinic_id,omitempty"`
	OldClinicId       int      `json:"old_clinic_id,omitempty"`
	PatientId         int      `json:"patient_id,omitempty"`
	InternalUserId    int      `json:"internal_user_id,omitempty"`
	InternalUserName  string   `json:"internal_user_name,omitempty"`
	InternalUserRole  string   `json:"internal_user_role,omitempty"`
	Role              string   `json:"role,omitempty"`
	CustomerList      []int    `json:"customer_list,omitempty"`
	SessionId         string   `json:"session_id,omitempty"`
	EmailLogInId      string   `json:"email_log_in_id,omitempty"`
	BetaProgramEnabled bool    `json:"beta_program_enabled,omitempty"`
	BetaPrograms      []string `json:"beta_programs,omitempty"`
	UserRoles         []string `json:"user_roles,omitempty"`
	jwt.StandardClaims
}

// GenerateJWTToken creates a new JWT token with the specified claims
func GenerateJWTToken(claims ExtendedTokenClaim) (string, error) {
	// Set the standard JWT claims
	currentTime := time.Now()
	expirationTime := currentTime.Add(time.Second * time.Duration(2700)) // Default 2700 seconds expiration

	// Update standard claims
	claims.StandardClaims = jwt.StandardClaims{
		IssuedAt:  currentTime.Unix(),
		ExpiresAt: expirationTime.Unix(),
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign the token with the secret
	tokenString, err := token.SignedString([]byte(common.Secrets.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}
	
	return tokenString, nil
}

// ParseExtendedJWTToken parses an extended JWT token and returns the claims
func ParseExtendedJWTToken(tokenString string) (*ExtendedTokenClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ExtendedTokenClaim{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(common.Secrets.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if claims, ok := token.Claims.(*ExtendedTokenClaim); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token or claims")
}

// GetJWTExpirationTimestamp returns the expiration timestamp for a new token
func GetJWTExpirationTimestamp() int64 {
	currentTime := time.Now()
	expirationTime := currentTime.Add(time.Second * time.Duration(2700)) // Default 2700 seconds expiration
	return expirationTime.Unix()
}