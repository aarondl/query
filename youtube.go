package query

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	rgxURL = regexp.MustCompile(
		`(?:http|https):\/\/(?:www\.)?(youtube.com|youtu.be)\/\S+`)
)

const (
	apiURLYoutube = `https://www.googleapis.com/youtube/v3/videos?part=snippet,contentDetails&id=%s&key=%s`
)

// YouTube will check to see if a message contains a YouTube uri, if so it will
// format a string with the title in it.
func YouTube(msg string, cfg *Config) (output string, err error) {
	link := rgxURL.FindStringSubmatch(msg)
	if len(link) == 0 {
		// Tell no one
		return "", nil
	}

	uri, err := url.Parse(link[0])
	if err != nil {
		return "", err
	}

	id := uri.Path
	if link[1] == "youtube.com" {
		id = uri.Query().Get("v")
	}

	// Must be an incomplete url
	if len(id) == 0 {
		return "", nil
	}

	apiURL := fmt.Sprintf(apiURLYoutube, url.QueryEscape(id), url.QueryEscape(cfg.GoogleYoutubeKey))
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return "", errors.New("failed to create youtube request")
	}

	client := http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("failed to perform youtube request")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	} else if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("bad status code: %d", resp.StatusCode)
	} else if resp.Body == nil {
		return "", errors.New("no response from api")
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ytResp youtubeListResponse
	if err := json.Unmarshal(b, &ytResp); err != nil {
		return "", err
	}

	// We should only ever have a single item
	item := ytResp.Items[0]

	duration := "Unknown"
	if dur, err := time.ParseDuration(strings.ToLower(strings.TrimPrefix(item.ContentDetails.Duration, "PT"))); err == nil {
		duration = dur.String()
	}

	output = fmt.Sprintf(
		`\x02YouTube (\x02%s\x02):\x02 %s`,
		duration,
		item.Snippet.Title,
	)

	return output, nil
}

type youtubeListResponse struct {
	Kind     string `json:"kind"`
	Etag     string `json:"etag"`
	PageInfo struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind    string `json:"kind"`
		Etag    string `json:"etag"`
		ID      string `json:"id"`
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			ChannelID   string    `json:"channelId"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				Default struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"default"`
				Medium struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"medium"`
				High struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"high"`
				Standard struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"standard"`
				Maxres struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"maxres"`
			} `json:"thumbnails"`
			ChannelTitle         string   `json:"channelTitle"`
			Tags                 []string `json:"tags"`
			CategoryID           string   `json:"categoryId"`
			LiveBroadcastContent string   `json:"liveBroadcastContent"`
			DefaultLanguage      string   `json:"defaultLanguage"`
			Localized            struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"localized"`
			DefaultAudioLanguage string `json:"defaultAudioLanguage"`
		} `json:"snippet"`
		ContentDetails struct {
			Duration        string `json:"duration"`
			Dimension       string `json:"dimension"`
			Definition      string `json:"definition"`
			Caption         string `json:"caption"`
			LicensedContent bool   `json:"licensedContent"`
			Projection      string `json:"projection"`
		} `json:"contentDetails"`
	} `json:"items"`
}
