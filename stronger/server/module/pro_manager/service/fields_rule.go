package service

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"server/global"
	"server/module/pro_conf/model"
	model2 "server/module/pro_manager/model"
	"server/module/pro_manager/model/request"
	"server/utils"
	"strconv"
	"strings"
)

func FieldsRuleSync(proCode string) (err error) {
	var proInformation model.SysProject
	err = global.GDb.Model(&model.SysProject{}).Where("code = ? ", proCode).Find(&proInformation).Error
	if err != nil {
		return
	}
	err = global.GDb.Model(&model2.SysFieldRule{}).Where("pro_id = ? ", proInformation.ID).Delete(&model2.SysFieldRule{}).Error
	if err != nil {
		return
	}
	var fileds []model.SysProField
	err = global.GDb.Model(&model.SysProField{}).Where("pro_id = ? ", proInformation.ID).Find(&fileds).Error
	if err != nil {
		return
	}
	var rule []model2.SysFieldRule
	for _, v := range fileds {
		var ruleItem model2.SysFieldRule
		ruleItem.ProId = v.ProId
		ruleItem.SysFieldConfId = v.ID
		ruleItem.SysFieldName = v.Name
		ruleItem.SysFieldCode = v.Code
		ruleItem.InputRule = "无"
		ruleItem.RulePicture = []string{}
		rule = append(rule, ruleItem)
	}
	err = global.GDb.Model(&model2.SysFieldRule{}).Debug().Create(&rule).Error
	return
}

func GetFieldsRuleList(info request.GFR) (err error, list interface{}, total int64) {
	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)

	var proInformation model.SysProject
	err = global.GDb.Model(&model.SysProject{}).Where("code = ? ", info.ProCode).Find(&proInformation).Error
	if err != nil {
		return
	}

	db := global.GDb.Model(&model2.SysFieldRule{})
	if info.FieldsName != "" {
		db = db.Where("sys_field_name LIKE ? ", "%"+info.FieldsName+"%")
	}
	if info.Rule != "" {
		db = db.Where("input_rule = ? ", info.Rule)
	}
	db = db.Where("pro_id = ? ", proInformation.ID)
	var fieldsRule []model2.SysFieldRule
	err = db.Count(&total).Limit(limit).Offset(offset).Debug().Find(&fieldsRule).Error
	return err, fieldsRule, total
}

func EditFieldsRule(fieldRule model2.SysFieldRule) error {
	return global.GDb.Model(&model2.SysFieldRule{}).
		Where("id = ? ", fieldRule.ID).Updates(&fieldRule).Error
}

func DeleteFieldsRule(ids []string) error {
	for _, id := range ids {
		var rule model2.SysFieldRule
		err := global.GDb.Model(&model2.SysFieldRule{}).Where("id = ? ", id).Find(&rule).Error
		if err != nil {
			return err
		}
		for _, pic := range rule.RulePicture {
			isExist, err := exists(pic)
			if err != nil {
				return err
			}
			if !isExist {
				continue
			}
			err = os.Remove(pic)
			if err != nil {
				return err
			}
		}
		rule.RulePicture = []string{}
		rule.InputRule = "无"
		err = global.GDb.Model(&model2.SysFieldRule{}).Where("id = ? ", rule.ID).Updates(&rule).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func UploadFieldsRule(info []*multipart.FileHeader, proCode []string, c *gin.Context) error {
	failmsg := ""
	for _, file := range info {
		filename := file.Filename
		dst := ""
		// 线上 uploads/file
		basicPath := global.GConfig.LocalUpload.FilePath + "/字段规则/" + proCode[0] + "/"
		sysFieldName := strings.Split(filename, ".")[0]
		// 设置文件需要保存的指定位置并设置保存的文件名字
		if strings.Split(filename, ".")[len(strings.Split(filename, "."))-1] != "png" || strings.Split(filename, ".")[len(strings.Split(filename, "."))-1] != "jpg" {
			dst = path.Join(basicPath, sysFieldName+".png")
		} else {
			dst = path.Join(basicPath, filename)
		}

		fmt.Println(filename + " has been saved in this path " + basicPath)
		fmt.Println("dst", dst)

		var syn model2.SysFieldRule
		var total int64

		err := global.GDb.Model(&model2.SysFieldRule{}).Where("sys_field_name = ? ", sysFieldName).
			Find(&syn).Count(&total).Error
		if err != nil {
			return err
		}
		if total == 0 {
			failmsg += filename + "没有对应的字段配置; "
			continue
		}
		isExist, err := exists(basicPath)
		if err != nil {
			return err
		}
		if !isExist {
			err = os.MkdirAll(basicPath, 0777)
			if err != nil {
				return err
			}
			err = os.Chmod(basicPath, 0777)
			if err != nil {
				return err
			}
		}
		// 上传文件到指定的路径
		saveErr := c.SaveUploadedFile(file, dst)
		if saveErr != nil {
			fmt.Println("saveErr", saveErr)
			continue
		}
		err = global.GDb.Model(&model2.SysFieldRule{}).Where("sys_field_name = ? ", sysFieldName).
			Updates(map[string]interface{}{
				"rule_picture": pq.StringArray{dst},
				"input_rule":   "有",
			}).Error
		if err != nil {
			return err
		}
	}
	if failmsg != "" {
		return errors.New(failmsg)
	}
	return nil
}

func ExportFieldsRule(info request.ExportFR) (error, string) {
	var proInformation model.SysProject
	err := global.GDb.Model(&model.SysProject{}).Where("code = ? ", info.ProCode).Find(&proInformation).Error
	if err != nil {
		return err, ""
	}
	var sysfieldRule []model2.SysFieldRule
	if info.Rule != "" {
		err = global.GDb.Model(model2.SysFieldRule{}).Where("pro_id = ? AND input_rule = ?", proInformation.ID, info.Rule).Find(&sysfieldRule).Error
	} else {
		err = global.GDb.Model(model2.SysFieldRule{}).Where("pro_id = ? ", proInformation.ID).Find(&sysfieldRule).Error
	}
	if err != nil {
		return err, ""
	}
	basicPath := global.GConfig.LocalUpload.FilePath + "字段规则导出/"
	bookName := info.ProCode + "-字段规则" + ".xlsx"
	err = utils.ExportBigExcel(basicPath, bookName, "", sysfieldRule)
	if err != nil {
		return err, ""
	}

	//插入图片
	xlsx, err := excelize.OpenFile(basicPath + bookName)
	if err != nil {
		return err, basicPath + bookName
	}
	for i, v := range sysfieldRule {
		xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), " ")
		xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), i+1)
		if v.InputRule != "有" {
			continue
		}
		_, err := os.Stat(v.RulePicture[0])
		if os.IsNotExist(err) {
			continue
		}
		file, err := ioutil.ReadFile(v.RulePicture[0])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("E" + strconv.Itoa(i+2))
		arr := strings.Split(v.RulePicture[0], ".")
		fmt.Println("ss", arr)
		if err := xlsx.AddPictureFromBytes("Sheet1", "E"+strconv.Itoa(i+2), "", v.SysFieldName, "."+arr[len(arr)-1], file); err != nil {
			fmt.Println(err)
		}
	}
	if err := xlsx.SaveAs(basicPath + bookName); err != nil {
		fmt.Println(err)
	}

	return nil, "files/字段规则导出/" + bookName
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
