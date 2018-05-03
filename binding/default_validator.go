// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"reflect"
	"sync"

	"gopkg.in/go-playground/validator.v8"
)

//定义默认的validator
type defaultValidator struct {
	once     sync.Once           //运行一次
	validate *validator.Validate //定义validate
}

//实例化一次
var _ StructValidator = &defaultValidator{}

//验证结构
func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	//如果是结构
	if kindOfData(obj) == reflect.Struct {
		//初始化一次
		v.lazyinit()
		//验证是否是结构体
		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}
	return nil
}

//注册一个key  Validator.func
func (v *defaultValidator) RegisterValidation(key string, fn validator.Func) error {
	v.lazyinit()
	return v.validate.RegisterValidation(key, fn)
}

//新建一个validate，初始化一次
func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		config := &validator.Config{TagName: "binding"}
		v.validate = validator.New(config)
	})
}

//获取数据的类型
func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
