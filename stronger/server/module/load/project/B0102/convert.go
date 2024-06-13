package B0102

import (
	"fmt"
	"regexp"
	"server/global"
	"server/module/pro_manager/model"
	"server/utils"
	"strconv"
	"strings"
	"sync"

	"go.uber.org/zap"
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

		imgPng := re.ReplaceAllString(picture, "-%d.png")
		global.GLog.Info("convert:picture:ii", zap.Any(picture, ii))
		downPic := projectBill.DownloadPath + picture
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
		matched, _ := regexp.MatchString(`(tif|TIF|TIFF|pdf)$`, picture)
		if !matched {
			projectBill.Pictures = append(projectBill.Pictures, picture)
			continue
		}
		cropPic := projectBill.DownloadPath + imgPng
		cmdConvert := fmt.Sprintf(`convert -resize 1600x1600 %s %s`, downPic, cropPic)
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

		//处理tif多页问题
		cmd = fmt.Sprintf("identify %s | wc -l", projectBill.DownloadPath+picture)
		global.GLog.Info("cmd", zap.Any("s", cmd))
		err, stdout, stderr := utils.ShellOut(cmd)
		if err != nil {
			global.GLog.Error("load-cmd-err identify", zap.Error(err))
			global.GLog.Error("load-cmd-stderr identify" + stderr)
		}
		global.GLog.Info("load-cmd-stdout identify::" + stdout)
		count, err := strconv.Atoi(strings.TrimSpace(stdout))
		if err != nil {
			global.GLog.Error("load-cmd-err count err", zap.Error(err))
		}
		global.GLog.Info("load-cmd-stdout identify count::", zap.Any("", count))
		for i := 0; i < count; i++ {
			global.GLog.Info("picture count", zap.Any(strconv.Itoa(i), re.ReplaceAllString(picture, "-"+strconv.Itoa(i)+".png")))
			projectBill.Pictures = append(projectBill.Pictures, re.ReplaceAllString(picture, "-"+strconv.Itoa(i)+".png"))
		}
	}
	//wg.Wait()
	return nil, projectBill
}
