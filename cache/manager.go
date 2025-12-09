package cache

import (
	"log"
	"time"

	"github.com/dgraph-io/ristretto"
	"santiirepair.dev/git-receipt/github"
)

type GitHubCacheManager struct {
	cache         *ristretto.Cache
	githubService *github.Service
}

func NewGitHubCacheManager(githubService *github.Service) *GitHubCacheManager {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e6,
		MaxCost:     1e6,
		BufferItems: 64,
		Metrics:     true,
	})

	if err != nil {
		log.Fatalf("Failed to create cache: %v", err)
	}

	return &GitHubCacheManager{
		cache:         cache,
		githubService: githubService,
	}
}

func (g *GitHubCacheManager) calculateTTL(username string, user *github.User) time.Duration {
	baseTTL := 30 * time.Minute

	if user.Followers > 1000 || user.PublicRepos > 50 {
		return 15 * time.Minute
	}

	if user.PublicRepos == 0 && user.Followers < 10 {
		return 2 * time.Hour
	}

	return baseTTL
}

func (g *GitHubCacheManager) GetUserData(username string) (*CachedUserData, error) {
	if cached, found := g.cache.Get(username); found {
		cachedData := cached.(*CachedUserData)
		if time.Since(cachedData.CachedAt) < cachedData.TTL {
			return cachedData, nil
		}
	}

	user, stats, err := g.githubService.GetUserStats(username)
	if err != nil {
		return nil, err
	}

	userData := &CachedUserData{
		User:     user,
		Stats:    stats,
		CachedAt: time.Now(),
	}

	userData.TTL = g.calculateTTL(username, user)
	cost := int64(len(user.Login) + len(user.Name) +
		len(stats.TopLanguages) + (stats.RepoCount * 50))

	g.cache.SetWithTTL(username, userData, cost, userData.TTL)

	return userData, nil
}

func (g *GitHubCacheManager) GetCacheMetrics() map[string]interface{} {
	metrics := g.cache.Metrics
	return map[string]interface{}{
		"hits":       metrics.Hits(),
		"misses":     metrics.Misses(),
		"ratio":      metrics.Ratio() * 100,
		"keys_added": metrics.KeysAdded(),
		"cost_added": metrics.CostAdded(),
	}
}

func (c *CachedUserData) GetMostActiveDay() string {
	if c.Stats == nil {
		return "Unknown"
	}
	return c.Stats.MostActiveDay
}

func (c *CachedUserData) GetTotalStars() int {
	if c.Stats == nil {
		return 0
	}
	return c.Stats.Stars
}

func (c *CachedUserData) GetTotalForks() int {
	if c.Stats == nil {
		return 0
	}
	return c.Stats.Forks
}

func (c *CachedUserData) GetTopLanguages() string {
	if c.Stats == nil {
		return "No data"
	}
	return c.Stats.TopLanguages
}

func (c *CachedUserData) GetCommits30d() int {
	if c.Stats == nil {
		return 0
	}
	return c.Stats.Commits30d
}
