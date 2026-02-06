package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// JSONB 自定义JSON类型
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// DiscoveryTask 数据源发现任务
type DiscoveryTask struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Topic           string    `gorm:"not null" json:"topic"`
	Market          string    `json:"market"`
	TargetCount     int       `json:"target_count"`
	SearchDepth     string    `json:"search_depth"` // quick/standard/deep
	Status          string    `gorm:"default:'pending'" json:"status"`
	Progress        int       `gorm:"default:0" json:"progress"`
	CompetitorsFound int      `gorm:"default:0" json:"competitors_found"`
	SourcesFound    int       `gorm:"default:0" json:"sources_found"`
	ResultData      JSONB     `gorm:"type:text" json:"result_data"`
	CreatedAt       time.Time `json:"created_at"`
	CompletedAt     *time.Time `json:"completed_at"`
}

// SearchCache 搜索缓存
type SearchCache struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Query        string    `gorm:"uniqueIndex;not null" json:"query"`
	SearchEngine string    `json:"search_engine"`
	Results      JSONB     `gorm:"type:text" json:"results"`
	CachedAt     time.Time `json:"cached_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	HitCount     int       `gorm:"default:0" json:"hit_count"`
}

// Competitor 竞品
type Competitor struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"not null" json:"name"`
	Company          string    `json:"company"`
	Website          string    `json:"website"`
	Category         string    `json:"category"`
	DiscoveryTaskID  *uint     `json:"discovery_task_id"`
	Confidence       float64   `json:"confidence"`
	Status           string    `gorm:"default:'active'" json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// DataSource 数据源
type DataSource struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	CompetitorID    uint      `gorm:"not null" json:"competitor_id"`
	SourceType      string    `json:"source_type"` // 官网/小红书/淘宝等
	URL             string    `json:"url"`
	Priority        int       `json:"priority"`
	QualityScore    float64   `json:"quality_score"`
	AutoDiscovered  bool      `gorm:"default:false" json:"auto_discovered"`
	Status          string    `gorm:"default:'active'" json:"status"`
	LastCrawlTime   *time.Time `json:"last_crawl_time"`
	Competitor      Competitor `gorm:"foreignKey:CompetitorID" json:"competitor,omitempty"`
}

// RawContent 原始内容
type RawContent struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	SourceID    uint       `gorm:"not null" json:"source_id"`
	ContentPath string     `json:"content_path"` // Markdown文件路径
	ContentHash string     `json:"content_hash"` // 内容哈希
	CrawlTime   time.Time  `json:"crawl_time"`
	Metadata    JSONB      `gorm:"type:text" json:"metadata"`
	DataSource  DataSource `gorm:"foreignKey:SourceID" json:"data_source,omitempty"`
}

// ParsedData 解析结果
type ParsedData struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	RawContentID  uint       `gorm:"not null" json:"raw_content_id"`
	DataType      string     `json:"data_type"` // product_info/features/pricing/reviews
	ExtractedData JSONB      `gorm:"type:text" json:"extracted_data"`
	Confidence    float64    `json:"confidence"`
	ParsedAt      time.Time  `json:"parsed_at"`
	RawContent    RawContent `gorm:"foreignKey:RawContentID" json:"raw_content,omitempty"`
}

// AnalysisReport 分析报告
type AnalysisReport struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ReportName  string    `json:"report_name"`
	ReportType  string    `json:"report_type"` // full/feature_compare/price_compare
	Competitors JSONB     `gorm:"type:text" json:"competitors"`
	ReportPath  string    `json:"report_path"`
	CreatedAt   time.Time `json:"created_at"`
}

// ChangeLog 变化日志
type ChangeLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CompetitorID uint      `gorm:"not null" json:"competitor_id"`
	ChangeType   string    `json:"change_type"`
	FieldName    string    `json:"field_name"`
	OldValue     string    `gorm:"type:text" json:"old_value"`
	NewValue     string    `gorm:"type:text" json:"new_value"`
	ImpactLevel  string    `json:"impact_level"` // 高/中/低
	DetectedAt   time.Time `json:"detected_at"`
	Notified     bool      `gorm:"default:false" json:"notified"`
	Competitor   Competitor `gorm:"foreignKey:CompetitorID" json:"competitor,omitempty"`
}

// MonitorTask 监控任务
type MonitorTask struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"not null" json:"name"`
	CompetitorIDs JSONB    `gorm:"type:text" json:"competitor_ids"`
	Frequency    string    `json:"frequency"` // daily/weekly/monthly
	AlertRules   JSONB     `gorm:"type:text" json:"alert_rules"`
	Status       string    `gorm:"default:'active'" json:"status"`
	LastRunTime  *time.Time `json:"last_run_time"`
	NextRunTime  time.Time `json:"next_run_time"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
