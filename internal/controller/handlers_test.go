package controller

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"regexp"
	"testing"
	"warehouse/internal/entity"
	"warehouse/internal/storage/repo"
	warehouse "warehouse/internal/usecase/warehouse"
)

func newLogger() *zap.SugaredLogger {
	atom := zap.NewAtomicLevel()
	atom.SetLevel(zapcore.Level(-1))
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	zapCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		os.Stdout,
		atom,
	)
	logger := zap.New(zapCore)
	logger = logger.With(zap.String("service", "warehouse"))
	log := logger.Sugar()
	return log
}

func TestWarehouseService_FetchWarehouseItemsHandler(t *testing.T) {
	tests := []struct {
		name             string
		request          *entity.RPCRequest
		expectedResponse *entity.RPCResult
		willReturnRows   entity.FetchWarehouseItemsResponse
	}{
		{
			name:             "OK",
			request:          &entity.RPCRequest{Data: `{"id":"uuidgeneration"}`},
			expectedResponse: &entity.RPCResult{Data: `{"items":[{"name":"name","size":"size","code":"code","amount":14,"reserved":10},{"name":"name2","size":"size2","code":"code2","amount":24,"reserved":14}]}`},
			willReturnRows: entity.FetchWarehouseItemsResponse{Items: []entity.ItemInfo{
				{
					Name:     "name",
					Size:     "size",
					Code:     "code",
					Amount:   14,
					Reserved: 10,
				},
				{
					Name:     "name2",
					Size:     "size2",
					Code:     "code2",
					Amount:   24,
					Reserved: 14,
				},
			}},
		},

		{
			name:             "Non existing warehouse id",
			request:          &entity.RPCRequest{Data: `{"id":"asdasdasdasdas"}`},
			expectedResponse: &entity.RPCResult{Data: `{"items":null}`},
			willReturnRows:   entity.FetchWarehouseItemsResponse{Items: []entity.ItemInfo{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := newLogger()
			con, mock := repo.NewDBMock(t)
			repository := repo.New(con)
			ws := &Server{
				logger:       log,
				experimentUC: warehouse.NewExperimentUsecase(repository, log),
			}

			response := &entity.RPCResult{}
			willReturnRows := sqlmock.NewRows([]string{"name", "size", "code", "amount", "reserved"})
			for _, row := range tt.willReturnRows.Items {
				willReturnRows.AddRow(row.Name, row.Size, row.Code, row.Amount, row.Reserved)
			}
			mock.ExpectQuery(regexp.QuoteMeta(repo.FetchItemsByWarehouseIdSql)).WithArgs(sqlmock.AnyArg()).WillReturnRows(willReturnRows)
			if err := ws.FetchWarehouseItemsHandler(tt.request, response); err != nil {
				t.Errorf("FetchWarehouseItemsHandler() error = %v", err)
			}

			assert.Equal(t, tt.expectedResponse, response)
		})
	}
}

func TestWarehouseService_ReserveItemHandler(t *testing.T) {
	tests := []struct {
		name             string
		request          *entity.RPCRequest
		expectedResponse *entity.RPCResult
		warehouses       entity.WarehousesWithItemUnreserved
	}{
		{
			name:             "OK",
			request:          &entity.RPCRequest{Data: `{"items":[{"code":"AN0145YT5ZKKINS","amount":8}]}`},
			expectedResponse: &entity.RPCResult{Data: `{"successes":null,"errors":["AN0145YT5ZKKINS"]}`},
			warehouses: entity.WarehousesWithItemUnreserved{Warehouses: []entity.WarehouseWithItem{
				{
					ID:         "719942d3-7831-420b-951f-5a6aa784ce11",
					Unreserved: 7,
				},
				{
					ID:         "0bf05571-a646-45a3-b43f-f4832495807f",
					Unreserved: 3,
				},
			}},
		},

		{
			name:             "Non existing item code",
			request:          &entity.RPCRequest{Data: `{"items":[{"code":"AN0145YT5ZKKIS","amount":8}]}`},
			expectedResponse: &entity.RPCResult{Data: `{"successes":null,"errors":["AN0145YT5ZKKIS"]}`},
			warehouses: entity.WarehousesWithItemUnreserved{Warehouses: []entity.WarehouseWithItem{
				{
					ID:         "719942d3-7831-420b-951f-5a6aa784ce11",
					Unreserved: 7,
				},
				{
					ID:         "0bf05571-a646-45a3-b43f-f4832495807f",
					Unreserved: 3,
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := newLogger()
			con, mock := repo.NewDBMock(t)
			repository := repo.New(con)
			ws := &Server{
				logger:       log,
				experimentUC: warehouse.NewExperimentUsecase(repository, log),
			}

			response := &entity.RPCResult{}
			willReturnRowsUnreserved := sqlmock.NewRows([]string{"unreserved"})
			willReturnRowsUnreserved.AddRow("10")

			willReturnRowsWarehouses := sqlmock.NewRows([]string{"id", "unreserved"})
			for _, war := range tt.warehouses.Warehouses {
				willReturnRowsWarehouses.AddRow(war.ID, war.Unreserved)
			}

			mock.ExpectBegin()
			mock.ExpectQuery(regexp.QuoteMeta(repo.FetchItemsByCodesSql)).WithArgs(sqlmock.AnyArg()).WillReturnRows(willReturnRowsUnreserved)
			mock.ExpectQuery(regexp.QuoteMeta(repo.FetchWarehousesWithItemUnreservedSql)).WithArgs(sqlmock.AnyArg()).WillReturnRows(willReturnRowsWarehouses)
			mock.ExpectExec(regexp.QuoteMeta(repo.ReserveItemInWarehouseSql)).WithArgs(driver.Value(7), driver.Value("AN0145YT5ZKKINS"), driver.Value("719942d3-7831-420b-951f-5a6aa784ce11")).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(regexp.QuoteMeta(repo.ReserveItemInWarehouseSql)).WithArgs(driver.Value(1), driver.Value("AN0145YT5ZKKINS"), driver.Value("0bf05571-a646-45a3-b43f-f4832495807f")).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			if err := ws.ReserveItemHandler(tt.request, response); err != nil {
				t.Errorf("ReserveItemHandler() error = %v", err)
			}

			assert.Equal(t, tt.expectedResponse, response)
		})
	}
}

func TestWarehouseService_UnreserveItemHandler(t *testing.T) {
	tests := []struct {
		name             string
		request          *entity.RPCRequest
		expectedResponse *entity.RPCResult
		warehouses       entity.WarehousesWithItemReserved
	}{
		{
			name:             "OK",
			request:          &entity.RPCRequest{Data: `{"items":[{"code":"AN0145YT5ZKKINS","amount":8}]}`},
			expectedResponse: &entity.RPCResult{Data: `{"successes":["AN0145YT5ZKKINS"],"errors":null}`},
			warehouses: entity.WarehousesWithItemReserved{Warehouses: []entity.WarehouseWithItem{
				{
					ID:       "719942d3-7831-420b-951f-5a6aa784ce11",
					Reserved: 7,
				},
				{
					ID:       "0bf05571-a646-45a3-b43f-f4832495807f",
					Reserved: 3,
				},
			}},
		},

		{
			name:             "Non existing item code",
			request:          &entity.RPCRequest{Data: `{"items":[{"code":"AN0145YT5ZKKIS","amount":8}]}`},
			expectedResponse: &entity.RPCResult{Data: `{"successes":null,"errors":["AN0145YT5ZKKIS"]}`},
			warehouses: entity.WarehousesWithItemReserved{Warehouses: []entity.WarehouseWithItem{
				{
					ID:       "719942d3-7831-420b-951f-5a6aa784ce11",
					Reserved: 7,
				},
				{
					ID:       "0bf05571-a646-45a3-b43f-f4832495807f",
					Reserved: 1,
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := newLogger()
			con, mock := repo.NewDBMock(t)
			repository := repo.New(con)
			ws := &Server{
				logger:       log,
				experimentUC: warehouse.NewExperimentUsecase(repository, log),
			}

			response := &entity.RPCResult{}

			willReturnRowsWarehouses := sqlmock.NewRows([]string{"id", "reserved"})
			for _, war := range tt.warehouses.Warehouses {
				willReturnRowsWarehouses.AddRow(war.ID, war.Reserved)
			}

			mock.ExpectBegin()
			mock.ExpectQuery(regexp.QuoteMeta(repo.FetchWarehousesWithItemReservedSql)).WithArgs(sqlmock.AnyArg()).WillReturnRows(willReturnRowsWarehouses)
			mock.ExpectExec(regexp.QuoteMeta(repo.UnreservedItemSql)).WithArgs(driver.Value(7), driver.Value("AN0145YT5ZKKINS"), driver.Value("719942d3-7831-420b-951f-5a6aa784ce11")).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(regexp.QuoteMeta(repo.UnreservedItemSql)).WithArgs(driver.Value(1), driver.Value("AN0145YT5ZKKINS"), driver.Value("0bf05571-a646-45a3-b43f-f4832495807f")).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			if err := ws.UnReserveItemHandler(tt.request, response); err != nil {
				t.Errorf("UnreserveItemHandler() error = %v", err)
			}

			assert.Equal(t, tt.expectedResponse, response)
		})
	}
}
