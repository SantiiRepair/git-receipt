package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"
)

type GitHubUser struct {
	Login       string `json:"login"`
	Name        string `json:"name"`
	PublicRepos int    `json:"public_repos"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
	Bio         string `json:"bio"`
	Location    string `json:"location"`
	CreatedAt   string `json:"created_at"`
	HTMLURL     string `json:"html_url"`
}

type GitHubRepo struct {
	Name            string `json:"name"`
	StargazersCount int    `json:"stargazers_count"`
	ForksCount      int    `json:"forks_count"`
	Language        string `json:"language"`
	UpdatedAt       string `json:"updated_at"`
}

type ReceiptData struct {
	Username          string
	FormattedDate     string
	OrderNumber       string
	CustomerName      string
	PublicRepos       int
	TotalStars        int
	TotalForks        int
	Followers         int
	Following         int
	TopLanguages      string
	MostActiveDay     string
	Commits30d        int
	ContributionScore int
	ServerName        string
	TimeString        string
	CouponCode        string
	AuthCode          string
	CardYear          int
}

var (
	servers         = []string{"Grace Hopper", "Alan Turing", "Ada Lovelace", "Tim Berners-Lee", "Linus Torvalds"}
	days            = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	templateContent string
	githubClient    *github.Client
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		githubClient = github.NewClient(tc)
	} else {
		githubClient = github.NewClient(nil)
	}
}

func loadTemplate() error {
	content, err := os.ReadFile("TEMPLATE")
	if err != nil {
		return fmt.Errorf("error loading template: %v", err)
	}
	templateContent = string(content)
	return nil
}

func generateRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func fetchGitHubUser(username string) (*GitHubUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, _, err := githubClient.Users.Get(ctx, username)
	if err != nil {
		return nil, err
	}

	return &GitHubUser{
		Login:       user.GetLogin(),
		Name:        user.GetName(),
		PublicRepos: user.GetPublicRepos(),
		Followers:   user.GetFollowers(),
		Following:   user.GetFollowing(),
		Bio:         user.GetBio(),
		Location:    user.GetLocation(),
		CreatedAt:   user.GetCreatedAt().Format(time.RFC3339),
		HTMLURL:     user.GetHTMLURL(),
	}, nil
}

func fetchGitHubRepos(username string) ([]GitHubRepo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		Sort:        "updated",
		Direction:   "desc",
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := githubClient.Repositories.List(ctx, username, opt)
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
			UpdatedAt:       repo.GetUpdatedAt().Format(time.RFC3339),
		})
	}

	return result, nil
}

func calculateMostActiveDay(username string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	repos, err := fetchGitHubRepos(username)
	if err != nil {
		return "Unknown"
	}

	dayCounts := make(map[string]int)
	since := time.Now().AddDate(0, -3, 0)

	for _, repo := range repos {
		commits, _, err := githubClient.Repositories.ListCommits(ctx, username, repo.Name, &github.CommitsListOptions{
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

func calculateStats(repos []GitHubRepo) (int, int, string) {
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

func generateReceiptHTML(data ReceiptData) string {
	html := templateContent
	html = strings.ReplaceAll(html, "{{.FormattedDate}}", data.FormattedDate)
	html = strings.ReplaceAll(html, "{{.OrderNumber}}", data.OrderNumber)
	html = strings.ReplaceAll(html, "{{.CustomerName}}", data.CustomerName)
	html = strings.ReplaceAll(html, "{{.Username}}", data.Username)
	html = strings.ReplaceAll(html, "{{.PublicRepos}}", fmt.Sprintf("%d", data.PublicRepos))
	html = strings.ReplaceAll(html, "{{.TotalStars}}", fmt.Sprintf("%d", data.TotalStars))
	html = strings.ReplaceAll(html, "{{.TotalForks}}", fmt.Sprintf("%d", data.TotalForks))
	html = strings.ReplaceAll(html, "{{.Followers}}", fmt.Sprintf("%d", data.Followers))
	html = strings.ReplaceAll(html, "{{.Following}}", fmt.Sprintf("%d", data.Following))
	html = strings.ReplaceAll(html, "{{.TopLanguages}}", data.TopLanguages)
	html = strings.ReplaceAll(html, "{{.MostActiveDay}}", data.MostActiveDay)
	html = strings.ReplaceAll(html, "{{.Commits30d}}", fmt.Sprintf("%d", data.Commits30d))
	html = strings.ReplaceAll(html, "{{.ContributionScore}}", fmt.Sprintf("%d", data.ContributionScore))
	html = strings.ReplaceAll(html, "{{.ServerName}}", data.ServerName)
	html = strings.ReplaceAll(html, "{{.TimeString}}", data.TimeString)
	html = strings.ReplaceAll(html, "{{.CouponCode}}", data.CouponCode)
	html = strings.ReplaceAll(html, "{{.AuthCode}}", data.AuthCode)
	html = strings.ReplaceAll(html, "{{.CardYear}}", fmt.Sprintf("%d", data.CardYear))
	return html
}

func main() {
	if err := loadTemplate(); err != nil {
		fmt.Printf("❌ Error loading template: %v\n", err)
		os.Exit(1)
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
		AppName:               "GitHub Receipt Service",
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept",
	}))

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path}\n",
	}))

	app.Get("/ping", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, _, err := githubClient.Users.Get(ctx, "github")
		githubStatus := "ok"
		if err != nil {
			githubStatus = "error"
		}

		return c.JSON(fiber.Map{
			"status":     "ok",
			"timestamp":  time.Now().Unix(),
			"github_api": githubStatus,
		})
	})

	app.Get("/:username", func(c *fiber.Ctx) error {
		username := c.Params("username")

		if username == "" {
			return c.Status(400).SendString("Username is required")
		}

		user, err := fetchGitHubUser(username)
		if err != nil {
			return c.Status(404).SendString(fmt.Sprintf("User '%s' not found: %v", username, err))
		}

		repos, err := fetchGitHubRepos(username)
		if err != nil {
			return c.Status(500).SendString(fmt.Sprintf("Error fetching repositories: %v", err))
		}

		mostActiveDay := calculateMostActiveDay(username)
		totalStars, totalForks, topLanguages := calculateStats(repos)
		contributionScore := user.PublicRepos*3 + user.Followers*2 + totalStars

		now := time.Now()
		data := ReceiptData{
			Username:          user.Login,
			FormattedDate:     now.Format("Monday, January 02, 2006"),
			OrderNumber:       fmt.Sprintf("%04d", rand.Intn(10000)),
			CustomerName:      user.Name,
			PublicRepos:       user.PublicRepos,
			TotalStars:        totalStars,
			TotalForks:        totalForks,
			Followers:         user.Followers,
			Following:         user.Following,
			TopLanguages:      topLanguages,
			MostActiveDay:     mostActiveDay,
			Commits30d:        rand.Intn(91) + 10,
			ContributionScore: contributionScore,
			ServerName:        servers[rand.Intn(len(servers))],
			TimeString:        now.Format("3:04:05 PM"),
			CouponCode:        generateRandomString(6),
			AuthCode:          fmt.Sprintf("%06d", rand.Intn(1000000)),
			CardYear:          now.Year(),
		}

		if data.CustomerName == "" {
			data.CustomerName = user.Login
		}

		html := generateReceiptHTML(data)
		c.Set("Content-Type", "text/html; charset=utf-8")
		c.Set("Cache-Control", "public, max-age=300")
		return c.SendString(html)
	})

	var port string
	if port := os.Getenv("PORT"); port == "" {
		port = "5000"
	}

	if err := app.Listen(":" + port); err != nil {
		os.Exit(1)
	}
}
