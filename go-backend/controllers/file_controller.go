package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"yunmeng-backend/models"

	"github.com/gin-gonic/gin"
)

type FileController struct{}

func NewFileController() *FileController {
	return &FileController{}
}

// UploadFile 上传文件
func (c *FileController) UploadFile(ctx *gin.Context) {
	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文件上传失败: " + err.Error(),
		})
		return
	}

	// 检查文件大小（限制为10MB）
	if file.Size > 10*1024*1024 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文件大小不能超过10MB",
		})
		return
	}

	// 检查文件类型
	allowedTypes := map[string]bool{
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".xls":  true,
		".xlsx": true,
		".ppt":  true,
		".pptx": true,
		".txt":  true,
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}

	ext := filepath.Ext(file.Filename)
	if !allowedTypes[ext] {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "不支持的文件类型",
		})
		return
	}

	// 创建上传目录
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("创建上传目录失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
		return
	}

	// 生成唯一文件名
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, file.Filename)
	filepath := filepath.Join(uploadDir, filename)

	// 保存文件
	if err := ctx.SaveUploadedFile(file, filepath); err != nil {
		log.Printf("保存文件失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "文件保存失败",
		})
		return
	}

	// 返回文件信息
	response := models.FileUploadResponse{
		FileName: file.Filename,
		FileURL:  "/uploads/" + filename,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文件上传成功",
		"data":    response,
	})
}
