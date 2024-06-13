package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/go-kivik/kivik/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"os"
	"path"
	"server/global"
	proConf "server/module/pro_conf/model"
	"server/module/pro_conf/model/request"
	"server/module/pro_conf/model/response"
	"server/utils"
	"sync"

	//"server/module/sys_base/model"
	"strings"
)

func GetTableTop(c *gin.Context, dbname string) (err error, list interface{}) {
	fmt.Println("dbname", dbname)
	Couchdb, err := CheckAndConnectionCouchdb(c, dbname)
	defer Couchdb.Close(c)
	if err != nil {
		return err, nil
	}
	find, err := Couchdb.Find(c, "{\"selector\": {\"_id\": {\"$regex\": \"^TableTop.*\"}}}")
	if err != nil {
		return err, nil
	}
	var tabletops response.TableTop
	//var tabletopDetail response.TableTopResp
	//tabletopResp := make([]response.TableTopResp, 0)
	top := make([]interface{}, 0)
	for find.Next() {
		if err = find.ScanDoc(&tabletops); err != nil {
			return err, nil
		}
		//tabletopDetail.Id = tabletops.Id
		//tabletopDetail.Rev = tabletops.Rev
		//tabletopDetail.Top = tabletops.Arr
		//tabletopResp = append(tabletopResp, tabletopDetail)
		//fmt.Println("sadas", tabletops.Arr)
		top = tabletops.Tabletop
	}
	return err, top
}

func GetConstTablesByProjectName(c *gin.Context, constTable request.ConstTableWithProject) (err error, constTables []response.ConstTableNameByProjectResp, count int) {
	//将请求结构体封装成查询结构体
	constQueryStruct := request.ConstQueryStruct{Selector: constTable}
	lpCouchdb, err := CheckAndConnectionCouchdb(c, "lp2")
	defer lpCouchdb.Close(c)
	if err != nil {
		return err, nil, 0
	}
	//"{\"selector\": {\"_id\": {\"$regex\": \"^TableTop.*\"}}}"
	fmt.Println("GetConstTablesByProjectName", constTable)
	constQueryStruct.Limit = 1000
	find, err := lpCouchdb.Find(c, constQueryStruct)
	if err != nil {
		return err, nil, 0
	}
	var constNameByProject response.ConstTableNameByProject
	var constNameByProjectDetail response.ConstTableNameByProjectResp
	ConstTableNameByProjectResp := make([]response.ConstTableNameByProjectResp, 0)
	for {
		if find.Next() {
			if err = find.ScanDoc(&constNameByProject); err != nil {
				return err, nil, count
			}
			constNameByProjectDetail.Id = constNameByProject.Id
			constNameByProjectDetail.Rev = constNameByProject.Rev
			constNameByProjectDetail.Conf = constNameByProject.Conf
			constNameByProjectDetail.ProCode = constNameByProject.ProCode
			constNameByProjectDetail.FileName = constNameByProject.FileName
			constNameByProjectDetail.DbName = constNameByProject.DbName
			constNameByProjectDetail.ChineseName = strings.Split(constNameByProject.ChineseName, "_")[2]
			ConstTableNameByProjectResp = append(ConstTableNameByProjectResp, constNameByProjectDetail)
		} else {
			break
		}
	}
	return err, ConstTableNameByProjectResp, len(ConstTableNameByProjectResp)
}

func PutConstTableByArr(c *gin.Context, constTable request.UpdateConstTableStructWithItems) (err error, table request.UpdateConstTableStructWithItems) {
	//直接保存传进来的数据需要Rev，否则会报update conflict
	var t request.UpdateConstTableStructWithItems
	Couchdb, err := CheckAndConnectionCouchdb(c, constTable.FileName)
	defer Couchdb.Close(c)
	if err != nil {
		return err, table
	}
	var newConst request.PutConst
	newConst.Id = constTable.Id
	newConst.Rev = constTable.Rev
	newConst.Docs = constTable.Docs
	rev, err := Couchdb.Put(c, constTable.Id, newConst, kivik.Options{})
	if err != nil {
		return err, t
	}
	constTable.Rev = rev
	return err, constTable
}

func DelConstTableLineById(c *gin.Context, R request.DelConstTableLineArr) (err error, rev map[string]string) {
	Couchdb, err := CheckAndConnectionCouchdb(c, R.DbName)
	defer Couchdb.Close(c)
	if err != nil {
		return err, nil
	}
	newRevItem := make(map[string]string, 0)
	for _, r := range R.Table {
		newRev, err := Couchdb.Delete(c, r.Id, r.Rev)
		if err != nil {
			return err, nil
		}
		newRevItem[r.Id] = newRev
	}
	return err, newRevItem
}

func PutConstTableWithExcel(projectCode string, relationship map[string]string) (failArray string, errMsg string) {
	//basicPath := global.GConfig.LocalUpload.FilePath + projectCode + "/常量表/"
	// 线上 uploads/file
	basicPath := global.GConfig.LocalUpload.FilePath + "const/" + projectCode + "/"
	// 本地
	//basicPath := "D:/Projects/tempfiles/B0108/常量表"
	wg := sync.WaitGroup{}
	for fn, dn := range relationship {

		wg.Add(1)
		go func(fn, dn string) {
			fmt.Println("fn:", fn, "dn:", dn)
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(fmt.Sprintf("上传常量表失败, Error: %s", r))
					global.GLog.Error(fmt.Sprintf("上传常量表失败, Error: %s", r))
					wg.Done()
				}
			}()
			var ConstTables []request.InsertConstTableStructWithItems

			isExists, err := global.GCouchdbClient.DBExists(context.TODO(), dn)
			if err != nil {
				panic(errors.New(fmt.Sprintf("常量表:%s, err:%s", fn, err.Error())))
			}
			if isExists {
				//panic(errors.New(fmt.Sprintf("常量表:%s, 已存在", fn)))
				// 删除当前常量表所对应的数据库
				err = DeleteDatabase(dn)
				if err != nil {
					panic(errors.New(fmt.Sprintf("覆盖失败, 常量表:%s, err:%s", fn, err.Error())))
				}
			}
			var ConstTableFromExcel request.InsertConstTableStructWithItems
			dst := path.Join(basicPath, fn)
			fmt.Println("文件路径", dst)
			xlsx, err := excelize.OpenFile(dst)
			if err != nil {
				panic(errors.New(fmt.Sprintf("常量表:%s, 读取失败excel, 无法存进couchdb, err:%s", fn, err.Error())))
			}
			ConstTableFromExcel.FilePath = dst
			ConstTableFromExcel.FileName = fn
			ConstTableFromExcel.DbName = dn
			ConstTableFromExcel.ChineseName = strings.Replace(fn, ".xlsx", "", -1)
			// 获取excel中具体的列的值
			rows := xlsx.GetRows("Sheet" + "1")
			Items := make([]request.Items, len(rows))
			TableTopItems := make([]request.TableTopItems, 1)
			// 循环刚刚获取到的表中的值
			for key, row := range rows {
				if key == 0 {
					//标题行hang为table_top
					TableTopItems[0].Tabletop = row
				} else {
					//其余的皆为子项
					uuids, _ := uuid.NewV4()
					Items[key].Id = "const::" + uuids.String()
					Items[key].Arr = row
					ConstTableFromExcel.Docs = append(ConstTableFromExcel.Docs, Items[key])
				}
			}
			ConstTableFromExcel.TableTop = TableTopItems
			ConstTableFromExcel.ProCode = projectCode
			ConstTables = append(ConstTables, ConstTableFromExcel)

			fmt.Println(fmt.Sprintf("%s,ConstTables完成,长度:%d", fn, len(ConstTables)))

			//连接lp2总数据库(里面存有所以常量表的基本信息)
			LpCouchdb, err := CheckAndConnectionCouchdb(context.TODO(), "lp2")
			defer LpCouchdb.Close(context.TODO())
			if err != nil {
				panic(errors.New(fmt.Sprintf("常量表:%s, err:%s", fn, err.Error())))
			}
			if len(ConstTables) == 0 {
				panic(errors.New(fmt.Sprintf("常量表:%s, 上传失败", fn)))
			}

			//先把基本信息存进lp2数据库中, 保存excel表格里的数据
			for _, constable := range ConstTables {
				//fmt.Println(constable.Docs[0], constable.Docs[len(constable.Docs)-1])
				//1. 把基本信息存进lp2数据库中
				fmt.Println(fmt.Sprintf("%s,把基本信息存进lp2数据库中", fn))
				uuids, e := uuid.NewV4()
				if e != nil {
					panic(errors.New(fmt.Sprintf("常量表:%s, err:%s", fn, "There is a error when the system create a uuid for constTableBaseInformation ")))
				}
				constTableBaseInformation := request.ConstTableBaseInformation{
					ProCode:     constable.ProCode,
					Conf:        "const",
					FileName:    constable.FileName,
					DbName:      constable.DbName,
					FilePath:    constable.FilePath,
					ChineseName: constable.ChineseName,
				}
				rev, err := LpCouchdb.Put(context.TODO(), "constTableBaseInformation::"+uuids.String(), constTableBaseInformation)
				if err != nil {
					panic(errors.New(fmt.Sprintf("常量表:%s, err:%s", fn, "There is a error when the system create a docs about constTableBaseInformation in Couchdb ")))
				}
				constTableBaseInformationId := "constTableBaseInformation::" + uuids.String()
				constTableBaseInformationRev := rev
				//2. 保存excel表格里的数据
				Couchdb, err := CheckAndCreateAndConnectCouchdb(context.TODO(), constable.DbName)
				if err != nil {
					err = DeleteDateInConstSum(constTableBaseInformationId, constTableBaseInformationRev)
					if err != nil {
						panic(errors.New(fmt.Sprintf("常量表:%s, 保存excel表格里的数据, err:%s", fn, err.Error())))
					}
				}
				//2.1 存中文表名
				fmt.Println(fmt.Sprintf("%s,存中文表名", fn))
				uuids, e = uuid.NewV4()
				if e != nil {
					fmt.Println("There is a error when the system create a uuid for ChineseTableName ")
					errMsg += " There is a error when the system create a uuid for ChineseTableName " + "\n"
					// 删除当前常量表所对应的数据库
					err = DeleteDatabase(constable.DbName)
					if err != nil {
						errMsg += err.Error() + "\n"
					}
					err = DeleteDateInConstSum(constTableBaseInformationId, constTableBaseInformationRev)
					if err != nil {
						errMsg += err.Error() + "\n"
					}
					panic(errors.New(fmt.Sprintf("常量表:%s, 存中文表名, err:%s", fn, errMsg)))
				}
				ChineseNameDoc := request.Docs{
					Content: constable.ChineseName,
				}
				_, err = Couchdb.Put(context.TODO(), "ChineseTableName::"+uuids.String(), ChineseNameDoc)
				if err != nil {
					errMsg += constable.FileName + ": There is a error when the system create a docs about ChineseTableName in Couchdb " + "\n"
					// 删除当前常量表所对应的数据库
					err = DeleteDatabase(constable.DbName)
					if err != nil {
						errMsg += err.Error() + "\n"
					}
					err = DeleteDateInConstSum(constTableBaseInformationId, constTableBaseInformationRev)
					if err != nil {
						errMsg += err.Error() + "\n"
					}
					panic(errors.New(fmt.Sprintf("常量表:%s, 存中文表名, err:%s", fn, errMsg)))
				}
				// 存标题行
				fmt.Println(fmt.Sprintf("%s,存标题行", fn))
				uuids, e = uuid.NewV4()
				if e != nil {
					fmt.Println("There is a error when the system create a uuid for TableTop ")
					errMsg += constable.FileName + ": There is a error when the system create a uuid for TableTop " + "\n"
					// 删除当前常量表所对应的数据库
					err = DeleteDatabase(constable.DbName)
					if err != nil {
						errMsg += err.Error() + "\n"
					}
					err = DeleteDateInConstSum(constTableBaseInformationId, constTableBaseInformationRev)
					if err != nil {
						errMsg += err.Error() + "\n"
					}
					panic(errors.New(fmt.Sprintf("常量表:%s, 存标题行, err:%s", fn, errMsg)))
				}
				_, err = Couchdb.Put(context.TODO(), "TableTop::"+uuids.String(), constable.TableTop[0])
				if err != nil {
					errMsg += constable.FileName + ": There is a error when the system create a docs about TableTop in Couchdb " + "\n"
					// 删除当前常量表所对应的数据库
					err = DeleteDatabase(constable.DbName)
					if err != nil {
						errMsg += err.Error() + "\n"
					}
					err = DeleteDateInConstSum(constTableBaseInformationId, constTableBaseInformationRev)
					if err != nil {
						errMsg += err.Error() + "\n"
					}
					panic(errors.New(fmt.Sprintf("常量表:%s, 存标题行, err:%s", fn, errMsg)))
				}

				//存常量数据, 批量存
				fmt.Println(fmt.Sprintf("%s,存常量数据", fn))
				if len(constable.Docs) < 5000 {
					_, err = Couchdb.BulkDocs(context.TODO(), constable.Docs)
					if err != nil {
						fmt.Println("err", err)
						// 删除当前常量表所对应的数据库
						err = DeleteDatabase(constable.DbName)
						if err != nil {
							errMsg += err.Error() + "\n"
						}
						err = DeleteDateInConstSum(constTableBaseInformationId, constTableBaseInformationRev)
						if err != nil {
							errMsg += err.Error() + "\n"
						}
						panic(errors.New(fmt.Sprintf("常量表:%s, 存常量数据, err:%s", fn, errMsg)))
					}
				} else {
					for i, j := 0, 5000; i < len(constable.Docs); {
						//fmt.Println(i, j)
						_, err = Couchdb.BulkDocs(context.TODO(), constable.Docs[i:j])
						if err != nil {
							fmt.Println("err", err)
							// 删除当前常量表所对应的数据库
							err = DeleteDatabase(constable.DbName)
							if err != nil {
								errMsg += err.Error() + "\n"
							}
							err = DeleteDateInConstSum(constTableBaseInformationId, constTableBaseInformationRev)
							if err != nil {
								errMsg += err.Error() + "\n"
							}
							panic(errors.New(fmt.Sprintf("常量表:%s, 存常量数据, err:%s", fn, errMsg)))
						}
						i = j
						if j+5000 >= len(constable.Docs) {
							j = len(constable.Docs)
						} else {
							j += 5000
						}
					}
				}

				//存常量数据, 一条一条存
				//for i := 0; i < len(constable.Docs); i++ {
				//	uuids, e = uuid.NewV4()
				//	if e != nil {
				//		fmt.Println("There is a error when the system create a uuid for ConstDate ")
				//		errMsg += constable.FileName + ": There is a error when the system create a uuid for ConstDate " + "\n"
				//		// 删除当前常量表所对应的数据库
				//		err = DeleteDatabase(constable.FileName)
				//		if err != nil {
				//			errMsg += err.Error() + "\n"
				//		}
				//		err = DeleteDateInConstSum(constTableBaseInformationId, constTableBaseInformationRev)
				//		if err != nil {
				//			errMsg += err.Error() + "\n"
				//		}
				//		break
				//	}
				//	fmt.Println("zzz", constable.Docs[i])
				//	_, err = Couchdb.Put(context.TODO(), "const::"+uuids.String(), constable.Docs[i])
				//	if err != nil {
				//		fmt.Println("There is a error when the system create a docs about ConstDate in Couchdb ")
				//		errMsg += constable.FileName + ": There is a error when the system create a docs about ConstDate in Couchdb " + "\n"
				//		// 删除当前常量表所对应的数据库
				//		err = DeleteDatabase(constable.FileName)
				//		if err != nil {
				//			errMsg += err.Error() + "\n"
				//			break
				//		}
				//		err = DeleteDateInConstSum(constTableBaseInformationId, constTableBaseInformationRev)
				//		if err != nil {
				//			errMsg += err.Error() + "\n"
				//		}
				//		break
				//	}
				//}
			}
			fmt.Println(fmt.Sprintf("%s,上传常量数据成功", fn))
			wg.Done()

		}(fn, dn)

	}

	wg.Wait()
	return failArray, errMsg
}

func CheckAndConnectionCouchdb(c context.Context, dbname string) (Couchdb *kivik.DB, err error) {
	isExists, err := global.GCouchdbClient.DBExists(c, dbname)
	if err != nil {
		return nil, err
	}
	if !isExists {
		return nil, errors.New(dbname + "不存在")
	}
	Couchdb = global.GCouchdbClient.DB(context.Background(), dbname)
	if Couchdb == nil {
		global.GLog.Error("CouchDB "+dbname+" 初始化失败", zap.Any("err", err))
	}
	fmt.Println("CouchDB " + dbname + " 初始化成功")
	return Couchdb, nil
}

func CheckAndCreateAndConnectCouchdb(c context.Context, dbname string) (Couchdb *kivik.DB, err error) {
	fmt.Println("CheckAndCreateAndConnectCouchdb", dbname)
	isExists, err := global.GCouchdbClient.DBExists(c, dbname)
	fmt.Println("a", isExists)
	if err != nil {
		return nil, err
	}
	if !isExists {
		err = global.GCouchdbClient.CreateDB(c, dbname)
		if err != nil {
			return nil, err
		}
	}
	Couchdb = global.GCouchdbClient.DB(context.Background(), dbname)
	if Couchdb == nil {
		global.GLog.Error("CouchDB "+dbname+" 初始化失败", zap.Any("err", err))
	}
	err = Couchdb.CreateIndex(context.Background(), dbname, "arr-json-index", `{"fields":["arr"]}`)
	if err != nil {
		global.GLog.Error("CouchDB "+dbname+" 创建索引失败", zap.Any("err", err))
	}
	fmt.Println("CouchDB " + dbname + " 初始化成功")
	return Couchdb, nil
}

func DeleteDatabase(dbname string) error {
	fmt.Println("DeleteDatabase dbname", dbname)
	isExists, err := global.GCouchdbClient.DBExists(context.TODO(), dbname)
	if err != nil {
		return err
	}
	if !isExists {
		return errors.New(dbname + "不存在无法删除")
	}
	err = global.GCouchdbClient.DestroyDB(context.TODO(), dbname)
	if err != nil {
		return err
	}
	//删除lp2常量信息汇总数据库下该要删除数据库的信息
	Couchdb, err := CheckAndConnectionCouchdb(context.TODO(), "lp2")
	defer Couchdb.Close(context.TODO())
	if err != nil {
		return err
	}
	query := "{\n   \"selector\": {\n      \"dbName\": \"" + dbname + "\"\n   }\n}"
	find, err := Couchdb.Find(context.TODO(), query)
	if err != nil {
		return err
	}
	var constNameByProject response.ConstTableNameByProject
	var constNameByProjectDetail response.ConstTableNameByProjectResp
	ConstTableNameByProjectResp := make([]response.ConstTableNameByProjectResp, 0)
	for {
		if find.Next() {
			if err = find.ScanDoc(&constNameByProject); err != nil {
				return err
			}
			constNameByProjectDetail.Id = constNameByProject.Id
			constNameByProjectDetail.Rev = constNameByProject.Rev
			constNameByProjectDetail.FilePath = constNameByProject.FilePath
			ConstTableNameByProjectResp = append(ConstTableNameByProjectResp, constNameByProjectDetail)
		} else {
			break
		}
	}
	for _, v := range ConstTableNameByProjectResp {
		_, err = Couchdb.Delete(context.TODO(), v.Id, v.Rev)
		if err != nil {
			return err
		}
		//删除xlsx表
		_ = os.Remove(v.FilePath)
	}

	return nil
}

// DeleteDateInConstSum 删除常量表基本信息
func DeleteDateInConstSum(id, rev string) error {
	Couchdb, err := CheckAndConnectionCouchdb(context.TODO(), "lp2")
	defer Couchdb.Close(context.TODO())
	if err != nil {
		return err
	}
	_, err = Couchdb.Delete(context.TODO(), id, rev)
	if err != nil {
		return err
	}
	return err
}

func ExportConstExcel(proCode, dbname string) (err error, path string) {
	isExists, err := global.GCouchdbClient.DBExists(context.TODO(), "lp2")
	if err != nil {
		return err, ""
	}
	if !isExists {
		return errors.New("lp2 is not exist"), ""
	}
	lpCouchDB, err := CheckAndConnectionCouchdb(context.TODO(), "lp2")
	if err != nil {
		return err, ""
	}
	query := "{\n   \"selector\": {\n      \"dbName\": \"" + dbname + "\"\n   }\n}"
	//query := "{\n   \"selector\": {\n      \"project\": \"" + proCode + "\"\n   }\n}"
	find, err := lpCouchDB.Find(context.TODO(), query)
	var constNameByProject response.ConstTableNameByProject
	for {
		if find.Next() {
			if err = find.ScanDoc(&constNameByProject); err != nil {
				return err, ""
			}
		} else {
			break
		}
	}
	exportConstDb, err := CheckAndConnectionCouchdb(context.TODO(), dbname)
	if err != nil {
		return err, ""
	}
	//excel
	excelLine := make([][]interface{}, 1)
	//获取表头
	fmt.Println("exportConstDb-获取表头")
	finds, err := exportConstDb.Find(context.TODO(), "{\"selector\": {\"_id\": {\"$regex\": \"^TableTop.*\"}}}")
	if err != nil {
		return err, ""
	}
	var tabletops response.TableTop
	top := make([]interface{}, 0)
	for finds.Next() {
		if err = finds.ScanDoc(&tabletops); err != nil {
			return err, ""
		}
		top = tabletops.Tabletop
	}
	excelLine[0] = top
	//获取常量内容
	row, err := exportConstDb.AllDocs(context.TODO(), kivik.Options{"include_docs": true})
	if err != nil {
		return err, ""
	}

	for {
		var constDate proConf.ConstTable
		if row.Next() {
			if err = row.ScanDoc(&constDate); err != nil {
				fmt.Println(err)
			}
			if strings.Index(constDate.Id, "TableTop") == -1 && strings.Index(constDate.Id, "ChineseTableName") == -1 {
				//fmt.Println("A", constDate.Arr)
				excelLine = append(excelLine, constDate.Arr)
			}
		} else {
			break
		}
	}
	fmt.Println("正在生成excel！")
	excelName := constNameByProject.ChineseName + "(常量表)"
	err = utils.ExportExcelTheMainEntrance(excelLine, excelName, proCode, "export-constTable")
	return err, "files/export-constTable" + "/" + proCode + "/" + excelName + ".xlsx"
}

// InsertConstLog 新增常量日志
func InsertConstLog(log proConf.ConstLog) error {
	return global.GDb.Model(&proConf.ConstLog{}).Create(&log).Error
}

// FetchConstLogPage 分页获取常量操作日志
func FetchConstLogPage(req proConf.ConstLogReq) (err error, total int64, list []proConf.ConstLog) {
	limit := req.PageSize
	offset := req.PageSize * (req.PageIndex - 1)
	db := global.GDb.Model(&proConf.ConstLog{}).
		Where("created_at BETWEEN ? AND ? and content like ?", req.StartTime, req.EndTime, "%"+req.Content+"%")
	if req.Type != "" {
		db = db.Where("type = ?", req.Type)
	}
	err = db.Count(&total).Error
	if err != nil {
		return err, total, list
	}
	err = db.Order("created_at").Limit(limit).Offset(offset).Find(&list).Error
	return err, total, list
}
