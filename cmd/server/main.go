package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"

	"github.com/rinconrj/golang-scraper/internal/contaja"
	"github.com/rinconrj/golang-scraper/internal/google"
)

var (
	client *http.Client
	docs   []contaja.Doc
)

func init() {
	jar, _ := cookiejar.New(nil)
	client = &http.Client{
		Jar: jar,
	}

}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/callback", handleCallback)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {

	cookies, csrfToken := contaja.GetTokens(client)

	logErr := contaja.ContajaLogin(client, csrfToken, cookies)
	if logErr != nil {
		fmt.Println("login error:", logErr)
	}

	docs, err := contaja.GetFiles(client)
	if err != nil {
		fmt.Println("get files from portal error:", err)
	}

	google.CreateEventFromDocs(docs, w, r)

}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	oauthConfig := google.GetOauthConfig()
	code := r.URL.Query().Get("code")
	fmt.Println("google code:", code)
	tok, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		fmt.Println("error: oauthConfig.Exchange")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	google.SaveToken(tok)
	google.CreateEventFromDocs(docs, w, r)
}
