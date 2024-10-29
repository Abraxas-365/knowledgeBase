package kbapi

import (
	"context"
	"time"

	kbsrv "github.com/Abraxas-365/opd/internal/kb/kbasesrv"
	"github.com/Abraxas-365/opd/internal/user"
	"github.com/Abraxas-365/toolkit/pkg/lucia"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// SetupRoutes sets up the API routes for the knowledge base service
func SetupRoutes(app *fiber.App, service *kbsrv.Service, authMiddleware *lucia.AuthMiddleware[*user.User]) {
	// Define the route for completing answers with metadata
	limiterGroup := app.Group("/chat", limiter.New(limiter.Config{
		Max:        100,           // Maximum number of requests per IP
		Expiration: 8 * time.Hour, // Time window for rate limiting
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Limit by IP
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).SendString("Rate limit exceeded")
		},
	}))
	limiterGroup.Post("/complete-answer", func(c *fiber.Ctx) error {
		type Request struct {
			UserMessage string  `json:"userMessage"`
			SessionID   *string `json:"sessionID,omitempty"`
		}

		var req Request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		output, err := service.CompleteAnswerWithMetadata(context.TODO(), req.UserMessage, req.SessionID)
		if err != nil {
			return err
		}

		return c.JSON(output)
	})

	// Route to generate a presigned PUT URL
	app.Post("/generate-presigned-url", authMiddleware.RequireAuth(), func(c *fiber.Ctx) error {
		type Request struct {
			Key string `json:"key"`
		}

		var req Request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		url, err := service.GeneratePutURL(req.Key)
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{"url": url})
	})

	// Route to list objects with pagination
	app.Get("/list-objects", authMiddleware.RequireAuth(), func(c *fiber.Ctx) error {
		type Request struct {
			PageSize          int32  `query:"pageSize"`
			ContinuationToken string `query:"continuationToken"`
		}

		var req Request
		if err := c.QueryParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
		}

		var continuationToken *string
		if req.ContinuationToken != "" {
			continuationToken = &req.ContinuationToken
		}

		files, nextToken, err := service.LisObjects(req.PageSize, continuationToken)
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"files":             files,
			"continuationToken": nextToken,
		})
	})

	// Route to delete an object
	app.Delete("/delete-object", authMiddleware.RequireAuth(), func(c *fiber.Ctx) error {
		type Request struct {
			Key string `json:"key"`
		}

		var req Request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if err := service.DeleteObject(req.Key); err != nil {
			return err
		}

		return c.JSON(fiber.Map{"message": "Object deleted successfully"})
	})

	// Endpoint to start the ingestion job for syncing knowledge base
	app.Post("/sync-knowledge-base", authMiddleware.RequireAuth(), func(c *fiber.Ctx) error {
		output, err := service.SyncKnowledgeBase(context.TODO())
		if err != nil {
			return err
		}
		return c.JSON(output)
	})

	// TODO: Implement the logic to get the status of the ingestion job
	// Endpoint to get the status of the ingestion job
	//
	//	app.Get("/ingestion-job-status", authMiddleware.RequireAuth(), func(c *fiber.Ctx) error {
	//		type Request struct {
	//			IngestionJobId string `query:"ingestionJobId"`
	//		}
	//
	//		var req Request
	//		if err := c.QueryParser(&req); err != nil {
	//			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
	//		}
	//
	//		// Implement call to check the ingestion job status using service.GetIngestionJobStatus (to be created in Service)
	//		status, err := service.GetIngestionJobStatus(req.IngestionJobId)
	//		if err != nil {
	//			return err
	//		}
	//
	//		return c.JSON(fiber.Map{
	//			"status": status,
	//		})
	//	})
}
