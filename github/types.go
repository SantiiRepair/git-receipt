package github

import "time"

type User struct {
	Login       string    `json:"login"`
	Name        string    `json:"name"`
	Bio         string    `json:"bio"`
	Location    string    `json:"location"`
	Followers   int       `json:"followers"`
	Following   int       `json:"following"`
	PublicRepos int       `json:"public_repos"`
	CreatedAt   time.Time `json:"created_at"`
	HTMLURL     string    `json:"html_url"`
	CachedAt    time.Time `json:"cached_at"`
}

type Stats struct {
	Stars         int    `json:"stars"`
	Forks         int    `json:"forks"`
	RepoCount     int    `json:"repo_count"`
	TopLanguages  string `json:"top_languages"`
	MostActiveDay string `json:"most_active_day"`
	Commits30d    int    `json:"commits_30d"`
}