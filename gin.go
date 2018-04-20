// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"html/template"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin/render"
)

//版本和固定的内存块
const (
	// Version is Framework's version.
	Version                = "v1.2"
	defaultMultipartMemory = 32 << 20 // 32 MB
)

//默认的错误信息
var (
	default404Body = []byte("404 page not found")
	default405Body = []byte("405 method not allowed")
	//默认引擎
	defaultAppEngine bool
)

//定义函数类型
type HandlerFunc func(*Context)

//定义函数链
type HandlersChain []HandlerFunc

//获取最后一个函数
// Last returns the last handler in the chain. ie. the last handler is the main own.
func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}

//路由基本信息
type RouteInfo struct {
	Method  string
	Path    string
	Handler string
}

//路由数组
type RoutesInfo []RouteInfo

// Engine is the framework's instance, it contains
// the muxer, middleware and configuration settings.
// Create an instance of Engine, by using New() or Default()
type Engine struct {
	RouterGroup //路由设置

	// Enables automatic redirection if the current
	// route can't be matched but a
	// handler for the path with (without)
	//the trailing slash exists.
	// For example if /foo/ is requested but a route
	//only exists for /foo, the
	// client is redirected to /foo with
	// http status code 301 for GET requests
	// and 307 for all other request methods.
	RedirectTrailingSlash bool //去除/跳转
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool //允许去除跳转

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool //是否运行处理
	ForwardedByClientIP    bool //客户端的ip

	// #726 #755 If enabled, it will thrust some
	//headers starting with
	// 'X-AppEngine...' for better integration with that PaaS.
	AppEngine bool //？

	// If enabled, the url.RawPath will be used to find parameters.
	UseRawPath bool //寻找参数?

	// If true, the path value will be unescaped.
	// If UseRawPath is false (by default), the UnescapePathValues effectively is true,
	// as url.Path gonna be used, which is already unescaped.
	UnescapePathValues bool //路径处理

	// Value of 'maxMemory' param that
	//is given to http.Request's ParseMultipartForm
	// method call.
	MaxMultipartMemory int64 //http请求的最大内存

	delims           render.Delims     //?
	secureJsonPrefix string            //?json？
	HTMLRender       render.HTMLRender //?模板引擎？
	FuncMap          template.FuncMap  //?模板方法
	allNoRoute       HandlersChain     //所有的无路由？
	allNoMethod      HandlersChain     //所有的无方法？
	noRoute          HandlersChain     //无路由？
	noMethod         HandlersChain     //无方法？
	pool             sync.Pool         //
	trees            methodTrees       //树
}

//一个默认的引擎？
var _ IRouter = &Engine{}

// New returns a new blank Engine instance without
//any middleware attached.
// By default the configuration is:
// - RedirectTrailingSlash:  true
// - RedirectFixedPath:      false
// - HandleMethodNotAllowed: false
// - ForwardedByClientIP:    true
// - UseRawPath:             false
// - UnescapePathValues:     true
func New() *Engine {
	//打印警告？
	debugPrintWARNINGNew()
	engine := &Engine{
		//定义RouteGroup
		RouterGroup: RouterGroup{
			Handlers: nil,  //默认处理器
			basePath: "/",  //默认路径
			root:     true, //主路径
		},
		//模板方法
		FuncMap:                template.FuncMap{},
		RedirectTrailingSlash:  true,                    //去除/
		RedirectFixedPath:      false,                   //?
		HandleMethodNotAllowed: false,                   //?
		ForwardedByClientIP:    true,                    //客户端ip
		AppEngine:              defaultAppEngine,        //默认引擎
		UseRawPath:             false,                   //使用原生的path？
		UnescapePathValues:     true,                    //去除path
		MaxMultipartMemory:     defaultMultipartMemory,  //默认最大内存
		trees:                  make(methodTrees, 0, 9), //函数树
		//模板分隔符
		delims: render.Delims{Left: "{{", Right: "}}"},
		//安全处理？
		secureJsonPrefix: "while(1);",
	}
	//设置group的默认engine
	engine.RouterGroup.engine = engine
	//对pool的new的初始化
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}
	//返回一个engine
	return engine
}

//返回一个默认的Engine，
// Default returns an Engine instance with the Logger and Recovery middleware already attached.
func Default() *Engine {
	debugPrintWARNINGDefault()
	engine := New() //初始化
	//使用Logger和recovery插件
	engine.Use(Logger(), Recovery())
	return engine
}

//新建一个Context，把当前的engin复制给当前的context
func (engine *Engine) allocateContext() *Context {
	return &Context{engine: engine}
}

//设置模板的左右划分字符
func (engine *Engine) Delims(left, right string) *Engine {
	engine.delims = render.Delims{Left: left, Right: right}
	return engine
}

//设置secureJsonPrefix
// SecureJsonPrefix sets the secureJsonPrefix used in Context.SecureJSON.
func (engine *Engine) SecureJsonPrefix(prefix string) *Engine {
	engine.secureJsonPrefix = prefix
	return engine
}

//设置模板函数
// LoadHTMLGlob loads HTML files identified by glob pattern
// and associates the result with HTML renderer.
func (engine *Engine) LoadHTMLGlob(pattern string) {
	left := engine.delims.Left
	right := engine.delims.Right

	//调试的话设置HTMLRender
	if IsDebugging() {
		debugPrintLoadTemplate(template.Must(template.New("").Delims(left, right).Funcs(engine.FuncMap).ParseGlob(pattern)))
		engine.HTMLRender = render.HTMLDebug{Glob: pattern, FuncMap: engine.FuncMap, Delims: engine.delims}
		return
	}

	templ := template.Must(template.New("").Delims(left, right).Funcs(engine.FuncMap).ParseGlob(pattern))
	engine.SetHTMLTemplate(templ)
}

//加载html文件
// LoadHTMLFiles loads a slice of HTML files
// and associates the result with HTML renderer.
func (engine *Engine) LoadHTMLFiles(files ...string) {
	if IsDebugging() {
		engine.HTMLRender = render.HTMLDebug{Files: files, FuncMap: engine.FuncMap, Delims: engine.delims}
		return
	}

	templ := template.Must(template.New("").Delims(engine.delims.Left, engine.delims.Right).Funcs(engine.FuncMap).ParseFiles(files...))
	engine.SetHTMLTemplate(templ)
}

//设置html模板
// SetHTMLTemplate associate a template with HTML renderer.
func (engine *Engine) SetHTMLTemplate(templ *template.Template) {
	if len(engine.trees) > 0 {
		debugPrintWARNINGSetHTMLTemplate()
	}

	engine.HTMLRender = render.HTMLProduction{Template: templ.Funcs(engine.FuncMap)}
}

//设置函数
// SetFuncMap sets the FuncMap used for template.FuncMap.
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.FuncMap = funcMap
}

//设置无route的处理函数
// NoRoute adds handlers for NoRoute. It return a 404 code by default.
func (engine *Engine) NoRoute(handlers ...HandlerFunc) {
	engine.noRoute = handlers
	engine.rebuild404Handlers()
}

//设置无方法的处理函数
// NoMethod sets the handlers called when... TODO.
func (engine *Engine) NoMethod(handlers ...HandlerFunc) {
	engine.noMethod = handlers
	engine.rebuild405Handlers()
}

//注册中间件。同时注册rebuild404Handlers和rebuild405Handlers
// Use attachs a global middleware to the router. ie. the middleware attached though Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func (engine *Engine) Use(middleware ...HandlerFunc) IRoutes {
	engine.RouterGroup.Use(middleware...)
	engine.rebuild404Handlers()
	engine.rebuild405Handlers()
	return engine
}

//把noroute合并到allNoRoute
func (engine *Engine) rebuild404Handlers() {
	engine.allNoRoute = engine.combineHandlers(engine.noRoute)
}

//把noMethod合并到noMethod
func (engine *Engine) rebuild405Handlers() {
	engine.allNoMethod = engine.combineHandlers(engine.noMethod)
}

//添加路由，
//方法，路径，处理器
func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
	assert1(path[0] == '/', "path must begin with '/'")
	assert1(method != "", "HTTP method can not be empty")
	assert1(len(handlers) > 0, "there must be at least one handler")
	//打印调试信息
	debugPrintRoute(method, path, handlers)
	//从树上获取路由
	root := engine.trees.get(method)
	if root == nil {
		//失败，添加一个新的
		root = new(node)
		engine.trees = append(engine.trees, methodTree{method: method, root: root})
	}
	//把处理器加入到里面
	root.addRoute(path, handlers)
}

//得到一个routes
// Routes returns a slice of registered routes, including some useful information, such as:
// the http method, path and the handler name.
func (engine *Engine) Routes() (routes RoutesInfo) {
	for _, tree := range engine.trees {
		routes = iterate("", tree.method, routes, tree.root)
	}
	return routes
}

//循环得到RoutesInfo
func iterate(path, method string, routes RoutesInfo, root *node) RoutesInfo {
	path += root.path
	if len(root.handlers) > 0 {
		routes = append(routes, RouteInfo{
			Method:  method,
			Path:    path,
			Handler: nameOfFunction(root.handlers.Last()),
		})
	}
	for _, child := range root.children {
		routes = iterate(path, method, routes, child)
	}
	return routes
}

//运行程序
// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) Run(addr ...string) (err error) {
	defer func() { debugPrintError(err) }()
	//获取地址
	address := resolveAddress(addr)
	debugPrint("Listening and serving HTTP on %s\n", address)
	//监听和运行代码
	err = http.ListenAndServe(address, engine)
	return
}

//运行tls
// RunTLS attaches the router to a http.Server and starts listening and serving HTTPS (secure) requests.
// It is a shortcut for http.ListenAndServeTLS(addr, certFile, keyFile, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) RunTLS(addr, certFile, keyFile string) (err error) {
	debugPrint("Listening and serving HTTPS on %s\n", addr)
	defer func() { debugPrintError(err) }()

	err = http.ListenAndServeTLS(addr, certFile, keyFile, engine)
	return
}

//运行unix
// RunUnix attaches the router to a http.Server and starts listening and serving HTTP requests
// through the specified unix socket (ie. a file).
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) RunUnix(file string) (err error) {
	debugPrint("Listening and serving HTTP on unix:/%s", file)
	defer func() { debugPrintError(err) }()

	os.Remove(file)
	listener, err := net.Listen("unix", file)
	if err != nil {
		return
	}
	defer listener.Close()
	err = http.Serve(listener, engine)
	return
}

//实现ServeHTTP接口
// ServeHTTP conforms to the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//获取context
	c := engine.pool.Get().(*Context)
	//初始化写
	c.writermem.reset(w)
	//放置请求
	c.Request = req
	c.reset()
	//处理http请求
	engine.handleHTTPRequest(c)
	//放置context
	engine.pool.Put(c)
}

//处理context
// HandleContext re-enter a context that has been rewritten.
// This can be done by setting c.Request.Path to your new target.
// Disclaimer: You can loop yourself to death with this, use wisely.
func (engine *Engine) HandleContext(c *Context) {
	c.reset()
	engine.handleHTTPRequest(c)
	engine.pool.Put(c)
}

//处理http请求
func (engine *Engine) handleHTTPRequest(c *Context) {
	//方法，路径，未转义
	httpMethod := c.Request.Method
	path := c.Request.URL.Path
	unescape := false
	if engine.UseRawPath && len(c.Request.URL.RawPath) > 0 {
		path = c.Request.URL.RawPath
		unescape = engine.UnescapePathValues
	}

	//获取素有的树
	// Find root of the tree for the given HTTP method
	t := engine.trees
	for i, tl := 0, len(t); i < tl; i++ {
		//方法判断
		if t[i].method == httpMethod {
			root := t[i].root
			// Find route in tree
			//路径，参数，转义
			handlers, params, tsr := root.getValue(path, c.Params, unescape)
			if handlers != nil {
				c.handlers = handlers
				c.Params = params
				c.Next()
				//处理请求头
				c.writermem.WriteHeaderNow()
				return
			}
			//连接状态
			if httpMethod != "CONNECT" && path != "/" {
				if tsr && engine.RedirectTrailingSlash {
					redirectTrailingSlash(c)
					return
				}
				if engine.RedirectFixedPath && redirectFixedPath(c, root, engine.RedirectFixedPath) {
					return
				}
			}
			break
		}
	}
	//调用不允许的方法
	if engine.HandleMethodNotAllowed {
		for _, tree := range engine.trees {
			if tree.method != httpMethod {
				if handlers, _, _ := tree.root.getValue(path, nil, unescape); handlers != nil {
					c.handlers = engine.allNoMethod
					serveError(c, 405, default405Body)
					return
				}
			}
		}
	}
	//设置处理器
	c.handlers = engine.allNoRoute
	serveError(c, 404, default404Body)
}

//头部？
var mimePlain = []string{MIMEPlain}

//错误处理
func serveError(c *Context, code int, defaultMessage []byte) {
	c.writermem.status = code
	c.Next()
	//处理错误
	if !c.writermem.Written() {
		if c.writermem.Status() == code {
			c.writermem.Header()["Content-Type"] = mimePlain
			c.Writer.Write(defaultMessage)
		} else {
			c.writermem.WriteHeaderNow()
		}
	}
}

//去除/和跳转
func redirectTrailingSlash(c *Context) {
	req := c.Request
	path := req.URL.Path
	code := 301 // Permanent redirect, request with GET method
	if req.Method != "GET" {
		code = 307
	}

	if length := len(path); length > 1 && path[length-1] == '/' {
		req.URL.Path = path[:length-1]
	} else {
		req.URL.Path = path + "/"
	}
	debugPrint("redirecting request %d: %s --> %s", code, path, req.URL.String())
	http.Redirect(c.Writer, req, req.URL.String(), code)
	c.writermem.WriteHeaderNow()
}

//跳转到fixed path
func redirectFixedPath(c *Context, root *node, trailingSlash bool) bool {
	req := c.Request
	path := req.URL.Path

	fixedPath, found := root.findCaseInsensitivePath(
		cleanPath(path),
		trailingSlash,
	)
	if found {
		code := 301 // Permanent redirect, request with GET method
		if req.Method != "GET" {
			code = 307
		}
		req.URL.Path = string(fixedPath)
		debugPrint("redirecting request %d: %s --> %s", code, path, req.URL.String())
		http.Redirect(c.Writer, req, req.URL.String(), code)
		c.writermem.WriteHeaderNow()
		return true
	}
	return false
}
