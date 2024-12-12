/* ==================================================================
* Package controller 提供了处理文件查询的功能。
* 该处理器用于接收并响应 GET 请求，根据查询参数中的路径（path）
* 查找并返回指定目录或文件的元数据信息。( json 格式！)
*
* 路由：`GET /query`
* 参数：
*   - path: 文件路径，可以是相对路径或绝对路径，
*           如果是相对路径，它会相对于当前工作目录进行查找。
*
* 该处理器会查询文件系统中指定路径下的文件，并返回该文件的元数据，包括：
*   - 文件名（name）
*   - 文件大小（size）
*   - 最后修改时间（mod_time）
*   - 文件扩展名（extension）
*
* 错误处理：
*   1. 如果缺少 `path` 参数，返回 400 错误。
*   2. 如果文件或目录无法找到，返回 404 错误，带有详细错误消息。
*
* 示例：
* 请求路径：`GET /query?path=./data/documents/`
*
* Author: Guo Sihua
* Date: 2024-12-04
 */

package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// FileMetadata 结构体用于返回文件的元数据
type FileMetadata struct {
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	ModTime   time.Time `json:"mod_time"`
	Extension string    `json:"extension"`
}

// HandleFileQuery 处理文件查询请求，返回指定路径下的文件元数据
func HandleFileQuery(w http.ResponseWriter, r *http.Request) {
	// 只允许 GET 请求
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is supported", http.StatusMethodNotAllowed)
		return
	}

	// 获取 URL 中的 path 参数
	pathParam := r.URL.Query().Get("path")
	if pathParam == "" {
		http.Error(w, "Path parameter is required", http.StatusBadRequest)
		return
	}

	// 处理路径，可以是相对路径、绝对路径或特殊根目录路径
	absolutePath, err := resolvePath(pathParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 获取该路径下的文件列表
	files, err := getFilesFromPath(absolutePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error querying files from path %s: %v", absolutePath, err)
		return
	}

	// 设置响应头为 JSON 格式
	w.Header().Set("Content-Type", "application/json")

	// 返回查询到的文件信息
	if err := json.NewEncoder(w).Encode(files); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// resolvePath 处理并解析路径，支持相对路径、绝对路径以及特殊路径
func resolvePath(path string) (string, error) {
	// 清理路径，防止目录遍历攻击
	cleanPath := filepath.Clean(path)

	// 如果路径是绝对路径，直接返回该路径
	if filepath.IsAbs(cleanPath) {
		return cleanPath, nil
	}

	// 如果是相对路径，拼接当前工作目录
	absolutePath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("error resolving relative path: %v", err)
	}

	// 返回解析后的绝对路径
	return absolutePath, nil
}

// getFilesFromPath 根据路径检索文件并返回文件的元数据
func getFilesFromPath(path string) ([]FileMetadata, error) {
	var fileMetadataList []FileMetadata

	// 检查路径是否存在
	dir, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open directory %s: %v", path, err)
	}
	defer dir.Close()

	// 获取该路径下所有文件和目录
	files, err := dir.Readdir(0)
	if err != nil {
		return nil, fmt.Errorf("unable to read directory %s: %v", path, err)
	}

	// 遍历目录中的文件
	for _, file := range files {
		// 如果是文件，则返回它的元数据
		if !file.IsDir() {
			fileMetadata := FileMetadata{
				Name:      file.Name(),
				Size:      file.Size(),
				ModTime:   file.ModTime(),
				Extension: filepath.Ext(file.Name()),
			}
			fileMetadataList = append(fileMetadataList, fileMetadata)
		}
	}

	// 如果没有找到文件
	if len(fileMetadataList) == 0 {
		return nil, fmt.Errorf("no files found in path %s", path)
	}

	// 返回文件列表
	return fileMetadataList, nil
}
