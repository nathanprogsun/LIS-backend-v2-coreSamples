package service

import (
	"context"
	"coresamples/common"
	"coresamples/dbutils"
	"coresamples/ent"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var PatientGuestLoginDenied = fmt.Errorf("patient not found, login denied")
var PatientFoundWithNoDoB = fmt.Errorf("patient record exists but does not have a date of birth")

type IPatientService interface {
	GuestPatientLogIn(accessionId string, birthday string, firstName string, lastName string, ctx context.Context) (string, *ent.Patient, error)
	LogPatientLogin(username string, ip string, loginPortal string, token string, err error, ctx context.Context) error
}

type PatientService struct {
	Service
	secret string
}

func newPatientService(dbClient *ent.Client, redisClient *common.RedisClient, secret string) IPatientService {
	return &PatientService{
		Service: InitService(dbClient, redisClient),
		secret:  secret,
	}
}

func (s *PatientService) GuestPatientLogIn(accessionId string, birthday string, firstName string, lastName string, ctx context.Context) (string, *ent.Patient, error) {
	patient, err := dbutils.GetGuestPatient(accessionId, birthday, firstName, lastName, s.dbClient, ctx)
	if patient == nil || err != nil {
		return "", nil, PatientGuestLoginDenied
	}
	if patient.PatientBirthdate == "" {
		return "", nil, PatientFoundWithNoDoB
	}
	now := time.Now()

	order, err := dbutils.GetOrderByAccessionId(accessionId, s.dbClient, ctx)
	if err != nil {
		return "", nil, err
	}

	sample, err := dbutils.GetSampleByAccessionId(accessionId, s.dbClient, ctx)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role":        "patient",
		"customer_id": sample.CustomerID,
		"patient_id":  patient.ID,
		"iat":         now.Unix(),
		"exp":         now.Add(2 * time.Hour).Unix(),
		"clinic_id":   order.ClinicID,
		"barcode":     accessionId,
	})
	tokenStr, err := token.SignedString([]byte(s.secret))
	return tokenStr, patient, err
}

func (s *PatientService) LogPatientLogin(
	username string,
	ip string,
	loginPortal string,
	token string,
	err error,
	ctx context.Context) error {
	if err != nil {
		return dbutils.CreateLoginHistory(username, ip, false, err.Error(), loginPortal, "", s.dbClient, ctx)
	}
	return dbutils.CreateLoginHistory(username, ip, true, "", loginPortal, token, s.dbClient, ctx)
}
