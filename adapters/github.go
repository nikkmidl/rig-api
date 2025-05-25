package adapters

import (
	"context"

	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	api *github.Client
}

func NewGHClient(token string) *GithubClient {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return &GithubClient{api: github.NewClient(tc)}
}

func (c *GithubClient) ListReposWithCollaborators(ctx context.Context, org string) (map[string][]string, error) {
	repos := make(map[string][]string)

	// Prapare options for requesting repositories
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	// Fetch repositories for the organization
	repoList, _, err := c.api.Repositories.ListByOrg(ctx, org, opt)
	if err != nil {
		return nil, err
	}

	// Iterate through each repository and fetch collaborators
	for _, repo := range repoList {
		collabs, _, err := c.api.Repositories.ListCollaborators(ctx, org, *repo.Name, nil)
		if err != nil {
			continue
		}
		var usernames []string
		for _, u := range collabs {
			usernames = append(usernames, *u.Login)
		}
		repos[*repo.Name] = usernames
	}

	return repos, nil
}
