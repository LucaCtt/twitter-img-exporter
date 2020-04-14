// Package twitter implements a basic Twitter API client.
package twitter

import (
	"fmt"
	"strings"
)

const (
	tokenEndpoint        string = "https://api.twitter.com/oauth2/token"
	userTimelineEndpoint string = "https://api.twitter.com/1.1/statuses/user_timeline.json"
)

// Client is a simple Twitter API client that automatically handles authentication
// and allows to get a user's timeline.
type Client struct {
	token string
}

// NewClient creates a new Twitter client with the given API key and API secret key.
func NewClient(key string, secret string) (*Client, error) {
	token, err := getToken(key, secret)
	if err != nil {
		return nil, fmt.Errorf("Get token failed: %w", err)
	}

	return &Client{token}, nil
}

// getToken makes a token request to the Twitter API.
//
// It requires the API key and the API secret key.
// If successful it will return the token.
func getToken(key string, secret string) (string, error) {
	var result apiAuthResp
	cfg := ReqCfg{
		method:    "POST",
		url:       tokenEndpoint,
		body:      strings.NewReader("grant_type=client_credentials"),
		basicAuth: &BasicAuth{key, secret},
		out:       &result,
	}

	if err := Request(&cfg); err != nil {
		return "", fmt.Errorf("Token request failed: %w", err)
	}

	return result.AccessToken, nil
}

// UserTimeline returns the tweet timeline of a user, identified by their screen name.
//
// The timeline does not include retweets or comments.
func (c *Client) UserTimeline(screenName string) ([]Tweet, error) {
	var result []Tweet
	cfg := ReqCfg{
		method: "GET",
		url:    fmt.Sprintf("%s?screen_name=%s", userTimelineEndpoint, screenName),
		token:  c.token,
		out:    &result,
	}

	if err := Request(&cfg); err != nil {
		return nil, fmt.Errorf("Timeline request failed: %w", err)
	}

	return result, nil
}
