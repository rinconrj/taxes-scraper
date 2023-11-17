package main

// Import necessary libraries
import (
	"context"
	"github.com/rinconrj/golang-scraper/internal/google"
	"log"
	"os"
	"os/signal"

	"github.com/rinconrj/golang-scraper/internal/contaja"
)

const tokFile = "token.json"

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
		Email:    "***REMOVED***",
		Password: "contaja12s3",
	}

	c := contaja.NewClient(cred)
	if err := c.ContajaLogin(); err != nil {
		return err
	}
	log.Println("logged on contaja")

	docs, err := c.GetFiles()
	if len(docs) < 1 {
		log.Println("New documents not found")
		return nil
	}
	if err != nil {
		return err
	}

	token, err := google.GetTokenFromFile(tokFile)
	if err != nil {
		return nil
	}

	client := google.NewClient(ctx, token)
	srv, err := client.NewService(ctx)
	if err != nil {
		return err
	}

	google.CreateEventFromDocs(srv, docs)

	return nil
}
