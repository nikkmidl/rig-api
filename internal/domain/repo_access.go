package domain

type RepoAccess struct {
	RepoName     string
	AllowedUsers []string
	BlockedUsers []string
}
