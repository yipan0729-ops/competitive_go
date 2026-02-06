package crawler

import (
	"fmt"
	"net/url"
	"strings"
)

// PlatformInfo 平台信息
type PlatformInfo struct {
	Name         string
	NeedsLogin   bool
	Priority     int // 1=Firecrawl, 2=Jina, 3=Playwright
	UserAgent    string
}

// IdentifyPlatform 识别URL对应的平台
func IdentifyPlatform(rawURL string) (*PlatformInfo, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("无效的URL: %w", err)
	}

	host := strings.ToLower(parsedURL.Host)

	// 微信公众号
	if strings.Contains(host, "mp.weixin.qq.com") {
		return &PlatformInfo{
			Name:       "微信公众号",
			NeedsLogin: false, // 短链接不需要登录
			Priority:   1,     // 优先用Firecrawl
			UserAgent:  "Mozilla/5.0 (Linux; Android 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36 MicroMessenger/8.0.0",
		}, nil
	}

	// 小红书
	if strings.Contains(host, "xiaohongshu.com") || strings.Contains(host, "xhslink.com") {
		return &PlatformInfo{
			Name:       "小红书",
			NeedsLogin: true, // 需要登录才能看完整内容
			Priority:   3,
			UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		}, nil
	}

	// 知乎
	if strings.Contains(host, "zhihu.com") {
		return &PlatformInfo{
			Name:       "知乎",
			NeedsLogin: false,
			Priority:   1,
			UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		}, nil
	}

	// 抖音
	if strings.Contains(host, "douyin.com") {
		return &PlatformInfo{
			Name:       "抖音",
			NeedsLogin: true,
			Priority:   3,
			UserAgent:  "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X)",
		}, nil
	}

	// 淘宝
	if strings.Contains(host, "taobao.com") || strings.Contains(host, "tmall.com") {
		return &PlatformInfo{
			Name:       "淘宝/天猫",
			NeedsLogin: false,
			Priority:   2,
			UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		}, nil
	}

	// 京东
	if strings.Contains(host, "jd.com") {
		return &PlatformInfo{
			Name:       "京东",
			NeedsLogin: false,
			Priority:   2,
			UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		}, nil
	}

	// B站
	if strings.Contains(host, "bilibili.com") {
		return &PlatformInfo{
			Name:       "哔哩哔哩",
			NeedsLogin: false,
			Priority:   1,
			UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		}, nil
	}

	// 微博
	if strings.Contains(host, "weibo.com") {
		return &PlatformInfo{
			Name:       "微博",
			NeedsLogin: false,
			Priority:   2,
			UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		}, nil
	}

	// 默认为普通网站
	return &PlatformInfo{
		Name:       "普通网站",
		NeedsLogin: false,
		Priority:   1,
		UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}, nil
}
