package utils

import (
	"fmt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"time"
)

var Period uint = 300               //过期时间 秒
var SecretSize uint = 64            //秘钥长度
var Digits = otp.DigitsEight        //生成现在的秘钥长度
var Algorithm = otp.AlgorithmSHA512 //加密算法

func Test() {
	//生成秘钥
	secret := Generate("userId")
	fmt.Println(secret)

	//获取当前code
	code := GetNowCode(secret)
	fmt.Println(code)

	i := 0
	for i < 600 {

		//获取当前code
		newCode := GetNowCode(secret)
		fmt.Println("code:", code, newCode)
		//校验code可不可以通过
		pass, _ := ValidateCode(code, secret)
		//if !pass {
		//	code = newCode
		//}
		newPass, _ := ValidateCode(newCode, secret)
		//fmt.Println(err)
		fmt.Println("----------------------------------------", i, pass, newPass)
		time.Sleep(time.Second)
		i++
	}
}

// Generate 生成秘钥
func Generate(code string) string {
	data := totp.GenerateOpts{
		Period:      Period,
		Digits:      Digits,
		Algorithm:   Algorithm,
		AccountName: "i-confluence.com",
		Issuer:      code,
		SecretSize:  SecretSize,
	}
	key, err := totp.Generate(data)
	if err != nil {
		fmt.Println(err)
	}
	return key.Secret()
}

// GetNowCode 获取当前秘钥
func GetNowCode(secret string) string {
	code, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period:    Period,
		Skew:      2,
		Digits:    Digits,
		Algorithm: Algorithm,
	})
	if err != nil {
		fmt.Println(err)
	}
	return code
}

// ValidateCode 校验code是否正常
func ValidateCode(code, secret string) (bool, error) {
	return totp.ValidateCustom(code, secret, time.Now(), totp.ValidateOpts{
		Period:    Period,
		Skew:      2,
		Digits:    Digits,
		Algorithm: Algorithm,
	})
}
