package controllers

import (
	"errors"
	"log"
	"reflect"
	"time"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/routines"
	"github.com/Pratham-Mishra04/interact/schemas"
	"github.com/Pratham-Mishra04/interact/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetViews(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(userID)

	viewsArr, count, err := utils.GetProfileViews(parsedUserID)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"viewsArr": viewsArr,
		"count":    count,
	})
}

func GetMe(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")

	var user models.User
	initializers.DB.
		Preload("Achievements").
		Preload("Projects").
		Preload("Posts").
		Preload("Posts.User").
		Preload("Memberships").
		Preload("Memberships.Project").
		First(&user, "id = ?", userID)

	var projects []models.Project
	for _, project := range user.Projects {
		if !project.IsPrivate {
			projects = append(projects, project)
		}
	}
	user.Projects = projects

	var memberships []models.Membership
	for _, membership := range user.Memberships {
		if !membership.Project.IsPrivate {
			memberships = append(memberships, membership)
		}
	}
	user.Memberships = memberships

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"user":    user,
	})
}

func GetUser(c *fiber.Ctx) error {
	username := c.Params("username")

	var user models.User
	initializers.DB.
		Preload("Achievements").
		Preload("Projects").
		Preload("Posts").
		Preload("Posts.User").
		Preload("Memberships").
		Preload("Memberships.Project").
		First(&user, "username = ?", username)

	if user.ID == uuid.Nil {
		return &fiber.Error{Code: 400, Message: "No user of this username found."}
	}

	var projects []models.Project
	for _, project := range user.Projects {
		if !project.IsPrivate {
			projects = append(projects, project)
		}
	}
	user.Projects = projects

	var memberships []models.Membership
	for _, membership := range user.Memberships {
		if !membership.Project.IsPrivate {
			memberships = append(memberships, membership)
		}
	}
	user.Memberships = memberships

	loggedInUserID := c.GetRespHeader("loggedInUserID")

	if user.ID.String() != loggedInUserID {
		routines.UpdateProfileViews(&user)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User Found",
		"user":    user,
	})
}

func GetMyLikes(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var likes []models.Like
	initializers.DB.
		Find(&likes, "user_id = ?", loggedInUserID)

	var likeIDs []string
	for _, like := range likes {
		if like.PostID != nil {
			likeIDs = append(likeIDs, like.PostID.String())
		} else if like.ProjectID != nil {
			likeIDs = append(likeIDs, like.ProjectID.String())
		} else if like.CommentID != nil {
			likeIDs = append(likeIDs, like.CommentID.String())
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User Found",
		"likes":   likeIDs,
	})
}

func UpdateMe(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")
	var user models.User
	if err := initializers.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No user of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	var reqBody schemas.UserUpdateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Request Body."}
	}

	picName, err := utils.SaveFile(c, "profilePic", "user/profilePics", true, 500, 500)
	if err != nil {
		return err
	}
	reqBody.ProfilePic = picName

	coverName, err := utils.SaveFile(c, "coverPic", "user/coverPics", true, 900, 400)
	if err != nil {
		return err
	}
	reqBody.CoverPic = coverName

	//TODO  make a routine for this
	if reqBody.ProfilePic != "" {
		err := utils.DeleteFile("user/profilePics", user.ProfilePic)
		if err != nil {
			log.Printf("Error while deleting user profile pic: %e", err)
		}
	}
	if reqBody.CoverPic != "" {
		err := utils.DeleteFile("user/coverPics", user.CoverPic)
		if err != nil {
			log.Printf("Error while deleting user cover pic: %e", err)
		}
	}

	updateUserValue := reflect.ValueOf(&reqBody).Elem()
	userValue := reflect.ValueOf(&user).Elem()

	for i := 0; i < updateUserValue.NumField(); i++ {
		field := updateUserValue.Type().Field(i)
		fieldName := field.Name

		if fieldValue := updateUserValue.Field(i); fieldValue.IsValid() && fieldValue.String() != "" {
			userField := userValue.FieldByName(fieldName)
			if userField.IsValid() && userField.CanSet() {
				userField.Set(fieldValue)
			}
		}
	}

	if err := initializers.DB.Save(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User updated successfully",
		"user":    user,
	})
}

func DeactivateMe(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")

	//TODO send email for verification

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No user of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	user.Active = false

	if err := initializers.DB.Save(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "User deactivated successfully",
	})
}

func UpdatePassword(c *fiber.Ctx) error {
	var reqBody struct {
		Password        string `json:"password"`
		NewPassword     string `json:"newPassword"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Validation Failed"}
	}

	if reqBody.NewPassword != reqBody.ConfirmPassword {
		return &fiber.Error{Code: 400, Message: "Passwords do not match."}
	}

	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var user models.User
	initializers.DB.First(&user, "id = ?", loggedInUserID)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)); err != nil {
		return &fiber.Error{Code: 400, Message: "Incorrect Password."}
	}

	//TODO send email for verification

	hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.NewPassword), 10)

	if err != nil {
		return helpers.AppError{Code: 500, Message: config.SERVER_ERROR, Err: err}
	}

	user.Password = string(hash)
	user.PasswordChangedAt = time.Now()

	if err := initializers.DB.Save(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return CreateSendToken(c, user, 200, "Password updated successfully")
}

func UpdateEmail(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")

	var reqBody struct {
		Email string `json:"email" validate:"required,email"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Validation Failed"}
	}

	var emailCheckUser models.User
	if err := initializers.DB.First(&emailCheckUser, "email = ?", reqBody.Email).Error; err == nil {
		return &fiber.Error{Code: 400, Message: "Email Address Already In Use."}
	}

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No user of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	user.Email = reqBody.Email
	user.Verified = false

	if err := initializers.DB.Save(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User updated successfully",
	})
}

func UpdatePhoneNo(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")

	var reqBody struct {
		PhoneNo string `json:"phoneNo"  validate:"e164"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Validation Failed"}
	}

	var phoneNoCheckUser models.User
	if err := initializers.DB.First(&phoneNoCheckUser, "phone_no = ?", reqBody.PhoneNo).Error; err == nil {
		return &fiber.Error{Code: 400, Message: "Phone Number Already In Use."}
	}

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No user of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	user.PhoneNo = reqBody.PhoneNo

	if err := initializers.DB.Save(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User updated successfully",
	})
}

func Deactive(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &fiber.Error{Code: 400, Message: "No user of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	user.Active = false
	user.DeactivatedAt = time.Now()

	if err := initializers.DB.Save(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, Err: err}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Account Deactived",
	})
}
