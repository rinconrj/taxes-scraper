package google

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rinconrj/golang-scraper/internal/contaja"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const credentialsFile = "credentials.json"

func NewClient(ctx context.Context, oauthConfig *oauth2.Config, token *oauth2.Token) *http.Client {
	return oauthConfig.Client(ctx, token)
}

func CreateEventFromDocs(docs []contaja.Doc, w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	srv, err := NewService(ctx, w, r)
	if err != nil {
		return
	}

	for _, v := range docs {
		fmt.Println("value", v)

		sd, err := timeParser(v.Vencimento, time.RFC3339, 17)
		if err != nil {
			fmt.Println("Error occurred:", err)
		}
		ed, err := timeParser(v.Vencimento, time.RFC3339, 20)
		if err != nil {
			fmt.Println("Error occurred:", err)
		}

		event := &calendar.Event{
			Summary:     fmt.Sprintf("%s %s", v.Descricao, 23, v.Competencia),
			Location:    "",
			Description: v.Actions,
			Start: &calendar.EventDateTime{
				DateTime: sd,
				TimeZone: "America/Sao_Paulo",
			},
			End: &calendar.EventDateTime{
				DateTime: ed,
				TimeZone: "America/Sao_Paulo",
			},
		}

		calendarID := "primary"
		_, err = CreateEvent(srv, calendarID, event)
		if err != nil {
			log.Fatalf("Unable to create event. %v\n", err)
		}

	}

}

func NewService(ctx context.Context, w http.ResponseWriter, r *http.Request) (*calendar.Service, error) {
	oauthConfig := GetOauthConfig()
	token, err := GetTokenFromFile(ctx)
	if err != nil {
		fmt.Println("local token not found, fetching from google")
		url := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
		http.Redirect(w, r, url, http.StatusFound)
		return nil, err
	}
	client := NewClient(ctx, oauthConfig, token)
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
		return nil, err
	}
	return srv, nil
}

func CreateEvent(service *calendar.Service, calendarID string, event *calendar.Event) (*calendar.Event, error) {
	createdEvent, err := service.Events.Insert(calendarID, event).Do()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Event created: %s\n", createdEvent.HtmlLink)
	return createdEvent, nil
}

func GetOauthConfig() *oauth2.Config {
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
	return oauthConfig
}

func timeParser(psdvalue string, layout string, addHours int) (string, error) {
	referenceLayout := "02/01/2006"
	parsed, err := time.Parse(referenceLayout, psdvalue)
	if err != nil {
		return "", err
	}
	parsed = parsed.Add(time.Duration(addHours) * time.Hour)
	finalFormat := parsed.Format(layout)
	return finalFormat, nil
}
