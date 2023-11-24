package google_test

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/rinconrj/golang-scraper/internal/google"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

const testTokenFile = "test_token.json"
const testCredentialsFile = "test_credentials.json"

var mokCredentials = &oauth2.Config{
	ClientID:     "test_client_id",
	ClientSecret: "test_client_secret",
	RedirectURL:  "test_redirect_url",
	Scopes:       []string{"test_scope"},
}

func TestGetTokenFromFile(t *testing.T) {
	token := &oauth2.Token{
		AccessToken: "TestToken",
	}
	generateFile(testTokenFile, token)
	defer os.Remove(testTokenFile)

	resultToken, err := google.GetTokenFromFile(testTokenFile)

	assert.NoError(t, err)
	assert.Equal(t, token.AccessToken, resultToken.AccessToken)
}

func TestSaveToken(t *testing.T) {
	token := &oauth2.Token{
		AccessToken: "TestToken",
	}

	err := google.SaveToken(token)
	assert.NoError(t, err)
	generateFile(testTokenFile, token)
	defer os.Remove(testTokenFile)

	resultToken, err := google.GetTokenFromFile(testTokenFile)

	assert.NoError(t, err)
	assert.Equal(t, token.AccessToken, resultToken.AccessToken)
}

func TestNewClient(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	tok := &oauth2.Token{}

	client := google.NewClient(ctx, tok, testCredentialsFile)

	assert.NotNil(client.Client)
}

func TestNewService(t *testing.T) {
	asserts := assert.New(t)
	ctx := context.Background()

	c := &google.Client{
		Client: &http.Client{},
	}

	srv, err := c.NewService(ctx)

	asserts.Nil(err)
	asserts.NotNil(srv)
}

func TestGetOauthConfig(t *testing.T) {
	asserts := assert.New(t)

	conf := google.GetOauthConfig(testCredentialsFile)

	asserts.NotNil(conf)
	asserts.NotNil(conf.Config)
}

func generateFile(filename string, content any) {
	var file []byte
	file, err := json.MarshalIndent(content, "", " ")
	if err != nil {
		return
	}
	_ = os.WriteFile(filename, file, 0644)
}
