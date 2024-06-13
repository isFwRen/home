/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/4/7 09:53
 */

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"reflect"
	"server/global"
	model2 "server/module/sys_base/model"
	"strconv"
)

// HttpRequest 发起http请求
func HttpRequest(url string, v interface{}) (error, []byte) {
	buf, err := json.Marshal(v)
	if err != nil {
		return err, nil
	}
	//client := &http.Client{Timeout: 10 * time.Second}
	client := &http.Client{}
	res, err := http.NewRequest("POST", url, bytes.NewBuffer(buf))
	if err != nil {
		return err, nil
	}
	res.Header.Set("Content-Type", "application/json; charset=utf-8")
	//res.Header.Set("x-user-id", "")
	//res.Header.Set("pro-code", "")
	//res.Header.Set("process", "")
	//res.Header.Set("x-token", "")
	//global.GLog.Info("HttpRequest", zap.Any("Do", "before"))
	resp, err := client.Do(res)
	//global.GLog.Info("HttpRequest", zap.Any("Do", "end"))
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return err, body
}

type Resp struct {
	List   []map[string]interface{} `json:"list"`
	Total  int64                    `json:"total"`
	Status int                      `json:"status"`
	Msg    string                   `json:"msg"`
}

type Req struct {
	ProCode    string            `json:"proCode"`
	Name       string            `json:"name"`
	QueryNames map[string]string `json:"queryNames"`
	RespNames  []string          `json:"respNames"`
	Sort       []map[string]int  `json:"sort"`
	model2.BasePageInfo
}

// FetchConst 常量匹配
//
// Parameters:
//
//	(string): 项目编码
//	(string): 常量名字
//	(string): 返回字段
//	(map[string]string): 查询字段和其值
//
// Returns:
//
//	(string): 匹配的值
//	(int64):  更新条数
//
// Description:
//
//	在常量服务查询匹配常量
func FetchConst(proCode, constName, respName string, queryName map[string]string) (val string, total int64) {
	resp := Resp{}
	//url := "http://127.0.0.1:30080/sys-const/page"
	//url := "http://127.0.0.1:8888/pro-manager/bill-list/page?proCode=B0108&pageIndex=1&pageSize=10&timeStart=2023-06-30T16:00:00.000Z&timeEnd=2023-07-17T15:59:59.000Z"
	req := Req{
		ProCode:    proCode,
		Name:       constName,
		QueryNames: queryName,
		RespNames:  []string{respName},
	}
	req.PageSize = 1
	req.PageIndex = 1
	for k, _ := range queryName {
		req.Sort = append(req.Sort, map[string]int{
			k: 1,
		})
	}
	//global.GLog.Info("请求参数", zap.Any("req", req))
	err, body := HttpRequest(global.GConfig.System.ConstUrl+"/sys-const/page", req)
	if err != nil {
		global.GLog.Error("HttpRequest", zap.Error(err))
		return "", 0
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		global.GLog.Error("Unmarshal", zap.Error(err))
		return "", 0
	}
	//global.GLog.Info("响应数据", zap.Any("body", resp))

	if len(resp.List) > 0 {
		if _, ok := resp.List[0][respName]; ok {
			switch reflect.ValueOf(resp.List[0][respName]).Kind() {
			case reflect.String:
				return resp.List[0][respName].(string), resp.Total
			case reflect.Bool:
				return strconv.FormatBool(resp.List[0][respName].(bool)), resp.Total
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return strconv.Itoa(resp.List[0][respName].(int)), resp.Total
			case reflect.Float32, reflect.Float64:
				return fmt.Sprintf("%.4f", resp.List[0][respName].(float64)), resp.Total
			}
		} else {
			global.GLog.Error("没有该字段", zap.Any(constName, respName))
			return "", 0
		}
		return resp.List[0][respName].(string), resp.Total
	}
	global.GLog.Error("resp.List is nil", zap.Any("", req))
	return "", 0
}
