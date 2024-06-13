package B0108

import (
	"fmt"
	"go.uber.org/zap"
	"regexp"
	"server/global"
	"server/module/pro_manager/model"
	"server/utils"
	"strconv"
	"sync"
)

var limitChan = make(chan bool, 50)

func Convert(proCode string, projectBill model.ProjectBill) (error, model.ProjectBill) {
	pictures := projectBill.Images
	projectBill.Pictures = []string{}
	re := regexp.MustCompile(`\.[^\.]+$`)
	var wg sync.WaitGroup
	for ii, picture := range pictures {
		//解密
		//decrypt

		imgPng := re.ReplaceAllString(picture, ".png")
		global.GLog.Info("convert:picture:ii", zap.Any(picture, ii))
		downPic := global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath + picture
		cropApic := global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath + "A" + imgPng
		cmd := fmt.Sprintf(`convert -resize x180 %s %s`, downPic, cropApic)
		global.GLog.Info("load", zap.Any("convert-cmd", cmd))
		wg.Add(1)
		limitChan <- true
		go func() {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					global.GLog.Error("转换缩列图异常", zap.Any("", err))
				}
			}()
			fmt.Println("123123", cmd)
			err, stdout, stderr := utils.ShellOut(cmd)
			if err != nil {
				global.GLog.Error("load-cmd-err", zap.Error(err))
				global.GLog.Error("load-cmd-stderr" + stderr)
				global.GLog.Error("load-cmd-stdout::" + stdout)
			}
			<-limitChan
		}()
		fmt.Println("limitChan0000", len(limitChan))
		//err, stdout, _ := utils.ShellOut(cmd)
		//if err != nil {
		//	global.GLog.Error("load-cmd-stdout", zap.Error(err))
		//}
		matched, _ := regexp.MatchString(`(tif|TIF|TIFF|pdf)$`, picture)
		if !matched {
			projectBill.Pictures = append(projectBill.Pictures, picture)
			continue
		}

		err, stdout, stderr := utils.ShellOut(`identify -format "%h" ` + downPic)
		if err != nil {
			global.GLog.Error("identify load-cmd-err", zap.Error(err))
			global.GLog.Error("identify load-cmd-stderr" + stderr)
		}
		global.GLog.Error("identify load-cmd-stdout::" + stdout)
		if h, _ := strconv.Atoi(stdout); h < 3000 {
			stdout = "1600"
		}
		cropPic := global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath + imgPng
		cmdConvert := fmt.Sprintf(`convert -resize 1600x%s %s %s`, stdout, downPic, cropPic)
		global.GLog.Info("cmdConvert", zap.Any("1600*"+stdout, cmdConvert))
		wg.Add(1)
		limitChan <- true
		go func() {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					global.GLog.Error("转换异常", zap.Any("", err))
				}
			}()
			fmt.Println("456456", cmdConvert)
			err, stdout, stderr := utils.ShellOut(cmdConvert)
			if err != nil {
				global.GLog.Error("load-cmd-err", zap.Error(err))
				global.GLog.Error("load-cmd-stderr" + stderr)
				global.GLog.Error("load-cmd-stdout::" + stdout)
			}
			<-limitChan
		}()
		//wg.Wait()
		//err, stdout, _ = utils.ShellOut(cmd)
		//global.GLog.Info("convert", zap.Any("resize", stdout))
		//if err != nil {
		//	global.GLog.Error("convert-resize", zap.Error(err))
		//}

		//加密
		//Encrypt

		projectBill.Pictures = append(projectBill.Pictures, imgPng)
	}
	wg.Wait()
	return nil, projectBill
}
