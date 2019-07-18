package query

import (
	"os"
	"testing"
)

func TestYouTube(t *testing.T) {
	t.Parallel()

	key := os.Getenv("YOUTUBE_API_KEY")
	if len(key) == 0 {
		t.Skip("env var YOUTUBE_API_KEY missing")
	}

	output, err := YouTube(`https://www.youtube.com/watch?v=kNcaiTM77cM`, &Config{GoogleYoutubeKey: key})
	if err != nil {
		t.Error(err)
	}

	if len(output) == 0 {
		t.Error("output was not set")
	}
	if output != "\x02YouTube (\x022m51s\x02):\x02 How To Catch Fish in the Sewer" {
		t.Error("output was wrong:", output)
	}
}

func TestYouTubeBE(t *testing.T) {
	t.Parallel()

	key := os.Getenv("YOUTUBE_API_KEY")
	if len(key) == 0 {
		t.Skip("env var YOUTUBE_API_KEY missing")
	}

	output, err := YouTube(`https://youtu.be/-_qpzFlpgpo`, &Config{GoogleYoutubeKey: key})
	if err != nil {
		t.Error(err)
	}

	if len(output) == 0 {
		t.Error("output was not set")
	}
	if output != "\x02YouTube (\x0237m51s\x02):\x02 Making metal crystals from Pepto-Bismol" {
		t.Error("output was wrong:", output)
	}
}
