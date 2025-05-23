package handler

import (
	"context"
	"coresamples/common"
	pb "coresamples/proto"
	"coresamples/service"
)

type SalesHandler struct {
	SalesService service.ISalesService
}

func (h *SalesHandler) GetSalesByTerritory(ctx context.Context, req *pb.Territory, resp *pb.SalesName) error {
	name, err := h.SalesService.GetSalesByTerritory(req.Zipcode, req.Country, req.State, ctx)
	resp.Name = name
	if err != nil {
		common.Error(err)
	}
	return nil
}
