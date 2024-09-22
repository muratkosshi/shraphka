package routes

import (
	"github.com/gobuffalo/buffalo"
	"sharaphka/app/interfaces/http"
)

func SetRoutes(app *buffalo.App, authHandler *http.AuthHandler) {
	app.POST("/login", authHandler.Login)
}
