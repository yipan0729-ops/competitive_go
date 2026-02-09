package handlers

import (
	"competitive-analyzer/ai"
	"competitive-analyzer/config"
	"competitive-analyzer/crawler"
	"competitive-analyzer/database"
	"competitive-analyzer/discovery"
	"competitive-analyzer/models"
	"competitive-analyzer/report"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
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
	TaskID              uint                `json:"task_id" binding:"required"`
	SelectedCompetitors []string            `json:"selected_competitors"`
	SelectedSources     map[string][]string `json:"selected_sources"`
	SaveAsConfig        bool                `json:"save_as_config"`
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

// CrawlBatchRequest 批量爬取请求
type CrawlBatchRequest struct {
	URLs       []URLItem `json:"urls" binding:"required"`
	Concurrent int       `json:"concurrent"` // 并发数，默认3
}

// URLItem URL项
type URLItem struct {
	URL        string `json:"url"`
	Competitor string `json:"competitor"`
	SourceType string `json:"source_type"`
}

// CrawlBatch 批量爬取
func (h *CrawlHandler) CrawlBatch(c *gin.Context) {
	var req CrawlBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认并发数（降低以避免403错误）
	concurrent := req.Concurrent
	if concurrent <= 0 {
		concurrent = 1 // 改为1，避免并发请求被封
	}
	if concurrent > 3 {
		concurrent = 3 // 降低最大并发，避免触发反爬虫
	}

	// 异步执行批量爬取
	go h.executeBatchCrawl(req.URLs, concurrent)

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"total_urls": len(req.URLs),
		"concurrent": concurrent,
		"message":    "批量爬取任务已启动",
	})
}

// executeBatchCrawl 执行批量爬取
func (h *CrawlHandler) executeBatchCrawl(urls []URLItem, concurrent int) {
	db := database.DB

	// 使用信号量控制并发
	sem := make(chan struct{}, concurrent)
	var wg sync.WaitGroup

	for i, urlItem := range urls {
		wg.Add(1)
		go func(index int, item URLItem) {
			defer wg.Done()
			sem <- struct{}{}        // 获取信号量
			defer func() { <-sem }() // 释放信号量

			// 添加渐进式延迟，避免同时发起过多请求
			time.Sleep(time.Duration(index) * 2 * time.Second)

			// 爬取（带重试）
			var result *crawler.CrawlResult
			var err error

			// 最多重试3次
			for retry := 0; retry < 3; retry++ {
				result, err = h.crawler.Crawl(item.URL)
				if err == nil {
					break // 成功，退出重试
				}

				if retry < 2 {
					// 失败，等待后重试
					waitTime := time.Duration(retry+1) * 5 * time.Second
					log.Printf("爬取失败 %s (重试 %d/2): %v，等待%v后重试...", item.URL, retry+1, err, waitTime)
					time.Sleep(waitTime)
				}
			}

			if err != nil {
				log.Printf("爬取最终失败 %s: %v", item.URL, err)
				return
			}

			// 保存
			saveResult, err := h.saver.Save(result, item.Competitor)
			if err != nil {
				log.Printf("保存失败 %s: %v", item.URL, err)
				return
			}

			// 保存到数据库
			var competitor models.Competitor
			db.FirstOrCreate(&competitor, models.Competitor{Name: item.Competitor})

			var dataSource models.DataSource
			db.FirstOrCreate(&dataSource, models.DataSource{
				CompetitorID: competitor.ID,
				URL:          item.URL,
				SourceType:   item.SourceType,
			})
			now := time.Now()
			dataSource.LastCrawlTime = &now
			db.Save(&dataSource)

			rawContent := &models.RawContent{
				SourceID:    dataSource.ID,
				ContentPath: saveResult.ContentPath,
				ContentHash: crawler.CalculateHash(result.Markdown),
				CrawlTime:   time.Now(),
				Metadata: models.JSONB{
					"title":    result.Title,
					"platform": result.Platform,
					"method":   result.Method,
					"url":      item.URL,
				},
			}
			db.Create(rawContent)

			log.Printf("爬取成功: %s", item.URL)
		}(i, urlItem)
	}

	wg.Wait()
	log.Println("批量爬取任务完成")
}

// AnalysisHandler AI分析处理器
type AnalysisHandler struct {
	llmClient            *ai.LLMClient
	productInfoExtractor *ai.ProductInfoExtractor
	swotAnalyzer         *ai.SWOTAnalyzer
}

// NewAnalysisHandler 创建AI分析处理器
func NewAnalysisHandler() *AnalysisHandler {
	cfg := config.AppConfig
	llmClient := ai.NewLLMClient(cfg.OpenAIAPIKey, cfg.LLMModel, cfg.LLMTemperature, cfg.LLMMaxTokens, cfg.LLMBaseURL)

	return &AnalysisHandler{
		llmClient:            llmClient,
		productInfoExtractor: ai.NewProductInfoExtractor(llmClient),
		swotAnalyzer:         ai.NewSWOTAnalyzer(llmClient),
	}
}

// AnalyzeCompetitorRequest 竞品分析请求
type AnalyzeCompetitorRequest struct {
	CompetitorID  uint   `json:"competitor_id" binding:"required"`
	MarketContext string `json:"market_context"` // 市场背景
}

// AnalyzeCompetitor 分析单个竞品
func (h *AnalysisHandler) AnalyzeCompetitor(c *gin.Context) {
	var req AnalyzeCompetitorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.DB

	// 获取竞品信息
	var competitor models.Competitor
	if err := db.First(&competitor, req.CompetitorID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "竞品不存在"})
		return
	}

	// 获取竞品的所有数据源
	var dataSources []models.DataSource
	db.Where("competitor_id = ?", req.CompetitorID).Find(&dataSources)

	if len(dataSources) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该竞品没有数据源，请先爬取内容"})
		return
	}

	// 获取原始内容
	var rawContents []models.RawContent
	sourceIDs := make([]uint, len(dataSources))
	for i, ds := range dataSources {
		sourceIDs[i] = ds.ID
	}
	db.Where("source_id IN ?", sourceIDs).Find(&rawContents)

	// 合并所有内容
	var allContent strings.Builder
	for _, rc := range rawContents {
		// 读取文件内容
		if rc.ContentPath != "" {
			contentBytes, err := os.ReadFile(rc.ContentPath)
			if err == nil {
				allContent.WriteString(string(contentBytes))
				allContent.WriteString("\n\n---\n\n")
			}
		}
	}

	content := allContent.String()
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有可分析的内容"})
		return
	}

	// 限制内容长度（防止Token超限）
	if len(content) > 50000 {
		content = content[:50000]
	}

	// 提取产品信息
	productInfo, err := h.productInfoExtractor.Extract(content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "产品信息提取失败: " + err.Error()})
		return
	}

	// SWOT分析
	swotAnalysis, err := h.swotAnalyzer.Analyze(competitor.Name, productInfo, req.MarketContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "SWOT分析失败: " + err.Error()})
		return
	}

	// 保存分析结果
	productInfoJSON, _ := json.Marshal(productInfo)
	swotJSON, _ := json.Marshal(swotAnalysis)

	parsedData := &models.ParsedData{
		RawContentID: rawContents[0].ID,
		DataType:     "product_info",
		ExtractedData: models.JSONB{
			"product_info":  string(productInfoJSON),
			"swot_analysis": string(swotJSON),
		},
		Confidence: 0.8,
		ParsedAt:   time.Now(),
	}
	db.Create(parsedData)

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"competitor":    competitor.Name,
		"product_info":  productInfo,
		"swot_analysis": swotAnalysis,
	})
}

// ReportHandler 报告处理器
type ReportHandler struct {
	reportGenerator *report.ReportGenerator
}

// NewReportHandler 创建报告处理器
func NewReportHandler() *ReportHandler {
	cfg := config.AppConfig
	llmClient := ai.NewLLMClient(cfg.OpenAIAPIKey, cfg.LLMModel, cfg.LLMTemperature, cfg.LLMMaxTokens, cfg.LLMBaseURL)

	return &ReportHandler{
		reportGenerator: report.NewReportGenerator(llmClient),
	}
}

// GenerateReportRequest 生成报告请求
type GenerateReportRequest struct {
	CompetitorIDs []uint `json:"competitor_ids" binding:"required"`
	Topic         string `json:"topic" binding:"required"`
	ReportName    string `json:"report_name"`
}

// GenerateReport 生成竞品报告
func (h *ReportHandler) GenerateReport(c *gin.Context) {
	var req GenerateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.DB

	// 收集所有竞品的分析数据
	var analysisData []report.CompetitorAnalysisData

	for _, competitorID := range req.CompetitorIDs {
		var competitor models.Competitor
		if err := db.First(&competitor, competitorID).Error; err != nil {
			continue
		}

		// 获取分析结果
		var parsedDataList []models.ParsedData
		db.Joins("JOIN raw_contents ON raw_contents.id = parsed_data.raw_content_id").
			Joins("JOIN data_sources ON data_sources.id = raw_contents.source_id").
			Where("data_sources.competitor_id = ?", competitorID).
			Find(&parsedDataList)

		if len(parsedDataList) == 0 {
			continue
		}

		// 解析产品信息和SWOT
		var productInfo *ai.ProductInfo
		var swotAnalysis *ai.SWOTAnalysis

		for _, pd := range parsedDataList {
			if productInfoStr, ok := pd.ExtractedData["product_info"].(string); ok {
				json.Unmarshal([]byte(productInfoStr), &productInfo)
			}
			if swotStr, ok := pd.ExtractedData["swot_analysis"].(string); ok {
				json.Unmarshal([]byte(swotStr), &swotAnalysis)
			}
		}

		// 获取原始内容
		var rawContents []models.RawContent
		db.Joins("JOIN data_sources ON data_sources.id = raw_contents.source_id").
			Where("data_sources.competitor_id = ?", competitorID).
			Find(&rawContents)

		analysisData = append(analysisData, report.CompetitorAnalysisData{
			Competitor:   &competitor,
			ProductInfo:  productInfo,
			SWOTAnalysis: swotAnalysis,
			RawContents:  rawContents,
		})
	}

	if len(analysisData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有可用的分析数据"})
		return
	}

	// 生成报告
	reportContent, err := h.reportGenerator.GenerateReport(analysisData, req.Topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "报告生成失败: " + err.Error()})
		return
	}

	// 保存报告
	reportName := req.ReportName
	if reportName == "" {
		reportName = fmt.Sprintf("%s_竞品分析报告_%s", req.Topic, time.Now().Format("20060102"))
	}

	reportPath := fmt.Sprintf("reports/%s.md", reportName)
	os.MkdirAll("reports", 0755)

	if err := os.WriteFile(reportPath, []byte(reportContent), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "报告保存失败: " + err.Error()})
		return
	}

	// 保存到数据库
	competitorIDs := req.CompetitorIDs
	competitorIDsJSON, _ := json.Marshal(competitorIDs)

	analysisReport := &models.AnalysisReport{
		ReportName:  reportName,
		ReportType:  "competitive_analysis",
		Competitors: models.JSONB{"ids": competitorIDsJSON},
		ReportPath:  reportPath,
		CreatedAt:   time.Now(),
	}
	db.Create(analysisReport)

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"report_id":   analysisReport.ID,
		"report_name": reportName,
		"report_path": reportPath,
		"competitors": len(analysisData),
	})
}

// AutomationHandler 全流程自动化处理器
type AutomationHandler struct {
	discoveryHandler *DiscoveryHandler
	crawlHandler     *CrawlHandler
	analysisHandler  *AnalysisHandler
	reportHandler    *ReportHandler
}

// NewAutomationHandler 创建自动化处理器
func NewAutomationHandler() *AutomationHandler {
	return &AutomationHandler{
		discoveryHandler: NewDiscoveryHandler(),
		crawlHandler:     NewCrawlHandler(),
		analysisHandler:  NewAnalysisHandler(),
		reportHandler:    NewReportHandler(),
	}
}

// AutoAnalysisRequest 自动分析请求
type AutoAnalysisRequest struct {
	Topic           string `json:"topic" binding:"required"`
	Market          string `json:"market"`
	CompetitorCount int    `json:"competitor_count"`
	Depth           string `json:"depth"`
	AutoCrawl       bool   `json:"auto_crawl"`      // 是否自动爬取，默认true
	AutoAnalyze     bool   `json:"auto_analyze"`    // 是否自动分析，默认true
	GenerateReport  bool   `json:"generate_report"` // 是否生成报告，默认true
}

// AutoAnalysis 全流程自动化分析
func (h *AutomationHandler) AutoAnalysis(c *gin.Context) {
	var req AutoAnalysisRequest
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

	// 修复：如果用户没有显式设置bool参数，默认启用所有功能
	// Go的JSON unmarshal会将未设置的bool字段设为false
	// 为了提供更好的用户体验，当三个参数都是false时（即用户可能没有设置），默认全部启用
	if !req.AutoCrawl && !req.AutoAnalyze && !req.GenerateReport {
		req.AutoCrawl = true
		req.AutoAnalyze = true
		req.GenerateReport = true
		log.Println("[自动化] 使用默认配置：启用爬取、分析和报告生成")
	}

	// 创建自动化任务记录
	db := database.DB
	task := &models.DiscoveryTask{
		Topic:       req.Topic,
		Market:      req.Market,
		TargetCount: req.CompetitorCount,
		SearchDepth: req.Depth,
		Status:      "processing",
		Progress:    0,
		CreatedAt:   time.Now(),
	}
	db.Create(task)

	// 异步执行全流程
	go h.executeAutoWorkflow(task.ID, req)

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"task_id":        task.ID,
		"status":         "processing",
		"workflow":       "discovery -> crawl -> analysis -> report",
		"estimated_time": 12000, // 预计10分钟
	})
}

// executeAutoWorkflow 执行自动化工作流
func (h *AutomationHandler) executeAutoWorkflow(taskID uint, req AutoAnalysisRequest) {
	db := database.DB

	var task models.DiscoveryTask
	db.First(&task, taskID)

	log.Printf("[自动化] 开始执行任务 #%d: %s", taskID, req.Topic)

	// 步骤1: 发现竞品
	task.Status = "discovering"
	task.Progress = 10
	db.Save(&task)

	searchResults, err := h.discoveryHandler.searchManager.SearchCompetitors(req.Topic, 10)
	if err != nil {
		task.Status = "failed"
		task.ResultData = models.JSONB{"error": "发现失败: " + err.Error()}
		db.Save(&task)
		return
	}

	// 提取竞品名称（简化处理）
	competitorNames := extractCompetitorNamesFromResults(searchResults, req.CompetitorCount)

	task.Progress = 30
	task.CompetitorsFound = len(competitorNames)
	db.Save(&task)

	log.Printf("[自动化] 发现竞品: %v", competitorNames)

	// 步骤2: 搜索数据源
	task.Status = "searching_sources"
	task.Progress = 40
	db.Save(&task)

	allURLs := []URLItem{}
	for _, competitorName := range competitorNames {
		sources, _ := h.discoveryHandler.searchManager.SearchDataSources(competitorName, []string{"官网", "产品功能"}, 3)

		for _, results := range sources {
			processed := discovery.ProcessSearchResults(results)
			for _, source := range processed {
				if len(allURLs) < req.CompetitorCount*3 { // 每个竞品最多3个URL
					allURLs = append(allURLs, URLItem{
						URL:        source.URL,
						Competitor: competitorName,
						SourceType: source.Type,
					})
				}
			}
		}
	}

	task.Progress = 50
	task.SourcesFound = len(allURLs)
	db.Save(&task)

	log.Printf("[自动化] 找到数据源: %d 个", len(allURLs))

	// 步骤3: 批量爬取（如果启用）
	if req.AutoCrawl && len(allURLs) > 0 {
		task.Status = "crawling"
		task.Progress = 60
		db.Save(&task)

		log.Printf("[自动化] 开始爬取 %d 个URL", len(allURLs))
		h.crawlHandler.executeBatchCrawl(allURLs, 1) // 使用并发1，避免403错误

		task.Progress = 75
		db.Save(&task)
	}

	// 步骤4: AI分析（如果启用）
	var analyzedCompetitorIDs []uint
	if req.AutoAnalyze {
		task.Status = "analyzing"
		task.Progress = 80
		db.Save(&task)

		log.Println("[自动化] 开始AI分析")

		for _, name := range competitorNames {
			var competitor models.Competitor
			if err := db.Where("name = ?", name).First(&competitor).Error; err != nil {
				continue
			}

			// 执行分析
			if err := h.analyzeCompetitorByID(competitor.ID, req.Market); err != nil {
				log.Printf("[自动化] 分析失败 %s: %v", name, err)
				continue
			}

			analyzedCompetitorIDs = append(analyzedCompetitorIDs, competitor.ID)
		}

		task.Progress = 90
		db.Save(&task)
	}

	// 步骤5: 生成报告（如果启用）
	var reportPath string
	if req.GenerateReport && len(analyzedCompetitorIDs) > 0 {
		task.Status = "generating_report"
		task.Progress = 95
		db.Save(&task)

		log.Println("[自动化] 生成报告")

		reportPath, err = h.generateReportForCompetitors(analyzedCompetitorIDs, req.Topic)
		if err != nil {
			log.Printf("[自动化] 报告生成失败: %v", err)
		}
	}

	// 完成
	task.Status = "completed"
	task.Progress = 100
	now := time.Now()
	task.CompletedAt = &now
	task.ResultData = models.JSONB{
		"competitors":    competitorNames,
		"urls_crawled":   len(allURLs),
		"analyzed_count": len(analyzedCompetitorIDs),
		"report_path":    reportPath,
	}
	db.Save(&task)

	log.Printf("[自动化] 任务完成 #%d", taskID)
}

// extractCompetitorNamesFromResults 从搜索结果提取竞品名称
func extractCompetitorNamesFromResults(results []discovery.SearchResult, limit int) []string {
	names := make(map[string]bool)
	nameList := []string{}

	for _, result := range results {
		if len(names) >= limit {
			break
		}
		
		// 改进：尝试从URL提取品牌名
		// 例如：https://www.canva.com → Canva
		// 例如：https://www.notion.so → Notion
		brandName := extractBrandFromURL(result.URL)
		
		// 如果能从URL提取品牌名，优先使用
		if brandName != "" && !names[brandName] {
			names[brandName] = true
			nameList = append(nameList, brandName)
		} else if !names[result.Title] && result.Title != "" {
			// 否则使用标题（临时方案）
			names[result.Title] = true
			nameList = append(nameList, result.Title)
		}
	}

	return nameList
}

// extractBrandFromURL 从URL提取品牌名
func extractBrandFromURL(urlStr string) string {
	// 简单的品牌名提取逻辑
	// www.canva.com → Canva
	// www.notion.so → Notion
	
	parts := strings.Split(urlStr, "//")
	if len(parts) < 2 {
		return ""
	}
	
	domain := strings.Split(parts[1], "/")[0]
	domain = strings.TrimPrefix(domain, "www.")
	domain = strings.Split(domain, ".")[0]
	
	// 首字母大写
	if len(domain) > 0 {
		return strings.ToUpper(domain[:1]) + domain[1:]
	}
	
	return ""
}

// analyzeCompetitorByID 分析竞品（内部方法）
func (h *AutomationHandler) analyzeCompetitorByID(competitorID uint, marketContext string) error {
	db := database.DB

	var competitor models.Competitor
	if err := db.First(&competitor, competitorID).Error; err != nil {
		return err
	}

	// 获取数据源和内容
	var dataSources []models.DataSource
	db.Where("competitor_id = ?", competitorID).Find(&dataSources)

	if len(dataSources) == 0 {
		return fmt.Errorf("没有数据源")
	}

	var rawContents []models.RawContent
	sourceIDs := make([]uint, len(dataSources))
	for i, ds := range dataSources {
		sourceIDs[i] = ds.ID
	}
	db.Where("source_id IN ?", sourceIDs).Find(&rawContents)

	// 合并内容
	var allContent strings.Builder
	for _, rc := range rawContents {
		if rc.ContentPath != "" {
			contentBytes, _ := os.ReadFile(rc.ContentPath)
			allContent.WriteString(string(contentBytes))
			allContent.WriteString("\n\n---\n\n")
		}
	}

	content := allContent.String()
	if len(content) > 50000 {
		content = content[:50000]
	}

	// 提取产品信息
	productInfo, err := h.analysisHandler.productInfoExtractor.Extract(content)
	if err != nil {
		return err
	}

	// SWOT分析
	swotAnalysis, err := h.analysisHandler.swotAnalyzer.Analyze(competitor.Name, productInfo, marketContext)
	if err != nil {
		return err
	}

	// 保存分析结果
	productInfoJSON, _ := json.Marshal(productInfo)
	swotJSON, _ := json.Marshal(swotAnalysis)

	parsedData := &models.ParsedData{
		RawContentID: rawContents[0].ID,
		DataType:     "product_info",
		ExtractedData: models.JSONB{
			"product_info":  string(productInfoJSON),
			"swot_analysis": string(swotJSON),
		},
		Confidence: 0.8,
		ParsedAt:   time.Now(),
	}
	db.Create(parsedData)

	return nil
}

// generateReportForCompetitors 为竞品生成报告（内部方法）
func (h *AutomationHandler) generateReportForCompetitors(competitorIDs []uint, topic string) (string, error) {
	db := database.DB

	var analysisData []report.CompetitorAnalysisData

	for _, competitorID := range competitorIDs {
		var competitor models.Competitor
		if err := db.First(&competitor, competitorID).Error; err != nil {
			continue
		}

		var parsedDataList []models.ParsedData
		db.Joins("JOIN raw_contents ON raw_contents.id = parsed_data.raw_content_id").
			Joins("JOIN data_sources ON data_sources.id = raw_contents.source_id").
			Where("data_sources.competitor_id = ?", competitorID).
			Find(&parsedDataList)

		var productInfo *ai.ProductInfo
		var swotAnalysis *ai.SWOTAnalysis

		for _, pd := range parsedDataList {
			if productInfoStr, ok := pd.ExtractedData["product_info"].(string); ok {
				json.Unmarshal([]byte(productInfoStr), &productInfo)
			}
			if swotStr, ok := pd.ExtractedData["swot_analysis"].(string); ok {
				json.Unmarshal([]byte(swotStr), &swotAnalysis)
			}
		}

		var rawContents []models.RawContent
		db.Joins("JOIN data_sources ON data_sources.id = raw_contents.source_id").
			Where("data_sources.competitor_id = ?", competitorID).
			Find(&rawContents)

		analysisData = append(analysisData, report.CompetitorAnalysisData{
			Competitor:   &competitor,
			ProductInfo:  productInfo,
			SWOTAnalysis: swotAnalysis,
			RawContents:  rawContents,
		})
	}

	reportContent, err := h.reportHandler.reportGenerator.GenerateReport(analysisData, topic)
	if err != nil {
		return "", err
	}

	reportName := fmt.Sprintf("%s_自动分析报告_%s", topic, time.Now().Format("20060102_150405"))
	reportPath := fmt.Sprintf("reports/%s.md", reportName)
	os.MkdirAll("reports", 0755)

	if err := os.WriteFile(reportPath, []byte(reportContent), 0644); err != nil {
		return "", err
	}

	// 保存到数据库
	competitorIDsJSON, _ := json.Marshal(competitorIDs)
	analysisReport := &models.AnalysisReport{
		ReportName:  reportName,
		ReportType:  "auto_competitive_analysis",
		Competitors: models.JSONB{"ids": competitorIDsJSON},
		ReportPath:  reportPath,
		CreatedAt:   time.Now(),
	}
	db.Create(analysisReport)

	return reportPath, nil
}
