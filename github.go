package query

import (
	"context"
	"errors"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GithubStars takes a user (aarondl) or repo (aarondl/query) and returns
// the number of stars.
func GithubStars(userOrRepo string, conf *Config) (count int, err error) {
	if len(userOrRepo) == 0 {
		return 0, errors.New("must supply a userOrRepo")
	}

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

	splits := strings.SplitN(userOrRepo, "/", 2)
	user := splits[0]

	if len(splits) == 2 {
		repoName := splits[1]

		repo, _, err := client.Repositories.Get(ctx, user, repoName)
		if err != nil {
			return 0, err
		}

		if repo.StargazersCount != nil {
			count += *repo.StargazersCount
		}

		return count, nil
	}

	for {
		pagedRepos, resp, err := client.Repositories.List(ctx, userOrRepo, opts)
		if err != nil {
			return 0, err
		}

		if len(pagedRepos) == 0 {
			if opts.Page == 0 {
				return 0, nil
			}
			break
		}

		repos = append(repos, pagedRepos...)

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	for _, r := range repos {
		if r.StargazersCount != nil {
			count += *r.StargazersCount
		}
	}

	return count, nil
}
