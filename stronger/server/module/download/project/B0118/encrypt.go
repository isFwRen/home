/**
 * @Author: xingqiyi
 * @Description: 加密
 * @Date: 2022/4/6 10:21
 */

package project

import (
	"encoding/base64"
	"go.uber.org/zap"
	"io/ioutil"
	"path"
	"server/global"
	"server/module/pro_manager/model"
	"server/utils"
	"strings"
)

// Encrypt 加密
func Encrypt(projectBill model.ProjectBill) (error, model.ProjectBill) {
	for _, image := range projectBill.Images {
		suffix := strings.Replace(path.Ext(image), ".", "", -1)
		filename := global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath + projectBill.BillName + "/" + image
		contentByte, err := ioutil.ReadFile(filename)
		if err != nil {
			global.GLog.Error("read file err", zap.Error(err))
			return err, projectBill
		}
		//文件转base64
		sEnc := base64.StdEncoding.EncodeToString(contentByte)
		completeBase64Img := "data:image/" + suffix + ";base64," + sEnc
		//加密
		//crypto := utils.CreateRandomString(16)
		//projectBill.Crypto = crypto
		encrypt, err := utils.AesEncrypt([]byte(completeBase64Img), []byte(filename))
		if err != nil {
			global.GLog.Error("encrypt file err", zap.Error(err))
			return err, model.ProjectBill{}
		}
		//加密数据写回文件
		err = ioutil.WriteFile(filename, encrypt, 0666)
		if err != nil {
			global.GLog.Error("write file err", zap.Error(err))
			return err, model.ProjectBill{}
		}
	}
	return nil, projectBill
}
