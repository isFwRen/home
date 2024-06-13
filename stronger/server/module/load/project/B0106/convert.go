package B0106

import (
	"fmt"
	"path"
	"regexp"
	"server/global"
	"server/module/pro_manager/model"
	"server/utils"
	"sync"

	"go.uber.org/zap"
)

var limitChan = make(chan bool, 50)

func Convert(proCode string, projectBill model.ProjectBill) (error, model.ProjectBill) {
	pictures := projectBill.Images
	projectBill.Pictures = []string{}
	projectBill.ImagesType = make([]string, len(projectBill.Images))

	re := regexp.MustCompile(`\.[^\.]+$`)
	var wg sync.WaitGroup
	for ii, picture := range pictures {
		//解密
		//decrypt

		ext := path.Ext(picture)
		imgPng := re.ReplaceAllString(picture, ext)
		matched, _ := regexp.MatchString(`(?i)(tif|tiff|pdf)$`, picture)
		if matched {
			imgPng = re.ReplaceAllString(picture, ".png")
		}
		global.GLog.Info("convert:picture:ii", zap.Any(picture, ii))
		downPic := projectBill.DownloadPath + picture

		//获取图片类型 清单报销单等
		// url := "http://192.168.0.48:18301/classifier"
		// if global.GConfig.System.Env == "test" {
		// 	url = "http://192.168.202.17:18001/classifier"
		// }
		// getImageClassifierCmd := fmt.Sprintf("curl -X POST -F image=@%v %v", downPic, url)
		// global.GLog.Info("getImageClassifierCmd", zap.Any("", getImageClassifierCmd))
		// err, stdout, stderr := utils.ShellOut(getImageClassifierCmd)
		// if err != nil {
		// 	global.GLog.Error("classifier-cmd-err", zap.Error(err))
		// 	global.GLog.Error("classifier-cmd-stderr" + stderr)
		// 	global.GLog.Error("classifier-cmd-stdout::" + stdout)
		// }
		// classArr := make([][]interface{}, 0)
		// err = json.Unmarshal([]byte(stdout), &classArr)
		// if err != nil {
		// 	global.GLog.Error("classifier-json-Unmarshal-err", zap.Error(err))
		// }
		// projectBill.ImagesType[ii] = stdout
		// if len(classArr) > 0 && len(classArr[0]) > 0 {
		// 	projectBill.ImagesType[ii] = classArr[0][0].(string)
		// }

		cropApic := projectBill.DownloadPath + "A" + imgPng
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
		if !matched {
			projectBill.Pictures = append(projectBill.Pictures, picture)
			continue
		}
		cropPic := projectBill.DownloadPath + imgPng
		cmdConvert := fmt.Sprintf(`convert -resize 1600x1600 %s %s`, downPic, cropPic)
		matched, _ = regexp.MatchString(`(pdf)$`, picture)
		if matched {
			cmdConvert = fmt.Sprintf(`convert -density 300 -quality 100 %s %s`, downPic, cropPic)
		}
		global.GLog.Info("cmdConvert", zap.Any("1600*1600", cmdConvert))
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
		wg.Wait()
		//err, stdout, _ = utils.ShellOut(cmd)
		//global.GLog.Info("convert", zap.Any("resize", stdout))
		//if err != nil {
		//	global.GLog.Error("convert-resize", zap.Error(err))
		//}

		//加密
		//Encrypt

		projectBill.Pictures = append(projectBill.Pictures, imgPng)
	}
	//wg.Wait()
	return nil, projectBill
}