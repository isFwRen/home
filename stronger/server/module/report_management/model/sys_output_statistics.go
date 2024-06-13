package model

import (
	"server/module/sys_base/model"
	"time"
)

//----------------------------------------------------------------产量明细-----------------------------------------------------------------------------//

type OutputStatistics struct {
	model.Model
	//ProCode    string    `json:"proCode" gorm:"comment:'项目编码'" excel:"项目编码"`
	Code       string    `json:"code" gorm:"comment:'工号'" excel:"工号"`
	NickName   string    `json:"nickName" gorm:"comment:'姓名'" excel:"姓名"`
	SubmitTime time.Time `json:"submitTime" gorm:"comment:'提交时间'" excel:"提交时间"`
	//------------------------------------------------------------------------------------初审----------------------------------------------------//
	// Op0SubmitTime              string  `json:"op0submitTime" gorm:"comment:'提交时间'" excel:"初审提交时间"`
	// Op0Time                    string  `json:"op0Time" gorm:"comment:'初审时间'" excel:"初审时间"`
	Op0CostTime       int64 `json:"op0CostTime" gorm:"comment:'初审花费时间'" excel:"初审花费时间"`
	Op0FieldNum       int   `json:"op0FieldNum" gorm:"comment:'初审字段数量'" excel:"初审字段数量"`
	Op0FieldCharacter int   `json:"op0FieldCharacter" gorm:"comment:'初审字符总量'" excel:"初审字符总量"`

	Op0FieldEffectiveCharacter int     `json:"op0FieldEffectiveCharacter" gorm:"comment:'初审有效字符总量'" excel:"初审有效字符总量"`
	Op0AccuracyRate            float64 `json:"op0AccuracyRate" gorm:"comment:'初审准确率'" excel:"初审准确率"`
	Op0InvoiceNum              int     `json:"op0InvoiceNum" gorm:"comment:'初审发票数量'" excel:"初审发票数量"`
	Op0BlockNum                int     `json:"op0BlockNum" gorm:"comment:'初审分块数量'" excel:"初审分块数量"`
	Op0BlockEfficiency         float64 `json:"op0BlockEfficiency" gorm:"comment:'初审分块效率'" excel:"初审分块效率"`
	Op0FieldEfficiency         float64 `json:"op0FieldEfficiency" gorm:"comment:'初审字符效率'" excel:"初审字符效率"`
	Op0QuestionMarkNumber      int     `json:"op0QuestionMarkNumber" gorm:"comment:'初审录入?数量'" excel:"初审录入?数量"`
	Op0QuestionMarkProportion  float64 `json:"op0QuestionMarkProportion" gorm:"comment:'初审录入?比例'" excel:"初审录入?比例"`
	//------------------------------------------------------------------------------------一码----------------------------------------------------//
	// Op1SubmitTime                               string  `json:"op1submitTime" gorm:"comment:'提交时间'" excel:"一码提交时间"`
	// Op1Time                                     string  `json:"op1Time" gorm:"comment:'一码时间'" excel:"一码时间"`
	Op1CostTime                int64 `json:"op1costTime" gorm:"comment:'一码花费时间'" excel:"一码花费时间"`
	Op1FieldNum                int   `json:"op1FieldNum" gorm:"comment:'一码字段数量'" excel:"一码字段数量"`
	Op1FieldCharacter          int   `json:"op1FieldCharacter" gorm:"comment:'一码字符总量'" excel:"一码字符总量"`
	Op1FieldEffectiveCharacter int   `json:"op1FieldEffectiveCharacter" gorm:"comment:'一码有效字符总量'" excel:"一码有效字符总量"`

	Op1NotExpenseAccountFieldCharacter          int     `json:"op1NotExpenseAccountFieldCharacter" gorm:"comment:'一码非报销单字符总量'" excel:"一码非报销单字符总量"`
	Op1NotExpenseAccountFieldEffectiveCharacter int     `json:"op1NotExpenseAccountFieldEffectiveCharacter" gorm:"comment:'一码非报销单有效字符总量'" excel:"一码非报销单有效字符总量"`
	Op1ExpenseAccountFieldCharacter             int     `json:"op1ExpenseAccountFieldCharacter" gorm:"comment:'一码报销单字符总量'" excel:"一码报销单字符总量"`
	Op1ExpenseAccountFieldEffectiveCharacter    int     `json:"op1ExpenseAccountFieldEffectiveCharacter" gorm:"comment:'一码报销单有效字符总量'" excel:"一码报销单有效字符总量"`
	Op1AccuracyRate                             float64 `json:"op1AccuracyRate" gorm:"comment:'一码准确率'" excel:"一码准确率"`
	Op1BlockNum                                 int     `json:"op1BlockNum" gorm:"comment:'一码分块数量'" excel:"一码分块数量"`
	Op1BlockEfficiency                          float64 `json:"op1BlockEfficiency" gorm:"comment:'一码分块效率'" excel:"一码分块效率"`
	Op1FieldEfficiency                          float64 `json:"op1FieldEfficiency" gorm:"comment:'一码字符效率'" excel:"一码字符效率"`
	Op1QuestionMarkNumber                       int     `json:"op1QuestionMarkNumber" gorm:"comment:'一码录入?数量'" excel:"一码录入?数量"`
	Op1QuestionMarkProportion                   float64 `json:"op1QuestionMarkProportion" gorm:"comment:'一码录入?比例'" excel:"一码录入?比例"`
	//------------------------------------------------------------------------------------二码----------------------------------------------------//
	// Op2SubmitTime                               string  `json:"op2submitTime" gorm:"comment:'提交时间'" excel:"二码提交时间"`
	// Op2Time                                     string  `json:"op2Time" gorm:"comment:'二码时间'" excel:"二码时间"`
	Op2CostTime                int64 `json:"op2costTime" gorm:"comment:'二码花费时间'" excel:"二码花费时间"`
	Op2FieldNum                int   `json:"op2FieldNum" gorm:"comment:'二码字段数量'" excel:"二码字段数量"`
	Op2FieldCharacter          int   `json:"op2FieldCharacter" gorm:"comment:'二码字符总量'" excel:"二码字符总量"`
	Op2FieldEffectiveCharacter int   `json:"op2FieldEffectiveCharacter" gorm:"comment:'二码有效字符总量'" excel:"二码有效字符总量"`

	Op2NotExpenseAccountFieldCharacter          int     `json:"op2NotExpenseAccountFieldCharacter" gorm:"comment:'二码非报销单字符总量'" excel:"二码非报销单字符总量"`
	Op2NotExpenseAccountFieldEffectiveCharacter int     `json:"op2NotExpenseAccountFieldEffectiveCharacter" gorm:"comment:'二码非报销单有效字符总量'" excel:"二码非报销单有效字符总量"`
	Op2ExpenseAccountFieldCharacter             int     `json:"op2ExpenseAccountFieldCharacter" gorm:"comment:'二码报销单字符总量'" excel:"二码报销单字符总量"`
	Op2ExpenseAccountFieldEffectiveCharacter    int     `json:"op2ExpenseAccountFieldEffectiveCharacter" gorm:"comment:'二码报销单有效字符总量'" excel:"二码报销单有效字符总量"`
	Op2AccuracyRate                             float64 `json:"op2AccuracyRate" gorm:"comment:'二码准确率'" excel:"二码准确率"`
	Op2BlockNum                                 int     `json:"op2BlockNum" gorm:"comment:'二码分块数量'" excel:"二码分块数量"`
	Op2BlockEfficiency                          float64 `json:"op2BlockEfficiency" gorm:"comment:'二码分块效率'" excel:"二码分块效率"`
	Op2FieldEfficiency                          float64 `json:"op2FieldEfficiency" gorm:"comment:'二码字符效率'" excel:"二码字符效率"`
	Op2QuestionMarkNumber                       int     `json:"op2QuestionMarkNumber" gorm:"comment:'二码录入?数量'" excel:"二码录入?数量"`
	Op2QuestionMarkProportion                   float64 `json:"op2QuestionMarkProportion" gorm:"comment:'二码录入?比例'" excel:"二码录入?比例"`
	//------------------------------------------------------------------------------------问题件----------------------------------------------------//
	// OpQSubmitTime              string  `json:"opQSubmitTime" gorm:"comment:'提交时间'" excel:"问题件提交时间"`
	// OpQTime                    string  `json:"opQTime" gorm:"comment:'问题件时间'" excel:"问题件时间"`
	OpQCostTime int64 `json:"opQCostTime" gorm:"comment:'问题件花费时间'" excel:"问题件花费时间"`
	OpQFieldNum int   `json:"opQFieldNum" gorm:"comment:'问题件字段数量'" excel:"问题件字段数量"`

	OpQFieldCharacter          int     `json:"opQFieldCharacter" gorm:"comment:'问题件字符总量'" excel:"问题件字符总量"`
	OpQFieldEffectiveCharacter int     `json:"opQFieldEffectiveCharacter" gorm:"comment:'问题件字符有效总量'" excel:"问题件字符有效总量"`
	OpQAccuracyRate            float64 `json:"opQAccuracyRate" gorm:"comment:'问题件准确率'" excel:"问题件准确率"`
	OpQBlockNum                int     `json:"opQBlockNum" gorm:"comment:'问题件分块数量'" excel:"问题件分块数量"`
	OpQBlockEfficiency         float64 `json:"opQBlockEfficiency" gorm:"comment:'问题件分块效率'" excel:"问题件分块效率"`
	OpQFieldEfficiency         float64 `json:"opQFieldEfficiency" gorm:"comment:'问题件字符效率'" excel:"问题件字符效率"`
	OpQQuestionMarkNumber      int     `json:"opQQuestionMarkNumber" gorm:"comment:'问题件录入?数量'" excel:"问题件录入?数量"`
	OpQQuestionMarkProportion  float64 `json:"opQQuestionMarkProportion" gorm:"comment:'问题件录入?比例'" excel:"问题件录入?比例"`
	//------------------------------------------------------------------------------------汇总----------------------------------------------------//
	// SummarySubmitTime              string  `json:"summarySubmitTime" gorm:"comment:'提交时间'" excel:"汇总-提交时间"`
	// SummaryTime                    string  `json:"summaryTime" gorm:"comment:'总时间'" excel:"汇总-总时间"`

	SummaryFieldNum int   `json:"summaryFieldNum" gorm:"comment:'汇总-字段数量'" excel:"汇总-字段数量"`
	SummaryCostTime int64 `json:"summaryCostTime" gorm:"comment:'汇总-花费时间'" excel:"汇总-花费时间"`

	SummaryFieldCharacter          int     `json:"summaryFieldCharacter" gorm:"comment:'总字符总量'" excel:"汇总-总字符总量"`
	SummaryFieldEffectiveCharacter int     `json:"summaryFieldEffectiveCharacter" gorm:"comment:'总字符有效总量'" excel:"汇总-总字符有效总量"`
	SummaryAccuracyRate            float64 `json:"summaryAccuracyRate" gorm:"comment:'总准确率'" excel:"汇总-总准确率"`
	SummaryBlockNum                int     `json:"summaryBlockNum" gorm:"comment:'总分块数量'" excel:"汇总-总分块数量"`
	SummaryBlockEfficiency         float64 `json:"summaryBlockEfficiency" gorm:"comment:'总分块效率'" excel:"汇总-总分块效率"`
	SummaryFieldEfficiency         float64 `json:"summaryFieldEfficiency" gorm:"comment:'总字符效率'" excel:"汇总-总字符效率"`
	SummaryQuestionMarkNumber      int     `json:"summaryQuestionMarkNumber" gorm:"comment:'总录入?数量'" excel:"汇总-总录入?数量"`
	SummaryQuestionMarkProportion  float64 `json:"summaryQuestionMarkProportion" gorm:"comment:'总录入?比例'" excel:"汇总-总录入?比例"`
	//判断是否处理过有效字符：
	ProcessedValidCharacters bool `json:"processedValidCharacters" gorm:"comment:'是否处理过有效字符'"`
}

//-----------------------------------------------------------产量全部-------------------------------------------------------------------------------------------//

type OutputStatisticsSummary struct {
	model.Model
	Code                 string    `json:"code" gorm:"comment:'工号'"`
	NickName             string    `json:"nickName" gorm:"comment:'姓名'"`
	SubmitTime           time.Time `json:"submitTime" gorm:"comment:'提交时间'"`
	ProCode              string    `json:"proCode" gorm:"comment:'项目编码'"`
	Mary                 int       `json:"mary" gorm:"comment:'汇总'"`
	Op0                  int       `json:"op0" gorm:"comment:'初审'"` //有效字符
	Op0InvoiceNum        int       `json:"op0InvoiceNum" gorm:"comment:'初审发票数量'"`
	Op1NotExpenseAccount int       `json:"op1NotExpenseAccount" gorm:"comment:'一码非报销单'"` //有效字符
	Op1ExpenseAccount    int       `json:"op1ExpenseAccount" gorm:"comment:'一码报销单'"`     //有效字符
	Op2NotExpenseAccount int       `json:"op2NotExpenseAccount" gorm:"comment:'二码非报销单'"` //有效字符
	Op2ExpenseAccount    int       `json:"op2ExpenseAccount" gorm:"comment:'二码报销单'"`     //有效字符
	Question             int       `json:"question" gorm:"comment:'问题件'"`                //有效字符
}

type OutputStatisticsSummaryList struct {
	Code             string    `json:"code" gorm:"comment:'工号'"`
	NickName         string    `json:"nickName" gorm:"comment:'姓名'"`
	SubmitTime       string    `json:"submitTime" gorm:"comment:'提交时间'"`
	AddUpToSomething float64   `json:"addUpToSomething"` //合计
	ProSummary       []Summary `json:"proSummary"`
}

type Summary struct {
	ProCode              string  `json:"proCode" gorm:"comment:'项目编码'"`
	Mary                 float64 `json:"mary" gorm:"comment:'汇总'"`
	Op0                  float64 `json:"op0" gorm:"comment:'初审'"`
	Op0InvoiceNum        float64 `json:"op0InvoiceNum" gorm:"comment:'初审发票数量'"`
	Op1NotExpenseAccount float64 `json:"op1NotExpenseAccount" gorm:"comment:'一码非报销单'"`
	Op1ExpenseAccount    float64 `json:"op1ExpenseAccount" gorm:"comment:'一码报销单'"`
	Op2NotExpenseAccount float64 `json:"op2NotExpenseAccount" gorm:"comment:'二码非报销单'"`
	Op2ExpenseAccount    float64 `json:"op2ExpenseAccount" gorm:"comment:'二码报销单'"`
	Question             float64 `json:"question" gorm:"comment:'问题件'"`
}

//--------------------------------------------------------------------------------项目折算比例-------------------------------------------------------------------//

type SysCorrected struct {
	model.Model
	ProCode              string                `json:"proCode" gorm:"comment:'项目编码'" editName:"项目编码"`
	Op0AsTheBlock        float64               `json:"op0AsTheBlock" gorm:"comment:'一码(按分块)'" editName:"一码(按分块)"`
	Op0AsTheInvoice      float64               `json:"op0AsTheInvoice" gorm:"comment:'一码(按发票)'" editName:"一码(按发票)"`
	Op1NotExpenseAccount float64               `json:"op1NotExpenseAccount" gorm:"comment:'一码非报销单'" editName:"一码非报销单"`
	Op1ExpenseAccount    float64               `json:"op1ExpenseAccount" gorm:"comment:'一码报销单'" editName:"一码报销单"`
	Op2NotExpenseAccount float64               `json:"op2NotExpenseAccount" gorm:"comment:'二码非报销单'" editName:"二码非报销单"`
	Op2ExpenseAccount    float64               `json:"op2ExpenseAccount" gorm:"comment:'二码报销单'" editName:"二码报销单"`
	Question             float64               `json:"question" gorm:"comment:'问题件'" editName:"问题件"`
	StartTime            time.Time             `json:"startTime" gorm:"comment:'开始时间'" editName:"开始时间"`
	Log                  []SysCorrectedEditLog `json:"log" gorm:"foreignKey:CorrectedID;comment:'日志'"`
}

type SysCorrectedRep struct {
	model.Model
	ProCode              string                `json:"proCode" gorm:"comment:'项目编码'"`
	Op0AsTheBlock        float64               `json:"op0AsTheBlock" gorm:"comment:'一码(按分块)'"`
	Op0AsTheInvoice      float64               `json:"op0AsTheInvoice" gorm:"comment:'一码(按发票)'"`
	Op1NotExpenseAccount float64               `json:"op1NotExpenseAccount" gorm:"comment:'一码非报销单'"`
	Op1ExpenseAccount    float64               `json:"op1ExpenseAccount" gorm:"comment:'一码报销单'"`
	Op2NotExpenseAccount float64               `json:"op2NotExpenseAccount" gorm:"comment:'二码非报销单'"`
	Op2ExpenseAccount    float64               `json:"op2ExpenseAccount" gorm:"comment:'二码报销单'"`
	Question             float64               `json:"question" gorm:"comment:'问题件'"`
	StartTime            string                `json:"startTime" gorm:"comment:'开始时间'"`
	Log                  []SysCorrectedEditLog `json:"log" gorm:"foreignKey:CorrectedID;comment:'日志'"`
}

type SysCorrectedEditLog struct {
	model.Model
	CorrectedID string `json:"correctedID" gorm:"comment:'修改比例日志id'"`
	EditName    string `json:"editName" gorm:"comment:'名称'"`
	BeforeEdit  string `json:"beforeEdit" gorm:"comment:'修改前内容'"`
	AfterEdit   string `json:"afterEdit" gorm:"comment:'修改后内容'"`
	EditedCode  string `json:"editedCode" gorm:"comment:'修改人工号'"`
	EditedName  string `json:"editedName" gorm:"comment:'修改人姓名'"`
}

type SysCorrectedEditLogReq struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	CorrectedID string `json:"correctedID" gorm:"comment:'修改比例日志id'"`
	EditName    string `json:"editName" gorm:"comment:'名称'"`
	BeforeEdit  string `json:"beforeEdit" gorm:"comment:'修改前内容'"`
	AfterEdit   string `json:"afterEdit" gorm:"comment:'修改后内容'"`
	EditedCode  string `json:"editedCode" gorm:"comment:'修改人工号'"`
	EditedName  string `json:"editedName" gorm:"comment:'修改人姓名'"`
}
