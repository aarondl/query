package query

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

const (
	googleURI = "https://www.googleapis.com/customsearch/v1?%s"
)

var (
	rgxTags = regexp.MustCompile(`<[^>]*>`)
)

// GoogleData is used to parse the response from Google.
type GoogleSearch struct {
	Items   []GoogleSearchItem      `json:"items"`
	Info    GoogleSearchInformation `json:"searchInformation"`
	Queries GoogleQueries           `json:"queries"`
}

// GoogleSearchItem is a search result item from a google search
type GoogleSearchItem struct {
	Title        string `json:"title"`
	Snippet      string `json:"snippet"`
	Link         string `json:"link"`
	DisplayLink  string `json:"displayLink"`
	FormattedURL string `json:"formattedUrl"`

	HTMLTitle        string `json:"htmlTitle"`
	HTMLSnippet      string `json:"htmlSnippet"`
	HTMLFormattedURL string `json:"htmlFormattedUrl"`

	CacheID string `json:"cacheId"`
	Kind    string `json:"kind"`
	// pagemap ignored
}

// GoogleSearchInformation
type GoogleSearchInformation struct {
	TotalResults          int64   `json:"integer"`
	SearchTime            float64 `json:"searchTime"`
	FormattedTotalResults string  `json:"formattedTotalResults"`
	FormattedSearchTime   string  `json:"formattedSearchTime"`
}

// GoogleQueries a set of queries involved in the current search query
type GoogleQueries struct {
	NextPage []GoogleQuery `json:"nextPage"`
	Request  []GoogleQuery `json:"request"`
}

// GoogleQuery is a description of a search query
type GoogleQuery struct {
	CX             string `json:"cx"`
	Title          string `json:"title"`
	TotalResults   string `json:"totalResults"`
	SearchTerms    string `json:searchTerms`
	Count          int    `json:"count"`
	StartIndex     int    `json:"startIndex"`
	InputEncoding  string `json:"inputEncoding"`
	OutputEncoding string `json:"outputEncoding"`
	Safe           string `json:"safe"`
}

// Google performs a query and returns a formatted result.
func Google(query string, conf *Config) (output string, err error) {
	params := make(url.Values)
	params.Set("cx", conf.GoogleSearchCXID)
	params.Set("key", conf.GoogleSearchAPIKey)
	params.Set("q", query)
	params.Set("num", "1")
	u := fmt.Sprintf(googleURI, params.Encode())

	fmt.Println(u)

	var resp *http.Response
	resp, err = http.Get(u)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("\x02Google: Query returned %d", resp.StatusCode), nil
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var results GoogleSearch
	if err = json.Unmarshal(b, &results); err != nil {
		return "", err
	}

	if len(results.Items) == 0 {
		return "\x02Google: No results found.\x02", nil
	}

	output = fmt.Sprintf(
		"\x02Google (\x02%d results\x02):\x02 %s - %s",
		results.Info.TotalResults,
		results.Items[0].Link,
		results.Items[0].Snippet,
	)

	return output, nil
}
