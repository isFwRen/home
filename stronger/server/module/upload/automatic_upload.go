/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/30 4:11 下午
 */

package upload

import (
	"errors"
	"fmt"
	"server/global"
	"server/module/pro_manager/model"
	service1 "server/module/upload/service"
	"sync"
)

var lock sync.Mutex

func AutomaticUpload(proCode string) (err error) {
	err, pro := service1.GetProIsAutoReturn(proCode)
	if err != nil {
		return err
	}
	if !pro.AutoReturn {
		global.GLog.Error("该项目不可自动回传")
		return errors.New("该项目不可自动回传")
	}

	err, bills := service1.GetIsAutoUploadBills(proCode)
	if err != nil {
		return err
	}

	for _, bill := range bills {
		req := model.ProCodeAndId{
			ID:      bill.ID,
			ProCode: bill.ProCode,
		}
		project, ok := global.GProConf[proCode]
		if !ok {
			global.GLog.Error(fmt.Sprintf("回传错误id:::%v,错误:::没有找到项目配置:::%v", bill.ID, err))
			return
		}
		err = BillUploadAdapter(req, project.UploadPaths)
		if err != nil {
			global.GLog.Error(fmt.Sprintf("回传错误id:::%v,错误:::%v", bill.ID, err))
			continue
		}
	}

	return err
}
