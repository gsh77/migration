/* ==================================================================
* Package controller 提供了处理文件删除的功能。
* 该处理器用于接收并响应 DELETE 请求，根据查询参数中的 `filename`
* 查找并删除指定的文件。
*
* 路由：`DELETE /delete`
* 参数：
*   - filename: 要删除的文件名（必须提供），
*               文件名会在拼接DataDir路径基础上进行查找。
*
* 该处理器会尝试删除指定路径下的文件，如果文件删除成功，返回 200 状态，
* 否则返回 404 错误。
*
* 错误处理：
*   1. 如果缺少 `filename` 参数，返回 400 错误。
*   2. 如果指定的文件不存在，返回 404 错误。
*   3. 如果删除操作失败，返回 500 错误。
*
* 示例:
*   DELETE /delete?filename=documents/file.txt
*
* Author: Guo Sihua
* Date: 2024-12-04
 */

package controller

import (
	"log"
	"migration/config"
	"net/http"
	"os"
	"path/filepath"
)

// HandleDelete 处理 DELETE 请求，用于删除指定文件
func HandleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE method is supported", http.StatusMethodNotAllowed)
		return
	}

	// 获取文件名
	fileName := r.URL.Query().Get("filename")
	if fileName == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	// 清理文件路径，防止路径遍历攻击
	fileName = filepath.Clean(fileName)

	// 形成文件的完整路径
	filePath := filepath.Join(config.DataDir, fileName)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		log.Printf("Error: file %s not found", fileName)
		return
	}

	// 删除文件
	err := os.Remove(filePath)
	if err != nil {
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		log.Printf("Error deleting file %s: %v", fileName, err)
		return
	}

	// 返回成功响应
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File deleted successfully"))
	log.Printf("File %s deleted successfully", fileName)
}
