package request

type ConstTableWithProject struct {
	ProCode string `json:"proCode"`
	Conf    string `json:"conf"`
}

type ConstQueryStruct struct {
	Selector ConstTableWithProject `json:"selector"`
	Limit    int                   `json:"limit"`
}

type ConstTableWithProjectIsExist struct {
	ProCode string `json:"proCode"`
	Conf    string `json:"conf"`
	Name    string `json:"name"`
}

type ConstQueryStructIsExist struct {
	Selector ConstTableWithProjectIsExist `json:"selector"`
}

type InsertConstTableStructWithItems struct {
	Id          string          `json:"_id"`
	ProCode     string          `json:"proCode"`
	Conf        string          `json:"conf"`
	FileName    string          `json:"fileName"`
	DbName      string          `json:"dbName"`
	FilePath    string          `json:"filePath"`
	ChineseName string          `json:"chineseName"`
	TableTop    []TableTopItems `json:"table_top"`
	Docs        []interface{}   `json:"docs"`
}

type Docs struct {
	Content string `json:"content"`
}

type TableTopItems struct {
	Tabletop []string `json:"tabletop"`
}

type Items struct {
	Id  string   `json:"_id"`
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

type UpdateConstTableStructWithItems struct {
	Id       string  `json:"_id"`
	Rev      string  `json:"_rev"`
	ProCode  string  `json:"proCode"`
	Conf     string  `json:"conf" example:"类型，默认为const"`
	FileName string  `json:"fileName"`
	Docs     []Items `json:"docs"` //需要更新的内容
}

type DelConstTableLineArr struct {
	DbName string              `json:"dbName"`
	Table  []DelConstTableLine `json:"table"`
}

type DelConstTableLine struct {
	Id  string `json:"_id"`
	Rev string `json:"_rev"`
}

type DelConstTable struct {
	DbName string `json:"dbName" form:"dbName"`
}

type IsExist struct {
	ProCode   string `json:"proCode"`
	ConstName string `json:"name"`
}

type PutConst struct {
	Id   string  `json:"_id"`
	Rev  string  `json:"_rev"`
	Docs []Items `json:"docs"`
}

type ConstRequest struct {
	ProCode      string            `json:"proCode"`
	Relationship map[string]string `json:"relationship"`
}
