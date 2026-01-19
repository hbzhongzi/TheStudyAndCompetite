package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Token   string      `json:"token,omitempty"`
}

type CompetitionRequest struct {
	Title           string `json:"title"`
	Type            string `json:"type"`
	Organizer       string `json:"organizer"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	Description     string `json:"description"`
	IsOpen          bool   `json:"is_open"`
	MaxParticipants int    `json:"max_participants"`
}

func main() {
	log.Println("=== 管理员权限测试工具 ===")

	baseURL := "http://localhost:8080/api"

	// 步骤1：测试管理员登录
	log.Println("步骤1: 测试管理员登录...")
	token, err := testAdminLogin(baseURL)
	if err != nil {
		log.Fatal("管理员登录失败:", err)
	}

	log.Println("✅ 管理员登录成功")
	log.Printf("Token: %s...", token[:20])

	// 步骤2：测试创建竞赛权限
	log.Println("\n步骤2: 测试创建竞赛权限...")
	err = testCreateCompetition(baseURL, token)
	if err != nil {
		log.Fatal("创建竞赛测试失败:", err)
	}

	log.Println("✅ 创建竞赛权限测试成功")
	log.Println("\n🎉 所有测试通过！管理员权限正常")
}

func testAdminLogin(baseURL string) (string, error) {
	loginData := LoginRequest{
		Username: "admin",
		Password: "123456",
		Role:     "admin",
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return "", fmt.Errorf("序列化登录数据失败: %v", err)
	}

	resp, err := http.Post(baseURL+"/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("登录请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return "", fmt.Errorf("解析登录响应失败: %v", err)
	}

	if loginResp.Code != 200 {
		return "", fmt.Errorf("登录失败: %s", loginResp.Message)
	}

	if loginResp.Token == "" {
		return "", fmt.Errorf("登录响应中没有Token")
	}

	return loginResp.Token, nil
}

func testCreateCompetition(baseURL, token string) error {
	competitionData := CompetitionRequest{
		Title:           "测试竞赛",
		Type:            "学术竞赛",
		Organizer:       "测试组织",
		StartTime:       time.Now().Format("2006-01-02 15:04:05"),
		EndTime:         time.Now().AddDate(0, 1, 0).Format("2006-01-02 15:04:05"),
		Description:     "这是一个测试竞赛",
		IsOpen:          true,
		MaxParticipants: 100,
	}

	jsonData, err := json.Marshal(competitionData)
	if err != nil {
		return fmt.Errorf("序列化竞赛数据失败: %v", err)
	}

	req, err := http.NewRequest("POST", baseURL+"/admin/competitions", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("创建竞赛请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode == 403 {
		return fmt.Errorf("权限不足 (403): %s", string(body))
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("创建竞赛失败 (状态码: %d): %s", resp.StatusCode, string(body))
	}

	log.Println("✅ 竞赛创建API调用成功")
	return nil
}
