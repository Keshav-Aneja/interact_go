package middlewares

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func checkOrgAccess(UserRole models.OrganizationRole, AuthorizedRole models.OrganizationRole) bool {
	if UserRole == models.Owner {
		return true
	} else if UserRole == models.Manager {
		return AuthorizedRole != models.Owner
	} else if UserRole == models.Member {
		return AuthorizedRole == models.Member
	}

	return false
}

func checkProjectAccess(UserRole models.ProjectRole, AuthorizedRole models.ProjectRole) bool {
	if UserRole == models.ProjectManager {
		return true
	} else if UserRole == models.ProjectEditor {
		return AuthorizedRole != models.ProjectManager
	} else if UserRole == models.ProjectMember {
		return AuthorizedRole == models.ProjectMember
	}

	return false
}

func OrgRoleAuthorization(Role models.OrganizationRole) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		loggedInUserID := c.GetRespHeader("loggedInUserID")
		orgID := c.Params("orgID")

		var orgMembership models.OrganizationMembership
		if err := initializers.DB.Preload("Organization").First(orgMembership, "organization_id = ? AND user_id=?", orgID, loggedInUserID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				var org models.Organization
				if err := initializers.DB.First(org, "user_id=?", loggedInUserID).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						return &fiber.Error{Code: 403, Message: "Cannot access this organization"}
					}
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
				}
				return c.Next()
			}
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}

		if !checkOrgAccess(orgMembership.Role, Role) {
			return &fiber.Error{Code: 403, Message: "You don't have the Permission to perform this action."}
		}

		c.Set("orgMemberID", c.GetRespHeader("loggedInUserID"))
		c.Set("loggedInUserID", orgMembership.Organization.UserID.String())

		return c.Next()
	}
}

func ProjectRoleAuthorization(Role models.ProjectRole) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		loggedInUserID := c.GetRespHeader("loggedInUserID")
		projectID := c.Params("projectID")

		var membership models.Membership
		if err := initializers.DB.Preload("Project").First(membership, "projectID = ? AND user_id=?", projectID, loggedInUserID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				var project models.Project
				if err := initializers.DB.First(project, "user_id=?", loggedInUserID).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						return &fiber.Error{Code: 403, Message: "Cannot access this project"}
					}
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
				}
				return c.Next()
			}
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
		}

		if !checkProjectAccess(membership.Role, Role) {
			return &fiber.Error{Code: 403, Message: "You don't have the Permission to perform this action."}
		}

		c.Set("projectMemberID", c.GetRespHeader("loggedInUserID"))
		c.Set("loggedInUserID", membership.Project.UserID.String())

		return c.Next()
	}
}
