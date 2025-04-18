package app

import (
	"taans/internal/handler"

	"github.com/go-chi/chi/v5"
)

func (app *Application) RegisterRoutes() {
	app.Router.Route("/", func(r chi.Router) {
		r.Post("/", handler.HandleNewMessage(app.Message))
	})
}
