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
)

type OrderService struct {
	secret string
}

type Order struct {
	OrderCreatedDate string `json:"order_created_date,omitempty"`
	ClinicID         int64  `json:"clinic_id,omitempty"`
}

type TubeRequirements struct {
	NumberOfTubes map[string]int32 `json:"noOfTubes,omitempty"`
}

var orderSvc *OrderService

func InitOrderService(secret string) {
	if orderSvc == nil {
		orderSvc = &OrderService{
			secret: secret,
		}
	}
}

func GetOrderService() *OrderService {
	return orderSvc
}

func (s *OrderService) GetLatestOrders(patientID int64) (map[string]*Order, error) {
	endpoint := common.EndpointsInfo.Order + "/orderTest/orderV2"

	input := struct {
		Patient []string `json:"patientIdList,omitempty"`
	}{
		Patient: []string{strconv.FormatInt(patientID, 10)},
	}

	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"internal_user_role": "admin",
	}).SignedString(s.secret)
	bf, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(bf))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		common.Error(fmt.Errorf("%v, request: %v, token: %s", err, input, token))
		return nil, err
	}
	if resp.StatusCode >= 400 {
		common.Error(fmt.Errorf("response status code %d: %s, %v", resp.StatusCode, endpoint, input))
		return nil, errors.New("response status code " + strconv.Itoa(resp.StatusCode))
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	output := map[string]*Order{}

	if err := json.Unmarshal(body, &output); err != nil {
		common.Error(fmt.Errorf("error parsing response body: %s", body))
		return nil, nil
	}
	return output, nil
}

func (s *OrderService) GetRequiredTubesBySampleId(sampleId int) (*TubeRequirements, error) {
	endpoint := common.EndpointsInfo.Order + fmt.Sprintf("/tubeInfo?sampleId=%d", sampleId)
	req, _ := http.NewRequest("GET", endpoint, nil)
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"internal_user_role": "admin",
	}).SignedString(s.secret)
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		common.Error(fmt.Errorf("%v, request: %s, token: %s", err, endpoint, token))
		return nil, err
	}
	if resp.StatusCode >= 400 {
		common.Error(fmt.Errorf("response status code %d: %s", resp.StatusCode, endpoint))
		return nil, errors.New("response status code " + strconv.Itoa(resp.StatusCode))
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	output := &TubeRequirements{}
	if err := json.Unmarshal(body, &output); err != nil {
		common.Error(fmt.Errorf("error parsing response body: %s", body))
		return nil, nil
	}
	return output, nil
}
