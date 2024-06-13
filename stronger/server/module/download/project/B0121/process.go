package B0121

import (
	"fmt"
	"go.uber.org/zap"
	"server/global"
	"time"
)

func Process() {
	projectBills := Scan()
	var err error
	for _, projectBill := range projectBills {
		projectBill.ProCode = global.GConfig.System.ProCode
		projectBill.DownloadAt = time.Now()
		projectBill.ScanAt = time.Now()
		projectBill.Status = 1

		err, projectBill = FetchBill(projectBill)
		if err != nil {
			global.GLog.Error("FetchBill", zap.Error(err))
			continue
		}

		err, projectBill = Decompress(projectBill)
		if err != nil {
			global.GLog.Error("Decompress", zap.Error(err))
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
			global.GLog.Error("save", zap.Error(err))
			continue
		}

		err, projectBill = Clean(projectBill)
		if err != nil {
			fmt.Println("Save err:", err)
			continue
		}
	}

}
