package controllers

import (
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	API "github.com/Pratham-Mishra04/interact/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
)

func GetTrendingPosts(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	var posts []models.Post

	searchedDB := API.Search(c, 2)(paginatedDB)

	if err := searchedDB.
		Preload("User").
		Select("*, (2 * no_likes + no_comments + 5 * no_shares) AS weighted_average").
		Order("weighted_average DESC").
		Find(&posts).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Failed to get the Trending Posts."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"posts":  posts,
	})
}

func GetTrendingOpenings(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	var openings []models.Opening

	searchedDB := API.Search(c, 3)(paginatedDB)

	if err := searchedDB.Preload("Project").Order("created_at DESC").Find(&openings).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Failed to get the Trending Openings."}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"openings": openings,
	})
}

func GetTrendingProjects(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	var projects []models.Project

	searchedDB := API.Search(c, 1)(paginatedDB)

	if err := searchedDB.Order("no_shares DESC").Find(&projects).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Failed to get the Trending Projects."}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"projects": projects,
	})
}

func GetRecommendedProjects(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	var projects []models.Project

	searchedDB := API.Search(c, 1)(paginatedDB)

	if err := searchedDB.Find(&projects).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Failed to get the Trending Projects."}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"projects": projects,
	})
}

func GetMostLikedProjects(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	var projects []models.Project

	searchedDB := API.Search(c, 1)(paginatedDB)

	if err := searchedDB.Order("no_likes DESC").Find(&projects).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Failed to get the Trending Projects."}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"projects": projects,
	})
}

func GetRecentlyAddedProjects(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)
	var projects []models.Project

	searchedDB := API.Search(c, 1)(paginatedDB)

	if err := searchedDB.Order("created_at DESC").Find(&projects).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Failed to get the Trending Projects."}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"projects": projects,
	})
}

func GetLastViewedProjects(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var projectViewed []models.LastViewed

	paginatedDB := API.Paginator(c)(initializers.DB)
	if err := paginatedDB.Order("timestamp DESC").Preload("Project").Where("user_id=?", loggedInUserID).Find(&projectViewed).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Failed to get the Last Viewed Projects."}
	}

	var projects []models.Project

	for _, projectView := range projectViewed {
		projects = append(projects, projectView.Project)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"projects": projects,
	})
}

func GetTrendingUsers(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)

	searchedDB := API.Search(c, 0)(paginatedDB)

	var users []models.User
	if err := searchedDB.Order("no_followers DESC").Find(&users).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Failed to get the Trending Projects."}
	}
	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"users":  users,
	})
}

func GetRecommendedUsers(c *fiber.Ctx) error {
	paginatedDB := API.Paginator(c)(initializers.DB)

	searchedDB := API.Search(c, 0)(paginatedDB)

	var users []models.User
	if err := searchedDB.Order("created_at DESC").Find(&users).Error; err != nil {
		return &fiber.Error{Code: 500, Message: "Failed to get the Trending Projects."}
	}
	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"users":  users,
	})
}
