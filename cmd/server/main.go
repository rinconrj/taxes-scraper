package main

// Import necessary libraries
import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	google2 "github.com/rinconrj/golang-scraper/internal/google"
	"github.com/spf13/viper"

	"github.com/rinconrj/golang-scraper/internal/contaja"
)

const tokFile = "token.json"
const CredentialsFile = "credentials.json"

var (
	email    = viperEnvVariable("EMAIL")
	password = viperEnvVariable("PASSWORD")
)


func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background())
	defer stop()

	cred := contaja.Credentials{
		Email:    email,
		Password: password,
	}

	c := contaja.NewServer(cred)
	if err := c.HTTPClient.ContajaLogin(); err != nil {
		return err
	}
	log.Println("logged on contaja")

	docs, err := c.HTTPClient.GetFiles()
	if len(docs) < 1 {
		log.Println("New documents not found")
		return nil
	}
	if err != nil {
		return err
	}

	token, err := google2.GetTokenFromFile(tokFile)
	if err != nil {
		config := google2.GetOauthConfig(CredentialsFile)
		config.FetchCode()
		return nil
	}

	client := google2.NewClient(ctx, token, CredentialsFile)
	srv, err := client.NewService(ctx)
	if err != nil {
		return err
	}

	events := contaja.ParseEvents(docs)

	google2.CreateEventFromDocs(srv, events)

	c.Stop()

	return nil
}

func viperEnvVariable(key string) string {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	if value == "" {
		log.Fatalf("Variable %s is empty", key)
	}

	return value
}
