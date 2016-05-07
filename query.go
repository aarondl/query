// Package query provides functions to query web interfaces.
package query

import "github.com/BurntSushi/toml"

// Config is the configuration for this thing.
type Config struct {
	WolframID          string `toml:"wolfram_id"`
	GeonamesID         string `toml:"geonames_id"`
	GoogleURIAPIKey    string `toml:"google_uri_api_key"`
	GoogleSearchAPIKey string `toml:"google_search_api_key"`
	GoogleSearchCXID   string `toml:"google_search_cx_id"`
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
