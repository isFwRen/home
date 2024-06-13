package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	//"gorm.io/gorm"
)

type UserRoleRelationShip struct {
	Model
	UserId string `json:"userId"`
	RoleId string `json:"roleId"`
}

type UserManagement struct {
	SysUser
	//SysProPermission []SysProPermission `json:"sysProPermission"`
	Role string `json:"role"`
}

type SysUser struct {
	Model
	Code           string    `bson:"username" json:"code" gorm:"comment:'工号'"`
	Password       string    `bson:"password" json:"-"  gorm:"comment:'用户登录密码'"`
	NickName       string    `bson:"nickname" json:"nickName" gorm:"default:'系统用户';comment:'姓名'"`
	HeaderImg      string    `bson:"headerImg" json:"headerImg" gorm:"default:'https://s4.ax1x.com/2022/01/17/7aIpHe.jpg';comment:'用户头像'"`
	Phone          string    `bson:"phone" json:"phone" gorm:"comment:'手机号'"`
	DingId         string    `bson:"dingId" json:"dingId" gorm:"comment:'钉钉id'"`
	IsMobile       bool      `bson:"isMobile" json:"isMobile" gorm:"是否手机端"`
	Referees       string    `bson:"recom_username" json:"referees" gorm:"推荐人"`
	EntryDate      time.Time `bson:"entryDate" json:"entryDate" gorm:"comment:'入职日期'"`
	MountGuardDate time.Time `bson:"mountGuardDate" json:"mountGuardDate" gorm:"上岗日期"`
	LeaveDate      time.Time `bson:"leaveDate" json:"leaveDate" gorm:"comment:'离职日期'"`
	Status         bool      `bson:"status" json:"status" gorm:"comment:'状态'"`
	Staff          string    `bson:"staffRole" json:"staff" gorm:"comment:'职位'"`
	Email          string    `bson:"email" json:"email" gorm:"comment:'邮箱'"`
	Sex            bool      `bson:"sex" json:"sex" gorm:"comment:'性别'"`
	//SysRoles       []SysRoles `bson:"sysRoles" json:"sysRoles" gorm:"foreignKey:UserId;comment:权限"`
	IDCard       string    `bson:"IDCard" json:"idCard" gorm:"comment:身份证"`
	Reason       string    `bson:"reason" json:"reason"`
	ObjectId     string    `bson:"_id" json:"objectId"`
	BankId       string    `bson:"bankId" json:"bankId"`
	BankNickName string    `bson:"bankNickName" json:"bankNickName"`
	BankName     string    `bson:"bankName" json:"bankName"`
	BankBranch   string    `bson:"bankBlanch" json:"bankBranch"`
	Address      string    `bson:"address" json:"address"`
	Profession   string    `bson:"profession" json:"profession"`
	Wechat       string    `bson:"wechat" json:"wechat"`
	Educational  string    `bson:"educational" json:"educational"`
	RoleId       string    `json:"roleId" gorm:"default:'926052346658553856';comment:'角色'"` //角色
	SysRoles     SysRoles  `json:"sysRoles" gorm:"foreignKey:RoleId;references:ID;"`        //角色
	Qrcode       string    `bson:"qrcode" json:"qrcode"`
	PasswordDate time.Time `bson:"passwordDate" json:"passwordDate" gorm:"comment:'密码更新日期'"`
}

type SysCaptcha struct {
	Model
	CaptchaId string `json:"captcha_id"`
}

type SysUserCopy struct {
	Model
	Code           string   `bson:"username" json:"code" gorm:"comment:'工号'"`
	Password       string   `bson:"password" json:"x-"  gorm:"comment:'用户登录密码'"`
	NickName       string   `bson:"nickname" json:"nickName" gorm:"default:'系统用户';comment:'姓名'"`
	HeaderImg      string   `bson:"headerImg" json:"headerImg" gorm:"default:'https://imgtu.com/i/T20KVP';comment:'用户头像'"`
	Phone          string   `bson:"phone" json:"phone" gorm:"comment:'手机号'"`
	DingId         string   `bson:"dingId" json:"dingId" gorm:"comment:'钉钉id'"`
	IsMobile       bool     `bson:"isMobile" json:"isMobile" gorm:"是否手机端"`
	Referees       string   `bson:"recom_username" json:"referees" gorm:"推荐人"`
	EntryDate      JsonTime `bson:"entryDate" json:"entryDate" gorm:"comment:'入职日期'"`
	MountGuardDate JsonTime `bson:"workDate" json:"mountGuardDate" gorm:"上岗日期"`
	LeaveDate      JsonTime `bson:"leaveDate" json:"leaveDate" gorm:"comment:'离职日期'"`
	Status         bool     `bson:"enabled" json:"status" gorm:"comment:'状态'"`
	Staff          string   `bson:"staffRole" json:"staff" gorm:"comment:'职位'"`
	Email          string   `bson:"email" json:"email" gorm:"comment:'邮箱'"`
	Sex            MySex    `bson:"sex" json:"sex" gorm:"comment:'性别'"`
	IDCard         string   `bson:"IDCard" json:"idCard" gorm:"comment:身份证"`
	Reason         string   `bson:"reason" json:"reason"`
	ObjectId       string   `bson:"_id" json:"objectId"`

	BankId       string   `bson:"bank_id" json:"bankId"`
	BankNickName string   `bson:"bank_nick_name" json:"bankNickName"`
	BankName     string   `bson:"bank_name" json:"bankName"`
	BankBranch   string   `bson:"bank_blanch" json:"bankBranch"`
	Address      string   `bson:"address" json:"address"`
	Profession   string   `bson:"profession" json:"profession"`
	Wechat       string   `bson:"wechat" json:"wechat"`
	Educational  string   `bson:"educational" json:"educational"`
	RoleId       string   `json:"roleId" gorm:"default:'926052346658553856';comment:'角色'"` //角色
	Roles        MyObject `bson:"role" json:"roles" gorm:"-"`
	UpdateAt     string   `bson:"update_at" json:"updateAt" gorm:"-"`
	Qrcode       string   `bson:"qrcode" json:"qrcode"`
}

func (v SysUserCopy) TableName() string {
	return "sys_users"
}

type JsonTime time.Time
type MySex bool
type MyObject []SysProPermission

const (
	timeFormat = "2006-01-02"
)

// UnmarshalJSON 实现json反序列化，从传递的字符串中解析成时间对象
func (t *JsonTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormat+`"`, string(data), time.Local)
	*t = JsonTime(now)
	return
}

// MarshalJSON 实现json序列化，将时间转换成字符串byte数组
func (t JsonTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormat)
	b = append(b, '"')
	return b, nil
}

// MarshalBSONValue mongodb是存储bson格式，因此需要实现序列化bsonvalue(这里不能实现MarshalBSON，MarshalBSON是处理Document的)，将时间转换成mongodb能识别的primitive.DateTime
func (t *JsonTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	targetTime := primitive.NewDateTimeFromTime(time.Time(*t))
	return bson.MarshalValue(targetTime)
}

// UnmarshalBSONValue 实现bson反序列化，从mongodb中读取数据转换成time.Time格式，这里用到了bsoncore中的方法读取数据转换成datetime然后再转换成time.Time
func (t *JsonTime) UnmarshalBSONValue(t2 bsontype.Type, data []byte) error {
	_, _, valid := bsoncore.ReadValue(data, t2)
	if valid == false {
		return errors.New(fmt.Sprintf("%s, %s, %s", "读取数据失败:", t2, data))
	}
	s, _, _ := bsoncore.ReadString(data)
	if s == "" {
		s = "0001-01-01"
	}
	pattern := "^\\d{4}-\\d{2}-\\d{2}$"
	pattern1 := "^\\d{4}-\\d{1}-\\d{2}$"
	pattern2 := "^\\d{4}-\\d{2}-\\d{1}$"
	pattern3 := "^\\d{4}-\\d{1}-\\d{1}$"
	pattern4 := "^\\d{6}$"
	pattern5 := "^\\d{2}-\\d{2}-\\d{2}$"
	pattern6 := "^\\d{4}/\\d{2}/\\d{2}$"
	pattern7 := "^\\d{8}$"
	pattern8 := "^\\d{4}/\\d{1}/\\d{2}$"
	pattern9 := "^\\d{4}/\\d{2}/\\d{1}$"
	pattern10 := "^\\d{4}/\\d{1}/\\d{1}$"
	layout := "2006-01-02"
	if match, _ := regexp.MatchString(pattern, s); match {
		layout = "2006-01-02"
	} else if match1, _ := regexp.MatchString(pattern1, s); match1 {
		layout = "2006-1-02"
	} else if match2, _ := regexp.MatchString(pattern2, s); match2 {
		layout = "2006-01-2"
	} else if match3, _ := regexp.MatchString(pattern3, s); match3 {
		layout = "2006-1-2"
	} else if match4, _ := regexp.MatchString(pattern4, s); match4 {
		layout = "060102"
	} else if match5, _ := regexp.MatchString(pattern5, s); match5 {
		layout = "06-01-02"
	} else if match6, _ := regexp.MatchString(pattern6, s); match6 {
		layout = "2006/01/02"
	} else if match7, _ := regexp.MatchString(pattern7, s); match7 {
		layout = "20060102"
	} else if match8, _ := regexp.MatchString(pattern8, s); match8 {
		layout = "2006/1/02"
	} else if match9, _ := regexp.MatchString(pattern9, s); match9 {
		layout = "2006/01/2"
	} else if match10, _ := regexp.MatchString(pattern10, s); match10 {
		layout = "2006/1/2"
	} else {
		s = "0001-01-01"
	}
	now, err := time.ParseInLocation(layout, s, time.Local)
	if err != nil {
		fmt.Println(err.Error())
		now, _ = time.ParseInLocation("2006-01-02", "0001-01-01", time.Local)
		*t = JsonTime(now)
	}
	*t = JsonTime(now)
	return nil
}

func (t JsonTime) Value() (driver.Value, error) {
	// JsonTime 转换成 time.Time 类型
	tTime := time.Time(t)
	return tTime.Format(yyyyMMddHHmmss), nil
}

func (t *JsonTime) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		// 字符串转成 time.Time 类型
		*t = JsonTime(vt)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}

func (t *JsonTime) String() string {
	return fmt.Sprintf("hhh:%s", time.Time(*t).String())
}

// UnmarshalBSONValue 实现bson反序列化，从mongodb中读取数据转换成bool格式
func (t *MySex) UnmarshalBSONValue(t2 bsontype.Type, data []byte) error {
	_, _, valid := bsoncore.ReadValue(data, t2)
	if valid == false {
		return errors.New(fmt.Sprintf("%s, %s, %s", "读取数据失败:", t2, data))
	}
	if t2.String() == "string" {
		s, _, _ := bsoncore.ReadString(data)
		if s == "女" {
			*t = MySex(false)
		} else {
			*t = MySex(true)
		}
	} else if t2.String() == "bool" {
		s, _, _ := bsoncore.ReadBoolean(data)
		*t = MySex(s)
	}
	return nil
}

func (t MySex) Value() (driver.Value, error) {
	// MySex 转换成 bool 类型
	sex := bool(t)
	return sex, nil
}

type ProRole struct {
	B0111  []string `json:"云南国寿理赔" bson:"云南国寿理赔"`
	B0108  []string `json:"太平理赔" bson:"太平理赔"`
	B0105  []string `json:"广西贵州国寿理赔" bson:"广西贵州国寿理赔"`
	B0104  []string `json:"海南陕西国寿理赔" bson:"海南陕西国寿理赔"`
	B0101  []string `json:"民生理赔" bson:"民生理赔"`
	B0106  []string `json:"陕西国寿理赔" bson:"陕西国寿理赔"`
	B0110  []string `json:"新疆国寿理赔" bson:"新疆国寿理赔"`
	B0113  []string `json:"百年理赔" bson:"百年理赔"`
	B0114  []string `json:"华夏理赔" bson:"华夏理赔"`
	B0115  []string `json:"新华理赔" bson:"新华理赔"`
	B0116  []string `json:"华夏人寿团险理赔" bson:"华夏人寿团险理赔"`
	B0118  []string `json:"中意理赔" bson:"中意理赔"`
	B0117  []string `json:"北大方正理赔" bson:"北大方正理赔"`
	B01031 []string `json:"新广西贵州国寿理赔" bson:"新广西贵州国寿理赔"`
	B0120  []string `json:"北京国寿理赔" bson:"北京国寿理赔"`
	B0121  []string `json:"百年人寿团险理赔" bson:"百年人寿团险理赔"`
	B0401  []string `json:"北大方正体检理赔" bson:"北大方正体检理赔"`
	B01151 []string `json:"新华理赔清单练习" bson:"新华理赔清单练习"`
}

// UnmarshalBSONValue 实现bson反序列化，从mongodb中读取数据转换成json格式
func (t *MyObject) UnmarshalBSONValue(t2 bsontype.Type, data []byte) error {
	_, _, valid := bsoncore.ReadValue(data, t2)
	if valid == false {
		return errors.New(fmt.Sprintf("%s, %s, %s", "读取数据失败:", t2, data))
	}
	var sysProPermissions []SysProPermission
	if t2.String() == "embedded document" {
		s, _, _ := bsoncore.ReadDocument(data)
		//fmt.Println(s)
		c := s.String()
		var proRole ProRole
		err := json.Unmarshal([]byte(c), &proRole)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		num := reflect.TypeOf(proRole).NumField()
		for i := 0; i < num; i++ {
			var sysProPermission SysProPermission
			code := reflect.TypeOf(proRole).Field(i).Name
			name := reflect.TypeOf(proRole).Field(i).Tag.Get("json")
			sysProPermission.ProCode = code
			sysProPermission.ProName = name
			//sysProPermission.ProId = global.ProCodeId[code]

			val := reflect.ValueOf(proRole).FieldByName(code)
			for _, s2 := range val.Interface().([]string) {
				switch s2 {
				case "op1":
					sysProPermission.HasOp0 = true
					sysProPermission.HasOp1 = true
				case "op2":
					sysProPermission.HasOp2 = true
				case "opQ":
					sysProPermission.HasOpq = true
				//case "out_net":
				//	sysProPermission.HasOutNet = true
				//case "in_net":
				//	sysProPermission.HasInNet = true
				case "pm":
					sysProPermission.HasPm = true
				}
			}
			if sysProPermission.HasOp0 ||
				sysProPermission.HasOp1 ||
				sysProPermission.HasOp2 ||
				sysProPermission.HasOpq ||
				sysProPermission.HasInNet ||
				sysProPermission.HasOutNet ||
				sysProPermission.HasPm {
				sysProPermissions = append(sysProPermissions, sysProPermission)
			}
		}

	} else {
		fmt.Println("读取数据类型失败：：：" + t2.String())
		s, _, _ := bsoncore.ReadString(data)
		fmt.Println(s)
		//return errors.New("读取数据类型失败")
	}
	*t = sysProPermissions
	return nil
}

func (t MyObject) Value() (driver.Value, error) {
	// MyObject 转换成 bool 类型
	//perm := []SysProPermission(t)
	return nil, nil
}
