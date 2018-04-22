// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package gin

import (
	"net/http"
	"os"
)

//FileSystem接口实现了对一系列命名文件的访问。
type onlyfilesFS struct {
	fs http.FileSystem
}

//File是被FileSystem接口的Open方法返回的接口类型
//可以被FileServer等函数用于文件访问服务。
type neuteredReaddirFile struct {
	http.File
}

// Dir returns a http.Filesystem that can be used by http.FileServer(). It is used internally
// in router.Static().
// if listDirectory == true, then it works the same as http.Dir() otherwise it returns
// a filesystem that prevents http.FileServer() to list the directory files.
func Dir(root string, listDirectory bool) http.FileSystem {
	//转换成http.Dir
	fs := http.Dir(root)
	//显示的话直接返回
	if listDirectory {
		return fs
	}
	//返回变量
	return &onlyfilesFS{fs}
}

//使用open方法，返回一个http.File的接口
// Open conforms to http.Filesystem.
func (fs onlyfilesFS) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	//打开成功返回一个http.FILE
	return neuteredReaddirFile{f}, nil
}

//读取目录
// Readdir overrides the http.File default implementation.
func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	// this disables directory listing
	return nil, nil
}
