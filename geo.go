package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	geoErrMsg = "Unable to find %s"
	geoURI    = "http://api.geonames.org/search?username=%s&q=%s&maxRows=1&type=json&orderby=relevance"
)

var placesLookup = map[string][]string{
	"oslo":     {"Norway", "Oslo", "Oslo", "Oslo"},
	"sandvika": {"Norway", "Akershus", "Bærum", "Sandvika"},
}

type geonameplace struct {
	CountryName string
	AdminName1  string
	Name        string
}

type geonamesdata struct {
	Geonames []geonameplace
}

type geoErr struct {
	query string
}

func (g geoErr) Error() string {
	return fmt.Sprintf(geoErrMsg, g.query)
}

func getLocation(query string, conf *Config) (country, state, city string, err error) {
	if len(conf.GeonamesID) == 0 {
		return country, state, city, errors.New("geo cannot be used without geonames_id in config")
	}

	resp, err := http.Get(fmt.Sprintf(geoURI, conf.GeonamesID, url.QueryEscape(query)))

	if err != nil {
		return
	}

	defer resp.Body.Close()

	r := json.NewDecoder(resp.Body)

	var data geonamesdata
	err = r.Decode(&data)

	if err != nil {
		return
	}

	if len(data.Geonames) == 0 {
		err = geoErr{query}
		return
	}

	country = data.Geonames[0].CountryName
	state = data.Geonames[0].AdminName1
	city = data.Geonames[0].Name

	return
}
