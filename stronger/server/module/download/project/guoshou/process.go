/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/4/12 16:18
 */

package guoshou

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"server/global"
	"server/module/download/const_data"
	"server/module/download/model"
	"server/module/download/service"
	model2 "server/module/pro_manager/model"
	"time"
)

func Process(param model.ImageFile) error {
	var bill model2.ProjectBill
	otherInfo, err := json.Marshal(param)
	if err != nil {
		return err
	}
	proCode := const_data.Num2code[param.BranchNo[:2]]
	if proCode == "" {
		return errors.New("未知案件机构:" + param.BranchNo)
	}
	bill.OtherInfo = string(otherInfo)
	bill.ProCode = proCode
	bill.DownloadAt = time.Now()
	bill.ScanAt = time.Now()
	bill.CreatedAt = time.Now()
	bill.Status = 1
	bill.Stage = 8 //先接收到系统  状态为待下载  定时下载文件
	bill.BillName = param.ClaimNo + "_" + param.ClaimTpaId
	bill.BillNum = param.ClaimNo + "_" + param.ClaimTpaId
	bill.BatchNum = param.ClaimTpaId
	bill.Agency = param.BranchNo
	bill.DownloadPath = fmt.Sprintf("%v%v/download/%v/%v/%v/%v/",
		global.GConfig.LocalUpload.FilePath,
		bill.ProCode,
		bill.CreatedAt.Format("2006"),
		bill.CreatedAt.Format("01"),
		bill.CreatedAt.Format("02"),
		bill.BillName)

	err, total := service.CountBillByBatchNum(proCode, bill.BatchNum)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		return err
	}
	if total > 0 {
		global.GLog.Error("已下载")
		return errors.New("已下载")
	}
	global.GLog.Info("DownloadPath", zap.Any("", bill.DownloadPath))
	err = os.MkdirAll(bill.DownloadPath, os.ModePerm)
	if err != nil {
		return err
	}

	//bill, err = FetchFiles(bill, param)
	//if err != nil {
	//	global.GLog.Error("FetchFiles", zap.Error(err))
	//	return err
	//}

	err, bill = service.InsertBill(proCode, bill)
	global.GLog.Info("file", zap.Any("推送成功_"+bill.ProCode, bill.BillNum))
	return err
}
