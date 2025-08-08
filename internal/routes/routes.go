package routes

import (
	"partiuFit/internal/app"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func RegisterRoutes(app *app.Application) *chi.Mux {
	app.Logger.Info("Registering routes")

	r := chi.NewRouter()

	// Security middlewares
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rate limiting: 100 requests per minute per IP
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	// Security headers
	r.Use(app.Middlewares.SecurityMiddleware.SecurityHeaders)

	// Standard middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
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
