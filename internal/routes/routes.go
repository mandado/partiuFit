package routes

import (
	"partiuFit/internal/app"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterRoutes(app *app.Application) *chi.Mux {
	app.Logger.Info("Registering routes")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(app.Middlewares.ErrorHandlerMiddleware.Handle)

	r.Get("/health", app.Handlers.HeathCheck)

	r.Group(func(r chi.Router) {
		r.Use(app.Middlewares.UserMiddleware.Authenticate)
		r.Use(app.Middlewares.UserMiddleware.RequireUser)

		r.Route("/workouts", func(r chi.Router) {
			r.Get("/", app.Handlers.WorkoutHandlers.GetWorkouts)
			r.Post("/", app.Handlers.WorkoutHandlers.CreateWorkout)
			r.Get("/{id}", app.Handlers.WorkoutHandlers.GetWorkoutByID)
			r.Put("/{id}", app.Handlers.WorkoutHandlers.UpdateWorkout)
			r.Delete("/{id}", app.Handlers.WorkoutHandlers.DeleteWorkout)
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", app.Handlers.UserHandlers.RegisterUser)
		r.Put("/", app.Handlers.UserHandlers.UpdateUser)
	})

	r.Route("/tokens", func(r chi.Router) {
		r.Post("/", app.Handlers.TokensHandlers.CreateToken)
	})

	return r
}
