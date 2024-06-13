package model

import (
	model2 "server/module/sys_base/model"
)

type SysProjectCache struct {
	model2.Model
	Name           string                `json:"name" form:"name" gorm:"comment:'项目名字'"`         //项目名字
	Type           string                `json:"type" form:"type" gorm:"项目类型"`                   //项目类型
	CacheTime      int                   `json:"cacheTime" form:"cacheTime" gorm:"缓存时间"`         //缓存时间
	AutoReturn     bool                  `json:"autoReturn" form:"autoReturn" gorm:"回传方式"`       //回传方式
	SaveDate       int                   `json:"saveDate" form:"saveDate" gorm:"数据保留天数"`         //数据保留天数
	RestartAt      string                `json:"restartAt" form:"restartAt" gorm:"重启时间"`         //重启时间
	Code           string                `json:"code" form:"code" gorm:"项目编码"`                   //项目编码
	DbHistory      string                `json:"dbHistory" form:"dbHistory" gorm:"数据库dsn"`       //数据库dsn
	DbTask         string                `json:"dbTask" form:"dbTask" gorm:"录入db"`               //录入db
	InAppPort      int                   `json:"inAppPort" form:"inAppPort" gorm:"内网录入端口"`       //录入端口
	OutAppPort     int                   `json:"outAppPort" form:"outAppPort" gorm:"外网录入端口"`     //录入端口
	BackEndPort    int                   `json:"backEndPort" form:"backEndPort" gorm:"录入后端端口"`   //录入后端端口
	InnerIp        string                `json:"innerIp" form:"innerIp" gorm:"内网ip"`             //内网ip
	OutIp          string                `json:"outIp" form:"outIp" gorm:"外网ip"`                 //外网ip
	MaxDownload    int                   `json:"maxDownload" form:"maxDownload" gorm:"一次最大下载个数"` //一次最大下载个数
	EditVersion    int                   `json:"editVersion" form:"editVersion" gorm:"版本"`       //一次最大下载个数
	Description    string                `json:"description" form:"description" gorm:"描述"`       //描述
	DownloadPaths  SysProDownloadPaths   `gorm:"foreignKey:ProId;references:ID;comment:下载路径"`
	UploadPaths    SysProDownloadPaths   `gorm:"foreignKey:ProId;references:id;comment:上传路径"`
	ConstTable     map[string][][]string `json:"constTable" gorm:"-"`     //常量
}

type Items struct {
	Arr []string `json:"arr"`
}

type ConstTableBaseInformation struct {
	ProCode     string `json:"proCode"`
	Conf        string `json:"conf"`
	FileName    string `json:"fileName"`
	DbName      string `json:"dbName"`
	FilePath    string `json:"filePath"`
	ChineseName string `json:"chineseName"`
}

type TableTop struct {
	Id       string   `json:"_id"`
	Rev      string   `json:"_rev"`
	Tabletop []string `json:"tabletop"`
}
