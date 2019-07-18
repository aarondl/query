package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const (
	bingURI = "https://api.cognitive.microsoft.com/bing/v7.0/search?%s"
)

type BingError struct {
	Type   string `json:"_type"`
	Errors []struct {
		Code      string `json:"code"`
		SubCode   string `json:"subCode"`
		Message   string `json:"message"`
		Parameter string `json:"parameter"`
	} `json:"errors"`
}

// This struct formats the answers provided by the Bing Web Search API.
type BingAnswer struct {
	Type         string `json:"_type"`
	QueryContext struct {
		OriginalQuery string `json:"originalQuery"`
	} `json:"queryContext"`
	WebPages struct {
		WebSearchURL          string `json:"webSearchUrl"`
		TotalEstimatedMatches int    `json:"totalEstimatedMatches"`
		Value                 []struct {
			ID               string    `json:"id"`
			Name             string    `json:"name"`
			URL              string    `json:"url"`
			IsFamilyFriendly bool      `json:"isFamilyFriendly"`
			DisplayURL       string    `json:"displayUrl"`
			Snippet          string    `json:"snippet"`
			DateLastCrawled  time.Time `json:"dateLastCrawled"`
			SearchTags       []struct {
				Name    string `json:"name"`
				Content string `json:"content"`
			} `json:"searchTags,omitempty"`
			About []struct {
				Name string `json:"name"`
			} `json:"about,omitempty"`
		} `json:"value"`
	} `json:"webPages"`
	Videos struct {
		WebSearchURL          string `json:"webSearchUrl"`
		TotalEstimatedMatches int    `json:"totalEstimatedMatches"`
		Value                 []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			ContentURL  string `json:"contentUrl"`
			HostPageURL string `json:"hostPageUrl"`
			Duration    string `json:"duration"`
		} `json:"value"`
	} `json:"videos"`
	RelatedSearches struct {
		ID    string `json:"id"`
		Value []struct {
			Text         string `json:"text"`
			DisplayText  string `json:"displayText"`
			WebSearchURL string `json:"webSearchUrl"`
		} `json:"value"`
	} `json:"relatedSearches"`
	RankingResponse struct {
		Mainline struct {
			Items []struct {
				AnswerType  string `json:"answerType"`
				ResultIndex int    `json:"resultIndex"`
				Value       struct {
					ID string `json:"id"`
				} `json:"value"`
			} `json:"items"`
		} `json:"mainline"`
		Sidebar struct {
			Items []struct {
				AnswerType string `json:"answerType"`
				Value      struct {
					ID string `json:"id"`
				} `json:"value"`
			} `json:"items"`
		} `json:"sidebar"`
	} `json:"rankingResponse"`
}

// Bing performs a query and returns a formatted result.
func Bing(query string, conf *Config) (output string, err error) {
	if len(conf.BingAPIKey) == 0 {
		return output, errors.New("cannot use bing search without bing_api_key")
	}

	params := make(url.Values)
	params.Set("answerCount", "1")
	params.Set("count", "1")
	params.Set("safeSearch", "Moderate")
	params.Set("q", query)
	u := fmt.Sprintf(bingURI, params.Encode())

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", conf.BingAPIKey)

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		if len(b) == 0 {
			return fmt.Sprintf("\x02Bing: Query returned %d", resp.StatusCode), nil
		}

		var errors BingError
		if err = json.Unmarshal(b, &errors); err != nil {
			return "", err
		}

		spew.Dump(errors)
		return fmt.Sprintf("\x02Bing: Query error %s", errors.Errors[0].Message), nil
	}

	var results BingAnswer
	if err = json.Unmarshal(b, &results); err != nil {
		return "", err
	}

	switch {
	case len(results.WebPages.Value) > 0:
		output := fmt.Sprintf("\x02Bing (\x02%d results\x02):\x02 %s - %s",
			results.WebPages.TotalEstimatedMatches,
			results.WebPages.Value[0].URL,
			results.WebPages.Value[0].Snippet)
		return output, nil
	case len(results.Videos.Value) > 0:
		duration := strings.ToLower(strings.TrimPrefix(results.Videos.Value[0].Duration, "PT"))

		output := fmt.Sprintf("\x02Bing (\x02%s\x02):\x02 %s - %s - %s",
			duration,
			results.Videos.Value[0].ContentURL,
			results.Videos.Value[0].Name,
			results.Videos.Value[0].Description)
		return output, nil
	default:
		return "\x02Bing: No results found.\x02", nil
	}
}
