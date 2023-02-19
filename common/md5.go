package common

import (
	"crypto/md5"
	"fmt"
)

// 将s进行md5加密
func MD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
