package external

import (
	"bytes"
	"coresamples/common"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
	"strconv"
	"time"
)

type ChargingService struct {
	secret string
}

type PaymentMethod struct {
	ID            int64  `json:"id,omitempty"`
	AccountID     int64  `json:"account_id,omitempty"`
	AccountType   string `json:"account_type,omitempty"`
	Type          string `json:"type,omitempty"`
	TokenPlatform string `json:"token_platform,omitempty"`
	PaymentToken  string `json:"payment_token,omitempty"`
	CardType      string `json:"card_type,omitempty"`
	ExpiryDate    string `json:"expiry_date,omitempty"`
	LastFour      string `json:"last_four,omitempty"`
	Subscription  bool   `json:"is_subscription,omitempty"`
	CustomerToken string `json:"customer_token,omitempty"`
}

type PaymentRecord struct {
	AccountID     string `json:"account_id,omitempty"`
	AccountType   string `json:"account_type,omitempty"`
	Amount        string `json:"amount,omitempty"`
	CustomerToken string `json:"customer_id_token,omitempty"`
	ExpiryDate    string `json:"expiry_date,omitempty"`
	LastFour      string `json:"last_four,omitempty"`
	PaymentToken  string `json:"payment_method_token,omitempty"`
	PaymentStatus string `json:"payment_status,omitempty"`
	PaymentType   string `json:"payment_type,omitempty"`
	PaymentSource string `json:"payment_source,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	TransactionID string `json:"transaction_id,omitempty"`
}

type Subscription struct {
	ID            int64         `json:"id,omitempty"`
	AccountID     int64         `json:"account_id,omitempty"`
	AccountType   string        `json:"account_type,omitempty"`
	Amount        float32       `json:"amount,omitempty"`
	Status        int32         `json:"status,omitempty"`
	Currency      string        `json:"currency,omitempty"`
	ChargeType    string        `json:"charge_type,omitempty"`
	ChargeTypeID  string        `json:"charge_type_id,omitempty"`
	StartAt       int64         `json:"start_at,omitempty"`
	EndAt         int64         `json:"end_at,omitempty"`
	Frequency     string        `json:"frequency,omitempty"`
	Interval      int32         `json:"interval,omitempty"`
	PaymentMethod PaymentMethod `json:"payment_method,omitempty"`
}

type SubscriptionInfo struct {
	ID            int64         `json:"id,omitempty"`
	AccountID     int64         `json:"account_id,omitempty"`
	AccountType   string        `json:"account_type,omitempty"`
	Amount        float32       `json:"amount,omitempty"`
	Status        int32         `json:"status"`
	Currency      string        `json:"currency,omitempty"`
	ChargeType    string        `json:"charge_type,omitempty"`
	ChargeTypeID  string        `json:"charge_type_id,omitempty"`
	StartAt       string        `json:"start_at,omitempty"`
	EndAt         string        `json:"end_at,omitempty"`
	LastRun       string        `json:"last_run,omitempty"`
	NextRun       string        `json:"next_run,omitempty"`
	CreatedAt     string        `json:"created_at,omitempty"`
	UpdatedAt     string        `json:"updated_at,omitempty"`
	Frequency     string        `json:"frequency,omitempty"`
	Interval      int32         `json:"interval,omitempty"`
	PaymentMethod PaymentMethod `json:"payment_method,omitempty"`
}

type PayRequest struct {
	AccountID     int64   `json:"account_id,omitempty"`
	AccountType   string  `json:"account_type,omitempty"`
	Amount        float32 `json:"amount,omitempty"`
	Currency      string  `json:"currency,omitempty"`
	ChargeType    string  `json:"charge_type,omitempty"`
	ChargeTypeID  string  `json:"charge_type_id,omitempty"`
	Type          string  `json:"type,omitempty"`
	TokenPlatform string  `json:"token_platform,omitempty"`
	PaymentSource string  `json:"payment_source,omitempty"`
	PaymentToken  string  `json:"payment_token,omitempty"`
	CustomerToken string  `json:"customer_token,omitempty"`
	Notes         string  `json:"notes,omitempty"`
}

type PayResponse struct {
	TransactionID string `json:"payment_transaction_id,omitempty"`
}

type Transactions struct {
	TransactionInfos []TransactionInfo `json:"transaction_infos,omitempty"`
}

type TransactionInfo struct {
	PaymentType string  `json:"payment_type,omitempty"`
	Amount      float32 `json:"amount,omitempty"`
	Type        string  `json:"type,omitempty"`
	LastFour    string  `json:"last_four,omitempty"`
	CreatedAt   string  `json:"created_at,omitempty"`
	Status      string  `json:"status,omitempty"`
	CardType    string  `json:"card_type,omitempty"`
}

var chargeSvc *ChargingService

func InitChargingService(secret string) {
	if chargeSvc == nil {
		chargeSvc = &ChargingService{
			secret: secret,
		}
	}
}

func GetChargingService() *ChargingService {
	return chargeSvc
}

func (s *ChargingService) GenerateToken(accountType string, accountID int64) (string, error) {
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

func (s *ChargingService) CreateSubscription(accountType string, accountID int64, amount float32, subscriptionID int, startAt time.Time, billingCycle int32, interval string, paymentID int64, status int32) error {
	// try deleting the subscription first in case there's a dead subscription
	endpoint := common.EndpointsInfo.Charging + "subscription"
	input := Subscription{
		AccountID:   accountID,
		AccountType: accountType,
		Amount:      amount, //res, err := external.GetOrderService().GetLatestOrders(678293)
		//if err != nil {
		//	t.Fatal(err)
		//}
		//common.Info(res["2047179"])
		Currency:      "usd",
		ChargeType:    "subscription",
		ChargeTypeID:  strconv.Itoa(subscriptionID),
		StartAt:       startAt.Unix(),
		Frequency:     interval,
		Interval:      billingCycle,
		Status:        status,
		PaymentMethod: PaymentMethod{ID: paymentID},
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

func (s *ChargingService) DeleteSubscription(subscriptionID int, accountType string, accountID int64) error {
	endpoint := common.EndpointsInfo.Charging + "subscription/" + strconv.Itoa(subscriptionID)
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
	return nil
}

func (s *ChargingService) PauseSubscription(subscriptionID int, accountType string, accountID int64) error {
	endpoint := common.EndpointsInfo.Charging + "subscription/edit"
	input := Subscription{
		AccountID:    accountID,
		ChargeTypeID: strconv.Itoa(subscriptionID),
		Status:       0,
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

func (s *ChargingService) ResumeSubscription(subscriptionID int, accountType string, accountID int64, startAt time.Time, paymentMethodID int64) error {
	endpoint := common.EndpointsInfo.Charging + "subscription/edit"
	input := Subscription{
		AccountID:     accountID,
		AccountType:   accountType,
		ChargeTypeID:  strconv.Itoa(subscriptionID),
		Status:        1,
		StartAt:       startAt.Unix(),
		PaymentMethod: PaymentMethod{ID: paymentMethodID},
	}
	bf, _ := json.Marshal(input)
	common.Infof(string(bf))
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

func (s *ChargingService) UpdatePaymentMethod(subscriptionID int, paymentID int64, accountType string, accountID int64) error {
	endpoint := common.EndpointsInfo.Charging + "subscription/edit"
	input := Subscription{
		AccountID:     accountID,
		AccountType:   accountType,
		ChargeTypeID:  strconv.Itoa(subscriptionID),
		PaymentMethod: PaymentMethod{ID: paymentID},
		Status:        1,
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

func (s *ChargingService) UpdateSubscriptionBillingCycle(subscriptionID int, accountType string, accountID int64, billingCycle int32, fee float32) error {
	endpoint := common.EndpointsInfo.Charging + "subscription/edit"
	input := Subscription{
		AccountID:    accountID,
		ChargeTypeID: strconv.Itoa(subscriptionID),
		Interval:     billingCycle,
		Amount:       fee,
		Frequency:    "monthly",
		Status:       1,
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
		common.Error(fmt.Errorf("response status code %d: %s, %s", resp.StatusCode, endpoint, string(bf)))
		return errors.New("response status code " + strconv.Itoa(resp.StatusCode))
	}
	return nil
}

func (s *ChargingService) UpdateSubscriptionFee(subscriptionID int, accountType string, accountID int64, fee float32) error {
	endpoint := common.EndpointsInfo.Charging + "subscription/edit"
	input := Subscription{
		AccountID:    accountID,
		ChargeTypeID: strconv.Itoa(subscriptionID),
		Amount:       fee,
		Status:       1,
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

func (s *ChargingService) GetSubscription(subscriptionID int, accountType string, accountID int64) (*SubscriptionInfo, error) {
	endpoint := common.EndpointsInfo.Charging + "subscription/" + strconv.Itoa(subscriptionID)
	req, _ := http.NewRequest("GET", endpoint, nil)
	token, err := s.GenerateToken(accountType, accountID)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		common.Error(fmt.Errorf("response status code %d: %s", resp.StatusCode, endpoint))
		return nil, errors.New("response status code " + strconv.Itoa(resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	output := &SubscriptionInfo{}

	if err := json.Unmarshal(body, output); err != nil {
		return nil, err
	}
	return output, err
}

func (s *ChargingService) CreatePaymentMethod(method *PaymentMethod) (int64, error) {
	if !method.Subscription {
		return -1, errors.New("method type should be subscription")
	}
	endpoint := common.EndpointsInfo.Charging + "/paymentMethod"
	bf, _ := json.Marshal(method)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(bf))
	token, err := s.GenerateToken(method.AccountType, method.AccountID)
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
		common.Error(fmt.Errorf("response status code %d: %s, %v", resp.StatusCode, endpoint, method))
		return -1, errors.New("response status code " + strconv.Itoa(resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return -1, err
	}

	output := &PaymentMethod{}

	if err := json.Unmarshal(body, output); err != nil {
		return -1, err
	}
	return output.ID, err
}

func (s *ChargingService) DeletePaymentMethod(paymentMethodID int64, accountType string, accountID int64) error {
	endpoint := common.EndpointsInfo.Charging +
		"paymentMethod/" + strconv.FormatInt(paymentMethodID, 10) +
		"?account_id=" + strconv.FormatInt(accountID, 10) +
		"&account_type=" + accountType
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
	return nil
}

func (s *ChargingService) GetPaymentMethods(accountType string, accountID int64) ([]*PaymentMethod, error) {
	endpoint := common.EndpointsInfo.Charging + "paymentMethod/" + strconv.FormatInt(accountID, 10) + "/type/clinic?subscription=1&method_type=card"
	req, _ := http.NewRequest("GET", endpoint, nil)
	token, err := s.GenerateToken(accountType, accountID)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		common.Error(fmt.Errorf("response status code %d: %s", resp.StatusCode, endpoint))
		return nil, errors.New("response status code " + strconv.Itoa(resp.StatusCode))
	}

	var res []*PaymentMethod
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res, err
}

func (s *ChargingService) OneTimeCharge(accountType string, accountID int64, amount float32, subscriptionID int, platform string, paymentToken string, customerToken string) (bool, int64, error) {
	success := false
	strAmount := decimal.NewFromFloat(float64(amount)).String()
	_, err := GetAccountingService().CreateCharge(accountID, accountType, strAmount, strconv.Itoa(subscriptionID))
	if err != nil {
		return success, 0, err
	}
	payEndpoint := common.EndpointsInfo.Charging + "transaction/pay"
	payReq := PayRequest{
		AccountID:     accountID,
		AccountType:   accountType,
		Amount:        amount,
		Currency:      "usd",
		ChargeType:    "subscription",
		ChargeTypeID:  strconv.FormatInt(int64(subscriptionID), 10),
		Type:          "card",
		TokenPlatform: platform,
		PaymentSource: "portal",
		PaymentToken:  paymentToken,
		CustomerToken: customerToken,
		Notes:         "one time charge for subscription",
	}
	bf, _ := json.Marshal(payReq)
	req, _ := http.NewRequest("POST", payEndpoint, bytes.NewBuffer(bf))
	token, err := s.GenerateToken(accountType, accountID)
	if err != nil {
		return success, 0, err
	}
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return success, 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		common.Error(fmt.Errorf("response status code %d: %s, %v", resp.StatusCode, payEndpoint, payReq))
		return false, 0, errors.New("response status code " + strconv.Itoa(resp.StatusCode))
	}
	success = true
	payResp := &PayResponse{}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return success, 0, err
	}
	if err := json.Unmarshal(body, payResp); err != nil {
		return success, 0, err
	}
	record := PaymentRecord{
		AccountID:     strconv.FormatInt(accountID, 10),
		AccountType:   "clinic",
		Amount:        strAmount,
		CustomerToken: customerToken,
		PaymentToken:  paymentToken,
		PaymentStatus: "1",
		PaymentType:   "cc",
		PaymentSource: platform,
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
		TransactionID: payResp.TransactionID,
	}
	paymentID, err := GetAccountingService().ApplyCharge(record, strconv.Itoa(subscriptionID), accountID, accountType)
	return success, paymentID, err
}

func (s *ChargingService) GetSubscriptionTransactionInfos(subscriptionID int, accountType string, accountID int64) (*Transactions, error) {
	endpoint := common.EndpointsInfo.Charging + "/transaction/transactionInfoV2"
	request := struct {
		ChargeType   string `json:"charge_type,omitempty"`
		ChargeTypeID string `json:"charge_type_id,omitempty"`
	}{
		ChargeType:   "subscription",
		ChargeTypeID: strconv.Itoa(subscriptionID),
	}
	bf, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(bf))
	token, err := s.GenerateToken(accountType, accountID)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		common.Error(fmt.Errorf("response status code %d: %s, %v", resp.StatusCode, endpoint, request))
		return nil, errors.New("response status code " + strconv.Itoa(resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	txns := &Transactions{}

	if err := json.Unmarshal(body, txns); err != nil {
		return nil, err
	}

	return txns, nil
}
