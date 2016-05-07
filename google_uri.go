package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	shortenURI = "https://www.googleapis.com/urlshortener/v1/url?fields=id,longUrl&key=%s"
)

// URLShortenResponse is the json response back from the server
type URLShortenResponse struct {
	ID      string `json:"id"`
	LongURL string `json:"longUrl"`
}

// URLShortenQuery is the json request to the server
type URLShortenQuery struct {
	LongURL string `json:"longUrl"`
}

// GetShortURL takes a long url and returns a shorter url from the Google API
func GetShortURL(longURL string, conf *Config) (short string, err error) {
	var resp *http.Response
	var body []byte

	urlQuery := URLShortenQuery{LongURL: longURL}
	if body, err = json.Marshal(urlQuery); err != nil {
		return "", fmt.Errorf("failed to marshal url query: %v", err)
	}

	resp, err = http.Post(
		fmt.Sprintf(shortenURI, conf.GoogleURIAPIKey),
		"application/json",
		bytes.NewReader(body),
	)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	var jsonObj URLShortenResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&jsonObj)
	if err != nil {
		return "", err
	}

	return jsonObj.ID, nil
}
