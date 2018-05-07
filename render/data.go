// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

//基础数据
import "net/http"

type Data struct {
	ContentType string //内容类型
	Data        []byte //数据
}

// Render (Data) writes data with custom ContentType.
func (r Data) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	_, err = w.Write(r.Data)
	return
}

//
func (r Data) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, []string{r.ContentType})
}
