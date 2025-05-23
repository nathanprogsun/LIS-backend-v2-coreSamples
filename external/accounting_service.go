package external

import (
	"bytes"
	"coresamples/common"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AccountingService struct {
	secret string
}

type AddCreditBody struct {
	AccountId   string `json:"account_id,omitempty"`
	AccountType string `json:"account_type,omitempty"`
	Amount      string `json:"amount,omitempty"`
	SourceType  string `json:"source_type,omitempty"`
	SourceId    string `json:"source_id,omitempty"`
	Status      string `json:"status,omitempty"`
	Cashable    bool   `json:"cashable,omitempty"`
	CustomerID  string `json:"customer_id,omitempty"`
}

type UpdateCreditBody struct {
	CreditID string `json:"credit_id,omitempty"`
	Status   string `json:"status,omitempty"`
}

type ChargeRecord struct {
	AccountID    string `json:"account_id,omitempty"`
	AccountType  string `json:"account_type,omitempty"`
	Amount       string `json:"amount,omitempty"`
	ChargeType   string `json:"charge_type,omitempty"`
	ChargeTypeID string `json:"charge_type_id,omitempty"`
	Notes        string `json:"notes,omitempty"`
	Quantity     string `json:"quantity,omitempty"`
	Type         string `json:"type,omitempty"`
	PaymentTerm  string `json:"payment_term,omitempty"`
	CreateAt     string `json:"create_at,omitempty"`
}

type ApplyChargeRecord struct {
	Record       PaymentRecord `json:"payment,omitempty"`
	ChargeType   string        `json:"charge_type,omitempty"`
	ChargeTypeID string        `json:"charge_type_id,omitempty"`
}

var accountingSvc *AccountingService

func InitAccountingService(secret string) {
	if accountingSvc == nil {
		accountingSvc = &AccountingService{
			secret: secret,
		}
	}
}

func GetAccountingService() *AccountingService {
	return accountingSvc
}

func (s *AccountingService) GenerateBearerToken(userID int64, clinicID int64, customerID int64, patientID int64) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role":        "customer",
		"userId":      userID,
		"clinic_id":   clinicID,
		"customer_id": customerID,
		"patient_id":  patientID,
		"iat":         now.Unix(),
		"exp":         now.AddDate(1, 0, 0).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}
	return "Bearer " + tokenString, nil
}

func (s *AccountingService) GenerateToken(accountType string, accountID int64) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"role": accountType,
		"iat":  now.Unix(),
		"exp":  now.AddDate(1, 0, 0).Unix(),
	}
	switch accountType {
	case "clinic":
		claims["clinic_id"] = accountID
	case "customer":
		claims["customer_id"] = accountID
	case "patient":
		claims["patient_id"] = accountID
	case "provider":
		claims["provider_id"] = accountID
	default:
		return "", errors.New("unrecognized account type" + accountType)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}
	return "Bearer " + tokenString, nil
}

func (s *AccountingService) AddCredits(credits string, orderID int64, customerID int64) (int64, error) {
	endpoint := common.EndpointsInfo.Accounting + "credit"
	input := AddCreditBody{
		AccountId:   strconv.FormatInt(customerID, 10),
		AccountType: "customer",
		Amount:      credits,
		SourceType:  "reward",
		SourceId:    strconv.FormatInt(orderID, 10),
		Status:      "2",
		Cashable:    true,
		CustomerID:  strconv.FormatInt(customerID, 10),
	}
	bf, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(bf))
	token, err := s.GenerateToken("customer", customerID)
	if err != nil {
		return -1, err
	}
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()
	if !strings.HasPrefix(resp.Status, "2") {
		common.Error(fmt.Errorf("response status code %d: %s, %v", resp.StatusCode, endpoint, input))
		return -1, errors.New("response status code " + strconv.Itoa(resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return -1, err
	}

	output := &struct {
		ID int64 `json:"id"`
	}{}

	if err := json.Unmarshal(body, output); err != nil {
		return -1, err
	}
	return output.ID, err
}

func (s *AccountingService) VoidCredits(creditId int64, accountID int64, accountType string) error {
	endpoint := common.EndpointsInfo.Accounting + "credit/" + strconv.FormatInt(creditId, 10)
	req, _ := http.NewRequest("DELETE", endpoint, nil)
	token, err := s.GenerateToken(accountType, accountID)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		common.Error(fmt.Errorf("response status code %d: %s", resp.StatusCode, endpoint))
		return errors.New("response status code " + strconv.Itoa(resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	output := &struct {
		Status string `json:"status"`
	}{}

	if err := json.Unmarshal(body, output); err != nil {
		return err
	}
	if output.Status != "success" {
		return errors.New("response is not success")
	}
	return nil
}

func (s *AccountingService) ActivateCredits(creditID int64, accountID int64, accountType string) error {
	endpoint := common.EndpointsInfo.Accounting + "credit/" + "update"
	input := UpdateCreditBody{
		Status:   "1",
		CreditID: strconv.FormatInt(creditID, 10),
	}
	bf, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(bf))
	token, err := s.GenerateToken(accountType, accountID)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		common.Error(fmt.Errorf("response status code %d: %s, %v", resp.StatusCode, endpoint, input))
		return errors.New("response status code " + strconv.Itoa(resp.StatusCode))
	}

	return nil
}

func (s *AccountingService) CreateCharge(accountID int64, accountType string, amount string, subscriptionID string) (int64, error) {
	endpoint := common.EndpointsInfo.Accounting + "charge"
	input := ChargeRecord{
		AccountID:    strconv.FormatInt(accountID, 10),
		AccountType:  accountType,
		Amount:       amount,
		ChargeType:   "subscription",
		ChargeTypeID: subscriptionID,
		Notes:        "one time subscription charge",
		Quantity:     "1",
		Type:         "surcharge",
		PaymentTerm:  "Net",
		CreateAt:     time.Now().Format(`"2006-01-02 15:01:05"`),
	}
	bf, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(bf))
	token, err := s.GenerateToken(accountType, accountID)
	if err != nil {
		return -1, err
	}
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		common.Error(fmt.Errorf("response status code %d: %s, %v, response: %v", resp.StatusCode, endpoint, input, resp))
		return -1, errors.New("response status code " + strconv.Itoa(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return -1, err
	}

	output := &struct {
		ID int64 `json:"id,omitempty"`
	}{}

	if err := json.Unmarshal(body, output); err != nil {
		return -1, err
	}
	return output.ID, nil
}

func (s *AccountingService) ApplyCharge(record PaymentRecord, subscriptionID string, accountID int64, accountType string) (int64, error) {
	endpoint := common.EndpointsInfo.Accounting + "apply/payment-apply"
	input := ApplyChargeRecord{
		Record:       record,
		ChargeTypeID: subscriptionID,
		ChargeType:   "subscription",
	}
	bf, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(bf))
	token, err := s.GenerateToken(accountType, accountID)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		common.Error(fmt.Errorf("response status code %d: %s, %v", resp.StatusCode, endpoint, input))
		return 0, errors.New("response status code " + strconv.Itoa(resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return 0, err
	}

	output := &struct {
		PaymentID int64 `json:"payment_id,omitempty"`
	}{}

	if err = json.Unmarshal(body, output); err != nil {
		return 0, err
	}
	return output.PaymentID, nil
}
