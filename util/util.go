package util

import (
	"bytes"
	"context"
	"coresamples/common"
	"coresamples/ent"
	"coresamples/ent/rbacroles"
	"coresamples/model"
	cryptorand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/exp/constraints"
)

type TokenClaim struct {
	UserId int `json:"userId,omitempty"`
	jwt.StandardClaims
}

// Swap for swapping (or copy) what's inside `from` to `to`,
// without knowing what's exactly inside the struct
func Swap(from, to interface{}) error {
	bytes, err := json.Marshal(from)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, to)
	return err
}

func StringEqualIgnoreCase(s1 string, s2 string) bool {
	return strings.EqualFold(s1, s2)
}

func Min[T constraints.Ordered](i1 T, i2 T) T {
	if i1 <= i2 {
		return i1
	}
	return i2
}

func EntOpen(driverName string, url string) (*ent.Client, error) {
	drv, err := entsql.Open(driverName, url)
	if err != nil {
		return nil, err
	}
	// Get the underlying sql.DB object of the driver.
	db := drv.DB()

	db.SetConnMaxLifetime(time.Second * 25)
	db.SetMaxOpenConns(30)
	return ent.NewClient(ent.Driver(drv)), nil
}

func ElementsUniqueInt32(array []int32) bool {
	elements := map[int32]bool{}
	for _, num := range array {
		if _, exist := elements[num]; exist {
			return false
		}
		elements[num] = true
	}
	return true
}

func GetRoleTypeEnum(name string) (rbacroles.Type, error) {
	types := []rbacroles.Type{
		rbacroles.TypeExternal,
		rbacroles.TypeClinic,
		rbacroles.TypeInternal,
	}
	for _, t := range types {
		if name == string(t) {
			return t, nil
		}
	}
	common.Error(fmt.Errorf("could not find error type: %s", name))
	return "", errors.New("role type not found")
}

func InterStringToString(a interface{}) string {
	return a.(string)
}

func IntArrayToInt32Array(arr []int) []int32 {
	var res []int32
	Swap(arr, &res)
	return res
}

func Int32ArrayToIntArray(arr []int32) []int {
	var res []int
	Swap(arr, &res)
	return res
}

func TubeTypeToSampleType(tubeType string) string {
	var sampleType string
	switch tubeType {
	case "EDTA_MN":
		sampleType = "EDTA"
	case "DNA_FINGERPRICK":
		sampleType = "DNA fingerprick"
	case "SST":
		sampleType = "Serum"
	case "METAL_FREE_URINE":
		sampleType = "Metal Free Urine"
	case "SODIUM_CITRATE_PLASMA":
		sampleType = "Sodium Citrate Plasma"
	case "BLOOD_FINGERPRICK":
		sampleType = "Blood fingerprick"
	case "STOOL":
		sampleType = "Stool"
	case "UNPRESERVED_STOOL":
		sampleType = "Unpreserved Stool"
	case "URINE":
		sampleType = "Urine"
	case "PLASMA":
		sampleType = "Plasma"
	case "BLOOD_MICROTUBE":
		sampleType = "Blood Microtube"
	case "COVID19_STOOL":
		sampleType = "Saliva"
	case "SALIVA":
		sampleType = "Saliva"
	case "FROZEN_SERUM":
		sampleType = "FROZEN SERUM"
	default:
		sampleType = tubeType
	}
	return sampleType
}

func ParseEventTime(eventTime string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339, eventTime)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid EventTime format: %w", err)
	}
	return parsedTime, nil
}

func ParseJWTToken(token string, secret string) (*TokenClaim, error) {
	parsedToken, _ := jwt.ParseWithClaims(token, &TokenClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if parsedToken != nil {
		if claims, ok := parsedToken.Claims.(*TokenClaim); ok && parsedToken.Valid {
			return claims, nil
		}
	}
	return nil, fmt.Errorf("unable to parse token")
}

func IsSuccessResponse(resp *model.NPIApiResponse) bool {
	if resp == nil {
		return false
	}
	return len(resp.Errors) == 0 && resp.ResultCount == 1 && len(resp.Results) == 1
}

// GenerateUUID generates a new UUID v4 string
func GenerateUUID() string {
	return uuid.NewString()
}

type NameParts struct {
	FirstName  string
	LastName   string
	MiddleName string
}

func SplitName(name string) NameParts {
	nameArray := strings.Fields(name) // Split by any whitespace
	var firstname, lastname, middlename string

	if len(nameArray) == 2 {
		firstname = nameArray[0]
		lastname = nameArray[1]
		middlename = ""
	} else if len(nameArray) >= 3 {
		firstname = nameArray[0]
		middlename = strings.ReplaceAll(nameArray[1], ".", "")
		lastname = nameArray[2]
	} else {
		firstname = name
		lastname = ""
		middlename = ""
	}

	return NameParts{
		FirstName:  firstname,
		LastName:   lastname,
		MiddleName: middlename,
	}
}

func NPIOnlineCheck(npiNumber string, ctx context.Context) (*model.NPIApiResponse, error) {
	url := fmt.Sprintf("https://npiregistry.cms.hhs.gov/api/?number=%s&version=2.1", npiNumber)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept-Encoding", "gzip, deflate")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var npiResp model.NPIApiResponse
	if err := json.Unmarshal(body, &npiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal NPI response: %w", err)
	}

	return &npiResp, nil
}

func IsNumericString(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func AssembleFullName(firstName, middleName, lastName string) string {
	var parts []string

	if strings.TrimSpace(firstName) != "" {
		parts = append(parts, strings.TrimSpace(firstName))
	}
	if strings.TrimSpace(middleName) != "" {
		parts = append(parts, strings.TrimSpace(middleName))
	}
	if strings.TrimSpace(lastName) != "" {
		parts = append(parts, strings.TrimSpace(lastName))
	}

	return strings.Join(parts, " ")
}

func EqualIntSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[int]int)
	for _, x := range a {
		m[x]++
	}
	for _, y := range b {
		if m[y] == 0 {
			return false
		}
		m[y]--
	}
	return true
}

func Contains(slice []string, item string) bool {
	itemLower := strings.ToLower(item) // Convert search item to lowercase

	for _, v := range slice {
		if strings.ToLower(v) == itemLower { // Convert slice element to lowercase
			return true
		}
	}
	return false
}
func SliceEqual(a, b []int32) bool {
	if len(a) != len(b) {
		return false
	}

	aMap := make(map[int32]bool)
	for _, val := range a {
		aMap[val] = true
	}

	for _, val := range b {
		if !aMap[val] {
			return false
		}
	}

	return true
}

// util.MustMarshalJSON is a helper that panics if marshal fails
func MustMarshalJSON(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal: %v", err))
	}
	return string(bytes)
}

// IsValidEmail checks if a string is a valid email format
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	return emailRegex.MatchString(email)
}

// GenerateOTP generates a random numeric OTP of the specified length
// Todo: Use third-party lib
func GenerateOTP(length int) string {
	const digits = "0123456789"
	result := make([]byte, length)

	randomBytes := make([]byte, length)
	if _, err := cryptorand.Read(randomBytes); err != nil {
		// Fallback to time-based generation if crypto/rand fails
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := range result {
			result[i] = digits[r.Intn(len(digits))]
		}
		return string(result)
	}

	for i := range result {
		result[i] = digits[randomBytes[i]%byte(len(digits))]
	}
	return string(result)
}

// GenerateSystemJWT generates a JWT token for system API calls
func GenerateSystemJWT(secret string) (string, error) {
	// Create claims for system use
	claims := TokenClaim{
		UserId: 78137, // System user ID from TypeScript implementation
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "LIS-Core",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and return
	return token.SignedString([]byte(secret))
}

// PostJSON sends a POST request with JSON body
func PostJSON(ctx context.Context, endpoint string, body interface{}, headers map[string]string) error {
	client := &http.Client{Timeout: 30 * time.Second}
	return PostJSONWithClient(ctx, client, endpoint, body, headers)
}

// PostJSONWithClient sends a POST request with JSON body using a provided HTTP client
func PostJSONWithClient(ctx context.Context, client *http.Client, endpoint string, body interface{}, headers map[string]string) error {
	// Marshal the body
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
