package main

import (
	"fmt"
	"log"
	"yunmeng-backend/utils"
)

func main() {
	fmt.Println("测试JWT Token生成和解析...")

	// 测试管理员Token
	adminToken, err := utils.GenerateToken(1, "admin", "admin")
	if err != nil {
		log.Fatalf("生成管理员Token失败: %v", err)
	}
	fmt.Printf("管理员Token: %s\n", adminToken)

	// 测试教师Token
	teacherToken, err := utils.GenerateToken(2, "teacher", "teacher")
	if err != nil {
		log.Fatalf("生成教师Token失败: %v", err)
	}
	fmt.Printf("教师Token: %s\n", teacherToken)

	// 测试学生Token
	studentToken, err := utils.GenerateToken(3, "student", "student")
	if err != nil {
		log.Fatalf("生成学生Token失败: %v", err)
	}
	fmt.Printf("学生Token: %s\n", studentToken)

	// 测试Token解析
	fmt.Println("\n测试Token解析...")

	// 解析管理员Token
	adminClaims, err := utils.ParseToken(adminToken)
	if err != nil {
		log.Fatalf("解析管理员Token失败: %v", err)
	}
	fmt.Printf("管理员Token解析结果 - 用户ID: %d, 用户名: %s, 角色: %s\n",
		adminClaims.UserID, adminClaims.Username, adminClaims.Role)

	// 解析教师Token
	teacherClaims, err := utils.ParseToken(teacherToken)
	if err != nil {
		log.Fatalf("解析教师Token失败: %v", err)
	}
	fmt.Printf("教师Token解析结果 - 用户ID: %d, 用户名: %s, 角色: %s\n",
		teacherClaims.UserID, teacherClaims.Username, teacherClaims.Role)

	// 解析学生Token
	studentClaims, err := utils.ParseToken(studentToken)
	if err != nil {
		log.Fatalf("解析学生Token失败: %v", err)
	}
	fmt.Printf("学生Token解析结果 - 用户ID: %d, 用户名: %s, 角色: %s\n",
		studentClaims.UserID, studentClaims.Username, studentClaims.Role)

	fmt.Println("\n权限测试完成！")
}
