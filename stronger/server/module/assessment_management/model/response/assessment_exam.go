package response

import "github.com/lib/pq"

// 开始考核返回结果
type RespStartExam struct {
	AssessmentProblem   string                `json:"assessmentProblem"`   //考核题目id
	AssessmentBlockList []RespAssessmentBlock `json:"assessmentBlockList"` //分块数组
}

type RespAssessmentBlock struct {
	AssessmentBlockData `mapstructure:",squash"`
	//BlockId              string                 `json:"blockId"`              //分块id
	//BlockImgUrl          string                 `json:"blockImgUrl"`          //分块题目图片url
	//Order                int                    `json:"order"`                //板块顺序
	AssessmentSingleList []RespAssessmentSingle `json:"assessmentSingleList" gorm:"foreignKey:BelongBlock;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` //题目数组
}

type RespAssessmentSingle struct {
	Id          string         `json:"assessmentSingleId"`            //单题id
	OrderNum    int            `json:"orderNum"`                      //题目序号
	BelongBlock string         `json:"belongBlock"`                   //所属分块
	Question    string         `json:"question"`                      // 题目
	ProblemType int            `json:"problemType"`                   //题目类型 1.填空 2.单选 3.多选
	IsRequire   int            `json:"isRequire"`                     //是否必填
	OptionList  pq.StringArray `json:"optionList" gorm:"type:text[]"` //选项
	SetScore    int            `json:"setScore"`                      //题目分值"
}

func (RespAssessmentSingle) TableName() string {
	return "assessment_problem_singles"
}

type AssessmentBlockData struct {
	ID       string `json:"id"`
	FilePath string `json:"blockImgUrl"`
	OrderNum string `json:"orderNum"`
}

type RespAssessCriterion struct {
	ProjectCode string `json:"projectCode"` //项目编码
	SetPoint    int    `json:"setPoint"`    //设置的考核标准
}

type SaveUserAssessmentSingle struct {
	Id             string         `json:"assessmentSingleId"`            //单题id
	OrderNum       int            `json:"orderNum"`                      //题目序号
	BelongBlock    string         `json:"belongBlock"`                   //所属分块
	Question       string         `json:"question"`                      // 题目
	ProblemType    int            `json:"problemType"`                   //题目类型 1.填空 2.单选 3.多选
	IsRequire      int            `json:"isRequire"`                     //是否必填
	OptionList     pq.StringArray `json:"optionList" gorm:"type:text[]"` //选项
	SetScore       int            `json:"setScore"`                      //题目分值"
	Answer         string         `json:"answer"`                        //填空答案
	AnswerList     pq.StringArray `json:"answerList" gorm:"type:text[]"` //选择题答案
	FilePath       string         `json:"filePath"`                      //题目图片路径
	UserAnswerList []string       `json:"userAnswerList"`                //用户答案
}
