/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023年08月07日10:33:53
 */

package B0122

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

var GlobalToken = "ILMP_ZHHLCLM" //区分外包商
var Tool = "/opt/jdk-20.0.2/bin/java -jar /opt/xincheng_des.jar"
var Key = "4Vj82gBO#&uFh2Sn" //加密秘钥
var KeyProd = "SBMLi2ViMGsOOe#&"

//var EncURL = "http://127.0.0.1:8889/v1/des/encrypt" //加密服务URL

type Resp struct {
	ResultCode int    `json:"resultCode"`
	ResultMsg  string `json:"resultMsg"`
	Data       struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
}

// BillUpload 回传单据
func BillUpload(reqParam model.ProCodeAndId, uploadPath model2.SysProDownloadPaths) error {
	//lock.Lock()
	//defer lock.Unlock()
	if global.GConfig.System.Env == "prod" {
		Key = KeyProd
	}
	if uploadPath.UploadRename == "" {
		return errors.New("加密url为空，UploadRename")
	}
	err, obj := service.GetProBillById(reqParam)
	if err != nil {
		return err
	}
	if obj.Stage != 3 && obj.Stage != 4 {
		return errors.New("回传查询单证失败,状态有误")
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	file := strings.Replace(dir, "bin", "", 1) +
		global.GConfig.LocalUpload.FilePath + obj.ProCode + "/upload_xml/" +
		fmt.Sprintf("%v/%v/%v/%v.json",
			obj.CreatedAt.Year(), int(obj.CreatedAt.Month()),
			obj.CreatedAt.Day(), obj.BillNum)

	uploadTime := time.Now()
	global.GLog.Info("timestamp", zap.Any("uploadTime", uploadTime))
	global.GLog.Info("file:::" + file)
	data, err := os.ReadFile(file)
	if err != nil {
		global.GLog.Error("Read JSON File", zap.Error(err))
		return err
	}
	var body map[string]interface{}
	err = json.Unmarshal(data, &body)
	if err != nil {
		global.GLog.Error("unmarshal json file error", zap.Error(err))
		return err
	}
	body["timestamp"] = uploadTime.Unix() * 1000
	strData, err := json.Marshal(body)
	if err != nil {
		global.GLog.Error("marshal json file error", zap.Error(err))
		return err
	}
	global.GLog.Info("s", zap.Any("strData", string(strData)))

	err = os.WriteFile(file, strData, 0666)
	if err != nil {
		global.GLog.Error("WriteFile file error", zap.Error(err))
		return err
	}

	//加密数据回传
	//cmd := fmt.Sprintf(`/opt/jdk-20.0.2/bin/java -jar /opt/xincheng_des.jar '%v' %v %v %v`, Key, GlobalToken, uploadTime.Unix()*1000, file)
	//cmd := fmt.Sprintf(`%v '%v' '%v' 0 `, Tool, file, Key)
	//curl --location --request POST 'http://127.0.0.1:13000/api/v1/xincheng/encrypt' --form 'file=@"/Users/mjl/jmeter.log"'
	cmd := fmt.Sprintf(`curl --location --request POST '%v/api/v1/xincheng/encrypt' --form 'file=@%v'`, uploadPath.UploadRename, file)
	global.GLog.Info("cmd", zap.Any("", cmd))
	err, stdout, stderr := utils.ShellOut(cmd)
	if err != nil {
		global.GLog.Error("err", zap.Error(err))
		return err
	}
	//if stderr != "" {
	//	global.GLog.Error("stderr", zap.Any("", stderr))
	//	return errors.New(stderr)
	//}

	global.GLog.Info("xin cheng upload", zap.Any("加密后数据", stdout))
	desResp := DesResp{}
	err = json.Unmarshal([]byte(stdout), &desResp)
	if err != nil {
		return err
	}

	global.GLog.Info("url:::" + uploadPath.Upload)

	err = os.WriteFile(file+".data", []byte(desResp.Data), 0666)
	if err != nil {
		global.GLog.Error("WriteFile file error", zap.Error(err))
		return err
	}

	//回传数据
	//stdout = "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"
	cmd = "curl --location -k --request POST '" + uploadPath.Upload + "' --header 'GLOBAL_TOKEN: " + GlobalToken + "' --header 'Content-Type: application/json;charset=UTF-8' -d @" + file + ".data"
	global.GLog.Info("cmd", zap.Any("curl", cmd))
	err, stdout, stderr = utils.ShellOut(cmd)
	if err != nil {
		global.GLog.Error("curl post", zap.Error(err))
		return err
	}
	global.GLog.Info("curl post", zap.Any("stderr", stderr))
	global.GLog.Info("curl post", zap.Any("stdout", stdout))

	//解密回传返回数据
	//cmd = fmt.Sprintf(`%v '%v' '%v' 1 `, Tool, stdout, Key)
	cmd = fmt.Sprintf(`curl --location --request POST '%v/api/v1/xincheng/decrypt' --form 'encData=%v'`, uploadPath.UploadRename, stdout)
	global.GLog.Info("cmd", zap.Any("解密", cmd))
	err, stdout, stderr = utils.ShellOut(cmd)
	if err != nil {
		global.GLog.Error("curl post", zap.Error(err))
		return err
	}
	global.GLog.Error("解密", zap.Any("stderr", stderr))
	global.GLog.Error("解密", zap.Any("stdout", stdout))

	resp := DesResp{}
	err = json.Unmarshal([]byte(stdout), &resp)
	if err != nil {
		return err
	}

	r := Resp{}
	err = json.Unmarshal([]byte(resp.Data), &r)
	if err != nil {
		global.GLog.Error("Unmarshal body err", zap.Error(err))
		return err
	}
	global.GLog.Info("Unmarshal body", zap.Any("Code", r.ResultCode))
	global.GLog.Info("Unmarshal body", zap.Any("Message", r.ResultMsg))

	if r.ResultCode == 0 && r.Data.Code == 0 {
		return service1.UpdateBill(reqParam, uploadTime)
	} else {
		return errors.New(r.ResultMsg + "," + r.Data.Message)
	}

}

type DesResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}
