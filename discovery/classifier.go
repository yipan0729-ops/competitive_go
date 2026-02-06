package discovery

import (
	"net/url"
	"regexp"
	"strings"
	"time"
)

// LinkClassifier 链接分类器
type LinkClassifier struct{}

// LinkCategory 链接分类
type LinkCategory struct {
	Type     string  // 官网首页/产品功能/定价页面等
	Priority int     // 优先级 1-4
	Score    float64 // 质量评分 0-10
}

// ClassifyLink 分类链接
func (c *LinkClassifier) ClassifyLink(link, title, description string) *LinkCategory {
	// 先尝试规则匹配
	category := c.classifyByRules(link)
	if category != nil {
		return category
	}

	// 基于标题和描述的关键词匹配
	return c.classifyByContent(title, description)
}

// classifyByRules 基于URL规则分类
func (c *LinkClassifier) classifyByRules(link string) *LinkCategory {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return nil
	}

	path := strings.ToLower(parsedURL.Path)
	host := strings.ToLower(parsedURL.Host)

	// 官网首页
	if regexp.MustCompile(`^https?://[\w-]+\.(com|cn|io|ai|net)/?$`).MatchString(link) {
		return &LinkCategory{Type: "官网首页", Priority: 1, Score: 9.0}
	}

	// 产品功能
	if strings.Contains(path, "features") || strings.Contains(path, "product") ||
		strings.Contains(path, "functions") {
		return &LinkCategory{Type: "产品功能", Priority: 1, Score: 8.5}
	}

	// 定价页面
	if strings.Contains(path, "pricing") || strings.Contains(path, "price") ||
		strings.Contains(path, "plans") {
		return &LinkCategory{Type: "定价页面", Priority: 1, Score: 9.0}
	}

	// 关于我们
	if strings.Contains(path, "about") || strings.Contains(path, "company") {
		return &LinkCategory{Type: "关于我们", Priority: 2, Score: 7.0}
	}

	// 帮助文档
	if strings.Contains(path, "docs") || strings.Contains(path, "help") ||
		strings.Contains(path, "support") || strings.Contains(path, "guide") {
		return &LinkCategory{Type: "帮助文档", Priority: 3, Score: 6.5}
	}

	// 博客
	if strings.Contains(path, "blog") || strings.Contains(path, "news") ||
		strings.Contains(path, "article") {
		return &LinkCategory{Type: "博客文章", Priority: 3, Score: 6.0}
	}

	// 用户评价平台
	if strings.Contains(host, "xiaohongshu.com") || strings.Contains(host, "zhihu.com") ||
		strings.Contains(host, "douban.com") {
		return &LinkCategory{Type: "用户评价", Priority: 2, Score: 8.0}
	}

	// 电商平台
	if strings.Contains(host, "taobao.com") || strings.Contains(host, "jd.com") ||
		strings.Contains(host, "tmall.com") {
		return &LinkCategory{Type: "电商", Priority: 2, Score: 7.5}
	}

	// 社交媒体
	if strings.Contains(host, "weixin.qq.com") || strings.Contains(host, "weibo.com") {
		return &LinkCategory{Type: "社交媒体", Priority: 3, Score: 6.5}
	}

	return nil
}

// classifyByContent 基于内容分类
func (c *LinkClassifier) classifyByContent(title, description string) *LinkCategory {
	content := strings.ToLower(title + " " + description)

	keywords := map[string]struct {
		category string
		priority int
		score    float64
	}{
		"官网":   {"官网首页", 1, 8.5},
		"功能":   {"产品功能", 1, 8.0},
		"特性":   {"产品功能", 1, 8.0},
		"价格":   {"定价页面", 1, 8.5},
		"定价":   {"定价页面", 1, 8.5},
		"套餐":   {"定价页面", 1, 8.0},
		"评价":   {"用户评价", 2, 7.5},
		"评测":   {"用户评价", 2, 7.5},
		"怎么样":  {"用户评价", 2, 7.0},
		"使用体验": {"用户评价", 2, 7.0},
		"教程":   {"博客文章", 3, 6.0},
		"案例":   {"博客文章", 3, 6.0},
	}

	for keyword, info := range keywords {
		if strings.Contains(content, keyword) {
			return &LinkCategory{
				Type:     info.category,
				Priority: info.priority,
				Score:    info.score,
			}
		}
	}

	return &LinkCategory{Type: "其他", Priority: 4, Score: 5.0}
}

// LinkScorer 链接评分器
type LinkScorer struct{}

// ScoreLink 为链接评分
func (s *LinkScorer) ScoreLink(result *SearchResult, category *LinkCategory) float64 {
	// 相关性评分（基于位置）
	relevanceScore := 1.0 - float64(result.Position-1)*0.05
	if relevanceScore < 0.5 {
		relevanceScore = 0.5
	}

	// 信息价值评分（来自分类）
	valueScore := category.Score / 10.0

	// 时效性评分（暂时给固定值，后续可以解析时间）
	freshnessScore := 0.8

	// 综合评分
	finalScore := relevanceScore*0.4 + valueScore*0.4 + freshnessScore*0.2

	return finalScore
}

// DataSourceInfo 数据源信息
type DataSourceInfo struct {
	URL          string  `json:"url"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Type         string  `json:"type"`
	Priority     int     `json:"priority"`
	QualityScore float64 `json:"quality_score"`
}

// ProcessSearchResults 处理搜索结果，分类和评分
func ProcessSearchResults(results []SearchResult) []*DataSourceInfo {
	classifier := &LinkClassifier{}
	scorer := &LinkScorer{}

	dataSources := []*DataSourceInfo{}

	for _, result := range results {
		// 分类
		category := classifier.ClassifyLink(result.URL, result.Title, result.Description)

		// 评分
		score := scorer.ScoreLink(&result, category)

		// 过滤低质量链接
		if score < 0.6 {
			continue
		}

		dataSources = append(dataSources, &DataSourceInfo{
			URL:          result.URL,
			Title:        result.Title,
			Description:  result.Description,
			Type:         category.Type,
			Priority:     category.Priority,
			QualityScore: score,
		})
	}

	return dataSources
}

// SearchCacheManager 搜索缓存管理器
type SearchCacheManager struct {
	cacheDuration time.Duration
}

// NewSearchCacheManager 创建缓存管理器
func NewSearchCacheManager(cacheDays int) *SearchCacheManager {
	return &SearchCacheManager{
		cacheDuration: time.Duration(cacheDays) * 24 * time.Hour,
	}
}

// IsExpired 检查缓存是否过期
func (m *SearchCacheManager) IsExpired(cachedAt time.Time) bool {
	return time.Since(cachedAt) > m.cacheDuration
}

// GetExpiresAt 获取过期时间
func (m *SearchCacheManager) GetExpiresAt() time.Time {
	return time.Now().Add(m.cacheDuration)
}
