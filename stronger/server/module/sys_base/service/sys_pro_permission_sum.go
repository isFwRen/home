package service

import (
	"fmt"
	"server/global"
	"server/module/sys_base/model"
)

func GetProPermissionSum() (err error, total int, Sum []model.ProPermissionSum) {
	var ProPermissionSum []model.ProPermissionSum
	err = global.GDb.Model(&model.SysProPermission{}).
		//Raw("select pro_code as procode, pro_name as proname, count(case when has_op0 then 1 end) AS op0,\n       count(case when has_op1 then 1 end) AS op1,\n       count(case when has_op2 then 1 end) AS op2,\n       count(case when has_opq then 1 end) AS opq,\n       count(case when sys_users.code !~ '^[P]' AND has_out_net then 1 end) AS outnet,\n       count(case when sys_users.code ~ '^[P]' AND has_in_net then 1 end) AS innet\nfrom sys_pro_permissions, sys_users where sys_pro_permissions.user_code = sys_users.code and sys_users.status = true\ngroup by sys_pro_permissions.pro_name, pro_code").
		Raw("select pro_code as procode, pro_name as proname, \n\t\t\t count(case when has_op0 then 1 end) AS op0,\n       count(case when has_op1 then 1 end) AS op1,\n       count(case when has_op2 then 1 end) AS op2,\n       count(case when has_opq then 1 end) AS opq\nfrom sys_pro_permissions\ngroup by sys_pro_permissions.pro_name, pro_code").
		Scan(&ProPermissionSum).Error
	if err != nil {
		return nil, 0, nil
	}
	fmt.Println(ProPermissionSum)
	return err, len(ProPermissionSum), ProPermissionSum
}
