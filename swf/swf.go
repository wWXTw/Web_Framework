package swf

import (
	"net/http"
)

// 函数定义
type HandleFunc func(*Context)

// 实现ServeHttp接口的结构
type Engine struct {
	router *Router
}

// 初始化创建
func New() *Engine {
	return &Engine{
		router: NewRouter(),
	}
}

// 往router中添加路由
func (engine *Engine) AddRoute(method string, pattern string, handler HandleFunc) {
	engine.router.AddRoute(method, pattern, handler)
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
	c := NewContext(w, req)
	engine.router.GetHandler(c)
}
