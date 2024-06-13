package service

import (
	"fmt"
	"server/module/assessment_management/model"
	practiceModel "server/module/practice/model"
	model2 "server/module/sys_base/model"
	trainingModule "server/module/training_management/model"
	"unsafe"
	model3 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/model"

	//"golang.org/x/image/colornames"
	"math/rand"
	"server/global"
	reqModel "server/module/assessment_management/model/request"
	respModel "server/module/assessment_management/model/response"
	"time"
)

// 开始考核 通过考核项目编码随机抽取考核题目
func StartExam(userId string, exam reqModel.ReqStartExam) (respExam respModel.RespStartExam, err error) {
	var problemIdList []string
	//获取正在使用的试题
	err = global.GDb.Table("assessment_problems").
		Where("project_code = ? and problem_status = 1 and assessment_criterion != 0", exam.ProjectCode).
		Select("id").Scan(&problemIdList).Error
	isExistRecord := int64(0)
	err = global.GDb.Table("assessment_records").Where("user_id = ?", userId).Count(&isExistRecord).Error
	if isExistRecord == 0 { //第一次参与考核
		assessmentRecord := model.AssessmentRecord{
			UserId:     userId,
			UsedExam:   []string{},
			UserStatus: 2,
			LastExamAt: time.Now(),
		}
		if err = global.GDb.Create(&assessmentRecord).Error; err != nil {
			return
		}
	}
	if err != nil {
		return
	}
	var usedExamList []string
	if err = global.GDb.Table("assessment_records").Where("user_id = ?", userId).
		Select("used_exam").Scan(&usedExamList).Error; err != nil {
		return
	}
	usedExamMap := make(map[string]int)
	for _, s := range usedExamList {
		usedExamMap[s] = 1
	}
	if len(problemIdList) > 0 {
		//随机试题
		rand.Seed(time.Now().UnixNano())
		problemOrder := rand.Intn(len(problemIdList))
		for _, ok := usedExamMap[problemIdList[problemOrder]]; ok; { //如果试题存在则循环，直到出现不存在
			problemOrder = rand.Intn(len(problemIdList))
		}
		usedExamList = append(usedExamList, problemIdList[problemOrder])

		//开始创建题目
		respExam.AssessmentProblem = problemIdList[problemOrder]

		// -------------------不使用连表时的老代码---------------------------
		//var blockList []respModel.RespAssessmentBlock
		//var blockDataList []respModel.AssessmentBlockData
		//global.GDb.Table("assessment_problem_blocks").
		//	Where("belong_problem = ?", problemIdList[problemOrder]).
		//	Select("id", "order_num", "file_Path").Order("order_num ASC").Scan(&blockDataList)
		//
		//for _, data := range blockDataList {
		//	var singleList []respModel.RespAssessmentSingle
		//	err = global.GDb.Table("assessment_problem_singles").Where("belong_block = ?", data.ID).
		//		Select("id", "question", "problem_type", "belong_block", "order_num",
		//			"is_require", "set_score", "option_list").Order("order_num ASC").Scan(&singleList).Error
		//	if err != nil {
		//		return
		//	}
		//	blockList = append(blockList, respModel.RespAssessmentBlock{
		//		AssessmentBlockData:  data,
		//		AssessmentSingleList: singleList,
		//	})
		//}
		// -------------------不使用连表时的老代码---------------------------

		//---------------连表获取单个题目------------------------------d//
		blockDataListPreload := []respModel.RespAssessmentBlock{}
		err = global.GDb.Table("assessment_problem_blocks").
			Where("belong_problem = ?", problemIdList[problemOrder]).
			Order("order_num ASC").
			Preload("AssessmentSingleList").
			Find(&blockDataListPreload).Error
		//testList = blockDataListPreload

		//---------------连表获取单个题目------------------------------d//

		respExam.AssessmentBlockList = blockDataListPreload
	}
	//返回题目
	return
}

// 结束考核 计算考核分数
//func EndExam(userId string, exam reqModel.ReqEndExam) (isPass bool, pointStr string, err error) {
//	blockMap := make(map[string]int)                ///map[分块id] 分块题目数量
//	singlesMap := make(map[string]map[int][]string) //map[题目id][题目类型]用户答案
//	var singleIdList []string
//	// 初始化map
//	for _, answer := range exam.AnswerList {
//		singleIdList = append(singleIdList, answer.Id)
//		_, ok := blockMap[answer.BelongBlock]
//		if ok {
//			blockMap[answer.BelongBlock] += 1
//		} else {
//			blockMap[answer.BelongBlock] = 1
//		}
//		//将题目放入map
//		singlesMap[answer.Id] = make(map[int][]string)
//		singlesMap[answer.Id][answer.ProblemType] = answer.OptionList
//		//answer.BelongBlock
//	}
//	var answerList []reqModel.AnswerFromDB
//	err = global.GDb.Table("assessment_problem_singles").Where("id in (?)", singleIdList).
//		Select("id", "belong_block", "problem_type", "is_require", "answer", "answer_list").
//		Scan(&answerList).Error
//	if err != nil {
//		return
//	}
//
//	blockTMap := make(map[string]int) ///map[分块id] 分块内回答正确的题目数量
//	for _, ans := range answerList {  //判断回答对错
//		strReq, ok := singlesMap[ans.Id][ans.ProblemType]
//
//		if _, e := blockTMap[ans.BelongBlock]; !e { //
//			blockTMap[ans.BelongBlock] = 0
//		}
//		if ok {
//			var strDB []string
//			switch ans.ProblemType {
//			case 1: //填空题
//				strDB = append(strDB, ans.Answer)
//			case 2: //单选题
//				strDB = append(strDB, ans.AnswerList...)
//			case 3: //多选题
//				strDB = append(strDB, ans.AnswerList...)
//			default:
//				return
//			}
//			if CompStringList(strReq, strDB) { //回答正确
//				blockTMap[ans.BelongBlock] += 1
//			} else {
//				if ans.IsRequire == 2 { //非必填项
//					blockTMap[ans.BelongBlock] += 1
//				}
//			}
//		}
//	}
//
//	rateTureBlock := 0.0 //所有分块总正确率
//	for s, i := range blockMap {
//		v, ok := blockTMap[s]
//		if ok {
//			rateTureBlock += float64(v) / float64(i)
//		} else {
//			return
//		}
//	}
//	var criterion int
//	if err = global.GDb.Table("assessment_problems").Select("assessment_criterion").
//		Where("id = ?", exam.AssessmentProblemId).Scan(&criterion).Error; err != nil {
//		return
//	}
//
//	point := ((rateTureBlock) / float64(len(blockMap))) * 100
//	isPass = int(point) >= criterion
//	pointStr = fmt.Sprintf("%.2f", point)
//
//	status := 2
//	if isPass {
//		status = 1
//	}
//	//计入考核记录
//	err = global.GDb.Table("assessment_records").Where("user_id = ?", userId).Updates(map[string]interface{}{
//		"user_status":  status,
//		"last_exam_at": time.Now(),
//	}).Error
//
//	//将试卷存入数据库
//	respExam := make(map[string]interface{})
//	err = CreatedExam(&respExam, singlesMap, exam.AssessmentProblemId)
//	if err != nil {
//		return
//	}
//	var assessmentUserAnswer model.AssessmentUserAnswer
//
//	assessmentUserAnswer.AnswerSheet = respExam
//	if err = global.GDb.Table("assessment_problems").Select("project_code").
//		Where("id = ?", exam.AssessmentProblemId).Scan(&assessmentUserAnswer.ProjectCode).Error; err != nil {
//		return
//	}
//
//	if err = global.GDb.Table("assessment_problems").Select("problem_Name").
//		Where("id = ?", exam.AssessmentProblemId).Scan(&assessmentUserAnswer.ProblemName).Error; err != nil {
//		return
//	}
//	assessmentUserAnswer.ProblemId = exam.AssessmentProblemId
//	assessmentUserAnswer.Score = pointStr
//	assessmentUserAnswer.Standard = criterion
//	assessmentUserAnswer.IsPass = status
//	assessmentUserAnswer.UserName = userId //******
//	assessmentUserAnswer.UserCode = userId //***********
//	if err = global.GDb.Create(&assessmentUserAnswer).Error; err != nil {
//		return
//	}
//	return
//}

// 判断考核时间是否已经过期
func IsOverDue(userId string) (bool, error) {
	var lastExamAt time.Time
	if err := global.GDb.Table("assessment_records").Select("last_exam_at").Where("user_id = ?", userId).Scan(&lastExamAt).Error; err != nil {
		return false, err
	}
	h, _ := time.ParseDuration("2h")
	lastExamAt = lastExamAt.Add(h)
	endTime := time.Now()
	if endTime.Before(lastExamAt) || endTime.Equal(lastExamAt) {
		return true, nil
	}
	return false, nil
}

func GetAssessmentCriterion(userInfo *model3.CustomClaims) (respList []respModel.RespAssessCriterion, err error) {
	var projectCodeList []string
	err = global.GDb.Table("sys_pro_permissions").
		Select("pro_code").
		Where("(has_op0 or has_op1 or has_op2 or has_opq or has_pm) and sys_pro_permissions.user_id = ?", userInfo.ID).
		Scan(&projectCodeList).Error
	if len(projectCodeList) != 0 {
		err = global.GDb.Table("assessment_criterions").Select("project_code", "set_point").
			Where("project_code in (?)", projectCodeList).Scan(&respList).Error
	}
	return
}

// 判断两个字符串数组元素是否相同，不比较元素顺序
func CompStringList(strA []string, strB []string) bool {
	strAMap := make(map[string]int)
	if len(strA) == len(strB) {
		for _, s := range strA {
			strAMap[s] = 1
		}
		for _, s := range strB {
			_, ok := strAMap[s]
			if ok {

			} else {
				return false
			}
		}
		return true
	}
	return false
}

// 创建试题
func CreatedExam(respExam *map[string]interface{}, singlesMap map[string]map[int][]string, problemId string) (err error) {
	(*respExam)["assessmentProblem"] = problemId
	var blockList []interface{}

	var blockDataList []respModel.AssessmentBlockData
	global.GDb.Table("assessment_problem_blocks").
		Where("belong_problem = ?", problemId).
		Select("id", "order_num", "file_Path").Order("order_num ASC").Scan(&blockDataList)

	for _, data := range blockDataList {
		var singleList []respModel.SaveUserAssessmentSingle
		err = global.GDb.Table("assessment_problem_singles").Where("belong_block = ?", data.ID).
			Select("id", "question", "problem_type", "belong_block", "order_num",
				"is_require", "set_score", "option_list", "answer", "answer_list").Order("order_num ASC").Scan(&singleList).Error
		if err != nil {
			return
		}

		for i, single := range singleList {
			strList, ok := singlesMap[single.Id][single.ProblemType]
			if ok {
				singleList[i].UserAnswerList = strList
				//err = global.GDb.Table("assessment_blocks").Select("file_path").
				//	Where("id = ?",single.BelongBlock).Scan(&singleList[i].FilePath).Error
				//if err!=nil {
				//	return
				//}
			}
		}
		blockListMap := make(map[string]interface{})
		blockListMap["assessmentBlockData"] = data
		blockListMap["assessmentSingleList"] = singleList
		blockList = append(blockList, blockListMap)
	}

	(*respExam)["assessmentBlockList"] = blockList
	return err
}

// 结束考核 计算考核分数
func EndExamScore(userInfo *model3.CustomClaims, exam reqModel.ReqEndExam) (isPass bool, pointStr string, err error) {
	blockMap := make(map[string]int)                ///map[分块id] 分块题目数量
	singlesMap := make(map[string]map[int][]string) //map[题目id][题目类型]用户答案
	var singleIdList []string
	// 初始化map
	for _, answer := range exam.AnswerList {
		singleIdList = append(singleIdList, answer.Id)
		_, ok := blockMap[answer.BelongBlock]
		if ok {
			blockMap[answer.BelongBlock] += 1
		} else {
			blockMap[answer.BelongBlock] = 1
		}
		//将题目放入map
		singlesMap[answer.Id] = make(map[int][]string)
		singlesMap[answer.Id][answer.ProblemType] = answer.OptionList
		//answer.BelongBlock
	}
	var answerList []reqModel.AnswerFromDB
	err = global.GDb.Table("assessment_problem_singles").Where("id in (?)", singleIdList).
		Select("id", "belong_block", "problem_type", "is_require", "answer", "answer_list").
		Scan(&answerList).Error
	if err != nil {
		return
	}

	blockTMap := make(map[string]int) ///map[分块id] 分块内回答正确的题目数量
	var ansIdList []string            //记录回到正确的题目id
	for _, ans := range answerList {  //判断回答对错
		strReq, ok := singlesMap[ans.Id][ans.ProblemType]

		if _, e := blockTMap[ans.BelongBlock]; !e { //
			blockTMap[ans.BelongBlock] = 0
		}
		if ok {
			var strDB []string
			switch ans.ProblemType {
			case 1: //填空题
				strDB = append(strDB, ans.Answer)
			case 2: //单选题
				strDB = append(strDB, ans.AnswerList...)
			case 3: //多选题
				strDB = append(strDB, ans.AnswerList...)
			default:
				return
			}
			if CompStringList(strReq, strDB) { //回答正确
				ansIdList = append(ansIdList, ans.Id)
				blockTMap[ans.BelongBlock] += 1
			} else {
				if ans.IsRequire == 2 { //非必填项
					blockTMap[ans.BelongBlock] += 1
				}
			}
		}
	}

	// 计算总得分
	totalScore := 0
	if len(ansIdList) != 0 {
		if err = global.GDb.Table("assessment_problem_singles").Select("SUM(set_score) as total_score").
			Where("id in (?)", ansIdList).Scan(&totalScore).Error; err != nil {
			return
		}
	}

	var criterion int
	if err = global.GDb.Table("assessment_problems").Select("assessment_criterion").
		Where("id = ?", exam.AssessmentProblemId).Scan(&criterion).Error; err != nil {
		return
	}

	//point := ((rateTureBlock) / float64(len(blockMap))) * 100

	tmpCount := (*int)(unsafe.Pointer(&totalScore))
	point := *tmpCount
	isPass = point >= criterion
	pointStr = fmt.Sprintf("%v", point)

	status := 2
	if isPass {
		status = 1
	}
	//计入考核记录
	err = global.GDb.Table("assessment_records").Where("user_id = ?", userInfo.ID).Updates(map[string]interface{}{
		"user_status":  status,
		"last_exam_at": time.Now(),
	}).Error

	//将试卷存入数据库
	respExam := make(map[string]interface{})
	err = CreatedExam(&respExam, singlesMap, exam.AssessmentProblemId)
	if err != nil {
		return
	}
	var assessmentUserAnswer model.AssessmentUserAnswer

	assessmentUserAnswer.AnswerSheet = respExam
	if err = global.GDb.Table("assessment_problems").Select("project_code").
		Where("id = ?", exam.AssessmentProblemId).Scan(&assessmentUserAnswer.ProjectCode).Error; err != nil {
		return
	}

	if err = global.GDb.Table("assessment_problems").Select("problem_Name").
		Where("id = ?", exam.AssessmentProblemId).Scan(&assessmentUserAnswer.ProblemName).Error; err != nil {
		return
	}
	assessmentUserAnswer.ProblemId = exam.AssessmentProblemId
	assessmentUserAnswer.Score = pointStr
	assessmentUserAnswer.Standard = criterion
	assessmentUserAnswer.IsPass = status
	//assessmentUserAnswer.UserName = userInfo.NickName
	assessmentUserAnswer.UserName = userInfo.Name
	assessmentUserAnswer.UserCode = userInfo.Code
	if err = global.GDb.Create(&assessmentUserAnswer).Error; err != nil {
		return
	}
	//TODO 把考核通过的人 添加到培训管理中 暂时还没有 	 这四个参数  后面补充
	//					"practice_yield_requirement": "",
	//					"accuracy_requirement":       "",
	//					"actual_practice_yield":      "",
	//					"actual_practice_accuracy":   "",
	//
	//1.查询所有项目考核
	var assessmentCriterionItem []respModel.RespAssessCriterion //考核标准
	if err = global.GDb.Model(&model.AssessmentCriterion{}).Find(&assessmentCriterionItem).Error; err != nil {
		return
	}
	//查询用户信息
	var user model2.SysUser
	global.GDb.Model(&model2.SysUser{}).Where("id  = ?", userInfo.ID).Find(&user)
	fmt.Println("===============assessCriterion.SetPoint========APP=  user", user)

	//查询产量
	for _, assessCriterion := range assessmentCriterionItem {
		if assessCriterion.ProjectCode == assessmentUserAnswer.ProjectCode {
			if point >= assessCriterion.SetPoint {
				//查询产量要求和准确率 model.PracticeSum{}
				var practiceAsk practiceModel.PracticeAsk
				var practiceSums []practiceModel.PracticeSum
				err = global.GDb.Model(&practiceModel.PracticeAsk{}).Where("pro_code = ?", assessmentUserAnswer.ProjectCode).Find(&practiceAsk).Error
				//err = global.GDb.Model(&practiceModel.PracticeSum{}).Where("code = ? AND  pro_code = ?", assessmentUserAnswer.UserCode, assessmentUserAnswer.ProjectCode).Find(&practiceSum).Error
				err := global.ProDbMap[assessmentUserAnswer.ProjectCode].Model(&practiceModel.PracticeSum{}).Where("code = ? AND  pro_code = ?", assessmentUserAnswer.UserCode, assessmentUserAnswer.ProjectCode).Order("created_at desc").Find(&practiceSums).Error
				if err != nil {
					return false, "", err
				}
				//考核达标的人添加到培训管理
				var training trainingModule.TrainingManagement
				training.UserCode = assessmentUserAnswer.UserCode
				training.UserName = assessmentUserAnswer.UserName
				if user.Sex == true {
					training.Gender = 2
				} else {
					training.Gender = 1
				}
				training.Phone = user.Phone
				training.ProjectCode = assessmentUserAnswer.ProjectCode
				training.EntryStartAt = user.EntryDate
				training.DutyAt = user.MountGuardDate
				training.PracticeYieldRequirement = practiceAsk.Character
				training.AccuracyRequirement = practiceAsk.AccuracyRate
				if len(practiceSums) > 0 {
					training.ActualPracticeYield = practiceSums[0].SummaryFieldEffectiveCharacter
					training.ActualPracticeAccuracy = practiceSums[0].SummaryAccuracyRate
				} else {
					training.ActualPracticeYield = 0
					training.ActualPracticeAccuracy = 0
				}

				training.Examine = assessCriterion.SetPoint
				training.ExamineScore = point
				training.AuditStatus = 1
				err = global.GDb.Model(&trainingModule.TrainingManagement{}).Create(&training).Error
				if err != nil {
					return false, "", err
				}
			}
		}
	}

	return
}
