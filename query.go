// Provides functions to query web interfaces.
package query

import (
	"encoding/json"
	"encoding/xml"
	"net/url"
	"bytes"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
	"github.com/BurntSushi/toml"
	"time"
)

// GoogleData is used to parse the response from Google.
type GoogleData struct {
	ResponseData *ResponseData
	ResponseStatus float64
}

// ResponseData is a substruct of GoogleData.
type ResponseData struct {
	Results []*Result
	Cursor  *Cursor
}

// Result is a substruct of GoogleData.
type Result struct {
	GsearchResultClass string
	UnescapedUrl       string
	Url                string
	VisibleUrl         string
	CacheUrl           string
	Title              string
	TitleNoFormatting  string
	Content            string
}

// Cursor is a substruct of GoogleData.
type Cursor struct {
	ResultCount string
}

// WolframData is used to parse the response from WolframAlpha.
type WolframData struct {
	XMLName     xml.Name `xml:"queryresult"`
	Success     bool     `xml:"success,attr"`
	ParseTiming float64  `xml:"parsetiming,attr"`
	Numpods     int      `xml:"numpods,attr"`
	Pods        []*Pod   `xml:"pod"`
	DidYouMeans []string `xml:"didyoumeans>didyoumean"`
}

// Pod is a substruct of WolframData.
type Pod struct {
	Title      string    `xml:"title,attr"`
	Id         string    `xml:"id,attr"`
	Primary    bool      `xml:"primary,attr"`
	Numsubpods int       `xml:"numsubpods,attr"`
	PlainTexts []string  `xml:"subpod>plaintext"`
}

// Config is the configuration for this thing.
type Config struct {
	WolframId string
}

// NewConfig loads the config file.
func NewConfig(file string) *Config {
	var conf Config
	_, err := toml.DecodeFile(file, &conf)
	if err != nil {
		return nil
	}
	return &conf
}

var (
	rgxTitle = regexp.MustCompile(`<meta name="title" content="(.*?)"`)
	rgxUrl   = regexp.MustCompile(
		`(http|https):\/\/(www\.)?(youtube.com|youtu.be)\/\S+`)
	rgxDur = regexp.MustCompile(
		`<meta itemprop="duration" content="PT([0-9]+M[0-9]+S)"`)
	rgxTags = regexp.MustCompile(`<[^>]*>`)
)

const (
	googleUri = "http://ajax.googleapis.com/ajax/services/search/web?v=1.0&q=%s"
	wolframUri = "http://api.wolframalpha.com/v2/query?format=plaintext&input=%s&appid=%s"
)

// Google performs a query and returns a formatted result.
func Google(query string) (output string, err error) {
	var resp *http.Response
	resp, err = http.Get(fmt.Sprintf(googleUri, url.QueryEscape(query)))
	if err != nil {
		return
	}

	defer resp.Body.Close()
	var jsonObj GoogleData
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&jsonObj)
	if err != nil {
		return
	}

	if jsonObj.ResponseData != nil && int(jsonObj.ResponseStatus) == 200 {
		if len(jsonObj.ResponseData.Results) == 0 {
			output = fmt.Sprintf("\x02Google: No results found.")
			return
		}

		result := jsonObj.ResponseData.Results[0]
		output = fmt.Sprintf(
			"\x02Google (\x02%s results\x02):\x02 %s - %s",
			jsonObj.ResponseData.Cursor.ResultCount,
			rgxTags.ReplaceAllString(result.UnescapedUrl, ""),
			html.UnescapeString(rgxTags.ReplaceAllString(result.Content, "")),
		)
	}

	return
}

// Wolfram performs a query and returns a formatted result.
func Wolfram(query string, conf *Config) (output string, err error) {
	var resp *http.Response
	requestUri := fmt.Sprintf(wolframUri, url.QueryEscape(query), conf.WolframId)
	resp, err = http.Get(requestUri)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	var xmlObj WolframData
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&xmlObj)
	if err != nil {
		return
	}

	// Handle cases of no results.
	if !xmlObj.Success {
		if len(xmlObj.DidYouMeans) > 0 {
			output = fmt.Sprintf("\x02Wolfram (\x02%.2fms\x02):\x02 Did you mean: %s",
				xmlObj.ParseTiming,
				xmlObj.DidYouMeans[0],
			)
		} else {
			output = fmt.Sprintf("\x02Wolfram (\x02%.2fms\x02):\x02 No results found.",
				xmlObj.ParseTiming)
		}
		return
	}

	// If there was no primary response fallback to link.
	if len(xmlObj.Pods) < 2 || len(xmlObj.Pods[1].PlainTexts[0]) == 0 {
		output = fmt.Sprintf(
			"\x02Wolfram (\x02%.2fms\x02):\x02 %s \x02=>\x02 http://www.wolframalpha.com/input/?i=%s",
			xmlObj.ParseTiming,
			xmlObj.Pods[0].PlainTexts[0],
			url.QueryEscape(query),
		)
		return
	}

	output = fmt.Sprintf(
		"\x02Wolfram (\x02%.2fms\x02):\x02 %s \x02=>\x02 %s",
		xmlObj.ParseTiming,
		xmlObj.Pods[0].PlainTexts[0],
		xmlObj.Pods[1].PlainTexts[0],
	)

	return
}

// YouTube will check to see if a message contains a YouTube uri, if so it will
// format a string with the title in it.
func YouTube(msg string) (output string, err error) {
	link := rgxUrl.FindString(msg)
	if len(link) == 0 {
		return
	}

	var resp *http.Response
	var body []byte

	resp, err = http.Get(link)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var buf = bytes.NewBufferString("\x02YouTube")
	if duration := rgxDur.FindSubmatch(body); duration != nil {
		dur, e := time.ParseDuration(string(bytes.ToLower(duration[1])))
		if e == nil {
			buf.WriteString(" (\x02")
			buf.WriteString(dur.String())
			buf.WriteString("\x02)")
		}
	}

	if title := rgxTitle.FindSubmatch(body); title != nil {
		buf.WriteString(":\x02 ")
		buf.WriteString(html.UnescapeString(string(title[1])))
		output = buf.String()
	}

	return
}
