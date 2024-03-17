package utils

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_ParseFloat(t *testing.T) {
	var myFloat64 float64 = 12345.6789
	// 将float64转换为字符串
	// 'f' 表示打印格式（不使用科学计数法）
	// -1 表示在小数点后打印尽可能多的小数位数，但这也可以设置为特定的位数，比如2
	// 64 表示这是一个64位的浮点数
	str := strconv.FormatFloat(myFloat64, 'f', -1, 64)
	fmt.Println("Float64 as string:", str)
}

func Test_ParseInteger(t *testing.T) {
	var (
		v32 int32 = 100
	)
	fmt.Println(FormatInteger(v32))
}
