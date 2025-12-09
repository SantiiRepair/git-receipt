package cache

import (
	"santiirepair.dev/git-receipt/github"
	"time"
)

type CachedUserData struct {
	User     *github.User
	Stats    *github.Stats
	CachedAt time.Time
	TTL      time.Duration
}
