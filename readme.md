---

# GoLang Scraper

## Description
This Go server is designed to interact with the Contaja and Google services. It logs into Contaja to fetch documents and then creates Google Calendar events based on these documents. The server handles OAuth2 callbacks for Google authentication.

## Installation

### Prerequisites
- Go installed on your system

- Contaja credential in the .env file:
```
EMAIL='email@email.com'
PASSWORD='secretpass'
```

- Set Up Google Calendar API:

1. Go to the Google Developers Console.
2. Create a new project or select an existing one.
3. Enable the Google Calendar API for your project.
4. Create credentials (OAuth client ID) for your application. Download the JSON file containing these credentials.
5. And save it in the root folder as "credentials.json"

### Setup
Clone the repository:
```bash
git clone https://github.com/rinconrj/golang-scraper.git
cd golang-scraper
```

## Usage

Run the server:
```bash
make run
```
The server will start on `localhost:8080`. It has a main endpoint:

- `/callback`: The OAuth2 callback endpoint for Google. It handles the authentication flow and token exchange.

## Endpoint

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