package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"server/global"
	"server/module/download/project"
	"server/module/download/service"
	"server/module/pro_manager/model"
	"strings"
	"time"
)

func FetchBill(projectBill model.ProjectBill) (error, model.ProjectBill) {
	/* 这是我的第一个简单的程序 */
	err, total := service.CountBill(projectBill.ProCode, projectBill.BillName)
	if err != nil {
		fmt.Println("errerr:", err)
		return err, projectBill
	}
	if total > 0 {
		fmt.Println("errerr:", err)
		_, projectBill := Clean(projectBill)
		return errors.New("已下载过"), projectBill
	}
	//  global.GConfig.LocalUpload.FilePath +
	fmt.Println("CreatedAtCreatedAt:", projectBill.CreatedAt)
	createdAt := time.Unix(projectBill.CreatedAt.Unix(), 0)
	fmt.Println("createdAt:", createdAt)
	projectBill.DownloadPath = projectBill.ProCode + "/download/" + createdAt.Format("2006") + "/" + createdAt.Format("01-02") + "/"
	fmt.Println("DownloadPathDownloadPath:", projectBill.DownloadPath)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	dir = strings.Replace(dir, "bin", "", 1)
	fmt.Println(dir)
	os.MkdirAll(dir+global.GConfig.LocalUpload.FilePath+projectBill.DownloadPath, os.ModePerm)
	// cmd := "curl -u myftp:myftp -f -s -k ftp://192.168.202.3/ZYLP/"
	downFile := dir + global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath + projectBill.BillName + ".zip"
	urlFile := projectBill.BillName + ".zip"
	cmd := fmt.Sprintf(global.GProConf[global.GConfig.System.ProCode].DownloadPaths.FetchBill, urlFile, downFile)
	//cmd := fmt.Sprintf("curl 'ftps://202.108.196.190:991/input/%s' -f -s -k -o '%s' -u ICSBPOzhuhaihuiliu:ICSBPOzhuhaihuiliu#123", urlFile, downFile)
	//cmd := fmt.Sprintf("curl 'ftp://192.168.0.50/ZYLP/%s' -f -s -k -o '%s' -u myftp:myftp", urlFile, downFile)
	//cmd := fmt.Sprintf("docker exec -it download_system /usr/bin/curl 'ftps://202.108.196.163:990/input/%s' -f -s -k -o '%s' -u ICSBPOzhuhaihuiliu:test", urlFile, downFile)
	//cmd := fmt.Sprintf("docker exec -it download_system /usr/bin/curl 'ftps://202.108.196.190:991/input/%s' -f -s -k -o '%s' -u ICSBPOzhuhaihuiliu:ICSBPOzhuhaihuiliu#123", urlFile, downFile)
	fmt.Println("cmdcmdcmdcmdcmd:", cmd)
	err, stdout, _ := project.SpecialShell(cmd)
	if err != nil {
		fmt.Println("errerr:", err, stdout)
		return err, projectBill
	}

	return err, projectBill
}
