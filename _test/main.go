package main

import (
	"fmt"
	"reflect"
)

type CustomInt int

func main() {
	// 获取 CustomInt 类型的切片类型
	t := reflect.SliceOf(reflect.TypeOf(CustomInt(0)))

	// 创建一个可寻址的切片值
	slice := reflect.New(t).Elem()

	// 使用 reflect.Append 追加元素
	elem := reflect.ValueOf(CustomInt(20))
	slice = reflect.Append(slice, elem)

	// 打印初始切片
	fmt.Println("Initial slice:", slice.Interface())

	// 使用 reflect.Value.Grow 扩展切片容量
	slice.Grow(10)

	// 打印扩展后的切片容量
	fmt.Println("Slice after expanding capacity:", slice.Interface())

	// 追加更多元素以展示动态增长
	for i := 1; i < 10; i++ {
		elem = reflect.ValueOf(CustomInt(20 + i))
		slice = reflect.Append(slice, elem)
	}

	// 打印结果
	fmt.Println("Slice after appending elements:", slice.Interface())
}
