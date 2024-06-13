package model

import (
	"server/module/sys_base/model"
	"time"
)

// TrainingManagement 培训管理 实体类
type TrainingManagement struct {
	model.Model
	UserCode                 string    `json:"userCode" form:"userCode" binding:"required" excel:"工号"`                                     //工号
	UserName                 string    `json:"userName" form:"userName" binding:"required" excel:"姓名"`                                     //姓名
	Gender                   int       `json:"gender" form:"gender" binding:"required" excel:"性别"`                                         //性别 1男 2女
	Phone                    string    `json:"phone" form:"phone" binding:"required" excel:"手机号"`                                          //手机号
	ProjectCode              string    `json:"projectCode" form:"projectCode" binding:"required" excel:"项目编号"`                             //项目编号
	EntryStartAt             time.Time `json:"entryStartAt" form:"entryStartAt" binding:"required" excel:"入职日期"`                           //入职日期
	DutyAt                   time.Time `json:"dutyAt" form:"dutyAt" binding:"required" excel:"上岗日期"`                                       //上岗日期
	PracticeYieldRequirement int       `json:"practiceYieldRequirement" form:"practiceYieldRequirement" binding:"required" excel:"练习产量要求"` //练习产量要求
	AccuracyRequirement      float64   `json:"accuracyRequirement" form:"accuracyRequirement" binding:"required" excel:"准确率要求"`            //准确率要求
	ActualPracticeYield      int       `json:"actualPracticeYield" form:"actualPracticeYield" binding:"required" excel:"实际练习产量"`           //实际练习产量
	ActualPracticeAccuracy   float64   `json:"actualPracticeAccuracy" form:"actualPracticeAccuracy" binding:"required" excel:"实际练习准确率"`    //实际练习准确率
	Examine                  int       `json:"examine" form:"examine" binding:"required" excel:"考核要求"`                                     //考核要求
	ExamineScore             int       `json:"examineScore" form:"examineScore" binding:"required" excel:"考核分数"`                           //考核分数
	AuditStatus              int       `json:"auditStatus" form:"auditStatus" binding:"required" excel:"审核状态"`                             //审核状态 1.待审核 2.审核通过  3.审核未通过
}
