/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/5/23 09:41
 */

package schedule

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"os"
	"path"
	"server/global"
	model2 "server/module/download/model"
	"server/module/schedule/model"
	"server/module/schedule/service"
	"server/utils"
	"strconv"
	"strings"
)

var (
	TOKEN               = "Bearer 5pyq6K6+572uYmFzZTY05qC85byP55qE546v5aKD5Y+Y6YePRFVSSUFOX0FQSV9UT0tFTuWtmOWcqOS4pemHjeeahOWuieWFqOmXrumimA=="
	DingTalkAccessToken = "1beb60dd0637a25afb678c48d65b5c2bad75d7044974162732e8ab9328e71eeb"
	DingTalkSecret      = "SEC694a5ef255b5ed9d6769d93f20e143f4da6296212117815daba77395090f8226"
)

func UpdateHospitalAndCatalogue() {
	if global.GConfig.System.Env == "prod" {
		DingTalkAccessToken = "dd65e48c001b9c8843343914bb2a677d3d790f70d10194c2af00317329b1d111"
		DingTalkSecret = "SECacee2d304965b17edc782af3ee3e870ee30eb37264f689e9ed20a03488d06785"
	}
	if !utils.HasItem(global.GConfig.System.ProArr, "B0106") &&
		!utils.HasItem(global.GConfig.System.ProArr, "B0103") &&
		!utils.HasItem(global.GConfig.System.ProArr, "B0110") {
		global.GLog.Error("更新国寿常量，该项目不在该服务器", zap.Any("B0103、B0106、B0110", global.GConfig.System.ProArr))
		return
	}
	global.GLog.Info("开始更新国寿常量")
	//获取数据库昨天日志表
	list, err := service.FetchLogsYesterday()
	if err != nil {
		global.GLog.Error("获取数据库昨天日志表失败，", zap.Error(err))
		return
	}
	global.GLog.Info("获取数据库昨天日志表", zap.Any("长度", len(list)))
	if len(list) == 0 {
		global.GLog.Error("无常量更新")
		return
	}

	proMsg := make(map[string]MsgItem)
	var proCodes []string
	//遍历列表
	for _, item := range list {
		//err, namePinYin := api.ChineseNameToPinyinName(item.ProCode, item.Name)
		//if err != nil {
		//	global.GLog.Error("转换中文失败id："+item.ID, zap.Error(err))
		//	continue
		//}
		msgItem := proMsg[item.ProCode]
		if item.Type == 1 {
			//更新医疗机构（医院）
			global.GLog.Info("开始更新医疗机构", zap.Any(item.ID, item.Name))
			//err = UpdateHospital(namePinYin, item)
			err = UpdateHospitalApp(item, &msgItem)
			if err != nil {
				global.GLog.Error("更新失败：id"+item.ID, zap.Error(err))
				continue
			}
		} else if item.Type == 2 {
			global.GLog.Info("开始更新医疗目录", zap.Any(item.ID, item.Name))
			//更新医疗目录
			err = UpdateCatalogueApp(item, &msgItem)
			if err != nil {
				global.GLog.Error("更新失败：id"+item.ID, zap.Error(err))
				//continue
			}
		} else {
			global.GLog.Error("常量log item有误，id:" + item.ID)
			continue
		}
		// 发布更新
		//curl -X PATCH -H "Content-Type: application/json" -d '{"proCode": "B0113","name":["B0113_百年理赔_百年理赔职业代码表"]}' https://example.com/sys-const/release/B0113
		releaseCmd := `curl -X PATCH --header 'Content-Type: application/json'  --header 'Authorization: ` +
			TOKEN + `' -d '{"proCode": "` +
			item.ProCode + `","name":["` +
			item.Name + `"]}' ` +
			global.GConfig.System.ConstUrl + `/sys-const/release/` +
			item.ProCode
		global.GLog.Info("发布更新 " + releaseCmd)
		err = ShellOut(releaseCmd)
		if err != nil {
			global.GLog.Error("发布更新失败：id"+item.ID, zap.Error(err))
			continue
		}

		//更新数据库
		err = service.UpdateStatusLog(item)
		if err != nil {
			global.GLog.Error("更新失败", zap.Error(err))
		}
		global.GLog.Info("更新医疗目录成功", zap.Any(item.ID, item.Name))
		proCodes = append(proCodes, item.ProCode)
		proMsg[item.ProCode] = msgItem
	}

	//钉钉通知信息
	//新增医疗机构信息:xx项目医疗机构新增y1，代码为x1;新增y2，代码为x2;（共a个医院已更新）
	//删除医疗机构信息:xx项目医疗机构删除y1，代码为x1;删除y2，代码为x2;（共a个医院已删除）
	//XXXX项目yy1、yy2、yy3共a个常量已更新;
	//XXXX项目新增XXXXX常量，请注意修改《数据库编码对应表》;
	//XXXX项目作废XXXXX常量，请注意在《数据库编码对应表》删除该目录信息；  （第一个下划线后的值为Y时）
	//XXXX项目的y1常量作废z1、z2、z3共n1条；y2常量作废z4、z5、z6共n2条；  （条件：第一个下划线后的值为N删除文件“是否作废”列中为Y）
	global.GLog.Info("----------------------------------------------------------------------------------------------------------------")
	global.GLog.Info("", zap.Any("proMsg", proMsg))
	global.GLog.Info("----------------------------------------------------------------------------------------------------------------")
	for code, item := range proMsg {
		msgArr := [6]string{}
		if len(item.AddHospital) > 0 {
			msgArr[0] = code + "项目医疗机构\n" + strings.Join(item.AddHospital, ";\n") + "\n（共" + strconv.Itoa(len(item.AddHospital)) + "个医院已更新）"
		}
		if len(item.DelHospital) > 0 {
			msgArr[1] = code + "项目医疗机构\n" + strings.Join(item.DelHospital, ";\n") + "\n（共" + strconv.Itoa(len(item.DelHospital)) + "个医院已删除）"
		}
		if len(item.UpdateCatalogue) > 0 {
			msgArr[2] = code + "项目" + strings.Join(item.UpdateCatalogue, "、") + "共" + strconv.Itoa(len(item.UpdateCatalogue)) + "个常量已更新;"
		}
		if len(item.AddCatalogue) > 0 {
			msgArr[3] = code + "项目新增" + strings.Join(item.AddCatalogue, "、") + "常量，请注意修改《数据库编码对应表》;"
		}
		if len(item.DelCatalogue) > 0 {
			msgArr[4] = code + "项目作废" + strings.Join(item.DelCatalogue, "、") + "常量，请注意在《数据库编码对应表》删除该目录信息;"
		}
		if len(item.DelCatalogueItem) > 0 {
			msgArr[5] = code + "项目的" + strings.Join(item.DelCatalogueItem, ";")
		}
		robot := utils.NewRobot(DingTalkAccessToken, DingTalkSecret)
		for _, msg := range msgArr {
			if msg != "" {
				err = robot.SendTextMessage(msg, []string{}, true)
				if err != nil {
					global.GLog.Error("通知失败,消息:" + msg)
					continue
				}
			}
		}
	}
	//刷一下管理常量
	//for _, proCode := range proCodes {
	//	err, constMap := utils.LoadConst(proCode)
	//	if err != nil {
	//		global.GLog.Error("刷新常量失败 proCode: "+proCode, zap.Error(err))
	//		continue
	//	}
	//	conf := global.GProConf[proCode]
	//	conf.ConstTable = constMap
	//	global.GProConf[proCode] = conf
	//	global.GLog.Info("刷新常量成功", zap.Any("", proCode))
	//}
}

// UpdateHospital 更新医疗机构（医院）
func UpdateHospital(namePinYin string, item model2.UpdateConstLog) error {
	//如果是1更新医疗机构（医院）
	var otherInfo model2.HospitalInfo
	err := json.Unmarshal([]byte(item.OtherInfo), &otherInfo)
	if err != nil {
		global.GLog.Error("更新失败id："+item.ID, zap.Error(err))
		return err
	}
	//B0103、B0106结构：arr = {医院名称 医院等级 医院编码}
	arrStr := fmt.Sprintf(`%v","%v`, otherInfo.HospitalName, otherInfo.HospitalCode)
	arr := []string{otherInfo.HospitalName, otherInfo.HospitalCode}
	if utils.RegIsMatch("^(B0110)$", item.ProCode) {
		//B0110结构：arr = {医院名称 医院编码}
		arr = []string{otherInfo.HospitalName, otherInfo.HospitalGrade, otherInfo.HospitalCode}
		arrStr = fmt.Sprintf(`%v","%v","%v`, otherInfo.HospitalName, otherInfo.HospitalGrade, otherInfo.HospitalCode)
	}
	if item.IsDeleted {
		//删除
		err = service.DeleteDoc(namePinYin, `{"selector": {"arr": ["`+arrStr+`"]}}`)
		if err != nil {
			global.GLog.Error("删除失败id："+item.ID, zap.Error(err))
			return err
		}
	} else {
		//新增
		err = service.InsertDoc(namePinYin, arr)
		if err != nil {
			global.GLog.Error("新增失败id："+item.ID, zap.Error(err))
			return err
		}

	}
	return err
}

// UpdateCatalogue 更新医疗目录
func UpdateCatalogue(namePinYin string, item model2.UpdateConstLog) error {
	var items []interface{}
	//第一行中文名称
	itemZHName := map[string]interface{}{
		"_id":     "ChineseTableName::" + uuid.NewV4().String(),
		"content": item.Name,
	}
	items = append(items, itemZHName)

	//第二行表头
	itemHead := map[string]interface{}{
		"_id":      "TableTop::" + uuid.NewV4().String(),
		"tabletop": dbTableTopArr,
	}
	items = append(items, itemHead)

	if !item.IsDeleted {
		//读取excel 获取是否废弃 == N 数据 插入数据
		list, err := GetExcelData(item.LocalUrl)
		//list, err := GetExcelData("/Users/mjl/git/stronger/server/files/const/hospital/650100_N_650100-20230523031213.xlsx")
		if err != nil {
			global.GLog.Error("获取目录表数据失败", zap.Error(err))
			return err
		}
		for _, info := range list {
			items = append(items, map[string]interface{}{
				"_id": info.Id,
				"arr": info.Arr,
			})
		}
	}

	//在lp2删掉一条数据
	//err := service.DeleteDoc("lp2", `{"selector": {"dbName": "`+namePinYin+`"}}`)
	//if err != nil {
	//	global.GLog.Error("删除lp2一条数据id："+item.ID, zap.Error(err))
	//	return err
	//}
	//删除整个目录表
	if has, _ := global.GCouchdbClient.DBExists(context.Background(), namePinYin); has {
		err := global.GCouchdbClient.DestroyDB(context.TODO(), namePinYin)
		if err != nil {
			global.GLog.Error("删除整个目录表id："+item.ID, zap.Error(err))
			return err
		}
	}

	//一条一条删
	//err := service.DeleteDoc(namePinYin, `{"selector": {"_id": {"$regex": "^const::"}},"limit": 1000000}`)
	//if err != nil {
	//	global.GLog.Error("删除lp2一条数据id："+item.ID, zap.Error(err))
	//	return err
	//}

	//上面删掉数据库了 重新建一个
	err := global.GCouchdbClient.CreateDB(context.Background(), namePinYin)
	if err != nil {
		global.GLog.Info("新建数据库失败", zap.Any(namePinYin, err))
		return err
	}

	//构造插入的数据
	var docs []interface{}
	for _, d := range items {
		doc, _ := json.Marshal(d)
		docs = append(docs, doc)
	}

	//插入数据 每次插入5000
	oneCount := 5000
	for i, j := 0, oneCount; i < len(docs); {
		_, err = global.GCouchdbClient.DB(context.Background(), namePinYin).BulkDocs(context.Background(), docs[i:j])
		if err != nil {
			global.GLog.Error("新增有误", zap.Error(err))
			return err
		}
		global.GLog.Info(fmt.Sprintf("新增一次成功 i:j -- %v:%v", i, j))
		i = j
		if j+oneCount >= len(docs) {
			j = len(docs)
		} else {
			j += oneCount
		}
	}
	global.GLog.Info("", zap.Any(namePinYin, "新增成功"))
	return err
}

// UpdateHospitalApp 更新医疗机构（医院）客户端
func UpdateHospitalApp(item model2.UpdateConstLog, msgItem *MsgItem) error {
	//如果是1更新医疗机构（医院）
	var otherInfo model2.HospitalInfo
	err := json.Unmarshal([]byte(item.OtherInfo), &otherInfo)
	if err != nil {
		global.GLog.Error("更新失败id："+item.ID, zap.Error(err))
		return err
	}
	//B0103、B0106结构： {医院名称 医院编码}
	//B0110结构： {医院名称 医院等级 医院编码}
	editJSON := `{"医院名称": "` + otherInfo.HospitalName + `", "医院编码": "` + otherInfo.HospitalCode + `"}`

	if utils.RegIsMatch("^(B0110)$", item.ProCode) {
		editJSON = `{"医院名称": "` + otherInfo.HospitalName + `", "医院编码": "` + otherInfo.HospitalCode + `","医院等级": "` + otherInfo.HospitalGrade + `"}`
	}
	if item.IsDeleted {
		// 删除医疗机构
		//curl -X POST -H "Content-Type: application/json" -d '{"item":editMap,"proCode":"B0113","name":"B0113_百年理赔_百年理赔职业代码表"}' https://example.com/sys-const/del-docs
		delCmd := `curl -X POST --header 'Content-Type: application/json'  --header 'Authorization: ` +
			TOKEN + `' -d '{"item":` +
			editJSON + `,"proCode":"` +
			item.ProCode + `","name":"` +
			item.Name + `"}' ` +
			global.GConfig.System.ConstUrl + `/sys-const/del-docs`
		global.GLog.Info("删除医疗机构 " + delCmd)
		err = ShellOut(delCmd)
		if err != nil {
			global.GLog.Error("删除失败id："+item.ID, zap.Error(err))
			return err
		}
		//推送钉钉消息 删除医疗机构信息
		msgItem.DelHospital = append(msgItem.DelHospital, "删除"+otherInfo.HospitalName+"，代码为"+otherInfo.HospitalCode)
	} else {
		// 新增医疗机构
		//curl -X PUT -H "Content-Type: application/json" -d '{"items":[editMap]}' https://example.com/sys-const/insert/B0113/B0113_百年理赔_百年理赔职业代码表
		//curl --location 'http://192.168.202.18:30080/sys-const/insert' --header 'Authorization: Bearer 5pyq6K6+572uYmFzZTY05qC85byP55qE546v5aKD5Y+Y6YePRFVSSUFOX0FQSV9UT0tFTuWtmOWcqOS4pemHjeeahOWuieWFqOmXrumimA==' --header 'Content-Type: application/json' --data '{"items": [{"医院名称": "广州为民康复医院有限公司","医院编码": "440100111565"}],"name": "B0110_新疆国寿理赔_医疗机构65","proCode": "B0110"}'
		insertCmd := `curl --location --request POST '` +
			global.GConfig.System.ConstUrl + `/sys-const/insert' --header 'Authorization:` +
			TOKEN + `' --header 'Content-Type: application/json' --data '{"items": [` +
			editJSON + `],"name": "` +
			item.Name + `","proCode": "` +
			item.ProCode + `"}'`
		global.GLog.Info("新增医疗机构 " + insertCmd)
		err = ShellOut(insertCmd)
		if err != nil {
			global.GLog.Error("新增失败id："+item.ID, zap.Error(err))
			return err
		}
		//推送钉钉消息 新增医疗机构信息
		msgItem.AddHospital = append(msgItem.AddHospital, "新增"+otherInfo.HospitalName+"，代码为"+otherInfo.HospitalCode)
	}
	return err
}

// UpdateCatalogueApp 更新医疗目录 客户端
func UpdateCatalogueApp(item model2.UpdateConstLog, msgItem *MsgItem) error {
	// 更新医疗目录
	//curl -X PUT -F "files=@/path/to/your/file.txt" https://example.com/sys-const/import/B0113
	//updateCmd := fmt.Sprintf("curl -X PUT --header 'Authorization: %v' -F \"files=@%v\" %v/sys-const/import/%v",
	//	TOKEN, item.LocalUrl, global.GConfig.System.ConstUrl, item.ProCode)
	//global.GLog.Info("更新医疗目录 " + updateCmd)
	//return ShellOut(updateCmd)

	//if !item.IsDeleted {
	//读取excel 获取是否废弃 == N 数据 插入数据
	//item.LocalUrl = "/Users/mjl/git/stronger/server/files/const/hospital/650100_N_650100-20230523031213.xlsx"
	tempDataPath := path.Dir(item.LocalUrl) + "/" + item.Name + ".ndjson"
	_, delItems, err := GetExcelDataApp(item.LocalUrl, tempDataPath)
	if err != nil {
		global.GLog.Error("读取excel失败："+item.ID, zap.Error(err))
		return err
	}
	//body := map[string]interface{}{
	//	"items":   items,
	//	"name":    item.Name,
	//	"proCode": item.ProCode,
	//}
	//data, err := json.Marshal(body)
	//if err != nil {
	//	global.GLog.Error("Marshal有误："+item.ID, zap.Error(err))
	//	return err
	//}
	//err = os.WriteFile(tempDataPath, data, 0666)
	//if err != nil {
	//	global.GLog.Error("创建临时文件失败："+item.ID, zap.Error(err))
	//	return err
	//}

	//updateCmd := fmt.Sprintf("curl -X PUT --header 'Authorization: %v' -F \"files=@%v\" %v/sys-const/import/%v",
	//	TOKEN, item.LocalUrl, global.GConfig.System.ConstUrl, item.ProCode)
	delCmd := `curl -X POST --header 'Content-Type: application/json'  --header 'Authorization: ` +
		TOKEN + `' -d '{"proCode": "` +
		item.ProCode + `","name":["` +
		item.Name + `"]}' ` +
		global.GConfig.System.ConstUrl + `/sys-const/del-tables`
	global.GLog.Info("删除医疗目录 " + delCmd)
	err = ShellOut(delCmd)
	if err != nil && err.Error() != "数据表不存在" {
		global.GLog.Error("删除失败："+item.ID, zap.Error(err))
		return err
	}
	if err != nil && err.Error() == "数据表不存在" {
		//推送钉钉消息 新增医疗目录常量
		msgItem.AddCatalogue = append(msgItem.AddCatalogue, item.Name)
	}

	insertCmd := fmt.Sprintf("curl -X PUT --header 'Authorization: %v' -F \"files=@%v\" %v/sys-const/import/%v",
		TOKEN, tempDataPath, global.GConfig.System.ConstUrl, item.ProCode)
	//insertCmd := `curl --location --request POST '` +
	//	global.GConfig.System.ConstUrl + `/sys-const/insert' --header 'Authorization:` +
	//	TOKEN + `' --header 'Content-Type: application/json' -d @` + tempDataPath
	global.GLog.Info("新增医疗目录 " + insertCmd)
	err = ShellOut(insertCmd)
	if err != nil {
		global.GLog.Error("新增失败id："+item.ID, zap.Error(err))
		//return err
	}

	if item.IsDeleted {
		//推送钉钉消息 作废医疗目录常量
		msgItem.DelCatalogue = append(msgItem.DelCatalogue, item.Name)
	} else {
		//推送钉钉消息 作废医疗目录常量项目
		//XXXX项目的y1常量作废z1、z2、z3共n1条；y2常量作废z4、z5、z6共n2条；  （条件：第一个下划线后的值为N删除文件“是否作废”列中为Y）
		//y1常量作废z1、z2、z3共n1条；
		msgItem.DelCatalogueItem = append(msgItem.DelCatalogueItem, item.Name+"常量作废"+strings.Join(delItems, "、")+"共"+strconv.Itoa(len(delItems))+"条")
	}
	//推送钉钉消息 更新医疗目录常量信息
	msgItem.UpdateCatalogue = append(msgItem.UpdateCatalogue, item.Name)

	//}
	return nil
}

func GetExcelDataApp(path, tempDataPath string) (items []map[string]string, delItems []string, err error) {
	// 打开文件并设置为追加模式 NDJSON
	file, err := os.OpenFile(tempDataPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return items, delItems, err
	}
	defer file.Close()
	// 清空文件内容
	err = file.Truncate(0)
	if err != nil {
		return items, delItems, err
	}

	f, err := excelize.OpenFile(path)
	if err != nil {
		return items, delItems, err
	}
	sheetName := f.GetSheetName(1)
	global.GLog.Info("sheetName", zap.Any("start read", sheetName))
	// 获取excel中具体的列的值
	rows := f.GetRows(sheetName)
	excelTableTopIndexArr := make([]int, len(excelTableTopArr))
	for i, row := range rows {
		if len(row) < len(excelTableTopArr) {
			global.GLog.Error("这一行太短了", zap.Any("index", i))
			continue
		}
		if i == 0 {
			//第一行是表头  获取需要更新的excel的列数 index
			for j, key := range excelTableTopArr {
				for index, r := range row {
					if key == r {
						excelTableTopIndexArr[j] = index
						fmt.Println(excelTableTopIndexArr)
					}
				}
			}
			continue
		}
		if row[excelTableTopIndexArr[len(excelTableTopIndexArr)-1]] == "N" {
			item := map[string]string{}
			for index, header := range dbTableTopArr {
				item[header] = row[excelTableTopIndexArr[index]]
			}
			bytes, _ := json.Marshal(&item)          // 将数据转换成字节切片
			_, _ = fmt.Fprintln(file, string(bytes)) // 将字节切片写入文件
			//items = append(items, item)
		} else {
			//推送钉钉消息 作废医疗目录常量项目
			//XXXX项目的y1常量作废z1、z2、z3共n1条；y2常量作废z4、z5、z6共n2条；  （条件：第一个下划线后的值为N删除文件“是否作废”列中为Y）
			delItems = append(delItems, row[2]) //z1、z2、z3
		}
	}
	global.GLog.Info("sheetName", zap.Any("end read", sheetName))
	return items, delItems, nil
}

func GetExcelData(path string) (list []model.TableInfo, err error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return list, err
	}
	sheetName := f.GetSheetName(1)
	global.GLog.Info("sheetName", zap.Any("", sheetName))
	// 获取excel中具体的列的值
	rows := f.GetRows(sheetName)
	var items []model.TableInfo
	excelTableTopIndexArr := make([]int, len(excelTableTopArr))
	for i, row := range rows {
		if len(row) < len(excelTableTopArr) {
			global.GLog.Error("这一行太短了", zap.Any("index", i))
			continue
		}
		if i == 0 {
			//第一行是表头  获取需要更新的excel的列数 index
			for j, key := range excelTableTopArr {
				for index, r := range row {
					if key == r {
						excelTableTopIndexArr[j] = index
						fmt.Println(excelTableTopIndexArr)
					}
				}
			}
			continue
		}
		arr := make([]string, len(excelTableTopArr))
		for l, index := range excelTableTopIndexArr {
			arr[l] = row[index]
		}
		//本地数据库：最高限价、医疗项目编码、项目名称、项目别名、单位、规格、剂型、自费比例
		//最后一个是 是否作废 == "N" 的才需要插入
		if arr[len(arr)-1] == "N" {
			item := model.TableInfo{
				Id:  "const::" + uuid.NewV4().String(),
				Arr: arr[:len(arr)-1], //最后一个 是否作废 不插入到数据库
			}
			items = append(items, item)
		}
	}
	return items, nil
}

var excelTableTopArr = []string{"最高限价", "项目编码*", "项目名称*", "项目别名", "单位", "规格", "剂型", "自费比例*", "是否作废"}
var dbTableTopArr = []string{"最高限价", "医疗项目编码", "项目名称", "项目别名", "单位", "规格", "剂型", "自费比例"}

type Resp struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func ShellOut(cmd string) error {
	var resp Resp
	cmd = cmd + " -m 60 -f"
	global.GLog.Info("cmd", zap.Any("", cmd))
	err, stdout, stderr := utils.ShellOut(cmd)
	global.GLog.Info("stdout", zap.Any("", stdout))
	if err != nil {
		return err
	}
	if stderr != "" {
		global.GLog.Error("stderr", zap.Any("", stderr))
	}
	err = json.Unmarshal([]byte(stdout), &resp)
	if err != nil {
		return err
	}
	if resp.Status != 200 {
		return errors.New(resp.Msg)
	}
	return nil
}

// MsgItem 钉钉通知信息结构体
type MsgItem struct {
	AddHospital      []string //新增医疗机构信息 eg:新增y1，代码为x1;
	DelHospital      []string //删除医疗机构信息 eg:删除y1，代码为x1;
	UpdateCatalogue  []string //更新医疗目录常量信息 eg:yy1、yy2、yy3
	AddCatalogue     []string //新增医疗目录常量 eg:XXXXX常量，XXXXX常量，
	DelCatalogue     []string //作废医疗目录常量 eg:XXXXX常量，XXXXX常量，
	DelCatalogueItem []string //作废医疗目录常量项目 eg:y1常量作废z1、z2、z3共n1条；y2常量作废z4、z5、z6共n2条；
}
