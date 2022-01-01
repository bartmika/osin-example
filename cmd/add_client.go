package cmd

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/bartmika/osin-example/internal/controllers"
	"github.com/bartmika/osin-example/internal/models"
	repo "github.com/bartmika/osin-example/internal/repositories"
	"github.com/bartmika/osin-example/internal/utils"
	"github.com/openshift/osin"
	"github.com/spf13/cobra"
)

var (
	addClientID          string
	addClientSecret      string
	addClientRedirectUri string
	addClientUserID      string
)

func init() {
	addClientCmd.Flags().StringVarP(&addClientID, "client_id", "a", "", "")
	addClientCmd.MarkFlagRequired("client_id")
	addClientCmd.Flags().StringVarP(&addClientSecret, "client_secret", "b", "", "-")
	addClientCmd.MarkFlagRequired("client_secret")
	addClientCmd.Flags().StringVarP(&addClientRedirectUri, "redirect_uri", "c", "", "-")
	addClientCmd.MarkFlagRequired("redirect_uri")
	addClientCmd.Flags().StringVarP(&addClientUserID, "user_id", "d", "", "-")
	addClientCmd.MarkFlagRequired("user_id")
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
	// Load up our database.
	db, err := utils.ConnectDB(
		databaseHost,
		databasePort,
		databaseUser,
		databasePassword,
		databaseName,
		"public",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Load up our repositories.
	ur := repo.NewUserRepo(db)

	userID, err := strconv.ParseUint(addClientUserID, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	// Lookup the user.
	u, err := ur.GetByID(context.Background(), userID)
	if err != nil {
		log.Fatal(err)
	}

	userData := &models.UserLite{
		ID:       u.ID,
		UUID:     u.UUID,
		TenantID: u.TenantID,
		Name:     u.Name,
		State:    u.State,
		RoleID:   u.RoleID,
		Timezone: u.Timezone,
		Language: u.Language,
	}

	// Set our client.
	oastore := controllers.NewOSINRedisStorage()
	oastore.CreateClient(&osin.DefaultClient{
		Id:          addClientID,
		Secret:      addClientSecret,
		RedirectUri: addClientRedirectUri,
		UserData:    userData,
	})

	log.Println("The client has been created")
}
