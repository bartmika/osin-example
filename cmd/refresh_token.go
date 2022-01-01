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

//
// var (
// 	signinEmail    string
// 	signinPassword string
// )

func init() {
	// refreshTokenCmd.Flags().StringVarP(&signinEmail, "email", "d", "", "Email of the user account")
	// refreshTokenCmd.MarkFlagRequired("email")
	// refreshTokenCmd.Flags().StringVarP(&signinPassword, "password", "e", "", "Password of the user account")
	rootCmd.AddCommand(refreshTokenCmd)
}

var refreshTokenCmd = &cobra.Command{
	Use:              "refresh_token -d -e",
	TraverseChildren: true,
	Short:            "Login a customer account",
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
	preq.SetBasicAuth("1234", "aabbccdd")

	pclient := &http.Client{}
	presp, err := pclient.Do(preq)
	if err != nil {
		log.Fatal(err)
	}
	defer presp.Body.Close()

	if presp.StatusCode != 200 {
		log.Fatal("Not 200!")
	}
	//
	// resp, err := http.PostForm(aURL, url.Values{
	// 	"grant_type":    {"refresh_token"},
	// 	"refresh_token": {refreshToken},
	// })
	//
	// if err != nil {
	// 	log.Fatalf("Post | An Error Occured %v", err)
	// }
	//
	// defer resp.Body.Close()

	// Read the response body
	responseBytes, err := ioutil.ReadAll(presp.Body)
	if err != nil {
		log.Fatalf("ReadAll | An Error Occured %v", err)
	}

	fmt.Println(utils.JsonPrettyPrint(string(responseBytes)))
}
