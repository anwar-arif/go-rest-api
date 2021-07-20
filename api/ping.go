package api

import (
	"go-rest-api/api/response"
	"go-rest-api/logger"
	"go-rest-api/model"
	"go-rest-api/service"
	"net/http"
)

type PingController struct {
	svc *service.Service
	lgr logger.StructLogger
}

func NewPingController(svc *service.Service, lgr logger.StructLogger) *PingController {
	return &PingController{
		svc: svc,
		lgr: lgr,
	}
}

func (pingCtrl *PingController) Ping(w http.ResponseWriter, r *http.Request) {
	ping := model.PingData{
		Method:      r.Method,
		ServiceName: "Go rest api",
	}
	_ = response.ServeJSON(w, http.StatusOK, nil, nil, "success", ping)
	return
}
