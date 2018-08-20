package query

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func GithubStars(user string, conf *Config) (output string, err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: conf.GithubAPIKey},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	var repos []*github.Repository
	opts := &github.RepositoryListOptions{}
	opts.Type = "public"
	opts.PerPage = 50

	for {
		pagedRepos, resp, err := client.Repositories.List(ctx, user, opts)
		if err != nil {
			return "", err
		}

		if len(pagedRepos) == 0 {
			if opts.Page == 0 {
				return "\x02Github Stars: No results found\x02", nil
			}
			break
		}

		repos = append(repos, pagedRepos...)

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	count := 0
	for _, r := range repos {
		if r.StargazersCount != nil {
			count += *r.StargazersCount
		}
	}

	return fmt.Sprintf("\x02Github Stars:\x02 %d", count), nil
}
