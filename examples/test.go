package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080"

func main() {
	fmt.Println("=== 自动化竞品调研工具测试 ===\n")

	// 1. 健康检查
	fmt.Println("1. 健康检查...")
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		fmt.Printf("❌ 服务未启动: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("✅ 服务正常运行\n")

	// 2. 测试智能发现
	fmt.Println("2. 测试智能数据源发现...")
	taskID := testDiscovery()
	if taskID > 0 {
		fmt.Printf("✅ 发现任务已创建，任务ID: %d\n", taskID)
		
		// 等待任务完成
		fmt.Println("   等待任务完成...")
		time.Sleep(5 * time.Second)
		
		// 查询任务状态
		status := getTaskStatus(taskID)
		fmt.Printf("   任务状态: %s, 进度: %d%%\n", status["status"], int(status["progress"].(float64)))
	}
	fmt.Println()

	// 3. 测试单个URL爬取
	fmt.Println("3. 测试单个URL爬取...")
	testCrawlSingle()
	fmt.Println()

	// 4. 获取竞品列表
	fmt.Println("4. 获取竞品列表...")
	testGetCompetitors()
	fmt.Println()

	fmt.Println("=== 测试完成 ===")
}

func testDiscovery() int {
	reqBody := map[string]interface{}{
		"topic":           "AI写作助手",
		"market":          "中国",
		"competitor_count": 3,
		"source_types":    []string{"官网", "评价"},
		"depth":           "quick",
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post(baseURL+"/api/discover/search", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ 请求失败: %v\n", err)
		return 0
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if taskID, ok := result["task_id"].(float64); ok {
		return int(taskID)
	}
	return 0
}

func getTaskStatus(taskID int) map[string]interface{} {
	url := fmt.Sprintf("%s/api/discover/status/%d", baseURL, taskID)
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result
}

func testCrawlSingle() {
	reqBody := map[string]interface{}{
		"url":         "https://notion.so",
		"competitor":  "Notion AI",
		"source_type": "官网",
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post(baseURL+"/api/crawl/single", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	if resp.StatusCode == 200 {
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		if success, ok := result["success"].(bool); ok && success {
			fmt.Printf("✅ 爬取成功\n")
			fmt.Printf("   内容路径: %s\n", result["content_path"])
			fmt.Printf("   图片数量: %.0f\n", result["image_count"].(float64))
			fmt.Printf("   标题: %s\n", result["title"])
		}
	} else {
		fmt.Printf("❌ 爬取失败: %s\n", string(body))
	}
}

func testGetCompetitors() {
	resp, err := http.Get(baseURL + "/api/competitors?page=1&page_size=10")
	if err != nil {
		fmt.Printf("❌ 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if total, ok := result["total"].(float64); ok {
		fmt.Printf("✅ 共有 %.0f 个竞品\n", total)
		
		if competitors, ok := result["competitors"].([]interface{}); ok && len(competitors) > 0 {
			fmt.Println("   竞品列表:")
			for i, comp := range competitors {
				if compMap, ok := comp.(map[string]interface{}); ok {
					fmt.Printf("   %d. %s (%s)\n", i+1, compMap["name"], compMap["status"])
				}
			}
		}
	}
}
