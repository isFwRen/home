/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/4/6 11:47
 */

package utils

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

//CreateRandomString 生成指定长度的字符串
func CreateRandomString(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}
