/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/4/13 11:32
 */

package guoshou

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"server/global"
	"server/module/download/const_data"
	"server/module/download/model"
	"server/module/download/project"
	"server/module/download/service"
	"server/utils"
	"time"
)

func FetchCatalogue(catalog model.Catalog) error {
	fileName := catalog.BranchNo + "_" + catalog.IsDeleted + "_" + catalog.CatalogCode + "-" + time.Now().Format("20060102030405") + ".xlsx"
	path := global.GConfig.LocalUpload.FilePath + global.PathUpdateConstCatalogue + time.Now().Format("20060102")
	isExist, _ := utils.PathExists(path)
	if !isExist {
		_ = utils.CreateDir(path)
	}
	cmd := fmt.Sprintf("curl -o %v/%v http://%v -x socks5://0.0.0.0:7070",
		path,
		fileName,
		catalog.Url+catalog.ResourceUrl)

	//存数据库log
	proCode := const_data.Num2code[catalog.BranchNo[:2]]
	if proCode == "" {
		return errors.New("未知案件机构:" + catalog.BranchNo)
	}
	byteData, err := json.Marshal(catalog)
	log := model.UpdateConstLog{
		Type:      2,
		Name:      proCode + "_" + global.GProConf[proCode].Name + "_医疗目录" + catalog.CatalogCode,
		ProCode:   proCode,
		RemoteUrl: catalog.Url + catalog.ResourceUrl,
		IsDeleted: model.IsDeletedMap[catalog.IsDeleted],
		LocalUrl:  path + "/" + fileName,
		OtherInfo: string(byteData),
	}
	err = service.InsertUpdateConstLog(log)
	if err != nil {
		global.GLog.Error("err", zap.Any("", err))
		return err
	}

	//下载文件
	global.GLog.Info("cmd", zap.Any("", cmd))
	//cmd := "curl -o %v http://%v -x socks5://0.0.0.0:7070"
	err, stdout, stderr := project.ShellOut(cmd)
	if err != nil {
		global.GLog.Info("stdout", zap.Any("", stdout))
		global.GLog.Error("stderr", zap.Any("", stderr))
		global.GLog.Error("err", zap.Any("", err))
	}
	global.GLog.Info("推送成功", zap.Any("catalogue", fileName))
	return err
}
