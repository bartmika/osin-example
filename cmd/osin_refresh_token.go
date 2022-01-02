package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/bartmika/osin-example/internal/utils"
	"github.com/spf13/cobra"
)

var (
	rtClientID     string
	rtClientSecret string
	rtAuthorizeURL string
	rtTokenURL     string
	rtRedirectURL  string
)

func init() {
	refreshTokenCmd.Flags().StringVarP(&rtClientID, "client_id", "a", "", "-")
	refreshTokenCmd.MarkFlagRequired("client_id")
	refreshTokenCmd.Flags().StringVarP(&rtClientSecret, "client_secret", "b", "", "-")
	refreshTokenCmd.MarkFlagRequired("client_secret")
	refreshTokenCmd.Flags().StringVarP(&rtAuthorizeURL, "authorize_uri", "c", "http://localhost:8000/authorize", "-")
	refreshTokenCmd.MarkFlagRequired("authorize_uri")
	refreshTokenCmd.Flags().StringVarP(&rtTokenURL, "token_url", "d", "http://localhost:8000/token", "-")
	refreshTokenCmd.MarkFlagRequired("token_url")
	refreshTokenCmd.Flags().StringVarP(&rtRedirectURL, "redirect_uri", "e", "http://localhost:8000/appauth/code", "-")
	rootCmd.AddCommand(refreshTokenCmd)
}

var refreshTokenCmd = &cobra.Command{
	Use:              "osin_refresh_token -d -e",
	TraverseChildren: true,
	Short:            "Refresh the token",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("\033[H\033[2J") // Clear screen
		doRunRefreshToken()
	},
}

type RefreshTokenRequest struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

func doRunRefreshToken() {
	accessToken := os.Getenv("OSIN_EXAMPLE_CLI_ACCESS_TOKEN")
	refreshToken := os.Getenv("OSIN_EXAMPLE_CLI_REFRESH_TOKEN")

	endpoint := fmt.Sprintf("http://127.0.0.1:8000/token")
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	// For debugging purposes only.

	log.Println(accessToken)
	log.Println(refreshToken)
	log.Println(url.QueryEscape(refreshToken))
	log.Println(endpoint)

	//
	// Submit the code.
	//

	preq, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	preq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	preq.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	preq.SetBasicAuth(rtClientID, rtClientSecret)

	pclient := &http.Client{}
	presp, err := pclient.Do(preq)
	if err != nil {
		log.Fatal(err)
	}
	defer presp.Body.Close()

	if presp.StatusCode != 200 {
		log.Fatal("Not 200!")
	}

	// Read the response body
	responseBytes, err := ioutil.ReadAll(presp.Body)
	if err != nil {
		log.Fatalf("ReadAll | An Error Occured %v", err)
	}

	fmt.Println(utils.JsonPrettyPrint(string(responseBytes)))
}
