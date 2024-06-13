/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/12/3 4:24 下午
 */

package api

import (
	"github.com/gin-gonic/gin"
	"server/global"
	"server/module/pro_manager/model"
	"server/module/pro_manager/service"
	api2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/api"
)

// SetLog 插入日志
func SetLog(c *gin.Context, billId string, proCode string, name string, before string, after string) error {
	user, err := api2.GetUserByToken(c)
	if err != nil {
		return err
	}
	var r = model.ResultDataLog{
		Name:      name,
		BillId:    billId,
		BeforeVal: before,
		AfterVal:  after,
		EditCode:  user.Code,
		EditName:  global.UserCodeName[user.Code],
	}
	err = service.InsertLog(proCode, r)
	if err != nil {
		return err
	}
	return nil
}
