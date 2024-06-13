package B0116

import (
	"encoding/json"
	"fmt"
	"regexp"
	"server/global"
	"server/module/pro_manager/model"
	"server/utils"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// var limitChan = make(chan bool, 50)

func Convert(proCode string, projectBill model.ProjectBill) (error, model.ProjectBill) {
	/* 这是我的第一个简单的程序 */
	pictures := projectBill.Images
	projectBill.Pictures = []string{}
	projectBill.ImagesType = make([]string, len(projectBill.Images))
	re := regexp.MustCompile(`\.[^\.]+$`)
	for ii, picture := range pictures {
		//解密
		//decrypt

		imgPng := "convert_" + re.ReplaceAllString(picture, ".png")
		fmt.Println("iiii", ii, picture)
		downPath := global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath
		downPic := downPath + picture

		//获取图片类型 清单报销单等
		url := "http://192.168.0.48:18301/classifier"
		if global.GConfig.System.Env == "test" {
			url = "http://192.168.202.17:18001/classifier"
		}
		getImageClassifierCmd := fmt.Sprintf("curl -X POST -F image=@%v %v", downPic, url)
		global.GLog.Info("getImageClassifierCmd", zap.Any("", getImageClassifierCmd))
		err, stdout, stderr := utils.ShellOut(getImageClassifierCmd)
		if err != nil {
			global.GLog.Error("classifier-cmd-err", zap.Error(err))
			global.GLog.Error("classifier-cmd-stderr" + stderr)
			global.GLog.Error("classifier-cmd-stdout::" + stdout)
		}
		classArr := make([][]interface{}, 0)
		err = json.Unmarshal([]byte(stdout), &classArr)
		if err != nil {
			global.GLog.Error("classifier-json-Unmarshal-err", zap.Error(err))
		}
		projectBill.ImagesType[ii] = stdout
		if len(classArr) > 0 && len(classArr[0]) > 0 {
			projectBill.ImagesType[ii] = classArr[0][0].(string)
		}

		cropApic := downPath + "A" + imgPng
		cmd := fmt.Sprintf(`convert -resize x180 %s %s`, downPic, cropApic)
		fmt.Println("cmdcmdcmd", cmd)
		// tif := projectBill.Path + picture
		// png := projectBill.Path + strings.Replace(picture, ".tif", ".png", 1)
		err, stdout, _ = utils.ShellOut(cmd)
		if err != nil {
			fmt.Println("convertconvert", err, stdout)
		}
		// matched, _ := regexp.MatchString(`(tif|TIF|TIFF|pdf)$`, picture)
		// if !matched {
		// 	projectBill.Pictures = append(projectBill.Pictures, picture)
		// 	continue
		// }
		cropPic := downPath + imgPng
		// cmd = fmt.Sprintf(`convert -resize 1600x1600 %s %s`, downPic, cropPic)
		// if RegIsMatch(picture, `pdf$`) {
		// 	pname := re.ReplaceAllString(cropPic, "*.png")
		// 	cmd = fmt.Sprintf(`convert -resize 1600x1600 %s %s;ls %s`, downPic, cropPic, pname)
		// }
		cmd = check_pic(downPic, cropPic)
		fmt.Println("-----------------------check_pic cmdcmdcmd --------------------------:", cmd)
		if cmd == "" {
			projectBill.Pictures = append(projectBill.Pictures, imgPng)
			cp_pic(downPic, cropPic)
			if !RegIsMatch(picture, `png$`) {
				cp_pic(downPic, downPath+re.ReplaceAllString(picture, ".png"))
			}
			continue
		}
		cp_pic(downPic, downPath+re.ReplaceAllString(picture, ".png"))
		err, stdout, _ = utils.ShellOut(cmd)
		fmt.Println("convertconvert", err, stdout)
		if err == nil && RegIsMatch(picture, `pdf$`) {
			pdfs := RegMatchAll(stdout, `.+`)
			fmt.Println("pdfs", pdfs)
			for _, pdf := range pdfs {
				pdf = strings.Replace(pdf, global.GConfig.LocalUpload.FilePath+projectBill.DownloadPath, "", 1)
				fmt.Println("pdf", pdf)
				projectBill.Pictures = append(projectBill.Pictures, pdf)
			}
		} else {
			projectBill.Pictures = append(projectBill.Pictures, imgPng)
		}
		//加密
		//Encrypt

	}
	// return errors.New("测试"), projectBill
	return nil, projectBill
}

// identify 302023030009248-030111-303056501-16.png
// 302023030009248-030111-303056501-16.png PNG 800x1143 800x1143+0+0 8-bit sRGB 1.44919MiB 0.000u 0:00.000

func cp_pic(downPic, cpPic string) {
	cmd := `cp ` + downPic + ` ` + cpPic
	if RegIsMatch(downPic, `(tif|TIF|TIFF|pdf)$`) {
		cmd = fmt.Sprintf(`convert -resize 1600x1600 %s %s;`, downPic, cpPic)
	}
	err, _, _ := utils.ShellOut(cmd)
	fmt.Println("--------赋值原图-----------", err, cmd)
}

func check_pic(downPic, cropPic string) string {
	cmd := `identify ` + downPic
	err, stdout, _ := utils.ShellOut(cmd)
	fmt.Println("convertconvert", err, stdout)
	if err != nil {
		fmt.Println("convertconvert", err, stdout)
	}
	arr := strings.Split(stdout, " ")
	fmt.Println("-------arr-----------", arr)
	if len(arr) < 7 {
		return ""
	}
	isLarge := false
	if strings.Index(arr[6], "MiB") != -1 {
		isLarge = true
	} else {
		size := ParseFloat(strings.Replace(arr[6], "B", "", 1))
		fmt.Println("-------size-----------", size)
		if (size / 1024) > 500 {
			isLarge = true
		}
	}
	fmt.Println("-------isLarge-----------", isLarge)
	if !isLarge {
		matched, _ := regexp.MatchString(`(tif|TIF|TIFF|pdf)$`, downPic)
		if !matched {
			return ""
		}
		if RegIsMatch(downPic, `pdf$`) {
			re := regexp.MustCompile(`\.[^\.]+$`)
			pname := re.ReplaceAllString(cropPic, "*.png")
			return fmt.Sprintf(`convert -resize 1600x1600 %s %s;ls %s`, downPic, cropPic, pname)
		}

	}
	if strings.Index(arr[2], "x") != -1 {
		wh := strings.Split(arr[2], "x")
		percent, _ := strconv.Atoi(wh[0])
		h, _ := strconv.Atoi(wh[1])
		if h < percent {
			percent = h
		}
		if percent > 1000 {
			percent = 100000 / percent
			return `nconvert -overwrite -xall -resize ` + strconv.Itoa(percent) + `% ` + strconv.Itoa(percent) + `% -colors 16  -quiet -out png -o ` + cropPic + ` ` + downPic
		}

	}

	return `nconvert -overwrite -xall -colors 16 -quiet -out png -o ` + cropPic + ` ` + downPic
	// cmd = `wine /usr/local/bin/nconvert.exe -overwrite -xall -resize 30% 30% -colors 16  -quiet -out png -o 1.png 322023040002299-030105-4364730-1.png`
}

func ParseFloat(value string) float64 {
	val, _ := strconv.ParseFloat(value, 64)
	return val
}

func RegMatchAll(value string, query string) []string {
	reg := regexp.MustCompile(query)
	return reg.FindAllString(value, -1)
}
