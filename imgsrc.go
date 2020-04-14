package main

import (
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/lucactt/twitter-img-exporter/twitter"
)

// Img represents an image.
// Since its content is a ReadCloser, it should be closed when not needed anymore.
type Img struct {
	Name    string
	Content io.ReadCloser
}

// ImgSrc represents a source from which images can be read.
type ImgSrc interface {
	Read() ([]Img, error)
}

// TwitterImgSrc is an ImgSrc that reads images from the tweets in an user's timeline.
// The user is identified by their screen name.
type TwitterImgSrc struct {
	Client     *twitter.Client
	ScreenName string
}

// Read extracts the images embedded in the tweets in an users timeline.
func (t *TwitterImgSrc) Read() ([]Img, error) {
	tweets, err := t.Client.UserTimeline(t.ScreenName)
	if err != nil {
		return nil, fmt.Errorf("Get timeline for user \"%s\" failed: %w", t.ScreenName, err)
	}

	var imgs []Img

	for i := 0; i < len(tweets); i++ {
		media, err := images(&tweets[i])
		if err != nil {
			return nil, fmt.Errorf("Get media of tweet failed: %w", err)
		}

		imgs = append(imgs, media...)
	}

	return imgs, nil
}

func images(t *twitter.Tweet) ([]Img, error) {
	var result []Img

	for i := 0; i < len(t.MediaURLs()); i++ {
		url := t.MediaURLs()[i]

		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("Get tweet image \"%s\" failed: %w", url, err)
		}

		result = append(result, Img{path.Base(url), resp.Body})
	}

	return result, nil
}
