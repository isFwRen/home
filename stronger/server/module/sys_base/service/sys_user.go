package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"server/global"
	"server/global/response"
	"server/module/sys_base/model"
	"server/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	uuid "github.com/satori/go.uuid"
)

func Login(u *model.SysUser, isIntranet bool) (err error, userInter *model.SysUser) {
	var user model.SysUser
	u.Password = utils.MD5V([]byte(u.Password))
	if !isIntranet {
		err = global.GDb.Where("phone = ? AND password = ? AND status = ?", u.Phone, u.Password, true).
			Preload("SysRoles", "status = ?", 1).First(&user).Error
		if !user.Status {
			return errors.New("This user is not exists or resignation "), userInter
		}
	} else {
		err = global.GDb.Where("code = ? AND password = ? AND status = ?", u.Code, u.Password, true).
			Preload("SysRoles", "status = ?", 1).First(&user).Error
		if !user.Status {
			return errors.New("This user is not exists or resignation, err: " + err.Error()), userInter
		}
	}
	return err, &user
}

func GetUserInfoList(info model.PageInfo) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.PageIndex - 1)
	db := global.GUserDb.Model(&model.SysUser{})
	var userList []model.SysUser
	err = db.Count(&total).Error
	//err = db.Limit(limit).Offset(offset).Preload("Authority").Find(&userList).Error
	err = db.Limit(limit).Offset(offset).Preload("Authority").Preload("SysRoles").Find(&userList).Error
	return err, userList, total
}

func UploadHeaderImg(uuid uuid.UUID, filePath string) (err error, userInter *model.SysUser) {
	var user model.SysUser
	err = global.GUserDb.Where("uuid = ?", uuid).First(&user).Update("header_img", filePath).First(&user).Error
	return err, &user
}

func GetRedisJWT(userId string) (err error, redisJWT string) {
	redisJWT, err = global.GRedis.Get("user_token:" + userId).Result()
	return err, redisJWT
}

func SetRedisJWT(jwt string, userId string, expiration time.Time) (err error) {
	// 此处过期时间等于jwt过期时间
	//str := strconv.Itoa(userId)
	now := time.Now().Local()
	//expiresDuration := time.Duration(global.GConfig.JWT.ExpiresAt * 1000 * 1000 * 1000)
	//newDay := now.Add(expiresDuration).Format("2006-01-02") + " 02:01:01"
	//expiration, err := time.ParseInLocation("2006-01-02 15:04:05", newDay, time.Local)
	//if err != nil {
	//	return err
	//}
	//global.GLog.Info("登录:token过期时间", zap.Any(newDay, expiration))
	err = global.GRedis.Set("user_token:"+userId, jwt, expiration.Sub(now)).Err()
	return err
}

func DelRedisJWT(userId string) (err error) {
	err = global.GRedis.Del("user_token:" + userId).Err()
	return err
}

func GetRedisSecret(userId string) (s string, err error) {
	s, err = global.GRedis.Get("user_secret:" + userId).Result()
	return s, err
}

func GetRedisSecretWithTTL(userId string) (isEx time.Duration) {
	ttl, _ := global.GRedis.TTL("user_secret:" + userId).Result()
	return ttl
}

func SetRedisSecret(jwt string, userId string) (err error) {
	now := time.Now().Local()
	newDay := now.AddDate(0, 0, 1).Format("2006-01-02") + " 02:01:01"
	expiration, err := time.ParseInLocation("2006-01-02 15:04:05", newDay, time.Local)
	if err != nil {
		return err
	}
	global.GLog.Info("登录", zap.Any("otp秘钥过期时间", expiration))
	err = global.GRedis.Set("user_secret:"+userId, jwt, expiration.Sub(now)).Err()
	return err
}

func DelRedisSecret(userId string) (err error) {
	err = global.GRedis.Del("user_secret:" + userId).Err()
	return err
}

func SetRedisPerm(permMap map[string]interface{}, userId string) (err error) {
	err = global.GRedis.HMSet("user_perms:"+userId, permMap).Err()
	return err
}

func GetRedisPerm(userId string, proCode string, process string) (err error, has []interface{}) {
	has, err = global.GRedis.HMGet("user_perms:"+userId, proCode+"_"+process).Result()
	return err, has
}

// GetUserById 根据id获取用户信息
func GetUserById(id string) (err error, u *model.SysUser) {
	var user model.SysUser
	err = global.GUserDb.Where("id = ?", id).Select("id,nick_name,code").First(&user).Error
	return err, &user
}

// UpdateUserQrcodeById 根据id获取更新用户二维码
func UpdateUserQrcodeById(id, qrcode string) error {
	return global.GUserDb.Model(&model.SysUser{}).Where("id = ?", id).Update("qrcode", qrcode).Error
}

// GetUserByIdStage 根据id和状态获取用户信息
func GetUserByIdStage(id string) (err error, u model.SysUser) {
	err = global.GUserDb.Where("id = ? and status = true", id).
		Select("id,nick_name,code,header_img,phone").First(&u).Error
	return err, u
}

// UserCount 查询是否有该员工
func UserCount(jobNumber, jobName string) (err error, count int64) {
	err = global.GUserDb.Model(&model.SysUser{}).Where("code = ? AND nick_name = ? AND status = 't'", jobNumber, jobName).Count(&count).Error
	return err, count
}

// UploadAvatarUser 修改用户头像信息
func UploadAvatarUser(avatarPath, jobNumber, jobName string) error {
	maps := map[string]interface{}{
		"header_img": avatarPath,
	}
	db := global.GUserDb.Model(&model.SysUser{})
	err := db.Where("code = ? AND nick_name = ? AND status = 't'", jobNumber, jobName).Updates(maps).Error
	return err
}

// UploadUserAvatar 修改用户头像
func UploadUserAvatar(c *gin.Context, file *multipart.FileHeader, jobNumber, jobName, filepath string) (err error, headerPath string) {
	t := time.Now()
	strTime := t.String()
	yearIta := strconv.Itoa(t.Year())
	month := strTime[5:7]
	dayIta := strconv.Itoa(t.Day())
	staffFileName := yearIta + "/" + month + "/" + dayIta + "/" + yearIta + "_" + month + "_" + dayIta + "_" + file.Filename
	file.Filename = staffFileName
	//检查是否存在该员工路径
	_, err = os.Stat(filepath)
	if err != nil && os.IsNotExist(err) {
		//员工路径不存在 - 直接创建 添加
		err = os.MkdirAll(filepath, os.ModePerm)
		if err != nil {
			global.GLog.Error("upload file fail:", zap.Any("err", err))
			return
		}
		dst := path.Join(filepath, staffFileName)
		err = c.SaveUploadedFile(file, dst)
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("上传失败，%v", err), c)
		}
		//用户头像地址
		avatarPath := filepath + staffFileName
		//地址存入员工数据库
		err = UploadAvatarUser(avatarPath, jobNumber, jobName)
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("修改失败，%v", err), c)
		} else {
			c.FileAttachment(avatarPath, jobNumber)

		}
		fmt.Println("---------------------------111-----*avatarPath=", avatarPath)
		return err, avatarPath
	} else {
		//存在 - 删除 添加
		err = os.RemoveAll(filepath)
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("删除失败，%v", err), c)
		}

		// 创建此路径
		err = os.MkdirAll(filepath, os.ModePerm)
		if err != nil {
			global.GLog.Error("upload file fail:", zap.Any("err", err))
			return
		}
		//上传头像
		dst := path.Join(filepath, staffFileName)
		err = c.SaveUploadedFile(file, dst)
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("上传失败，%v", err), c)
		}

		//用户头像地址
		avatarPath := filepath + staffFileName
		//地址存入员工数据库
		err = UploadAvatarUser(avatarPath, jobNumber, jobName)
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("修改失败，%v", err), c)
		} else {
			c.FileAttachment(avatarPath, jobNumber)
		}
		fmt.Println("--------------------------------*avatarPath=", avatarPath)
		return err, avatarPath
	}
	return err, ""
}

func CheckPasswordTime(u *model.SysUser) (err error) {
	passwordDate := u.PasswordDate
	if passwordDate.IsZero() {
		return errors.New("初次记录请修改密码")
	}
	bb := passwordDate.Add(90 * 24 * time.Hour)
	if time.Now().After(bb) {
		return errors.New("您的登陆密码已使用超过三个月，为确保您的账号安全，请修改后登录")
	}
	return nil
}
