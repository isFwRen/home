package task

import (
	"fmt"
	"server/module/load/service"
	sumModel "server/module/report_management/model"
	"server/module/task/response"
	"time"
)

var TaskListOpNum response.TaskOpNum
var UserDayOutputs []sumModel.OutputStatistics

func TaskList() {
	// fmt.Println("----------CacheFieldConf--------------------:", CacheFieldConf)
	for true {
		// TaskListOpNum.Op0
		_, TaskListOpNum.Op0 = service.CountBlockOpNum(TaskProCode, "op0")
		_, TaskListOpNum.Op1 = service.CountBlockOpNum(TaskProCode, "op1")
		_, TaskListOpNum.Op2 = service.CountBlockOpNum(TaskProCode, "op2")
		_, TaskListOpNum.Opq = service.CountBlockOpNum(TaskProCode, "opq")
		submitTime, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
		_, UserDayOutputs = service.DayOutput(TaskProCode, submitTime)
		fmt.Println("----------分块数量--------------------:", TaskListOpNum)
		<-time.After(TaskListSecond * time.Second)
	}
}
