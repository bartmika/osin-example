package cmd

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/openshift/osin"
	"github.com/spf13/cobra"

	"github.com/bartmika/osin-example/internal/controllers"
	"github.com/bartmika/osin-example/internal/models"
	repo "github.com/bartmika/osin-example/internal/repositories"
	"github.com/bartmika/osin-example/internal/utils"
)

var (
	addApplicationName        string
	addApplicationDescription string
	addApplicationWebsiteURL  string
	addApplicationScope       string
	addApplicationRedirectURL string
	addApplicationImageURL    string
	addApplicationTenantID    string
)

func init() {
	addApplicationCmd.Flags().StringVarP(&addApplicationName, "name", "a", "", "")
	addApplicationCmd.MarkFlagRequired("name")
	addApplicationCmd.Flags().StringVarP(&addApplicationDescription, "description", "b", "", "-")
	addApplicationCmd.MarkFlagRequired("description")
	addApplicationCmd.Flags().StringVarP(&addApplicationWebsiteURL, "website_url", "c", "", "-")
	addApplicationCmd.MarkFlagRequired("website_url")
	addApplicationCmd.Flags().StringVarP(&addApplicationScope, "scope", "d", "", "-")
	addApplicationCmd.MarkFlagRequired("scope")
	addApplicationCmd.Flags().StringVarP(&addApplicationRedirectURL, "redirect_url", "e", "", "-")
	addApplicationCmd.MarkFlagRequired("redirect_url")
	addApplicationCmd.Flags().StringVarP(&addApplicationImageURL, "image_url", "f", "", "-")
	addApplicationCmd.MarkFlagRequired("image_url")
	addApplicationCmd.Flags().StringVarP(&addApplicationTenantID, "tenant_id", "g", "", "-")
	addApplicationCmd.MarkFlagRequired("tenant_id")
	rootCmd.AddCommand(addApplicationCmd)
}

var addApplicationCmd = &cobra.Command{
	Use:              "add_application",
	TraverseChildren: true,
	Short:            "Create a third-party application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("\033[H\033[2J") // Clear screen
		doRunAddApplication()
	},
}

func doRunAddApplication() {
	tenantID, err := strconv.ParseUint(addApplicationTenantID, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

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
	ar := repo.NewApplicationRepo(db)

	// Create the clientID and clientSecret values.
	clientID := utils.RandomBase16String(16)
	clientSecret := utils.RandomBase16String(255)

	// Create our application in our database.

	m := &models.Application{
		TenantID:     tenantID,
		UUID:         uuid.NewString(),
		Name:         addApplicationName,
		Description:  addApplicationDescription,
		WebsiteURL:   addApplicationWebsiteURL,
		Scope:        addApplicationScope,
		RedirectURL:  addApplicationRedirectURL,
		ImageURL:     addApplicationImageURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		State:        models.ApplicationRunningState,
	}
	err = ar.Insert(context.Background(), m)
	if err != nil {
		log.Fatal(err)
	}
	m, err = ar.GetByUUID(context.Background(), m.UUID)

	// Create our oAuth 2.0 client in the storage.

	//     string
	// addApplicationTenantID    string
	oastore := controllers.NewOSINRedisStorage()
	oastore.CreateClient(&osin.DefaultClient{
		Id:          clientID,
		Secret:      clientSecret,
		RedirectUri: addApplicationRedirectURL,
	})

	// Return the results from the database.
	log.Println("The application has been created")
	log.Println("TenantID:", m.TenantID)
	log.Println("ID:", m.ID)
	log.Println("UUID:", m.UUID)
	log.Println("Name:", m.Name)
	log.Println("Description:", m.Description)
	log.Println("Scope:", m.Scope)
	log.Println("RedirectURL:", m.RedirectURL)
	log.Println("ImageURL:", m.ImageURL)
	log.Println("ClientID:", clientID)
	log.Println("ClientSecret:", clientSecret)

}
