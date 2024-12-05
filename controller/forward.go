package controller

import (
	"bytes"
	"io"
	"log"
	"migration/config"
	"migration/utils"
	"net/http"
)

func HandleForward(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is supported", http.StatusMethodNotAllowed)
		return
	}
	// 获取目标地址（转发目标）
	targetURL := r.URL.Query().Get("target")
	if targetURL == "" {
		http.Error(w, "Target URL is required for forwarding", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fileName := r.URL.Query().Get("fileName")
	if fileName == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	file, err := utils.LoadFile(config.DataDir, fileName)
	if err != nil {
		http.Error(w, "Failed to load file", http.StatusInternalServerError)
		log.Printf("Error loading file: %v", err)
		return
	}
	defer file.Close()

	// 将文件内容读取到一个缓冲区，以便转发
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	if err != nil {
		http.Error(w, "Failed to read file content", http.StatusInternalServerError)
		log.Printf("Error reading file content: %v", err)
		return
	}

	// 向目标服务器发送 POST 请求
	resp, err := http.Post(targetURL+"/post?filename="+fileName, "application/octet-stream", buf)
	if err != nil {
		http.Error(w, "Failed to forward data to target", http.StatusInternalServerError)
		log.Printf("Failed to forward data to %s: %v", targetURL, err)
		return
	}
	defer resp.Body.Close()
	log.Printf("Data successfully forwarded to %s with status %d", targetURL, resp.StatusCode)

	// 将目标服务器的响应直接写回客户端
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error copying response body: %v", err)
	}
}
