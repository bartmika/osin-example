package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"golang.org/x/oauth2/clientcredentials"
)

var (
	ccClientID     string
	ccClientSecret string
	ccAuthorizeURL string
	ccTokenURL     string
	ccRedirectURL  string
)

func init() {
	clientCredentialCmd.Flags().StringVarP(&ccClientID, "client_id", "a", "", "-")
	clientCredentialCmd.MarkFlagRequired("client_id")
	clientCredentialCmd.Flags().StringVarP(&ccClientSecret, "client_secret", "b", "", "-")
	clientCredentialCmd.MarkFlagRequired("client_secret")
	clientCredentialCmd.Flags().StringVarP(&ccAuthorizeURL, "authorize_uri", "c", "http://localhost:8000/authorize", "-")
	clientCredentialCmd.MarkFlagRequired("authorize_uri")
	clientCredentialCmd.Flags().StringVarP(&ccTokenURL, "token_url", "d", "http://localhost:8000/token", "-")
	clientCredentialCmd.MarkFlagRequired("token_url")
	clientCredentialCmd.Flags().StringVarP(&ccRedirectURL, "redirect_uri", "e", "http://localhost:8000/appauth/code", "-")
	clientCredentialCmd.MarkFlagRequired("redirect_uri")
	rootCmd.AddCommand(clientCredentialCmd)
}

var clientCredentialCmd = &cobra.Command{
	Use:   "osin_client_credential",
	Short: "Get token from client credentials",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Beginning ClientCredential Based Authorization")

		cfg := clientcredentials.Config{
			ClientID:     ccClientID,
			ClientSecret: ccClientSecret,
			Scopes:       []string{"all"},
			TokenURL:     ccTokenURL,
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
