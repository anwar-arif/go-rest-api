package api

import (
	"encoding/json"
	"fmt"
	"go-rest-api/api/controller"
	"go-rest-api/config"
	"go-rest-api/infra/mongo"
	"io"
	"net/http"
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
func NewApiRouter(cfgDBTable *config.Table, db *mongo.Mongo, logger logger.StructLogger) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger(lgr))
	router.Use(middleware.Headers)
	router.Use(middleware.Cors())
	router.Use(chimiddleware.Timeout(30 * time.Second))

	router.NotFound(NotFoundHandler)
	router.MethodNotAllowed(MethodNotAllowed)

	router.Route("/", func(r chi.Router) {
		r.Mount("/brands", brandsRouter(controller.NewBrandsController(cfgDBTable, db, logger)))
		r.Mount("/users", usersRouter(controller.NewUsersController(cfgDBTable, db, logger)))
	})
	return router
}

func NewPingRouter(logger logger.StructLogger) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger(lgr))
	router.Use(middleware.Headers)
	router.Use(middleware.Cors())
	router.Use(chimiddleware.Timeout(30 * time.Second))

	router.NotFound(NotFoundHandler)
	router.MethodNotAllowed(MethodNotAllowed)

	router.Route("/", func(r chi.Router) {
		r.Mount("/", pingRouter(controller.NewPingController(logger)))
	})

	return router
}

// NewSystemRouter ...
func NewSystemRouter(sysCtrl *controller.SystemController) http.Handler {
	lgr.Println("NewSystemRouter")
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

func parseSlugFromUrlParameter(r *http.Request) (string, error) {
	slug := chi.URLParam(r, "slug")
	if len(slug) < 1 {
		return "", fmt.Errorf("slug is required")
	}

	return slug, nil
}

func getQueryParamString(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
