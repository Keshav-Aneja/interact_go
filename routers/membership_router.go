package routers

import (
	"github.com/Pratham-Mishra04/interact/controllers"
	"github.com/Pratham-Mishra04/interact/middlewares"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
)

func MembershipRouter(app *fiber.App) {
	membershipRoutes := app.Group("/membership", middlewares.Protect)
	membershipRoutes.Post("/project/:projectID", middlewares.ProjectRoleAuthorization(models.ProjectManager), controllers.AddMember)
	membershipRoutes.Patch("/:membershipID", middlewares.ProjectRoleAuthorization(models.ProjectManager), controllers.ChangeMemberRole)
	membershipRoutes.Delete("project/:projectID", controllers.LeaveProject)
	membershipRoutes.Delete("/:membershipID", middlewares.ProjectRoleAuthorization(models.ProjectManager), controllers.RemoveMember)
}
