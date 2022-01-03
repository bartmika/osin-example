package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"golang.org/x/oauth2"
)

var (
	passClientID     string
	passClientSecret string
	passAuthorizeURL string
	passTokenURL     string
	passRedirectURL  string
	passPassword     string
	passEmail        string
)

func init() {
	passwordCmd.Flags().StringVarP(&passClientID, "client_id", "a", "", "-")
	passwordCmd.MarkFlagRequired("client_id")
	passwordCmd.Flags().StringVarP(&passClientSecret, "client_secret", "b", "", "-")
	passwordCmd.MarkFlagRequired("client_secret")
	passwordCmd.Flags().StringVarP(&passAuthorizeURL, "authorize_uri", "c", "", "-")
	passwordCmd.MarkFlagRequired("authorize_uri")
	passwordCmd.Flags().StringVarP(&passTokenURL, "token_url", "d", "", "-")
	passwordCmd.MarkFlagRequired("token_url")
	passwordCmd.Flags().StringVarP(&passRedirectURL, "redirect_uri", "e", "", "-")
	passwordCmd.MarkFlagRequired("redirect_uri")
	passwordCmd.Flags().StringVarP(&passEmail, "email", "f", "", "Email of the user account")
	passwordCmd.MarkFlagRequired("email")
	passwordCmd.Flags().StringVarP(&passPassword, "password", "g", "", "Password of the user account")
	passwordCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(passwordCmd)
}

var passwordCmd = &cobra.Command{
	Use:   "osin_password",
	Short: "Get token from username and password",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Beginning Password Based Authorization")

		client := &oauth2.Config{
			ClientID:     passClientID,
			ClientSecret: passClientSecret,
			Scopes:       []string{"all"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  passAuthorizeURL,
				TokenURL: passTokenURL,
			},
			RedirectURL: passRedirectURL,
		}

		// NOTE: https://pkg.go.dev/golang.org/x/oauth2#Config.PasswordCredentialsToken
		token, err := client.PasswordCredentialsToken(context.Background(), passEmail, passPassword)
		if err != nil {
			log.Fatal(err)
		}

		// NOTE: https://pkg.go.dev/golang.org/x/oauth2#Token
		log.Println("AccessToken", token.AccessToken)
		log.Println("TokenType", token.TokenType)
		log.Println("RefreshToken", token.RefreshToken)
		log.Println("Expiry", token.Expiry)
		log.Println("UserData|TenantID", token.Extra("tenant_id"))
		log.Println("UserData|UserID", token.Extra("user_id"))
		log.Println("UserData|UserUUID", token.Extra("user_uuid"))

		// Output message.
		fmt.Printf("Please run in your console:\n\nexport OSIN_EXAMPLE_CLI_ACCESS_TOKEN=%s\n\n", token.AccessToken)
		fmt.Printf("export OSIN_EXAMPLE_CLI_REFRESH_TOKEN=%s\n\n", token.RefreshToken)
	},
}
