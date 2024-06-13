/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/4/6 10:21
 */

package B0118

import (
	"encoding/base64"
	"go.uber.org/zap"
	"io/ioutil"
	"path"
	"server/global"
	"server/module/pro_manager/model"
	"server/utils"
	"strings"
	"sync"
)

// Encrypt 加密
func Encrypt(projectBill model.ProjectBill) (error, model.ProjectBill) {

	//秘钥
	//crypto := utils.CreateRandomString(16)
	//projectBill.Crypto = crypto
	var wg sync.WaitGroup

	//原图
	wg.Add(1)
	go func() {
		defer wg.Done()
		encImages(projectBill, false)
	}()

	//缩列图
	wg.Add(1)
	go func() {
		defer wg.Done()
		encImages(projectBill, true)
	}()

	wg.Wait()

	return nil, projectBill
}

func encImages(projectBill model.ProjectBill, isThumbnail bool) {
	images := projectBill.Pictures
	prefix := ""
	if isThumbnail {
		images = projectBill.Images
		prefix = "A"
	}
	for _, image := range images {
		image = prefix + image
		err := encImage(projectBill, image, image)
		global.GLog.Error(err.Error())
	}
}

func encImage(projectBill model.ProjectBill, crypto string, image string) error {
	suffix := strings.Replace(path.Ext(image), ".", "", -1)
	filename := global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath + projectBill.BillName + "/" + image
	contentByte, err := ioutil.ReadFile(filename)
	if err != nil {
		global.GLog.Error("read file err", zap.Error(err))
		return err
	}
	//文件转base64
	sEnc := base64.StdEncoding.EncodeToString(contentByte)
	completeBase64Img := "data:image/" + suffix + ";base64," + sEnc
	//加密
	encrypt, err := utils.AesEncrypt([]byte(completeBase64Img), []byte(crypto))
	if err != nil {
		global.GLog.Error("encrypt file err", zap.Error(err))
		return err
	}
	//加密数据写回文件
	err = ioutil.WriteFile(filename, encrypt, 0666)
	if err != nil {
		global.GLog.Error("write file err", zap.Error(err))
		return err
	}
	return nil
}

func decrypt(projectBill model.ProjectBill, block string, image string) error {
	//读取加密文件
	//解密
	//base64转文件
	return nil
}
