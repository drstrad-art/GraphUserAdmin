package main

import (
	"fmt"
	"os"

	"GraphUserAdmin/internal/auth"
	"GraphUserAdmin/internal/config"

	"github.com/spf13/cobra"
)

var (
	configPath string
	cfg        *config.Config
	token      string
	version    = "1.0.0"
	buildDate  = "unknown"
	verbose    bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gua",
		Short: "GraphUserAdmin - Microsoft 365 User & License Management using Microsoft Graph API",
		Long: `GraphUserAdmin (gua) is a command-line tool for managing Microsoft 365 users, licenses, and groups
using the Microsoft Graph REST API with client credentials authentication.`,
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip authentication for help command
			if cmd.Name() == "help" {
				return nil
			}

			// Load configuration
			var err error
			cfg, err = config.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Authenticate
			if verbose {
				fmt.Println("Authenticating with Microsoft Graph...")
				fmt.Printf("Tenant ID: %s\n", cfg.TenantID)
				fmt.Printf("Client ID: %s\n", cfg.ClientID)
			} else {
				fmt.Println("Authenticating with Microsoft Graph...")
			}
			token, err = auth.GetAccessToken(cfg.TenantID, cfg.ClientID, cfg.ClientSecret)
			if err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}
			if verbose {
				fmt.Println("✓ Authentication successful!")
				fmt.Printf("Token length: %d characters\n", len(token))
			} else {
				fmt.Println("✓ Authentication successful!")
			}
			fmt.Println()

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.json", "Path to configuration file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output for debugging")

	// Setup all commands (users, licenses, groups)
	setupCommands(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
