package handlers

import (
	"gitea/pcp-inariam/inariam/core/services/api/responses"
	"net/http"

	"github.com/labstack/echo/v4"
)

// @Summary		Health check endpoint for the API
// @Description	Verify if the server is running
// @ID				healthcheck
// @Tags			Misc
// @Success		200		{boolean}	bool
// @Failure		404		{object}	responses.Error
// @Router			/health [get]
func HealthCheck(ctx echo.Context) error {
	return responses.MessageResponse(ctx, http.StatusOK, "Server is working")
}
