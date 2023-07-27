package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/gofiber/fiber/v2"
)

func NotificationRouter(app *fiber.App) {
	notificationRoutes := app.Group("/notifications", middlewares.Protect)
	notificationRoutes.Get("/", controllers.GetNotifications)
	notificationRoutes.Delete("/:notificationID", controllers.DeleteNotification)
}
