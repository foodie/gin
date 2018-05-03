// Copyright 2017 Bo-Yi Wu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// +build !jsoniter
//json wapper一下
package json

import "encoding/json"

//定义内部变量
var (
	//Marshal函数返回v的json编码
	Marshal = json.Marshal
	//类似Marshal但会使用缩进将输出格式化
	MarshalIndent = json.MarshalIndent
	//创建一个从r读取并解码json对象的*Decoder，
	//解码器有自己的缓冲，并可能超前读取部分json数据。
	NewDecoder = json.NewDecoder
)
