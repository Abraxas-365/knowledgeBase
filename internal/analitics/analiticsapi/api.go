package analiticsapi

import (
	"time"

	"github.com/Abraxas-365/opd/internal/analitics/analiticssrv"
	"github.com/Abraxas-365/opd/internal/user"
	"github.com/Abraxas-365/toolkit/pkg/errors"
	"github.com/Abraxas-365/toolkit/pkg/lucia"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App,
	service *analiticssrv.Service,
	authMiddleware *lucia.AuthMiddleware[*user.User],
) {
	app.Get("/analytics", authMiddleware.RequireAuth(), getAnalytics(service))
}

func getAnalytics(service *analiticssrv.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse query parameters for date range
		startDateStr := c.Query("start_date")
		endDateStr := c.Query("end_date")

		var startDate, endDate *time.Time

		// Parse start date if provided
		if startDateStr != "" {
			parsedDate, err := time.Parse("2006-01-02", startDateStr)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid start_date format. Use YYYY-MM-DD",
				})
			}
			startDate = &parsedDate
		}

		// Parse end date if provided
		if endDateStr != "" {
			parsedDate, err := time.Parse("2006-01-02", endDateStr)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid end_date format. Use YYYY-MM-DD",
				})
			}
			endDate = &parsedDate
		}

		// Validate date range if both dates are provided
		if startDate != nil && endDate != nil && endDate.Before(*startDate) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "end_date cannot be before start_date",
			})
		}

		// Get analytics
		analytics, err := service.GetAllAnalitics(c.Context(), startDate, endDate)
		if err != nil {
			switch {
			case errors.IsNotFound(err):
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": err.Error(),
				})
			case errors.IsDatabaseError(err):
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Database error occurred",
				})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to fetch analytics",
				})
			}
		}

		return c.JSON(fiber.Map{
			"data": analytics,
		})
	}
}
