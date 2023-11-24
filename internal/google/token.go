package google

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
)

const TokFile = "token.json"

func GetTokenFromFile(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	tok := &oauth2.Token{}
	if err := json.NewDecoder(f).Decode(tok); err != nil {
		return nil, err
	}
	return tok, err
}

func SaveToken(token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", TokFile)
	f, err := os.OpenFile(TokFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Printf("Unable to cache oauth token: %v", err)
		return err
	}
	defer func() { _ = f.Close() }()
	if err := json.NewEncoder(f).Encode(token); err != nil {
		return err
	}
	return nil
}
