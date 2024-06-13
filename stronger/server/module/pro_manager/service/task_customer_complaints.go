package service

import (
	"server/global"
	"server/module/pro_manager/model"
	modelBase "server/module/sys_base/model"
)

func GetCustomerComplaints(info model.CustomerComplaints, uid string) (err error, list interface{}, total int64) {
	var U modelBase.SysUser
	err = global.GDb.Model(&modelBase.SysUser{}).Where("id = ? ", uid).Find(&U).Error
	if err != nil {
		return err, nil, 0
	}

	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
	db := global.GDb.Debug().Model(&model.Quality{})
	//if info.ProCode != "all" {
	//	db = db.Where("pro_code LIKE ? ", "%"+info.ProCode+"%")
	//}
	//

	//db = db.Raw(" SELECT * FROM (SELECT * FROM \"qualities\"\n WHERE pro_code LIKE '%" + info.ProCode + "%'\n    AND bill_name LIKE '%" + info.BillName + "%'\n    AND wrong_field_name LIKE '%" + info.WrongFieldName + "%'\n    AND op0_responsible_code LIKE '%" + U.Code + "%'\n    OR op1_responsible_code LIKE '%" + U.Code + "%'\n    OR op2_responsible_code LIKE '%" + U.Code + "%'\n    OR opq_responsible_code LIKE '%" + U.Code + "%') a WHERE a.month = '" + info.Month + "'")
	//SELECT * FROM (SELECT * FROM "qualities"
	//WHERE op0_responsible_code LIKE '%%'
	//OR op1_responsible_code LIKE '%%'
	//OR op2_responsible_code LIKE '%%'
	//OR opq_responsible_code LIKE '%%') a WHERE a.month = '' and a.pro_code LIKE '%%' and a.bill_name LIKE '%%' and a.wrong_field_name LIKE '%%'
	sql := "SELECT * FROM (SELECT * FROM \"qualities\"\n WHERE op0_responsible_code LIKE '%" + U.Code + "%'\n    OR op1_responsible_code LIKE '%" + U.Code + "%'\n    OR op2_responsible_code LIKE '%" + U.Code + "%'\n    OR opq_responsible_code LIKE '%" + U.Code + "%') a "
	if info.Month != "" {
		sql += " WHERE a.month = '" + info.Month + "'"
	}

	if info.ProCode != "" {
		sql += "and a.pro_code LIKE '%" + info.ProCode + "%'"
	}

	if info.BillName != "" {
		sql += "and a.bill_name LIKE '%" + info.BillName + "%'"
	}

	if info.WrongFieldName != "" {
		sql += "and a.wrong_field_name LIKE '%" + info.WrongFieldName + "%'"
	}

	db = db.Raw(sql)
	var lr []model.Quality
	err = db.Limit(limit).Offset(offset).Find(&lr).Error
	if err != nil {
		return nil, nil, total
	}
	if err != nil {
		return nil, nil, total
	}

	return nil, lr, int64(len(lr))
}
