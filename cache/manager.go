package cache

import (
	"log"
	"time"

	"github.com/dgraph-io/ristretto"
	"santiirepair.dev/git-receipt/github"
)

type GitHubCacheManager struct {
	cache         *ristretto.Cache
	githubService *github.GitHubService
}

func NewGitHubCacheManager(githubService *github.GitHubService) *GitHubCacheManager {
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

func (g *GitHubCacheManager) calculateTTL(username string, user *github.GitHubUser) time.Duration {
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
			log.Printf("✅ Cache HIT for user: %s", username)
			return cachedData, nil
		}

		log.Printf("⚠️ Cache EXPIRED for user: %s", username)
	}

	log.Printf("🔍 Cache MISS for user: %s", username)

	userData, err := g.fetchFreshUserData(username)
	if err != nil {
		return nil, err
	}

	userData.TTL = g.calculateTTL(username, userData.User)
	userData.CachedAt = time.Now()

	cost := int64(len(userData.User.Login) + len(userData.User.Name) +
		len(userData.TopLanguages) + (len(userData.Repos) * 100))

	g.cache.SetWithTTL(username, userData, cost, userData.TTL)

	log.Printf("💾 Cached user: %s for %v", username, userData.TTL)

	return userData, nil
}

func (g *GitHubCacheManager) fetchFreshUserData(username string) (*CachedUserData, error) {
	user, repos, err := g.githubService.GetUserAndRepos(username)
	if err != nil {
		return nil, err
	}

	mostActiveDay := g.githubService.CalculateMostActiveDay(username)
	totalStars, totalForks, topLanguages := g.githubService.CalculateStats(repos)

	return &CachedUserData{
		User:          user,
		Repos:         repos,
		MostActiveDay: mostActiveDay,
		TotalStars:    totalStars,
		TotalForks:    totalForks,
		TopLanguages:  topLanguages,
		CachedAt:      time.Now(),
	}, nil
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
