package main

import (
	"context"
	"fmt"

	"github.com/Abraxas-365/opd/internal/analitics/analiticsapi"
	analyticsinfra "github.com/Abraxas-365/opd/internal/analitics/analiticsinfra"
	"github.com/Abraxas-365/opd/internal/analitics/analiticssrv"
	"github.com/Abraxas-365/opd/internal/chatuser/chatuserinfra"
	"github.com/Abraxas-365/opd/internal/chatuser/chatusersrv"
	"github.com/Abraxas-365/opd/internal/interaction/interactioninfra"
	"github.com/Abraxas-365/opd/internal/interaction/interactionsrv"
	"github.com/Abraxas-365/opd/internal/kb/kbapi"
	"github.com/Abraxas-365/opd/internal/kb/kbasesrv"
	"github.com/Abraxas-365/opd/internal/kb/kbinfra"
	"github.com/Abraxas-365/opd/internal/user"
	"github.com/Abraxas-365/opd/internal/user/userapi"
	"github.com/Abraxas-365/opd/internal/user/userinfra"
	"github.com/Abraxas-365/opd/internal/user/usersrv"
	"github.com/Abraxas-365/opd/pkg/conf"
	"github.com/Abraxas-365/toolkit/pkg/errors"
	"github.com/Abraxas-365/toolkit/pkg/lucia"
	"github.com/Abraxas-365/toolkit/pkg/lucia/luciastore"
	"github.com/Abraxas-365/toolkit/pkg/s3client"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/bedrockagent"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jmoiron/sqlx"
)

func main() {

	conf := conf.Load()

	db, err := sqlx.Connect("postgres", conf.DatabaseURL)
	if err != nil {
		panic(err)
	}

	userRepo := userinfra.NewUserStore(db)
	userSrv := usersrv.NewService(userRepo)
	analrepo := analyticsinfra.NewAnalyticsStore(db)
	sessionStore := luciastore.NewStoreFromConnection(db)
	authSrv := lucia.NewAuthService[*user.User](userSrv, sessionStore)
	chatUserRepo := chatuserinfra.NewChatUserStore(db)
	chatUserSrv := chatusersrv.New(chatUserRepo)
	analSrv := analiticssrv.NewService(analrepo)

	interactionRepo := interactioninfra.NewInteractionStore(db)
	interactionSrv := interactionsrv.New(interactionRepo)

	// Initialize Google OAuth provider
	googleProvider := lucia.NewGoogleProvider(
		conf.GoogleClientID,
		conf.GoogleClientSecret,
		conf.GoogleRedirectURI,
	)
	authSrv.RegisterProvider("google", googleProvider)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		panic("unable to load SDK config: " + err.Error())
	}

	client := bedrockagentruntime.NewFromConfig(cfg)

	repo := kbinfra.NewStore(db)
	s3client, err := s3client.NewS3Client("vendy", s3client.WithRegion("us-east-1"))
	if err != nil {
		panic(err)
	}

	brClient := bedrockagent.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})))

	// Then modify the kbService initialization to include the brClient:
	kbSerive := kbsrv.New(client, brClient, repo, s3client, *userSrv, *chatUserSrv, *interactionSrv)

	app := fiber.New()
	authMiddleware := lucia.NewAuthMiddleware(authSrv)
	app.Use(authMiddleware.SessionMiddleware())

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     conf.AllowOrigins,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	kbapi.SetupRoutes(app, kbSerive, authMiddleware)
	userapi.SetupRoutes(app, userSrv, authMiddleware)
	analiticsapi.SetupRoutes(app, analSrv, authMiddleware)

	// Google OAuth routes
	app.Get("/login/google", func(c *fiber.Ctx) error {
		authURL, state, err := authSrv.GetAuthURL("google")
		if err != nil {
			return err
		}
		c.Cookie(&fiber.Cookie{
			Name:     "oauth_state",
			Value:    state,
			HTTPOnly: true,
			Secure:   true,
		})
		return c.Redirect(authURL)
	})

	app.Get("/login/google/callback", func(c *fiber.Ctx) error {
		state := c.Cookies("oauth_state")
		if state == "" || state != c.Query("state") {
			return errors.ErrUnauthorized("Invalid state")
		}

		code := c.Query("code")
		if code == "" {
			return errors.ErrBadRequest("Missing code")
		}

		session, err := authSrv.HandleCallback(c.Context(), "google", code)
		if err != nil {
			return err
		}

		// Set session cookie
		lucia.SetSessionCookie(c, session)

		// Return session ID in JSON response for the frontend to access
		res := c.JSON(fiber.Map{
			"session_id": session.ID, // Assuming session has an ID field
			"message":    "Login successful",
		})

		fmt.Println(res)

		return c.Redirect(conf.RedirectAfterLogin)
	})
	// Logout route
	app.Post("/logout", func(c *fiber.Ctx) error {
		session := lucia.GetSession(c)
		if session != nil {
			if err := authSrv.DeleteSession(c.Context(), session.ID); err != nil {
				return err
			}
		}
		lucia.ClearSessionCookie(c)
		return c.SendString("Logged out successfully")
	})

	// Start server
	app.Listen(conf.Port)
}
