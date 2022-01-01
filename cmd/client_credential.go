package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"golang.org/x/oauth2/clientcredentials"
)

func init() {
	rootCmd.AddCommand(clientCredentialCmd)
}

var clientCredentialCmd = &cobra.Command{
	Use:   "client_credential",
	Short: "Get token from client credentials",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Beginning ClientCredential Based Authorization")

		cfg := clientcredentials.Config{
			ClientID:     "1234",
			ClientSecret: "aabbccdd",
			Scopes:       []string{"all"},
			TokenURL:     "http://localhost:8000/token",
		}

		// https://github.com/go-oauth2/oauth2/blob/b208c14e621016995debae2fa7dc20c8f0e4f6f8/example/client/client.go#L116
		token, err := cfg.Token(context.Background())
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
