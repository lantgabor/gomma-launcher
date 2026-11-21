package moddb

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type ModInfo struct {
	Filename    string
	MD5Hash     string
	DownloadURL string
}

type Scraper struct {
	UserAgent  string
	HTTPClient *http.Client
}

func NewScraper(client *http.Client) *Scraper {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &Scraper{
		UserAgent:  "GommaLauncher/1.0",
		HTTPClient: client,
	}
}

func (s *Scraper) GetModInfo(iurl string) (*ModInfo, error) {
	info := &ModInfo{}

	baseURL, err := url.Parse(iurl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mod url %s: %w", iurl, err)
	}

	c := colly.NewCollector(
		colly.UserAgent(s.UserAgent),
		colly.AllowedDomains("www.moddb.com", "moddb.com"),
	)
	c.SetClient(s.HTTPClient)

	var downloadErr error
	c.OnHTML("body", func(b *colly.HTMLElement) {
		downloadLink := strings.TrimSpace(b.ChildAttr(`#downloadmirrorstoggle`, "href"))
		if downloadLink == "" {
			downloadErr = fmt.Errorf("no download link found on %s", iurl)
			return
		}
		ref, err := url.Parse(downloadLink)
		if err != nil {
			downloadErr = fmt.Errorf("failed to parse download link %q: %w", downloadLink, err)
			return
		}
		info.DownloadURL, downloadErr = s.GetDownloadURL(baseURL.ResolveReference(ref).String())
	})

	c.OnHTML("div.row.clear", func(d *colly.HTMLElement) {
		h5 := d.ChildText("h5")
		span := strings.TrimSpace(d.ChildText("span"))
		switch h5 {
		case "Filename":
			info.Filename = span
		case "MD5 Hash":
			info.MD5Hash = span
		}
	})

	if err := c.Visit(iurl); err != nil {
		return nil, fmt.Errorf("failed to scrape %s: %w", iurl, err)
	}

	if downloadErr != nil {
		return nil, fmt.Errorf("failed to get download url %s: %w", iurl, downloadErr)
	}
	if info.DownloadURL == "" {
		return nil, fmt.Errorf("no download url found on %s", iurl)
	}

	return info, nil
}

func (s *Scraper) GetDownloadURL(dlStartURL string) (string, error) {
	mirrorURL, err := s.extractDownloadURL(dlStartURL)
	if err != nil {
		return "", err
	}

	finalURL, err := s.followRedirect(mirrorURL)
	if err != nil {
		return "", err
	}

	return finalURL, nil
}

func (s *Scraper) extractDownloadURL(dlStartURL string) (string, error) {
	req, err := http.NewRequest("GET", dlStartURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed creating request for %s: %w", dlStartURL, err)
	}
	req.Header.Set("User-Agent", s.UserAgent)

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed fetching download start page %s: %w", dlStartURL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed reading download start page %s: %w", dlStartURL, err)
	}

	html := string(body)

	start := strings.Index(html, `"/downloads/mirror/`)
	if start == -1 {
		return "", fmt.Errorf("download link not found when requesting %s", dlStartURL)
	}
	start++

	end := strings.IndexByte(html[start:], '"')
	if end == -1 {
		return "", fmt.Errorf("download link not found when requesting %s", dlStartURL)
	}

	return "https://www.moddb.com" + html[start:start+end], nil
}

func (s *Scraper) followRedirect(targetURL string) (string, error) {
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed creating request for %s: %w", targetURL, err)
	}
	req.Header.Set("User-Agent", s.UserAgent)

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed requesting %s: %w", targetURL, err)
	}
	resp.Body.Close()

	return resp.Request.URL.String(), nil
}
