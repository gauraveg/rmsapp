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
		// todo add signup api
		v1.Post("/login", handlers.UserLogin)
		v1.Group(func(router chi.Router) {
			router.Use(middlewares.Authenticate)

			router.Route("/admin", func(admin chi.Router) {
				admin.Use(middlewares.ShouldHaveRole("admin"))
				admin.Post("/create-sub-admin", handlers.CreateUser)
				admin.Get("/get-sub-admins", handlers.GetSubAdminsByAdmin)
				admin.Post("/create-user", handlers.CreateUser)
				admin.Get("/get-users", handlers.GetUsersByAdmin)
				admin.Post("/create-restaurant", handlers.CreateRestaurant)
				admin.Get("/get-restaurants", handlers.GetRestaurantsByAdmin)
				admin.Route("/{restaurantId}", func(restId chi.Router) {
					restId.Post("/create-dish", handlers.CreateDish)
				})
				admin.Get("/get-all-dishes", handlers.GetAllDishesByAdmin)
			})

			router.Route("/sub-admin", func(subAdmin chi.Router) {
				subAdmin.Use(middlewares.ShouldHaveRole("sub-admin"))
				subAdmin.Post("/create-user", handlers.CreateUser)
				subAdmin.Get("/get-users", handlers.GetUsersBySubAdmin)
				subAdmin.Post("/create-restaurant", handlers.CreateRestaurant)
				subAdmin.Get("/get-restaurants", handlers.GetRestaurantsBySubAdmin)
				subAdmin.Route("/{restaurantId}", func(restId chi.Router) {
					restId.Post("/create-dish", handlers.CreateDish)
				})
				subAdmin.Get("/get-all-dishes", handlers.GetAllDishesBySubAdmin)
			})

			router.Route("/user", func(user chi.Router) {
				user.Use(middlewares.ShouldHaveRole("user"))
				user.Get("/get-all-restaurants", handlers.GetRestaurantsByAdmin)
				user.Get("/get-all-dishes", handlers.GetAllDishesByAdmin)
				user.Route("/{restaurantId}", func(restId chi.Router) {
					restId.Get("/dishes-by-restaurant", handlers.GetDishesByRestId)
				})
			})

			router.Post("/logout", handlers.UserLogout)
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
