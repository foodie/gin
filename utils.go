// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"encoding/xml"
	"net/http"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
)

//默认的key
const BindKey = "_gin-gonic/gin/bindkey"

//设置_gin-gonic/gin/bindkey的值
func Bind(val interface{}) HandlerFunc {
	value := reflect.ValueOf(val)
	//如果说是指针，报错
	if value.Kind() == reflect.Ptr {
		panic(`Bind struct can not be a pointer. Example:
	Use: gin.Bind(Struct{}) instead of gin.Bind(&Struct{})
`)
	}
	//获取类型
	typ := value.Type()

	//设置绑定的值
	return func(c *Context) {
		//本方法返回val当前持有的值
		obj := reflect.New(typ).Interface()
		if c.Bind(obj) == nil {
			c.Set(BindKey, obj)
		}
	}
}

//把http.HandlerFunc转换成HandlerFunc
// WrapF is a helper function for wrapping http.HandlerFunc
// Returns a Gin middleware
func WrapF(f http.HandlerFunc) HandlerFunc {
	return func(c *Context) {
		f(c.Writer, c.Request)
	}
}

//把http.Handler转换成 HandlerFunc
// WrapH is a helper function for wrapping http.Handler
// Returns a Gin middleware
func WrapH(h http.Handler) HandlerFunc {
	return func(c *Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

//定义一个大的map
// H is a shortcup for map[string]interface{}
type H map[string]interface{}

//写入xml
// MarshalXML allows type H to be used with xml.Marshal.
func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	//xml 的名字
	start.Name = xml.Name{
		Space: "",
		Local: "map",
	}
	//向底层写入一个token，解析错误
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	//循环调用
	for key, value := range h {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}
	//结束标签
	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

//自己的assert
func assert1(guard bool, text string) {
	if !guard {
		panic(text)
	}
}

//根据' '或者;截断代码
func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

//如果custom为空，wildcard不为空返回wildcard
//否则返回custom
func chooseData(custom, wildcard interface{}) interface{} {
	if custom == nil {
		if wildcard == nil {
			panic("negotiation config is invalid")
		}
		return wildcard
	}
	return custom
}

func parseAccept(acceptHeader string) []string {
	parts := strings.Split(acceptHeader, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if index := strings.IndexByte(part, ';'); index >= 0 {
			part = part[0:index]
		}
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}

//最后一个char
func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

/**
返回调用的函数名

返回一个表示调用栈标识符pc对应的调用栈的*Func；
如果该调用栈标识符没有对应的调用栈，函数会返回nil。
每一个调用栈必然是对某个函数的调用。
**/
func nameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

//绝对路径，相对路径，如果relativePath后有/返回的也是有/的
func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}

//对addr进行处理，默认返回:8080
func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			debugPrint("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		debugPrint("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too much parameters")
	}
}
