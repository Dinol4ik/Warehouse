package controller

import (
	"context"
	"encoding/json"
	"warehouse/internal/entity"
)

func (s *Server) ReserveItemHandler(request *entity.RPCRequest, response *entity.RPCResult) error {
	s.logger.Infof("Received request on || ReserveItemHandler || with parameters: %v", *request)
	params := entity.ReserveItemRequest{}
	err := json.Unmarshal([]byte(request.Data), &params)
	if err != nil {
		return err
	}

	ctx := context.Background()
	res, err := s.warehouseUC.ReserveItem(ctx, params)
	marshal, err := json.Marshal(res)
	if err != nil {
		return err
	}

	*response = entity.RPCResult{Data: string(marshal)}
	return nil
}

func (s *Server) FetchWarehouseItemsHandler(request *entity.RPCRequest, response *entity.RPCResult) error {
	s.logger.Infof("Received request on || FetchWarehouseItemsHandler || with parameters: %v", *request)
	params := entity.FetchWarehouseItemsRequest{}
	err := json.Unmarshal([]byte(request.Data), &params)
	if err != nil {
		return err
	}

	ctx := context.Background()
	res, err := s.warehouseUC.FetchItemsWarehouse(ctx, params)
	if err != nil {
		return err
	}

	marshal, err := json.Marshal(res)
	if err != nil {
		return err
	}

	*response = entity.RPCResult{Data: string(marshal)}

	return nil
}

func (s *Server) UnReserveItemHandler(request *entity.RPCRequest, response *entity.RPCResult) error {
	s.logger.Infof("Received request on || UnreserveItemHandler || with parameters: %v", *request)

	params := entity.UnReserveItemRequest{}
	err := json.Unmarshal([]byte(request.Data), &params)
	if err != nil {
		return err
	}

	ctx := context.Background()
	res, err := s.warehouseUC.UnReserveItem(ctx, params)
	if err != nil {
		return err
	}

	marshal, err := json.Marshal(res)
	if err != nil {
		return err
	}

	*response = entity.RPCResult{Data: string(marshal)}
	return nil
}
func (s *Server) FetchItemsByCodesHandler(request *entity.RPCRequest, response *entity.RPCResult) error {
	s.logger.Infof("Received request on || FetchItemsByCodesHandler || with parameters: %v", *request)

	params := entity.FetchItemsByCodesRequest{}
	err := json.Unmarshal([]byte(request.Data), &params)
	if err != nil {
		return err
	}

	ctx := context.Background()
	res, err := s.warehouseUC.FetchItemsByCodes(ctx, params)
	if err != nil {
		return err
	}

	marshal, err := json.Marshal(res)
	if err != nil {
		return err
	}

	*response = entity.RPCResult{Data: string(marshal)}
	return nil
}
