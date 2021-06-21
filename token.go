package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
)

type Credential struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(authPath string, config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	var tokFile string

	if len(authPath) == 0 {
		tokFile = fmt.Sprintf("%s/%s/%s", os.Getenv("HOME"), defaultWorkDirPrefix, tokenName)
	} else {
		tokFile = fmt.Sprintf("%s/%s", authPath, tokenName)
	}

	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// get client_id and client_secret from file
func getCredential(authPath string) (*Credential, error) {
	var credFile string

	if len(authPath) == 0 {
		credFile = fmt.Sprintf("%s/%s/%s", os.Getenv("HOME"), defaultWorkDirPrefix, credentialName)
	} else {
		credFile = fmt.Sprintf("%s/%s", authPath, credentialName)
	}

	cred, err := credentialFromFile(credFile)
	if err != nil {
		return nil, err
	}

	return cred, nil
}

// create oauto2 config from client ID and secret
func getOauthConfig(clientID string, clientSecret string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{drive.DriveScope},
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Retrieves a credential from a local file.
func credentialFromFile(file string) (*Credential, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	cred := &Credential{}
	err = json.NewDecoder(f).Decode(cred)
	return cred, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
