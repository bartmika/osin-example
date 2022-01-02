package cmd

import (
	"fmt"
	"log"

	"github.com/bartmika/osin-example/internal/controllers"
	"github.com/openshift/osin"
	"github.com/spf13/cobra"
)

var (
	addClientID          string
	addClientSecret      string
	addClientRedirectUri string
)

func init() {
	addClientCmd.Flags().StringVarP(&addClientID, "client_id", "a", "", "")
	addClientCmd.MarkFlagRequired("client_id")
	addClientCmd.Flags().StringVarP(&addClientSecret, "client_secret", "b", "", "-")
	addClientCmd.MarkFlagRequired("client_secret")
	addClientCmd.Flags().StringVarP(&addClientRedirectUri, "redirect_uri", "c", "", "-")
	addClientCmd.MarkFlagRequired("redirect_uri")
	rootCmd.AddCommand(addClientCmd)
}

var addClientCmd = &cobra.Command{
	Use:              "add_client",
	TraverseChildren: true,
	Short:            "Create a client",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("\033[H\033[2J") // Clear screen
		doRunAddClient()
	},
}

func doRunAddClient() {
	// Set our client.
	oastore := controllers.NewOSINRedisStorage()
	oastore.CreateClient(&osin.DefaultClient{
		Id:          addClientID,
		Secret:      addClientSecret,
		RedirectUri: addClientRedirectUri,
	})

	log.Println("The client has been created")
}
