package wiki

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const baseURL = "https://bindingofisaacrebirth.fandom.com/api.php"

type SearchItem struct {
	Title   string
	Snippet string
}

type PageDetail struct {
	Title     string
	Extract   string
	Thumbnail string
}

type searchResponse struct {
	Query struct {
		Search []struct {
			Title   string `json:"title"`
			Snippet string `json:"snippet"`
		} `json:"search"`
	} `json:"query"`
}

type pageResponse struct {
	Query struct {
		Pages map[string]struct {
			Title     string `json:"title"`
			Extract   string `json:"extract"`
			Thumbnail struct {
				Source string `json:"source"`
			} `json:"thumbnail"`
		} `json:"pages"`
	} `json:"query"`
}

func Search(query string) ([]SearchItem, error) {
	params := url.Values{}
	params.Set("action", "query")
	params.Set("list", "search")
	params.Set("srsearch", query)
	params.Set("srlimit", "12")
	params.Set("format", "json")

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	items := make([]SearchItem, 0, len(result.Query.Search))
	for _, s := range result.Query.Search {
		items = append(items, SearchItem{
			Title:   s.Title,
			Snippet: cleanSnippet(s.Snippet),
		})
	}
	return items, nil
}

func GetPage(title string) (*PageDetail, error) {
	params := url.Values{}
	params.Set("action", "query")
	params.Set("titles", title)
	params.Set("prop", "extracts|pageimages")
	params.Set("exintro", "false")
	params.Set("explaintext", "true")
	params.Set("exchars", "1500")
	params.Set("pithumbsize", "400")
	params.Set("format", "json")

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result pageResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	for _, page := range result.Query.Pages {
		detail := &PageDetail{
			Title:     page.Title,
			Extract:   formatExtract(page.Extract),
			Thumbnail: page.Thumbnail.Source,
		}
		return detail, nil
	}

	return nil, fmt.Errorf("page not found: %s", title)
}

func cleanSnippet(s string) string {
	s = html.UnescapeString(s)
	re := regexp.MustCompile(`<[^>]+>`)
	return re.ReplaceAllString(s, "")
}

func formatExtract(s string) string {
	s = strings.TrimSpace(s)
	paragraphs := strings.Split(s, "\n")
	result := make([]string, 0)
	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, "<p>"+p+"</p>")
		}
	}
	return strings.Join(result, "\n")
}
