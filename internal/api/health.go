package api

import (
	"net/http"
	"github.com/danmaina/HttpResponse"
)

type healthResponse struct {
	Status string `json:"status"`
}

// HealthHandler returns a simple 200 OK for k8s health checks
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	res := handlers.Response{
		Status: http.StatusOK,
		Error:  nil,
		Body:   healthResponse{Status: "UP"},
	}
	res.ReturnResponse(w)
}
