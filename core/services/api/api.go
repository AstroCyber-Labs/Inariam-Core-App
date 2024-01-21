// Package api provides an API structure and functions for running an Echo web server.
package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"gitea/pcp-inariam/inariam/core/config"
	"gitea/pcp-inariam/inariam/pkgs/log"
)

// API represents the main structure for the API.
type API struct {
	Config *config.Config
	Echo   *echo.Echo
	DB     *gorm.DB
}

// New creates a new instance of the API with the provided configuration.
func New(config *config.Config) *API {
	e := echo.New()
	e.HideBanner = true

	return &API{
		Config: config,
		Echo:   e,
		DB:     nil,
	}
}

// RunServer starts the Echo web server based on the API configuration.
func (api *API) RunServer() error {
	log.Logger.Infof("Server is running on port %d\n", api.Config.APIConfig.Port)
	return api.Echo.Start(fmt.Sprintf(":%d", api.Config.APIConfig.Port))
}
