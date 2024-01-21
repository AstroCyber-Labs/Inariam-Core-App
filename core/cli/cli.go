package cli

import (
	"gitea/pcp-inariam/inariam/core/config"

	"gitea/pcp-inariam/inariam/pkgs/log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Inariam",
	Short: "Your app's short description",
	Long:  `A longer description of your application.`,
}

func Execute() {
	color.Green("[+] Booting up....")
	color.Yellow(`
  ___                  _                 
 |_ _|_ __   __ _ _ __(_) __ _ _ __ ___  
  | || '_ \ / _\ | __| |  | | '_ \ _ \ 
  | || | | | (_| | |  | | (_| | | | | | |
 |___|_| |_|\__,_|_|  |_|\__,_|_| |_| |_| is here ! v0.0.1



`)
	if err := rootCmd.Execute(); err != nil {
		log.Logger.Panicln(err)
	}
}

func init() {

	cobra.OnInitialize(config.InitConfig)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(configCmd)

	userCmd.AddCommand(newSignUpCmd())
	rootCmd.AddCommand(userCmd)

}
