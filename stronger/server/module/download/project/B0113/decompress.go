/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023年03月10日17:07:08
 */

package B0113

import (
	"fmt"
	"go.uber.org/zap"
	"regexp"
	"server/global"
	"server/module/download/project"
	"server/module/pro_manager/model"
	"strings"
)

func Decompress(projectBill model.ProjectBill) (error, model.ProjectBill) {
	zip := projectBill.DownloadPath + projectBill.Files[0]
	cmd := fmt.Sprintf("7za x %v -y -aoa -bb -o%v", zip, projectBill.DownloadPath)
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
		if line == "" {
			continue
		}
		projectBill.Images = append(projectBill.Images, line)
	}

	global.GLog.Info("images", zap.Any("", projectBill.Images))

	return err, projectBill
}
