package entity

type RPCRequest struct {
	Data string
}

type RPCResult struct {
	Data string
}

type UnReserveItemRequest struct {
	Method string `json:"method"`
	Items  []Item `json:"items"`
}

type UnReserveItemResponse struct {
	Success []string `json:"successes"`
	Errors  []string `json:"errors"`
}

type ReserveItemRequest struct {
	Method string `json:"method"`
	Items  []Item `json:"items"`
}

type ReserveItemResponse struct {
	Success []string `json:"successes"`
	Errors  []string `json:"errors"`
}

type Item struct {
	Code   string `json:"code"`
	Amount int    `json:"amount"`
}

type WarehousesWithItemUnreserved struct {
	Warehouses []WarehouseWithItem `json:"warehouses"`
}

type WarehousesWithItemReserved struct {
	Warehouses []WarehouseWithItem `json:"warehouses"`
}

type WarehouseWithItem struct {
	ID         string `json:"id"`
	Reserved   int    `json:"reserved"`
	Unreserved int    `json:"unreserved"`
}

type FetchWarehouseItemsRequest struct {
	Method string `json:"method"`
	ID     string `json:"id"`
}

type FetchWarehouseItemsResponse struct {
	Items []ItemInfo `json:"items"`
}

type ItemInfo struct {
	Name               string `json:"name" db:"name"`
	Size               string `json:"size" db:"size"`
	Code               string `json:"code" db:"code"`
	Amount             int    `json:"amount" db:"amount"`
	Reserved           int    `json:"reserved" db:"reserved"`
	WarehouseID        string `json:"warehouseId,omitempty" db:"warehouse_id"`
	WarehouseAvailable bool   `json:"warehouseAvailable,omitempty" db:"is_available"`
}

type FetchItemsByCodesRequest struct {
	Codes []string `json:"codes"`
}

type FetchItemsByCodesResponse struct {
	Items []ItemInfo `json:"items"`
}
