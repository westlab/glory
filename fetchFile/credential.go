package fetchFile

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func getClient(config *oauth2.Config) (*http.Client, error) {
	tokenFile := "credential/token.json"
	tok, err := getTokenFromFile(tokenFile)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		saveToken(tokenFile, tok)
	}
	return config.Client(context.Background(), tok), nil
}

func getTokenFromFile(fileName string) (*oauth2.Token, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web %v", err)
	}
	return tok, nil
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// GetCredential はGoogle Drive APIのクライアントを返す
func GetCredential() (*drive.Service, error) {
	b, err := ioutil.ReadFile("credential/credentials.json")
	if err != nil {
		return nil, fmt.Errorf("read credentials error: %w", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %w", err)
	}
	client, err := getClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create API client: %w", err)
	}
	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))

	return srv, nil

}
