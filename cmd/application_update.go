package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bartmika/osin-example/internal/idos"
	"github.com/bartmika/osin-example/internal/utils"
)

// Example Usage:
// go run main.go application_update -a=3 -b="Third Party" -c="Demonstration purposes only" -d=http://demo.com -e=all -f=http://127.0.0.1:8001/appauth/code -g=https://g.com/img.png

var (
	applicationUpdateID          string
	applicationUpdateName        string
	applicationUpdateDescription string
	applicationUpdateWebsiteURL  string
	applicationUpdateScope       string
	applicationUpdateRedirectURL string
	applicationUpdateImageURL    string
)

func init() {
	applicationUpdateCmd.Flags().StringVarP(&applicationUpdateID, "id", "a", "", "")
	applicationUpdateCmd.Flags().StringVarP(&applicationUpdateName, "name", "b", "", "")
	applicationUpdateCmd.Flags().StringVarP(&applicationUpdateDescription, "description", "c", "", "-")
	applicationUpdateCmd.Flags().StringVarP(&applicationUpdateWebsiteURL, "website_url", "d", "", "-")
	applicationUpdateCmd.Flags().StringVarP(&applicationUpdateScope, "scope", "e", "", "-")
	applicationUpdateCmd.Flags().StringVarP(&applicationUpdateRedirectURL, "redirect_url", "f", "", "-")
	applicationUpdateCmd.Flags().StringVarP(&applicationUpdateImageURL, "image_url", "g", "", "-")
	rootCmd.AddCommand(applicationUpdateCmd)
}

var applicationUpdateCmd = &cobra.Command{
	Use:              "application_update",
	TraverseChildren: true,
	Short:            "Update a third-party application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("\033[H\033[2J") // Clear screen
		doRunApplicationUpdate()
	},
}

func doRunApplicationUpdate() {
	aUrl := "http://127.0.0.1:8000" + strings.Replace("/api/v1/application/xxx", "xxx", applicationUpdateID, -1)
	data := idos.ApplicationUpdateRequestIDO{
		Name:        applicationUpdateName,
		Description: applicationUpdateName,
		WebsiteURL:  applicationUpdateWebsiteURL,
		Scope:       applicationUpdateScope,
		RedirectURL: applicationUpdateRedirectURL,
		ImageURL:    applicationUpdateImageURL,
	}
	dataBytes, _ := json.Marshal(data)
	requestBodyBuf := bytes.NewBuffer(dataBytes)

	// Update a Bearer string by appending string access token
	accessToken := os.Getenv("OSIN_EXAMPLE_CLI_ACCESS_TOKEN")
	var bearer = "Bearer " + accessToken

	client := &http.Client{}
	req, _ := http.NewRequest("PUT", aUrl, requestBodyBuf)

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
