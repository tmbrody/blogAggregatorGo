package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// Define a struct that represents the structure of an RSS feed
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type CDATA struct {
	Value string `xml:",cdata"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       CDATA  `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
}

// Function to fetch data from an RSS feed URL and return it as a Go struct
func FetchRSSFeed(url string) (RSS, error) {
	response, err := http.Get(url)
	if err != nil {
		return RSS{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return RSS{}, fmt.Errorf("HTTP status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return RSS{}, err
	}

	var rss RSS
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return RSS{}, err
	}

	return rss, nil
}
