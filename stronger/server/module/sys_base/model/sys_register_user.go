package model

import (
	"golang.org/x/crypto/ssh"
	"time"
	//"gorm.io/gorm"
)

type SysRegisterUser struct {
	//gorm.Model
	Model
	Username    string       `json:"userName" gorm:"comment:'用户登录名'"`
	Password    string       `json:"-"  gorm:"comment:'用户登录密码'"`
	NickName    string       `json:"nickName" gorm:"default:'系统用户';comment:'用户昵称'" `
	HeaderImg   string       `json:"headerImg" gorm:"default:'http://qmplusimg.henrongyi.top/head.png';comment:'用户头像'"`
	Authority   SysAuthority `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"`
	AuthorityId string       `json:"authorityId" gorm:"default:888;comment:'用户角色ID'"`
	Phone       string       `json:"phone" gorm:"comment:'手机号'"`
	DingId      string       `json:"dingId" gorm:"comment:'钉钉id'"`
	EntryDate   time.Time    `json:"entryDate" gorm:"comment:'入职日期'"`
	LeaveDate   time.Time    `json:"leaveDate" gorm:"comment:'离职日期'"`
	Status      bool         `json:"status" gorm:"comment:'在职or离职'"`
	Staff       string       `json:"staff" gorm:"comment:'角色'"`
	Email       string       `json:"email" gorm:"comment:'邮箱'"`
	Sex         bool         `json:"sex" gorm:"comment:'性别'"`
	SysRoles    []SysRoles   `json:"sysRoles" gorm:"foreignKey:id;comment:权限"`
}

type Cli struct {
	IP         string      //IP地址
	Username   string      //用户名
	Password   string      //密码
	Port       int         //端口号
	Client     *ssh.Client //ssh客户端
	LastResult string      //最近一次Run的结果
}
