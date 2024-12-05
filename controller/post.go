package controller

import (
	"log"
	"migration/config"
	"migration/utils"
	"net/http"
	"path/filepath"
)

// HandlePost 处理大文件的上传
func HandlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is supported", http.StatusMethodNotAllowed)
		return
	}

	// 获取文件名（从 URL 参数或者 Header 中获取）
	fileName := r.URL.Query().Get("filename")
	if fileName == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}
	fileName = filepath.Clean(fileName)

	//调用文件保存函数
	written, err := utils.SaveFile(config.DataDir, fileName, r.Body)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		log.Printf("Error saving file: %v", err)
		return
	}
	//// 存储到内存
	//buf := new(bytes.Buffer)
	//written, err := io.Copy(buf, r.Body)
	//if err != nil {
	//	http.Error(w, "Failed to read file content", http.StatusInternalServerError)
	//	log.Printf("Error reading file content: %v", err)
	//	return
	//}
	//log.Printf("Received %d bytes from client", written)

	// 返回成功响应
	log.Printf("File %s uploaded successfully, size: %d bytes", fileName, written)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("File uploaded successfully"))
}
