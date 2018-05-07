// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import "net/http"

//定义接口，处理http.ResponseWriter
type Render interface {
	Render(http.ResponseWriter) error
	WriteContentType(w http.ResponseWriter)
}

//定义各种空对象
var (
	_ Render     = JSON{}
	_ Render     = IndentedJSON{}
	_ Render     = SecureJSON{}
	_ Render     = XML{}
	_ Render     = String{}
	_ Render     = Redirect{}
	_ Render     = Data{}
	_ Render     = HTML{}
	_ HTMLRender = HTMLDebug{}
	_ HTMLRender = HTMLProduction{}
	_ Render     = YAML{}
	_ Render     = MsgPack{}
)

//设置header
//如果Content-Type 不存在，写Content-Type
func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
