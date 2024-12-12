package main

import (
	"log"
	"migration/controller"
	"net/http"
)

func main() {
	// 注册路由
	http.HandleFunc("/post", controller.HandlePost)
	http.HandleFunc("/get", controller.HandleGet)
	http.HandleFunc("/post/forward", controller.HandlePostAndForward)
	http.HandleFunc("/query", controller.HandleFileQuery)
	http.HandleFunc("/delete", controller.HandleDelete)

	// 启动 HTTP 服务器
	port := ":8081"
	log.Printf("Server is running on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
