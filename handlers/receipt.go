package handlers

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"santiirepair.dev/git-receipt/cache"
	"santiirepair.dev/git-receipt/utils"
)

type ReceiptHandler struct {
	cacheManager    *cache.GitHubCacheManager
	templateContent string
	servers         []string
}

func NewReceiptHandler(cacheManager *cache.GitHubCacheManager, templateContent string, servers []string) *ReceiptHandler {
	return &ReceiptHandler{
		cacheManager:    cacheManager,
		templateContent: templateContent,
		servers:         servers,
	}
}

func (h *ReceiptHandler) GenerateReceipt(username string) (string, error) {
	cachedData, err := h.cacheManager.GetUserData(username)
	if err != nil {
		return "", err
	}

	user := cachedData.User
	contributionScore := user.PublicRepos*3 + user.Followers*2 + cachedData.Stats.Stars

	now := time.Now()
	data := ReceiptData{
		Username:          user.Login,
		FormattedDate:     now.Format("Monday, January 02, 2006"),
		OrderNumber:       fmt.Sprintf("%04d", rand.Intn(10000)),
		CustomerName:      user.Name,
		PublicRepos:       user.PublicRepos,
		TotalStars:        cachedData.Stats.Stars,
		TotalForks:        cachedData.Stats.Forks,
		Followers:         user.Followers,
		Following:         user.Following,
		TopLanguages:      cachedData.Stats.TopLanguages,
		MostActiveDay:     cachedData.Stats.MostActiveDay,
		Commits30d:        cachedData.Stats.Commits30d,
		ContributionScore: contributionScore,
		ServerName:        h.servers[rand.Intn(len(h.servers))],
		TimeString:        now.Format("3:04:05 PM"),
		CouponCode:        utils.GenerateRandomString(6),
		AuthCode:          fmt.Sprintf("%06d", rand.Intn(1000000)),
		CardYear:          now.Year(),
	}

	if data.CustomerName == "" {
		data.CustomerName = user.Login
	}

	html := h.generateReceiptHTML(data)
	return html, nil
}

func (h *ReceiptHandler) generateReceiptHTML(data ReceiptData) string {
	html := h.templateContent
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
