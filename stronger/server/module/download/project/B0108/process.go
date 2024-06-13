package project

import (
	"go.uber.org/zap"
	"server/global"
	"strings"
	"time"
)

func Process() {
	if strings.Index(global.GConfig.System.DownloadPath, "Inquiry") != -1 {
		//客户消息，通知函2
		customerNotices := CustomerNoticeScan()
		for _, customerNotice := range customerNotices {
			err := CustomerNoticeFetchBill(customerNotice)
			if err != nil {
				global.GLog.Error("Process-CustomerNoticeFetchBill", zap.Error(err))
				continue
			}
		}

	} else if strings.Index(global.GConfig.System.DownloadPath, "claim_error") != -1 {
		//通知函
		projectBills := NoticeScan()
		for _, projectBill := range projectBills {
			err := NoticeFetchBill(global.GConfig.System.ProCode, projectBill)
			if err != nil {
				global.GLog.Error("Process-NoticeFetchBill", zap.Error(err))
				continue
			}
		}

	} else {
		//正常下单
		projectBills := Scan()
		var err error
		for _, projectBill := range projectBills {
			projectBill.ProCode = global.GConfig.System.ProCode
			global.GLog.Info("process", zap.Any("projectBill", projectBill))
			projectBill.DownloadAt = time.Now()
			projectBill.ScanAt = time.Now()
			projectBill.Status = 1

			err, projectBill = FetchBill(projectBill)
			if err != nil {
				global.GLog.Error("Process-FetchBill", zap.Error(err))
				continue
			}

			err, projectBill = Decompress(projectBill)
			if err != nil {
				global.GLog.Error("Process-Decompress", zap.Error(err))
				continue
			}

			//加密图片
			//err, projectBill = Encrypt(projectBill)
			//if err != nil {
			//	fmt.Println("Encrypt err:", err)
			//	continue
			//}

			err, projectBill = Save(projectBill)
			if err != nil {
				global.GLog.Error("Process-Save", zap.Error(err))
				continue
			}

			err, projectBill = Clean(projectBill)
			if err != nil {
				global.GLog.Error("Process-Clean", zap.Error(err))
				continue
			}
		}
	}
}
