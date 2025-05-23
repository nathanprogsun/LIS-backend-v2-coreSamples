package handler

import (
	pb "coresamples/proto"
	"coresamples/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrderHTTPHandler struct {
	orderService service.IOrderService
}

func NewOrderHTTPHandler(service service.IOrderService) *OrderHTTPHandler {
	return &OrderHTTPHandler{
		orderService: service,
	}
}

func (h *OrderHTTPHandler) UpdateOrderKitStatus(c *gin.Context) {
	req := &pb.UpdateOrderKitStatusRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":         400,
			"errorMessage": err.Error(),
		})
		return
	}

	// Call the service to update the order's kit status
	err := h.orderService.UpdateOrderKitStatusByAccessionId(req.AccessionId, req.KitStatus, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":         500,
			"errorMessage": err.Error(),
		})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Order kit status updated successfully",
	})
}
