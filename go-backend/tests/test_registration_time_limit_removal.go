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

// 测试数据结构
type CompetitionCreateRequest struct {
	Title             string     `json:"title"`
	Type              string     `json:"type"`
	Organizer         string     `json:"organizer"`
	RegistrationStart *time.Time `json:"registration_start"`
	RegistrationEnd   *time.Time `json:"registration_end"`
	StartTime         *time.Time `json:"start_time"`
	EndTime           *time.Time `json:"end_time"`
	Description       string     `json:"description"`
	IsOpen            bool       `json:"is_open"`
	MaxParticipants   int        `json:"max_participants"`
}

type CompetitionResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ID                uint       `json:"id"`
		Title             string     `json:"title"`
		RegistrationStart *time.Time `json:"registration_start"`
		RegistrationEnd   *time.Time `json:"registration_end"`
		StartTime         *time.Time `json:"start_time"`
		EndTime           *time.Time `json:"end_time"`
	} `json:"data"`
}

func main() {
	fmt.Println("=== 测试报名时间限制解除 ===")
	fmt.Println()

	// 获取管理员token
	token, err := getAdminToken()
	if err != nil {
		fmt.Printf("获取管理员token失败: %v\n", err)
		return
	}

	fmt.Printf("获取到管理员token: %s\n", token)
	fmt.Println()

	// 测试1: 设置过去的报名开始时间
	fmt.Println("测试1: 设置过去的报名开始时间")
	testPastRegistrationStart(token)
	fmt.Println()

	// 测试2: 设置过去的报名截止时间
	fmt.Println("测试2: 设置过去的报名截止时间")
	testPastRegistrationEnd(token)
	fmt.Println()

	// 测试3: 设置过去的报名时间段
	fmt.Println("测试3: 设置过去的报名时间段")
	testPastRegistrationPeriod(token)
	fmt.Println()

	// 测试4: 设置立即开始的报名
	fmt.Println("测试4: 设置立即开始的报名")
	testImmediateRegistration(token)
	fmt.Println()

	// 测试5: 验证时间逻辑关系
	fmt.Println("测试5: 验证时间逻辑关系")
	testTimeLogicValidation(token)
	fmt.Println()

	fmt.Println("=== 测试完成 ===")
}

// 获取管理员token
func getAdminToken() (string, error) {
	loginData := map[string]string{
		"username": "admin",
		"password": "admin123",
	}

	jsonData, _ := json.Marshal(loginData)
	resp, err := http.Post(baseURL+"/api/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if result["code"] == float64(200) {
		return result["data"].(map[string]interface{})["token"].(string), nil
	}

	return "", fmt.Errorf("登录失败: %s", result["message"])
}

// 测试1: 设置过去的报名开始时间
func testPastRegistrationStart(token string) {
	pastTime := time.Now().Add(-24 * time.Hour)      // 24小时前
	futureTime := time.Now().Add(7 * 24 * time.Hour) // 7天后

	req := CompetitionCreateRequest{
		Title:             "测试竞赛-过去报名开始时间",
		Type:              "校级",
		Organizer:         "测试主办方",
		RegistrationStart: &pastTime,
		RegistrationEnd:   &futureTime,
		StartTime:         &futureTime,
		EndTime:           &futureTime,
		Description:       "测试设置过去的报名开始时间",
		IsOpen:            true,
		MaxParticipants:   100,
	}

	success := createCompetition(token, req)
	if success {
		fmt.Println("✅ 成功设置过去的报名开始时间")
	} else {
		fmt.Println("❌ 设置过去的报名开始时间失败")
	}
}

// 测试2: 设置过去的报名截止时间
func testPastRegistrationEnd(token string) {
	pastTime := time.Now().Add(-12 * time.Hour)      // 12小时前
	futureTime := time.Now().Add(7 * 24 * time.Hour) // 7天后

	req := CompetitionCreateRequest{
		Title:             "测试竞赛-过去报名截止时间",
		Type:              "校级",
		Organizer:         "测试主办方",
		RegistrationStart: &pastTime,
		RegistrationEnd:   &pastTime,
		StartTime:         &futureTime,
		EndTime:           &futureTime,
		Description:       "测试设置过去的报名截止时间",
		IsOpen:            true,
		MaxParticipants:   100,
	}

	success := createCompetition(token, req)
	if success {
		fmt.Println("✅ 成功设置过去的报名截止时间")
	} else {
		fmt.Println("❌ 设置过去的报名截止时间失败")
	}
}

// 测试3: 设置过去的报名时间段
func testPastRegistrationPeriod(token string) {
	pastStart := time.Now().Add(-48 * time.Hour)     // 48小时前
	pastEnd := time.Now().Add(-24 * time.Hour)       // 24小时前
	futureTime := time.Now().Add(7 * 24 * time.Hour) // 7天后

	req := CompetitionCreateRequest{
		Title:             "测试竞赛-过去报名时间段",
		Type:              "校级",
		Organizer:         "测试主办方",
		RegistrationStart: &pastStart,
		RegistrationEnd:   &pastEnd,
		StartTime:         &futureTime,
		EndTime:           &futureTime,
		Description:       "测试设置过去的报名时间段",
		IsOpen:            true,
		MaxParticipants:   100,
	}

	success := createCompetition(token, req)
	if success {
		fmt.Println("✅ 成功设置过去的报名时间段")
	} else {
		fmt.Println("❌ 设置过去的报名时间段失败")
	}
}

// 测试4: 设置立即开始的报名
func testImmediateRegistration(token string) {
	now := time.Now()
	futureTime := time.Now().Add(7 * 24 * time.Hour) // 7天后

	req := CompetitionCreateRequest{
		Title:             "测试竞赛-立即开始报名",
		Type:              "校级",
		Organizer:         "测试主办方",
		RegistrationStart: &now,
		RegistrationEnd:   &futureTime,
		StartTime:         &futureTime,
		EndTime:           &futureTime,
		Description:       "测试设置立即开始的报名",
		IsOpen:            true,
		MaxParticipants:   100,
	}

	success := createCompetition(token, req)
	if success {
		fmt.Println("✅ 成功设置立即开始的报名")
	} else {
		fmt.Println("❌ 设置立即开始的报名失败")
	}
}

// 测试5: 验证时间逻辑关系
func testTimeLogicValidation(token string) {
	now := time.Now()
	futureTime := time.Now().Add(7 * 24 * time.Hour) // 7天后

	// 测试错误的逻辑关系：报名开始时间晚于报名截止时间
	req := CompetitionCreateRequest{
		Title:             "测试竞赛-错误时间逻辑",
		Type:              "校级",
		Organizer:         "测试主办方",
		RegistrationStart: &futureTime, // 报名开始时间晚于报名截止时间
		RegistrationEnd:   &now,        // 报名截止时间早于报名开始时间
		StartTime:         &futureTime,
		EndTime:           &futureTime,
		Description:       "测试错误的时间逻辑关系",
		IsOpen:            true,
		MaxParticipants:   100,
	}

	success := createCompetition(token, req)
	if !success {
		fmt.Println("✅ 正确阻止了错误的时间逻辑关系")
	} else {
		fmt.Println("❌ 应该阻止错误的时间逻辑关系")
	}
}

// 创建竞赛
func createCompetition(token string, req CompetitionCreateRequest) bool {
	jsonData, _ := json.Marshal(req)

	client := &http.Client{}
	request, _ := http.NewRequest("POST", baseURL+"/api/admin/competitions", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(request)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result CompetitionResponse
	json.Unmarshal(body, &result)

	fmt.Printf("响应状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应内容: %s\n", string(body))

	return result.Code == 200
}

// 辅助函数：格式化时间显示
func formatTime(t *time.Time) string {
	if t == nil {
		return "未设置"
	}
	return t.Format("2006-01-02 15:04:05")
}
