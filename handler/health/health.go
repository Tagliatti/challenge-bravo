package health

import (
	"net/http"
	"time"
)

type Health struct {
}

func NewHealthHandler() *Health {
	return &Health{}
}

func (h *Health) Handle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(time.Now().UTC().Format(time.RFC3339)))
}
