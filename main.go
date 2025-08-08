package main

import (
	"fmt"
	"net/http"
	"webframe/swf"
)

func main() {
	r := swf.New()
	// 设置routers
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})
	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "[Head:%q] = %q\n", k, v)
		}
	})
	// 启动
	r.Run(":5555")
}
