// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

//默认的状态
const (
	noWritten     = -1  //无写操作
	defaultStatus = 200 //状态码
)

//基本的接口
type ResponseWriter interface {
	//header() writer() writeHeader()
	http.ResponseWriter //包含http接口
	// Hijack()
	http.Hijacker //可以让HTTP处理器接管该连接。
	//Flush()
	http.Flusher //Flush将缓冲中的所有数据发送到客户端
	// CloseNotify返回一个通道，该通道会在客户端连接丢失时接收到唯一的值
	//CloseNotify() <-chan bool
	http.CloseNotifier //可以让用户检测下层的连接是否停止。

	//http状态
	// Returns the HTTP response status code of the current request.
	Status() int

	//返回数据的字节数
	// Returns the number of bytes already written into the response http body.
	// See Written()
	Size() int

	//写入一个string
	// Writes the string into the response body.
	WriteString(string) (int, error)

	//是否可写
	// Returns true if the response body was already written.
	Written() bool

	//往头部写内容
	// Forces to write the http header (status code + headers).
	WriteHeaderNow()
}

//内部的写
//对http.ResponseWriter 进行封装
type responseWriter struct {
	http.ResponseWriter     //包含http.ResponseWrite接口
	size                int //尺寸
	status              int //状态码
}

//内部定义一个空的
var _ ResponseWriter = &responseWriter{}

//初始化写
func (w *responseWriter) reset(writer http.ResponseWriter) {
	w.ResponseWriter = writer
	w.size = noWritten
	w.status = defaultStatus
}

//实现http.ResponseWriter 接口
//设置header的状态
func (w *responseWriter) WriteHeader(code int) {
	if code > 0 && w.status != code {
		//是否可写
		if w.Written() {
			debugPrint("[WARNING] Headers were already written. Wanted to override status code %d with %d", w.status, code)
		}
		//设置状态码
		w.status = code
	}
}

//实现http.ResponseWriter 接口
//如果不可写，可以直接写，写状态码
//
func (w *responseWriter) WriteHeaderNow() {
	if !w.Written() {
		w.size = 0
		w.ResponseWriter.WriteHeader(w.status)
	}
}

//设置header
//写入数据，记录写入的字节数
func (w *responseWriter) Write(data []byte) (n int, err error) {
	w.WriteHeaderNow()
	//写入数据
	n, err = w.ResponseWriter.Write(data)
	w.size += n
	return
}

//写入string
func (w *responseWriter) WriteString(s string) (n int, err error) {
	w.WriteHeaderNow()
	//通过io往ResponseWriter里面写数据
	n, err = io.WriteString(w.ResponseWriter, s)
	w.size += n
	return
}

//status
func (w *responseWriter) Status() int {
	return w.status
}

//size
func (w *responseWriter) Size() int {
	return w.size
}

//判断是否可写
func (w *responseWriter) Written() bool {
	return w.size != noWritten
}

//主要用来初始化size
// Hijack implements the http.Hijacker interface.
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.size < 0 {
		w.size = 0
	}
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

//使用默认的通知代码
//完成的通知
// CloseNotify implements the http.CloseNotify interface.
func (w *responseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

//使用默认的flush代码
// Flush implements the http.Flush interface.
func (w *responseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}
