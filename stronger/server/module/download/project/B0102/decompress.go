package B0102

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"server/global"
	"server/module/download/project"
	"server/module/pro_manager/model"
	"server/utils"
	"strings"

	"go.uber.org/zap"
)

func Decompress(projectBill model.ProjectBill) (error, model.ProjectBill) {
	fileName := projectBill.BillName + ".zip"
	zip := projectBill.DownloadPath + fileName
	//cmd := fmt.Sprintf("7za x %v -y -aoa -bb -o%v", zip, projectBill.DownloadPath)
	cmd := fmt.Sprintf("unzip -O GBK %v -d %v", zip, projectBill.DownloadPath)
	global.GLog.Info("", zap.Any("cmd", cmd))
	err, stdout, _ := project.ShellOut(cmd)
	global.GLog.Info("", zap.Any("stdout", stdout))
	if err != nil {
		global.GLog.Error("7za", zap.Error(err))
		return err, projectBill
	}
	reg := regexp.MustCompile(`[\r\n]+`)
	lines := reg.Split(stdout, -1)
	global.GLog.Info("lines", zap.Any("", lines))
	for _, line := range lines {
		if strings.Index(line, "inflating: ") == -1 {
			continue
		}
		global.GLog.Info("line", zap.Any("", line))
		line = strings.TrimSpace(strings.ReplaceAll(line, "inflating: ", ""))
		global.GLog.Info("line1111", zap.Any("", line))
		if line == "" {
			continue
		}
		if strings.Index(line, ".xml") != -1 {
			continue
		}
		arr := strings.Split(line, "/")
		projectBill.Images = append(projectBill.Images, arr[len(arr)-1])
	}

	global.GLog.Info("images", zap.Any("", projectBill.Images))

	filename := projectBill.DownloadPath + "/" + utils.Substr(projectBill.BillName, 0, 20) + ".xml"
	contentByte, err := ioutil.ReadFile(filename)
	if err != nil {
		global.GLog.Error("read xml", zap.Error(err))
		return err, projectBill
	}
	data := string(contentByte)
	//global.GLog.Info("", zap.Any("xml data", data))

	// accidentAreaOption := utils.GetNodeValue(data, "accidentAreaOption")
	// claimInCode := utils.GetNodeValue(data, "claimInCode")
	// global.GLog.Info("accidentAreaOption", zap.Any("", accidentAreaOption))
	// global.GLog.Info("claimInCode", zap.Any("", claimInCode))
	projectBill.BatchNum = projectBill.BillName
	projectBill.BillNum = utils.Substr(projectBill.BillName, 0, 20)
	projectBill.OtherInfo = data
	case_type := utils.GetNodeValue(data, "case_type")
	if case_type == "W" {
		projectBill.SaleChannel = "微理赔"
	} else if case_type == "M" {
		projectBill.SaleChannel = "移动理赔"
	} else if case_type == "Z" {
		projectBill.SaleChannel = "纸质"
	}
	projectBill.Agency = utils.GetNodeValue(data, "branch_name")
	// c, ok := global.GProConf["B0102"].ConstTable["B0102_信诚理赔_机构号代码表"]
	// if ok {
	// 	for _, arr := range c {
	// 		if branch_no == strings.TrimSpace(arr[1]) {
	// 			projectBill.Agency = strings.TrimSpace(arr[0])
	// 		}
	// 	}
	// }

	return err, projectBill
}
