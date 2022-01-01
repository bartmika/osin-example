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
	regTenantId  int
	regFirstName string
	regLastName  string
	regEmail     string
	regPassword  string
	regLanguage  string
)

func init() {
	registerCmd.Flags().StringVarP(&regFirstName, "first_name", "b", "", "First name of the user account")
	registerCmd.MarkFlagRequired("first_name")
	registerCmd.Flags().StringVarP(&regLastName, "last_name", "c", "", "Last name of the user account")
	registerCmd.MarkFlagRequired("last_name")
	registerCmd.Flags().StringVarP(&regEmail, "email", "d", "", "Email of the user account")
	registerCmd.MarkFlagRequired("email")
	registerCmd.Flags().StringVarP(&regPassword, "password", "e", "", "Password of the user account")
	registerCmd.MarkFlagRequired("password")
	registerCmd.Flags().StringVarP(&regLanguage, "language", "f", "", "Language of the user account")
	registerCmd.MarkFlagRequired("language")
	rootCmd.AddCommand(registerCmd)
}

var registerCmd = &cobra.Command{
	Use:              "register -a -b -c -d -e -f",
	TraverseChildren: true,
	Short:            "Register a customer account",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("\033[H\033[2J") // Clear screen
		doRunRegister()
	},
}

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Language  string `json:"language"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

func doRunRegister() {
	aUrl := "http://127.0.0.1:8000/api/v1/register"
	data := &RegisterRequest{
		FirstName: regFirstName,
		LastName:  regLastName,
		Password:  regPassword,
		Email:     regEmail,
		Language:  regLanguage,
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

	var responseData RegisterResponse

	// De-serialize bytes into our struct object.
	err = json.Unmarshal(responseBytes, &responseData)
	if err != nil {
		log.Println(string(responseBytes))
		log.Fatalf("Unmarshal | An Error Occured %v", err)
	}

	// Output message.
	fmt.Printf(responseData.Message)
}
