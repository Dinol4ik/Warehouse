package controller

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"net/rpc"
	"warehouse/internal/usecase"
)

type Server struct {
	logger      *zap.SugaredLogger
	warehouseUC usecase.Warehouse
}

func NewServer(logger zap.SugaredLogger, wUC usecase.Warehouse) *Server {
	return &Server{
		logger:      &logger,
		warehouseUC: wUC,
	}
}

func (s *Server) Run() error {
	r := mux.NewRouter()
	server := rpc.NewServer()
	err := server.RegisterName("WarehouseService", s)
	if err != nil {
		s.logger.Fatalf("WarehouseService.Run Error: %s", err)
	}
	r.Use()
	r.Handle("/rpc", server)

	return http.ListenAndServe("127.0.0.1:8080", r)
}
