package project

import (
	"fmt"
	"time"
)

func Process() {
	projectBills := Scan()
	var err error
	for ii, projectBill := range projectBills {
		projectBill.ProCode = "B0118"
		fmt.Println("projectBill", ii, projectBill)
		projectBill.DownloadAt = time.Now()
		projectBill.ScanAt = time.Now()
		projectBill.Status = 1

		err, projectBill = FetchBill(projectBill)
		if err != nil {
			fmt.Println("FetchBill err:", err)
			continue
		}

		err, projectBill = Decompress(projectBill)
		if err != nil {
			fmt.Println("Decompress err:", err)
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
			fmt.Println("Save err:", err)
			continue
		}

		err, projectBill = Clean(projectBill)
		if err != nil {
			fmt.Println("Save err:", err)
			continue
		}
	}

}
