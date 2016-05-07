package query

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
)

const (
	wolframURI = "http://api.wolframalpha.com/v2/query?format=plaintext&input=%s&appid=%s"
)

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
	Title      string   `xml:"title,attr"`
	ID         string   `xml:"id,attr"`
	Primary    bool     `xml:"primary,attr"`
	Numsubpods int      `xml:"numsubpods,attr"`
	PlainTexts []string `xml:"subpod>plaintext"`
}

// Wolfram performs a query and returns a formatted result.
func Wolfram(query string, conf *Config) (output string, err error) {
	var resp *http.Response
	requestURI := fmt.Sprintf(wolframURI, url.QueryEscape(query), conf.WolframID)
	resp, err = http.Get(requestURI)
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
