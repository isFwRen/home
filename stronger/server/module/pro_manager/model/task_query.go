package model

type CaseMessage struct {
	BillId        string  `json:"billId"`
	ProCode       string  `json:"proCode"`                            //项目编号
	BillName      string  `json:"billName"`                           //单据号来源内部
	BillNum       string  `json:"billNum"`                            //单据号来源客户
	Agency        string  `json:"agency"`                             //机构号
	ScanAt        string  `json:"scanAt"`                             //扫描时间
	LastAt        string  `json:"last_at"`                            //最晚回传时间
	RemainderAt   string  `json:"remainderAt"`                        //剩余时间
	Second        float64 `json:"second"`                             //剩余时间(秒), 给前端排序
	Stage         string  `json:"stage"`                              //录入状态
	StickLevel    int     `json:"stickLevel" from:"stickLevel"`       //加急件
	AppCompleteAt string  `json:"appCompleteAt" from:"appCompleteAt"` //录入完成时间
	SaleChannel   string  `json:"saleChannel" form:"saleChannel"`     //销售渠道

}

// TaskList 不要随便修改！
type TaskList struct {
	Op0notassign int `json:"op0NotAssign"` //初审待分配
	Op0assign    int `json:"op0Assign"`    //初审已分配
	Op0cache     int `json:"op0Cache"`     //初审缓存区
	//--------------------------------一码-----------------------------------
	Op1notassignexpenseaccount    int `json:"op1NotAssignExpenseAccount"`    //一码非报销单待分配
	Op1assignexpenseaccount       int `json:"op1AssignExpenseAccount"`       //一码非报销单已分配
	Op1cacheexpenseaccount        int `json:"op1CacheExpenseAccount"`        //一码非报销单缓存区
	Op1notassignnotexpenseaccount int `json:"op1NotAssignNotExpenseAccount"` //一码报销单待分配
	Op1assignnotexpenseaccount    int `json:"op1AssignNotExpenseAccount"`    //一码报销单已分配
	Op1cachenotexpenseaccount     int `json:"op1CacheNotExpenseAccount"`     //一码报销单缓存区
	//--------------------------------二码-----------------------------------
	Op2notassignexpenseaccount    int `json:"op2NotAssignExpenseAccount"`    //二码非报销单待分配
	Op2assignexpenseaccount       int `json:"op2AssignExpenseAccount"`       //二码非报销单已分配
	Op2cacheexpenseaccount        int `json:"op2CacheExpenseAccount"`        //二码非报销单缓存区
	Op2notassignnotexpenseaccount int `json:"op2NotAssignNotExpenseAccount"` //二码报销单待分配
	Op2assignnotexpenseaccount    int `json:"op2AssignNotExpenseAccount"`    //二码报销单已分配
	Op2cachenotexpenseaccount     int `json:"op2CacheNotExpenseAccount"`     //二码报销单缓存区
	//--------------------------------问题件-----------------------------------
	Opqnotassign int      `json:"opqNotAssign"` //问题件待分配
	Opqassign    int      `json:"opqAssign"`    //问题件已分配
	Opqcache     int      `json:"opqCache"`     //问题件缓存区
	Urgent       []string `json:"urgent"`       //紧急件
	Priority     []string `json:"priority"`     //优先件
	Agency       []string `json:"agency"`       //机构号
	Num          int64    `json:"num"`          //单子数量
}

type TaskListDetailNotAssign struct {
	BlockId    string `json:"blockId"`    //分块ID
	BlockName  string `json:"blockName"`  //分块名称
	BillNum    string `json:"billNum"`    //案件号
	Agency     string `json:"agency"`     //机构
	TempType   string `json:"tempType"`   //模板类型
	TaskAssign string `json:"taskAssign"` //任务分配
}

type TaskListDetailAssignAndCache struct {
	BlockId     string `json:"blockId"`     //分块ID
	BlockName   string `json:"blockName"`   //分块名称
	BillNum     string `json:"billNum"`     //案件号
	Agency      string `json:"agency"`      //机构
	Op1Code     string `json:"op1Code"`     //1码人员编号(0为系统录入)
	Op1ApplyAt  string `json:"op1ApplyAt"`  //1码领取时间
	Op1SubmitAt string `json:"op1SubmitAt"` //1码提交时间
	Op2Code     string `json:"op2Code"`     //2码人员编号
	Op2ApplyAt  string `json:"op2ApplyAt"`  //2码领取时间
	Op2SubmitAt string `json:"op2SubmitAt"` //2码提交时间
	OpqCode     string `json:"opqCode"`     //问题件人员编号
	OpqApplyAt  string `json:"opqApplyAt"`  //问题件领取时间
	OpqSubmitAt string `json:"opqSubmitAt"` //问题件提交时间
	Op0Code     string `json:"op0Code"`     //初审人员编号
	Op0ApplyAt  string `json:"op0ApplyAt"`  //初审领取时间
	Op0SubmitAt string `json:"op0SubmitAt"` //初审提交时间
}
