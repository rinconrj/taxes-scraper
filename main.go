package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)


type Doc struct {
	Id int;
	Empresa_id int;
	Documento string;
	Code string;
	Vencimento string;
	Competencia string;
	Descricao string;
}

type Response struct {
	Data []Doc;
	Draw int;
	RecordsFiltered int;
	RecordsTotal int;
}

func main() {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	// Step 1: Get the login page to fetch the CSRF token
	resp, err := client.Get("https://app.contaja.com.br/login")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	cookies := extractCookies(resp)
	fmt.Println("Extracted Cookies:", cookies)

	body, _ := io.ReadAll(resp.Body)
	csrfToken := extractCSRFToken(string(body))

	// Step 2: Log in using the credentials and CSRF token
	data := url.Values{}
	data.Set("email", "***REMOVED***")
	data.Set("password", "***REMOVED***")
	data.Set("_token", csrfToken)

	req, _ := http.NewRequest("POST", "https://app.contaja.com.br/login", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", cookies)

	// Perform the login request
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Check login success and proceed
	if resp.StatusCode == http.StatusOK {
		// getFiles()
		if err != nil {
			fmt.Println("last response error",err)
		}

		docs, err := getFiles()
		if err !=nil {
			return
		}
		createEventFromDocs(docs)
	}

}


func getFiles() ([]Doc, error) {
	url := "https://app.contaja.com.br/tributos-folhas/get-tributos-folhas?draw=3&columns%5B0%5D%5Bdata%5D=documento&columns%5B0%5D%5Bname%5D=documento&columns%5B0%5D%5Bsearchable%5D=true&columns%5B0%5D%5Borderable%5D=true&columns%5B0%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B0%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B1%5D%5Bdata%5D=created_at&columns%5B1%5D%5Bname%5D=created_at&columns%5B1%5D%5Bsearchable%5D=true&columns%5B1%5D%5Borderable%5D=true&columns%5B1%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B1%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B2%5D%5Bdata%5D=competencia&columns%5B2%5D%5Bname%5D=competencia&columns%5B2%5D%5Bsearchable%5D=true&columns%5B2%5D%5Borderable%5D=true&columns%5B2%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B2%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B3%5D%5Bdata%5D=vencimento&columns%5B3%5D%5Bname%5D=vencimento&columns%5B3%5D%5Bsearchable%5D=true&columns%5B3%5D%5Borderable%5D=true&columns%5B3%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B3%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B4%5D%5Bdata%5D=destinatario&columns%5B4%5D%5Bname%5D=destinatario&columns%5B4%5D%5Bsearchable%5D=true&columns%5B4%5D%5Borderable%5D=true&columns%5B4%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B4%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B5%5D%5Bdata%5D=descricao&columns%5B5%5D%5Bname%5D=descricao&columns%5B5%5D%5Bsearchable%5D=true&columns%5B5%5D%5Borderable%5D=true&columns%5B5%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B5%5D%5Bsearch%5D%5Bregex%5D=false&columns%5B6%5D%5Bdata%5D=actions&columns%5B6%5D%5Bname%5D=actions&columns%5B6%5D%5Bsearchable%5D=true&columns%5B6%5D%5Borderable%5D=true&columns%5B6%5D%5Bsearch%5D%5Bvalue%5D=&columns%5B6%5D%5Bsearch%5D%5Bregex%5D=false&order%5B0%5D%5Bcolumn%5D=1&order%5B0%5D%5Bdir%5D=desc&order%5B1%5D%5Bcolumn%5D=2&order%5B1%5D%5Bdir%5D=asc&start=0&length=10&search%5Bvalue%5D=&search%5Bregex%5D=false&competencia=&vencimento=&envio=&destinatario=%20&_=1700066512820"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Cookie", "remember_web_3dc7a913ef5fd4b890ecabe3487085573e16cf82=eyJpdiI6InFuRmtBWmI0V01VVkVRaXRnN0h6TkE9PSIsInZhbHVlIjoibldIZzBsblI2R2dTQ3BPTWNYUUc2dEs0TUN4YTVUQk5BSnF3WUI4SGNHXC90VzJzNTR2Um9lMnIrNEsrVXR5SzRzT3BcL0VzbWwzMWRaS2RYalIwOVl6Y05GOXRHK0pZcWR3dWRVaEU4Q1pBTjU1cG9qMCtNZndwTFRZNFRyellzMEgyNWpCNzh3aWdiYUNMNUN2WDQ2Mm9ZZXl6bHp2RVVudGpWQlRJMEl4XC80PSIsIm1hYyI6ImM2NjM3Nzg3MjYwMzA2ZGEzMTFmZTlhZGE0YjNhMzhiOGQwMjUxZTMzNzY5ZTFhOTAwMzQ1ODJmOTY3ZWE3ZjcifQ%3D%3D; XSRF-TOKEN=eyJpdiI6IlZTY3hyMEJJU3FhaWZoMGRzV01vbmc9PSIsInZhbHVlIjoiUzcxdHFLRWNlTmttNUlxdjYzc0o0V1FWZHBYcEVDdkw4d3B6WUZ5b052RXQxVnFcLzhRcStETnZBc0t5Z2dHcHciLCJtYWMiOiI3NDRiZjA1YmUzMDg0ZWExMmI5NTIzNThkOTNhMDgyYzkyY2NmYWU2YmExN2Y2YzkxMmVlMzAxMjIyMjE2YjkxIn0%3D; contaja_session=eyJpdiI6IkhiZlcrMzVGaU92WHdmN1NrNWVwZ2c9PSIsInZhbHVlIjoiQ09lRXJcL1hKV1wvb29oenNaeFFINitXQjBqNE8xSFwvZ3g3MnBaUXpCS2tWQ2RQaWVoSFdxc05BdnhiWXVMXC9ER1EiLCJtYWMiOiI4NTNiYzZkNjVlNmE1NWIwZWMwN2VmMmJhZjcxNmM5ODg1ZmJjNTU5MmM3NGU0YjY5NTk2MmU3MGI2YTBjNGQwIn0%3D")
	req.Header.Add("User-Agent", "insomnia/2023.5.8")
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error to fetch the files:", err)
		return nil, err
	}

	defer res.Body.Close()

	var v Response

	dec := json.NewDecoder(res.Body)

	err = dec.Decode(&v)
	if err != nil {
		fmt.Println("error:",err)
		return nil, err
	}

	return v.Data, err
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