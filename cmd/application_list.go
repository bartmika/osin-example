package cmd

import (
	// "os"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/bartmika/osin-example/internal/utils"
)

var (
	applicationListPageToken string = "100"
	applicationListPageSize  string = "0"
)

// EXAMPLE USAGE
// go run main.go application_list

func init() {
	listapplicationsCmd.Flags().StringVarP(&applicationListPageToken, "page_token", "p", "0", "The page we are on.")
	listapplicationsCmd.Flags().StringVarP(&applicationListPageSize, "page_size", "s", "100", "The number to paginate the results by")
	rootCmd.AddCommand(listapplicationsCmd)
}

var listapplicationsCmd = &cobra.Command{
	Use:              "application_list",
	TraverseChildren: true,
	Short:            "List in a paginated manner",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("\033[H\033[2J") // Clear screen
		doRunListapplications()
	},
}

func doRunListapplications() {
	// Generate the base URL with the URL parameters to do specific filtering of our data.
	aUrl := applicationAddress + "/api/v1/applications" + "?page_size=" + applicationListPageSize + "&page_token=" + applicationListPageToken

	client := &http.Client{}
	req, _ := http.NewRequest("GET", aUrl, nil)

	// Create a Bearer string by appending string access token
	accessToken := os.Getenv("OSIN_EXAMPLE_CLI_ACCESS_TOKEN")
	var bearer = "Bearer " + accessToken

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	defer resp.Body.Close()

	// Read the response body
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Print a pretty formatted JSON output.
	fmt.Println(utils.JsonPrettyPrint(string(responseBytes)))
}
