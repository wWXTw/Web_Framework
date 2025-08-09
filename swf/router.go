package swf

import (
	"log"
	"net/http"
)

// Router结构
type Router struct {
	handlers map[string]HandleFunc
}

// 创建Router函数
func NewRouter() *Router {
	return &Router{
		handlers: make(map[string]HandleFunc),
	}
}

// 加入新route函数
func (r *Router) AddRoute(method string, pattern string, handler HandleFunc) {
	log.Printf("%s - %s has been added\n", method, pattern)
	key := method + "%" + pattern
	r.handlers[key] = handler
}

func (r *Router) GetHandler(c *Context) {
	key := c.Method + "%" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 Not Found:%q\n", c.Path)
	}
}
