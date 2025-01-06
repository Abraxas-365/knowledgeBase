package userapi

import (
	"context"
	"fmt"
	"log"

	"github.com/Abraxas-365/opd/internal/user"
	"github.com/Abraxas-365/opd/internal/user/usersrv"
	"github.com/Abraxas-365/toolkit/pkg/lucia"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up the API routes for the user service
func SetupRoutes(app *fiber.App, service *usersrv.Service, authMiddleware *lucia.AuthMiddleware[*user.User]) {

	app.Get("/users/me", authMiddleware.RequireAuth(), func(c *fiber.Ctx) error {
		log.Println("Accessing /users/me endpoint")
		log.Printf("Request headers: %+v", c.GetReqHeaders())

		session := lucia.GetSession(c)
		userID, err := session.UserIDToString()
		log.Printf("Authenticated user ID: %s", session.ID)
		user, err := service.GetUser(c.Context(), userID)
		if err != nil {
			log.Printf("Error fetching user details: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to fetch user details: %v", err)})
		}

		log.Println("Successfully fetched user details")
		log.Printf("User details: %+v", user)

		return c.JSON(user)
	})

	app.Get("/users", authMiddleware.RequireAuth(), func(c *fiber.Ctx) error {
		page, pageSize := 1, 10 // Default values
		pageParam := c.Query("page")
		pageSizeParam := c.Query("pageSize")

		// Parse pagination params if provided
		if pageParam != "" {
			page = c.QueryInt("page")
		}
		if pageSizeParam != "" {
			pageSize = c.QueryInt("pageSize")
		}

		users, err := service.GetUsers(context.TODO(), page, pageSize)
		if err != nil {
			return err
		}

		return c.JSON(users)
	})

	app.Get("/users/not-admin", authMiddleware.RequireAuth(), func(c *fiber.Ctx) error {
		page, pageSize := 1, 10
		if p := c.Query("page"); p != "" {
			page = c.QueryInt("page")
		}
		if ps := c.Query("pageSize"); ps != "" {
			pageSize = c.QueryInt("pageSize")
		}

		users, err := service.GetNotAdminUsers(context.TODO(), page, pageSize)
		if err != nil {
			return err
		}

		return c.JSON(users)
	})

	app.Get("/users/admin", authMiddleware.RequireAuth(), func(c *fiber.Ctx) error {
		page, pageSize := 1, 10
		if p := c.Query("page"); p != "" {
			page = c.QueryInt("page")
		}
		if ps := c.Query("pageSize"); ps != "" {
			pageSize = c.QueryInt("pageSize")
		}

		users, err := service.GetUsersAdminRole(context.TODO(), page, pageSize)
		if err != nil {
			return err
		}

		return c.JSON(users)
	})

	app.Post("/users/promote-to-admin", authMiddleware.RequireAuth(), func(c *fiber.Ctx) error {
		type Request struct {
			UserID string `json:"userID"`
		}

		var req Request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		err := service.PromoteUserToAdmin(context.TODO(), req.UserID)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)
	})

	app.Delete("/users/:id", authMiddleware.RequireAuth(), func(c *fiber.Ctx) error {
		userID := c.Params("id")

		err := service.DeleteUser(context.TODO(), userID)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	})

	app.Get("/users/whitelist", authMiddleware.RequireAuth(), func(c *fiber.Ctx) error {
		blacklist, err := service.GetWhitelist(context.TODO())
		if err != nil {
			return err
		}

		return c.JSON(blacklist)
	})

	app.Post("/users/whitelist", authMiddleware.RequireAuth(), func(c *fiber.Ctx) error {
		type Request struct {
			Email string `json:"email"`
		}

		var req Request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		err := service.AddToWhitelist(context.TODO(), req.Email)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)
	})

	app.Delete("/users/d/whitelist", authMiddleware.RequireAuth(), func(c *fiber.Ctx) error {
		type Request struct {
			Email string `json:"email"`
		}

		var req Request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		err := service.RemoveFromWhitelist(context.TODO(), req.Email)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)
	})
}
