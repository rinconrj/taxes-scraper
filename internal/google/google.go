package google

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
)

const credentialsFile = "credentials.json"

type Configer struct {
	Config *oauth2.Config
}

type Client struct {
	Client *http.Client
}

func NewClient(ctx context.Context, tok *oauth2.Token) *Client {
	conf := GetOauthConfig()

	return &Client{
		Client: conf.Config.Client(ctx, tok),
	}
}

func CreateEventFromDocs(srv *calendar.Service, events []*calendar.Event) {
	for _, v := range events {
		calendarID := "primary"
		_, err := CreateEvent(srv, calendarID, v)
		if err != nil {
			log.Fatalf("Unable to create event. %v\n", err)
		}
	}

}

func (c Client) NewService(ctx context.Context) (*calendar.Service, error) {
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Println("Unable to retrieve Calendar client: %v", err)
		return nil, err
	}
	return srv, nil
}

func (conf Configer) FetchCode() {
	url := conf.Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.RedirectHandler(url, http.StatusFound)
}

func CreateEvent(service *calendar.Service, calendarID string, event *calendar.Event) (*calendar.Event, error) {
	createdEvent, err := service.Events.Insert(calendarID, event).Do()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Event created: %s\n", createdEvent.HtmlLink)

	return createdEvent, nil
}

func GetOauthConfig() *Configer {
	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		panic(err)
	}
	oauthConfig, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		panic(err)
	}
	return &Configer{
		Config: oauthConfig,
	}
}
