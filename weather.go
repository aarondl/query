package query

import (
	"fmt"
	"strings"

	"github.com/Islandstone/yr"
)

const (
	weatherURI    = "http://www.yr.no/place/%s/%s/%s/forecast.xml"
	weatherNorURI = "http://www.yr.no/place/%s/%s/%s/%s/forecast.xml"
)

// Weather provides weather information from yr.no
func Weather(query string, conf *Config) (output string, err error) {
	var data *yr.WeatherData
	var URL, city, country, state string

	if subURI, ok := placesLookup[strings.ToLower(query)]; ok {
		country = subURI[0]
		state = subURI[1]
		county := subURI[2]
		city = subURI[3]

		URL = fmt.Sprintf(weatherNorURI, country, state, county, city)
	} else {
		country, state, city, err = getLocation(query, conf)

		if err != nil {
			if e, ok := err.(geoErr); ok {
				return fmt.Sprintf("\x02Weather (\x02YR.no\x02):\x02 %v", e), nil
			}
			return "", err
		}

		URL = fmt.Sprintf(weatherURI, country, state, city)
	}

	data, err = yr.LoadFromURL(URL)

	if err != nil {
		return
	}

	output = fmt.Sprintf(
		"\x02Weather (\x02YR.no\x02):\x02 %s, %s \x02=>\x02 %s, %d \u00B0C",
		city,
		country,
		data.Current().Symbol.Name,
		data.Current().Temperature.Value,
	)

	return
}
