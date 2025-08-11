package swf

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 简化构造json的数据结构
type H map[string]interface{}

// 上下文的数据结构
type Context struct {
	// request与responsewriter
	W   http.ResponseWriter
	Req *http.Request
	// 请求信息
	Path   string
	Method string
	Params map[string]string
	// 返回信息
	StatusCode int
	// 中间件
	handlers []HandleFunc
	index    int
}

// 创建上下文
func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		W:      w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// 获取参数函数
func (c *Context) Param(key string) string {
	v := c.Params[key]
	return v
}

// URL查询函数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 表单查询函数
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 设置返回值函数
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.W.WriteHeader(code)
}

// 设置Header函数
func (c *Context) SetHeader(key string, value string) {
	c.W.Header().Set(key, value)
}

// 构造HTML响应
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.W.Write([]byte(html))
}

// 构造JSON响应
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	// 对json进行编码
	encoder := json.NewEncoder(c.W)
	if err := encoder.Encode(obj); err != nil {
		// 编码失败返回错误信息
		http.Error(c.W, err.Error(), 500)
	}
}

// 构造String响应
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.W.Write([]byte(fmt.Sprintf(format, values...)))
}

// 构造Data响应
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.W.Write(data)
}

// 中间件处理函数
func (c *Context) Next() {
	c.index++
	length := len(c.handlers)
	for ; c.index < length; c.index++ {
		c.handlers[c.index](c)
	}
}
