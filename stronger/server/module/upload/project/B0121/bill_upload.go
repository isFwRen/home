/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/12/6 17:43
 */

package B0121

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"server/global"
	model2 "server/module/pro_conf/model"
	"server/module/pro_manager/model"
	"server/module/pro_manager/service"
	service1 "server/module/upload/service"
	"server/utils"
	"strings"
	"time"
)

//var lock sync.Mutex

// BillUpload 回传单据
func BillUpload(reqParam model.ProCodeAndId, uploadPath model2.SysProDownloadPaths) error {
	//lock.Lock()
	//defer lock.Unlock()
	err, obj := service.GetProBillById(reqParam)
	if err != nil {
		return err
	}
	if obj.Stage != 3 && obj.Stage != 4 {
		return errors.New("回传查询单证失败,状态有误")
	}
	//time.Sleep(time.Second * 10)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	dir = strings.Replace(dir, "bin", "", 1)
	xmlFilePath := dir + global.GConfig.LocalUpload.FilePath + obj.ProCode + "/upload_xml/" +
		fmt.Sprintf("%v/%v/%v/",
			obj.CreatedAt.Year(), int(obj.CreatedAt.Month()),
			obj.CreatedAt.Day())
	global.GLog.Info(obj.BillNum)
	//sh := "curl -T %v%v.xml 'ftp://192.168.202.3/lipei2.0/%v.temp' -u 'myftp:myftp' -s -x socks5://127.0.0.1:8090"
	sh := fmt.Sprintf(uploadPath.Upload, xmlFilePath+obj.BillNum+".xml")
	global.GLog.Info("上传cmd:::", zap.Any("sh", sh))
	err, s, s2 := utils.ShellOut(sh)
	global.GLog.Error("上传错误err", zap.Error(err))
	global.GLog.Info("上传回显", zap.Any("std out", s))
	global.GLog.Error("上传回显", zap.Any("std err", s2))
	if err != nil {
		return err
	}
	err = service1.UpdateBill(reqParam, time.Now())

	return err
	//lock.Unlock()
}
