package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"golang.org/x/oauth2"
)

var (
	loginPassword string
	loginEmail    string
)

func init() {
	passwordCmd.Flags().StringVarP(&loginEmail, "email", "d", "", "Email of the user account")
	passwordCmd.MarkFlagRequired("email")
	passwordCmd.Flags().StringVarP(&loginPassword, "password", "e", "", "Password of the user account")
	passwordCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(passwordCmd)
}

var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "Get token from username and password",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Beginning Password Based Authorization")

		client := &oauth2.Config{
			ClientID:     "1234",
			ClientSecret: "aabbccdd",
			Scopes:       []string{"all"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "http://localhost:8000/authorize",
				TokenURL: "http://localhost:8000/token",
			},
			RedirectURL: "http://localhost:8000/appauth/code",
		}

		// NOTE: https://pkg.go.dev/golang.org/x/oauth2#Config.PasswordCredentialsToken
		token, err := client.PasswordCredentialsToken(context.Background(), loginEmail, loginPassword)
		if err != nil {
			log.Fatal(err)
		}

		// NOTE: https://pkg.go.dev/golang.org/x/oauth2#Token
		log.Println("AccessToken", token.AccessToken)
		log.Println("TokenType", token.TokenType)
		log.Println("RefreshToken", token.RefreshToken)
		log.Println("Expiry", token.Expiry)
		log.Println("UserID", token.Extra("custom_parameter"))

	},
}
