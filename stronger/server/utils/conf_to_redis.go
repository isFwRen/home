package utils

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"server/global"
	"server/module/pro_conf/model"
)

//将项目配置加载到Redis key-value
func AddProjectConfToRedis() {
	global.GLog.Info("开始将项目配置加载到Redis")
	db := global.GDb.Model(&model.SysProject{})
	var sysProjectList []model.SysProject
	err := db.Order("id desc").Find(&sysProjectList).Error
	if err == nil {
		for _, v := range sysProjectList {
			///添加一条
			value, err1 := json.Marshal(v)
			if err1 != nil {
				global.GLog.Error("json marshal err")
				return
			}
			global.GRedis.Do("set", v.Code+"_pro_conf", value)
		}
	}
	getProjectConfFromRedis("B1210_pro_conf")
}

//废用
func getProjectConfFromRedis(name string) interface{} {
	global.GLog.Info("开始从Redis读取数据")

	val2, err := global.GRedis.Get(name).Result()
	if err == redis.Nil {
		global.GLog.Error("key2 does not exists")
	} else if err != nil {
		panic(err)
	}

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(val2), &dat); err == nil {
		global.GLog.Info("==============json str 转map=======================")
	}
	return dat
}

//////////////////////////////////////////////////////////////////////////////////////////

//将项目配置加载到Redis  key-list
func AddProjectConfToRedisList(name string) (returnRes []string) {
	var result []string
	global.GLog.Info("开始将项目配置加载到Redis")
	db := global.GDb.Model(&model.SysProject{})
	var sysProjectList []model.SysProject
	err := db.Order("id desc").Find(&sysProjectList).Error
	if err == nil {
		for _, v := range sysProjectList {
			///添加一条
			value, err1 := json.Marshal(v)
			if err1 != nil {
				global.GLog.Error("json marshal err")
				return
			}
			global.GRedis.RPush(name, value).Err()
		}

		rLen, _ := global.GRedis.LLen(name).Result()
		result, _ = global.GRedis.LRange(name, 0, rLen-1).Result()
	}
	return result
}

//更新项目配置加载到Redis
//name : key
func UpdateProjectConfToRedis(name string) (result []string) {
	//del kv
	var rLen, _ = global.GRedis.LLen(name).Result()
	for i := 0; i < int(rLen); i++ {
		global.GRedis.LPop(name)
	}

	//add kv
	return AddProjectConfToRedisList(name)
}
