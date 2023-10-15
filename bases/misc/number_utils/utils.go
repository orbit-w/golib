package number_utils

import "github.com/orbit-w/golib/v1/bases/misc/common"

/*
   @Time: 2023/8/22 00:17
   @Author: david
   @File: utils
*/

func Min[T common.Integer](a, b T) T {
	if a < b {
		return a
	}
	return b
}
