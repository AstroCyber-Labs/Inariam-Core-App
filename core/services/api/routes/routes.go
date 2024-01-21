package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"

	"gitea/pcp-inariam/inariam/core/services/api"
	"gitea/pcp-inariam/inariam/core/services/api/handlers"
	"gitea/pcp-inariam/inariam/core/services/api/handlers/cloud/aws"
	"gitea/pcp-inariam/inariam/core/services/api/handlers/cloud/gcp"
	"gitea/pcp-inariam/inariam/core/services/api/responses"
	"gitea/pcp-inariam/inariam/pkgs/log"
)

// ConfigureRoutes this function configure all the routes present in the http api service
func ConfigureRoutes(httpApi *api.API) {

	authHandler := handlers.NewAuthHandler(httpApi)
	awsHandler := aws.NewAwsHandler(httpApi)
	gcpHandler := gcp.NewGCPHandler(httpApi)

	httpApi.Echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)
			return nil
		},
	}))

	httpApi.Echo.Use(middleware.CORS())

	httpApi.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
	httpApi.Echo.GET("/", func(c echo.Context) error {
		return responses.MessageResponse(c, http.StatusOK, "Hello there")
	})
	authGroup := httpApi.Echo.Group("/auth")

	httpApi.Echo.GET("/health", handlers.HealthCheck)
	authGroup.POST("/login", authHandler.StartSignInProcess)
	authGroup.POST("/verify-user", authHandler.VerifyUser)
	authGroup.POST("/confirm-mfa", authHandler.ConfirmMFACode)
	authGroup.POST("/complete-signin", authHandler.CompleteSignIn)
	authGroup.POST("/confirmation-resend", authHandler.ResendConfirmatioNEmail)
	authGroup.POST("/activate-mfa", authHandler.GetMFADeviceCode)

	awsIam := httpApi.Echo.Group("/aws/iam")

	awsIamGroup := awsIam.Group("/groups")
	awsIamGroup.GET("/", awsHandler.ListGroups)
	awsIamGroup.GET("/:id", awsHandler.GetGroup)
	awsIamGroup.POST("/", awsHandler.CreateGroup)
	awsIamGroup.PUT("/:id", awsHandler.UpdateGroup)
	awsIamGroup.DELETE("/:id", awsHandler.DeleteGroup)

	awsIamUser := awsIam.Group("/users")
	awsIamUser.GET("/", awsHandler.ListUsers)
	awsIamUser.GET("/:id", awsHandler.GetUser)
	awsIamUser.POST("/", awsHandler.CreateUser)
	awsIamUser.PUT("/:id", awsHandler.UpdateUser)
	awsIamUser.DELETE("/:id", awsHandler.DeleteUser)

	awsIamRole := awsIam.Group("/roles")
	awsIamRole.GET("/", awsHandler.ListRoles)
	awsIamRole.GET("/:id", awsHandler.GetRole)
	awsIamRole.POST("/", awsHandler.CreateRole)
	awsIamRole.PUT("/:id", awsHandler.UpdateRole)
	awsIamRole.DELETE("/:id", awsHandler.DeleteRole)

	awsIamPolicy := awsIam.Group("/policies")
	awsIamPolicy.GET("/", awsHandler.ListPolicies)
	awsIamPolicy.GET("/:arn", awsHandler.GetPolicy)
	awsIamPolicy.POST("/", awsHandler.CreatePolicy)
	awsIamPolicy.PUT("/:arn", awsHandler.UpdatePolicy)
	awsIamPolicy.DELETE("/:arn", awsHandler.DeletePolicy)

	gcpIam := httpApi.Echo.Group("/gcp/iam")

	gcpIamGroup := gcpIam.Group("/groups")
	gcpIamGroup.GET("/", gcpHandler.ListGroups)
	gcpIamGroup.GET("/:name", gcpHandler.GetGroup)
	gcpIamGroup.POST("/", gcpHandler.CreateGroup)
	gcpIamGroup.PUT("/:name", gcpHandler.UpdateGroup)
	gcpIamGroup.DELETE("/:name", gcpHandler.DeleteGroup)

	gcpIamRole := gcpIam.Group("/roles")
	gcpIamRole.GET("/", gcpHandler.ListRoles)
	gcpIamRole.GET("/:name", gcpHandler.GetRole)
	gcpIamRole.POST("/", gcpHandler.CreateIamRole)
	gcpIamRole.PUT("/:id", gcpHandler.UpdateIamRole)
	gcpIamRole.DELETE("/:name", gcpHandler.DeleteIamRole)

	// gcpIamPolicies := gcpIam.Group("/policies")
	// gcpIamPolicies.GET("/", gcpHandler.)
	// gcpIamPolicies.GET("/:id", gcpHandler.GetRole)
	// gcpIamPolicies.POST("/", gcpHandler.CreateIamRole)
	// gcpIamPolicies.PUT("/:id", gcpHandler.UpdateIamRole)
	// gcpIamPolicies.DELETE("/:id", gcpHandler.DeleteIamRole)

	gcpIamServiceAccounts := gcpIam.Group("/service-accounts")
	gcpIamServiceAccounts.GET("/", gcpHandler.ListServiceAccounts)
	gcpIamServiceAccounts.GET("/:name", gcpHandler.GetServiceAccount)
	gcpIamServiceAccounts.POST("/", gcpHandler.CreateServiceAccount)
	// gcpIamServiceAccounts.PUT("/:id", gcpHandler.UpdateServiceAccount)
	gcpIamServiceAccounts.POST("/:name/enable", gcpHandler.EnableServiceAccount)
	gcpIamServiceAccounts.POST("/:name/disable", gcpHandler.DisableServiceAccount)
	gcpIamServiceAccounts.DELETE("/:name", gcpHandler.DeleteServiceAccount)

	gcpIamPolicies := gcpIam.Group("/policies")
	gcpIamPolicies.POST("/", gcpHandler.SetPolicy)
	gcpIamPolicies.GET("/", gcpHandler.GetPolicy)
	gcpIamPolicies.DELETE("/", gcpHandler.DeletePolicy)
}
