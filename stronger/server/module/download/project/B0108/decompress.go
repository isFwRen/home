package project

import (
	"fmt"
	"go.uber.org/zap"
	"regexp"
	"server/global"
	"server/module/download/project"
	"server/module/pro_manager/model"
	"server/utils"
	"strings"
)

func Decompress(projectBill model.ProjectBill) (error, model.ProjectBill) {
	zip := global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath + projectBill.BillName + ".zip"
	file := global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath
	cmd := fmt.Sprintf(`7za x "%s" -y -aoa -bb -o"%s"`, zip, file)
	global.GLog.Info("Decompress", zap.Any("cmd", cmd))
	err, stdout, _ := project.ShellOut(cmd)
	global.GLog.Info("Decompress-stdout", zap.Any("", stdout))
	if err != nil {
		global.GLog.Error("Decompress-err", zap.Error(err))
		return err, projectBill
	}
	reg := regexp.MustCompile(`[\r\n]+`)
	lines := reg.Split(stdout, -1)
	for _, line := range lines {
		if !strings.HasPrefix(line, "- "+projectBill.BillName) {
			continue
		}
		line = strings.ReplaceAll(line, "- ", "")
		lineArr := strings.Split(line, "/")
		if lineArr == nil {
			continue
		}
		if projectBill.SaleChannel == "理赔" {
			if len(lineArr) < 3 || lineArr[2] == "" {
				continue
			}
			projectBill.Images = append(projectBill.Images, lineArr[2])
			projectBill.BillNum = lineArr[1]
		}
		if projectBill.SaleChannel == "秒赔" {
			if len(lineArr) < 2 || lineArr[1] == "" || strings.HasSuffix(lineArr[1], ".dat") {
				continue
			}
			projectBill.Images = append(projectBill.Images, lineArr[1])
			projectBill.BillNum = strings.Split(lineArr[1], "_")[0]
		}
	}
	projectBill.BillName = projectBill.BillName + projectBill.BillNum
	global.GLog.Info("Images", zap.Any("projectBill.Images", projectBill.Images))

	err, agencies := utils.GetRedisAgency(projectBill.ProCode)
	if err != nil {
		return err, model.ProjectBill{}
	}
	//包含
	if hasItem(agencies, projectBill.Agency) {
		projectBill.StickLevel = 1
	}
	return err, projectBill
}

func hasItem(arr []string, item string) bool {
	for _, i2 := range arr {
		if i2 == item {
			return true
		}
	}
	return false
}
