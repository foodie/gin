// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ginS

import (
	"html/template"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var once sync.Once
var internalEngine *gin.Engine

//初始化引擎一次
func engine() *gin.Engine {
	once.Do(func() {
		internalEngine = gin.Default()
	})
	return internalEngine
}

//加载Glob
func LoadHTMLGlob(pattern string) {
	engine().LoadHTMLGlob(pattern)
}

//加载html files
func LoadHTMLFiles(files ...string) {
	engine().LoadHTMLFiles(files...)
}

//设置html模板
func SetHTMLTemplate(templ *template.Template) {
	engine().SetHTMLTemplate(templ)
}

//设置no router方法
// NoRoute adds handlers for NoRoute. It return a 404 code by default.
func NoRoute(handlers ...gin.HandlerFunc) {
	engine().NoRoute(handlers...)
}

//设置no method方法
// NoMethod sets the handlers called when... TODO
func NoMethod(handlers ...gin.HandlerFunc) {
	engine().NoMethod(handlers...)
}

//设置group
// Group creates a new router group. You should add all the routes that have common middlwares or the same path prefix.
// For example, all the routes that use a common middlware for authorization could be grouped.
func Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return engine().Group(relativePath, handlers...)
}

//处理
func Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return engine().Handle(httpMethod, relativePath, handlers...)
}

//post处理
// POST is a shortcut for router.Handle("POST", path, handle)
func POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return engine().POST(relativePath, handlers...)
}

//get处理
// GET is a shortcut for router.Handle("GET", path, handle)
func GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return engine().GET(relativePath, handlers...)
}

//delete
// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return engine().DELETE(relativePath, handlers...)
}

//patch
// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func PATCH(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return engine().PATCH(relativePath, handlers...)
}

//put处理
// PUT is a shortcut for router.Handle("PUT", path, handle)
func PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return engine().PUT(relativePath, handlers...)
}

//选项处理
// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func OPTIONS(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return engine().OPTIONS(relativePath, handlers...)
}

//头部处理
// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func HEAD(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return engine().HEAD(relativePath, handlers...)
}

//任意处理
func Any(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return engine().Any(relativePath, handlers...)
}

//静态文件
func StaticFile(relativePath, filepath string) gin.IRoutes {
	return engine().StaticFile(relativePath, filepath)
}

// Static serves files from the given file system root.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use :
//     router.Static("/static", "/var/www")
//静态文件
func Static(relativePath, root string) gin.IRoutes {
	return engine().Static(relativePath, root)
}

//静态文件
func StaticFS(relativePath string, fs http.FileSystem) gin.IRoutes {
	return engine().StaticFS(relativePath, fs)
}

//使用中间件
// Use attachs a global middleware to the router. ie. the middlewares attached though Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func Use(middlewares ...gin.HandlerFunc) gin.IRoutes {
	return engine().Use(middlewares...)
}

//按照地址运行
// Run : The router is attached to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine undefinitelly unless an error happens.
func Run(addr ...string) (err error) {
	return engine().Run(addr...)
}

//运行https
// RunTLS : The router is attached to a http.Server and starts listening and serving HTTPS requests.
// It is a shortcut for http.ListenAndServeTLS(addr, certFile, keyFile, router)
// Note: this method will block the calling goroutine undefinitelly unless an error happens.
func RunTLS(addr string, certFile string, keyFile string) (err error) {
	return engine().RunTLS(addr, certFile, keyFile)
}

//运行unix
// RunUnix : The router is attached to a http.Server and starts listening and serving HTTP requests
// through the specified unix socket (ie. a file)
// Note: this method will block the calling goroutine undefinitelly unless an error happens.
func RunUnix(file string) (err error) {
	return engine().RunUnix(file)
}
