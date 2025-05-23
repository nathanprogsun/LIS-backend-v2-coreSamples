package handler

import (
	"coresamples/common"
	pb "coresamples/proto"
	"coresamples/service"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const PNSGuestLogin = "PNS guest"

type PatientHTTPHandler struct {
	patientService service.IPatientService
}

func NewPatientHTTPHandler(service service.IPatientService) *PatientHTTPHandler {
	return &PatientHTTPHandler{
		patientService: service,
	}
}

func (h *PatientHTTPHandler) PatientGuestLogIn(c *gin.Context) {
	req := &pb.PatientGuestLoginRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":         400,
			"errorMessage": err.Error(),
		})
		return
	}
	username := fmt.Sprintf("%s_%s_%s_%s", req.PatientFirstName, req.PatientLastName, req.PatientBirthdate, req.PatientAccessionId)
	token, patient, err := h.patientService.GuestPatientLogIn(req.PatientAccessionId, req.PatientBirthdate, req.PatientFirstName, req.PatientLastName, c)
	if err != nil {
		var code int32
		if errors.Is(err, service.PatientFoundWithNoDoB) {
			code = 422 // Custom code for missing DOB
		} else {
			code = 403 // Forbidden for mismatched information
		}
		c.JSON(http.StatusForbidden, gin.H{
			"code":         code,
			"errorMessage": err.Error(),
		})
		err = h.patientService.LogPatientLogin(username, c.ClientIP(), PNSGuestLogin, token, err, c)
		if err != nil {
			common.Error(err)
		}
		return
	}

	// Extract exp from the token
	parsedToken, _, _ := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	exp := int64(parsedToken.Claims.(jwt.MapClaims)["exp"].(float64))

	resp := &pb.PatientGuestLoginResponse{
		Token:      token,
		PatientId:  int32(patient.ID),
		Expiration: exp,
		Code:       200,
	}

	err = h.patientService.LogPatientLogin(username, c.ClientIP(), PNSGuestLogin, token, err, c)
	if err != nil {
		common.Error(err)
	}
	c.JSON(http.StatusOK, resp)
}
