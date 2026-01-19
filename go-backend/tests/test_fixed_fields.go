package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const dsn = "root:123456@tcp(localhost:3306)/cloud_dream_system?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	fmt.Println("=== 测试修复后的字段名 ===")
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

	// 测试查询语句
	fmt.Println("测试查询语句...")

	// 测试1: 查询users表
	fmt.Println("测试查询users表:")
	query1 := `SELECT id, username, create_time FROM users LIMIT 1`
	rows, err := db.Query(query1)
	if err != nil {
		fmt.Printf("❌ 查询users表失败: %v\n", err)
	} else {
		fmt.Println("✅ 查询users表成功")
		rows.Close()
	}

	// 测试2: 查询projects表
	fmt.Println("测试查询projects表:")
	query2 := `SELECT id, title, created_at FROM projects LIMIT 1`
	rows, err = db.Query(query2)
	if err != nil {
		fmt.Printf("❌ 查询projects表失败: %v\n", err)
	} else {
		fmt.Println("✅ 查询projects表成功")
		rows.Close()
	}

	// 测试3: 查询competitions表
	fmt.Println("测试查询competitions表:")
	query3 := `SELECT id, title, created_at FROM competitions LIMIT 1`
	rows, err = db.Query(query3)
	if err != nil {
		fmt.Printf("❌ 查询competitions表失败: %v\n", err)
	} else {
		fmt.Println("✅ 查询competitions表成功")
		rows.Close()
	}

	// 测试4: 查询system_logs表
	fmt.Println("测试查询system_logs表:")
	query4 := `SELECT id, log_type, created_at FROM system_logs LIMIT 1`
	rows, err = db.Query(query4)
	if err != nil {
		fmt.Printf("❌ 查询system_logs表失败: %v\n", err)
	} else {
		fmt.Println("✅ 查询system_logs表成功")
		rows.Close()
	}

	fmt.Println("\n=== 测试完成 ===")
}
