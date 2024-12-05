package controller

import (
	"log"
	"net/http"
)

// HandlePostAndForward 处理 POST 请求并实时转发数据
func HandlePostAndForward(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is supported", http.StatusMethodNotAllowed)
		return
	}
	// 获取转发文件名称
	fileName := r.URL.Query().Get("fileName")
	if fileName == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}
	// 获取目标地址（转发目标）
	targetURL := r.URL.Query().Get("target")
	if targetURL == "" {
		http.Error(w, "Target URL is required for forwarding", http.StatusBadRequest)
		return
	}
	targetURL = targetURL + "/post?filename=" + fileName
	defer r.Body.Close()

	// 将post的请求体加载为新的请求体
	newRequest, err := http.NewRequest(http.MethodPost, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create new request", http.StatusInternalServerError)
		log.Printf("Error creating new request: %v", err)
		return
	}
	resp, err := http.Post(targetURL, "application/octet-stream", newRequest.Body)
	if err != nil {
		http.Error(w, "Failed to forward data to target", http.StatusInternalServerError)
		log.Printf("Failed to forward data to %s: %v", targetURL, err)
		return
	}
	defer resp.Body.Close()
	log.Printf("Data successfully forwarded to %s with status %d", targetURL, resp.StatusCode)

	// 关闭写端并结束
	w.WriteHeader(http.StatusOK)
	log.Printf("Data transfer completed")
}
