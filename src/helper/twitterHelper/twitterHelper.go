package twitterHelper

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Credentials struct {
	ApiToken          string
	ApiTokenSecret    string
	AccessToken       string
	AccessTokenSecret string
}

func GetClient(creds *Credentials) (*twitter.Client, error) {

	config := oauth1.NewConfig(creds.ApiToken, creds.ApiTokenSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	_, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	return client, nil
}
