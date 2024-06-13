package response

type GetTaskList struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type GetVariousStateBillResponse struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type SetPriorityOrganizationNumberResponse struct {
	OrganizationNumber string `json:"organization_number"`
	//UpdateResult       []*mongo.UpdateResult `json:"update_result"`
}

type GetCaseDetailsResponse struct {
	List  interface{} `json:"list"`
	Total int         `json:"total"`
}

type GetCaseDetailsBlockResponse struct {
	ProCode     string `json:"proCode"`
	BlockId     string `json:"blockId"`
	BillID      string `json:"billID" gorm:"comment:'单ID'"`                //单ID
	Name        string `json:"name" gorm:"模板分块名字"`                         //模板分块名字
	Code        string `json:"code" gorm:"模板分块编码"`                         //模板分块编码
	Stage       string `json:"stage" gorm:"comment:'状态(op0|op1|op2|opQ)'"` //状态
	Status      string `json:"status" gorm:"comment:'1(初审)|2'"`            //1(初审)|2
	Op1Code     string `json:"op1Code" gorm:"comment:'1码人员编号(0为系统录入)'"`    //1码人员编号(0为系统录入)
	Op1ApplyAt  string `json:"op1ApplyAt" gorm:"comment:'1码领取时间'"`         //1码领取时间
	Op1SubmitAt string `json:"op1SubmitAt" gorm:"comment:'1码提交时间'"`        //1码提交时间
	Op2Code     string `json:"op2Code" gorm:"comment:'2码人员编号'"`            //2码人员编号
	Op2ApplyAt  string `json:"op2ApplyAt" gorm:"comment:'2码领取时间'"`         //2码领取时间
	Op2SubmitAt string `json:"op2SubmitAt" gorm:"comment:'2码提交时间'"`        //2码提交时间
	OpqCode     string `json:"opqCode" gorm:"comment:'问题件人员编号'"`           //问题件人员编号
	OpqApplyAt  string `json:"opqApplyAt" gorm:"comment:'问题件领取时间'"`        //问题件领取时间
	OpqSubmitAt string `json:"opqSubmitAt" gorm:"comment:'问题件提交时间'"`       //问题件提交时间
	Op0Code     string `json:"op0Code" gorm:"comment:'初审人员编号'"`            //初审人员编号
	Op0ApplyAt  string `json:"op0ApplyAt" gorm:"comment:'初审领取时间'"`         //初审领取时间
	Op0SubmitAt string `json:"op0SubmitAt" gorm:"comment:'初审提交时间'"`        //初审提交时间
}

type GetCaseDetailsFieldResponse struct {
	ProCode  string `json:"proCode"`
	BillId   string `json:"billId"`
	BlockId  string `json:"blockId"`
	FieldId  string `json:"fieldId"`
	Name     string `json:"name" form:"name" gorm:"comment:'字段名字'"`             //字段名字
	Code     string `json:"code" form:"code" gorm:"comment:'字段编码'"`             //字段编码
	Op0Value string `json:"op0Value" gorm:"comment:'初审内容'"`                     //1码内容
	Op0Input string `json:"op0Input" gorm:"comment:'初审录入状态(yes|no|no_if|ocr)'"` //1码录入状态(yes|no|no_if|ocr)
	Op1Value string `json:"op1Value" gorm:"comment:'1码内容'"`                     //1码内容
	Op1Input string `json:"op1Input" gorm:"comment:'1码录入状态(yes|no|no_if|ocr)'"` //1码录入状态(yes|no|no_if|ocr)
	Op2Value string `json:"op2Value" gorm:"comment:'2码内容'"`                     //2码内容
	Op2Input string `json:"op2Input" gorm:"comment:'2码录入状态'"`                   //2码录入状态
	OpqValue string `json:"opqValue" gorm:"comment:'问题件内容'"`                    //问题件内容
	OpqInput string `json:"opqInput" gorm:"comment:'问题件录入状态'"`                  //问题件录入状态
}
