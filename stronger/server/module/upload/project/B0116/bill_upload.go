package B0116

import (
	"errors"
	"fmt"
	"os"
	"server/global"
	model2 "server/module/pro_conf/model"
	"server/module/pro_manager/model"
	"server/module/pro_manager/service"
	service1 "server/module/upload/service"
	"server/utils"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
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

	countNum := utils.GetNodeValue(obj.OtherInfo, "countNum")
	num, _ := strconv.ParseInt(countNum, 10, 64)

	err, bNum := service1.GetNumByBatchNum(reqParam, obj.BatchNum)
	if err != nil {
		return err
	}
	// if obj.Stage != 3 && obj.Stage != 4 {
	// 	return errors.New("回传查询单证失败,状态有误")
	// }

	// if utils.RegIsMatch(`(节点值不能为空|晚于当前系统时间|出院日期早于入院日期)`, obj.WrongNote) {
	// 	return errors.New("无法回传，请根据导出校验修改案件数据！")
	// }

	// _, fc189 := lservice.SelectBillFields(obj.ProCode, obj.ID, -1, "fc189")
	// // for _, field := range obj.

	// if utils.RegIsMatch(`(差额)`, obj.WrongNote) && utils.RegIsMatch(`^(A|)$`, fc189.ResultValue) {
	// 	return errors.New("无法回传，请根据导出校验修改案件数据！")
	// }
	//压缩

	//time.Sleep(time.Second * 10)
	xmlFilePath := global.GConfig.LocalUpload.FilePath + obj.ProCode + "/upload_xml/" +
		fmt.Sprintf("%v/%v/%v/",
			obj.CreatedAt.Year(), int(obj.CreatedAt.Month()),
			obj.CreatedAt.Day())

	global.GLog.Info(obj.BillNum)
	global.GLog.Info(obj.BatchNum)

	//CSB0108RC0287000
	buf, err := os.ReadFile(xmlFilePath + obj.BillNum + ".xml")
	newCase := utils.GetNode(string(buf), "newcase")
	if newCase == nil {
		return errors.New("案件报文不完整，无法回传")
	}

	zipCmd := fmt.Sprintf("zip -j %v.zip %v.xml", xmlFilePath+obj.BatchNum, xmlFilePath+obj.BillNum)
	global.GLog.Info("压缩cmd:::", zap.Any("sh", zipCmd))
	err, s, s2 := utils.ShellOut(zipCmd)
	global.GLog.Error("压缩错误err", zap.Error(err))
	global.GLog.Info("压缩回显", zap.Any("std out", s))
	global.GLog.Error("压缩回显", zap.Any("std err", s2))

	if bNum != num-1 || strings.Index(obj.UploadAt.String(), "0001-01-01") == -1 {
		err = service1.UpdateStageBill(reqParam)
		return err
	}

	//sh := "curl -T %v%v.xml 'ftp://192.168.202.3/lipei2.0/%v.temp' -u 'myftp:myftp' -s -x socks5://127.0.0.1:8090"
	sh := fmt.Sprintf(uploadPath.Upload, xmlFilePath, obj.BatchNum, obj.BatchNum)
	global.GLog.Info("上传cmd:::", zap.Any("sh", sh))
	err, s, s2 = utils.ShellOut(sh)
	global.GLog.Error("上传错误err", zap.Error(err))
	global.GLog.Info("上传回显", zap.Any("std out", s))
	global.GLog.Error("上传回显", zap.Any("std err", s2))
	if err != nil {
		return err
	}
	sh = fmt.Sprintf(uploadPath.UploadRename, obj.BatchNum, obj.BatchNum)
	global.GLog.Info("上传重命名cmd:::", zap.Any("sh", sh))
	//curl ftp://192.168.202.3/ -u 'myftp:myftp' -Q 'CWD lipei2.0/' -Q 'RNFR ./%v.temp' -Q 'RNTO ./%v.xml' -s -x socks5://127.0.0.1:8090
	err, s, s2 = utils.ShellOut(sh)
	global.GLog.Error("上传重命名错误err", zap.Error(err))
	global.GLog.Info("上传重命名回显", zap.Any("std out", s))
	global.GLog.Error("上传重命名回显", zap.Any("std err", s2))
	if err != nil {
		return err
	}

	// err = service1.UpdateBill(reqParam, time.Now())
	err = service1.UpdateBillByBatch(reqParam.ProCode, obj.BatchNum, time.Now())

	return err
	//lock.Unlock()
}
