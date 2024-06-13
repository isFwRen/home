package api

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"regexp"
	"server/global"
	"server/global/response"
	"server/middleware"
	service2 "server/module/pro_conf/service"
	"server/module/sys_base"
	"server/module/sys_base/model"
	request2 "server/module/sys_base/model/request"
	sys_base4 "server/module/sys_base/model/response"
	"server/module/sys_base/service"
	"server/utils"
	"time"
	"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/api"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

// TaskLogin
// @Tags SysUser
// @Summary 录入系统--用户登录
// @Produce  application/json
// @Param data body request2.Logins true "用户登录接口"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登陆成功"}"
// @Router /sys-base/sys-login/login-task [post]
func TaskLogin(c *gin.Context) {
	var L request2.Logins
	_ = c.ShouldBindJSON(&L)
	if !L.IsIntranet {
		UserVerify := utils.Rules{
			//"CaptchaId": {utils.NotEmpty()},
			// "Captcha":  {utils.NotEmpty()},
			"Phone":    {utils.NotEmpty()},
			"Password": {utils.NotEmpty()},
		}
		UserVerifyErr := utils.Verify(L, UserVerify)
		if UserVerifyErr != nil {
			response.FailWithMessage(UserVerifyErr.Error(), c)
			return
		}

		err, a := utils.GetRedisCaptcha(L.Phone)
		if err == nil && a == "" {
			response.FailWithMessage("验证码已过期!", c)
			return
		}
		fmt.Println("global.GConfig.System.ProCode", global.GConfig.System.ProCode)
		var UserPermission model.SysProPermission
		var userx []model.SysUser
		err = global.GDb.Model(&model.SysUser{}).Where("phone = ? AND status = 'true'", L.Phone).Find(&userx).Error
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("登录失败"), c)
			return
		}
		if len(userx) > 1 {
			response.FailWithMessage(fmt.Sprintf("登录失败, 存在同工号在职的员工"), c)
			return
		}
		err = global.GDb.Model(&model.SysProPermission{}).Where("user_code = ? AND pro_code = ? ", userx[0].Code, global.GConfig.System.ProCode).Find(&UserPermission).Error
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("登录失败"), c)
			return
		}
		reg1 := regexp.MustCompile("^[P]")
		ok := false
		if reg1.MatchString(L.Username) {
			ok = true
		} else {
			if UserPermission.HasOutNet {
				ok = true
			}
		}

		if ok {
			// if GetLoginCode(userx[0].Qrcode) == L.Captcha {
			fmt.Println("global.GStore.Verify", global.GStore.Verify(L.CaptchaId, L.Captcha, true))
			U := &model.SysUser{Phone: L.Phone, Password: L.Password, Status: true}
			if err, user := service.Login(U, L.IsIntranet); err != nil {
				response.FailWithMessage(fmt.Sprintf("用户名/密码错误"), c)
			} else {
				err = utils.DelRedisCaptcha(L.Phone)
				if err != nil {
					response.FailWithMessage("验证码缓存删除失败!", c)
				}
				err = service.CheckPasswordTime(user)
				if err != nil {
					response.FailWithMessage(err.Error(), c)
					return
				}
				tokenNext(c, *user)
			}
			// } else {
			// 	response.FailWithMessage("验证码错误", c)
			// }
		} else {
			response.FailWithMessage(fmt.Sprintf("登录失败, 该用户没有外网登录的权限"), c)
			return
		}

	} else {
		UserVerify := utils.Rules{
			"Code":     {utils.NotEmpty()},
			"Password": {utils.NotEmpty()},
		}
		UserVerifyErr := utils.Verify(L, UserVerify)
		if UserVerifyErr != nil {
			response.FailWithMessage(UserVerifyErr.Error(), c)
			return
		}

		var UserPermission model.SysProPermission
		err := global.GDb.Model(&model.SysProPermission{}).Where("user_code = ? AND pro_code = ? ", L.Username, global.GConfig.System.ProCode).Find(&UserPermission).Error
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("登录失败"), c)
			return
		}
		reg1 := regexp.MustCompile("^[P]")
		ok := false
		if reg1.MatchString(L.Username) {
			if UserPermission.HasInNet {
				ok = true
			}
		} else {
			ok = true
		}

		if ok {
			U := &model.SysUser{Code: L.Username, Password: L.Password, Status: true}
			if err, user := service.Login(U, L.IsIntranet); err != nil {
				response.FailWithMessage(fmt.Sprintf("用户名/密码错误"), c)
			} else {
				err = service.CheckPasswordTime(user)
				if err != nil {
					response.FailWithMessage(err.Error(), c)
					return
				}
				tokenNext(c, *user)
			}
		} else {
			response.FailWithMessage(fmt.Sprintf("登录失败, 该用户没有内网登录的权限"), c)
			return
		}
	}
}

func GetUserQrCode(c *gin.Context) {
	var L request2.Logins
	_ = c.ShouldBindJSON(&L)
	UserVerify := utils.Rules{
		//"CaptchaId": {utils.NotEmpty()},
		// "Captcha":  {utils.NotEmpty()},
		"Phone":    {utils.NotEmpty()},
		"Password": {utils.NotEmpty()},
	}
	UserVerifyErr := utils.Verify(L, UserVerify)
	if UserVerifyErr != nil {
		response.FailWithMessage(UserVerifyErr.Error(), c)
		return
	}

	// err, a := utils.GetRedisCaptcha(L.Phone)
	// if err == nil && a == "" {
	// 	response.FailWithMessage("验证码已过期!", c)
	// 	return
	// }
	// fmt.Println("------------aaaa---------------", a, L.Captcha)

	// if a == L.Captcha {
	// 	fmt.Println("global.GStore.Verify", global.GStore.Verify(L.CaptchaId, L.Captcha, true))
	U := &model.SysUser{Phone: L.Phone, Password: L.Password, Status: true}
	if err, user := service.Login(U, L.IsIntranet); err != nil {
		response.FailWithMessage(fmt.Sprintf("用户名/密码错误"), c)
	} else {
		// err = utils.DelRedisCaptcha(L.Phone)
		if err != nil {
			response.FailWithMessage("验证码缓存删除失败!", c)
		}
		qrcode := ""
		if user.Qrcode != "" {
			qrcode = "otpauth://totp/" + user.Code + ":i-confluence.com?algorithm=SHA512&digits=6&issuer=" + user.Code + "&period=60&secret=" + user.Qrcode
		} else {
			key, err := GetQrCode(user.Code)
			fmt.Println("--------------111111111------------------", key, err, key.Secret())
			if err != nil {
				response.FailWithMessage("获取二维码失败!", c)
				return
			}
			err2 := service.UpdateUserQrcodeById(user.ID, key.Secret())
			fmt.Println("--------------err2err2err2------------------", err2)
			if err2 != nil {
				response.FailWithMessage("获取二维码失败!", c)
				return
			}
			qrcode = key.URL()
		}
		fmt.Println("--------------GetUserQrCode------------------", qrcode)
		response.OkWithData(qrcode, c)
	}
	// } else {
	// 	response.FailWithMessage("验证码错误", c)
	// }

	// if user.Qrcode != "" {
	// 	qrcode = user.Qrcode
	// } else {
	// 	key, err := GetQrCode()
	// 	if err != nil {
	// 		err2 := service.UpdateUserQrcodeById(user.ID, key.Secret())
	// 		if err2 != nil {
	// 			qrcode = key.URL()
	// 		}
	// 	}
	// }
	// fmt.Println("--------------GetUserQrCode------------------", qrcode)
	// return qrcode
}

func GetQrCode(code string) (*otp.Key, error) {
	data := totp.GenerateOpts{
		Period:      60,
		Digits:      6,
		Algorithm:   otp.AlgorithmSHA512,
		AccountName: "i-confluence.com",
		Issuer:      code,
		SecretSize:  20,
	}
	return totp.Generate(data)
}

func GetLoginCode(value string) string {
	code, err := totp.GenerateCodeCustom(value, time.Now(), totp.ValidateOpts{
		Period:    60,
		Skew:      2,
		Digits:    6,
		Algorithm: otp.AlgorithmSHA512,
	})
	fmt.Println("---------totptotptotptotp----------------", value, code, err)
	return code
}

// Login
// @Tags SysUser
// @Summary 用户登录
// @Produce  application/json
// @Param data body request2.Logins true "用户登录接口"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登陆成功"}"
// @Router /sys-base/sys-login/login [post]
func Login(c *gin.Context) {
	var L request2.Logins
	_ = c.ShouldBindJSON(&L)
	if !L.IsIntranet {
		UserVerify := utils.Rules{
			//"CaptchaId": {utils.NotEmpty()},
			// "Captcha":  {utils.NotEmpty()},
			"Phone":    {utils.NotEmpty()},
			"Password": {utils.NotEmpty()},
		}
		UserVerifyErr := utils.Verify(L, UserVerify)
		if UserVerifyErr != nil {
			response.FailWithMessage(UserVerifyErr.Error(), c)
			return
		}

		// err, a := utils.GetRedisCaptcha(L.Phone)
		// if err == nil && a == "" {
		// 	response.FailWithMessage("验证码已过期!", c)
		// 	return
		// }

		// if a == L.Captcha {
		fmt.Println("global.GStore.Verify", global.GStore.Verify(L.CaptchaId, L.Captcha, true))
		U := &model.SysUser{Phone: L.Phone, Password: L.Password, Status: true}
		if err, user := service.Login(U, L.IsIntranet); err != nil {
			response.FailWithMessage(fmt.Sprintf("用户名/密码错误"), c)
			return
		} else {
			// err = utils.DelRedisCaptcha(L.Phone)
			// if err != nil {
			// 	response.FailWithMessage("验证码缓存删除失败!", c)
			// }
			// if GetLoginCode(user.Qrcode) == L.Captcha {
			err = service.CheckPasswordTime(user)
			if err != nil {
				response.FailWithMessage(err.Error(), c)
				return
			}
			tokenNext(c, *user)
			// } else {
			// 	response.FailWithMessage("验证码错误", c)
			// 	return
			// }
		}
		// } else {
		// 	response.FailWithMessage("验证码错误", c)
		// }
	} else {
		UserVerify := utils.Rules{
			"Code":     {utils.NotEmpty()},
			"Password": {utils.NotEmpty()},
		}
		UserVerifyErr := utils.Verify(L, UserVerify)
		if UserVerifyErr != nil {
			response.FailWithMessage(UserVerifyErr.Error(), c)
			return
		}
		U := &model.SysUser{Code: L.Username, Password: L.Password, Status: true}
		if err, user := service.Login(U, L.IsIntranet); err != nil {
			response.FailWithMessage(fmt.Sprintf("用户名/密码错误"), c)
		} else {
			err = service.CheckPasswordTime(user)
			if err != nil {
				response.FailWithMessage(err.Error(), c)
				return
			}
			tokenNext(c, *user)
		}
	}
}

// 登录以后签发jwt
func tokenNext(c *gin.Context, user model.SysUser) {
	err := storePermToRedis(user)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("权限获取失败,%s", err.Error()), c)
		return
	}
	j := &middleware.JWT{
		SigningKey: []byte(global.GConfig.JWT.SigningKey), // 唯一签名
	}
	if user.SysRoles.ID == "" {
		response.FailWithMessage("获取角色失败或者角色已停用", c)
		return
	}
	//fmt.Println("tokeNext", auth)
	global.GLog.Info("global.GConfig.JWT.ExpiresAt", zap.Any("s", global.GConfig.JWT.ExpiresAt))

	now := time.Now().Local()
	expiresDuration := time.Duration(global.GConfig.JWT.ExpiresAt * 1000 * 1000 * 1000)
	newDay := now.Add(expiresDuration).Format("2006-01-02") + " 02:01:01"
	expiration, err := time.ParseInLocation("2006-01-02 15:04:05", newDay, time.Local)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("生成过期时间失败,%s", err.Error()), c)
		return
	}
	global.GLog.Info("登录:token过期时间", zap.Any(newDay, expiration))

	clams := request2.CustomClaims{
		//UUID:        user.UUID,
		ID:       user.ID,
		Phone:    user.Phone,
		NickName: user.NickName,
		DingId:   user.DingId,
		Code:     user.Code,
		RoleId:   user.RoleId,
		RoleName: user.SysRoles.Name,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000, // 签名生效时间
			ExpiresAt: expiration.Unix(),        // 过期时间
			Issuer:    "xingqiyi",               // 签名的发行者
		},
	}
	fmt.Println("asd", clams.StandardClaims)
	fmt.Println("time.Now().Unix()", time.Now().Unix())
	fmt.Println("expiration.Unix()", expiration.Unix())
	err1, token := service.GetRedisJWT(user.ID)
	if err1 != nil {
		token, err = j.CreateToken(clams)
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("获取token失败,%s", err.Error()), c)
			return
		}
	}
	//过期了就重新签证
	if _, err = j.ParseToken(token); err == middleware.TokenExpired {
		token, err = j.CreateToken(clams)
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("获取token失败,%s", err.Error()), c)
			return
		}
	}

	if !global.GConfig.System.UseMultipoint {
		response.OkWithData(sys_base4.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: clams.StandardClaims.ExpiresAt * 1000,
		}, c)
		return
	}
	if err = service.SetRedisJWT(token, user.ID, expiration); err != nil {
		response.FailWithMessage(fmt.Sprintf("设置登录状态失败,%s", err.Error()), c)
		return
	}

	secret := ""
	ttl := service.GetRedisSecretWithTTL(user.ID)
	global.GLog.Info("ttl", zap.Any("s", ttl))
	if ttl < time.Minute {
		secret = utils.Generate(user.ID)
		err = service.SetRedisSecret(secret, user.ID)
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("设置登录状态失败1,%s", err.Error()), c)
			return
		}
	} else {
		secret, err = service.GetRedisSecret(user.ID)
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("设置登录状态失败2,%s", err.Error()), c)
			return
		}
	}

	response.OkWithData(sys_base4.LoginResponse{
		User:      user,
		Token:     token,
		Secret:    secret,
		ExpiresAt: clams.StandardClaims.ExpiresAt * 1000,
	}, c)
}

func storePermToRedis(u model.SysUser) error {
	err, perms := service.GetAllPermissionByUId(u.ID)
	if err != nil {
		return err
	}
	var permMap = make(map[string]interface{})
	for _, perm := range perms {
		permMap[perm.ProCode+"_op0"] = perm.HasOp0
		permMap[perm.ProCode+"_op1"] = perm.HasOp1
		permMap[perm.ProCode+"_op2"] = perm.HasOp2
		permMap[perm.ProCode+"_opq"] = perm.HasOpq
		permMap[perm.ProCode+"_HasInNet"] = perm.HasInNet
		permMap[perm.ProCode+"_HasOutNet"] = perm.HasOutNet
		permMap[perm.ProCode+"_pm"] = perm.HasPm
	}
	if len(permMap) == 0 {
		return nil
	}
	err = service.SetRedisPerm(permMap, u.ID)
	return err
}

type UserHeaderImg struct {
	HeaderImg multipart.File `json:"headerImg"`
}

// GetUserList
// @Tags SysBase
// @Summary 分页获取用户列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.PageInfo true "分页获取用户列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /user/getUserList [post]
func GetUserList(c *gin.Context) {
	var pageInfo model.PageInfo
	_ = c.ShouldBindJSON(&pageInfo)
	PageVerifyErr := utils.Verify(pageInfo, utils.CustomizeMap["PageVerify"])
	if PageVerifyErr != nil {
		response.FailWithMessage(PageVerifyErr.Error(), c)
		return
	}
	err, list, total := service.GetUserInfoList(pageInfo)
	//global.G_LOG.Info(zap.Int64("123",total));
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取数据失败，%v", err), c)
	} else {
		response.OkWithData(sys_base4.PageResult{
			List:      list,
			Total:     total,
			PageIndex: pageInfo.PageIndex,
			PageSize:  pageInfo.PageSize,
		}, c)
	}
}

// UploadUserImage
// @Tags SysBase
// @Summary 用户管理--用户上传头像
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param headerImg formData file true "用户上传头像"
// @Param username formData string true "用户名"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"上传成功"}"
// @Router /user/uploadHeaderImg [post]
func UploadUserImage(c *gin.Context) {
	form, err := c.MultipartForm()
	files := form.File["file"]
	sysProTempId := c.PostForm("sysProTempId")
	pathArr := make([]string, 0)
	for _, file := range files {
		_, sysProTempBIdpath := utils.SaveImgFile(c, file, file.Filename, global.GConfig.LocalUpload.FilePath+"/SysUserImage/"+sysProTempId+"/")
		pathArr = append(pathArr, sysProTempBIdpath)
	}
	//intId, _ := strconv.Atoi(sysProTempId)
	reRow := service2.UpdateImages(pathArr, sysProTempId)
	if reRow != 1 {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// Sync
// @Tags SysBase
// @Summary 人员管理--同步新增的数据
// @Auth xingqiyi
// @Date 2022/1/4 10:09 上午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sys-user/sync [get]
func Sync(c *gin.Context) {
	var maxId time.Time
	err := global.GDb.Model(&model.MaxUserId{}).Where("id = '1'").Select("updated_at").First(&maxId).Error
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
		return
	}
	maxTime := maxId.Format("2006-01-02 15:04:05")
	global.GLog.Warn("old maxId:::" + maxTime)
	u, r := sys_base.Move(maxTime, false)
	//u, r := sys_base.Move("", true)
	response.OkWithData(map[string]int64{"user": u, "role": r}, c)
}

// ChangeRole
// @Tags SysUser
// @Summary 人员管理--修改用户角色
// @Auth xingqiyi
// @Date 2022/3/29 15:03
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request2.ChangeRoleForm true "ChangeRoleForm"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sys-base/user-management/sys-user/change-role [post]
func ChangeRole(c *gin.Context) {
	var changeRoleForm request2.ChangeRoleForm
	_ = c.ShouldBindJSON(&changeRoleForm)
	verify := utils.Rules{
		"Id":     {utils.NotEmpty()},
		"RoleId": {utils.NotEmpty()},
	}
	verifyErr := utils.Verify(changeRoleForm, verify)
	if verifyErr != nil {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}
	err := service.DelRedisJWT(changeRoleForm.Id)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("修改失败，%v", err), c)
		return
	}
	row := service.ChangeRole(changeRoleForm)
	if row != 1 {
		response.FailWithMessage(fmt.Sprintf("修改失败，%v", row), c)
	} else {
		response.OkWithMessage("修改成功", c)
	}
}

// SysProPermissionExport
// @Tags 项目管理--人员管理
// @Summary 项目管理--导出项目权限
// @Auth xingqiyi
// @Date 2022年11月16日14:46:15
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode       query   string      true    "项目编码全部则为：all"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sys-base/user-management/sys-pro-permission/export [get]
func SysProPermissionExport(c *gin.Context) {
	proCode := c.Query("proCode")
	if proCode == "" {
		proCode = "all"
	}
	list, err := service.GetPermissionByProCode(proCode)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("导出失败，%v", err.Error()), c)
		return
	}
	//fmt.Println(list)
	var sysProPermissionExports []model.SysProPermissionExport
	for _, permission := range list {
		if _, ok := global.UserCodeName[permission.UserCode]; !ok {
			global.GLog.Error("没有这个user permission.UserCode：" + permission.UserCode)
			continue
		}
		var sysProPermissionExport = model.SysProPermissionExport{
			ProCode:   permission.ProCode,
			UserCode:  permission.UserCode,
			UserName:  global.UserCodeName[permission.UserCode],
			HasOp0:    permission.HasOp0,
			HasOp1:    permission.HasOp1,
			HasOp2:    permission.HasOp2,
			HasOpq:    permission.HasOpq,
			HasPm:     permission.HasPm,
			HasInNet:  permission.HasInNet,
			HasOutNet: permission.HasOutNet,
		}
		sysProPermissionExports = append(sysProPermissionExports, sysProPermissionExport)
	}
	path := global.GConfig.LocalUpload.FilePath + global.PathUserPermissionDownload
	// 尝试创建此路径
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		global.GLog.Error("upload file fail:", zap.Any("err", err))
		return
	}
	name := fmt.Sprintf("理赔2.0项目权限表-%v.xlsx", time.Now().Format("20060102"))
	err = utils.ExportBigExcel(path, name, "sheet", sysProPermissionExports)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(path+name, c)
	//c.FileAttachment(path+name, name)
}

// SysProPermissionImport
// @Tags 项目管理--人员管理
// @Summary 项目管理--导入项目权限
// @Auth xingqiyi
// @Date 2022年11月16日14:48:42
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce application/json
// @Param file formData file true "权限excel"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sys-base/user-management/sys-pro-permission/import [post]
func SysProPermissionImport(c *gin.Context) {
	form, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("上传的excel有问题呀,%s", err.Error()), c)
		return
	}
	path := global.GConfig.LocalUpload.FilePath + global.PathUserPermissionUpload
	// 尝试创建此路径
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		global.GLog.Error("upload file fail:", zap.Any("err", err))
		return
	}
	dst := path + time.Now().Format("20060102150405") + form.Filename
	global.GLog.Info("文件路径", zap.Any("dst", dst))
	err = c.SaveUploadedFile(form, dst)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("excel保存有问题呀,%s", err.Error()), c)
		return
	}
	xlsx, err := excelize.OpenFile(dst)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("excel读取有问题呀,%s", err.Error()), c)
		return
	}
	rows, err := xlsx.GetRows(xlsx.GetSheetName(xlsx.GetActiveSheetIndex()))
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("excel获取行有问题呀,%s", err.Error()), c)
		return
	}
	//fmt.Println(rows)
	var relationMap = map[string]bool{
		"t":     true,
		"true":  true,
		"TRUE":  true,
		"1":     true,
		"f":     false,
		"false": false,
		"FALSE": false,
		"0":     false,
	}
	var sysProPermissions []model.SysProPermission
	//项目编码	工号	姓名	初审	一码	二码	问题件	PM	内网	外网
	for i, row := range rows {
		//if _, ok := global.UserCodeAndId[row[1]]; !ok {
		//	global.GLog.Error("跳过工号" + row[1])
		//	continue
		//}
		var sysProPermission model.SysProPermission
		if i == 0 {
			if row[0] != "项目编码" || row[1] != "工号" || row[2] != "姓名" || row[3] != "初审" ||
				row[4] != "一码" || row[5] != "二码" || row[6] != "问题件" || row[7] != "PM" ||
				row[8] != "内网" || row[9] != "外网" {
				response.FailWithMessage("excel表头有误", c)
				return
			}
		} else {
			err, u := service.FetchUserIDAndObjectIDByCode(row[1])
			if err != nil {
				global.GLog.Error("没找到工号："+row[1], zap.Error(err))
				continue
			}
			sysProPermission.HasOp0 = relationMap[row[3]]
			sysProPermission.HasOp1 = relationMap[row[4]]
			sysProPermission.HasOp2 = relationMap[row[5]]
			sysProPermission.HasOpq = relationMap[row[6]]
			sysProPermission.HasPm = relationMap[row[7]]
			sysProPermission.HasInNet = relationMap[row[8]]
			sysProPermission.HasOutNet = relationMap[row[9]]
			sysProPermission.ProCode = row[0]
			sysProPermission.ProId = global.ProCodeId[row[0]]
			sysProPermission.ProName = global.GProConf[row[0]].Name
			sysProPermission.UserCode = row[1]
			sysProPermission.UserId = u.ID
			sysProPermission.ObjectId = u.ObjectId
			sysProPermissions = append(sysProPermissions, sysProPermission)
		}
	}
	r := service.UpdateProPermission(sysProPermissions)

	if r == 0 {
		response.FailWithMessage(fmt.Sprintf("导入数据失败，%v", err), c)
	} else {
		response.Ok(c)
	}
}

// GetUserInfo
// @Tags SysUser
// @Summary 获取用户信息
// @Description
// @Date 2023年08月16日15:26:09
// @Security ApiKeyAuth
// @Security UserID
// @Accept json
// @Produce json
// @Success 200 {string} sys_base4.LoginResponse
// @Router /sys-base/sys-login/user-info [post]
func GetUserInfo(c *gin.Context) {
	customClaims, err := api.GetUserByToken(c)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithoutCode(err, "失败", c)
		return
	}
	if c.GetHeader(global.XUserId) != customClaims.ID {
		response.FailWithoutCode(err, "失败", c)
		return
	}
	err, user := service.GetUserByIdStage(customClaims.ID)
	if err != nil {
		response.FailWithoutCode(err, "失败", c)
		return
	}
	response.OkWithData(sys_base4.LoginResponse{
		User:      user,
		Token:     c.GetHeader(global.XToken),
		ExpiresAt: customClaims.StandardClaims.ExpiresAt * 1000,
	}, c)
}

// UploadUserAvatar
// @Tags SysUser
// @Summary 员工信息-上传用户头像
// @Auth Tuesday
// @Date 2023/8/22 10:30 上午
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce application/json
// @Param file	formData	file	true	"file"
// @Param jobNumber	formData	string	true    "工号"
// @Param jobName	formData	string	true    "姓名"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sys-user/upload-user/avatar [post]
func UploadUserAvatar(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 32<<20)
	//获取用户头像
	file, err := c.FormFile("file")
	//工号和姓名
	jobNumber := c.PostForm("jobNumber")
	jobName := c.PostForm("jobName")
	if jobNumber == "" || jobName == "" {
		response.FailWithMessage("工号或姓名不能为空,请输入正确工号或姓名", c)
		return
	}
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//判断文件后缀名是否合法
	ext := path.Ext(file.Filename)
	suffixName := map[string]bool{
		".jpg":  true,
		".JPG":  true,
		".png":  true,
		".PNG":  true,
		".tif":  true,
		".TIF":  true,
		".jpge": true,
		".JPGE": true,
	}
	_, legal := suffixName[ext]
	if !legal {
		response.FailWithMessage("上传的文件不合法,文件后缀名必须是 .jpg 或 .png 或 .tif 或 .jpge", c)
		return
	}
	//查询是否有该员工
	err, count := service.UserCount(jobNumber, jobName)
	if count <= 0 {
		response.FailWithMessage("没有该用户,请输入正确的工号或姓名并确认在职", c)
		return
	}
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err, headerPath := service.UploadUserAvatar(c, file, jobNumber, jobName, global.GConfig.LocalUpload.FilePath+global.PathUpdateUserHeaderImg+jobNumber+"/")
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("修改失败，%v", err), c)
		return
	} else {
		response.OkWithData(headerPath, c)
	}

}
