# SWF - A Tiny Gin-like Web Framework in Go
> 受Gin启发编写的轻量级Web框架

### ✨ 框架特点  
- 为每个路由设置一颗前缀树(Tries),支持动态路由(`:params`)与通配符(`*wildcard`)
- 支持RouterGroup前缀拼接,更易维护中大型路由
- Context接口封装,且支持`String`/`JSON`/`HTML`/`Data`等API
- 支持全局与分组中间件功能,内置`Logger`,`Recovery`等实用中间件
- 支持静态文件目录映射
- 支持`text/html`(含`FuncMap`)渲染HTML

### 📦 目录结构  
WebFramework/  
├── main.go &nbsp; #入口文件  
├── go.mod  
├── static &nbsp; #默认静态文件存放位置  
└── swf/  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├── swf.go &nbsp; #Engine&RouterGroup:核心入口,分组接口,中间件接口,静态文件映射接口,HTML模版绑定接口  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├── router.go &nbsp; #Router:输入pattern与Tries结点之间的中间层,提供handler  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├── tries.go &nbsp; #Tries:前缀树结点的设置与插入、查询函数,使其支持动态路由与通配符  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├── context.go &nbsp; #统一封装Context上下文,支持多种API响应,中间件Next()调度  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├── logger.go &nbsp; #Logger中间件,打印操作状态  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;└── recovery.go &nbsp; #Recovery中间件,panic安全恢复  

### 🧱 主要方法
> 路由功能
- `(group *RouterGroup) GET(component string, handler HandleFunc)` | `swf.go`
- `(group *RouterGroup) POST(component string, handler HandleFunc)` | `swf.go`  
  面向用户的GET与POST注册路由函数,都是通过AddRoute()进行的  
  由于RouterGroup结构体内部嵌入了对Engine的引用,因此无论是在全局Engine层,还是在具体的路由组层,都能够调用统一的路由注册方法
- `(*Router) AddRoute(method, pattern string, handler HandleFunc)` | `router.go`  
  解析pattern在前缀树中进行注册,同时以Method%pattern为key注册handler哈希表
- `(*Router) GetRoute(method, path string) (*Node, map[string]string)` | `router.go`  
  在前缀树上找出结点与参数字典
- `(*Router) GetHandler(c *Context)` | `router.go`  
  将得到的结点的handler放入执行链中,并将参数字典记录到Context中
- `(group *RouterGroup) Group(prefix string) *RouterGroup` | `swf.go`  
  面向用户的创建新RouterGroup函数,通过前缀拼接确定新prefix
> 上下文功能
- 上下文进行统一封装结构
```bash
type Context struct {
	// engine实例
	engine *Engine
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
```
- 请求相关函数  
`(*Context) Param(key string) string`  #从参数字典里查询参数  
`(*Context) Query(key string) string`  #从URL中查询  
`(*Context) PostForm(key string) string`  #从表单数据中查询  
- 响应相关函数  
`(*Context) String(code int, format string, values ...any)`  
`(*Context) JSON(code int, obj any)`  
`(*Context) HTML(code int, name string, data any)`  
`(*Context) Data(code int, data []byte)`  
> 中间件功能
- `(group *RouterGroup) Use(middle ...HandleFunc)` | `swf.go`  
面向用户的中间件添加函数,将handler添加入中间件队列中,context会读取这个队列  
- `(c *Context) Next()` | `context.go`  
context依次执行队列中的函数  
> 静态文件功能
- `(group *RouterGroup) Static(relativePath string, root string)` | `swf.go`  
对外提供的用户接口,自动将URL前缀(如/assets)映射到本地目录(如./static),并把静态处理函数注册到路由表
- `(group *RouterGroup) CreateStaticHandler(relativePath string, fs http.FileSystem) HandleFunc` | `swf.go`  
生成一个静态文件请求处理函数，负责剥离 URL 前缀、检查文件是否存在  
> HTML模板功能
- `(engine *Engine) SetFuncMap(funcMap template.FuncMap)` | `swf.go`  
对外提供的用户接口,用户提供HTML模版的FuncMap
- `(engine *Engine) LoadHTMLGlob(pattern string)` | `swf.go`  
对外提供的用户接口,将HTML模版存入engine结构中,在接收HTML响应式进行调用

### 🚀 快速开始
```bash
git clone https://github.com/wWXTw/Web_Framework.git
cd Web_Framework
go run main.go
# visit: http://localhost:5555  (默认端口设置为5555)
