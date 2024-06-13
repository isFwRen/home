/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/5/23 15:45
 */

package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"server/global"
	model3 "server/module/download/model"
	"server/module/schedule/model"
)

//DeleteDoc 删除条数据
func DeleteDoc(dbName, query string) error {
	couchdb := global.GCouchdbClient.DB(context.Background(), dbName)
	if couchdb == nil {
		global.GLog.Error("没有这个couchdb name:" + dbName)
		return errors.New("没有该数据库:" + dbName)
	}
	//`{"selector": {"arr": ["阿瓦提县塔木托格拉克乡卫生院","08","652900000162"]}}`
	rows, err := couchdb.Find(context.Background(), query)
	if err != nil {
		return err
	}
	for rows.Next() {
		var t model.TableInfo
		err = rows.ScanDoc(&t)
		if err != nil {
			return err
		}
		newRev, err := couchdb.Delete(context.Background(), t.Id, t.Rev)
		if err != nil {
			return err
		}
		fmt.Println("删除成功", newRev)
	}
	return rows.Err()
}

//InsertDoc 新增一行数据
func InsertDoc(dbName string, arr []string) error {
	data, err := json.Marshal(map[string]interface{}{"arr": arr})
	couchdb := global.GCouchdbClient.DB(context.Background(), dbName)
	if couchdb == nil {
		global.GLog.Error("没有这个couchdb name:" + dbName)
		return errors.New("没有该数据库:" + dbName)
	}
	rev, err := couchdb.Put(context.TODO(), "const::"+uuid.NewV4().String(), data)
	if err != nil {
		return err
	}
	fmt.Println(rev)
	return err
}

func UpdateStatusLog(log model3.UpdateConstLog) error {
	return global.GDb.Model(&model3.UpdateConstLog{}).
		Where("id = ? and is_updated = false", log.ID).
		Updates(map[string]interface{}{
			"is_updated": true,
		}).Error
}
