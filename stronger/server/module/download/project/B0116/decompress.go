package B0116

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"server/global"
	"server/module/download/project"
	"server/module/pro_manager/model"
	"server/utils"
	"strconv"

	"go.uber.org/zap"
)

func Decompress(projectBill model.ProjectBill) (error, model.ProjectBill) {
	zip := global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath + projectBill.BatchNum + ".zip"
	file := global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath
	cmd := fmt.Sprintf(`7za x "%s" -y -aoa -bb -o"%s"`, zip, file)
	global.GLog.Info("", zap.Any("cmd", cmd))
	err, stdout, _ := project.ShellOut(cmd)
	// global.GLog.Info("", zap.Any("stdout", stdout))
	if err != nil {
		global.GLog.Error("7za", zap.Error(err))
		return err, projectBill
	}
	cmd = fmt.Sprintf(`ls %s`, file+projectBill.BatchNum+"/")
	global.GLog.Info("", zap.Any("cmd", cmd))
	err, stdout, _ = project.ShellOut(cmd)
	// fmt.Println("------cmd------", err, cmd, stdout)
	reg := regexp.MustCompile(`[\r\n]+`)
	lines := reg.Split(stdout, -1)
	//global.GLog.Info("lines", zap.Any("", lines))
	for _, line := range lines {
		if utils.RegIsMatch(`xml$`, line) || line == "" {
			continue
		}
		// fmt.Println("------line------", line)
		projectBill.Files = append(projectBill.Files, line)
		//global.GLog.Info("line", zap.Any("", line))
		// line = strings.ReplaceAll(line, "- ", "")
		// //global.GLog.Info("line1111", zap.Any("", line))
		// lineArr := strings.Split(line, "/")
		// //global.GLog.Info("lineArr", zap.Any("", lineArr))
		// if lineArr == nil || len(lineArr) == 0 || lineArr[0] == "" {
		// 	continue
		// }
		// global.GLog.Info("lineArr", zap.Any("", lineArr))
		// projectBill.Images = append(projectBill.Images, lineArr[0])
	}

	global.GLog.Info("images", zap.Any("", projectBill.Images))
	filename := global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath + "/" + projectBill.BatchNum + "/" + projectBill.BatchNum + ".xml"
	contentByte, err := ioutil.ReadFile(filename)
	if err != nil {
		global.GLog.Error("read xml", zap.Error(err))
		return err, projectBill
	}
	data := string(contentByte)
	//global.GLog.Info("", zap.Any("xml data", data))

	//accidentAreaOption := utils.GetNodeValue(data, "accidentAreaOption")
	// claimInCode := utils.GetNodeValue(data, "claimInCode")
	//global.GLog.Info("accidentAreaOption", zap.Any("", accidentAreaOption))
	// global.GLog.Info("claimInCode", zap.Any("", claimInCode))
	//projectBill.BatchNum = accidentAreaOption
	// projectBill.BatchNum = strings.Replace(projectBill.Files[1], ".zip", "", -1)
	// projectBill.BillNum = claimInCode
	projectBill.OtherInfo = data

	countNum := utils.GetNodeValue(data, "countNum")
	fmt.Println("------countNum------", countNum, len(projectBill.Files))
	rgtBpoClass := utils.GetNodeValue(data, "rgtBpoClass")
	projectBill.InsuranceType = rgtBpoClass
	num, _ := strconv.Atoi(countNum)
	if num != len(projectBill.Files) {
		return errors.New("数量异常，请检查；"), projectBill
	}
	// if len(policyBranchCodeArr) > 0 {
	// 	projectBill.Agency = policyBranchCodeArr[0]
	// }

	return err, projectBill
}
