package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/bartmika/osin-example/internal/idos"
	"github.com/bartmika/osin-example/internal/utils"
)

// Example Usage:
// go run main.go application_create -a="Third Party" -b="Demonstration purposes only" -c=http://demo.com -d=all -e=http://127.0.0.1:8001/appauth/code -f=https://g.com/img.png -g=1

var (
	applicationCreateName        string
	applicationCreateDescription string
	applicationCreateWebsiteURL  string
	applicationCreateScope       string
	applicationCreateRedirectURL string
	applicationCreateImageURL    string
	applicationCreateTenantID    string
)

func init() {
	applicationCreateCmd.Flags().StringVarP(&applicationCreateName, "name", "a", "", "")
	applicationCreateCmd.Flags().StringVarP(&applicationCreateDescription, "description", "b", "", "-")
	applicationCreateCmd.Flags().StringVarP(&applicationCreateWebsiteURL, "website_url", "c", "", "-")
	applicationCreateCmd.Flags().StringVarP(&applicationCreateScope, "scope", "d", "", "-")
	applicationCreateCmd.Flags().StringVarP(&applicationCreateRedirectURL, "redirect_url", "e", "", "-")
	applicationCreateCmd.Flags().StringVarP(&applicationCreateImageURL, "image_url", "f", "", "-")
	applicationCreateCmd.Flags().StringVarP(&applicationCreateTenantID, "tenant_id", "g", "", "-")
	rootCmd.AddCommand(applicationCreateCmd)
}

var applicationCreateCmd = &cobra.Command{
	Use:              "application_create",
	TraverseChildren: true,
	Short:            "Create a third-party application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("\033[H\033[2J") // Clear screen
		doRunapplicationCreate()
	},
}

func doRunapplicationCreate() {
	aUrl := applicationAddress + "/api/v1/applications"
	data := idos.ApplicationCreateRequestIDO{
		Name:        applicationCreateName,
		Description: applicationCreateName,
		WebsiteURL:  applicationCreateWebsiteURL,
		Scope:       applicationCreateScope,
		RedirectURL: applicationCreateRedirectURL,
		ImageURL:    applicationCreateImageURL,
	}
	dataBytes, _ := json.Marshal(data)
	requestBodyBuf := bytes.NewBuffer(dataBytes)

	// Create a Bearer string by appending string access token
	accessToken := os.Getenv("OSIN_EXAMPLE_CLI_ACCESS_TOKEN")
	var bearer = "Bearer " + accessToken

	client := &http.Client{}
	req, _ := http.NewRequest("POST", aUrl, requestBodyBuf)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Set data format.
	req.Header.Set("Content-Type", "application/json")

	// Send req using http Client
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	defer resp.Body.Close()

	// Read the response body
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll | An Error Occured %v", err)
	}

	fmt.Println(utils.JsonPrettyPrint(string(responseBytes)))
}
