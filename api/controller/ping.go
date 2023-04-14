package controller

import (
	"go-rest-api/api/response"
	"go-rest-api/logger"
	"go-rest-api/model"
	"net/http"
)

type PingController struct {
	lgr logger.StructLogger
}

func NewPingController(lgr logger.StructLogger) *PingController {
	return &PingController{
		lgr: lgr,
	}
}

func (pingCtrl *PingController) Ping(w http.ResponseWriter, r *http.Request) {
	ping := model.PingData{
		Method:      r.Method,
		ServiceName: "Go rest api",
	}
	_ = response.Serve(w, http.StatusOK, "success", ping)
	return
}
