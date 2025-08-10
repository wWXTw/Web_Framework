package swf

import (
	"log"
	"net/http"
	"strings"
)

// Router结构
type Router struct {
	// 每种方式(GET/POST/...)对应一颗前缀树
	roots    map[string]*Node
	handlers map[string]HandleFunc
}

// 创建Router函数
func NewRouter() *Router {
	return &Router{
		roots:    make(map[string]*Node),
		handlers: make(map[string]HandleFunc),
	}
}

// 解析URL.Path函数
func ParsePath(path string) []string {
	res := make([]string, 0)
	parts := strings.Split(path, "/")
	for _, v := range parts {
		if v != "" {
			res = append(res, v)
			// 通配符之后路由无意义
			if v[0] == '*' {
				break
			}
		}
	}
	return res
}

// 加入新route函数
func (r *Router) AddRoute(method string, pattern string, handler HandleFunc) {
	parts := ParsePath(pattern)
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &Node{}
	}
	key := method + "%" + pattern
	// 在前缀树中插入路由结点
	r.roots[method].InsertTrie(pattern, parts, 0)
	r.handlers[key] = handler
	log.Printf("%s - %s has been added\n", method, pattern)
}

// 获取route函数
func (r *Router) GetRoute(method string, path string) (*Node, map[string]string) {
	parts := ParsePath(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	// 没有前缀树直接返回空
	if !ok {
		return nil, nil
	}
	// 找寻对应路由结点
	node := root.QueryTrie(parts, 0)
	if node != nil {
		nodeParts := ParsePath(node.pattern)
		for index, part := range nodeParts {
			if part[0] == ':' {
				// 如果查询出来是动态则进行记录 (如:lang:"a" -> 输入是/d/a但只查询出来/:lang/a则进行记录)
				params[part[1:]] = parts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				// 同理如果是通配符则全部进行记录
				params[part[1:]] = strings.Join(parts[index:], "/")
				break
			}
		}
		return node, params
	}
	return nil, nil
}

// 处理请求函数
func (r *Router) GetHandler(c *Context) {
	node, params := r.GetRoute(c.Method, c.Path)
	if node != nil {
		// 结点存在记录参数并响应请求
		c.Params = params
		// 利用node.pattern还原注册路由
		key := c.Method + "%" + node.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 Not Found:%q\n", c.Path)
	}
}
