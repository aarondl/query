package query

import (
	"bytes"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

var (
	rgxTitle = regexp.MustCompile(`<meta name="title" content="(.*?)"`)
	rgxURL   = regexp.MustCompile(
		`(http|https):\/\/(www\.)?(youtube.com|youtu.be)\/\S+`)
	rgxDur = regexp.MustCompile(
		`<meta itemprop="duration" content="PT([0-9]+M[0-9]+S)"`)
)

// YouTube will check to see if a message contains a YouTube uri, if so it will
// format a string with the title in it.
func YouTube(msg string) (output string, err error) {
	link := rgxURL.FindString(msg)
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
