package handler

import (
	"crtexBalance/internal/service"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "crtexBalance/docs"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", h.getBalance).Methods("POST")
	r.HandleFunc("/topup", h.replenishmentBalance).Methods("POST")
	r.HandleFunc("/transfer", h.transfer).Methods("POST")

	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	return r
}
