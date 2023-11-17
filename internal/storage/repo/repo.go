package repo

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"warehouse/internal/entity"
)

type Repo struct {
	db *sqlx.DB
}

type querier interface {
	SelectContext(context.Context, interface{}, string, ...interface{}) error
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	GetContext(context.Context, interface{}, string, ...interface{}) error
}

// New -.
func New(db *sqlx.DB) *Repo {
	return &Repo{db}
}

func (r *Repo) GetRepo() *Repo {
	return r
}

const GetItemUnreservedByCodeSql = `SELECT SUM(amount) - SUM(reserved) AS unreserved FROM item
										INNER JOIN item_warehouse iw ON item.code = iw.item_code
										INNER JOIN warehouse w ON w.uuid = iw.warehouse_id
										WHERE code=$1 AND w.is_available GROUP BY iw.item_code`

func (r *Repo) GetItemUnreservedByCode(ctx context.Context, code string) (int, error) {
	var unreserved int

	var q querier = r.db
	if tx := extractTx(ctx); tx != nil {
		q = tx
	}

	err := q.GetContext(
		ctx,
		&unreserved,
		GetItemUnreservedByCodeSql,
		code,
	)

	if err != nil {
		return 0, err
	}

	return unreserved, err
}

const UnreservedItemSql = `UPDATE item_warehouse SET reserved = reserved-$1 WHERE item_code=$2 AND warehouse_id=$3`

func (r *Repo) UnReserveItem(ctx context.Context, itemCode string, quantity int, warehouseId string) error {
	var q querier = r.db
	if tx := extractTx(ctx); tx != nil {
		q = tx
	}

	_, err := q.ExecContext(
		ctx,
		UnreservedItemSql,
		quantity,
		itemCode,
		warehouseId,
	)

	if err != nil {
		return err
	}

	return err
}

const ReserveItemInWarehouseSql = `UPDATE item_warehouse SET reserved = reserved+$1
									WHERE item_code=$2
									AND warehouse_id=$3`

func (r *Repo) ReserveItem(ctx context.Context, itemCode string, amount int, warehouseId string) error {
	var q querier = r.db
	if tx := extractTx(ctx); tx != nil {
		q = tx
	}

	_, err := q.ExecContext(
		ctx,
		ReserveItemInWarehouseSql,
		amount,
		itemCode,
		warehouseId,
	)
	if err != nil {
		return err
	}

	return err
}

const FetchItemsByWarehouseIdSql = `SELECT name, size, code, (amount - reserved) as amount
									FROM item_warehouse
									INNER JOIN item iw ON iw.code = item_code
									WHERE warehouse_id = $1`

func (r *Repo) FetchItemsByWarehouseId(ctx context.Context, warehouseId string) (entity.FetchWarehouseItemsResponse, error) {
	res := entity.FetchWarehouseItemsResponse{}

	var q querier = r.db
	if tx := extractTx(ctx); tx != nil {
		q = tx
	}

	err := q.SelectContext(ctx,
		&res.Items,
		FetchItemsByWarehouseIdSql,
		warehouseId,
	)

	if err != nil {
		return res, err
	}

	return res, nil
}

const FetchWarehousesWithItemReservedSql = `SELECT warehouse_id AS id, reserved
											FROM item_warehouse iw
											INNER JOIN warehouse w ON w.uuid = iw.warehouse_id
											WHERE item_code = $1
											  AND reserved >= 0
											  AND w.is_available
											ORDER BY reserved DESC`

func (r *Repo) FetchWarehousesWithItemReserved(ctx context.Context, itemCode string) (entity.WarehousesWithItemReserved, error) {
	var res entity.WarehousesWithItemReserved

	var q querier = r.db
	if tx := extractTx(ctx); tx != nil {
		q = tx
	}

	err := q.SelectContext(ctx,
		&res.Warehouses,
		FetchWarehousesWithItemReservedSql,
		itemCode,
	)
	if err != nil {
		return res, err
	}

	return res, nil
}

const FetchWarehousesWithItemUnreservedSql = `SELECT warehouse_id AS id, amount - reserved AS unreserved
		FROM item_warehouse iw
				 INNER JOIN warehouse w ON w.uuid = iw.warehouse_id
		WHERE item_code = $1
		  AND amount - reserved > 0
		  AND w.is_available
		ORDER BY reserved DESC`

func (r *Repo) FetchWarehousesWithItemUnreserved(ctx context.Context, itemCode string) (entity.WarehousesWithItemUnreserved, error) {
	res := entity.WarehousesWithItemUnreserved{}

	var q querier = r.db
	if tx := extractTx(ctx); tx != nil {
		q = tx
	}

	err := q.SelectContext(ctx,
		&res.Warehouses,
		FetchWarehousesWithItemUnreservedSql,
		itemCode,
	)
	if err != nil {
		return res, err
	}

	return res, nil
}

const FetchItemsByCodesSql = `SELECT item.name, size, code, amount, reserved, w.is_available, w.uuid AS warehouse_id FROM item 
		                                    INNER JOIN item_warehouse iw ON item.code = iw.item_code
		                                    INNER JOIN warehouse w ON w.uuid = iw.warehouse_id
		                                    WHERE item_code = ANY($1)`

func (r *Repo) FetchItemsByCodes(ctx context.Context, codes []string) (entity.FetchItemsByCodesResponse, error) {
	res := entity.FetchItemsByCodesResponse{}

	var q querier = r.db
	if tx := extractTx(ctx); tx != nil {
		q = tx
	}

	err := q.SelectContext(ctx,
		&res.Items,
		FetchItemsByCodesSql,
		pq.Array(codes),
	)

	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *Repo) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *Repo) CommitTx(ctx context.Context) error {
	tx := extractTx(ctx)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) RollbackTx(ctx context.Context) error {
	tx := extractTx(ctx)
	err := tx.Rollback()
	if err != nil {
		return err
	}
	return nil
}

type txKey struct{}

func InjectTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return nil
}
