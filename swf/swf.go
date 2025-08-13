package swf

import (
	"net/http"
	"path"
	"strings"
)

// 函数定义
type HandleFunc func(*Context)

// 实现ServeHttp接口的engine结构
type Engine struct {
	// 嵌入式字段 (类似于继承的功能 | 这样既可以直接注册也可以分组注册)
	*RouterGroup
	router *Router
	// 保存所有的Group
	groups []*RouterGroup
}

// 路由组的结构
type RouterGroup struct {
	// 组前缀
	prefix string
	// 组映射的中间件
	middlewares []HandleFunc
	// support nesting...?
	parent *RouterGroup
	// 每个Group保存一个Engine对象
	engine *Engine
}

// 初始化创建
func New() *Engine {
	engine := &Engine{
		router: NewRouter(),
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 创建新Group函数
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		// 前缀拼接
		prefix:      group.prefix + prefix,
		middlewares: make([]HandleFunc, 0),
		parent:      group,
		engine:      engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 往router中添加路由 (Engine与RouterGroup都能使用)
func (group *RouterGroup) AddRoute(method string, component string, handler HandleFunc) {
	pattern := group.prefix + component
	group.engine.router.AddRoute(method, pattern, handler)
}

// Get函数 (Engine与RouterGroup都能使用)
func (group *RouterGroup) GET(component string, handler HandleFunc) {
	group.AddRoute("GET", component, handler)
}

// Post函数 (Engine与RouterGroup都能使用)
func (group *RouterGroup) POST(component string, handler HandleFunc) {
	group.AddRoute("POST", component, handler)
}

// 启动框架函数
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

// 设置中间件
func (group *RouterGroup) Use(middle ...HandleFunc) {
	group.middlewares = append(group.middlewares, middle...)
}

// 重写ServeHTTP接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(w, req)
	// 根据请求从Group设置的中间件中获取对应的放入上下文Context中
	middlewares := make([]HandleFunc, 0)
	for _, group := range engine.groups {
		if isMember := strings.HasPrefix(req.URL.Path, group.prefix); isMember {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	// 将中间件放入Context中
	c.handlers = middlewares
	engine.router.GetHandler(c)
}

// 创建静态文件handler
func (group *RouterGroup) CreateStaticHandler(relativePath string, fs http.FileSystem) HandleFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	// fileServer为一个新的的handler(与engine类似)实例,功能是将URL的前缀去掉并在fs对应的目录下进行查找
	// fileServer本质上是一个包装好的handler处理器
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(ctx *Context) {
		file := ctx.Param("filepath")
		// 不能打开/文件不存在返回错误
		if _, err := fs.Open(file); err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}
		// 先去除前缀,再由内部的http.FileServer(fs)进行处理
		fileServer.ServeHTTP(ctx.W, ctx.Req)
	}
}

// 用户的静态文件映射接口
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.CreateStaticHandler(relativePath, http.Dir(root))
	// 设置URL
	pattern := relativePath + "/*filepath"
	// 注册路由
	group.GET(pattern, handler)
}
