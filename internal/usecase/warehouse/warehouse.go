package warehouse

import (
	"context"
	"go.uber.org/zap"
	"warehouse/internal/entity"
	"warehouse/internal/storage"
	"warehouse/internal/storage/repo"
)

type UseCase struct {
	repo   storage.Repo
	logger *zap.SugaredLogger
}

func NewWarehouse(r storage.Repo, logger *zap.SugaredLogger) *UseCase {
	return &UseCase{
		repo:   r,
		logger: logger,
	}
}

func (u *UseCase) ReserveItem(ctx context.Context, params entity.ReserveItemRequest) (
	entity.ReserveItemResponse, error) {
	res := entity.ReserveItemResponse{}
	for _, item := range params.Items {

		tx, err := u.repo.BeginTx(ctx)
		ctx = repo.InjectTx(ctx, tx)
		if err != nil {
			u.repo.RollbackTx(ctx)
			return entity.ReserveItemResponse{}, err
		}
		unreserved, err := u.repo.GetItemUnreservedByCode(ctx, item.Code)
		if err != nil {
			res.Errors = append(res.Errors, item.Code)
			u.repo.RollbackTx(ctx)
			continue
		}

		needReserve := item.Amount
		if unreserved > 0 {
			var warehouses entity.WarehousesWithItemUnreserved
			warehouses, err = u.repo.FetchWarehousesWithItemUnreserved(ctx, item.Code)
			if err != nil {
				res.Errors = append(res.Errors, item.Code)
				return entity.ReserveItemResponse{}, err
			}
			for _, warehouse := range warehouses.Warehouses {
				needReserve, err = u.ReserveWarehouseItem(ctx, item.Code, needReserve, warehouse)
				if err != nil {
					u.logger.Errorf("UseCase.reserveItemWarehouse Error: %s", err)
					continue
				}
				if needReserve == 0 {
					res.Success = append(res.Success, item.Code)
				}
			}
		}

		if unreserved < item.Amount || needReserve != 0 {
			u.repo.RollbackTx(ctx)
			res.Errors = append(res.Errors, item.Code)
			continue
		}

		err = u.repo.CommitTx(ctx)
		if err != nil {
			return entity.ReserveItemResponse{}, err
		}
	}
	return res, nil
}

func (u *UseCase) ReserveWarehouseItem(ctx context.Context, itemCode string, needReserve int,
	warehouse entity.WarehouseWithItem) (int, error) {
	if warehouse.Unreserved >= needReserve {
		err := u.repo.ReserveItem(ctx, itemCode, needReserve, warehouse.ID)
		if err != nil {
			return needReserve, err
		}

	} else {
		err := u.repo.ReserveItem(ctx, itemCode, warehouse.Unreserved, warehouse.ID)
		if err != nil {
			return needReserve, err
		}
		return needReserve - warehouse.Unreserved, nil
	}

	return 0, nil
}

func (u *UseCase) FetchItemsWarehouse(ctx context.Context, params entity.FetchWarehouseItemsRequest) (
	entity.FetchWarehouseItemsResponse, error) {
	res, err := u.repo.FetchItemsByWarehouseId(ctx, params.ID)
	if err != nil {
		return entity.FetchWarehouseItemsResponse{}, err
	}

	return res, nil
}

func (u *UseCase) UnReserveItem(ctx context.Context, params entity.UnReserveItemRequest) (
	entity.UnReserveItemResponse, error) {
	res := entity.UnReserveItemResponse{}
	for _, item := range params.Items {

		tx, err := u.repo.BeginTx(ctx)
		ctx = repo.InjectTx(ctx, tx)
		if err != nil {
			u.repo.RollbackTx(ctx)
			return entity.UnReserveItemResponse{}, err
		}

		needUnReserved := item.Amount
		warehouses, err := u.repo.FetchWarehousesWithItemReserved(ctx, item.Code)
		if err != nil {
			res.Errors = append(res.Errors, item.Code)
			u.repo.RollbackTx(ctx)
			return entity.UnReserveItemResponse{}, err
		}
		for _, warehouse := range warehouses.Warehouses {
			if needUnReserved > 0 {
				needUnReserved, err = u.UnReserveItemWarehouse(ctx, item.Code, needUnReserved, warehouse)
				if err != nil {
					u.logger.Errorf("UseCase.reserveItemWarehouse Error: %s", err)
					continue
				}
			} else {
				break
			}
		}
		if needUnReserved == 0 {
			res.Success = append(res.Success, item.Code)
		} else {
			u.repo.RollbackTx(ctx)
			res.Errors = append(res.Errors, item.Code)
			continue
		}

		err = u.repo.CommitTx(ctx)
		if err != nil {
			return entity.UnReserveItemResponse{}, err
		}

	}
	return res, nil
}

func (u *UseCase) UnReserveItemWarehouse(ctx context.Context, itemCode string,
	needUnReserved int, warehouse entity.WarehouseWithItem) (int, error) {
	if warehouse.Reserved >= needUnReserved {
		err := u.repo.UnReserveItem(ctx, itemCode, needUnReserved, warehouse.ID)
		if err != nil {
			return needUnReserved, err
		}

	} else {
		err := u.repo.UnReserveItem(ctx, itemCode, warehouse.Reserved, warehouse.ID)
		if err != nil {
			return needUnReserved, err
		}
		return needUnReserved - warehouse.Reserved, nil
	}

	return 0, nil
}
func (u *UseCase) FetchItemsByCodes(ctx context.Context, params entity.FetchItemsByCodesRequest) (entity.FetchItemsByCodesResponse, error) {

	res, err := u.repo.FetchItemsByCodes(ctx, params.Codes)
	if err != nil {
		return res, err
	}

	return res, nil
}
