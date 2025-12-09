package github

import (
	"context"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"
)

type Service struct {
	client *github.Client
	authed bool
}

func NewService() *Service {
	client := github.NewClient(nil)
	authed := false

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		client = github.NewClient(oauth2.NewClient(context.Background(), ts))
		authed = true
	}

	return &Service{client: client, authed: authed}
}

func (s *Service) GetUserStats(username string) (*User, *Stats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, _, err := s.client.Users.Get(ctx, username)
	if err != nil {
		return nil, nil, err
	}

	repos, err := s.getRepos(ctx, username)
	if err != nil {
		return nil, nil, err
	}

	statsChan := make(chan *Stats, 1)
	go func() {
		statsChan <- s.calcStats(ctx, username, repos)
	}()

	gitUser := &User{
		Login:       user.GetLogin(),
		Name:        user.GetName(),
		Bio:         user.GetBio(),
		Location:    user.GetLocation(),
		Followers:   user.GetFollowers(),
		Following:   user.GetFollowing(),
		PublicRepos: user.GetPublicRepos(),
		CreatedAt:   user.GetCreatedAt().Time,
		HTMLURL:     user.GetHTMLURL(),
		CachedAt:    time.Now(),
	}

	stats := <-statsChan
	return gitUser, stats, nil
}

func (s *Service) getRepos(ctx context.Context, username string) ([]*github.Repository, error) {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		Type:        "owner",
		Sort:        "updated",
		Direction:   "desc",
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := s.client.Repositories.List(ctx, username, opt)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepos, nil
}

func (s *Service) calcStats(ctx context.Context, username string, repos []*github.Repository) *Stats {
	stats := &Stats{}
	langs := make(map[string]int)
	dayCounts := make(map[string]int)
	commitCounts := make(map[string]int)

	since30d := time.Now().AddDate(0, 0, -30)
	since90d := time.Now().AddDate(0, -3, 0)

	for _, repo := range repos {
		stats.Stars += repo.GetStargazersCount()
		stats.Forks += repo.GetForksCount()
		stats.RepoCount++

		if lang := repo.GetLanguage(); lang != "" {
			langs[lang]++
		}

		if repo.GetUpdatedAt().After(since90d) {
			commits, _, err := s.client.Repositories.ListCommits(ctx, username,
				repo.GetName(), &github.CommitsListOptions{
					Author:      username,
					Since:       since90d,
					ListOptions: github.ListOptions{PerPage: 100},
				})
			if err != nil {
				continue
			}

			for _, commit := range commits {
				if commit.Commit != nil && commit.Commit.Author != nil && commit.Commit.Author.Date != nil {
					date := commit.Commit.Author.Date.Time

					day := date.Weekday().String()
					dayCounts[day]++

					dateStr := date.Format("2006-01-02")
					if date.After(since30d) {
						commitCounts[dateStr]++
					}
				}
			}
		}
	}

	stats.TopLanguages = getTopLanguages(langs, 3)
	stats.MostActiveDay = getMaxKey(dayCounts)
	stats.Commits30d = sumMapValues(commitCounts)

	return stats
}

func getTopLanguages(langs map[string]int, limit int) string {
	if len(langs) == 0 {
		return "No data"
	}

	type langCount struct {
		lang  string
		count int
	}
	var counts []langCount
	for lang, count := range langs {
		counts = append(counts, langCount{lang, count})
	}

	sort.Slice(counts, func(i, j int) bool { return counts[i].count > counts[j].count })

	var top []string
	for i := 0; i < len(counts) && i < limit; i++ {
		top = append(top, counts[i].lang)
	}
	return strings.Join(top, ", ")
}

func getMaxKey(m map[string]int) string {
	if len(m) == 0 {
		return "Unknown"
	}
	maxKey := ""
	maxVal := 0
	for k, v := range m {
		if v > maxVal {
			maxVal = v
			maxKey = k
		}
	}
	return maxKey
}

func sumMapValues(m map[string]int) int {
	total := 0
	for _, v := range m {
		total += v
	}
	return total
}

func (s *Service) CheckAPI() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _, err := s.client.Users.Get(ctx, "github")
	return err == nil, err
}

func (s *Service) Authed() bool { return s.authed }
