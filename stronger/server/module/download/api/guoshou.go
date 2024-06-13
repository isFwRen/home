/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/4/12 14:54
 */

package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wxnacy/wgo/file"
	"go.uber.org/zap"
	"net/http"
	"server/global"
	"server/global/response"
	"server/module/download/const_data"
	"server/module/download/model"
	"server/module/download/project/guoshou"
	"server/module/download/service"
	"time"
)

// UploadFile
// @Tags sys-download
// @Summary 下载--下载单据
// @Auth xingqiyi
// @Date 2023年04月12日14:59:29
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ImageFile true "推送实体类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"推送成功"}"
// @Router /v1/gs/upload-file [post]
func UploadFile(c *gin.Context) {
	var imageFile model.ImageFile
	err := c.ShouldBindJSON(&imageFile)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	global.GLog.Info("imageFile", zap.Any("", imageFile))
	err = guoshou.Process(imageFile)
	if err != nil {
		global.GLog.Error("err", zap.Error(err))
		c.JSON(http.StatusOK, map[string]interface{}{
			"resultCode": "500",
			"resultMsg":  "推送失败",
			"claimNo":    imageFile.ClaimNo,
			"claimTpaId": imageFile.ClaimTpaId,
		})
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{
			"resultCode": "0000",
			"resultMsg":  nil,
			"claimNo":    imageFile.ClaimNo,
			"claimTpaId": imageFile.ClaimTpaId,
		})
	}
}

// UploadHospital
// @Tags sys-download
// @Summary 下载--下载医院机构
// @Auth xingqiyi
// @Date 2023年04月12日14:59:29
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Hospital true "推送实体类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"推送成功"}"
// @Router /v1/gs/upload-hospital [post]
func UploadHospital(c *gin.Context) {
	var h model.Hospital
	err := c.ShouldBindJSON(&h)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := json.Marshal(h)
	if err != nil {
		global.GLog.Error("err", zap.Error(err))
		c.JSON(http.StatusOK, map[string]interface{}{
			"resultCode": "500",
			"resultMsg":  "推送失败",
		})
		return
	}
	path := fmt.Sprintf("%v%v/hospital-%v.json",
		global.GConfig.LocalUpload.FilePath+global.PathUpdateConstHospital,
		time.Now().Format("20060102"), time.Now().Format("20060102030405"))

	//存数据库log
	var hospitalLogs []model.UpdateConstLog
	for _, hospital := range h.HospitalInfoList {
		proCode := const_data.Num2code[hospital.BranchNo[:2]]
		if proCode == "" {
			global.GLog.Error("err：proCode，hospital.BranchNo" + proCode + "——————" + hospital.BranchNo)
			continue
		}
		byteData, err := json.Marshal(hospital)
		if err != nil {
			global.GLog.Error("err", zap.Error(err))
			c.JSON(http.StatusOK, map[string]interface{}{
				"resultCode": "500",
				"resultMsg":  "推送失败",
			})
			return
		}
		hospitalLog := model.UpdateConstLog{
			Type:      1,
			Name:      proCode + "_" + global.GProConf[proCode].Name + "_医疗机构" + hospital.BranchNo[:2],
			ProCode:   proCode,
			IsDeleted: model.IsDeletedMap[hospital.IsDeleted],
			LocalUrl:  path,
			OtherInfo: string(byteData),
		}
		hospitalLogs = append(hospitalLogs, hospitalLog)
	}
	_ = service.InsertUpdateConstLogs(hospitalLogs)
	//2023年06月09日10:57:47  没有也要返回接收成功
	//if err != nil {
	//	global.GLog.Error("err", zap.Error(err))
	//	c.JSON(http.StatusOK, map[string]interface{}{
	//		"resultCode": "500",
	//		"resultMsg":  "推送失败",
	//	})
	//}

	//保存json文件
	err = file.SaveFile(path, string(data))
	if err != nil {
		global.GLog.Error("err", zap.Error(err))
		c.JSON(http.StatusOK, map[string]interface{}{
			"resultCode": "500",
			"resultMsg":  "推送失败",
		})
	} else {
		global.GLog.Info("hospital", zap.Any("推送成功", path))
		c.JSON(http.StatusOK, map[string]interface{}{
			"resultCode": "0000",
			"resultMsg":  nil,
		})
	}
}

// UploadCatalogue
// @Tags sys-download
// @Summary 下载--下载医疗目录
// @Auth xingqiyi
// @Date 2023年04月12日14:59:29
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Catalog true "推送实体类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"推送成功"}"
// @Router /v1/gs/upload-catalogue [post]
func UploadCatalogue(c *gin.Context) {
	var catalog model.Catalog
	err := c.ShouldBindJSON(&catalog)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = guoshou.FetchCatalogue(catalog)
	if err != nil {
		global.GLog.Error("err", zap.Error(err))
		c.JSON(http.StatusOK, map[string]interface{}{
			"resultCode": "500",
			"resultMsg":  "推送失败",
		})
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{
			"resultCode": "0000",
			"resultMsg":  nil,
		})
	}
}
