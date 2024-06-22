package tool

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

func IsBlankString(str string) bool {
	for _, r := range str {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return len(str) == 0 || strings.TrimSpace(str) == ""
}
func MD5Str(inputStr string) string {
	// 创建一个新的MD5哈希实例
	hasher := md5.New()

	// 将字符串转换为字节切片并写入哈希对象
	hasher.Write([]byte(inputStr))

	// 计算哈希值
	hashBytes := hasher.Sum(nil)

	// 将字节切片转换为16进制字符串表示形式
	return hex.EncodeToString(hashBytes)
}

// 可见ASCII字符集（不包括空格）
const visibleChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"

// 随机字符串生成函数
func GenerateRandomVisibleString(length int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, length)
	for i := range b {
		b[i] = visibleChars[rand.Intn(len(visibleChars))]
	}
	return string(b)
}
