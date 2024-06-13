/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/5/7 11:55
 */

package utils

import (
	"context"
	"errors"
	"github.com/go-kivik/kivik/v3"
	"go.uber.org/zap"
	"os"
	"runtime/debug"
	"server/global"
	"server/module/pro_conf/model"
	"sync"
)

var maxGNum = 30 //限制协程并发量

// LoadConst 加载常量
func LoadConst(proCode string) (error, map[string][][]string) {
	constMap := make(map[string][][]string, 0)
	return nil, constMap
	if RegIsMatch("^(B0103|B0106|B0110)$", proCode) {
		global.GLog.Info("加载常量", zap.Any("不加载", proCode))
		return nil, constMap
	}
	if !HasItem(global.GConfig.System.ProArr, proCode) {
		return nil, constMap
	}
	//获取常量表数据库基本信息
	c := context.Background()
	baseDB, err := checkAndConnectionCouchdb(c, "lp2")
	if err != nil {
		return err, constMap
	}
	findProInfo, err := baseDB.Find(c, "{\"selector\": {\"proCode\": \""+proCode+"\"},\"limit\": 100 }")
	if err != nil {
		return err, constMap
	}
	var tableBaseInfos []model.ConstTableBaseInformation
	for findProInfo.Next() {
		var tableBaseInfo model.ConstTableBaseInformation
		err = findProInfo.ScanDoc(&tableBaseInfo)
		if err != nil {
			return err, constMap
		}
		tableBaseInfos = append(tableBaseInfos, tableBaseInfo)
	}
	var wg sync.WaitGroup
	ch := make(chan bool, maxGNum)
	for _, v := range tableBaseInfos {
		ch <- true
		wg.Add(1)
		info := v
		go func() {
			defer func() {
				wg.Done()
				<-ch
				if err := recover(); err != nil {
					global.GLog.Error("加载常量错误Stack", zap.Any(info.ProCode+"--"+info.ChineseName, string(debug.Stack())))
					global.GLog.Error("加载常量错误error", zap.Any(info.ProCode+"--"+info.ChineseName, err))
					os.Exit(0)
				}
			}()
			err = getOneConstTable(c, info, &constMap)
			if err != nil {
				global.GLog.Error("", zap.Error(err))
				global.GLog.Error("加载常量错误", zap.Any(info.ProCode, info.ChineseName))
				os.Exit(0)
			}
		}()
	}
	wg.Wait()
	return err, constMap
}

func checkAndConnectionCouchdb(c context.Context, dbname string) (Couchdb *kivik.DB, err error) {
	isExists, err := global.GCouchdbClient.DBExists(c, dbname)
	if err != nil {
		return nil, err
	}
	if !isExists {
		return nil, errors.New(dbname + "不存在")
	}
	Couchdb = global.GCouchdbClient.DB(c, dbname)
	if Couchdb == nil {
		global.GLog.Error("CouchDB "+dbname+" 初始化失败", zap.Any("err", err))
		return nil, err
	}
	global.GLog.Info("加载常量开始:::" + dbname + " 初始化成功")
	return Couchdb, nil
}

func getOneConstTable(c context.Context, tableBaseInfo model.ConstTableBaseInformation, constMap *map[string][][]string) error {
	//获取常量表
	eachConstDB, err := checkAndConnectionCouchdb(c, tableBaseInfo.DbName)
	//findConstInfo, err := eachConstDB.Find(c, "{\"selector\": {\"_id\": {\"$regex\": \"^const.*\"}}}")
	findConstInfo, err := eachConstDB.AllDocs(c, kivik.Options{"include_docs": true})
	if err != nil {
		return err
	}
	var constInfoAll [][]string
	for findConstInfo.Next() {
		var constInfo model.Items
		err = findConstInfo.ScanDoc(&constInfo)
		if err != nil {
			return err
		}
		if constInfo.Arr == nil {
			continue
		}
		constInfoAll = append(constInfoAll, constInfo.Arr)
	}
	(*constMap)[tableBaseInfo.ChineseName] = constInfoAll

	////获取常量表头信息
	findTableTop, err := eachConstDB.Find(c, "{\"selector\": {\"_id\": {\"$regex\": \"^TableTop.*\"}}}")
	if err != nil {
		return err
	}
	var tabletops model.TableTop
	var top []string
	var head [][]string
	for findTableTop.Next() {
		if err = findTableTop.ScanDoc(&tabletops); err != nil {
			return err
		}
		top = tabletops.Tabletop
		head = append(head, top)
	}
	(*constMap)[tableBaseInfo.ChineseName+"_head"] = head
	global.GLog.Info("proCode", zap.Any(tableBaseInfo.DbName, top))
	global.GLog.Info("加载常量完成:::" + tableBaseInfo.ChineseName)
	return nil
}
