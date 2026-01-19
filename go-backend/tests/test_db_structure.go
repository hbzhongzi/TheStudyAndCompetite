package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const dsn = "root:123456@tcp(localhost:3306)/cloud_dream_system?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	fmt.Println("=== 检查数据库表结构 ===")
	fmt.Println()

	// 连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	fmt.Println("✅ 数据库连接成功")
	fmt.Println()

	// 检查所有表的字段结构
	fmt.Println("检查所有表的字段结构...")

	// 检查users表
	fmt.Println("\n--- users表 ---")
	checkTableFields(db, "users")

	// 检查projects表
	fmt.Println("\n--- projects表 ---")
	checkTableFields(db, "projects")

	// 检查competitions表
	fmt.Println("\n--- competitions表 ---")
	checkTableFields(db, "competitions")

	// 检查system_logs表
	fmt.Println("\n--- system_logs表 ---")
	checkTableFields(db, "system_logs")

	// 检查其他相关表
	fmt.Println("\n--- user_profiles表 ---")
	checkTableFields(db, "user_profiles")

	fmt.Println("\n--- roles表 ---")
	checkTableFields(db, "roles")

	fmt.Println("\n--- user_roles表 ---")
	checkTableFields(db, "user_roles")

	fmt.Println("\n=== 检查完成 ===")
}

// 检查单个表的字段
func checkTableFields(db *sql.DB, tableName string) {
	query := fmt.Sprintf("DESCRIBE %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("❌ 查询%s表结构失败: %v\n", tableName, err)
		return
	}
	defer rows.Close()

	fmt.Printf("%s表字段:\n", tableName)
	for rows.Next() {
		var field, typ, null, key, defaultVal, extra sql.NullString
		if err := rows.Scan(&field, &typ, &null, &key, &defaultVal, &extra); err != nil {
			fmt.Printf("❌ 读取字段信息失败: %v\n", err)
			continue
		}
		fmt.Printf("  - %s: %s\n", field.String, typ.String)
	}
}
