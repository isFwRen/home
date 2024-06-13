package response

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"server/global"
	"server/utils"
)

type ExportResponse struct {
	Code int    `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"msg"`
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type Response2 struct {
	Code int         `json:"code"`
	List interface{} `json:"list"`
	Msg  string      `json:"msg"`
}

const (
	ERROR    = 400
	SUCCESS  = 200
	NotLogin = 401
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Result2(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response2{
		code,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "操作成功", c)
}

func GetOkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "获取成功", c)
}

func GetFailWithData(data interface{}, c *gin.Context) {
	Result(ERROR, data, "获取失败", c)
}

func OKWithLogin(date interface{}, c *gin.Context) {
	Result(SUCCESS, date, "登录成功", c)
}

func OkDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func OkDetailed2(data interface{}, message string, c *gin.Context) {
	Result2(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func FailWithParamErr(err error, c *gin.Context) {
	Result(ERROR, err.Error(), global.ParamErr.Error(), c)
}

func FailWithDetailed(code int, data interface{}, message string, c *gin.Context) {
	Result(code, data, message, c)
}
func FailWithoutCode(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}

//处理雪花算法生成id序列化精度不足
func coverDate(data interface{}) (reData interface{}) {
	j, _ := json.Marshal(data)
	reg := regexp.MustCompile(`ID\":(\d{16,20}),"`)
	l := len(reg.FindAllString(string(j), -1)) //正则匹配16-20位的数字，如果找到了就开始正则替换并解析
	if l != 0 {
		//fmt.Printf("\n正则替换前的数据%+v", data)
		var mapResult map[string]interface{}
		str := reg.ReplaceAllString(string(j), `ID": "${1}","`)
		err := json.Unmarshal([]byte(str), &mapResult)
		if err != nil {
			panic(fmt.Sprintf("coverDate>>>>%v", err))
		}
		data = &mapResult
	}
	return data
}

func OkDetailedEncrypt(data interface{}, message string, c *gin.Context) {
	key := "xingqiyistronger"
	marshal, err := json.Marshal(data)
	if err != nil {
		return
	}
	encrypt, err := utils.AesEncrypt(marshal, []byte(key))
	if err != nil {
		return
	}
	Result(SUCCESS, encrypt, message, c)
}
