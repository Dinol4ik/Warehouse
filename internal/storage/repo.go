package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"warehouse/internal/entity"
)

type Repo interface {
	ReserveItem(context.Context, string, int, string) error
	UnReserveItem(context.Context, string, int, string) error
	FetchWarehousesWithItemUnreserved(context.Context, string) (entity.WarehousesWithItemUnreserved, error)
	GetItemUnreservedByCode(context.Context, string) (int, error)
	FetchWarehousesWithItemReserved(context.Context, string) (entity.WarehousesWithItemReserved, error)
	FetchItemsByWarehouseId(context.Context, string) (entity.FetchWarehouseItemsResponse, error)
	FetchItemsByCodes(context.Context, []string) (entity.FetchItemsByCodesResponse, error)
	BeginTx(context.Context) (*sqlx.Tx, error)
	CommitTx(context.Context) error
	RollbackTx(context.Context) error
}
