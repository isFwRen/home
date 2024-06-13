package service

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
	"os"
	"server/global"
	"server/module/pro_manager/model"
	"server/module/pro_manager/model/request"
	"server/utils"
	"strings"
)

func GetReimbursementFormTemplate(info request.GetRFT) (err error, list []map[string]interface{}, total int64) {
	db := global.ProDbMap[info.ProCode]
	if db == nil {
		return global.ProDbErr, nil, 0
	}
	db = db.Model(&model.PicInformation{})
	if info.Name != "" {
		db = db.Where("name LIKE ? ", "%"+info.Name+"%")
	}
	var l []model.PicInformation
	err = db.Debug().Order("updated_at desc").Find(&l).Count(&total).Error
	for i, _ := range l {
		l[i].Size = decimal.NewFromFloat(l[i].Size / float64(1024)).RoundCeil(2).InexactFloat64()
	}

	list = []map[string]interface{}{}
	for _, rule := range l {
		list = append(list, utils.Struct2Map(rule))
	}

	return err, list, total
}

func AddReimbursementFormTemplate(proCode string, pics []model.PicInformation) error {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr
	}
	//err := db.Debug().Model(&model.PicInformation{}).Create(&pics).Error

	err := db.Debug().Model(&model.PicInformation{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		UpdateAll: true,
	}).Create(&pics).Error

	return err
}

func DeleteReimbursementFormTemplate(rm request.Rm) (err error) {
	var R model.PicInformation
	db := global.ProDbMap[rm.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	for _, id := range rm.Ids {
		err = global.ProDbMap[rm.ProCode].Model(&model.PicInformation{}).Where("id = ?", id).Delete(&R).Error
		if err != nil {
			return err
		}
	}
	return err
}

func ReNameReimbursementFormTemplate(rn request.RFTRename) error {
	var temp model.PicInformation
	db := global.ProDbMap[rn.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	db = db.Model(&model.PicInformation{})
	err := db.Where("id = ? ", rn.Id).Find(&temp).Error
	if err != nil {
		return err
	}
	oldName := temp.Name
	newName := rn.Name + "." + temp.Types
	newPath := strings.Replace(temp.Path, oldName+"."+temp.Types, newName, -1)
	err = os.Rename(temp.Path, newPath)
	if err != nil {
		return err
	}
	temp.Name = newName
	temp.Path = newPath
	err = db.Where("id = ? ", rn.Id).Updates(&temp).Error
	return err
}
