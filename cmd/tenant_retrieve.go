package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bartmika/osin-example/internal/utils"
	"github.com/spf13/cobra"
)

// EXAMPLE USAGE
// go run main.go tenant_retrieve --id=1

var (
	tenantRetrieveID string
)

func init() {
	tenantRetrieveCmd.Flags().StringVarP(&tenantRetrieveID, "id", "a", "0", "Id of the tenant")
	tenantRetrieveCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(tenantRetrieveCmd)
}

var tenantRetrieveCmd = &cobra.Command{
	Use:              "tenant_retrieve -d -e",
	TraverseChildren: true,
	Short:            "Login a customer account",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("\033[H\033[2J") // Clear screen
		doRunTenantRetrieve()
	},
}

func doRunTenantRetrieve() {
	aUrl := "http://127.0.0.1:8000" + strings.Replace("/api/v1/tenant/xxx", "xxx", tenantRetrieveID, -1)

	// Create a Bearer string by appending string access token
	accessToken := os.Getenv("OSIN_EXAMPLE_CLI_ACCESS_TOKEN")
	var bearer = "Bearer " + accessToken

	client := &http.Client{}
	req, _ := http.NewRequest("GET", aUrl, nil)

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
