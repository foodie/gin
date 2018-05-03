// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import "net/http"

//默认的内存使用32M
const defaultMemory = 32 * 1024 * 1024

//定义formbind
type formBinding struct{}

//form数据
type formPostBinding struct{}

//multipart数据
type formMultipartBinding struct{}

//获取名字
func (formBinding) Name() string {
	return "form"
}

func (formBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	//请求的主体作为multipart/form-data解析
	req.ParseMultipartForm(defaultMemory)
	//解析form，验证form
	if err := mapForm(obj, req.Form); err != nil {
		return err
	}
	return validate(obj)
}

func (formPostBinding) Name() string {
	return "form-urlencoded"
}

//解析表单
func (formPostBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := mapForm(obj, req.PostForm); err != nil {
		return err
	}
	return validate(obj)
}

//form-data
func (formMultipartBinding) Name() string {
	return "multipart/form-data"
}

//基本的校验
func (formMultipartBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseMultipartForm(defaultMemory); err != nil {
		return err
	}
	if err := mapForm(obj, req.MultipartForm.Value); err != nil {
		return err
	}
	return validate(obj)
}
