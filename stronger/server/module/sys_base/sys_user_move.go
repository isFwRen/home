/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/1/4 10:12 上午
 */

package sys_base

import (
	"context"
	"fmt"
	"server/global"
	"server/module/sys_base/model"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

const (
	pageSize      int64 = 1  //每页大小 500 点击同步的话这个值必须为1
	gCount        int   = 10 //协程数量 7
	eachGDealPage int64 = 13 //每个协程处理的页数 13
)

func Move(maxId string, isFirst bool) (int64, int64) {
	now := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	global.GLog.Warn("MongoDB.Url:::" + global.GConfig.MongoDB.Url)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(global.GConfig.MongoDB.Url))
	if err != nil {
		global.GLog.Error("1111" + err.Error())
	}
	collection := client.Database("sys_user").Collection("users")

	//获取新的最大的id
	//获取最新时间
	var maxResult model.SysUserCopy
	err = collection.FindOne(ctx,
		bson.D{
			{"nickname", bson.D{{"$ne", nil}}},
			{"username", bson.D{{"$ne", "-1"}}},
		},
		//options.FindOne().SetSort(bson.D{{"_id", -1}}),
		options.FindOne().SetSort(bson.D{{"update_at", -1}}),
	).Decode(&maxResult)
	if err != nil {
		global.GLog.Error("9999" + err.Error())
		return 0, 0
	}

	//数据mongo复制到pg
	var wg sync.WaitGroup
	var countUser, countRole int64
	for l := 0; l < gCount; l++ {
		j := int64(l)
		wg.Add(1)
		countUserChan := make(chan int64, 45)
		countRoleChan := make(chan int64, 45)
		go func(j *int64) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					global.GLog.Error("复制异常", zap.Any("", err))
				}
			}()
			var i int64
			for i = eachGDealPage * (*j); i < eachGDealPage*((*j)+1); i++ {
				s := i * pageSize
				var rowUser int64 = 0
				var rowRole int64 = 0
				findAndAdd(collection, ctx, &s, maxId, isFirst, &rowUser, &rowRole)
				global.GLog.Warn("新增条数：：：" + fmt.Sprintf("%d----%d", rowUser, rowRole))
				countUserChan <- rowUser
				countRoleChan <- rowRole
			}
			close(countUserChan)
			close(countRoleChan)
		}(&j)
		go func() {
			for {
				n, ok := <-countUserChan
				n1, ok1 := <-countRoleChan
				if !ok && !ok1 {
					break
				}
				countUser = countUser + n
				countRole = countRole + n1
			}
		}()
	}

	wg.Wait()

	//更新新的最大的id
	err = global.GDb.Model(model.MaxUserId{}).Where("id = '1'").Updates(map[string]interface{}{
		"max_id":     maxResult.ObjectId,
		"updated_at": maxResult.UpdateAt,
	}).Error
	if err != nil {
		global.GLog.Error(err.Error())
	}

	ss := time.Since(now).Milliseconds()

	fmt.Println(fmt.Sprintf("countUser::%d-----countRole::%d", countUser, countRole))
	fmt.Println("花费：：：" + strconv.FormatInt(ss, 10))
	return countUser, countRole
}

func findAndAdd(collection *mongo.Collection, ctx context.Context, pageIndex *int64, maxId string,
	isFirst bool, rowUser *int64, rowRole *int64) {
	filterBson := bson.D{
		{"nickname", bson.D{{"$ne", nil}}},
		{"username", bson.D{{"$ne", "-1"}}},
	}
	if !isFirst {
		//var b primitive.ObjectID
		//err := b.UnmarshalText([]byte(maxId))
		//if err != nil {
		//	global.GLog.Error("1111" + err.Error())
		//}
		filterBson = bson.D{
			{"nickname", bson.D{{"$ne", nil}}},
			//{"_id", bson.D{{"$gt", bsonx.ObjectID(b)}}},
			{"update_at", bson.D{{"$gt", maxId}}},
		}
	}

	cur, err := collection.Find(ctx, filterBson,
		options.Find().SetSkip(*pageIndex).SetLimit(pageSize).SetSort(bson.D{{"_id", 1}}),
	)
	if err != nil {
		global.GLog.Error("22222" + err.Error())
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err = cur.Close(ctx)
		if err != nil {
			global.GLog.Error(err.Error())
		}
	}(cur, ctx)
	var results1 []model.SysUserCopy
	if err = cur.All(context.TODO(), &results1); err != nil {
		global.GLog.Error(err.Error())
	}
	if results1 == nil {
		return
	}
	//点击同步的话pageSize=1
	if len(results1) == 1 && maxId != "" {
		// fmt.Println("-------------------111111--------------------:", results1[0])
		// 在`object_id`冲突时，将列更新为新值
		*rowUser = global.GDb.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "object_id"}},
			//DoUpdates: clause.AssignmentColumns([]string{"code", "password","nick_name"}),
			DoUpdates: clause.AssignmentColumns([]string{"code", "password", "nick_name", "status", "phone",
				"is_mobile", "entry_date", "mount_guard_date", "leave_date", "referees", "qrcode"}),
		}).Select("id", "code", "password", "nick_name", "header_img", "phone",
			"ding_id", "is_mobile", "referees", "entry_date", "mount_guard_date", "leave_date",
			"status", "staff", "email", "sex", "id_card", "reason", "object_id", "bank_id",
			"bank_nick_name", "bank_name", "bank_branch", "address", "profession", "wechat",
			"educational", "qrcode").Create(&results1[0]).RowsAffected
		return
	}
	//上面是同步的代码
	//下面是第一次move_data的
	tx := global.GDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			global.GLog.Error("", zap.Any("", r))
		}
	}()
	//fmt.Println(results1)
	//results2 := userCopyToUser(results1)
	if len(results1) > 0 {
		// fmt.Println("-------------------22222222222--------------------:", results1)
		*rowUser = tx.
			Select("id", "code", "password", "nick_name", "header_img", "phone",
				"ding_id", "is_mobile", "referees", "entry_date", "mount_guard_date", "leave_date",
				"status", "staff", "email", "sex", "id_card", "reason", "object_id", "bank_id",
				"bank_nick_name", "bank_name", "bank_branch", "address", "profession", "wechat",
				"educational", "role_id", "qrcode").
			Create(&results1).RowsAffected
	}
	var sysProPermissions []model.SysProPermission
	for _, userCopy := range results1 {
		//fmt.Println(userCopy.Roles)
		for _, role := range userCopy.Roles {
			var sysProPermission model.SysProPermission
			role.UserCode = userCopy.Code
			role.UserId = userCopy.ID
			role.ObjectId = userCopy.ObjectId
			role.ProId = global.ProCodeId[role.ProCode]
			if role.UserCode[:1] == "P" {
				role.HasOutNet = true
				role.HasInNet = false
			} else {
				role.HasOutNet = false
				role.HasInNet = true
			}
			sysProPermission = role
			sysProPermissions = append(sysProPermissions, sysProPermission)
		}
	}
	if len(sysProPermissions) > 0 {
		*rowRole = tx.Model(model.SysProPermission{}).Create(&sysProPermissions).RowsAffected
	}
	err = tx.Commit().Error

	if err != nil {
		global.GLog.Error(err.Error())
	}
	if err = cur.Err(); err != nil {
		global.GLog.Error("4444" + err.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			global.GLog.Error("all", zap.Any("", r))
		}
	}()
	//global.GLog.Warn("新增条数：：：" + fmt.Sprintf("%d----%d", rowUser, rowRole))
}

func userCopyToUser(user []model.SysUserCopy) (users []model.SysUser) {
	for _, sysUser := range user {
		var user1 model.SysUser
		user1.ID = sysUser.ID
		user1.Code = sysUser.Code
		user1.Password = sysUser.Password
		user1.NickName = sysUser.NickName
		user1.HeaderImg = sysUser.HeaderImg
		user1.Phone = sysUser.Phone
		user1.DingId = sysUser.DingId
		user1.IsMobile = sysUser.IsMobile
		user1.Referees = sysUser.Referees
		user1.EntryDate = time.Time(sysUser.EntryDate)
		user1.MountGuardDate = time.Time(sysUser.MountGuardDate)
		user1.LeaveDate = time.Time(sysUser.LeaveDate)
		user1.Status = sysUser.Status
		user1.Staff = sysUser.Staff
		user1.Email = sysUser.Email
		user1.Sex = bool(sysUser.Sex)
		user1.IDCard = sysUser.IDCard
		user1.Reason = sysUser.Reason
		user1.ObjectId = sysUser.ObjectId
		user1.BankId = sysUser.BankId
		user1.BankNickName = sysUser.BankNickName
		user1.BankName = sysUser.BankName
		user1.BankBranch = sysUser.BankBranch
		user1.Address = sysUser.Address
		user1.Profession = sysUser.Profession
		user1.Wechat = sysUser.Wechat
		user1.Educational = sysUser.Educational
		user1.RoleId = sysUser.RoleId
		users = append(users, user1)
	}
	return users
}
