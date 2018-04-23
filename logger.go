// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mattn/go-isatty"
)

//定义颜色的字符串
var (
	green        = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white        = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow       = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red          = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue         = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta      = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan         = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset        = string([]byte{27, 91, 48, 109})
	disableColor = false
)

//禁用颜色
// DisableConsoleColor disables color output in the console.
func DisableConsoleColor() {
	disableColor = true
}

//报错日志
// ErrorLogger returns a handlerfunc for any error type.
func ErrorLogger() HandlerFunc {
	return ErrorLoggerT(ErrorTypeAny)
}

//处理任意的错误类型
// ErrorLoggerT returns a handlerfunc for a given error type.
func ErrorLoggerT(typ ErrorType) HandlerFunc {
	//
	return func(c *Context) {
		c.Next() //处理下一个
		//获取error，线上error
		errors := c.Errors.ByType(typ)
		if len(errors) > 0 {
			c.JSON(-1, errors)
		}
	}
}

//写日志，中间件
// Logger instances a Logger middleware that will write the logs to gin.DefaultWriter.
// By default gin.DefaultWriter = os.Stdout.
func Logger() HandlerFunc {
	return LoggerWithWriter(DefaultWriter)
}

// LoggerWithWriter instance a Logger middleware with the specified writter buffer.
// Example: os.Stdout, a file opened in write mode, a socket...
func LoggerWithWriter(out io.Writer, notlogged ...string) HandlerFunc {
	//是否是iterm
	isTerm := true

	if w, ok := out.(*os.File); !ok || //不是文件
		(os.Getenv("TERM") == "dumb" || //环境是dumb
			(!isatty.IsTerminal(w.Fd()) && //不是终端
				!isatty.IsCygwinTerminal(w.Fd()))) || //不是终端
		disableColor { //不显示颜色
		isTerm = false
	}
	//需要跳过
	var skip map[string]struct{}

	//记录不写日志的路径
	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *Context) {
		// Start timer
		//开始时间，路径，查询串
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		//处理情况请求
		// Process request
		c.Next()

		//不是需要跳过的
		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// Stop timer
			//结束时间
			end := time.Now()
			//访问时间
			latency := end.Sub(start)

			//客户端ip
			clientIP := c.ClientIP()
			//方法
			method := c.Request.Method
			//状态码
			statusCode := c.Writer.Status()
			//设置各种颜色
			var statusColor, methodColor, resetColor string
			if isTerm {
				statusColor = colorForStatus(statusCode)
				methodColor = colorForMethod(method)
				resetColor = reset
			}
			//获取注释内容
			comment := c.Errors.ByType(ErrorTypePrivate).String()

			//得到path
			if raw != "" {
				path = path + "?" + raw
			}

			//输出日志
			fmt.Fprintf(out, "[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %s\n%s",
				end.Format("2006/01/02 - 15:04:05"),
				statusColor, statusCode, resetColor,
				latency,
				clientIP,
				methodColor, method, resetColor,
				path,
				comment,
			)
		}
	}
}

//根据状态码获取对应的颜色
func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

//根据字符串获取对应的颜色
func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}
