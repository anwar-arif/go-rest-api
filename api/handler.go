package api

import (
	"go-rest-api/api/controller"
	"go-rest-api/api/middleware"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func brandsRouter(ctrl controller.BrandsController) http.Handler {
	h := chi.NewRouter()
	h.Group(func(r chi.Router) {
		r.Get("/", ctrl.ListBrand)
		r.Post("/", ctrl.AddBrand)
	})

	return h
}

func usersRouter(ctrl controller.UsersController) http.Handler {
	h := chi.NewRouter()
	h.Group(func(r chi.Router) {
		r.Post("/signup", ctrl.CreateUser)
		r.Post("/login", ctrl.LogIn)
	})
	h.Group(func(r chi.Router) {
		r.Use(middleware.AuthenticatedOnly)
		r.Get("/users", ctrl.GetByEmail)
		r.Post("/logout", ctrl.LogOut)
	})
	return h
}

func pingRouter(ctrl *controller.PingController) http.Handler {
	h := chi.NewRouter()
	h.Group(func(r chi.Router) {
		r.Get("/", ctrl.Ping)
	})
	return h
}

func healthRouter(ctrl *controller.SystemController) http.Handler {
	log.Println("healthRouter")
	h := chi.NewRouter()
	h.Group(func(r chi.Router) {
		r.Get("/system", ctrl.SystemCheck)
	})
	return h
}
