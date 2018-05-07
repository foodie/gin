// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"fmt"
	"net/http"
)

//以对应的状态码
type Redirect struct {
	Code     int
	Request  *http.Request
	Location string
}

//Redirect回复请求一个重定向地址urlStr和状态码code。
//该重定向地址可以是相对于请求r的相对地址。
func (r Redirect) Render(w http.ResponseWriter) error {
	if (r.Code < 300 || r.Code > 308) && r.Code != 201 {
		panic(fmt.Sprintf("Cannot redirect with status code %d", r.Code))
	}
	http.Redirect(w, r.Request, r.Location, r.Code)
	return nil
}

func (r Redirect) WriteContentType(http.ResponseWriter) {}
