/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/3 3:47 下午
 */

package service

import (
	"fmt"
	"server/global"
	"server/module/pro_manager/model"
	"time"
)

//GetTaskBillByPage 根据一堆过滤条件获取内存库案件列表
func GetTaskBillByPage(billListSearch model.BillListSearch) (err error, total int64, list interface{}) {
	//limit := billListSearch.PageSize
	//offset := billListSearch.PageSize * (billListSearch.PageIndex - 1)
	var projectBills []model.ProjectBill
	db := global.GTaskDb
	if db == nil {
		return global.ProDbErr, 0, projectBills
	}

	// 查询数据
	rows, err := db.Query("SELECT * FROM user_info")
	if err != nil {
		return err, 0, nil
	}

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created time.Time
		err = rows.Scan(&uid, &username, &department, &created)
		if err != nil {
			return err, 0, nil
		}
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}
	return err, total, projectBills
}
