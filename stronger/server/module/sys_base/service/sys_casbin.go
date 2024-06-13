package service

import (
	"errors"
	"server/module/sys_base/model"
	sys_base2 "server/module/sys_base/model/request"

	//"github.com/casbin/casbin"
	"github.com/casbin/casbin/util"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"server/global"
	"strings"
)

// @title    UpdateCasbin
// @description   update casbin authority, 更新casbin权限
// @auth                     （2020/04/05  20:22）
// @param     authorityId      string
// @param     casbinInfos      []CasbinInfo
// @return                     error

func UpdateCasbin(authorityId string, casbinInfos []sys_base2.CasbinInfo) error {
	ClearCasbin(0, authorityId)
	rules := [][]string{}
	for _, v := range casbinInfos {
		cm := model.CasbinModel{
			ID:     0,
			PType:  "p",
			RoleId: authorityId,
			Path:   v.Path,
			Method: v.Method,
		}
		rules = append(rules, []string{cm.RoleId, cm.Path, cm.Method})
	}
	e := Casbin()
	success, _ := e.AddPolicies(rules)
	if success == false {
		return errors.New("存在相同api,添加失败,请联系管理员")
	}
	return nil
}

// @title    AddCasbin
// @description   add casbin authority, 添加权限
// @auth                     （2020/04/05  20:22）
// @param     cm              model.CasbinModel
// @return                    bool

func AddCasbin(cm model.CasbinModel) (bool, error) {
	e := Casbin()
	return e.AddPolicy(cm.RoleId, cm.Path, cm.Method)
}

// @title    UpdateCasbinApi
// @description   update casbin apis, API更新随动
// @auth                     （2020/04/05  20:22）
// @param     oldPath          string
// @param     newPath          string
// @param     oldMethod        string
// @param     newMethod        string
// @return                     error

func UpdateCasbinApi(oldPath string, newPath string, oldMethod string, newMethod string) error {
	err := global.GDb.Table("casbin_rule").Model(&model.CasbinModel{}).Where("v1 = ? AND v2 = ?", oldPath, oldMethod).Updates(map[string]interface{}{
		"v1": newPath,
		"v2": newMethod,
	}).Error
	return err
}

// @title    GetPolicyPathByAuthorityId
// @description   get policy path by authorityId, 获取权限列表
// @auth                     （2020/04/05  20:22）
// @param     authorityId     string
// @return                    []string

func GetPolicyPathByAuthorityId(authorityId string) (pathMaps []sys_base2.CasbinInfo) {
	e := Casbin()
	list := e.GetFilteredPolicy(0, authorityId)
	for _, v := range list {
		pathMaps = append(pathMaps, sys_base2.CasbinInfo{
			Path:   v[1],
			Method: v[2],
		})
	}
	return pathMaps
}

// @title    ClearCasbin
// @description   清除匹配的权限
// @auth                     （2020/04/05  20:22）
// @param     v               int
// @param     p               string
// @return                    bool

func ClearCasbin(v int, p ...string) bool {
	e := Casbin()
	success, _ := e.RemoveFilteredPolicy(v, p...)
	return success

}

// @title    Casbin
// @description   store to DB, 持久化到数据库  引入自定义规则
// @auth                     （2020/04/05  20:22）

func Casbin() *casbin.Enforcer {
	p := global.GConfig.Postgresql
	dsn := "user=" + p.Username + " password=" + p.Password + " host=" + p.Host + " dbname=" + p.Dbname + " port=" + p.Port + " " + p.Config
	global.GLog.Info(dsn)
	//a, _ := gormadapter.NewAdapter(global.G_CONFIG.System.DbType, "postgresql://postgres:root@localhost:5432/stronger?sslmode=disable", true)
	a, _ := gormadapter.NewAdapter(global.GConfig.System.DbType, dsn, true)
	e, _ := casbin.NewEnforcer(global.GConfig.Casbin.ModelPath, a)
	e.AddFunction("ParamsMatch", ParamsMatchFunc)
	_ = e.LoadPolicy()
	return e
}

// @title    ParamsMatch
// @description   customized rule, 自定义规则函数
// @auth                     （2020/04/05  20:22）
// @param     fullNameKey1    string
// @param     key2            string
// @return                    bool

func ParamsMatch(fullNameKey1 string, key2 string) bool {
	key1 := strings.Split(fullNameKey1, "?")[0]
	// 剥离路径后再使用casbin的keyMatch2
	return util.KeyMatch2(key1, key2)
}

// @title    ParamsMatchFunc
// @description   customized function, 自定义规则函数
// @auth                     （2020/04/05  20:22）
// @param     args            ...interface{}
// @return                    interface{}
// @return                    error

func ParamsMatchFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)

	return ParamsMatch(name1, name2), nil
}
