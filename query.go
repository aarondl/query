// Package query provides functions to query web interfaces.
package query

import "github.com/BurntSushi/toml"

// Config is the configuration for this thing.
type Config struct {
	WolframID          string `toml:"wolfram_id"`
	GeonamesID         string `toml:"geonames_id"`
	GoogleURLAPIKey    string `toml:"google_url_api_key"`
	GoogleSearchAPIKey string `toml:"google_search_api_key"`
	GoogleSearchCXID   string `toml:"google_search_cx_id"`
	GithubAPIKey       string `toml:"github_api_key"`
	WundergroundAPIKey string `toml:"wunderground_api_key"`
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
