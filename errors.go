// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin/json"
)

//定义错误类型
type ErrorType uint64

//错误类型的值
const (
	//bing错误
	ErrorTypeBind ErrorType = 1 << 63 // used when c.Bind() fails
	//private错误
	ErrorTypeRender ErrorType = 1 << 62 // used when c.Render() fails
	//private 错误
	ErrorTypePrivate ErrorType = 1 << 0
	//public 错误
	ErrorTypePublic ErrorType = 1 << 1
	//任意错误
	ErrorTypeAny ErrorType = 1<<64 - 1
	//nu错误
	ErrorTypeNu = 2
)

//错误类型
type Error struct {
	Err  error       //定义error
	Type ErrorType   //错误类型
	Meta interface{} //处理元数据
}

//定义Error数组
type errorMsgs []*Error

//为啥要这么做呢？
var _ error = &Error{}

//设置类型
func (msg *Error) SetType(flags ErrorType) *Error {
	msg.Type = flags
	return msg
}

//设置元数据
func (msg *Error) SetMeta(data interface{}) *Error {
	msg.Meta = data
	return msg
}

//返回json字符串
func (msg *Error) JSON() interface{} {
	json := H{} //定义基本的串
	//获取meta的结构
	if msg.Meta != nil {
		value := reflect.ValueOf(msg.Meta)

		//获取基本类型
		switch value.Kind() {
		case reflect.Struct:
			return msg.Meta //结构的话直接返回
		case reflect.Map:
			//如果是map直接放入
			for _, key := range value.MapKeys() {
				json[key.String()] = value.MapIndex(key).Interface()
			}
		default: //其他类型
			json["meta"] = msg.Meta
		}
	}
	if _, ok := json["error"]; !ok {
		json["error"] = msg.Error()
	}
	return json
}

//对msg进行json编码
// MarshalJSON implements the json.Marshaller interface.
func (msg *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.JSON())
}

//返回error
// Error implements the error interface
func (msg Error) Error() string {
	return msg.Err.Error()
}

//type&flag 是否大于0
func (msg *Error) IsType(flags ErrorType) bool {
	return (msg.Type & flags) > 0
}

//
// ByType returns a readonly copy filtered the byte.
// ie ByType(gin.ErrorTypePublic) returns a slice of errors with type=ErrorTypePublic.
func (a errorMsgs) ByType(typ ErrorType) errorMsgs {
	if len(a) == 0 {
		return nil
	}
	//nil和任意类型直接返回
	if typ == ErrorTypeAny {
		return a
	}
	//返回符合类型的result
	var result errorMsgs
	for _, msg := range a {
		if msg.IsType(typ) {
			result = append(result, msg)
		}
	}
	return result
}

//返回最后一个error
// Last returns the last error in the slice. It returns nil if the array is empty.
// Shortcut for errors[len(errors)-1].
func (a errorMsgs) Last() *Error {
	if length := len(a); length > 0 {
		return a[length-1]
	}
	return nil
}

// Errors returns an array will all the error messages.
// Example:
// 		c.Error(errors.New("first"))
// 		c.Error(errors.New("second"))
// 		c.Error(errors.New("third"))
// 		c.Errors.Errors() // == []string{"first", "second", "third"}
//返回一个string slice
func (a errorMsgs) Errors() []string {
	if len(a) == 0 {
		return nil
	}
	errorStrings := make([]string, len(a))
	for i, err := range a {
		errorStrings[i] = err.Error()
	}
	return errorStrings
}

//对errorMsgs json串
func (a errorMsgs) JSON() interface{} {
	switch len(a) {
	case 0:
		return nil
	case 1:
		return a.Last().JSON()
	default:
		json := make([]interface{}, len(a))
		for i, err := range a {
			json[i] = err.JSON()
		}
		return json
	}
}

//json压缩
func (a errorMsgs) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.JSON())
}

//对 errormsgs进行json化
func (a errorMsgs) String() string {
	if len(a) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for i, msg := range a {
		fmt.Fprintf(&buffer, "Error #%02d: %s\n", i+1, msg.Err)
		if msg.Meta != nil {
			fmt.Fprintf(&buffer, "     Meta: %v\n", msg.Meta)
		}
	}
	return buffer.String()
}
