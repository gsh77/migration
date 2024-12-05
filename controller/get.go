package controller

import (
	"io"
	"log"
	"migration/config"
	"migration/utils"
	"net/http"
	"path/filepath"
)

// HandleGet 处理 GET 请求，用于返回存储的文件
func HandleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is supported", http.StatusMethodNotAllowed)
		return
	}

	// 获取文件名
	fileName := r.URL.Query().Get("filename")
	if fileName == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}
	fileName = filepath.Clean(fileName)

	// 调用文件加载函数
	file, err := utils.LoadFile(config.DataDir, fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		log.Printf("Error loading file %s: %v", fileName, err)
		return
	}
	defer file.Close()

	// 设置响应头
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)

	// 将文件内容写入响应
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Failed to send file", http.StatusInternalServerError)
		log.Printf("Error sending file: %v", err)
	}
}
