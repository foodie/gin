// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"html/template"
	"net/http"
)

//切分字段
type Delims struct {
	Left  string
	Right string
}

//新建一个htmlRender
type HTMLRender interface {
	Instance(string, interface{}) Render
}

/**
Delims方法用于设置action的分界字符串，
应用于之后的Parse、ParseFiles、ParseGlob方法。
嵌套模板定义会继承这种分界符设置。
空字符串分界符表示相应的默认分界符：{{或}}。
返回值就是t，以便进行链式调用。
**/
//html基本的
type HTMLProduction struct {
	Template *template.Template
	Delims   Delims
}

//html debug
type HTMLDebug struct {
	Files   []string
	Glob    string
	Delims  Delims
	FuncMap template.FuncMap
}

//html
type HTML struct {
	Template *template.Template
	Name     string
	Data     interface{}
}

//基本的text/html
var htmlContentType = []string{"text/html; charset=utf-8"}

//通过HTMLProduction创建
func (r HTMLProduction) Instance(name string, data interface{}) Render {
	return HTML{
		Template: r.Template,
		Name:     name,
		Data:     data,
	}
}

//通过HTMLDebug创建
func (r HTMLDebug) Instance(name string, data interface{}) Render {
	return HTML{
		Template: r.loadTemplate(),
		Name:     name,
		Data:     data,
	}
}

//处理文件和Glob
func (r HTMLDebug) loadTemplate() *template.Template {
	if r.FuncMap == nil {
		r.FuncMap = template.FuncMap{}
	}
	//ParseGlob方法解析filenames指定的文件里的模板定义并将解析结果与t关联
	if len(r.Files) > 0 {
		return template.Must(template.New("").Delims(r.Delims.Left, r.Delims.Right).Funcs(r.FuncMap).ParseFiles(r.Files...))
	}
	//ParseFiles方法解析匹配pattern的文件里的模板定义并将解析结果与t关联。
	if r.Glob != "" {
		return template.Must(template.New("").Delims(r.Delims.Left, r.Delims.Right).Funcs(r.FuncMap).ParseGlob(r.Glob))
	}
	panic("the HTML debug render was created without files or glob pattern")
}

//解析和执行文件
func (r HTML) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	if r.Name == "" {
		return r.Template.Execute(w, r.Data)
	}
	return r.Template.ExecuteTemplate(w, r.Name, r.Data)
}

func (r HTML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, htmlContentType)
}
