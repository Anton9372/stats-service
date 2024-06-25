package metric

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"stats-service/pkg/logging"
)

const (
	URL = "/api/heartbeat"
)

type Handler struct {
	Logger *logging.Logger
}

func NewHandler(logger *logging.Logger) *Handler {
	return &Handler{
		Logger: logger,
	}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, URL, h.Heartbeat)
}

func (h *Handler) Heartbeat(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(204)
}
