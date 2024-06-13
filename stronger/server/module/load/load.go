package load

import (
	"fmt"
	"server/module/load/project/B0102"
	"server/module/load/project/B0103"
	"server/module/load/project/B0106"
	B0108 "server/module/load/project/B0108"
	"server/module/load/project/B0110"
	"server/module/load/project/B0113"
	"server/module/load/project/B0114"
	"server/module/load/project/B0116"
	B0118 "server/module/load/project/B0118"
	"server/module/load/project/B0121"
	"server/module/load/project/B0122"
	"server/module/pro_manager/model"
	// "server/module/load/"
)

func ProLoadFunc(proCode string, projectBill model.ProjectBill) error {
	err := error(nil)
	switch proCode {
	case "B0118":
		fmt.Println("B0118")
		err, projectBill = B0118.Convert(proCode, projectBill)
		if err != nil {
			return err
		}
		err, projectBill = B0118.Crop(proCode, projectBill)
	case "B0108":
		fmt.Println("B0108")
		err, projectBill = B0108.Convert(proCode, projectBill)
		err, projectBill = B0108.Crop(proCode, projectBill)
	case "B0114":
		fmt.Println("B0114")
		err, projectBill = B0114.Convert(proCode, projectBill)
		err, projectBill = B0114.Crop(proCode, projectBill)
	case "B0113":
		fmt.Println("B0113")
		err, projectBill = B0113.Convert(proCode, projectBill)
		err, projectBill = B0113.Crop(proCode, projectBill)
	case "B0121":
		fmt.Println("B0121")
		err, projectBill = B0121.Convert(proCode, projectBill)
		err, projectBill = B0121.Crop(proCode, projectBill)
	case "B0106":
		fmt.Println("B0106")
		err, projectBill = B0106.Convert(proCode, projectBill)
		err, projectBill = B0106.Crop(proCode, projectBill)
	case "B0103":
		fmt.Println("B0103")
		err, projectBill = B0103.Convert(proCode, projectBill)
		err, projectBill = B0103.Crop(proCode, projectBill)
	case "B0110":
		fmt.Println("B0110")
		err, projectBill = B0110.Convert(proCode, projectBill)
		err, projectBill = B0110.Crop(proCode, projectBill)
	case "B0122":
		fmt.Println("B0122")
		err, projectBill = B0122.Convert(proCode, projectBill)
		err, projectBill = B0122.Crop(proCode, projectBill)
	case "B0116":
		fmt.Println("B0116")
		err, projectBill = B0116.Convert(proCode, projectBill)
		err, projectBill = B0116.Crop(proCode, projectBill)
	case "B0102":
		fmt.Println("B0102")
		err, projectBill = B0102.Convert(proCode, projectBill)
		err, projectBill = B0102.Crop(proCode, projectBill)
	default:
		fmt.Println("不存在的加载进程项目编号:", proCode)
	}
	return err
}
