package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go-rest-api/api/response"
	"go-rest-api/infra"
)

// SystemController ..
type SystemController struct {
	db infra.DB
}

// NewSystemController ..
func NewSystemController(db infra.DB) *SystemController {
	return &SystemController{
		db: db,
	}
}

func (s *SystemController) SystemCheck(w http.ResponseWriter, r *http.Request) {
	if err := s.ConnCheck(); err != nil {
		_ = response.Serve(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.ServeJSONData(w, "ok", http.StatusOK)
	return
}

func (s *SystemController) ApiCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("ApiCheck")
	if err := s.ConnCheck(); err != nil {
		_ = response.Serve(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.ServeJSONData(w, "ok", http.StatusOK)
	return
}

func (s *SystemController) ConnCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	log.Println("db ping")
	if err := s.db.Ping(ctx); err != nil {
		return fmt.Errorf("mongo conn error: %v", err)
	}

	return nil
}
