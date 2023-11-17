package proxy

import (
	"context"
	"fmt"
	"github.com/rinconrj/golang-scraper/internal/contaja"
	"github.com/rinconrj/golang-scraper/internal/google"
	"net/http"
)

func Start(docs []contaja.Doc) error {
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		config := google.Configer{Config: google.GetOauthConfig()}
		config.FetchCode()

		code := r.URL.Query().Get("code")
		fmt.Println("google code:", code)

		tok, err := config.Config.Exchange(ctx, code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if err := google.SaveToken(tok); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		client := google.NewClient(ctx, tok)
		srv, err := client.NewService(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		google.CreateEventFromDocs(srv, docs)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return err
	}
	return nil
}
