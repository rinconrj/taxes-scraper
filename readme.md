---

# GoLang Scraper

## Description
This Go server is designed to interact with the Contaja and Google services. It logs into Contaja to fetch documents and then creates Google Calendar events based on these documents. The server handles OAuth2 callbacks for Google authentication.

## Installation

### Prerequisites
- Go installed on your system
- Contaja and Google API credentials set up

### Setup
Clone the repository and build the server:
```bash
git clone https://github.com/rinconrj/golang-scraper.git
cd golang-scraper
go build
```

## Usage

Run the server:
```bash
./golang-scraper
```
The server will start on `localhost:8080`. It has two main endpoints:

- `/`: The main site endpoint. It logs into Contaja, fetches documents, and creates Google Calendar events based on these documents.
- `/callback`: The OAuth2 callback endpoint for Google. It handles the authentication flow and token exchange.

## Endpoints

### `GET /`
This endpoint performs several actions:
1. Retrieves authentication tokens from Contaja.
2. Logs into Contaja using the fetched tokens.
3. Fetches documents from Contaja.
4. Creates Google Calendar events based on the fetched documents.

### `GET /callback`
This endpoint handles the OAuth2 callback for Google. It performs the following:
1. Receives an authentication code from Google.
2. Exchanges the code for a token.
3. Saves the received token for further use.
4. Creates Google Calendar events based on documents previously fetched from Contaja.

## Error Handling
Errors at any step (login, file fetching, OAuth2 exchange) are logged to the console. An internal server error is returned in case of OAuth2 exchange failure.

## Contributing
Contributions to improve the server or extend its functionality are welcome. Please follow the standard Git workflow - fork the repository, make changes, and submit a pull request.

---