package controllers

import (
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/Pratham-Mishra04/interact/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetProject(c *fiber.Ctx) error {

	postID := c.Params("id")

	var project models.Project
	if err := initializers.DB.Preload("User").Select("id, username, name, profile_pic").First(&project, "id = ?", postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Project of this ID found."}
		}
		return &fiber.Error{Code: 500, Message: "Database Error."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"project": project,
	})
}

func AddProject(c *fiber.Ctx) error {
	var reqBody schemas.ProjectCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	if err := helpers.Validate[schemas.ProjectCreateSchema](reqBody); err != nil {
		return err
	}

	parsedID, err := uuid.Parse(c.GetRespHeader("loggedInUser"))
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Error Parsing the Loggedin User ID."}
	}

	picName, err := utils.SaveFile(c, "coverPic", "projects/coverPics", true, true, 900, 400)
	if err != nil {
		return err
	}
	reqBody.CoverPic = picName

	newProject := models.Project{
		UserID:      parsedID,
		Title:       reqBody.Title,
		Tagline:     reqBody.Tagline,
		CoverPic:    reqBody.CoverPic,
		Description: reqBody.Description,
		Tags:        reqBody.Tags,
		Category:    reqBody.Category,
		IsPrivate:   reqBody.IsPrivate,
		Links:       reqBody.Links,
	}

	result := initializers.DB.Create(&newProject)

	if result.Error != nil {
		return &fiber.Error{Code: 500, Message: "Internal Server Error while creating project"}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Project Added",
		"project": newProject,
	})
}
