package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-rest-api/service"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"go-rest-api/api/middleware"
	"go-rest-api/logger"
)

var lgr logger.Logger

// SetLogger ..
func SetLogger(l logger.Logger) {
	lgr = l
}

// NewApiRouter ..
func NewApiRouter(svc *service.Service, logger logger.StructLogger) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger(lgr))
	router.Use(middleware.Headers)
	router.Use(middleware.Cors())
	router.Use(chimiddleware.Timeout(30 * time.Second))

	router.NotFound(NotFoundHandler)
	router.MethodNotAllowed(MethodNotAllowed)

	router.Route("/", func(r chi.Router) {
		r.Mount("/brands", brandsRouter(NewBrandsController(svc, logger)))
	})
	return router
}

func NewPingRouter(svc *service.Service, logger logger.StructLogger) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger(lgr))
	router.Use(middleware.Headers)
	router.Use(middleware.Cors())
	router.Use(chimiddleware.Timeout(30 * time.Second))

	router.NotFound(NotFoundHandler)
	router.MethodNotAllowed(MethodNotAllowed)

	router.Route("/", func(r chi.Router) {
		r.Mount("/", pingRouter(NewPingController(svc, logger)))
	})

	return router
}

// NewSystemRouter ...
func NewSystemRouter(sysCtrl *SystemController) http.Handler {
	log.Println("NewSystemRouter")
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger(lgr))
	router.Use(middleware.Headers)
	router.Use(middleware.Cors())
	router.Use(chimiddleware.Timeout(30 * time.Second))
	router.Route("/", func(r chi.Router) {
		r.Mount("/health", healthRouter(sysCtrl))
	})
	return router
}

// NotFoundHandler handles when no routes match
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// MethodNotAllowed handles when no routes match
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func parseJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func parseSkipLimit(r *http.Request, def, max int) (int64, int64, int64, error) {
	q := r.URL.Query()
	var page, skip, limit int64 = 1, 0, int64(def)
	pageQ := q.Get("page")
	if pageQ != "" {
		pageInt, err := strconv.Atoi(pageQ)
		if err != nil {
			return page, skip, limit, errors.New("failed to parse page")
		}
		page = int64(pageInt)
	}
	limitQ := q.Get("limit")

	if limitQ != "" {
		limitInt, err := strconv.Atoi(limitQ)
		if err != nil {
			return page, skip, limit, errors.New("failed to parse limit")
		}
		limit = int64(limitInt)
	}

	if limit > int64(max) {
		limit = int64(max)
	}
	skip = limit * (page - 1)
	if skip < 0 {
		skip = 0
	}
	return page, skip, limit, nil
}

func parseSlugFromUrlParameter(r *http.Request) (string, error) {
	slug := chi.URLParam(r, "slug")
	if len(slug) < 1 {
		return "", fmt.Errorf("slug is required")
	}

	return slug, nil
}

func getNextPreviousPager(path string, page, limit int64) (*string, *string) {
	var previous, next string
	if page-1 > 0 {
		previous = fmt.Sprintf("%s?limit=%d&page=%d", path, limit, page-1)
	}
	next = fmt.Sprintf("%s?limit=%d&page=%d", path, limit, page+1)
	return &previous, &next
}

func getQueryParamString(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
