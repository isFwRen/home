package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mime/multipart"
	"os"
	"path"
	"server/global"
	"strings"
)

// @title    PathExists
// @description   文件目录是否存在
// @auth                     （2020/04/05  20:22）
// @param     path            string
// @return    err             error

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// @title    createDir
// @description   批量创建文件夹
// @auth                     （2020/04/05  20:22）
// @param     dirs            string
// @return    err             error

func CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist, err := PathExists(v)
		if err != nil {
			return err
		}
		if !exist {
			global.GLog.Debug("create directory ", zap.Any("v", v))
			err = os.MkdirAll(v, os.ModePerm)
			if err != nil {
				//global.G_LOG.Error("create directory", zap.Any("v", v), " error:", err)
				global.GLog.Error("create directory", zap.Any("err", err))
			}
		}
	}
	return err
}

//保存图片文件
func SaveImgFile(c *gin.Context, f *multipart.FileHeader, fileName string, filepath string) (reErr error, rePath string) {
	fileExt := strings.ToLower(path.Ext(f.Filename))
	if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" {
		return errors.New("文件类型不对"), ""
	}
	isExist, _ := PathExists(filepath)
	if !isExist {
		_ = CreateDir(filepath)
	}
	err := c.SaveUploadedFile(f, filepath+fileName)
	return err, filepath + fileName
}
