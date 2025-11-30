package cache

import (
	"santiirepair.dev/git-receipt/github"
	"time"
)

type CachedUserData struct {
	User          *github.GitHubUser  `json:"user"`
	Repos         []github.GitHubRepo `json:"repos"`
	MostActiveDay string              `json:"most_active_day"`
	TotalStars    int                 `json:"total_stars"`
	TotalForks    int                 `json:"total_forks"`
	TopLanguages  string              `json:"top_languages"`
	CachedAt      time.Time           `json:"cached_at"`
	TTL           time.Duration       `json:"ttl"`
}
