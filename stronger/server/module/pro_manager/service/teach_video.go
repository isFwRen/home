package service

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/lib/pq"
	"io"
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

func BlockSync(proCode string) (err error) {
	var proinformation model.SysProject
	err = global.GDb.Model(&model.SysProject{}).Where("code = ? ", proCode).Find(&proinformation).Error
	if err != nil {
		return err
	}
	var spt []model.SysProTemplate
	err = global.GDb.Model(&model.SysProTemplate{}).Where("pro_id = ? ", proinformation.ID).Find(&spt).Error
	if err != nil {
		return err
	}
	idsSpt := make([]string, 0)
	for _, v := range spt {
		idsSpt = append(idsSpt, v.ID)
	}
	var sptb []model.SysProTempB
	err = global.GDb.Model(&model.SysProTempB{}).Where("pro_temp_id in ? ", idsSpt).Find(&sptb).Error
	if err != nil {
		return err
	}
	err = global.GDb.Model(&model2.TeachVideo{}).Where("pro_id = ? ", proinformation.ID).Delete(&model2.TeachVideo{}).Error
	if err != nil {
		return err
	}
	var teachVideos []model2.TeachVideo
	for _, v := range sptb {
		var item model2.TeachVideo
		item.SysBlockConfId = v.ID
		item.SysBlockCode = v.Code
		item.SysBlockName = v.Name
		item.ProId = proinformation.ID
		item.InputRule = "无"
		item.Video = pq.StringArray{}
		teachVideos = append(teachVideos, item)
	}
	err = global.GDb.Model(&model2.TeachVideo{}).Create(teachVideos).Error
	return
}

func GetTeachVideoList(info request.TV) (err error, list []map[string]interface{}, total int64) {
	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)

	var proInformation model.SysProject
	err = global.GDb.Model(&model.SysProject{}).Where("code = ? ", info.ProCode).Find(&proInformation).Error
	if err != nil {
		return
	}

	db := global.GDb.Model(&model2.TeachVideo{})
	if info.BlockName != "" {
		db = db.Where("sys_block_name LIKE ? ", "%"+info.BlockName+"%")
	}
	if info.Rule != "" {
		db = db.Where("input_rule = ? ", info.Rule)
	}
	db = db.Where("pro_id = ? ", proInformation.ID)
	var tv []model2.TeachVideo
	var tvReturn []model2.TeachVideos
	err = db.Count(&total).Limit(limit).Offset(offset).Where("array_length(video,1) > 0").Debug().Find(&tv).Error
	err = copier.Copy(&tvReturn, &tv)
	if err != nil {
		return err, nil, 0
	}
	for _, v := range tv {
		if len(v.Video) == 0 {
			continue
		}
		for j, v1 := range tvReturn {
			if v.ID == v1.ID {
				//fmt.Println(strings.Split(v.Video[0], "/")[len(strings.Split(v.Video[0], "/"))-1])
				for _, vid := range v.Video {
					videoName := strings.Split(vid, "/")[len(strings.Split(v.Video[0], "/"))-1]
					vs := make(map[string]string, 0)
					vs["name"] = videoName
					vs["path"] = vid
					tvReturn[j].Video = append(tvReturn[j].Video, vs)
				}
			}
		}
	}

	list = []map[string]interface{}{}
	for _, rule := range tvReturn {
		list = append(list, utils.Struct2Map(rule))
	}
	return err, list, total
}
func GetTeachVideoListCM(info request.TV) (err error, list interface{}, total int64) {
	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)

	var proInformation model.SysProject
	err = global.GDb.Model(&model.SysProject{}).Where("code = ? ", info.ProCode).Find(&proInformation).Error
	if err != nil {
		return
	}

	db := global.GDb.Model(&model2.TeachVideo{})
	if info.BlockName != "" {
		db = db.Where("sys_block_name LIKE ? ", "%"+info.BlockName+"%")
	}
	if info.Rule != "" {
		db = db.Where("input_rule = ? ", info.Rule)
	}
	db = db.Where("pro_id = ? ", proInformation.ID)
	var tv []model2.TeachVideo
	var tvReturn []model2.TeachVideos
	err = db.Count(&total).Limit(limit).Offset(offset).Debug().Find(&tv).Error
	err = copier.Copy(&tvReturn, &tv)
	if err != nil {
		return err, nil, 0
	}
	for _, v := range tv {
		if len(v.Video) == 0 {
			continue
		}
		for j, v1 := range tvReturn {
			if v.ID == v1.ID {
				//fmt.Println(strings.Split(v.Video[0], "/")[len(strings.Split(v.Video[0], "/"))-1])
				for _, vid := range v.Video {
					videoName := strings.Split(vid, "/")[len(strings.Split(v.Video[0], "/"))-1]
					vs := make(map[string]string, 0)
					vs["name"] = videoName
					vs["path"] = vid
					tvReturn[j].Video = append(tvReturn[j].Video, vs)
				}
			}
		}
	}
	return err, tvReturn, total
}

func EditTeachVideo(newTv model2.TeachVideo) error {
	return global.GDb.Model(&model2.TeachVideo{}).
		Where("id = ? ", newTv.ID).Updates(&newTv).Error
}

func DeleteTeachVideo(ids []string) error {
	for _, id := range ids {
		var rule model2.TeachVideo
		err := global.GDb.Model(&model2.TeachVideo{}).Where("id = ? ", id).Find(&rule).Error
		if err != nil {
			return err
		}
		for _, pic := range rule.Video {
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
		rule.Video = []string{}
		rule.InputRule = "无"
		err = global.GDb.Model(&model2.TeachVideo{}).Where("id = ? ", rule.ID).Updates(&rule).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func ExportTeachVideo(info request.ExportTV) (err error, arr []string) {
	var proInformation model.SysProject
	err = global.GDb.Model(&model.SysProject{}).Where("code = ? ", info.ProCode).Find(&proInformation).Error
	if err != nil {
		return err, nil
	}
	var tvs []model2.TeachVideo
	err = global.GDb.Model(model2.TeachVideo{}).Where("pro_id = ? ", proInformation.ID).Find(&tvs).Error
	if err != nil {
		return err, nil
	}
	basicPath := global.GConfig.LocalUpload.FilePath + "教学视频导出/"
	bookName := info.ProCode + "-教学视频" + ".xlsx"
	err = utils.ExportBigExcel(basicPath, bookName, "", tvs)
	if err != nil {
		return err, nil
	}

	xlsx, err := excelize.OpenFile(basicPath + bookName)
	if err != nil {
		return err, nil
	}
	for i, _ := range tvs {
		xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), " ")
		xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), i+1)
	}
	if err = xlsx.SaveAs(basicPath + bookName); err != nil {
		fmt.Println(err)
	}

	//打包视频文件
	zp, err := os.Create(global.GConfig.LocalUpload.FilePath + "教学视频导出/" + info.ProCode + "-教学视频" + ".zip")
	if err != nil {
		return nil, nil
	}
	defer zp.Close()
	zipWriter := zip.NewWriter(zp)
	for _, v := range tvs {
		if len(v.Video) == 0 {
			continue
		}
		for _, v1 := range v.Video {
			_, err := os.Stat(v1)
			if os.IsNotExist(err) {
				continue
			}
			f, err := os.Open(v1)
			if err != nil {
				return nil, nil
			}
			//fmt.Println(v.Video[0])
			//fmt.Println(strings.Replace(global.GConfig.LocalUpload.FilePath+"教学视频/"+proInformation.Code+"/", "./", "", -1))
			//name := strings.Replace(v.Video[0], global.GConfig.LocalUpload.FilePath+"/教学视频/"+proInformation.Code+"/", "", -1)
			//fmt.Println("1", name)
			name := strings.Replace(v1, strings.Replace(global.GConfig.LocalUpload.FilePath+"教学视频/"+proInformation.Code+"/", "./", "", -1), "", -1)
			//fmt.Println("2", name)
			w1, err := zipWriter.Create(name)
			if err != nil {
				panic(err)
			}
			if _, err := io.Copy(w1, f); err != nil {
				panic(err)
			}
		}
	}
	zipWriter.Close()
	arr = append(arr, "files/教学视频导出/"+bookName)
	arr = append(arr, "files/教学视频导出/"+info.ProCode+"-教学视频"+".zip")
	return
}

func UploadTeachVideo(info []*multipart.FileHeader, proCode []string, c *gin.Context) error {
	failmsg := ""
	for _, file := range info {
		filename := file.Filename

		// 线上 uploads/file
		basicPath := global.GConfig.LocalUpload.FilePath + "/教学视频/" + proCode[0] + "/"
		// 设置文件需要保存的指定位置并设置保存的文件名字
		dst := path.Join(basicPath, filename)
		fmt.Println(filename + " has been saved in this path " + basicPath)
		fmt.Println("dst", dst)

		var syn model2.TeachVideo
		var total int64
		blockName := strings.Split(filename, ".")[0]
		err := global.GDb.Model(&model2.TeachVideo{}).Where("sys_block_name = ? ", blockName).
			Find(&syn).Count(&total).Error
		if err != nil {
			return err
		}
		if total == 0 {
			failmsg += filename + "没有对应的分块教学配置; "
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
		err = global.GDb.Model(&model2.TeachVideo{}).Where("sys_block_name = ? ", blockName).
			Updates(map[string]interface{}{
				"video":      pq.StringArray{dst},
				"input_rule": "有",
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
