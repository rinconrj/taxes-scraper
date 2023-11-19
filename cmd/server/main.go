package main

// Import necessary libraries
import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/viper"

	google2 "github.com/rinconrj/golang-scraper/internal/google"

	"github.com/rinconrj/golang-scraper/internal/contaja"
)

const tokFile = "token.json"

var (
	email    = viperEnvVariable("EMAIL")
	password = viperEnvVariable("PASSWORD")
)

// Main function runs when the program starts
func main() {
	err := Run()
	log.Fatal(err)
	// Start the server on port 8080
}

func Run() error {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		oscall := <-channel
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	cred := contaja.Credentials{
		Email:    email,
		Password: password,
	}

	c := contaja.NewServer(nil, cred)
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
		config := google2.GetOauthConfig()
		config.FetchCode()
		return nil
	}

	client := google2.NewClient(ctx, token)
	srv, err := client.NewService(ctx)
	if err != nil {
		return err
	}

	events := contaja.ParseEvents(docs)

	google2.CreateEventFromDocs(srv, events)

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
