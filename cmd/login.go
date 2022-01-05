package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// EXAMPLE USAGE:
// go run main.go login -d=demo@demo.com -e=123password

var (
	signinEmail    string
	signinPassword string
)

func init() {
	loginCmd.Flags().StringVarP(&signinEmail, "email", "d", "", "Email of the user account")
	loginCmd.MarkFlagRequired("email")
	loginCmd.Flags().StringVarP(&signinPassword, "password", "e", "", "Password of the user account")
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:              "login -d -e",
	TraverseChildren: true,
	Short:            "Login a customer account",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("\033[H\033[2J") // Clear screen
		doRunLogin()
	},
}

type LoginRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginResponse struct {
	TenantID         uint64 `json:"tenant_id"`
	TenantSchemaName string `json:"tenant_schema_name"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	RoleID           int8   `json:"role_id"`
	Language         string `json:"language"`

	// https://pkg.go.dev/golang.org/x/oauth2#Token
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type,omitempty"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry,omitempty"`
}

func doRunLogin() {
	aUrl := "http://127.0.0.1:8000/api/v1/login"
	data := &LoginRequest{
		Password: signinPassword,
		Email:    signinEmail,
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
	// fmt.Println(string(responseBytes)) // Return what was returned from API.

	var responseData LoginResponse

	// De-serialize bytes into our struct object.
	err = json.Unmarshal(responseBytes, &responseData)
	if err != nil {
		log.Println(string(responseBytes))
		log.Fatalf("Unmarshal | An Error Occured %v", err)
	}

	// Output message.
	fmt.Println("Raw:", responseData)
	fmt.Println("FirstName:", responseData.FirstName)
	fmt.Println("LastName:", responseData.LastName)
	fmt.Println("Email:", responseData.Email)
	fmt.Println("RoleID:", responseData.RoleID)
	fmt.Println("TenantID:", responseData.TenantID)
	fmt.Println("Language:", responseData.Language)
	fmt.Println("AccessToken:", responseData.AccessToken)
	fmt.Println("TokenType:", responseData.TokenType)
	fmt.Println("RefreshToken:", responseData.RefreshToken)
	fmt.Println("Expiry:", responseData.Expiry)

	// Output message.
	fmt.Printf("Please run in your console:\n\nexport OSIN_EXAMPLE_CLI_ACCESS_TOKEN=%s\n\n", responseData.AccessToken)
	fmt.Printf("export OSIN_EXAMPLE_CLI_REFRESH_TOKEN=%s\n\n", responseData.RefreshToken)
}
