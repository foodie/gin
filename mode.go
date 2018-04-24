// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"io"
	"os"

	"github.com/gin-gonic/gin/binding"
)

//gin的模式
const ENV_GIN_MODE = "GIN_MODE"

//定义基本的模式 debug,release,test三种模式
const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)

//设置三种类型的状态码 0 debug 1 release 2 testcode
const (
	debugCode = iota
	releaseCode
	testCode
)

// DefaultWriter is the default io.Writer used the Gin for debug output and
// middleware output like Logger() or Recovery().
// Note that both Logger and Recovery provides custom ways to configure their
// output io.Writer.
// To support coloring in Windows use:
// 		import "github.com/mattn/go-colorable"
// 		gin.DefaultWriter = colorable.NewColorableStdout()

//默认的输出和错误输出位置
var DefaultWriter io.Writer = os.Stdout
var DefaultErrorWriter io.Writer = os.Stderr

//默认是调试模式
var ginMode = debugCode
var modeName = DebugMode

//初始化设置模式
func init() {
	mode := os.Getenv(ENV_GIN_MODE)
	SetMode(mode)
}

//设置模式
func SetMode(value string) {
	switch value {
	case DebugMode, "":
		ginMode = debugCode
	case ReleaseMode:
		ginMode = releaseCode
	case TestMode:
		ginMode = testCode
	default:
		panic("gin mode unknown: " + value)
	}
	if value == "" {
		value = DebugMode
	}
	modeName = value
}

// 过滤
func DisableBindValidation() {
	binding.Validator = nil
}

//是否运行json解码
func EnableJsonDecoderUseNumber() {
	binding.EnableDecoderUseNumber = true
}

//返回当前的模式
func Mode() string {
	return modeName
}
