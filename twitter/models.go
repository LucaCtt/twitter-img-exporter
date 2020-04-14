package twitter

import (
	"fmt"
	"strings"
)

// apiAuthResp represents the response returned by
//the Twitter API to a successful token request.
type apiAuthResp struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

// apiErr represents the error returned by the Twitter API.
type apiErr struct {
	Errors []struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"errors"`
}

// err builds a golang error from an API error.
func (err *apiErr) err() error {
	var msgs strings.Builder

	for i := 0; i < len(err.Errors); i++ {
		msgs.WriteString(fmt.Sprintf("Code: %d, Message: %s; ", err.Errors[i].Code, err.Errors[i].Message))
	}

	return fmt.Errorf(msgs.String())
}

// Tweet represents a minimal tweet structure that
// contains only the bare minimum data required by the project.
type Tweet struct {
	Entities struct {
		Media []struct {
			URL string `json:"media_url_https"`
		} `json:"media"`
	} `json:"entities"`
}

// MediaURLs returns a slice containing the URLs of all the media in the tweet.
func (t *Tweet) MediaURLs() []string {
	var result []string
	for i := 0; i < len(t.Entities.Media); i++ {
		result = append(result, t.Entities.Media[i].URL)
	}

	return result
}
