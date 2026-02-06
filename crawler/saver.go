package crawler

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ContentSaver 内容保存器
type ContentSaver struct {
	StoragePath string
}

// SaveResult 保存结果
type SaveResult struct {
	ContentPath string   `json:"content_path"`
	ImagePaths  []string `json:"image_paths"`
	Title       string   `json:"title"`
	URL         string   `json:"url"`
}

// NewContentSaver 创建内容保存器
func NewContentSaver(storagePath string) *ContentSaver {
	return &ContentSaver{
		StoragePath: storagePath,
	}
}

// Save 保存爬取的内容
func (s *ContentSaver) Save(result *CrawlResult, competitorName string) (*SaveResult, error) {
	// 创建保存目录：日期_竞品名_标题
	date := time.Now().Format("20060102")
	title := sanitizeFilename(result.Title)
	if title == "" {
		title = "untitled"
	}

	dirName := fmt.Sprintf("%s_%s_%s", date, sanitizeFilename(competitorName), title)
	savePath := filepath.Join(s.StoragePath, dirName)

	if err := os.MkdirAll(savePath, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}

	// 下载图片并替换链接
	markdown, imagePaths, err := s.downloadImages(result.Markdown, savePath, result.Platform)
	if err != nil {
		// 图片下载失败不影响主流程
		markdown = result.Markdown
	}

	// 添加元数据
	metadata := fmt.Sprintf(`---
title: %s
source: %s
platform: %s
url: %s
crawl_method: %s
crawl_time: %s
---

`, result.Title, competitorName, result.Platform, result.URL, result.Method, time.Now().Format("2006-01-02 15:04:05"))

	fullContent := metadata + markdown

	// 保存Markdown文件
	contentPath := filepath.Join(savePath, "content.md")
	if err := os.WriteFile(contentPath, []byte(fullContent), 0644); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	return &SaveResult{
		ContentPath: contentPath,
		ImagePaths:  imagePaths,
		Title:       result.Title,
		URL:         result.URL,
	}, nil
}

// downloadImages 下载图片并返回新的markdown内容
func (s *ContentSaver) downloadImages(markdown, savePath, platform string) (string, []string, error) {
	imagePaths := []string{}
	imageCount := 1

	// 正则匹配多种图片格式
	patterns := []string{
		`!\[([^\]]*)\]\(([^)]+)\)`,                                     // Markdown格式
		`https?://[^\s<>"]+?\.(?:jpg|jpeg|png|gif|webp)(?:\?[^\s<>"]*)`, // 直接URL
	}

	newMarkdown := markdown

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(markdown, -1)

		for _, match := range matches {
			var imageURL string
			if len(match) == 3 {
				// Markdown格式
				imageURL = match[2]
			} else if len(match) == 1 {
				// 直接URL
				imageURL = match[0]
			} else {
				continue
			}

			// 下载图片
			localPath, err := s.downloadImage(imageURL, savePath, imageCount, platform)
			if err != nil {
				continue
			}

			imagePaths = append(imagePaths, localPath)

			// 替换markdown中的图片链接
			newMarkdown = strings.ReplaceAll(newMarkdown, imageURL, localPath)
			imageCount++
		}
	}

	return newMarkdown, imagePaths, nil
}

// downloadImage 下载单张图片
func (s *ContentSaver) downloadImage(imageURL, savePath string, index int, platform string) (string, error) {
	// 解析URL
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return "", err
	}

	// 获取文件扩展名
	ext := filepath.Ext(parsedURL.Path)
	if ext == "" {
		ext = ".jpg"
	}

	// 生成本地文件名
	filename := fmt.Sprintf("img_%02d%s", index, ext)
	localPath := filepath.Join(savePath, filename)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		return "", err
	}

	// 根据平台设置Referer
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	if strings.Contains(platform, "小红书") {
		req.Header.Set("Referer", "https://www.xiaohongshu.com/")
	} else if strings.Contains(platform, "知乎") {
		req.Header.Set("Referer", "https://www.zhihu.com/")
	}

	// 下载图片
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载图片失败: %d", resp.StatusCode)
	}

	// 保存到本地
	out, err := os.Create(localPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return filename, nil
}

// sanitizeFilename 清理文件名中的非法字符
func sanitizeFilename(name string) string {
	// 移除或替换非法字符
	name = strings.TrimSpace(name)
	name = regexp.MustCompile(`[<>:"/\\|?*]`).ReplaceAllString(name, "_")
	name = regexp.MustCompile(`\s+`).ReplaceAllString(name, "_")

	// 限制长度
	if len(name) > 50 {
		name = name[:50]
	}

	return name
}

// CalculateHash 计算内容哈希
func CalculateHash(content string) string {
	hash := md5.Sum([]byte(content))
	return fmt.Sprintf("%x", hash)
}
