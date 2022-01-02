package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	rtaGrantType    string
	rtaRefreshToken string
)

func init() {
	refreshTokenAPICmd.Flags().StringVarP(&rtaGrantType, "grant_type", "a", "refresh_token", "-")
	refreshTokenAPICmd.Flags().StringVarP(&rtaRefreshToken, "refresh_token", "b", "", "-")
	refreshTokenAPICmd.MarkFlagRequired("refresh_token")

	rootCmd.AddCommand(refreshTokenAPICmd)
}

var refreshTokenAPICmd = &cobra.Command{
	Use:              "refresh_token -d -e",
	TraverseChildren: true,
	Short:            "Login a customer account",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("\033[H\033[2J") // Clear screen
		doRunRefreshTokenAPI()
	},
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope_type"`
	TokenType    string `json:"token_type"`
}

func doRunRefreshTokenAPI() {
	aUrl := "http://127.0.0.1:8000/api/v1/refresh-token"
	data := &RefreshTokenRequest{
		RefreshToken: rtaRefreshToken,
		GrantType:    rtaGrantType,
	}
	dataBin, _ := json.Marshal(data)
	requestBodyBuf := bytes.NewBuffer(dataBin)

	resp, err := http.Post(aUrl, "application/json", requestBodyBuf)
	if err != nil {
		log.Fatalf("Post | An Error Occured %v", err)
	}

	defer resp.Body.Close()

	// Read the response body
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll | An Error Occured %v", err)
	}
	fmt.Println("-->", string(responseBytes)) // Return what was returned from API.

	var responseData RefreshTokenResponse

	// De-serialize bytes into our struct object.
	err = json.Unmarshal(responseBytes, &responseData)
	if err != nil {
		log.Println(string(responseBytes))
		log.Fatalf("Unmarshal | An Error Occured %v", err)
	}

	// Output message.
	fmt.Println("Raw:", responseData)
	fmt.Println("AccessToken:", responseData.AccessToken)
	fmt.Println("ExpiresIn:", responseData.ExpiresIn)
	fmt.Println("RefreshToken:", responseData.RefreshToken)
	fmt.Println("Scope:", responseData.Scope)
	fmt.Println("TokenType:", responseData.TokenType)

	// Output message.
	fmt.Printf("Please run in your console:\n\nexport OSIN_EXAMPLE_CLI_ACCESS_TOKEN=%s\n\n", responseData.AccessToken)
	fmt.Printf("export OSIN_EXAMPLE_CLI_REFRESH_TOKEN=%s\n\n", responseData.RefreshToken)
}
