package cmd

import (
	"fmt"
	"os"

	// homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

var (
	databaseHost                    string
	databasePort                    string
	databaseUser                    string
	databasePassword                string
	databaseName                    string
	applicationSigningKey           string
	applicationAddress              string
	applicationFrontendClientID     string
	applicationFrontendClientSecret string
	applicationFrontendReturnURL    string
)

// Initialize function will be called when every command gets called.
func init() {
	// Get our environment variables which will used to configure our application and save across all the sub-commands.
	rootCmd.PersistentFlags().StringVar(&databaseHost, "dbHost", os.Getenv("OSIN_DB_HOST"), "The address of database.")
	rootCmd.PersistentFlags().StringVar(&databasePort, "dbPort", os.Getenv("OSIN_DB_PORT"), "The port of database.")
	rootCmd.PersistentFlags().StringVar(&databaseUser, "dbUser", os.Getenv("OSIN_DB_USER"), "The database user.")
	rootCmd.PersistentFlags().StringVar(&databasePassword, "dbPassword", os.Getenv("OSIN_DB_PASSWORD"), "The database password.")
	rootCmd.PersistentFlags().StringVar(&databaseName, "dbName", os.Getenv("OSIN_DB_NAME"), "The database name.")
	rootCmd.PersistentFlags().StringVar(&applicationSigningKey, "appSignKey", os.Getenv("OSIN_APP_SIGNING_KEY"), "The signing key.")
	rootCmd.PersistentFlags().StringVar(&applicationAddress, "applicationAddress", os.Getenv("OSIN_APP_ADDRESS"), "Application address")
	rootCmd.PersistentFlags().StringVar(&applicationFrontendClientID, "clientID", os.Getenv("OSIN_APP_FRONTEND_CLIENT_ID"), "Client ID")
	rootCmd.PersistentFlags().StringVar(&applicationFrontendClientSecret, "clientSecret", os.Getenv("OSIN_APP_FRONTEND_CLIENT_SECRET"), "Client Secret")
	rootCmd.PersistentFlags().StringVar(&applicationFrontendReturnURL, "returnURL", os.Getenv("OSIN_APP_FRONTEND_RETURN_URL"), "Client Return URL")
}

var rootCmd = &cobra.Command{
	Use:   "osgin-example",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing.
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
