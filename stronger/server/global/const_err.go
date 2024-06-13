/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/11 3:00 下午
 */

package global

import "errors"

var ProDbErr = errors.New("没有连接项目数据库")
var RoleErr = errors.New("权限不足")
var NoLogin = errors.New("未登录或非法访问")
var ParamErr = errors.New("参数有误")
var NotPro = errors.New("没有找到该项目")

var LogErr = errors.New("插入日志失败")
