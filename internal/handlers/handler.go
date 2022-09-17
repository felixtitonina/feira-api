package handlers

import (
	"github.com/jeffersonto/feira-api/internal/service"
)

type Handler struct {
	Service service.FairService
}

func NewHandler(service service.FairService) Handler {
	return Handler{Service: service}
}