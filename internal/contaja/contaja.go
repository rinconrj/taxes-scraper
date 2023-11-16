package contaja

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// type Doc struct {
// 	Id          int    `json:"id"`
// 	Empresa_id  int    `json:"empresa_id"`
// 	Documento   string `json:"document_name"`
// 	Code        string `json:"code"`
// 	Vencimento  string `json:"expires"`
// 	Competencia string `json:"title"`
// 	Descricao   string `json:"desc"`
// }

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

// type ContajaResponse struct {
// 	Data            []Doc
// 	Draw            int `json:"draw"`
// 	RecordsFiltered int `json:"records_filtered"`
// 	RecordsTotal    int `json:"records_totals"`
// }

type ContajaResponse struct {
	Data            []Doc
	Draw            int
	RecordsFiltered int
	RecordsTotal    int
}

func ContajaLogin(c *http.Client, csrfToken string, cookies string) error {
	data := url.Values{}
	data.Set("email", "***REMOVED***")
	data.Set("password", "***REMOVED***")
	data.Set("_token", csrfToken)

	req, _ := http.NewRequest("POST", "https://app.contaja.com.br/login", strings.NewReader(data.Encode()))
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
	html, err := c.Get("https://app.contaja.com.br/login")
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
	url := "https://app.contaja.com.br/tributos-folhas/get-tributos-folhas?draw=2&columns%5B0%5D%5Bdata%5D=documento&columns%5B0%5D%5Bname%5D=documento&columns%5B0%5D%5Bsearchable%5D=true&columns%5B0%5D%5Borderable%5D=true&columns%5B0%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B0%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B1%5D%5Bdata%5D=created_at&columns%5B1%5D%5Bname%5D=created_at&columns%5B1%5D%5Bsearchable%5D=true&columns%5B1%5D%5Borderable%5D=true&columns%5B1%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B1%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B2%5D%5Bdata%5D=competencia&columns%5B2%5D%5Bname%5D=competencia&columns%5B2%5D%5Bsearchable%5D=true&columns%5B2%5D%5Borderable%5D=true&columns%5B2%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B2%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B3%5D%5Bdata%5D=vencimento&columns%5B3%5D%5Bname%5D=vencimento&columns%5B3%5D%5Bsearchable%5D=true&columns%5B3%5D%5Borderable%5D=true&columns%5B3%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B3%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B4%5D%5Bdata%5D=destinatario&columns%5B4%5D%5Bname%5D=destinatario&columns%5B4%5D%5Bsearchable%5D=true&columns%5B4%5D%5Borderable%5D=true&columns%5B4%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B4%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B5%5D%5Bdata%5D=descricao&columns%5B5%5D%5Bname%5D=descricao&columns%5B5%5D%5Bsearchable%5D=true&columns%5B5%5D%5Borderable%5D=true&columns%5B5%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B5%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B6%5D%5Bdata%5D=actions&columns%5B6%5D%5Bname%5D=actions&columns%5B6%5D%5Bsearchable%5D=true&columns%5B6%5D%5Borderable%5D=true&columns%5B6%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B6%5D%5Bsearch%5D%5Bregex%5D=false&order%5B0%5D%5Bcolumn%5D=1&order%5B0%5D%5Bdir%5D=desc&order%5B1%5D%5Bcolumn%5D=2&order%5B1%5D%5Bdir%5D=asc&start=0&length=10&search%5Bvalue%5D=&search%5Bregex%5D=false&competencia=&vencimento=&envio=&status=true&destinatario=%20&_=1700165002794"
	req, _ := http.NewRequest("GET", url, nil)
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

func parseData(data string) (Doc, error) {
	parts := strings.Split(data, " ")
	if len(parts) != 7 {
		return Doc{}, fmt.Errorf("invalid data")
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return Doc{}, err
	}
	empresa_id, err := strconv.Atoi(parts[1])
	if err != nil {
		return Doc{}, err
	}

	return Doc{
		Id:          id,
		Empresa_id:  empresa_id,
		Documento:   parts[2],
		Code:        parts[3],
		Vencimento:  parts[4],
		Competencia: parts[5],
		Descricao:   parts[6],
	}, nil
}
