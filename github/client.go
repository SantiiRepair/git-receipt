package github

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"
)

type GitHubService struct {
	client   *github.Client
	withAuth bool
}

func NewGitHubService() *GitHubService {
	token := ""
	client := github.NewClient(nil)
	withAuth := false

	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		client = github.NewClient(tc)
		withAuth = true
	}

	return &GitHubService{
		client:   client,
		withAuth: withAuth,
	}
}

func (g *GitHubService) WithAuth() bool {
	return g.withAuth
}

func (g *GitHubService) GetUserAndRepos(username string) (*GitHubUser, []GitHubRepo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user, _, err := g.client.Users.Get(ctx, username)
	if err != nil {
		return nil, nil, err
	}

	repos, err := g.fetchUserRepos(username)
	if err != nil {
		return nil, nil, err
	}

	gitHubUser := &GitHubUser{
		Login:       user.GetLogin(),
		Name:        user.GetName(),
		PublicRepos: user.GetPublicRepos(),
		Followers:   user.GetFollowers(),
		Following:   user.GetFollowing(),
		Bio:         user.GetBio(),
		Location:    user.GetLocation(),
		CreatedAt:   user.GetCreatedAt().Time,
		HTMLURL:     user.GetHTMLURL(),
		CachedAt:    time.Now(),
	}

	return gitHubUser, repos, nil
}

func (g *GitHubService) fetchUserRepos(username string) ([]GitHubRepo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		Sort:        "updated",
		Direction:   "desc",
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := g.client.Repositories.List(ctx, username, opt)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	var result []GitHubRepo
	for _, repo := range allRepos {
		result = append(result, GitHubRepo{
			Name:            repo.GetName(),
			StargazersCount: repo.GetStargazersCount(),
			ForksCount:      repo.GetForksCount(),
			Language:        repo.GetLanguage(),
			UpdatedAt:       repo.GetUpdatedAt().Time,
		})
	}

	return result, nil
}

func (g *GitHubService) CalculateMostActiveDay(username string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	repos, err := g.fetchUserRepos(username)
	if err != nil {
		return "Unknown"
	}

	dayCounts := make(map[string]int)
	since := time.Now().AddDate(0, -3, 0)

	for _, repo := range repos {
		commits, _, err := g.client.Repositories.ListCommits(ctx, username, repo.Name, &github.CommitsListOptions{
			Author: username,
			Since:  since,
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		})

		if err != nil {
			continue
		}

		for _, commit := range commits {
			if commit.Commit != nil && commit.Commit.Author != nil && commit.Commit.Author.Date != nil {
				day := commit.Commit.Author.Date.Weekday().String()
				dayCounts[day]++
			}
		}
	}

	if len(dayCounts) == 0 {
		return "Unknown"
	}

	maxDay := ""
	maxCount := 0
	for day, count := range dayCounts {
		if count > maxCount {
			maxCount = count
			maxDay = day
		}
	}

	return maxDay
}

func (g *GitHubService) CalculateStats(repos []GitHubRepo) (int, int, string) {
	totalStars := 0
	totalForks := 0
	languages := make(map[string]int)

	if len(repos) > 0 {
		for _, repo := range repos {
			totalStars += repo.StargazersCount
			totalForks += repo.ForksCount
			if repo.Language != "" {
				languages[repo.Language]++
			}
		}
	}

	topLanguages := "No data"
	if len(languages) > 0 {
		type langCount struct {
			lang  string
			count int
		}
		var langCounts []langCount
		for lang, count := range languages {
			langCounts = append(langCounts, langCount{lang, count})
		}

		sort.Slice(langCounts, func(i, j int) bool {
			return langCounts[i].count > langCounts[j].count
		})

		var topLangs []string
		for i := 0; i < len(langCounts) && i < 3; i++ {
			topLangs = append(topLangs, langCounts[i].lang)
		}

		if len(topLangs) > 0 {
			topLanguages = strings.Join(topLangs, ", ")
		}
	}

	return totalStars, totalForks, topLanguages
}

func (g *GitHubService) CheckAPIStatus() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, _, err := g.client.Users.Get(ctx, "github")
	return err == nil, err
}
