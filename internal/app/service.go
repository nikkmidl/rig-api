package app

import (
	"context"

	"github.com/nikkmidl/rig-api/adapters"
	"github.com/nikkmidl/rig-api/internal/domain"
	"github.com/nikkmidl/rig-api/pkg/opa"
)

type Service struct {
	gh  *adapters.GithubClient
	opa *opa.Evaluator
}

func New(gh *adapters.GithubClient, opa *opa.Evaluator) *Service {
	return &Service{gh: gh, opa: opa}
}

func (s *Service) GetAccessInfo(ctx context.Context, org string) ([]domain.RepoAccess, error) {
	data, err := s.gh.ListReposWithCollaborators(ctx, org)
	if err != nil {
		return nil, err
	}

	var results []domain.RepoAccess
	for repo, users := range data {
		var allowed, blocked []string
		// Check each user against the OPA policy
		for _, u := range users {
			err, ok := s.opa.IsBlocked(u)
			if err != nil {
				return nil, err
			}
			if ok {
				blocked = append(blocked, u)
			} else {
				allowed = append(allowed, u)
			}
		}
		// Append the results for this repository
		results = append(results, domain.RepoAccess{
			RepoName:     repo,
			AllowedUsers: allowed,
			BlockedUsers: blocked,
		})
	}
	return results, nil
}
