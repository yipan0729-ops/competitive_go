package discovery

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// SearchEngine 搜索引擎接口
type SearchEngine interface {
	Search(query string, numResults int) ([]SearchResult, error)
	Name() string
}

// SearchResult 搜索结果
type SearchResult struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Position    int    `json:"position"`
}

// SerperSearchEngine Serper搜索引擎
type SerperSearchEngine struct {
	APIKey string
}

func (s *SerperSearchEngine) Name() string {
	return "serper"
}

func (s *SerperSearchEngine) Search(query string, numResults int) ([]SearchResult, error) {
	if s.APIKey == "" {
		return nil, errors.New("Serper API Key未配置")
	}

	requestBody := map[string]interface{}{
		"q":   query,
		"num": numResults,
		"gl":  "cn", // 中国地区
		"hl":  "zh-cn",
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://google.serper.dev/search", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-KEY", s.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Serper返回错误 %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// 提取organic结果
	organic, ok := result["organic"].([]interface{})
	if !ok {
		return nil, errors.New("Serper响应格式错误")
	}

	results := []SearchResult{}
	for i, item := range organic {
		if i >= numResults {
			break
		}

		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		title, _ := itemMap["title"].(string)
		link, _ := itemMap["link"].(string)
		snippet, _ := itemMap["snippet"].(string)

		results = append(results, SearchResult{
			Title:       title,
			URL:         link,
			Description: snippet,
			Position:    i + 1,
		})
	}

	return results, nil
}

// GoogleSearchEngine Google搜索引擎
type GoogleSearchEngine struct {
	APIKey   string
	EngineID string
}

func (g *GoogleSearchEngine) Name() string {
	return "google"
}

func (g *GoogleSearchEngine) Search(query string, numResults int) ([]SearchResult, error) {
	if g.APIKey == "" || g.EngineID == "" {
		return nil, errors.New("Google API Key或Engine ID未配置")
	}

	apiURL := fmt.Sprintf(
		"https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s&num=%d",
		g.APIKey, g.EngineID, url.QueryEscape(query), numResults,
	)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google返回错误 %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	items, ok := result["items"].([]interface{})
	if !ok {
		return nil, errors.New("Google响应格式错误")
	}

	results := []SearchResult{}
	for i, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		title, _ := itemMap["title"].(string)
		link, _ := itemMap["link"].(string)
		snippet, _ := itemMap["snippet"].(string)

		results = append(results, SearchResult{
			Title:       title,
			URL:         link,
			Description: snippet,
			Position:    i + 1,
		})
	}

	return results, nil
}

// BingSearchEngine Bing搜索引擎
type BingSearchEngine struct {
	APIKey string
}

func (b *BingSearchEngine) Name() string {
	return "bing"
}

func (b *BingSearchEngine) Search(query string, numResults int) ([]SearchResult, error) {
	if b.APIKey == "" {
		return nil, errors.New("Bing API Key未配置")
	}

	apiURL := fmt.Sprintf(
		"https://api.bing.microsoft.com/v7.0/search?q=%s&count=%d&mkt=zh-CN",
		url.QueryEscape(query), numResults,
	)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", b.APIKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bing返回错误 %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	webPages, ok := result["webPages"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Bing响应格式错误")
	}

	values, ok := webPages["value"].([]interface{})
	if !ok {
		return nil, errors.New("Bing响应格式错误")
	}

	results := []SearchResult{}
	for i, item := range values {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		name, _ := itemMap["name"].(string)
		url, _ := itemMap["url"].(string)
		snippet, _ := itemMap["snippet"].(string)

		results = append(results, SearchResult{
			Title:       name,
			URL:         url,
			Description: snippet,
			Position:    i + 1,
		})
	}

	return results, nil
}
