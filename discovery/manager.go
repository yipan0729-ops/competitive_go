package discovery

import (
	"fmt"
	"sync"
)

// QueryGenerator 查询生成器
type QueryGenerator struct{}

// GenerateCompetitorQueries 生成竞品发现查询
func (q *QueryGenerator) GenerateCompetitorQueries(topic string) []string {
	return []string{
		fmt.Sprintf("%s 竞品", topic),
		fmt.Sprintf("%s 对比", topic),
		fmt.Sprintf("%s 替代品", topic),
		fmt.Sprintf("best %s alternatives", topic),
		fmt.Sprintf("%s vs", topic),
		fmt.Sprintf("%s 排行榜", topic),
		fmt.Sprintf("%s 推荐", topic),
		fmt.Sprintf("top %s tools", topic),
	}
}

// GenerateDataSourceQueries 生成数据源查询
func (q *QueryGenerator) GenerateDataSourceQueries(competitorName string) map[string][]string {
	return map[string][]string{
		"官网": {
			fmt.Sprintf("%s 官网", competitorName),
			fmt.Sprintf("%s official website", competitorName),
		},
		"产品功能": {
			fmt.Sprintf("%s features", competitorName),
			fmt.Sprintf("%s 功能介绍", competitorName),
			fmt.Sprintf("%s 产品", competitorName),
		},
		"定价": {
			fmt.Sprintf("%s pricing", competitorName),
			fmt.Sprintf("%s 价格", competitorName),
			fmt.Sprintf("%s 套餐", competitorName),
		},
		"用户评价": {
			fmt.Sprintf("%s 评价 site:xiaohongshu.com", competitorName),
			fmt.Sprintf("%s 怎么样 site:zhihu.com", competitorName),
			fmt.Sprintf("%s reviews", competitorName),
		},
		"电商": {
			fmt.Sprintf("%s site:taobao.com", competitorName),
			fmt.Sprintf("%s site:jd.com", competitorName),
		},
		"社交媒体": {
			fmt.Sprintf("%s 公众号", competitorName),
			fmt.Sprintf("%s 微博", competitorName),
		},
	}
}

// SearchManager 搜索管理器
type SearchManager struct {
	engines       []SearchEngine
	queryGen      *QueryGenerator
	maxWorkers    int
}

// NewSearchManager 创建搜索管理器
func NewSearchManager(engines []SearchEngine) *SearchManager {
	return &SearchManager{
		engines:    engines,
		queryGen:   &QueryGenerator{},
		maxWorkers: 5,
	}
}

// SearchCompetitors 搜索竞品
func (m *SearchManager) SearchCompetitors(topic string, maxResults int) ([]SearchResult, error) {
	queries := m.queryGen.GenerateCompetitorQueries(topic)

	// 并发搜索
	resultChan := make(chan []SearchResult, len(queries))
	errorChan := make(chan error, len(queries))
	semaphore := make(chan struct{}, m.maxWorkers)

	var wg sync.WaitGroup

	for _, query := range queries {
		wg.Add(1)
		go func(q string) {
			defer wg.Done()
			semaphore <- struct{}{}        // 获取信号量
			defer func() { <-semaphore }() // 释放信号量

			// 使用第一个可用的搜索引擎
			for _, engine := range m.engines {
				results, err := engine.Search(q, maxResults)
				if err == nil {
					resultChan <- results
					return
				}
			}
			errorChan <- fmt.Errorf("所有搜索引擎都失败")
		}(query)
	}

	// 等待所有goroutine完成
	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	// 收集结果
	allResults := []SearchResult{}
	for results := range resultChan {
		allResults = append(allResults, results...)
	}

	// 去重
	uniqueResults := m.deduplicateResults(allResults)

	if len(uniqueResults) == 0 {
		return nil, fmt.Errorf("未找到任何搜索结果")
	}

	return uniqueResults, nil
}

// SearchDataSources 搜索数据源
func (m *SearchManager) SearchDataSources(competitorName string, sourceTypes []string, maxPerType int) (map[string][]SearchResult, error) {
	allQueries := m.queryGen.GenerateDataSourceQueries(competitorName)

	// 过滤需要的数据源类型
	queries := make(map[string][]string)
	if len(sourceTypes) == 0 {
		queries = allQueries
	} else {
		for _, sourceType := range sourceTypes {
			if q, ok := allQueries[sourceType]; ok {
				queries[sourceType] = q
			}
		}
	}

	results := make(map[string][]SearchResult)
	var mu sync.Mutex
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, m.maxWorkers)

	for sourceType, queryList := range queries {
		for _, query := range queryList {
			wg.Add(1)
			go func(st, q string) {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				// 使用第一个可用的搜索引擎
				for _, engine := range m.engines {
					searchResults, err := engine.Search(q, maxPerType)
					if err == nil && len(searchResults) > 0 {
						mu.Lock()
						results[st] = append(results[st], searchResults...)
						mu.Unlock()
						break
					}
				}
			}(sourceType, query)
		}
	}

	wg.Wait()

	// 去重每个类型的结果
	for sourceType, res := range results {
		results[sourceType] = m.deduplicateResults(res)
	}

	return results, nil
}

// deduplicateResults 去重搜索结果
func (m *SearchManager) deduplicateResults(results []SearchResult) []SearchResult {
	seen := make(map[string]bool)
	unique := []SearchResult{}

	for _, result := range results {
		if !seen[result.URL] {
			seen[result.URL] = true
			unique = append(unique, result)
		}
	}

	return unique
}
