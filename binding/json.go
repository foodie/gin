// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"net/http"

	"github.com/gin-gonic/gin/json"
)

//UseNumber方法将dec设置为当接收端
//是interface{}接口时将json数字解码为Number类型
var EnableDecoderUseNumber = false

type jsonBinding struct{}

//返回名字
func (jsonBinding) Name() string {
	return "json"
}

//读取数据，解压数据
func (jsonBinding) Bind(req *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(req.Body)
	if EnableDecoderUseNumber {
		decoder.UseNumber()
	}
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	//验证数据？
	return validate(obj)
}
