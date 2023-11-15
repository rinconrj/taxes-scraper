package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Replace with your OAuth 2.0 Client ID credentials
const credentialsFile = "/credentials.json"

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

// getTokenFromWeb requests a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
			log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
			log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// saveToken saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
			log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getClient(config *oauth2.Config) *http.Client {
    // The file token.json stores the user's access and refresh tokens.
    tokFile := "token.json"
    tok, err := tokenFromFile(tokFile)
    if err != nil {
        tok = getTokenFromWeb(config)
        saveToken(tokFile, tok)
    }
    return config.Client(context.Background(), tok)
}

// Function to create a new event
func CreateEvent(service *calendar.Service, calendarID string, event *calendar.Event) (*calendar.Event, error) {
    createdEvent, err := service.Events.Insert(calendarID, event).Do()
    if err != nil {
        return nil, err
    }
    fmt.Printf("Event created: %s\n", createdEvent.HtmlLink)
    return createdEvent, nil
}

func getGoogleService() (*calendar.Service, error) {
	cxt := context.Background()
	b, err := os.ReadFile(credentialsFile)
	if err != nil {
			log.Fatalf("Unable to read client secret file: %v", err)
			return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
			return nil, err
	}
	client := getClient(config)

	srv, err := calendar.NewService(cxt, option.WithHTTPClient(client))
	if err != nil {
			log.Fatalf("Unable to retrieve Calendar client: %v", err)
			return nil, err
	}

	return srv, nil
}

func createEventFromDocs(doc []Doc) {

	srv, err := getGoogleService()
	if err != nil {
		return
	}
	//TODO

    // Create an event
    event := &calendar.Event{
        Summary:     "Google I/O 2023",
        Location:    "800 Howard St., San Francisco, CA 94103",
        Description: "A chance to hear more about Google's developer products.",
        Start: &calendar.EventDateTime{
            DateTime: "2023-05-28T09:00:00-07:00",
            TimeZone: "America/Los_Angeles",
        },
        End: &calendar.EventDateTime{
            DateTime: "2023-05-28T17:00:00-07:00",
            TimeZone: "America/Los_Angeles",
        },
    }

    calendarID := "primary"
    _, err = CreateEvent(srv, calendarID, event)
    if err != nil {
        log.Fatalf("Unable to create event. %v\n", err)
    }
}
