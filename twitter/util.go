package twitter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// BasicAuth contains the parameters required for HTTP basic authentication.
type BasicAuth struct {
	username string
	password string
}

// ReqCfg contains the parameters required to perform a request.
type ReqCfg struct {
	method    string
	url       string
	body      io.Reader
	basicAuth *BasicAuth
	token     string
	out       interface{}
}

// Request performs an HTTP request specificed using a reqCfg.
//
// This method supports both basic HTTP and Bearer token authentication.
// If the token is not empty, it will be used over basic HTTP auth.
func Request(cfg *ReqCfg) error {
	req, err := http.NewRequest(cfg.method, cfg.url, cfg.body)
	if err != nil {
		// If err is not nil it must be a developer error, so panic is ok
		panic(err)
	}

	// If the token is not specified, resort to basic auth
	if cfg.token == "" {
		// keys and secret must be escaped.
		req.SetBasicAuth(url.QueryEscape(cfg.basicAuth.username), url.QueryEscape(cfg.basicAuth.password))
	} else {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cfg.token))
	}

	// The charset here is very important: if omitted, the API will return an error
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != 200 {
		var result apiErr
		if decoder.Decode(&result) != nil {
			panic(err)
		}

		return result.err()
	}

	if decoder.Decode(&cfg.out) != nil {
		panic(err)
	}

	return nil
}
