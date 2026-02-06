package crawler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// CrawlResult 爬取结果
type CrawlResult struct {
	Success    bool              `json:"success"`
	Markdown   string            `json:"markdown"`
	Title      string            `json:"title"`
	URL        string            `json:"url"`
	Platform   string            `json:"platform"`
	Method     string            `json:"method"` // firecrawl/jina/playwright
	Metadata   map[string]string `json:"metadata"`
	Error      string            `json:"error,omitempty"`
}

// Crawler 爬虫接口
type Crawler interface {
	Crawl(url string, platform *PlatformInfo) (*CrawlResult, error)
	Name() string
}

// FirecrawlCrawler Firecrawl爬虫
type FirecrawlCrawler struct {
	APIKey string
}

func (f *FirecrawlCrawler) Name() string {
	return "firecrawl"
}

func (f *FirecrawlCrawler) Crawl(url string, platform *PlatformInfo) (*CrawlResult, error) {
	if f.APIKey == "" {
		return nil, errors.New("Firecrawl API Key未配置")
	}

	// Firecrawl v2 API
	requestBody := map[string]interface{}{
		"url": url,
		"formats": []string{"markdown"},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("构造请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.firecrawl.dev/v1/scrape", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+f.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Firecrawl返回错误: %s", string(body))
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查是否成功
	success, ok := result["success"].(bool)
	if !ok || !success {
		return nil, errors.New("Firecrawl爬取失败")
	}

	// 提取数据
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Firecrawl响应格式错误")
	}

	markdown, _ := data["markdown"].(string)
	title := ""
	if metadata, ok := data["metadata"].(map[string]interface{}); ok {
		title, _ = metadata["title"].(string)
	}

	// 验证内容
	if len(markdown) < 100 {
		return nil, errors.New("内容过短，可能是验证页面")
	}

	if strings.Contains(strings.ToLower(markdown), "验证") ||
		strings.Contains(strings.ToLower(markdown), "captcha") {
		return nil, errors.New("触发了验证码")
	}

	return &CrawlResult{
		Success:  true,
		Markdown: markdown,
		Title:    title,
		URL:      url,
		Platform: platform.Name,
		Method:   "firecrawl",
		Metadata: map[string]string{
			"api": "firecrawl-v2",
		},
	}, nil
}

// JinaCrawler Jina Reader爬虫
type JinaCrawler struct{}

func (j *JinaCrawler) Name() string {
	return "jina"
}

func (j *JinaCrawler) Crawl(url string, platform *PlatformInfo) (*CrawlResult, error) {
	// Jina Reader API
	jinaURL := "https://r.jina.ai/" + url

	req, err := http.NewRequest("GET", jinaURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Accept", "text/markdown")
	req.Header.Set("User-Agent", platform.UserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Jina返回错误: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	markdown := string(body)

	// 验证内容
	if len(markdown) < 100 {
		return nil, errors.New("内容过短")
	}

	// 从markdown中提取标题（通常第一行是# 标题）
	title := ""
	lines := strings.Split(markdown, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			title = strings.TrimPrefix(line, "# ")
			break
		}
	}

	return &CrawlResult{
		Success:  true,
		Markdown: markdown,
		Title:    title,
		URL:      url,
		Platform: platform.Name,
		Method:   "jina",
		Metadata: map[string]string{
			"api": "jina-reader",
		},
	}, nil
}

// ThreeLayerCrawler 三层策略爬虫
type ThreeLayerCrawler struct {
	Crawlers []Crawler
}

// NewThreeLayerCrawler 创建三层策略爬虫
func NewThreeLayerCrawler(firecrawlKey string) *ThreeLayerCrawler {
	crawlers := []Crawler{}

	// 第一层：Firecrawl（如果配置了）
	if firecrawlKey != "" {
		crawlers = append(crawlers, &FirecrawlCrawler{APIKey: firecrawlKey})
	}

	// 第二层：Jina（免费）
	crawlers = append(crawlers, &JinaCrawler{})

	// 第三层：Playwright（暂未实现，需要浏览器环境）
	// TODO: 实现Playwright爬虫

	return &ThreeLayerCrawler{
		Crawlers: crawlers,
	}
}

// Crawl 使用三层策略爬取
func (t *ThreeLayerCrawler) Crawl(url string) (*CrawlResult, error) {
	// 识别平台
	platform, err := IdentifyPlatform(url)
	if err != nil {
		return nil, err
	}

	var lastError error

	// 依次尝试每个爬虫
	for _, crawler := range t.Crawlers {
		result, err := crawler.Crawl(url, platform)
		if err == nil && result.Success {
			return result, nil
		}
		lastError = err
	}

	if lastError != nil {
		return nil, fmt.Errorf("所有爬虫都失败了，最后一个错误: %w", lastError)
	}

	return nil, errors.New("所有爬虫都失败了")
}
