package ai

import (
	"encoding/json"
	"fmt"
)

// CompetitorExtractor 竞品提取器
type CompetitorExtractor struct {
	llmClient *LLMClient
}

// CompetitorInfo 竞品信息
type CompetitorInfo struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
	Reason     string  `json:"reason"`
}

// NewCompetitorExtractor 创建竞品提取器
func NewCompetitorExtractor(llmClient *LLMClient) *CompetitorExtractor {
	return &CompetitorExtractor{
		llmClient: llmClient,
	}
}

// ExtractCompetitors 从内容中提取竞品
func (e *CompetitorExtractor) ExtractCompetitors(topic, content string) ([]CompetitorInfo, error) {
	systemPrompt := `你是一位专业的市场研究分析师。请从提供的文章内容中提取所有相关产品/工具的名称。

要求：
1. 只提取明确提到的产品名称，不要臆测
2. 排除通用名词（如"AI工具"、"软件"等）
3. 每个产品给出置信度评分（0-1）
4. 按照JSON格式输出

输出格式：
{
    "competitors": [
        {"name": "产品名", "confidence": 0.95, "reason": "文中明确提到为XX工具"},
        {"name": "产品名2", "confidence": 0.80, "reason": "推测可能相关"}
    ]
}`

	userPrompt := fmt.Sprintf(`主题：%s

文章内容：
%s

请提取相关竞品。`, topic, content)

	response, err := e.llmClient.CompletionWithJSON(systemPrompt, userPrompt)
	if err != nil {
		return nil, err
	}

	// 解析结果
	competitorsData, ok := response["competitors"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("响应格式错误")
	}

	competitors := []CompetitorInfo{}
	for _, item := range competitorsData {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		name, _ := itemMap["name"].(string)
		confidence, _ := itemMap["confidence"].(float64)
		reason, _ := itemMap["reason"].(string)

		competitors = append(competitors, CompetitorInfo{
			Name:       name,
			Confidence: confidence,
			Reason:     reason,
		})
	}

	return competitors, nil
}

// ProductInfoExtractor 产品信息提取器
type ProductInfoExtractor struct {
	llmClient *LLMClient
}

// ProductInfo 产品信息
type ProductInfo struct {
	ProductName   string              `json:"product_name"`
	Company       string              `json:"company"`
	Tagline       string              `json:"tagline"`
	TargetUsers   []string            `json:"target_users"`
	FoundingYear  string              `json:"founding_year"`
	TeamSize      string              `json:"team_size"`
	Funding       string              `json:"funding"`
	CoreFeatures  []FeatureInfo       `json:"core_features"`
	Pricing       PricingInfo         `json:"pricing"`
	Confidence    string              `json:"confidence"`
}

// FeatureInfo 功能信息
type FeatureInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Unique      bool   `json:"unique"`
}

// PricingInfo 价格信息
type PricingInfo struct {
	Model string          `json:"model"`
	Tiers []PricingTier   `json:"tiers"`
	Trial TrialInfo       `json:"trial"`
}

// PricingTier 价格层级
type PricingTier struct {
	Name         string   `json:"name"`
	Price        float64  `json:"price"`
	BillingCycle string   `json:"billing_cycle"`
	Features     []string `json:"features"`
	Limitations  []string `json:"limitations"`
}

// TrialInfo 试用信息
type TrialInfo struct {
	Available bool   `json:"available"`
	Duration  string `json:"duration"`
}

// NewProductInfoExtractor 创建产品信息提取器
func NewProductInfoExtractor(llmClient *LLMClient) *ProductInfoExtractor {
	return &ProductInfoExtractor{
		llmClient: llmClient,
	}
}

// Extract 提取产品信息
func (e *ProductInfoExtractor) Extract(content string) (*ProductInfo, error) {
	systemPrompt := `你是一位专业的产品分析师。请从以下内容中提取竞品的产品信息。

请按照以下JSON格式输出：
{
    "product_name": "产品名称",
    "company": "公司名称",
    "tagline": "产品定位/slogan",
    "target_users": ["目标用户群1", "目标用户群2"],
    "founding_year": "成立年份",
    "team_size": "团队规模",
    "funding": "融资阶段",
    "core_features": [
        {
            "name": "功能名称",
            "description": "功能描述",
            "category": "基础功能/核心功能/高级功能",
            "unique": true
        }
    ],
    "pricing": {
        "model": "订阅制/买断制/免费+增值",
        "tiers": [
            {
                "name": "套餐名称",
                "price": 0,
                "billing_cycle": "月付/年付",
                "features": ["功能1"],
                "limitations": ["限制1"]
            }
        ],
        "trial": {
            "available": true,
            "duration": "14天"
        }
    },
    "confidence": "高/中/低"
}

注意：
1. 如果信息缺失，字段值设为null或空数组
2. 对于不确定的信息，在confidence字段标注"低"
3. 提取时保持客观，避免主观评价`

	userPrompt := fmt.Sprintf(`内容：\n%s\n\n请提取产品信息。`, content)

	response, err := e.llmClient.CompletionWithJSON(systemPrompt, userPrompt)
	if err != nil {
		return nil, err
	}

	// 转换为ProductInfo结构
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	var productInfo ProductInfo
	if err := json.Unmarshal(jsonBytes, &productInfo); err != nil {
		return nil, err
	}

	return &productInfo, nil
}

// SWOTAnalyzer SWOT分析器
type SWOTAnalyzer struct {
	llmClient *LLMClient
}

// SWOTAnalysis SWOT分析结果
type SWOTAnalysis struct {
	Competitor         string             `json:"competitor"`
	Strengths          []SWOTItem         `json:"strengths"`
	Weaknesses         []SWOTItem         `json:"weaknesses"`
	Opportunities      []OpportunityItem  `json:"opportunities"`
	Threats            []ThreatItem       `json:"threats"`
	OverallAssessment  string             `json:"overall_assessment"`
	StrategicSuggestions []string         `json:"strategic_suggestions"`
}

// SWOTItem SWOT项目
type SWOTItem struct {
	Point    string `json:"point"`
	Evidence string `json:"evidence"`
	Impact   string `json:"impact"` // 高/中/低
}

// OpportunityItem 机会项目
type OpportunityItem struct {
	Point   string `json:"point"`
	Context string `json:"context"`
	Action  string `json:"action"`
}

// ThreatItem 威胁项目
type ThreatItem struct {
	Point   string `json:"point"`
	Context string `json:"context"`
	Action  string `json:"action"`
}

// NewSWOTAnalyzer 创建SWOT分析器
func NewSWOTAnalyzer(llmClient *LLMClient) *SWOTAnalyzer {
	return &SWOTAnalyzer{
		llmClient: llmClient,
	}
}

// Analyze 进行SWOT分析
func (a *SWOTAnalyzer) Analyze(competitorName string, productInfo *ProductInfo, marketContext string) (*SWOTAnalysis, error) {
	systemPrompt := `你是一位专业的战略分析师。请对给定的竞品进行SWOT分析。

输出JSON格式：
{
    "competitor": "竞品名称",
    "strengths": [
        {"point": "优势点", "evidence": "证据", "impact": "高"}
    ],
    "weaknesses": [
        {"point": "劣势点", "evidence": "证据", "impact": "高"}
    ],
    "opportunities": [
        {"point": "机会点", "context": "背景", "action": "建议行动"}
    ],
    "threats": [
        {"point": "威胁点", "context": "背景", "action": "应对建议"}
    ],
    "overall_assessment": "总体评估",
    "strategic_suggestions": ["建议1", "建议2"]
}`

	productInfoJSON, _ := json.MarshalIndent(productInfo, "", "  ")

	userPrompt := fmt.Sprintf(`竞品名称：%s

产品信息：
%s

市场背景：
%s

请进行SWOT分析。`, competitorName, string(productInfoJSON), marketContext)

	response, err := a.llmClient.CompletionWithJSON(systemPrompt, userPrompt)
	if err != nil {
		return nil, err
	}

	// 转换为SWOTAnalysis结构
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	var swotAnalysis SWOTAnalysis
	if err := json.Unmarshal(jsonBytes, &swotAnalysis); err != nil {
		return nil, err
	}

	return &swotAnalysis, nil
}
