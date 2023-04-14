package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func brandsRouter(ctrl BrandsController) http.Handler {
	h := chi.NewRouter()
	h.Group(func(r chi.Router) {
		r.Get("/", ctrl.ListBrand)
		r.Post("/", ctrl.AddBrand)
	})

	return h
}

func usersRouter(ctrl UsersController) http.Handler {
	h := chi.NewRouter()
	h.Group(func(r chi.Router) {
		r.Post("/signup", ctrl.CreateUser)
		r.Post("/login", ctrl.LogIn)
		r.Get("/", ctrl.GetByEmail)
	})
	return h
}

func pingRouter(ctrl *PingController) http.Handler {
	h := chi.NewRouter()
	h.Group(func(r chi.Router) {
		r.Get("/", ctrl.Ping)
	})
	return h
}

func healthRouter(ctrl *SystemController) http.Handler {
	log.Println("healthRouter")
	h := chi.NewRouter()
	h.Group(func(r chi.Router) {
		//r.Get("/api", ctrl.apiCheck)
		r.Get("/system", ctrl.systemCheck)
	})
	return h
}
