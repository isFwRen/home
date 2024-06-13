package response

type CorrectedResponse struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type UpdateOutputResponse struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type PageResult struct {
	List      interface{} `json:"list"`
	Total     int64       `json:"total"`
	Top       []string    `json:"top"`
	PageIndex int         `json:"pageIndex"`
	PageSize  int         `json:"pageSize"`
}

//----------------------------------------------------------------产量明细-----------------------------------------------------------------------------//

type OutputStatisticsRes struct {
	//ProCode    string `json:"proCode" gorm:"comment:'项目编码'" excel:"项目编码"`
	SubmitTime string `json:"submitTime" gorm:"comment:'提交时间'" excel:"日期"`
	Code       string `json:"code" gorm:"comment:'工号'" excel:"工号"`
	NickName   string `json:"nickName" gorm:"comment:'姓名'" excel:"姓名"`

	//------------------------------------------------------------------------------------汇总----------------------------------------------------//
	// SummarySubmitTime              string  `json:"summarySubmitTime" gorm:"comment:'提交时间'" excel:"汇总-提交时间"`
	// SummaryTime                    string  `json:"summaryTime" gorm:"comment:'总时间'" excel:"汇总-总时间"`
	SummaryFieldCharacter          int64   `json:"summaryFieldCharacter" gorm:"comment:'总字符总量'" excel:"字符总量"`
	SummaryFieldEffectiveCharacter int64   `json:"summaryFieldEffectiveCharacter" gorm:"comment:'总字符有效总量'" excel:"有效字符总量"`
	SummaryAccuracyRate            string  `json:"summaryAccuracyRate" gorm:"comment:'总准确率'" excel:"准确率"`
	SummaryCostTime                string  `json:"summaryCostTime" gorm:"comment:'花费时间'" excel:"时间"`
	SummaryBlockNum                int64   `json:"summaryBlockNum" gorm:"comment:'总分块数量'" excel:"分块数量"`
	SummaryBlockEfficiency         float64 `json:"summaryBlockEfficiency" gorm:"comment:'总分块效率'" excel:"分块效率"`
	SummaryFieldEfficiency         float64 `json:"summaryFieldEfficiency" gorm:"comment:'总字符效率'" excel:"字符效率"`
	SummaryQuestionMarkNumber      int64   `json:"summaryQuestionMarkNumber" gorm:"comment:'总录入?数量'" excel:"录入?数量"`
	SummaryQuestionMarkProportion  string  `json:"summaryQuestionMarkProportion" gorm:"comment:'总录入?比例'" excel:"录入?比例"`

	SummaryFieldNum int64 `json:"summaryFieldNum" gorm:"comment:'字段数量'" `
	//------------------------------------------------------------------------------------初审----------------------------------------------------//
	// Op0SubmitTime              string  `json:"op0submitTime" gorm:"comment:'提交时间'" excel:"初审提交时间"`
	// Op0Time                    string  `json:"op0Time" gorm:"comment:'初审时间'" excel:"初审时间"`
	Op0InvoiceNum      int64   `json:"op0InvoiceNum" gorm:"comment:'初审发票数量'" excel:"发票数量"`
	Op0BlockNum        int64   `json:"op0BlockNum" gorm:"comment:'初审分块数量'" excel:"分块数量"`
	Op0BlockEfficiency float64 `json:"op0BlockEfficiency" gorm:"comment:'初审分块效率'" excel:"分块效率"`
	Op0AccuracyRate    string  `json:"op0AccuracyRate" gorm:"comment:'初审准确率'" excel:"准确率"`
	Op0CostTime        string  `json:"op0CostTime" gorm:"comment:'花费时间'" excel:"时间"`
	Op0FieldNum        int64   `json:"op0FieldNum" gorm:"comment:'字段数量'" excel:"字段数量"`

	Op0FieldCharacter          int64   `json:"op0FieldCharacter" gorm:"comment:'初审字符总量'" `
	Op0FieldEffectiveCharacter int64   `json:"op0FieldEffectiveCharacter" gorm:"comment:'初审有效字符总量'"`
	Op0FieldEfficiency         float64 `json:"op0FieldEfficiency" gorm:"comment:'初审字符效率'" `
	Op0QuestionMarkNumber      int64   `json:"op0QuestionMarkNumber" gorm:"comment:'初审录入?数量'" `
	Op0QuestionMarkProportion  string  `json:"op0QuestionMarkProportion" gorm:"comment:'初审录入?比例'" `
	//------------------------------------------------------------------------------------一码----------------------------------------------------//
	// Op1SubmitTime                               string  `json:"op1submitTime" gorm:"comment:'提交时间'" excel:"一码提交时间"`
	// Op1Time                                     string  `json:"op1Time" gorm:"comment:'一码时间'" excel:"一码时间"`
	Op1FieldCharacter          int64   `json:"op1FieldCharacter" gorm:"comment:'初审字符总量'" excel:"字符总量"`
	Op1FieldEffectiveCharacter int64   `json:"op1FieldEffectiveCharacter" gorm:"comment:'初审有效字符总量'" excel:"有效字符总量"`
	Op1AccuracyRate            string  `json:"op1AccuracyRate" gorm:"comment:'一码准确率'" excel:"准确率"`
	Op1CostTime                string  `json:"op1CostTime" gorm:"comment:'花费时间'" excel:"时间"`
	Op1BlockNum                int64   `json:"op1BlockNum" gorm:"comment:'一码分块数量'" excel:"分块数量"`
	Op1BlockEfficiency         float64 `json:"op1BlockEfficiency" gorm:"comment:'一码分块效率'" excel:"分块效率"`
	Op1FieldEfficiency         float64 `json:"op1FieldEfficiency" gorm:"comment:'一码字符效率'" excel:"字符效率"`
	Op1FieldNum                int64   `json:"op1FieldNum" gorm:"comment:'字段数量'" excel:"字段数量"`
	Op1QuestionMarkNumber      int64   `json:"op1QuestionMarkNumber" gorm:"comment:'一码录入?数量'" excel:"录入?数量"`
	Op1QuestionMarkProportion  string  `json:"op1QuestionMarkProportion" gorm:"comment:'一码录入?比例'" excel:"录入?比例"`

	Op1NotExpenseAccountFieldCharacter          int64 `json:"op1NotExpenseAccountFieldCharacter" gorm:"comment:'一码非报销单字符总量'"`
	Op1NotExpenseAccountFieldEffectiveCharacter int64 `json:"op1NotExpenseAccountFieldEffectiveCharacter" gorm:"comment:'一码非报销单有效字符总量'"`
	Op1ExpenseAccountFieldCharacter             int64 `json:"op1ExpenseAccountFieldCharacter" gorm:"comment:'一码报销单字符总量'" `
	Op1ExpenseAccountFieldEffectiveCharacter    int64 `json:"op1ExpenseAccountFieldEffectiveCharacter" gorm:"comment:'一码报销单有效字符总量'"`
	//------------------------------------------------------------------------------------二码----------------------------------------------------//
	// Op2SubmitTime                               string  `json:"op2submitTime" gorm:"comment:'提交时间'" excel:"二码提交时间"`
	// Op2Time                                     string  `json:"op2Time" gorm:"comment:'二码时间'" excel:"二码时间"`
	Op2FieldCharacter          int64   `json:"op2FieldCharacter" gorm:"comment:'初审字符总量'" excel:"字符总量"`
	Op2FieldEffectiveCharacter int64   `json:"op2FieldEffectiveCharacter" gorm:"comment:'初审有效字符总量'" excel:"有效字符总量"`
	Op2AccuracyRate            string  `json:"op2AccuracyRate" gorm:"comment:'二码准确率'" excel:"准确率"`
	Op2CostTime                string  `json:"op2CostTime" gorm:"comment:'花费时间'" excel:"时间"`
	Op2BlockNum                int64   `json:"op2BlockNum" gorm:"comment:'二码分块数量'" excel:"分块数量"`
	Op2BlockEfficiency         float64 `json:"op2BlockEfficiency" gorm:"comment:'二码分块效率'" excel:"分块效率"`
	Op2FieldEfficiency         float64 `json:"op2FieldEfficiency" gorm:"comment:'二码字符效率'" excel:"字符效率"`
	Op2FieldNum                int64   `json:"op2FieldNum" gorm:"comment:'字段数量'" excel:"字段数量"`
	Op2QuestionMarkNumber      int64   `json:"op2QuestionMarkNumber" gorm:"comment:'二码录入?数量'" excel:"二码录入?数量"`
	Op2QuestionMarkProportion  string  `json:"op2QuestionMarkProportion" gorm:"comment:'二码录入?比例'" excel:"二码录入?比例"`

	Op2NotExpenseAccountFieldCharacter          int64 `json:"op2NotExpenseAccountFieldCharacter" gorm:"comment:'二码非报销单字符总量'" `
	Op2NotExpenseAccountFieldEffectiveCharacter int64 `json:"op2NotExpenseAccountFieldEffectiveCharacter" gorm:"comment:'二码非报销单有效字符总量'"`
	Op2ExpenseAccountFieldCharacter             int64 `json:"op2ExpenseAccountFieldCharacter" gorm:"comment:'二码报销单字符总量'"`
	Op2ExpenseAccountFieldEffectiveCharacter    int64 `json:"op2ExpenseAccountFieldEffectiveCharacter" gorm:"comment:'二码报销单有效字符总量'"`
	//------------------------------------------------------------------------------------问题件----------------------------------------------------//
	// OpQSubmitTime              string  `json:"opQSubmitTime" gorm:"comment:'提交时间'" excel:"问题件提交时间"`
	// OpQTime                    string  `json:"opQTime" gorm:"comment:'问题件时间'" excel:"问题件时间"`
	OpQFieldCharacter          int64   `json:"opQFieldCharacter" gorm:"comment:'问题件字符总量'" excel:"字符总量"`
	OpQFieldEffectiveCharacter int64   `json:"opQFieldEffectiveCharacter" gorm:"comment:'问题件字符有效总量'" excel:"有效字符总量"`
	OpQAccuracyRate            string  `json:"opQAccuracyRate" gorm:"comment:'问题件准确率'" excel:"准确率"`
	OpQCostTime                string  `json:"opQCostTime" gorm:"comment:'花费时间'" excel:"时间"`
	OpQBlockNum                int64   `json:"opQBlockNum" gorm:"comment:'问题件分块数量'" excel:"分块数量"`
	OpQBlockEfficiency         float64 `json:"opQBlockEfficiency" gorm:"comment:'问题件分块效率'" excel:"分块效率"`
	OpQFieldEfficiency         float64 `json:"opQFieldEfficiency" gorm:"comment:'问题件字符效率'" excel:"字符效率"`
	OpQFieldNum                int64   `json:"opQFieldNum" gorm:"comment:'字段数量'" excel:"字段数量"`
	OpQQuestionMarkNumber      int64   `json:"opQQuestionMarkNumber" gorm:"comment:'问题件录入?数量'" excel:"录入?数量"`
	OpQQuestionMarkProportion  string  `json:"opQQuestionMarkProportion" gorm:"comment:'问题件录入?比例'" excel:"录入?比例"`
}
