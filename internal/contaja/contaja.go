package contaja

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var user = os.Getenv("CONTAJA_USER")
var pass = os.Getenv("CONTAJA_PASSWORD")
var login_url = os.Getenv("CONTAJA_LOGIN")
var query = os.Getenv("QUERY")

type Doc struct {
	Id          int    `json:"id"`
	Empresa_id  int    `json:"empresa_id"`
	Documento   string `json:"documento"`
	Code        string `json:"code"`
	Vencimento  string `json:"vencimento"`
	Competencia string `json:"competencia"`
	Descricao   string `json:"descricao"`
	Actions     string `json:"actions"`
}

type ContajaResponse struct {
	Data            []Doc
	Draw            int
	RecordsFiltered int
	RecordsTotal    int
}

func ContajaLogin(c *http.Client, csrfToken string, cookies string) error {
	data := url.Values{}
	data.Set("email", user)
	data.Set("password", pass)
	data.Set("_token", csrfToken)

	req, _ := http.NewRequest("POST", login_url, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", cookies)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return resp.Request.Context().Err()
}

func GetTokens(c *http.Client) (string, string) {
	html, err := c.Get(login_url)
	if err != nil {
		panic(err)
	}
	defer html.Body.Close()

	cookies := extractCookies(html)
	fmt.Println("Extracted Cookies:", cookies)

	body, _ := io.ReadAll(html.Body)
	csrfToken := extractCSRFToken(string(body))
	fmt.Println("Extracted csrfToken:", csrfToken)

	return cookies, csrfToken
}

func GetFiles(c *http.Client) ([]Doc, error) {
	req, _ := http.NewRequest("GET", query, nil)
	req.Header.Add("User-Agent", "insomnia/2023.5.8")
	req.Header.Add("Content-Type", "application/json")

	res, err := c.Do(req)
	if err != nil {
		fmt.Println("Error to fetch the files:", err)
		return nil, err
	}
	defer res.Body.Close()

	var v ContajaResponse

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&v)
	if err != nil {
		fmt.Println("Error decoding files:", err)
		return nil, err
	}
	docs := []Doc{}
	docs = append(docs, v.Data...)

	return docs, err
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

func extractCookies(resp *http.Response) string {
	cookies := resp.Header["Set-Cookie"]
	return strings.Join(cookies, "; ")
}
