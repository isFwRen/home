package request

import (
	"github.com/lib/pq"
	//"time"
)

// 开始考核
type ReqStartExam struct {
	ProjectCode string `json:"projectCode"` //项目编码
	//StartTime   time.Time `json:"startTime"`   //开始时间
}

type ReqEndExam struct {
	AssessmentProblemId string `json:"assessmentProblemId"` //考核试题id
	//EndTime             time.Time   `json:"endTime"`             //结束时间
	AnswerList []ReqAnswer `json:"answerList"` //回答数组
}

type ReqAnswer struct {
	Id          string   `json:"assessmentSingleId"` //单题id
	BelongBlock string   `json:"belongBlock"`        //所属分块
	ProblemType int      `json:"problemType"`        //题目类型 1.填空 2.单选 3.多选
	IsRequire   int      `json:"isRequire"`          //是否必填
	OptionList  []string `json:"optionList"`         //选项mak
	SetScore    int      `json:"setScore"`           //题目分值"
}

type AnswerFromDB struct {
	Id          string         `json:"id"`
	BelongBlock string         `json:"belongBlock"`
	ProblemType int            `json:"problemType"`
	IsRequire   int            `json:"isRequire"`
	Answer      string         `json:"answer"`
	AnswerList  pq.StringArray `json:"answerList" gorm:"type:text[]"`
}

type UserAnswer struct {
	Id             string         `json:"id"`
	BelongBlock    string         `json:"belongBlock"`
	ProblemType    int            `json:"problemType"`
	IsRequire      int            `json:"isRequire"`
	Answer         string         `json:"answer"`
	AnswerList     pq.StringArray `json:"answerList" gorm:"type:text[]"`
	UserAnswerList []string       `json:"userAnswerList"`
}
