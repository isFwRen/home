package project

import (
	"fmt"
	"server/module/download/service"
	"server/module/pro_manager/model"
)

func Save(projectBill model.ProjectBill) (error, model.ProjectBill) {
	/* 这是我的第一个简单的程序 */
	fmt.Println("projectBill:", projectBill)
	projectBill.DownloadPath = projectBill.DownloadPath + projectBill.BillName + "/"
	err, aaa := service.InsertBill(projectBill.ProCode, projectBill)
	if err != nil {
		fmt.Println("errerr:", err)
		return err, projectBill
	}

	fmt.Println("aaaaaaaaa:", aaa)

	return err, projectBill
}
