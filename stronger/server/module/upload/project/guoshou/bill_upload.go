/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/12/6 17:43
 */

package guoshou

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"server/global"
	model2 "server/module/pro_conf/model"
	"server/module/pro_manager/model"
	"server/module/pro_manager/service"
	service1 "server/module/upload/service"
	"server/utils"
	"strings"
	"time"

	"go.uber.org/zap"
)

//var lock sync.Mutex

// BillUpload 回传单据
func BillUpload(reqParam model.ProCodeAndId, uploadPath model2.SysProDownloadPaths) error {
	//lock.Lock()
	//defer lock.Unlock()
	err, obj := service.GetProBillById(reqParam)
	if err != nil {
		return err
	}
	if obj.Stage != 3 && obj.Stage != 4 {
		return errors.New("回传查询单证失败,状态有误")
	}
	if uploadPath.UploadRename == "" {
		return errors.New("回传查询单证失败,路径获取token的配置【UploadRename】为空")
	}
	//time.Sleep(time.Second * 10)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	dir = strings.Replace(dir, "bin", "", 1)
	jsonFile := dir + global.GConfig.LocalUpload.FilePath + obj.ProCode + "/upload_xml/" +
		fmt.Sprintf("%v/%v/%v/%v.json",
			obj.CreatedAt.Year(), int(obj.CreatedAt.Month()),
			obj.CreatedAt.Day(),
			obj.BillNum)
	global.GLog.Info(jsonFile)
	cmd := `curl -H "Content-Type:application/json" -X POST  -d @%v -x socks5://127.0.0.1:7070 http://110.20.14.175:9002/clbps/eomsmicroapp/services/rest/claimEomsService/receiveBpo`
	//secret, err := service1.GetGuoShouCacheSecret()
	//if err != nil {
	//	global.GLog.Error("获取缓存的secret失败", zap.Error(err))
	//	secret, err = GetToken(uploadPath.UploadRename)
	//	if err != nil {
	//		global.GLog.Error("获取secret失败", zap.Error(err))
	//		return err
	//	}
	//}
	//sh := "curl -T %v%v.xml 'ftp://192.168.202.3/lipei2.0/%v.temp' -u 'myftp:myftp' -s -x socks5://127.0.0.1:8090"
	// todo 调用回传
	//fmt.Println(secret)
	//if uploadPath.Upload == "" {
	//	return errors.New("回传查询单证失败,回传路径配置【Upload】为空")
	//}
	//curl -H 'Content-Type:application/json' -H 'Accept-Charset:utf-8' -H 'contentType:utf-8'
	//  -H 'Authorization:Bearer %v' -X POST
	//   -d @%v https://110.18.69.36:443/api/ivmsout/bpo/receivercptinfo
	cmd = fmt.Sprintf(cmd, jsonFile)
	err, s, s2 := utils.ShellOut(cmd)
	global.GLog.Error("回传命令", zap.Any("cmd", cmd))
	global.GLog.Error("上传错误err", zap.Error(err))
	global.GLog.Info("上传回显", zap.Any("std out", s))
	global.GLog.Error("上传回显", zap.Any("std err", s2))
	if err != nil {
		return err
	}

	uploadResp := Resp{}
	err = json.Unmarshal([]byte(s), &uploadResp)
	if err != nil {
		return err
	}
	if uploadResp.ResultCode != "0000" {
		return errors.New(uploadResp.ResultCode + uploadResp.ResultMsg)
	}

	err = service1.UpdateBill(reqParam, time.Now())

	return err
	//lock.Unlock()
}

// Resp 回传的结构体
type Resp struct {
	ResultCode string `json:"resultCode"`
	ResultMsg  string `json:"resultMsg"`
}

// TokenResp 获取token返回值的结构体
type TokenResp struct {
	Code   string `json:"code"`
	Result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	} `json:"result"`
}

type UploadResp struct {
	Code   string `json:"code"`
	Msg    string `json:"msg"`
	Result Resp   `json:"result"`
}

// GetToken 获取token
func GetToken(sh string) (string, error) {
	//sh := "curl -k --location --request POST 'https://110.18.13.74:8443/oauth2/token' --header 'Accept-Charset: utf-8' --header 'contentType: utf-8' --header 'Content-Type: application/json' --data-raw '{\"app_key\": \"hmkxg9xauivqmrpjky\",\"app_secret\": \"824253c9ec71444ab547d948600c44be\"}'"
	err, s, s2 := utils.ShellOut(sh)
	if err != nil {
		global.GLog.Error("上传错误err", zap.Error(err))
		return "", err
	}
	global.GLog.Error("上传回显", zap.Any("std err", s2))
	global.GLog.Info("上传回显", zap.Any("std out", s))
	tokenResp := TokenResp{}
	err = json.Unmarshal([]byte(s), &tokenResp)
	if err != nil {
		return "", err
	}
	resErrMap := map[string]string{
		"40001": "未正确传递AppKey及AppSecret",
		"40002": "AppKey或AppSecret不正确",
		"40302": "AppKey已禁用",
	}
	if errMsg, ok := resErrMap[tokenResp.Code]; ok {
		return "", errors.New(errMsg)
	}
	// fmt.Println("---------------AccessToken---------------", tokenResp.Result.AccessToken)
	// fmt.Println("---------------ExpiresIn---------------", tokenResp.Result.ExpiresIn)

	err = service1.SetGuoShouCacheSecret(tokenResp.Result.AccessToken, tokenResp.Result.ExpiresIn*1000*1000*1000)
	if err != nil {
		return "", err
	}
	return tokenResp.Result.AccessToken, nil
}
