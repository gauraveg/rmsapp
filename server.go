package main

import (
	"net/http"

	"github.com/gauraveg/rmsapp/handlers"
	"github.com/gauraveg/rmsapp/middlewares"
	"github.com/gauraveg/rmsapp/utils"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	chi.Router
	server *http.Server
}

func RmsRouters() *Server {
	mainRouter := chi.NewRouter()
	mainRouter.Use(middlewares.CommonMiddleware()...)

	mainRouter.Route("/v1", func(v1 chi.Router) {
		v1.Get("/check", func(w http.ResponseWriter, r *http.Request) {
			utils.ResponseWithJson(w, http.StatusOK, map[string]string{
				"status": "ok",
			})
		})

		v1.Post("/login", handlers.UserLogin)
		v1.Group(func(router chi.Router) {
			router.Use(middlewares.Authenticate)
			router.Get("/info", func(w http.ResponseWriter, r *http.Request) {
				utils.ResponseWithJson(w, http.StatusOK, map[string]string{
					"Service":        "Restaurent Management System",
					"Users roles":    "Admin, Sub-admin, Users",
					"Fetch Data for": "Restaurents, Dishes, registered users",
				})
			})
			router.Post("/logout", handlers.UserLogout)

			router.Route("/admin", func(admin chi.Router) {
				admin.Use(middlewares.ShouldHaveRole("admin"))
				admin.Post("/create-sub-admin", handlers.CreateUser)
				admin.Get("/get-sub-admins", handlers.GetSubAdmins)
				admin.Post("/create-user", handlers.CreateUser)
				admin.Get("/get-users", handlers.GetUsersByAdmin)
				admin.Post("/create-restaurent", handlers.CreateRestaurent)
				admin.Get("/get-restaurents", handlers.GetRestaurentsByAdmin)
			})

			router.Route("/sub-admin", func(subadmin chi.Router) {
				subadmin.Use(middlewares.ShouldHaveRole("sub-admin"))
				subadmin.Post("/create-user", handlers.CreateUser)
				subadmin.Get("/get-users", handlers.GetUsersBySubAdmin)
				subadmin.Post("/create-restaurent", handlers.CreateRestaurent)
				subadmin.Get("/get-restaurents", handlers.GetRestaurentsBySubAdmin)
			})
		})
	})

	return &Server{
		Router: mainRouter,
	}
}

func (svc *Server) Run(port string) error {
	svc.server = &http.Server{
		Addr:    ":" + port,
		Handler: svc.Router,
	}
	return svc.server.ListenAndServe()
}
