package handler

import (
	"context"

	"github.com/nikkmidl/rig-api/internal/app"
	proto "github.com/nikkmidl/rig-api/proto"
)

type Handler struct {
	proto.AccessServiceServer
	svc *app.Service
}

func NewHandler(svc *app.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ListRepos(ctx context.Context, req *proto.ListReposRequest) (*proto.ListReposResponse, error) {
	accessList, err := h.svc.GetAccessInfo(ctx, req.OrgName)

	if err != nil {
		return nil, err
	}

	var repos []*proto.RepoAccessInfo

	for _, item := range accessList {
		repos = append(repos, &proto.RepoAccessInfo{
			RepoName:     item.RepoName,
			AllowedUsers: item.AllowedUsers,
			BlockedUsers: item.BlockedUsers,
		})
	}

	return &proto.ListReposResponse{Repos: repos}, nil
}
