package usecase

import (
	"context"
	"warehouse/internal/entity"
)

type Warehouse interface {
	ReserveItem(context.Context, entity.ReserveItemRequest) (entity.ReserveItemResponse, error)
	FetchItemsWarehouse(context.Context, entity.FetchWarehouseItemsRequest) (entity.FetchWarehouseItemsResponse, error)
	UnReserveItem(context.Context, entity.UnReserveItemRequest) (entity.UnReserveItemResponse, error)
	FetchItemsByCodes(context.Context, entity.FetchItemsByCodesRequest) (entity.FetchItemsByCodesResponse, error)
}
