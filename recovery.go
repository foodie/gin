// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http/httputil"
	"runtime"
)

//?
var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// a middleware
// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() HandlerFunc {
	return RecoveryWithWriter(DefaultErrorWriter)
}

// RecoveryWithWriter returns a middleware for a given writer that recovers from any panics and writes a 500 if there was one.
func RecoveryWithWriter(out io.Writer) HandlerFunc {
	//日志类
	var logger *log.Logger
	if out != nil {
		logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
	}
	return func(c *Context) {
		//最后的时候调用，用于错误统一处理，数据的返回
		defer func() {
			//用来恢复错误
			if err := recover(); err != nil {
				if logger != nil {
					//?
					stack := stack(3)
					//DumpRequest返回req的和
					//被服务端接收到时一样的有线表示，可选地包括请求的主体，用于debug。
					httprequest, _ := httputil.DumpRequest(c.Request, false)
					//输出错误
					logger.Printf("[Recovery] panic recovered:\n%s\n%s\n%s%s", string(httprequest), err, stack, reset)
				}
				//500错误
				c.AbortWithStatus(500)
			}
		}()
		//下一个
		c.Next()
	}
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	//定义buffer
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		//报告当前go程调用栈所执行的函数的文件和行号信息。
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		//打印基本的信息
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		//如果不是上一个文件
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			//使用\n分行
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		//，打印出错的行
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	//返回buf数据
	return buf.Bytes()
}

//清除lines[n]的空格
// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

//获取函数名字
// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	//FuncForPC返回一个表示调用栈标识符pc对应的调用栈的*Func；
	//如果该调用栈标识符没有对应的调用栈，函数会返回nil。
	//每一个调用栈必然是对某个函数的调用
	fn := runtime.FuncForPC(pc)
	//为空
	if fn == nil {
		return dunno
	}
	//获取基本的名字
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	//得到最后的/获取名字
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	//通过.获取名字
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	//替换.为.
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
