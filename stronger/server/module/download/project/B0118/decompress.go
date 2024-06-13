package project

import (
	"fmt"
	"io/ioutil"
	"server/global"
	"server/module/download/project"
	"server/module/pro_manager/model"
	"server/utils"
	"strings"
)

func Decompress(projectBill model.ProjectBill) (error, model.ProjectBill) {
	/* 这是我的第一个简单的程序 */
	zip := global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath + projectBill.BillName + ".zip"
	file := global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath
	cmd := fmt.Sprintf(`7za x "%s" -y -aoa -o"%s"`, zip, file)
	fmt.Printf("cmdcmdcmd", cmd)
	err, stdout, _ := project.ShellOut(cmd)
	if err != nil {
		fmt.Printf("7za", err, stdout)
		return err, projectBill
	}
	filename := global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath + projectBill.BillName + "/" + projectBill.BillName + ".xml"
	// filename := `D:\aaa.txt`
	contentByte, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("ReadFile err", err)
		return err, projectBill
	}
	data := string(contentByte)

	// fmt.Println("ioutil.ReadFile", err)
	fmt.Println("ioutil.ReadFile data", data)

	billNum := GetNodeValue(data, "claimNum")
	priority := GetNodeValue(data, "priority")
	images := GetNodeValue(data, "images")
	desc := GetNodeValue(data, "desc")
	saleChannel := GetNodeValue(data, "specialAgree")
	fmt.Println("billNum", billNum)
	fmt.Println("priority", priority)
	fmt.Println("images", images)
	fmt.Println("desc", desc)
	//fmt.Println("saleChannel", saleChannel)
	projectBill.Stage = 1
	//1是紧急 2是优先
	if priority == "1" {
		projectBill.StickLevel = 1
	}
	if priority == "2" {
		projectBill.StickLevel = 2
	}
	if priority == "0" {
		projectBill.InsuranceType = "一般"
	}
	if priority == "5" {
		projectBill.InsuranceType = "补录"
	}

	//CSB0118RC0318000
	projectBill.SaleChannel = saleChannel
	projectBill.BillNum = billNum
	projectBill.Agency = priority
	projectBill.Images = strings.Split(images, ";")
	fmt.Println("ImagesImages", projectBill.Images)

	err, agencies := utils.GetRedisAgency(projectBill.ProCode)
	if err != nil {
		return err, model.ProjectBill{}
	}
	//包含
	if hasItem(agencies, projectBill.Agency) {
		projectBill.StickLevel = 1
	}

	// projectBill.Code = billNum
	// projectBill.ProjectCode = "BC0001"
	// projectBill.Path = "/download/"

	return err, projectBill
}

func GetNodeValue(xmlValue string, nodeName string) string {
	sValue := ""
	beginNode := strings.Index(xmlValue, "<"+nodeName+">") + len(nodeName) + 2
	endNode := strings.Index(xmlValue, "</"+nodeName+">")
	fmt.Println("endNodeendNodeendNode", nodeName, beginNode, endNode)
	if beginNode != -1 && endNode != -1 {
		sValue = xmlValue[beginNode:endNode]
	}
	return sValue
}

func hasItem(arr []string, item string) bool {
	for _, i2 := range arr {
		if i2 == item {
			return true
		}
	}
	return false
}
