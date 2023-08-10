package util

import (
	"math/rand"
	"time"
)

// RandomString 生成随机字符串
func RandomString(n int) string {
	var letters = []byte("asdfghjkzxcvbnmqwertyuioASDFGHJKLZXCVBNMQWERTYUIO")
	result := make([]byte, n)
	//初始化随机数生成器
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
