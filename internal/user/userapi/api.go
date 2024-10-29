package userapi

import (
	"context"

	"github.com/Abraxas-365/opd/internal/user"
	"github.com/Abraxas-365/opd/internal/user/usersrv"
	"github.com/Abraxas-365/toolkit/pkg/lucia"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up the API routes for the user service
func SetupRoutes(app *fiber.App, service *usersrv.Service, authMiddleware *lucia.AuthMiddleware[*user.User]) {
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

	app.Get("/users/blacklist", authMiddleware.RequireAuth(), func(c *fiber.Ctx) error {
		blacklist, err := service.GetBlacklist(context.TODO())
		if err != nil {
			return err
		}

		return c.JSON(blacklist)
	})

	app.Post("/users/blacklist", authMiddleware.RequireAuth(), func(c *fiber.Ctx) error {
		type Request struct {
			Email string `json:"email"`
		}

		var req Request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		err := service.AddToBlacklist(context.TODO(), req.Email)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)
	})

	app.Delete("/users/blacklist", authMiddleware.RequireAuth(), func(c *fiber.Ctx) error {
		type Request struct {
			Email string `json:"email"`
		}

		var req Request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		err := service.RemoveFromBlacklist(context.TODO(), req.Email)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)
	})
}
