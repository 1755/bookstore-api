package cmd

import (
	"fmt"

	"github.com/1755/bookstore-api/api"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Run the API server",
	Long:  "This command runs the API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			return fmt.Errorf("failed to get config path: %w", err)
		}

		app, cleanup, err := api.InjectApplication(api.ConfigPath(configPath))
		if err != nil {
			return fmt.Errorf("failed to create application: %w", err)
		}
		defer cleanup()

		if err := app.Run(); err != nil {
			return fmt.Errorf("failed to run application: %w", err)
		}

		return nil
	},
}
