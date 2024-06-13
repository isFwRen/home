package model

type TabsModel struct {
	AutoAddIdModel
	Title       string `json:"title" gorm:"标题"`
	Name 		string `json:"name" gorm:"类似索引(定位到当前Tabs内容)"`
	Path        string `json:"path" gorm:"路径"`
	QueryPath   string `json:"queryPath" gorm:"带参路径"`
}

