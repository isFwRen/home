package model

import (
	"github.com/lib/pq"
	//"gorm.io/driver/postgres"
	//"github.com/go-gorm/postgres"
	"server/module/sys_base/model"
	"time"
)

type AssessmentRecord struct {
	model.Model
	UserId     string         `json:"userId"`                      //用户id
	UsedExam   pq.StringArray `json:"usedExam" gorm:"type:text[]"` //考核过的题目
	UserStatus int            `json:"userStatus"`                  //考核转态  1.通过 2.未通过
	LastExamAt time.Time      `json:"lastExamAt"`                  //最后一次考核时间
}

type AssessmentUserAnswer struct {
	model.Model
	ProblemName string                 `json:"problemName"` //考核题名称
	ProblemId   string                 `json:"problemId"`   //考核题id
	UserCode    string                 `json:"userCode"`
	UserName    string                 `json:"userName"`
	ProjectCode string                 `json:"projectCode"`
	Score       string                 `json:"score"`
	AnswerSheet map[string]interface{} `json:"answerSheet" gorm:"type:jsonb"`
	Standard    int                    `json:"standard"`
	IsPass      int                    `json:"isPass"`
}

// 考核标准
type AssessmentCriterion struct {
	model.Model
	ProjectCode string `json:"projectCode" from:"projectCode"` //项目ID
	SetPoint    int    `json:"setPoint" from:"setPoint"`       // 考核标准分数
}
