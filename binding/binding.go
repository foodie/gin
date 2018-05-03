// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"net/http"

	"gopkg.in/go-playground/validator.v8"
)

//定义结构的基本类型
const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	//其他类型
	MIMEPROTOBUF = "application/x-protobuf"
	MIMEMSGPACK  = "application/x-msgpack"
	MIMEMSGPACK2 = "application/msgpack"
)

//基本的接口，名字和绑定
type Binding interface {
	Name() string
	Bind(*http.Request, interface{}) error
}

//验证结构
type StructValidator interface {
	// ValidateStruct can receive any kind of type and it should never panic, even if the configuration is not right.
	// If the received type is not a struct, any validation should be skipped and nil must be returned.
	// If the received type is a struct or pointer to a struct, the validation should be performed.
	// If the struct is not valid or the validation itself fails, a descriptive error should be returned.
	// Otherwise nil must be returned.
	//基本方法
	ValidateStruct(interface{}) error

	// RegisterValidation adds a validation Func to a Validate's map of validators denoted by the key
	// NOTE: if the key already exists, the previous validation function will be replaced.
	// NOTE: this method is not thread-safe it is intended that these all be registered prior to any validation
	//注册validator方法
	RegisterValidation(string, validator.Func) error
}

//基本的验证
var Validator StructValidator = &defaultValidator{}

//定义其他类型的验证方法
var (
	JSON          = jsonBinding{}
	XML           = xmlBinding{}
	Form          = formBinding{}
	Query         = queryBinding{}
	FormPost      = formPostBinding{}
	FormMultipart = formMultipartBinding{}
	ProtoBuf      = protobufBinding{}
	MsgPack       = msgpackBinding{}
)

//默认的处理
//根据header里面的数据返回对应的验证类型
func Default(method, contentType string) Binding {
	if method == "GET" {
		return Form
	}

	switch contentType {
	case MIMEJSON:
		return JSON
	case MIMEXML, MIMEXML2:
		return XML
	case MIMEPROTOBUF:
		return ProtoBuf
	case MIMEMSGPACK, MIMEMSGPACK2:
		return MsgPack
	default: //case MIMEPOSTForm, MIMEMultipartPOSTForm:
		return Form
	}
}

//验证数据的结构
func validate(obj interface{}) error {
	if Validator == nil {
		return nil
	}
	return Validator.ValidateStruct(obj)
}
