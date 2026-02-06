package handlers

import (
	"competitive-analyzer/ai"
	"competitive-analyzer/config"
	"competitive-analyzer/crawler"
	"competitive-analyzer/database"
	"competitive-analyzer/discovery"
	"competitive-analyzer/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// DiscoveryHandler 数据源发现处理器
type DiscoveryHandler struct {
	searchManager *discovery.SearchManager
	llmClient     *ai.LLMClient
}

// NewDiscoveryHandler 创建处理器
func NewDiscoveryHandler() *DiscoveryHandler {
	cfg := config.AppConfig

	// 初始化搜索引擎
	engines := []discovery.SearchEngine{}

	if cfg.SerperAPIKey != "" {
		engines = append(engines, &discovery.SerperSearchEngine{APIKey: cfg.SerperAPIKey})
	}
	if cfg.GoogleAPIKey != "" && cfg.GoogleEngineID != "" {
		engines = append(engines, &discovery.GoogleSearchEngine{
			APIKey:   cfg.GoogleAPIKey,
			EngineID: cfg.GoogleEngineID,
		})
	}
	if cfg.BingAPIKey != "" {
		engines = append(engines, &discovery.BingSearchEngine{APIKey: cfg.BingAPIKey})
	}

	searchManager := discovery.NewSearchManager(engines)
	llmClient := ai.NewLLMClient(cfg.OpenAIAPIKey, cfg.LLMModel, cfg.LLMTemperature, cfg.LLMMaxTokens, cfg.LLMBaseURL)

	return &DiscoveryHandler{
		searchManager: searchManager,
		llmClient:     llmClient,
	}
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Topic           string   `json:"topic" binding:"required"`
	Market          string   `json:"market"`
	CompetitorCount int      `json:"competitor_count"`
	SourceTypes     []string `json:"source_types"`
	Depth           string   `json:"depth"` // quick/standard/deep
}

// Search 开始数据源发现
func (h *DiscoveryHandler) Search(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认值
	if req.Depth == "" {
		req.Depth = "standard"
	}
	if req.CompetitorCount == 0 {
		req.CompetitorCount = 5
	}

	// 创建任务
	task := &models.DiscoveryTask{
		Topic:       req.Topic,
		Market:      req.Market,
		TargetCount: req.CompetitorCount,
		SearchDepth: req.Depth,
		Status:      "processing",
		Progress:    0,
		CreatedAt:   time.Now(),
	}

	if err := database.DB.Create(task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建任务失败"})
		return
	}

	// 异步执行搜索
	go h.executeSearch(task, req)

	c.JSON(http.StatusOK, gin.H{
		"task_id":        task.ID,
		"status":         "processing",
		"progress":       0,
		"estimated_time": 60,
	})
}

// executeSearch 执行搜索（后台任务）
func (h *DiscoveryHandler) executeSearch(task *models.DiscoveryTask, req SearchRequest) {
	db := database.DB

	// 更新进度：10% - 开始搜索竞品
	task.Progress = 10
	db.Save(task)

	// 1. 搜索竞品
	maxResults := 10
	if req.Depth == "deep" {
		maxResults = 20
	} else if req.Depth == "quick" {
		maxResults = 5
	}

	searchResults, err := h.searchManager.SearchCompetitors(req.Topic, maxResults)
	if err != nil {
		task.Status = "failed"
		task.ResultData = models.JSONB{"error": err.Error()}
		db.Save(task)
		return
	}

	// 更新进度：40% - 提取竞品名称
	task.Progress = 40
	db.Save(task)

	// 2. 从搜索结果中提取竞品（这里简化处理，实际应该访问文章用LLM提取）
	// 为了演示，我们直接从标题中提取可能的竞品名称
	competitorNames := h.extractCompetitorNames(searchResults, req.Topic)

	if len(competitorNames) > req.CompetitorCount {
		competitorNames = competitorNames[:req.CompetitorCount]
	}

	// 更新进度：60% - 搜索数据源
	task.Progress = 60
	task.CompetitorsFound = len(competitorNames)
	db.Save(task)

	// 3. 为每个竞品搜索数据源
	allDataSources := make(map[string][]*discovery.DataSourceInfo)

	for _, competitorName := range competitorNames {
		sources, err := h.searchManager.SearchDataSources(competitorName, req.SourceTypes, 5)
		if err != nil {
			continue
		}

		// 处理和评分
		for sourceType, results := range sources {
			processed := discovery.ProcessSearchResults(results)
			key := fmt.Sprintf("%s_%s", competitorName, sourceType)
			allDataSources[key] = processed
		}
	}

	// 统计数据源总数
	totalSources := 0
	for _, sources := range allDataSources {
		totalSources += len(sources)
	}

	// 更新进度：100% - 完成
	task.Progress = 100
	task.Status = "completed"
	task.CompetitorsFound = len(competitorNames)
	task.SourcesFound = totalSources
	now := time.Now()
	task.CompletedAt = &now

	// 保存结果
	task.ResultData = models.JSONB{
		"competitors":  competitorNames,
		"data_sources": allDataSources,
	}

	db.Save(task)
}

// extractCompetitorNames 从搜索结果中提取竞品名称（简化版）
func (h *DiscoveryHandler) extractCompetitorNames(results []discovery.SearchResult, topic string) []string {
	// 这里应该用LLM提取，为了简化，我们从标题中提取
	// 实际应该访问对比类文章，用LLM提取竞品列表
	names := make(map[string]bool)

	for _, result := range results {
		// 简单的关键词提取逻辑
		// 实际项目中应该用LLM
		if len(names) >= 10 {
			break
		}
		// 这里只是占位符，实际需要更智能的提取
		names[result.Title] = true
	}

	nameList := []string{}
	for name := range names {
		nameList = append(nameList, name)
	}

	return nameList
}

// GetStatus 获取任务状态
func (h *DiscoveryHandler) GetStatus(c *gin.Context) {
	taskID := c.Param("task_id")

	var task models.DiscoveryTask
	if err := database.DB.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":             task.Status,
		"progress":           task.Progress,
		"competitors_found":  task.CompetitorsFound,
		"data_sources_found": task.SourcesFound,
		"result":             task.ResultData,
		"created_at":         task.CreatedAt,
		"completed_at":       task.CompletedAt,
	})
}

// ConfirmRequest 确认请求
type ConfirmRequest struct {
	TaskID              uint                            `json:"task_id" binding:"required"`
	SelectedCompetitors []string                        `json:"selected_competitors"`
	SelectedSources     map[string][]string             `json:"selected_sources"`
	SaveAsConfig        bool                            `json:"save_as_config"`
}

// Confirm 用户确认并保存配置
func (h *DiscoveryHandler) Confirm(c *gin.Context) {
	var req ConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.DB

	// 获取任务
	var task models.DiscoveryTask
	if err := db.First(&task, req.TaskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	// 保存竞品和数据源到数据库
	for _, competitorName := range req.SelectedCompetitors {
		// 创建竞品记录
		competitor := &models.Competitor{
			Name:            competitorName,
			DiscoveryTaskID: &task.ID,
			Confidence:      0.9,
			Status:          "active",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		if err := db.Create(competitor).Error; err != nil {
			continue
		}

		// 保存数据源
		if sources, ok := req.SelectedSources[competitorName]; ok {
			for _, sourceURL := range sources {
				dataSource := &models.DataSource{
					CompetitorID:   competitor.ID,
					URL:            sourceURL,
					Priority:       1,
					QualityScore:   0.8,
					AutoDiscovered: true,
					Status:         "active",
				}
				db.Create(dataSource)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "配置已保存",
	})
}

// CrawlHandler 爬取处理器
type CrawlHandler struct {
	crawler *crawler.ThreeLayerCrawler
	saver   *crawler.ContentSaver
}

// NewCrawlHandler 创建爬取处理器
func NewCrawlHandler() *CrawlHandler {
	cfg := config.AppConfig
	return &CrawlHandler{
		crawler: crawler.NewThreeLayerCrawler(cfg.FirecrawlAPIKey),
		saver:   crawler.NewContentSaver(cfg.StoragePath),
	}
}

// CrawlSingleRequest 单个爬取请求
type CrawlSingleRequest struct {
	URL        string `json:"url" binding:"required"`
	Competitor string `json:"competitor" binding:"required"`
	SourceType string `json:"source_type"`
}

// CrawlSingle 爬取单个URL
func (h *CrawlHandler) CrawlSingle(c *gin.Context) {
	var req CrawlSingleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 爬取
	result, err := h.crawler.Crawl(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 保存
	saveResult, err := h.saver.Save(result, req.Competitor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败: " + err.Error()})
		return
	}

	// 保存到数据库
	db := database.DB

	// 查找或创建竞品
	var competitor models.Competitor
	db.FirstOrCreate(&competitor, models.Competitor{Name: req.Competitor})

	// 查找或创建数据源
	var dataSource models.DataSource
	db.FirstOrCreate(&dataSource, models.DataSource{
		CompetitorID: competitor.ID,
		URL:          req.URL,
	})
	now := time.Now()
	dataSource.LastCrawlTime = &now
	db.Save(&dataSource)

	// 保存原始内容
	rawContent := &models.RawContent{
		SourceID:    dataSource.ID,
		ContentPath: saveResult.ContentPath,
		ContentHash: crawler.CalculateHash(result.Markdown),
		CrawlTime:   time.Now(),
		Metadata: models.JSONB{
			"title":    result.Title,
			"platform": result.Platform,
			"method":   result.Method,
		},
	}
	db.Create(rawContent)

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"content_path": saveResult.ContentPath,
		"image_count":  len(saveResult.ImagePaths),
		"title":        saveResult.Title,
	})
}

// GetCompetitors 获取竞品列表
func GetCompetitors(c *gin.Context) {
	db := database.DB

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var competitors []models.Competitor
	var total int64

	db.Model(&models.Competitor{}).Count(&total)
	db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&competitors)

	c.JSON(http.StatusOK, gin.H{
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"competitors": competitors,
	})
}

// GetDataSources 获取数据源列表
func GetDataSources(c *gin.Context) {
	db := database.DB

	competitorID := c.Query("competitor_id")

	var dataSources []models.DataSource
	query := db.Preload("Competitor")

	if competitorID != "" {
		query = query.Where("competitor_id = ?", competitorID)
	}

	query.Find(&dataSources)

	c.JSON(http.StatusOK, gin.H{
		"data_sources": dataSources,
	})
}
