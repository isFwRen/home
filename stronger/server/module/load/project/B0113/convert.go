package B0113

import (
	"fmt"
	"path"
	"regexp"
	"server/global"
	"server/module/pro_manager/model"
	"server/utils"
	"strings"
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
		matched, _ := regexp.MatchString(`(tif|TIF|TIFF|pdf)$`, picture)
		if matched {
			imgPng = re.ReplaceAllString(picture, ".png")
		}
		global.GLog.Info("convert:picture:ii", zap.Any(picture, ii))
		downPic := projectBill.DownloadPath + picture
		cropApic := projectBill.DownloadPath + strings.Replace(imgPng, "/", "/A", -1)
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
			re := regexp.MustCompile(`\.[^\.]+$`)
			pname := re.ReplaceAllString(cropPic, "*.png")
			cmdConvert = fmt.Sprintf(`convert -resize 1600x1600 %s %s;ls %s`, downPic, cropPic, pname)
			// cmdConvert = fmt.Sprintf(`convert -density 300 -quality 100 %s %s`, downPic, cropPic)
		}
		global.GLog.Info("cmdConvert", zap.Any("1600*1600", cmdConvert))
		// wg.Add(1)
		// limitChan <- true
		// go func() {
		// 	defer wg.Done()
		// 	defer func() {
		// 		if err := recover(); err != nil {
		// 			global.GLog.Error("转换异常", zap.Any("", err))
		// 		}
		// 	}()
		// 	fmt.Println("456456", cmdConvert)
		// 	err, stdout, stderr := utils.ShellOut(cmdConvert)
		// 	if err != nil {
		// 		global.GLog.Error("load-cmd-err", zap.Error(err))
		// 		global.GLog.Error("load-cmd-stderr" + stderr)
		// 		global.GLog.Error("load-cmd-stdout::" + stdout)
		// 	}
		// 	<-limitChan
		// }()
		// wg.Wait()
		//err, stdout, _ = utils.ShellOut(cmd)
		//global.GLog.Info("convert", zap.Any("resize", stdout))
		//if err != nil {
		//	global.GLog.Error("convert-resize", zap.Error(err))
		//}

		err, stdout, stderr := utils.ShellOut(cmdConvert)
		fmt.Println("---------cmdConvert-------------", cmdConvert, matched)
		if err != nil {
			global.GLog.Error("load-cmd-err", zap.Error(err))
			global.GLog.Error("load-cmd-stderr" + stderr)
			global.GLog.Error("load-cmd-stdout::" + stdout)
			return err, projectBill
		}

		//加密
		//Encrypt
		if matched {
			pdfs := RegMatchAll(stdout, `.+`)
			fmt.Println("pdfs", pdfs)
			for _, pdf := range pdfs {
				pdf = strings.Replace(pdf, projectBill.DownloadPath, "", 1)
				fmt.Println("pdf", pdf)
				projectBill.Pictures = append(projectBill.Pictures, pdf)
			}
		} else {
			projectBill.Pictures = append(projectBill.Pictures, imgPng)
		}

	}
	//wg.Wait()
	return nil, projectBill
}

func RegMatchAll(value string, query string) []string {
	reg := regexp.MustCompile(query)
	return reg.FindAllString(value, -1)
}