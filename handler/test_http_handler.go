package handler

import (
	pb "coresamples/proto"
	"coresamples/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TestHTTPHandler struct {
	service service.ITestService
}

func NewTestHTTPHandler(service service.ITestService) *TestHTTPHandler {
	return &TestHTTPHandler{service: service}
}

// GetTest handles fetching tests by IDs.
func (h *TestHTTPHandler) GetTest(c *gin.Context) {
	testIDs, err := parseIntSlice(c.QueryArray("test_ids"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": "invalid test_ids"})
		return
	}
	tests, err := h.service.GetTest(testIDs, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tests)
}

// GetTestField fetches test fields based on IDs and detail names.
func (h *TestHTTPHandler) GetTestField(c *gin.Context) {
	testIDs, err := parseIntSlice(c.QueryArray("test_ids"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": "invalid test_ids"})
		return
	}
	testDetailNames := c.QueryArray("test_detail_names")
	tests, err := h.service.GetTestField(testIDs, testDetailNames, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tests)
}

// CreateTest creates a new test.
func (h *TestHTTPHandler) CreateTest(c *gin.Context) {
	var req pb.CreateTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": "invalid request body"})
		return
	}
	test, err := h.service.CreateTest(&req, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, test)
}

// GetTestIDsFromTestCodes fetches test IDs by their codes.
func (h *TestHTTPHandler) GetTestIDsFromTestCodes(c *gin.Context) {
	testCodes := c.QueryArray("test_codes")
	if len(testCodes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": "test_codes cannot be empty"})
		return
	}
	ids, err := h.service.GetTestIDsFromTestCodes(testCodes, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ids)
}

// GetDuplicateAssayGroupTest fetches duplicate assay group test IDs.
func (h *TestHTTPHandler) GetDuplicateAssayGroupTest(c *gin.Context) {
	testID, err := strconv.Atoi(c.Query("test_id"))
	if err != nil || testID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": "invalid test_id"})
		return
	}
	ids, err := h.service.GetDuplicateAssayGroupTest(testID, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ids)
}

// parseIntSlice converts a slice of strings to a slice of integers.
func parseIntSlice(values []string) ([]int, error) {
	var result []int
	for _, v := range values {
		num, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		result = append(result, num)
	}
	return result, nil
}
