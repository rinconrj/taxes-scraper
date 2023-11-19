package contaja

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rinconrj/golang-scraper/internal/google"
	"google.golang.org/api/calendar/v3"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

const loginURL = "https://app.contaja.com.br/login"
const query = "https://app.contaja.com.br/tributos-folhas/get-tributos-folhas?draw=2&columns%5B0%5D%5Bdata%5D=documento&columns%5B0%5D%5Bname%5D=documento&columns%5B0%5D%5Bsearchable%5D=true&columns%5B0%5D%5Borderable%5D=true&columns%5B0%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B0%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B1%5D%5Bdata%5D=created_at&columns%5B1%5D%5Bname%5D=created_at&columns%5B1%5D%5Bsearchable%5D=true&columns%5B1%5D%5Borderable%5D=true&columns%5B1%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B1%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B2%5D%5Bdata%5D=competencia&columns%5B2%5D%5Bname%5D=competencia&columns%5B2%5D%5Bsearchable%5D=true&columns%5B2%5D%5Borderable%5D=true&columns%5B2%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B2%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B3%5D%5Bdata%5D=vencimento&columns%5B3%5D%5Bname%5D=vencimento&columns%5B3%5D%5Bsearchable%5D=true&columns%5B3%5D%5Borderable%5D=true&columns%5B3%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B3%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B4%5D%5Bdata%5D=destinatario&columns%5B4%5D%5Bname%5D=destinatario&columns%5B4%5D%5Bsearchable%5D=true&columns%5B4%5D%5Borderable%5D=true&columns%5B4%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B4%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B5%5D%5Bdata%5D=descricao&columns%5B5%5D%5Bname%5D=descricao&columns%5B5%5D%5Bsearchable%5D=true&columns%5B5%5D%5Borderable%5D=true&columns%5B5%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B5%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B6%5D%5Bdata%5D=actions&columns%5B6%5D%5Bname%5D=actions&columns%5B6%5D%5Bsearchable%5D=true&columns%5B6%5D%5Borderable%5D=true&columns%5B6%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B6%5D%5Bsearch%5D%5Bregex%5D=false&order%5B0%5D%5Bcolumn%5D=1&order%5B0%5D%5Bdir%5D=desc&order%5B1%5D%5Bcolumn%5D=2&order%5B1%5D%5Bdir%5D=asc&start=0&length=10&search%5Bvalue%5D=&search%5Bregex%5D=false&competencia=&vencimento=&envio=&status=true&destinatario=%20&_=1700165002794"

type Doc struct {
	ID          int    `json:"id"`
	EmpresaID   int    `json:"empresaID"`
	Documento   string `json:"documento"`
	Code        string `json:"code"`
	Vencimento  string `json:"vencimento"`
	Competencia string `json:"competencia"`
	Descricao   string `json:"descricao"`
	Actions     string `json:"actions"`
}

type Response struct {
	Data            []Doc
	Draw            int
	RecordsFiltered int
	RecordsTotal    int
}

type Credentials struct {
	Email    string
	Password string
}

type Client struct {
	client      *http.Client
	credentials Credentials
}
type Server struct {
	URL        string
	Listener   net.Listener
	closed     bool
	HTTPClient *Client // Embed the Client in the Server struct
	wg         sync.WaitGroup
	mu         sync.Mutex
	Config     *http.Server
}

func NewClient(credentials Credentials) *Client {
	jar, _ := cookiejar.New(nil)
	c := &http.Client{
		Jar: jar,
	}

	return &Client{
		client:      c,
		credentials: credentials,
	}
}

func NewServer(handler http.Handler, credentials Credentials) *Server {
	client := NewClient(credentials)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", client.HandleCallback)

	s := &Server{
		HTTPClient: client,
		Config:     &http.Server{Handler: mux},
	}
	s.Start()

	return s
}

func (s *Server) Start() {
	s.Listener, _ = net.Listen("tcp", ":8080")
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		http.Serve(s.Listener, nil)
	}()
}

func (s *Server) Stop() {
	s.Listener.Close()
	s.wg.Wait()
}

func (c *Client) HandleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	config := google.GetOauthConfig()

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
	docs, err := c.GetFiles()

	events := ParseEvents(docs)

	google.CreateEventFromDocs(srv, events)
}

func (c *Client) GetTokens() (string, string, error) {
	html, err := c.client.Get(loginURL)
	if err != nil {
		panic(err)
	}
	defer func() { _ = html.Body.Close() }()

	cookies := extractCookies(html)
	fmt.Println("Extracted Cookies:", cookies)

	body, err := io.ReadAll(html.Body)
	if err != nil {
		return "", "", err
	}
	csrfToken := extractCSRFToken(string(body))
	log.Println("Extracted csrfToken:", csrfToken)

	return cookies, csrfToken, nil
}

func (c *Client) ContajaLogin() error {
	cookies, csrfToken, err := c.GetTokens()
	fmt.Println("cookies fetched:", cookies, csrfToken)
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Set("email", c.credentials.Email)
	data.Set("password", c.credentials.Password)
	data.Set("_token", csrfToken)

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(data.Encode()))
	req.Header.Add("Cookie", cookies)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusOK {
		if isLogged(string(body)) {
			return nil
		}
		return fmt.Errorf("Login failed")
	}

	return res.Request.Context().Err()
}

func (c *Client) GetFiles() ([]Doc, error) {
	res, err := c.client.Get(query)
	if err != nil {
		fmt.Println("Error to fetch the files:", err)
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	var v Response

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&v)
	if err != nil {
		body, _ := io.ReadAll(res.Body)
		fmt.Println(string(body))
		fmt.Println("Error decoding files:", err)
		return nil, err
	}

	var docs []Doc
	docs = append(docs, v.Data...)

	return docs, err
}

func ParseEvents(docs []Doc) []*calendar.Event {
	var events []*calendar.Event
	for _, v := range docs {
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
		events = append(events, event)
	}
	return events
}

func extractCSRFToken(html string) string {
	r := regexp.MustCompile(`name="_token" value="(.+?)"`)
	matches := r.FindStringSubmatch(html)
	if len(matches) < 2 {
		fmt.Println("CSRF token not found")
		return ""
	}
	return matches[1]
}

func isLogged(html string) bool {
	r := regexp.MustCompile(`name="class" value="m-login__body"`)
	matches := r.FindStringSubmatch(html)
	return len(matches) > 2
}

func extractCookies(resp *http.Response) string {
	cookies := resp.Header["Set-Cookie"]
	return strings.Join(cookies, "; ")
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
