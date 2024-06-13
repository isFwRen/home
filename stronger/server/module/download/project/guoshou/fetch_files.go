/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/4/12 15:47
 */

package guoshou

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"server/global"
	model2 "server/module/download/model"
	"server/module/download/project"
	"server/module/pro_manager/model"
)

func FetchFiles(bill model.ProjectBill) (model.ProjectBill, error) {
	cache, ok := global.GProConf[bill.ProCode]
	if !ok {
		err := errors.New("没有该项目[" + bill.ProCode + "]配置")
		global.GLog.Error("", zap.Error(err))
		return bill, err
	}
	if cache.DownloadPaths.FetchFile == "" {
		err := errors.New("该项目[" + bill.ProCode + "]下载配置为空配置")
		global.GLog.Error("", zap.Error(err))
		return bill, err
	}
	var param model2.ImageFile
	err := json.Unmarshal([]byte(bill.OtherInfo), &param)
	if err != nil {
		return bill, err
	}
	for _, image := range param.CmsImageInfoList {
		fileName := image.Id + "." + image.FileType
		bill.Images = append(bill.Images, fileName)
		cmd := fmt.Sprintf(cache.DownloadPaths.FetchFile,
			bill.DownloadPath+fileName,
			param.Url+image.ImageUrl)
		global.GLog.Info("cmd", zap.Any("", cmd))
		//cmd := "curl -o %v http://1%v -x socks5://0.0.0.0:7070"
		err, stdout, stderr := project.ShellOut(cmd)
		global.GLog.Info("stdout", zap.Any("", stdout))
		global.GLog.Error("stderr", zap.Any("", stderr))
		if err != nil {
			global.GLog.Error("下载图片失败", zap.Any(image.Id, err))
			bill.Status = 4
			bill.WrongNote = "该案件影像异常" + image.Id
			//return bill, err
		}

	}

	//利用协程下载
	//gCount := len(param.CmsImageInfoList)
	//var wg sync.WaitGroup
	//for i := 0; i < gCount; i++ {
	//	j := int64(i)
	//	wg.Add(1)
	//	go func(k *int64) {
	//		defer func() {
	//			if err := recover(); err != nil {
	//				bill.Status = 4
	//				bill.WrongNote = "该案件影像异常"
	//				global.GLog.Error("下载图片失败", zap.Any("", err))
	//			}
	//			wg.Done()
	//		}()
	//		image := param.CmsImageInfoList[*k]
	//		fileName := image.Id + "." + image.ImageType
	//		bill.Images = append(bill.Images, fileName)
	//		cmd := fmt.Sprintf(cache.DownloadPaths.FetchFile,
	//			bill.DownloadPath+fileName,
	//			param.Url+image.ImageUrl)
	//		global.GLog.Info("cmd", zap.Any("", cmd))
	//		//cmd := "curl -o %v http://1%v -x socks5://0.0.0.0:7070"
	//		err, stdout, stderr := project.ShellOut(cmd)
	//		global.GLog.Info("stdout", zap.Any("", stdout))
	//		global.GLog.Error("stderr", zap.Any("", stderr))
	//		if err != nil {
	//			panic(err)
	//		}
	//	}(&j)
	//}
	//wg.Wait()
	return bill, nil
}
