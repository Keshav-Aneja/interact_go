package organization_routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/validators"
	"github.com/gofiber/fiber/v2"
)

func AuthRouter(app *fiber.App) {
	oauthRoutes := app.Group("/org")
	oauthRoutes.Post("/signup", validators.UserCreateValidator, controllers.OrganizationSignUp)
}
