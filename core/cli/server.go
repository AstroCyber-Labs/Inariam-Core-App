package cli

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"

	"gitea/pcp-inariam/inariam/core/config"
	"gitea/pcp-inariam/inariam/core/services/api"
	"gitea/pcp-inariam/inariam/core/services/api/routes"
	"gitea/pcp-inariam/inariam/pkgs/log"
)

const (
	ErrorRunningServer = "error running server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run Inariam server",
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.New()
		if err != nil {
			log.Logger.Panicf("error loading configuration %w", err)
		}

		api := api.New(cfg)

		routes.ConfigureRoutes(api)

		data, err := json.MarshalIndent(api.Echo.Routes(), "", "  ")
		if err != nil {
			log.Logger.Panicln(err)
		}
		os.WriteFile("routes.json", data, 0644)

		err = api.RunServer()
		if err != nil {
			log.Logger.Panicln(err.Error())
		}
	},
}
