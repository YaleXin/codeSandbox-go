package tool

import (
	"reflect"
)

// IsStructEmpty 判断任意结构体变量是否为空
func IsStructEmpty(s interface{}) bool {
	// 获取变量的反射值
	rv := reflect.ValueOf(s)
	// 如果变量是指针，获取其指向的元素
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	// 判断是否为结构体
	if rv.Kind() != reflect.Struct {
		return false
	}
	// 创建相同类型的零值实例
	zeroStruct := reflect.New(rv.Type()).Elem()
	// 比较变量与零值实例是否相等
	return reflect.DeepEqual(rv.Interface(), zeroStruct.Interface())
}
