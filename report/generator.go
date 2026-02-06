package report

import (
	"competitive-analyzer/ai"
	"competitive-analyzer/models"
	"fmt"
	"strings"
	"time"
)

// ReportGenerator 报告生成器
type ReportGenerator struct {
	llmClient *ai.LLMClient
}

// NewReportGenerator 创建报告生成器
func NewReportGenerator(llmClient *ai.LLMClient) *ReportGenerator {
	return &ReportGenerator{
		llmClient: llmClient,
	}
}

// CompetitorAnalysisData 竞品分析数据
type CompetitorAnalysisData struct {
	Competitor   *models.Competitor
	ProductInfo  *ai.ProductInfo
	SWOTAnalysis *ai.SWOTAnalysis
	RawContents  []models.RawContent
}

// GenerateReport 生成完整报告
func (g *ReportGenerator) GenerateReport(data []CompetitorAnalysisData, topic string) (string, error) {
	var report strings.Builder

	// 报告头部
	report.WriteString(g.generateHeader(topic, len(data)))

	// 一、执行摘要
	report.WriteString("\n## 一、执行摘要\n\n")
	report.WriteString(g.generateExecutiveSummary(data))

	// 二、竞品概览
	report.WriteString("\n## 二、竞品概览\n\n")
	report.WriteString(g.generateCompetitorOverview(data))

	// 三、功能对比分析
	report.WriteString("\n## 三、功能对比分析\n\n")
	report.WriteString(g.generateFeatureComparison(data))

	// 四、价格策略分析
	report.WriteString("\n## 四、价格策略分析\n\n")
	report.WriteString(g.generatePricingAnalysis(data))

	// 五、SWOT分析
	report.WriteString("\n## 五、SWOT分析\n\n")
	for _, item := range data {
		report.WriteString(fmt.Sprintf("### %s\n\n", item.Competitor.Name))
		if item.SWOTAnalysis != nil {
			report.WriteString(g.formatSWOT(item.SWOTAnalysis))
		}
		report.WriteString("\n")
	}

	// 六、战略建议
	report.WriteString("\n## 六、战略建议\n\n")
	report.WriteString(g.generateStrategicSuggestions(data))

	// 附录
	report.WriteString("\n## 附录\n\n")
	report.WriteString(g.generateAppendix(data))

	return report.String(), nil
}

// generateHeader 生成报告头部
func (g *ReportGenerator) generateHeader(topic string, competitorCount int) string {
	return fmt.Sprintf(`# 竞品分析报告

**分析主题**: %s  
**分析时间**: %s  
**竞品数量**: %d个  
**生成方式**: 自动化分析

---

`, topic, time.Now().Format("2006-01-02"), competitorCount)
}

// generateExecutiveSummary 生成执行摘要
func (g *ReportGenerator) generateExecutiveSummary(data []CompetitorAnalysisData) string {
	var summary strings.Builder

	summary.WriteString("### 市场概况\n\n")
	summary.WriteString(fmt.Sprintf("本次分析共涵盖 %d 个竞品，采集了多渠道数据源，包括官网、用户评价、电商平台等。\n\n", len(data)))

	summary.WriteString("### 主要发现\n\n")
	for i, item := range data {
		if item.ProductInfo != nil {
			summary.WriteString(fmt.Sprintf("%d. **%s**: %s\n", i+1, item.Competitor.Name, item.ProductInfo.Tagline))
		}
	}
	summary.WriteString("\n")

	summary.WriteString("### 核心建议\n\n")
	summary.WriteString("1. 关注竞品的差异化功能，寻找市场空白区\n")
	summary.WriteString("2. 优化价格策略，提升性价比竞争力\n")
	summary.WriteString("3. 加强用户体验，提高用户满意度\n\n")

	return summary.String()
}

// generateCompetitorOverview 生成竞品概览
func (g *ReportGenerator) generateCompetitorOverview(data []CompetitorAnalysisData) string {
	var overview strings.Builder

	overview.WriteString("### 竞品基本信息\n\n")
	overview.WriteString("| 竞品 | 公司 | 产品定位 | 目标用户 | 发展阶段 |\n")
	overview.WriteString("|------|------|----------|----------|----------|\n")

	for _, item := range data {
		if item.ProductInfo != nil {
			targetUsers := strings.Join(item.ProductInfo.TargetUsers, ", ")
			overview.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				item.Competitor.Name,
				item.ProductInfo.Company,
				item.ProductInfo.Tagline,
				targetUsers,
				item.ProductInfo.Funding,
			))
		}
	}
	overview.WriteString("\n")

	return overview.String()
}

// generateFeatureComparison 生成功能对比
func (g *ReportGenerator) generateFeatureComparison(data []CompetitorAnalysisData) string {
	var comparison strings.Builder

	comparison.WriteString("### 功能对比矩阵\n\n")

	// 收集所有功能
	allFeatures := make(map[string]bool)
	for _, item := range data {
		if item.ProductInfo != nil {
			for _, feature := range item.ProductInfo.CoreFeatures {
				allFeatures[feature.Name] = true
			}
		}
	}

	// 生成对比表格
	comparison.WriteString("| 功能 |")
	for _, item := range data {
		comparison.WriteString(fmt.Sprintf(" %s |", item.Competitor.Name))
	}
	comparison.WriteString("\n")

	comparison.WriteString("|------|")
	for range data {
		comparison.WriteString("------|")
	}
	comparison.WriteString("\n")

	for feature := range allFeatures {
		comparison.WriteString(fmt.Sprintf("| %s |", feature))
		for _, item := range data {
			hasFeature := false
			if item.ProductInfo != nil {
				for _, f := range item.ProductInfo.CoreFeatures {
					if f.Name == feature {
						hasFeature = true
						break
					}
				}
			}
			if hasFeature {
				comparison.WriteString(" ✅ |")
			} else {
				comparison.WriteString(" ❌ |")
			}
		}
		comparison.WriteString("\n")
	}
	comparison.WriteString("\n")

	return comparison.String()
}

// generatePricingAnalysis 生成价格分析
func (g *ReportGenerator) generatePricingAnalysis(data []CompetitorAnalysisData) string {
	var analysis strings.Builder

	analysis.WriteString("### 价格体系对比\n\n")
	analysis.WriteString("| 竞品 | 定价模式 | 起步价 | 专业版 | 企业版 | 试用 |\n")
	analysis.WriteString("|------|----------|--------|--------|--------|------|\n")

	for _, item := range data {
		if item.ProductInfo != nil && len(item.ProductInfo.Pricing.Tiers) > 0 {
			pricing := item.ProductInfo.Pricing
			tiers := pricing.Tiers

			startPrice := "免费"
			if len(tiers) > 0 && tiers[0].Price > 0 {
				startPrice = fmt.Sprintf("¥%.0f/%s", tiers[0].Price, tiers[0].BillingCycle)
			}

			proPrice := "-"
			if len(tiers) > 1 {
				proPrice = fmt.Sprintf("¥%.0f/%s", tiers[1].Price, tiers[1].BillingCycle)
			}

			entPrice := "-"
			if len(tiers) > 2 {
				entPrice = fmt.Sprintf("¥%.0f/%s", tiers[2].Price, tiers[2].BillingCycle)
			}

			trial := "无"
			if pricing.Trial.Available {
				trial = pricing.Trial.Duration
			}

			analysis.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
				item.Competitor.Name,
				pricing.Model,
				startPrice,
				proPrice,
				entPrice,
				trial,
			))
		}
	}
	analysis.WriteString("\n")

	return analysis.String()
}

// formatSWOT 格式化SWOT分析
func (g *ReportGenerator) formatSWOT(swot *ai.SWOTAnalysis) string {
	var result strings.Builder

	result.WriteString("**优势 (Strengths)**\n\n")
	for _, item := range swot.Strengths {
		result.WriteString(fmt.Sprintf("- %s (影响: %s)\n  - 证据: %s\n", item.Point, item.Impact, item.Evidence))
	}
	result.WriteString("\n")

	result.WriteString("**劣势 (Weaknesses)**\n\n")
	for _, item := range swot.Weaknesses {
		result.WriteString(fmt.Sprintf("- %s (影响: %s)\n  - 证据: %s\n", item.Point, item.Impact, item.Evidence))
	}
	result.WriteString("\n")

	result.WriteString("**机会 (Opportunities)**\n\n")
	for _, item := range swot.Opportunities {
		result.WriteString(fmt.Sprintf("- %s\n  - 背景: %s\n  - 建议: %s\n", item.Point, item.Context, item.Action))
	}
	result.WriteString("\n")

	result.WriteString("**威胁 (Threats)**\n\n")
	for _, item := range swot.Threats {
		result.WriteString(fmt.Sprintf("- %s\n  - 背景: %s\n  - 应对: %s\n", item.Point, item.Context, item.Action))
	}
	result.WriteString("\n")

	result.WriteString(fmt.Sprintf("**总体评估**: %s\n\n", swot.OverallAssessment))

	return result.String()
}

// generateStrategicSuggestions 生成战略建议
func (g *ReportGenerator) generateStrategicSuggestions(data []CompetitorAnalysisData) string {
	var suggestions strings.Builder

	suggestions.WriteString("### 产品策略\n\n")
	suggestions.WriteString("1. 强化差异化功能，打造核心竞争力\n")
	suggestions.WriteString("2. 优化用户体验，降低使用门槛\n")
	suggestions.WriteString("3. 快速迭代，响应市场需求\n\n")

	suggestions.WriteString("### 定价策略\n\n")
	suggestions.WriteString("1. 采用灵活的定价模式，满足不同客户需求\n")
	suggestions.WriteString("2. 提供免费试用，降低用户决策成本\n")
	suggestions.WriteString("3. 建立清晰的价值阶梯，引导用户升级\n\n")

	suggestions.WriteString("### 营销策略\n\n")
	suggestions.WriteString("1. 多渠道布局，扩大品牌影响力\n")
	suggestions.WriteString("2. 内容营销为主，建立专业形象\n")
	suggestions.WriteString("3. 社区运营，培养用户忠诚度\n\n")

	return suggestions.String()
}

// generateAppendix 生成附录
func (g *ReportGenerator) generateAppendix(data []CompetitorAnalysisData) string {
	var appendix strings.Builder

	appendix.WriteString("### 数据来源\n\n")
	for _, item := range data {
		appendix.WriteString(fmt.Sprintf("**%s**\n", item.Competitor.Name))
		for _, content := range item.RawContents {
			if url, ok := content.Metadata["url"].(string); ok {
				appendix.WriteString(fmt.Sprintf("- %s\n", url))
			}
		}
		appendix.WriteString("\n")
	}

	appendix.WriteString(fmt.Sprintf("### 分析时间\n\n生成时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	return appendix.String()
}

// ExportToPDF 导出为PDF（TODO）
func (g *ReportGenerator) ExportToPDF(markdown string, outputPath string) error {
	// TODO: 实现Markdown转PDF
	// 可以使用 chromedp 或者调用外部工具如 pandoc
	return fmt.Errorf("PDF导出功能待实现")
}

// ExportToHTML 导出为HTML（TODO）
func (g *ReportGenerator) ExportToHTML(markdown string, outputPath string) error {
	// TODO: 实现Markdown转HTML
	return fmt.Errorf("HTML导出功能待实现")
}
