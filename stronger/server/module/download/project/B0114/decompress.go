/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/12/8 14:29
 */

package B0114

import (
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"regexp"
	"server/global"
	"server/module/download/project"
	"server/module/pro_manager/model"
	"server/utils"
	"strings"
)

func Decompress(projectBill model.ProjectBill) (error, model.ProjectBill) {
	zip := global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath + projectBill.Files[1]
	file := global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath
	cmd := fmt.Sprintf(`7za x "%s" -y -aoa -bb -o"%s"`, zip, file)
	global.GLog.Info("", zap.Any("cmd", cmd))
	err, stdout, _ := project.ShellOut(cmd)
	global.GLog.Info("", zap.Any("stdout", stdout))
	if err != nil {
		global.GLog.Error("7za", zap.Error(err))
		return err, projectBill
	}
	reg := regexp.MustCompile(`[\r\n]+`)
	lines := reg.Split(stdout, -1)
	//global.GLog.Info("lines", zap.Any("", lines))
	for _, line := range lines {
		if !strings.HasPrefix(line, "- ") {
			continue
		}
		//global.GLog.Info("line", zap.Any("", line))
		line = strings.ReplaceAll(line, "- ", "")
		//global.GLog.Info("line1111", zap.Any("", line))
		lineArr := strings.Split(line, "/")
		//global.GLog.Info("lineArr", zap.Any("", lineArr))
		if lineArr == nil || len(lineArr) == 0 || lineArr[0] == "" {
			continue
		}
		global.GLog.Info("lineArr", zap.Any("", lineArr))
		projectBill.Images = append(projectBill.Images, lineArr[0])
	}

	global.GLog.Info("images", zap.Any("", projectBill.Images))
	filename := global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath + "/" + projectBill.Files[0]
	contentByte, err := ioutil.ReadFile(filename)
	if err != nil {
		global.GLog.Error("read xml", zap.Error(err))
		return err, projectBill
	}
	data := string(contentByte)
	//global.GLog.Info("", zap.Any("xml data", data))

	//accidentAreaOption := utils.GetNodeValue(data, "accidentAreaOption")
	claimInCode := utils.GetNodeValue(data, "claimInCode")
	//global.GLog.Info("accidentAreaOption", zap.Any("", accidentAreaOption))
	global.GLog.Info("claimInCode", zap.Any("", claimInCode))
	//projectBill.BatchNum = accidentAreaOption
	projectBill.BatchNum = strings.Replace(projectBill.Files[1], ".zip", "", -1)
	projectBill.BillNum = claimInCode
	projectBill.OtherInfo = data

	policyBranchCodeArr := utils.GetNodeData(data, "policyBranchCode")
	if len(policyBranchCodeArr) > 0 {
		projectBill.Agency = policyBranchCodeArr[0]
	}

	return err, projectBill
}
