package request

type SysProjectUpdateRecord struct {
	AutoReturn  bool   `json:"autoReturn" form:"autoReturn" gorm:"回传方式"`
	CacheTime   int    `json:"cacheTime" form:"cacheTime" gorm:"缓存时间"`
	Id          string `json:"id" form:"id" gorm:"id"`
	RestartAt   string `json:"restartAt" form:"restartAt" gorm:"重启时间"`
	SaveDate    int    `json:"saveDate" form:"saveDate" gorm:"数据保留天数"`
	EditVersion int    `json:"editVersion" form:"editVersion" gorm:"version版本"`
	Type        string `json:"type" form:"type" gorm:"所属行业"`
	Code        string `json:"code" form:"code" gorm:"项目编码"`
}
