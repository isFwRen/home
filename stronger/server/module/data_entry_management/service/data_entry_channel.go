package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"server/global"
	"server/module/data_entry_management/model"
	load "server/module/load/model"
	proConf "server/module/pro_conf/model"
	sybase "server/module/sys_base/model"
)

func GetDataEntryChannelInformation(ProPermission []sybase.SysProPermission) (error, []model.EntryChannel) {
	var EC []model.EntryChannel
	errMsg := ""
	fmt.Println(global.ProDbMap)
	for _, permission := range ProPermission {
		db := global.ProDbMap[permission.ProCode+"_task"]
		if db == nil {
			errMsg += permission.ProCode + ":" + global.ProDbErr.Error() + "\n ;"
			continue
		}
		err, ecitem := Statistics(permission.HasOp0, permission.HasOp1, permission.HasOp2, permission.HasOpq, db)
		if err != nil {
			return err, nil
		}

		//录入网址
		var p proConf.SysProject
		err = global.GDb.Model(&proConf.SysProject{}).Where("code = ? ", permission.ProCode).Find(&p).Error
		if err != nil {
			return err, nil
		}
		ecitem.InnerIp = p.InnerIp
		ecitem.OutIp = p.OutIp
		ecitem.OutAppPort = p.OutAppPort
		ecitem.InAppPort = p.InAppPort
		ecitem.BackEndPort = p.BackEndPort
		ecitem.ProCode = permission.ProCode
		EC = append(EC, ecitem)
	}
	return errors.New(errMsg), EC
}

func Statistics(hasOp0, hasOp1, hasOp2, hasOpq bool, db *gorm.DB) (error, model.EntryChannel) {
	var total int64
	var ECitem model.EntryChannel
	if hasOp0 {
		err := db.Model(&load.ProjectBlock{}).Where("op0_stage = 'op0'").Count(&total).Error
		if err != nil {
			return err, ECitem
		}
		ECitem.Op0Num = total
	}
	if hasOp1 {
		err := db.Model(&load.ProjectBlock{}).Where("op1_stage = 'op1'").Count(&total).Error
		if err != nil {
			return err, ECitem
		}
		ECitem.Op1Num = total
	}
	if hasOp2 {
		err := db.Model(&load.ProjectBlock{}).Where("op2_stage = 'op2'").Count(&total).Error
		if err != nil {
			return err, ECitem
		}
		ECitem.Op2Num = total
	}
	if hasOpq {
		err := db.Model(&load.ProjectBlock{}).Where("opq_stage = 'opq'").Count(&total).Error
		if err != nil {
			return err, ECitem
		}
		ECitem.OpqNum = total
	}
	return nil, ECitem
}
