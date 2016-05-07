package query

import (
	"fmt"
	"testing"
)

func TestGoogleSearchParsing(t *testing.T) {
	cfg := NewConfig("../ultimateq/query.toml")

	fmt.Println(Google("hello", cfg))
}
