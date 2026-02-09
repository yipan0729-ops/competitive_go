package main

import (
	"competitive-analyzer/config"
	"competitive-analyzer/database"
	"competitive-analyzer/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	if err := database.InitDB(cfg.DBPath); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 设置Gin模式
	gin.SetMode(cfg.GinMode)

	// 创建路由
	r := gin.Default()

	// CORS中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"version": "1.0.0",
		})
	})

	// API路由组
	api := r.Group("/api")
	{
		// 数据源发现模块
		discoveryHandler := handlers.NewDiscoveryHandler()
		discover := api.Group("/discover")
		{
			discover.POST("/search", discoveryHandler.Search)
			discover.GET("/status/:task_id", discoveryHandler.GetStatus)
			discover.POST("/confirm", discoveryHandler.Confirm)
		}

		// 爬取模块
		crawlHandler := handlers.NewCrawlHandler()
		crawl := api.Group("/crawl")
		{
			crawl.POST("/single", crawlHandler.CrawlSingle)
			crawl.POST("/batch", crawlHandler.CrawlBatch) // 批量爬取
		}

		// 竞品管理
		competitors := api.Group("/competitors")
		{
			competitors.GET("", handlers.GetCompetitors)
			competitors.GET("/:id/sources", handlers.GetDataSources)
		}

		// AI分析模块
		analysisHandler := handlers.NewAnalysisHandler()
		analyze := api.Group("/analyze")
		{
			analyze.POST("/competitor", analysisHandler.AnalyzeCompetitor)
		}

		// 报告模块
		reportHandler := handlers.NewReportHandler()
		reportAPI := api.Group("/report")
		{
			reportAPI.POST("/generate", reportHandler.GenerateReport)
		}

		// 全流程自动化
		automationHandler := handlers.NewAutomationHandler()
		auto := api.Group("/auto")
		{
			auto.POST("/analysis", automationHandler.AutoAnalysis) // 一键分析
		}
	}

	// 启动服务器
	log.Printf("服务器启动在端口 %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
