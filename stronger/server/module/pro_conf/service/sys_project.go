package service

import (
	"errors"
	"fmt"
	"server/global"
	"server/module/pro_conf/model"
	requestProConf "server/module/pro_conf/model/request"
	"server/module/pro_conf/model/response"
	model2 "server/module/pro_manager/model"
	sybase "server/module/sys_base/model"
	requestSysBase "server/module/sys_base/model/request"
	"strconv"
	"strings"
	"time"
	model3 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/model"

	"github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// @title    SysProject
// @description   AddSysProject, 新增项目配置
// @auth                     （2020年10月20日14:58:59）
// @param     sysProject      model.SysProject
// @return    err             error

func AddSysProject(sysProject model.SysProject, customClaims *model3.CustomClaims) (err error) {
	tx := global.GDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	var reSysProject model.SysProject
	// 判断是否有该项目
	err1 := global.GDb.Where("name = ? AND code = ? AND type = ?", sysProject.Name, sysProject.Code, sysProject.Type).First(&reSysProject).Error
	if !errors.Is(err1, gorm.ErrRecordNotFound) {
		return errors.New("已存在该项目配置")
	} else {
		if err = tx.Create(&sysProject).Error; err != nil {
			tx.Rollback()
			return err
		}

		//默认添加一个导出模板
		var sysExport model.SysExport
		sysExport.ProName = sysProject.Name
		sysExport.ProId = sysProject.ID
		sysExport.TempVal = "init"
		sysExport.XmlType = "utf-8"
		if err = tx.Model(model.SysExport{}).Create(&sysExport).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err = tx.Model(&sybase.SysProPermission{}).Where("pro_code = ?", sysProject.Code).
			Update("pro_id", sysProject.ID).Error; err != nil {
			tx.Rollback()
			return err
		}
		var row int64
		if tx.Model(&sybase.SysProPermission{}).Where("user_id = ? and pro_code = ?", customClaims.ID, sysProject.Code).Count(&row); row == 0 {
			err = tx.Model(&sybase.SysProPermission{}).Create(&sybase.SysProPermission{
				ProCode:   sysProject.Code,
				ProId:     sysProject.ID,
				ProName:   sysProject.Name,
				HasOp0:    true,
				HasOp1:    true,
				HasOp2:    true,
				HasOpq:    true,
				HasInNet:  true,
				HasOutNet: true,
				UserCode:  customClaims.Code,
				UserId:    customClaims.ID,
				ObjectId:  "",
				HasPm:     true,
			}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		//新增两条下载路径
		d := []model.SysProDownloadPaths{
			{DownloadType: 1, Scan: "1", IsDownload: true, IsUpload: true, MaxDownload: 5, MaxConnections: 5, ProId: sysProject.ID, ProCode: sysProject.Code, ProName: sysProject.Name},
			{DownloadType: 0, Scan: "2", IsDownload: false, IsUpload: false, MaxDownload: 5, MaxConnections: 5, ProId: sysProject.ID, ProCode: sysProject.Code, ProName: sysProject.Name},
		}
		err = tx.Model(&model.SysProDownloadPaths{}).Create(&d).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		//更新到Redis
		//value, err1 := json.Marshal(sysProject)
		//if err1 != nil {
		//	global.GLog.Error("ProConf json marshal err")
		//	return
		//}
		//global.GRedis.RPush(global.GProConf, value).Err()
	}

	return tx.Commit().Error
	//return err

}

// @title    SysProject
// @description   GetSysProject, 获取项目配置 分页
// @auth      星期一          （2020年10月21日09:00:56）
// @param     info            request.SysProjectRecordSearch
// @return    err error, list interface{}, total int64

func GetSysProjectByPage(info requestSysBase.SysProjectRecordSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.PageIndex - 1)
	// 创建db
	db := global.GDb.Model(&model.SysProject{})
	var sysProjectList []model.SysProject
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Name != "" {
		db = db.Where("name = ?", info.Name)
	}
	err = db.Count(&total).Error
	err = db.Order("id desc").Limit(limit).Offset(offset).Preload("SysProTemplate").Find(&sysProjectList).Error
	return err, sysProjectList, total
}

// @title    SysProject
// @description   GetSysProjectList, 获取所有项目配置
// @auth      星期一          （2020年10月21日09:00:56）
// @return   err error, list interface{}, total int64

func GetSysProjectList() (err error, list interface{}, total int64) {
	// 创建db
	db := global.GDb.Model(&model.SysProject{})
	var sysProjectList []model.SysProject
	err = db.Count(&total).Error
	err = db.Order("created_at desc").Preload("SysProTemplate", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at asc")
	}).Find(&sysProjectList).Error

	//rLen, err := global.G_REDIS.LLen("router").Result()
	//lists, _ := global.G_REDIS.LRange("router", 0, rLen-1).Result()
	//
	//reLists := make([]model.SysProject,0)
	//for _, v := range lists {
	//	var sysPro model.SysProject
	//	json.Unmarshal([]byte(v),&sysPro)
	//	reLists = append(reLists, sysPro)
	//}
	return err, sysProjectList, total
}

// @title    SysProject
// @description   UpdateSysProjectById, 更新项目配置
// @auth      星期一          （2020年10月21日09:00:56）
// @param    sysProjectRq request.SysProjectUpdateRecord
// @return   rows int64

func UpdateSysProjectById(sysProjectRq requestProConf.SysProjectUpdateRecord) (err error) {
	// 创建db
	var SysProject model.SysProject
	SysProject.ID = sysProjectRq.Id
	SysProject.EditVersion = sysProjectRq.EditVersion
	db := global.GDb.Model(&model.SysProject{})
	m1 := map[string]interface{}{
		//"name":sysProject.Name,
		"auto_return":  sysProjectRq.AutoReturn,
		"cache_time":   sysProjectRq.CacheTime,
		"restart_at":   sysProjectRq.RestartAt,
		"save_date":    sysProjectRq.SaveDate,
		"edit_version": sysProjectRq.EditVersion + 1,
		"type":         sysProjectRq.Type,
	}
	row := db.Model(&SysProject).Where("id = ? and edit_version = ?", sysProjectRq.Id, sysProjectRq.EditVersion).
		Updates(m1).RowsAffected
	global.GLog.Info("更新配置", zap.Any("行数", row))
	if row == 0 {
		return errors.New("没有更新到配置")
	}

	//proDb := global.ProDbMap[sysProjectRq.Code]
	//if db == nil {
	//	return global.ProDbErr
	//}
	//row = proDb.Model(&model2.ProjectBill{}).Where("stage != ? and stage != ? ", 5, 6).
	//	Update("is_auto_upload", sysProjectRq.AutoReturn).RowsAffected
	//global.GLog.Info("更新单据", zap.Any("行数", row))
	return nil
}

// @title    SysProject
// @description   delete SysProject
// @auth                     （2020/04/05  20:22）
// @param     ids request.ReqIds
// @return    reRows int64

func RmSysProjectByIds(ids requestSysBase.ReqIds) (reRows int64) {
	rows := global.GDb.Delete(&[]model.SysProject{}, ids.Ids).RowsAffected
	return rows
}

// @title    SysProTemplate
// @description   AddSysProTemplate, 新增项目模板配置
// @auth                     （2020年10月20日14:58:59）
// @param     sysProTemplate model.SysProTemplate
// @return    err             error

func AddSysProTemplate(sysProTemplate model.SysProTemplate) (err error) {
	var reSysProTemplate model.SysProTemplate
	// 判断是否有该项目
	err1 := global.GDb.Where("name = ? AND pro_id = ? ", sysProTemplate.Name, sysProTemplate.ProId).First(&reSysProTemplate).Error
	if !errors.Is(err1, gorm.ErrRecordNotFound) {
		return errors.New("已存在该项目模板配置")
	} else {

		err = global.GDb.Create(&sysProTemplate).Error
		if err != nil {
			return err
		}
		////更新到Redis
		//value, err1 := json.Marshal(sysProTemplate)
		//if err1 != nil {
		//	global.G_LOG.Error("sysProTemplate json marshal err" )
		//	return
		//}
		//global.G_REDIS.RPush(global.GProTempConf, value).Err()
	}
	return err
}

// @title    SysProject
// @description   GetSysProjectList, 获取所有项目配置
// @auth      星期一          （2020年10月23日11:24:29）
// @param     id            string
// @return    err error, list interface{}

func GetSysProTemplateById(id int) (err error, list interface{}) {
	var reSysProTemplate response.SysProTemplate
	err = global.GDb.Where("id = ?", id).First(&reSysProTemplate).Error
	//err = global.G_DB.Where("id = ?", id).Preload("SysProTempB").First(&reSysProTemplate).Error
	return err, reSysProTemplate
}

// @title    SysProject
// @description   GetSysProTempListByProId, 根据项目id获取所有模板
// @auth      星期一          （2020年10月27日11:18:52）
// @param     proId            	string
// @return    err error,list   interface{}

func GetSysProTempListByProId(proId string) (err error, list interface{}) {
	var sysProTemplateList []model.SysProTemplate
	err = global.GDb.Model(&model.SysProTemplate{}).Where("pro_id = ?", proId).Find(&sysProTemplateList).Error
	//err = global.G_DB.Where("id = ?", id).Preload("SysProTempB").First(&reSysProTemplate).Error
	return err, sysProTemplateList
}

// @title    SysProject
// @description   GetSysProTempBlockByTempId, 根据模板id获取项目模板分块
// @auth      星期一          （2020年10月23日16:11:30）
// @param     info 			request.InfoSearch
// @return    err error, list interface{}, total int64

func GetSysProTempBlockByTempId(tempId string, name string) (err error, list interface{}, total int64, max int) {
	//limit := info.PageSize
	//offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GDb.Model(&model.SysProTempB{}).Where("pro_temp_id = ?", tempId)
	var sysProTempBList []model.SysProTempB
	// 如果有条件搜索 下方会自动创建搜索语句
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	err = db.Count(&total).Error
	//err = db.Order("id").Limit(limit).Offset(offset).Find(&sysProTempBList).Error
	//err = db.Order("code").Find(&sysProTempBList).Error
	err = db.Order("my_order").Find(&sysProTempBList).Error

	global.GDb.Raw("SELECT MAX( sys_pro_temp_bs.my_order ) AS max FROM sys_pro_temp_bs WHERE pro_temp_id = ?", tempId).Scan(&max)

	return err, sysProTempBList, total, max + 1
}

// @title    SysProject
// @description   	UpdateImages, 更新模板图片
// @auth      星期一          （2020年10月23日16:11:30）
// @param     images []string,sysProTempId int
// @return    reRows

func UpdateImages(images []string, sysProTempId string) (reRows int64) {
	reRow := global.GDb.Model(&model.SysProTemplate{}).Where("id = ?", sysProTempId).Updates(model.SysProTemplate{
		Images: images,
	}).RowsAffected
	return reRow
}

// @title    SysProject
// @description   RmAllAndAddSysProTempBlockByTempId, 删除所有然后插入新的
// @auth      星期一          （2020年10月27日10:07:15）
// @param     pro            sysProjectBlocks
// @return    err error,sysProTempBList []model.sys_pro_temp

func RmAllAndAddSysProTempBlockByTempId(sysProjectBlocks requestProConf.SysProjectBlocks) (err error, sysProTempBList []model.SysProTempB) {
	//func RmAllAndAddSysProTempBlockByTempId(sysProjectBlocks []model.SysProTempB,tempId int) (err error) {
	//	global.G_DB.Transaction(func(tx *gorm.DB) error {
	//		if err := tx.Where("pro_temp_id = ?", sysProjectBlocks.TempId).Delete(&model.SysProTempB{}).Error; err != nil {
	//			return err
	//		}
	//
	//		if err := tx.Omit("id").Create(sysProjectBlocks.SysProTempBArr).Error; err != nil {
	//			return err
	//		}
	//
	//		// 返回 nil 提交事务
	//		return nil
	//	})
	//	return nil
	tx := global.GDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var reSysProTempBList []model.SysProTempB
	if err := tx.Error; err != nil {
		return err, reSysProTempBList
	}

	if err := tx.Where("pro_temp_id = ?", sysProjectBlocks.TempId).Delete(&model.SysProTempB{}).Error; err != nil {
		return err, reSysProTempBList
	}

	if err := tx.Create(&sysProjectBlocks.SysProTempBArr).Error; err != nil {
		return err, reSysProTempBList
	}

	return tx.Commit().Error, sysProjectBlocks.SysProTempBArr
}

// @title    SysProject
// @description   AddBlock, 新增block配置
// @auth                     （2020年10月20日14:58:59）
// @param     SysProTempB      model.SysProTempB
// @return    err             error

func AddBlock(sysProTempB model.SysProTempB) (err error) {
	tx := global.GDb.Begin()

	max := "bc000"
	err = tx.Raw("SELECT MAX( sys_pro_temp_bs.code ) AS max FROM sys_pro_temp_bs WHERE pro_temp_id = ?", sysProTempB.ProTempId).Scan(&max).Error
	if err != nil {
		global.GLog.Error("出错 %v")
	}
	//if sysProTempB.Name == "未定义" {
	//	sysProTempB.FreeTime = 3600
	//}
	maxInt, _ := strconv.Atoi(strings.Replace(max, "bc", "", -1))
	sysProTempB.Code = appendBc(maxInt + 1)

	err = tx.Create(&sysProTempB).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	var blockTeachVideo model2.TeachVideo
	var proInformation model.SysProTemplate
	err = tx.Model(&model.SysProTemplate{}).Where("id = ? ", sysProTempB.ProTempId).Find(&proInformation).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	blockTeachVideo.ProId = proInformation.ID
	blockTeachVideo.SysBlockCode = sysProTempB.Code
	blockTeachVideo.SysBlockName = sysProTempB.Name
	blockTeachVideo.SysBlockConfId = sysProTempB.ID
	blockTeachVideo.InputRule = "无"
	blockTeachVideo.Video = pq.StringArray{}

	err = tx.Model(&model2.TeachVideo{}).Create(&blockTeachVideo).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return err
}

func appendBc(maxInt int) (max string) {
	switch {
	case maxInt < 10:
		return "bc00" + strconv.Itoa(maxInt)
	case maxInt < 100:
		return "bc0" + strconv.Itoa(maxInt)
	default:
		return "bc" + strconv.Itoa(maxInt)
	}
}

// @title    SysProject
// @description   EditBlock, 修改block配置
// @auth                     （2020年10月20日14:58:59）
// @param     SysProTempB      model.SysProTempB
// @return    err             error

func EditBlock(sysProTempB model.SysProTempB) (err error) {
	err = global.GDb.Where("id = ?", sysProTempB.ID).Updates(&sysProTempB).Updates(map[string]interface{}{
		"f_eight":        sysProTempB.FEight,
		"relation":       sysProTempB.Relation,
		"is_mobile":      sysProTempB.IsMobile,
		"free_time":      sysProTempB.FreeTime,
		"name":           sysProTempB.Name,
		"ocr":            sysProTempB.Ocr,
		"is_loop":        sysProTempB.IsLoop,
		"is_competitive": sysProTempB.IsCompetitive,
	}).Error
	return err
}

// @title    SysProject
// @description   	UpdateBlockCoordinate, 分块截图配置
// @auth      星期一          （2020年10月27日15:44:40）
// @param     sysReqBlockCoordinate	request.SysReqBlockCoordinate
// @return    reRow int64

func UpdateBlockCoordinate(sysReqBlockCoordinate requestProConf.SysReqBlockCoordinate) (reRow int64) {
	if sysReqBlockCoordinate.CoordinateType {
		reRow = global.GDb.Model(&model.SysProTempB{}).Where("id = ?", sysReqBlockCoordinate.BlockId).Updates(map[string]interface{}{
			"pic_page":     sysReqBlockCoordinate.PicPage,
			"w_coordinate": sysReqBlockCoordinate.Coordinate,
		}).RowsAffected
	} else {
		reRow = global.GDb.Model(&model.SysProTempB{}).Where("id = ?", sysReqBlockCoordinate.BlockId).Updates(map[string]interface{}{
			"m_pic_page":     sysReqBlockCoordinate.PicPage,
			"m_w_coordinate": sysReqBlockCoordinate.Coordinate,
		}).RowsAffected
	}
	return reRow
}

// @title    SysProject
// @description   	RmAllAndAddBlockRelations, 删除所有然后插入新的分块关系
// @auth      星期一          （2020年10月28日15:02:32）
// @param     sysBlocksRelationsObj request.SysBlocksRelations
// @return    err , []model.TempBlockRelation

func RmAllAndAddBlockRelations(sysBlocksRelationsObj requestProConf.SysBlocksRelations) (err error, tempBlockRelationList []model.TempBlockRelation) {
	tx := global.GDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			global.GLog.Error(fmt.Sprintf("RmAllAndAddBlockRelations:::%v", r))
			tx.Rollback()
		}
	}()
	var reSysProTempBList []model.TempBlockRelation
	if err := tx.Error; err != nil {
		return err, reSysProTempBList
	}

	if err := tx.Where("temp_b_id = ? AND my_type = ?", sysBlocksRelationsObj.BlockId, sysBlocksRelationsObj.MyType).Delete(&model.TempBlockRelation{}).Error; err != nil {
		return err, reSysProTempBList
	}

	if len(sysBlocksRelationsObj.TempBlockRelationArr) > 0 {
		if err := tx.Create(&sysBlocksRelationsObj.TempBlockRelationArr).Error; err != nil {
			return err, reSysProTempBList
		}
	}

	return tx.Commit().Error, sysBlocksRelationsObj.TempBlockRelationArr
}

// @title    SysProject
// @description   	RmAllAndAddBlockRelations, 删除所有然后插入新的分块关系
// @auth      星期一          （2020年10月28日15:02:32）
// @param     bId string
// @param     myType string
// @return    err , []model.TempBlockRelation

func GetBlockRelationsByBIdWithType(bId string, myType string) (err error, tempBlockRelationList []model.TempBlockRelation) {

	var tempBlockRelation []model.TempBlockRelation
	err = global.GDb.Model(&model.TempBlockRelation{}).Where("temp_b_id = ? AND my_type = ?", bId, myType).Find(&tempBlockRelation).Error
	return err, tempBlockRelation
}

// @title    SysProject
// @description   	CopyTemp, 复制模板
// @auth      星期一          （2020年10月30日16:17:40）
// @param     tId 		string
// @param     tempName 	string
// @return    err

func CopyTemp(tId string, tempName string) (err error, tempId string) {
	tx := global.GDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			global.GLog.Error(fmt.Sprintf("%v", r))
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		panic(err)
	}

	//查询模板
	var sysProTemplate response.SysProTemplate
	//global.G_LOG.Info("sysProTemplate id:" + string(tId))

	if global.GDb.Where("id = ?", tId).First(&sysProTemplate).Error != nil {
		return errors.New("template record not found"), ""
	}

	var ccc int64
	global.GDb.Model(model.SysProTemplate{}).Where("name = ? AND pro_id = ?", tempName, sysProTemplate.ProId).Count(&ccc)
	if ccc > 0 {
		return errors.New("有重复的模板"), ""
	}
	//global.G_LOG.Info("sysProTemplate name:" + sysProTemplate.Name)

	//插入新模板
	sysProTemplate.Name = tempName
	sysProTemplate.ID = ""
	//if global.G_DB.Omit("id").Create(&sysProTemplate).Error != nil {
	if tx.Create(&sysProTemplate).Error != nil {
		return err, ""
	}
	//global.G_LOG.Info("sysProTemplate new id:" + string(sysProTemplate.ID))
	//global.G_LOG.Info("sysProTemplate new name:" + sysProTemplate.Name)

	//查分块和关系
	var sysBlockWithRelation []response.SysProTempB
	err = global.GDb.Where("pro_temp_id = ?", tId).
		Preload("OneRelations", "my_type = ?", "0").
		Preload("TwoRelations", "my_type = ?", "1").
		Preload("ThreeRelations", "my_type = ?", "2").
		Preload("TempBFRelation").
		Find(&sysBlockWithRelation).Error

	//插入分块  将新分块id存到切片
	var newIdIndexBlockList []response.NewIdIndexBlock
	for i, proBlock := range sysBlockWithRelation {
		var newIdIndexBlock response.NewIdIndexBlock
		proBlock.ProTempId = sysProTemplate.ID
		bTemp := sysBlockWithRelation[i]
		//global.G_LOG.Info("b######################")
		//global.G_LOG.Info("block old id:" + string(proBlock.ID))
		//err = global.G_DB.Model(model.SysProTempB{}).Omit("id").Create(&proBlock.SysProTempB).Error
		proBlock.SysProTempB.ID = ""
		err = tx.Model(model.SysProTempB{}).Create(&proBlock.SysProTempB).Error
		if err != nil {
			panic(err)
		}

		//插入关系
		newIdIndexBlock.NewId = proBlock.ID
		ReNewIdIndexBlock(&newIdIndexBlock, &bTemp, sysBlockWithRelation, 0)
		ReNewIdIndexBlock(&newIdIndexBlock, &bTemp, sysBlockWithRelation, 1)
		ReNewIdIndexBlock(&newIdIndexBlock, &bTemp, sysBlockWithRelation, 2)
		newIdIndexBlockList = append(newIdIndexBlockList, newIdIndexBlock)

		//插入新的分块字段关系
		for j, tempBFRelationObj := range proBlock.TempBFRelation {
			tempBFRelationObj.TempBId = proBlock.ID
			tempBFRelationObj.UpdatedAt = time.Now().Add(time.Millisecond * time.Duration(j))
			tempBFRelationObj.CreatedAt = time.Now()
			tempBFRelationObj.ID = ""
			fmt.Println(tempBFRelationObj.UpdatedAt)
			//err = global.G_DB.Model(model.TempBFRelation{}).
			err = tx.Model(model.TempBFRelation{}).
				Create(&tempBFRelationObj).Error
			if err != nil {
				panic(err)
			}
		}
	}

	//插入新的关系
	for _, newIdIndexBlockItem := range newIdIndexBlockList {
		for j, oldRelation := range newIdIndexBlockItem.BlockRelaList {
			var newRelation model.TempBlockRelation
			newRelation.TempBId = newIdIndexBlockItem.NewId
			newRelation.PreBId = newIdIndexBlockList[newIdIndexBlockItem.OldIndex[j]].NewId
			newRelation.PreBName = oldRelation.PreBName
			newRelation.PreBCode = oldRelation.PreBCode
			newRelation.MyType = oldRelation.MyType
			newRelation.UpdatedAt = time.Now()
			newRelation.UpdatedAt = time.Now()
			//err = global.G_DB.Model(model.TempBlockRelation{}).Create(&newRelation).Error
			err = tx.Model(model.TempBlockRelation{}).Create(&newRelation).Error
			if err != nil {
				panic(err)
			}
		}
	}

	return tx.Commit().Error, sysProTemplate.ID
}

func ReNewIdIndexBlock(newIdIndexBlock *response.NewIdIndexBlock, bTemp *response.SysProTempB, sysBlockWithRelation []response.SysProTempB, myType int) {
	var listArr []model.TempBlockRelation

	switch myType {
	case 0:
		listArr = bTemp.OneRelations
	case 1:
		listArr = bTemp.TwoRelations
	case 2:
		listArr = bTemp.ThreeRelations
	}
	for _, oneObj := range listArr {
		newIdIndexBlock.BlockRelaList = append(newIdIndexBlock.BlockRelaList, oneObj)
		for j, relationObj := range sysBlockWithRelation {
			if relationObj.ID == oneObj.PreBId {
				newIdIndexBlock.OldIndex = append(newIdIndexBlock.OldIndex, j)
				break
			}
		}

	}
}

// @title    SysProject
// @description   AddFields, 新增项目字段配置
// @auth                     （2020年10月20日14:58:59）
// @param     sysProField      model.SysProField
// @return    err             error

func AddFields(sysProField model.SysProField) (err error) {
	var reSysProField model.SysProField
	var fieldsRule model2.SysFieldRule
	// 判断是否有该项目
	tx := global.GDb.Begin()

	err1 := tx.Where("(name = ? or code = ? ) AND pro_id = ?", sysProField.Name, sysProField.Code, sysProField.ProId).First(&reSysProField).Error
	if !errors.Is(err1, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return errors.New("字段名称重复")
	} else {
		var myOrder int
		tx.Raw("select COALESCE(max(my_order),0)  from sys_pro_fields where pro_id = ?", sysProField.ProId).Scan(&myOrder)
		sysProField.MyOrder = myOrder + 1
		err = tx.Create(&sysProField).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Where("(name = ? or code = ? ) AND pro_id = ?", sysProField.Name, sysProField.Code, sysProField.ProId).
			First(&reSysProField).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		fieldsRule.ProId = reSysProField.ProId
		fieldsRule.SysFieldConfId = reSysProField.ID
		fieldsRule.SysFieldName = reSysProField.Name
		fieldsRule.SysFieldCode = reSysProField.Code
		fieldsRule.InputRule = "无"
		fieldsRule.RulePicture = []string{}
		err = tx.Create(&fieldsRule).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return err
}

// @title    SysProject
// @description   GetSysFieldsByPage, 获取项目字段配置分页
// @auth      星期一          （2020年11月04日10:09:47）
// @param     info            request.SysProjectRecordSearch
// @return    err error, list interface{}, total int64

func GetSysFieldsByPage(info requestSysBase.SysProFieldsSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.PageIndex - 1)
	// 创建db
	db := global.GDb.Model(&model.SysProField{}).Where("pro_id = ?", info.ProId)
	var sysProFieldList []model.SysProField
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.InputProcess != 0 {
		db = db.Where("input_process = ?", info.InputProcess)
	}
	err = db.Count(&total).Error
	err = db.Order("code").Limit(limit).Offset(offset).Find(&sysProFieldList).Error
	return err, sysProFieldList, total
}

// @title    SysProject
// @description   UpdateSysFieldsById, 更新字段
// @auth      星期一          （2020年11月04日14:56:19）
// @param    sysProjectRq request.SysProjectUpdateRecord
// @return   rows int64

func UpdateSysFieldsById(sysProField model.SysProField) (reErr error) {
	// 创建db
	sysProField.UpdatedAt = time.Now()
	tx := global.GDb.Begin()
	err := tx.Model(&sysProField).
		Where("id = ?", sysProField.ID).
		Updates(map[string]interface{}{
			"name":            sysProField.Name,
			"my_order":        sysProField.MyOrder,
			"fix_value":       sysProField.FixValue,
			"spec_char":       sysProField.SpecChar,
			"default_val":     sysProField.DefaultVal,
			"check_date":      sysProField.CheckDate,
			"val_change":      sysProField.ValChange,
			"question_change": sysProField.QuestionChange,
			"val_insert":      sysProField.ValInsert,
			"ignore_if":       sysProField.IgnoreIf,
			"prompt":          sysProField.Prompt,
			"max_len":         sysProField.MaxLen,
			"min_len":         sysProField.MinLen,
			"fix_len":         sysProField.FixLen,
			"max_val":         sysProField.MaxVal,
			"min_val":         sysProField.MinVal,
			"validations":     sysProField.Validations,
			"input_process":   sysProField.InputProcess,
		}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(&model.TempBFRelation{}).
		Where("f_id = ?", sysProField.ID).Omit("updated_at").
		Updates(map[string]interface{}{
			"f_name": sysProField.Name,
		}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// @title    SysProject
// @description   RmSysProjectById, 删除字段
// @auth      星期一          （2020年11月04日15:10:48）
// @param    reqId string
// @return   rows int64

func RmSysProjectById(idsIntReq requestSysBase.ReqIds) (rows int64) {
	tx := global.GDb.Begin()
	err := tx.Where("sys_field_conf_id in ? ", idsIntReq.Ids).Delete(&model2.SysFieldRule{}).Error
	if err != nil {
		tx.Rollback()
		return 0
	}

	rows = tx.Delete(&[]model.SysProField{}, idsIntReq.Ids).RowsAffected
	if rows < 1 {
		tx.Rollback()
		return 0
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return 0
	}
	return rows
}

// @title    SysProject
// @description   DelBlockByIds, 删除分块
// @auth      星期一          （2020年11月04日15:10:48）
// @param    reqId string
// @return   rows int64

func DelBlockByIds(idsIntReq requestSysBase.ReqIds) (rows int64) {
	tx := global.GDb.Begin()
	err := tx.Model(&model2.TeachVideo{}).Where("sys_block_conf_id in ? ", idsIntReq.Ids).Delete(&model2.TeachVideo{}).Error
	if err != nil {
		return 0
	}
	rows = tx.Delete(&[]model.SysProTempB{}, idsIntReq.Ids).RowsAffected
	if rows < 1 {
		tx.Rollback()
		return 0
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return 0
	}
	return rows
}

// @title    SysProject
// @description   GetSysFieldsList, 根据项目名称获取所有字段
// @auth      星期一          （2020年11月04日15:31:59）
// @param     proName 			string
// @return    err error, list interface{}

func GetSysFieldsList(proId string) (err error, list interface{}) {
	var sysProFieldList []model.SysProField
	err = global.GDb.Select("name, id, code").
		Where("pro_id = ?", proId).
		Find(&sysProFieldList).Error
	return err, sysProFieldList
}

// @title    SysProject
// @description   	RmSysFieldsById, 删除所有然后插入新的分块字段关系
// @auth      星期一          （2020年11月05日16:36:07）
// @param     sysBFRelations request.SysBFRelations
// @return    err , []model.TempBFRelation

func RmSysFieldsById(sysBFRelations requestProConf.SysBFRelations) (err error, reTempBFRelationList []model.TempBFRelation) {
	tx := global.GDb.Begin()
	var reTempBFRelation []model.TempBFRelation

	if err := tx.Where("temp_b_id = ?", sysBFRelations.BlockId).Delete(&model.TempBFRelation{}).Error; err != nil {
		tx.Rollback()
		return err, reTempBFRelation
	}
	if len(sysBFRelations.TempBFRelationArr) > 0 {
		for i, _ := range sysBFRelations.TempBFRelationArr {
			sysBFRelations.TempBFRelationArr[i].UpdatedAt = time.Now().Add(time.Millisecond * time.Duration(i))
		}
		err = tx.Create(&sysBFRelations.TempBFRelationArr).Error
		if err != nil {
			tx.Rollback()
			return err, reTempBFRelation
		}
	}

	return tx.Commit().Error, sysBFRelations.TempBFRelationArr
}

// @title    SysProject
// @description   	AddExports, 增加新的节点
// @auth      星期一          （2020年11月12日10:12:48）
// @param     sysExportNode model.SysExportNode
// @return    err

func AddExportNode(sysExportNode model.SysExportNode) (err error) {
	if sysExportNode.ExportId != "" {
		var myOrder int
		err = global.GDb.Raw("select COALESCE(max(my_order),0)  from sys_export_nodes where export_id = ?", sysExportNode.ExportId).Scan(&myOrder).Error
		if err != nil {
			return err
		}
		sysExportNode.MyOrder = myOrder + 1
	}
	err = global.GDb.Model(&model.SysExportNode{}).Create(&sysExportNode).Error
	return err
}

// @title    SysProject
// @description   	UpdateExportNodes, 根据id更新节点
// @auth      星期一          （2020年11月12日10:23:42）
// @param     sysExportNode model.SysExportNode
// @return    row int64

func UpdateExportNodes(sysExportNode model.SysExportNode) (err error) {
	tx := global.GDb.Begin()
	var o int
	tx.Model(&sysExportNode).Where("id = ?", sysExportNode.ID).Select("my_order").First(&o)
	global.GLog.Info("my_order:::")
	fmt.Println(o)
	fmt.Println(sysExportNode.MyOrder)
	if o != sysExportNode.MyOrder {
		var orderTemp int
		if o < sysExportNode.MyOrder {
			orderTemp = o
		} else {
			orderTemp = sysExportNode.MyOrder
		}
		err = tx.Model(&model.SysExportNode{}).
			Where("my_order >= ? and export_id = ?", orderTemp, sysExportNode.ExportId).
			Update("my_order", gorm.Expr("my_order + ?", 1)).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Model(sysExportNode).Updates(map[string]interface{}{
		"name":              sysExportNode.Name,
		"one_fields":        sysExportNode.OneFields,
		"two_fields":        sysExportNode.TwoFields,
		"three_fields":      sysExportNode.ThreeFields,
		"fixed_value":       sysExportNode.FixedValue,
		"remark":            sysExportNode.Remark,
		"my_type":           sysExportNode.MyType,
		"my_order":          sysExportNode.MyOrder,
		"one_fields_name":   sysExportNode.OneFieldsName,
		"two_fields_name":   sysExportNode.TwoFieldsName,
		"three_fields_name": sysExportNode.ThreeFieldsName,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func SearchTempApp(t model.TempConfig) (ip string, pwd string, shell string, err error) {
	var temp model.TempConfig
	err = global.GDb.Where("temp_name=?", t.TempName).First(&temp).Error
	ip = temp.TempUrl
	pwd = temp.TempPassword
	shell = temp.TempApp
	return ip, pwd, shell, err
}

func SearchTempMain(t model.TempConfig) (ip string, pwd string, shell string, err error) {
	var temp model.TempConfig
	err = global.GDb.Where("temp_name=?", t.TempName).First(&temp).Error
	ip = temp.TempUrl
	pwd = temp.TempPassword
	shell = temp.TempMain
	return ip, pwd, shell, err
}

/**
 * @title SysProject
 * @description AddInspection,增加审核配置
 * @Auth xingqiyi （2020/11/23 下午5:36）
 * @param     sysInspection model.SysInspection
 * @return err
 */

func AddInspection(sysInspection model.SysInspection) (err error) {
	err = global.GDb.Model(&model.SysInspection{}).Create(&sysInspection).Error
	return err
}

/**
 * @title SysProject
 * @description AddInspection,根据id更新审核配置
 * @Auth xingqiyi （2020/11/23 下午5:51）
 * @param     sysInspection model.SysInspection
 * @return err
 */

func UpdateInspection(sysInspection model.SysInspection) (row int64) {
	row = global.GDb.Model(sysInspection).
		Updates(map[string]interface{}{
			"my_order":      sysInspection.MyOrder,
			"xml_node_name": sysInspection.XmlNodeName,
			"xml_node_code": sysInspection.XmlNodeCode,
			"only_input":    sysInspection.OnlyInput,
			"not_input":     sysInspection.NotInput,
			"max_len":       sysInspection.MaxLen,
			"min_len":       sysInspection.MinLen,
			"max_val":       sysInspection.MaxVal,
			"min_val":       sysInspection.MinVal,
			"validation":    sysInspection.Validation,
		}).RowsAffected
	return row
}

// @title    SysProject
// @description   	UpdateExports, 根据id更新导出配置
// @auth      星期一          （2020年11月12日10:23:42）
// @param     sysExport model.SysExport
// @return    row int64

func UpdateExport(sysExport model.SysExport) (row int64) {
	row = global.GDb.Model(&model.SysExport{}).Where("pro_id = ?", sysExport.ProId).Omit("created_at").Updates(sysExport).RowsAffected
	return row
}

/**
 * @title SysProject
 * @description 批量删除Inspection,删除审核配置
 * @Auth xingqiyi （2020/11/23 下午8:25）
 * @param     ids request.ReqIds
 * @return err
 */

func RmSysInspectionByIds(ids requestSysBase.ReqIds) (reRows int64) {
	rows := global.GDb.Delete(&[]model.SysInspection{}, ids.Ids).RowsAffected
	return rows
}

// @title    SysProject
// @description   	AddExports, 增加新的节点
// @auth      星期一          （2020年11月12日10:12:48）
// @param     sysExport model.SysExport
// @return    err

func AddExport(sysExport model.SysExport) (err error) {
	err = global.GDb.Model(&model.SysExport{}).Create(&sysExport).Error
	return err
}

/**
 * @title SysProject
 * @description AddInspection,增加审核配置
 * @Auth xingqiyi （2020/12/22 下午3:39）
 * @param     sysInspection model.SysInspection
 * @return err
 */

func GetExportAndNodesByProId(search requestSysBase.SysExportNodesSearch) (err error, sysExport response.SysExport, total int64) {
	//var sysExport model.SysExport
	err = global.GDb.Model(&model.SysExport{}).Where("pro_id = ?", search.ProId).Order("updated_at desc").First(&sysExport).Error

	limit := search.PageSize
	offset := search.PageSize * (search.PageIndex - 1)
	// 创建db
	db := global.GDb.Model(&model.SysExportNode{}).Where("export_id = ?", sysExport.ID)
	var sysExportNodes []response.SysExportNode
	// 如果有条件搜索 下方会自动创建搜索语句
	if search.Name != "" {
		db = db.Where("name LIKE ?", "%"+search.Name+"%")
	}

	if search.FieldLike != "" {
		db = db.Where("one_fields LIKE ? or two_fields LIKE ? or three_fields LIKE ? or one_fields_name LIKE ? or two_fields_name LIKE ? or three_fields_name LIKE ? ",
			"%"+search.FieldLike+"%",
			"%"+search.FieldLike+"%",
			"%"+search.FieldLike+"%",
			"%"+search.FieldLike+"%",
			"%"+search.FieldLike+"%",
			"%"+search.FieldLike+"%")
	}
	err = db.Count(&total).Error
	err = db.Order("my_order asc").Limit(limit).Offset(offset).Find(&sysExportNodes).Error
	//Preload("OneSysProField").
	//Preload("TwoSysProField").
	//Preload("ThreeSysProField").Find(&sysExportNodes).Error
	sysExport.NodeList = sysExportNodes
	return err, sysExport, total
}

/**
 * @title SysProject
 * @description RmSysExportNodeByIds,根据id删除
 * @Auth xingqiyi （2020/12/23 上午11:58）
 * @param     ids request.ReqIds
 * @return err
 */

func RmSysExportNodeByIds(ids requestSysBase.ReqIds) (err error, rows int64) {
	tx := global.GDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var nodes []model.SysExportNode
	tx.Model(&model.SysExportNode{}).Where("id in (?)", ids.Ids).Find(&nodes)
	rows = tx.Delete(&[]model.SysExportNode{}, ids.Ids).RowsAffected
	global.GLog.Error("删除条数", zap.Int64("", rows))
	for _, node := range nodes {
		tx.Model(&model.SysExportNode{}).
			Where("export_id = ? and my_order > ? ", node.ExportId, node.MyOrder).
			Update("my_order", gorm.Expr("my_order - ?", 1))
	}
	return tx.Commit().Error, rows
}

/**
 * @title SysProject
 * @description InsertTo,将序号插到序号前
 * @Auth xingqiyi （2020/12/25 下午1:53）
 * @param     insertTo request.InsertTo
 * @return err error
 */

func ChangeOrder(insertTo requestProConf.ChangeOrder) (err error) {
	tx := global.GDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err = errors.New("传入了同一个order")
	if insertTo.StartOrder < insertTo.EndOrder {
		//向下插入

		//将开始到结束的order -1
		err = tx.Model(&model.SysExportNode{}).
			Where("export_id = ? and my_order > ? and my_order <= ?", insertTo.ExportId, insertTo.StartOrder, insertTo.EndOrder).
			Update("my_order", gorm.Expr("my_order - ?", 1)).
			Error
	} else if insertTo.StartOrder > insertTo.EndOrder {
		//向上插入

		//将开始到结束的order -1
		err = tx.Model(&model.SysExportNode{}).
			Where("export_id = ? and my_order >= ? and my_order < ?", insertTo.ExportId, insertTo.EndOrder, insertTo.StartOrder).
			Update("my_order", gorm.Expr("my_order + ?", 1)).
			Error
	} else {
		return err
	}

	//将开始移动的（StartOrder）更改order为 endOrder
	err = tx.Model(&model.SysExportNode{}).
		Where("export_id = ? and id = ? ", insertTo.ExportId, insertTo.StartId).
		Update("my_order", insertTo.EndOrder).
		Error

	return tx.Commit().Error
}

// BlockChangeOrder 分块交换顺序
func BlockChangeOrder(insertTo requestProConf.SwitchOrder) (err error) {
	tx := global.GDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err = tx.Model(&model.SysProTempB{}).
		Where("id = ? ", insertTo.StartId).
		Update("my_order", insertTo.EndOrder).
		Error

	err = tx.Model(&model.SysProTempB{}).
		Where("id = ? ", insertTo.EndId).
		Update("my_order", insertTo.StartOrder).
		Error

	return tx.Commit().Error
}

/**
 * @title SysProject
 * @description ExportNodeByExportId,根据exportId导出到Excel
 * @Auth xingqiyi （2020/12/28 下午5:35）
 * @param     exportId string
 * @return err
 */

func ExportNodeByExportId(exportId string) (sysExportNodeList []model.SysExportNode, err error) {
	err = global.GDb.Model(&model.SysExportNode{}).Where("export_id = ?", exportId).
		Order("my_order").Find(&sysExportNodeList).Error
	return sysExportNodeList, err
}

/**
 * @title SysProject
 * @description GetExportById,根据id获取export
 * @Auth xingqiyi （2020/12/29 上午11:57）
 * @param     exportId string
 * @return err
 */

func GetExportById(exportId string) (sysExport model.SysExport, err error) {
	err = global.GDb.Model(&model.SysExport{}).Where("id = ?", exportId).First(&sysExport).Error
	return sysExport, err
}

/**
 * @title SysProject
 * @description UpdateExportXmlValById,根据id更新xml导出模板
 * @Auth xingqiyi （2020/12/29 下午2:11）
 * @param     exportId string
 * @param     xmlVal string
 * @return err
 */

func UpdateExportXmlValById(exportId string, xmlVal string) (err error) {
	err = global.GDb.Model(&model.SysExport{}).
		Where("id = ? ", exportId).
		Update("temp_val", xmlVal).
		Error
	return err
}

/**
 * @title SysProject
 * @description GetExportNodeById,根据id获取exportNode
 * @Auth xingqiyi （2020/12/29 下午2:25）
 * @param     exportNodeId string
 * @return err
 */

func GetExportNodeById(exportNodeId string) (err error, sysExportNode model.SysExportNode) {
	err = global.GDb.Model(&model.SysExportNode{}).
		Where("id = ? ", exportNodeId).
		First(&sysExportNode).
		Error
	return err, sysExportNode
}

/**
 * @title SysProject
 * @description GetTempBFRelationByBId,根据blockid获取关系
 * @Auth xingqiyi （2020/12/30 下午1:57）
 * @param     bId string
 * @return err
 */

func GetTempBFRelationByBId(bId string) (err error, tempBFRelationList []model.TempBFRelation) {
	err = global.GDb.Model(&model.TempBFRelation{}).
		Where("temp_b_id = ? ", bId).
		Order("updated_at asc").
		Find(&tempBFRelationList).
		Error
	return err, tempBFRelationList
}

/**
 * @title SysProject
 * @description AddInspection,增加质检配置
 * @Auth xingqiyi （2021/1/4 下午5:28）
 * @param     quality model.SysQuality
 * @return err
 */

func AddSysQuality(quality model.SysQuality) (err error) {
	var count int64
	global.GDb.Model(&model.SysQuality{}).
		Where("pro_id = ? and parent_xml_node_name = ? and xml_node_name = ?", quality.ProId, quality.ParentXmlNodeName, quality.XmlNodeName).
		Count(&count)
	if count <= 0 {
		return global.GDb.Model(&model.SysQuality{}).Create(&quality).Error
	}
	return errors.New("已存在一条一样的记录")
}

/**
 * @title SysProject
 * @description RmAddSysQualityByIds,批量删除质检配置
 * @Auth xingqiyi （2021/1/5 上午9:41）
 * @param     sysInspection model.SysInspection
 * @return err
 */

func RmAddSysQualityByIds(ids requestSysBase.ReqIds) (reRows int64) {
	rows := global.GDb.Delete(&[]model.SysQuality{}, ids.Ids).RowsAffected
	return rows
}

/**
 * @title SysProject
 * @description UpdateSysQuality,更新质检配置
 * @Auth xingqiyi （2021/1/5 上午11:07）
 * @param     quality model.SysQuality
 * @return row int64
 */

func UpdateSysQuality(quality model.SysQuality) (row int64) {
	//var count int64
	var arr []model.SysQuality
	global.GDb.Model(&model.SysQuality{}).
		Where("pro_id = ? and parent_xml_node_name = ? and xml_node_name = ?", quality.ProId, quality.ParentXmlNodeName, quality.XmlNodeName).
		Find(&arr)
	if len(arr) == 1 {
		if arr[0].ID != quality.ID {
			return 0
		}
	}
	return global.GDb.Model(&model.SysQuality{}).
		Where("pro_id = ? and id = ?", quality.ProId, quality.ID).
		Omit("created_at").
		Omit("created_by").
		Updates(quality).RowsAffected
	//return errors.New("已存在一条一样的记录")
}

/**
 * @title SysProject
 * @description GetSysQualityByPage,分页查询质检配置
 * @Auth xingqiyi （2021/1/5 上午11:29）
 * @param     info request.SearchSysQuality
 * @return err error, list interface{}, total int64
 */

func GetSysQualityByPage(info requestSysBase.SearchSysQuality) (err error, list interface{}, total int64, maxOrder int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.PageIndex - 1)
	// 创建查询db
	db := global.GDb.Model(&model.SysQuality{}).Where("pro_id = ?", info.ProId)
	var sysQualityList []model.SysQuality
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.ParentXmlNodeName != "" {
		db = db.Where("parent_xml_node_name = ?", info.ParentXmlNodeName)
	}
	if info.XmlNodeName != "" {
		db = db.Where("xml_node_name = ?", info.XmlNodeName)
	}
	if info.FieldName != "" {
		db = db.Where("field_name = ?", info.FieldName)
	}
	if info.BelongType > 0 {
		db = db.Where("belong_type = ?", info.BelongType)
	}
	err = db.Count(&total).Error
	err = db.Order("my_order").Limit(limit).Offset(offset).Find(&sysQualityList).Error
	db.Raw("select MAX(my_order) from sys_qualities where pro_id = ?", info.ProId).Scan(&maxOrder)
	return err, sysQualityList, total, maxOrder + 1
}

/**
 * @title SysProject
 * @description GetMaxCodeSysFieldsByProName,获取项目字段code最大值
 * @Auth xingqiyi （2021/1/6 上午11:10）
 * @param     proId string
 * @return err
 */

func GetMaxCodeSysFieldsByProName(proId string) (err error, maxCode string) {
	var code string
	var count int64
	global.GDb.Model(model.SysProField{}).Where("pro_id = ?", proId).Count(&count)
	if count == 0 {
		code = "fc000"
	} else {
		err = global.GDb.Raw("select MAX(code) from sys_pro_fields where pro_id = ?", proId).Scan(&code).Error
	}
	return err, code
}

/**
 * @title SysProject
 * @description GetGroupByProId,根据项目id获取所有分类
 * @Auth xingqiyi （2021/1/7 上午10:49）
 * @param     proId string
 * @return err
 */

func GetGroupByProId(proId int) (err error, myType []int) {
	err = global.GDb.Raw("select belong_type from sys_qualities where pro_id = ? group by belong_type", proId).Scan(&myType).Error
	return err, myType
}

/**
 * @title SysProject
 * @description GetByTypeAndProId,根据项目id和类型获取质检配置
 * @Auth xingqiyi （2021/1/7 上午10:49）
 * @param     proId string
 * @return err
 */

func GetByTypeAndProId(item int, proId string) (err error, sysQualityList []model.SysQuality) {
	err = global.GDb.Raw("select * from sys_qualities where pro_id = ?  and belong_type = ?", proId, item).Scan(&sysQualityList).Error
	return err, sysQualityList
}

/**
 * @title SysProject
 * @description GetByTypeAndProId,根据项目id和类型获取质检配置
 * @Auth xingqiyi （2021/1/7 上午10:49）
 * @param     proId string
 * @return err
 */

func GetByTypeAndProIdAndIds(item []int, proId int) (err error, sysQualityList []model.SysQuality) {
	err = global.GDb.Raw("select * from sys_qualities where pro_id = ?  and belong_type in ", proId, item).Scan(&sysQualityList).Error
	return err, sysQualityList
}

/**
 * @title SysProject
 * @description GetSysInspectionByPage,分页获取审核配置
 * @Auth xingqiyi （2021/1/5 上午11:29）
 * @param     info request.SearchSysInspection
 * @return err error, list interface{}, total int64
 */

func GetSysInspectionByPage(info requestSysBase.SearchSysInspection) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.PageIndex - 1)
	// 创建查询db
	db := global.GDb.Model(&model.SysInspection{}).Where("pro_id = ?", info.ProId)
	if info.XmlNodeCode != "" {
		db = db.Where("xml_node_code Like ?", "%"+info.XmlNodeCode+"%")
	}
	if info.XmlNodeName != "" {
		db = db.Where("xml_node_name Like ?", "%"+info.XmlNodeName+"%")
	}
	var sysInspection []model.SysInspection
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	err = db.Order("created_at desc").Limit(limit).Offset(offset).Find(&sysInspection).Error
	return err, sysInspection, total
}

// GetSysQualityFormat 获取格式化的质检配置
func GetSysQualityFormat(proId string) (err error, o []model.SysQuality) {
	err = global.GDb.Model(&model.SysQuality{}).Where("pro_id = ?", proId).Find(&o).Error
	return err, o
}

func RefreshIndex() {
	var list []model.SysExportNode
	global.GDb.Model(&model.SysExportNode{}).Where("export_id = '907280101312299008'").Order("my_order asc").Find(&list)
	tx := global.GDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	for i, n := range list {
		err := tx.Model(&model.SysExportNode{}).Where("id = ?", n.ID).Updates(map[string]interface{}{
			"my_order": i + 1,
		}).Error
		if err != nil {
			tx.Rollback()
		}
	}
	tx.Commit()
}
