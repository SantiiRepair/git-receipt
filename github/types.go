package github

import "time"

type GitHubUser struct {
	Login       string    `json:"login"`
	Name        string    `json:"name"`
	PublicRepos int       `json:"public_repos"`
	Followers   int       `json:"followers"`
	Following   int       `json:"following"`
	Bio         string    `json:"bio"`
	Location    string    `json:"location"`
	CreatedAt   time.Time `json:"created_at"`
	HTMLURL     string    `json:"html_url"`
	CachedAt    time.Time `json:"cached_at"`
}

type GitHubRepo struct {
	Name            string    `json:"name"`
	StargazersCount int       `json:"stargazers_count"`
	ForksCount      int       `json:"forks_count"`
	Language        string    `json:"language"`
	UpdatedAt       time.Time `json:"updated_at"`
}
