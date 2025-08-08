package swf

import (
	"fmt"
	"net/http"
)

// 函数定义
type HandleFunc func(http.ResponseWriter, *http.Request)

// 实现ServeHttp接口的结构
type Engine struct {
	// router为一个根据不同req获取不同函数的哈希表
	router map[string]HandleFunc
}

// 构造函数
func New() *Engine {
	return &Engine{
		router: make(map[string]HandleFunc),
	}
}

// 加入新route函数
func (engine *Engine) AddRoute(method string, pattern string, handler HandleFunc) {
	key := method + "%" + pattern
	engine.router[key] = handler
}

// Get函数
func (engine *Engine) GET(pattern string, handler HandleFunc) {
	engine.AddRoute("GET", pattern, handler)
}

// Post函数
func (engine *Engine) POST(pattern string, handler HandleFunc) {
	engine.AddRoute("POST", pattern, handler)
}

// 启动框架函数
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

// 重写ServeHTTP接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "%" + req.URL.Path
	// 检验这个router是否存储在routers之中
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 Not Found:%q\n", req.URL.Path)
	}
}
