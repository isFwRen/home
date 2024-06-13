package model

import (
	"server/module/sys_base/model"
	"time"

	"github.com/lib/pq"
)

//循环分块是同一个分块，用field的BlockIndex区分
//初审切出非初审分块 field.blockIndex = block.zero

type ProjectBlock struct {
	// ID
	model.Model
	BillID      string         `json:"billID" gorm:"comment:'单ID'"`                              //单ID
	Temp        string         `json:"temp" gorm:"comment:'类型'"`                                 //类型
	Name        string         `json:"name" gorm:"模板分块名字"`                                       //模板分块名字
	Code        string         `json:"code" gorm:"模板分块编码"`                                       //模板分块编码
	FEight      bool           `json:"fEight" gorm:"是否f8提交"`                                     //是否f8提交
	Ocr         string         `json:"ocr" gorm:"ocr流程"`                                         //ocr流程
	FreeTime    int            `json:"freeTime" gorm:"释放时间"`                                     //释放时间
	IsLoop      bool           `json:"isLoop" gorm:"是否循环分块"`                                     //是否循环分块
	PicPage     int            `json:"picPage" gorm:"图片页码，第几张图片"`                                //图片页码，第几张图片
	IsMobile    bool           `json:"isMobile" gorm:"移动端可录入"`                                   //移动端可录入
	WCoordinate pq.StringArray `json:"wCoordinate" gorm:"type:varchar(100)[];comment:'web截图位置'"` //web截图位置
	MCoordinate pq.StringArray `json:"mCoordinate" gorm:"type:varchar(100)[];comment:'移动端截图位置'"` //移动端截图位置
	MPicPage    int            `json:"mPicPage" gorm:"移动端图片页码，第几张图片"`                            //移动端图片页码，第几张图片
	PreBCode    pq.StringArray `json:"preBCode" gorm:"type:varchar(100)[];前置分块编码"`               //前置分块编码
	LinkBCode   pq.StringArray `json:"linkBCode" gorm:"type:varchar(100)[];参考分块编码"`              //参考分块编码
	Op1PreBCode pq.StringArray `json:"Op1PreBCode" gorm:"type:varchar(100)[];一码前置分块编码"`          //一码前置分块编码
	//
	Op0Stage string `json:"op0Stage" gorm:"comment:'状态(op0|op0Cache|op|op1|op2|op1)'"` //状态
	// Mark    string `json:"mark" gorm:"comment:'状态(1|2|done)'"`                     //状态
	Status  int    `json:"status" gorm:"comment:'1(初审)|2'"` //1(初审)|2
	Zero    int    `json:"zero" gorm:"comment:'下标'"`        //下标
	Picture string `json:"picture" gorm:"图片"`
	Level   int    `json:"level" gorm:"comment:'优先级'"`

	Op1Stage      string `json:"op1Stage" gorm:"comment:'状态1(opCache|op1|op1Cache|done)'"` //状态
	Op2Stage      string `json:"op2Stage" gorm:"comment:'状态2'"`                            //状态
	OpqStage      string `json:"opqStage" gorm:"comment:'状态q'"`                            //状态
	IsCompetitive bool   `json:"isCompetitive" gorm:"comment:'状态q'"`

	Op1Code     string    `json:"op1Code" gorm:"comment:'1码人员编号(0为系统录入)'"` //1码人员编号(0为系统录入)
	Op1ApplyAt  time.Time `json:"op1ApplyAt" gorm:"comment:'1码领取时间'"`      //1码领取时间
	Op1SubmitAt time.Time `json:"op1SubmitAt" gorm:"comment:'1码提交时间'"`     //1码提交时间
	Op2Code     string    `json:"op2Code" gorm:"comment:'2码人员编号'"`         //2码人员编号
	Op2ApplyAt  time.Time `json:"op2ApplyAt" gorm:"comment:'2码领取时间'"`      //2码领取时间
	Op2SubmitAt time.Time `json:"op2SubmitAt" gorm:"comment:'2码提交时间'"`     //2码提交时间
	OpqCode     string    `json:"opqCode" gorm:"comment:'问题件人员编号'"`        //问题件人员编号
	OpqApplyAt  time.Time `json:"opqApplyAt" gorm:"comment:'问题件领取时间'"`     //问题件领取时间
	OpqSubmitAt time.Time `json:"opqSubmitAt" gorm:"comment:'问题件提交时间'"`    //问题件提交时间
	Op0Code     string    `json:"op0Code" gorm:"comment:'初审人员编号'"`         //初审人员编号
	Op0ApplyAt  time.Time `json:"op0ApplyAt" gorm:"comment:'初审领取时间'"`      //初审领取时间
	Op0SubmitAt time.Time `json:"op0SubmitAt" gorm:"comment:'初审提交时间'"`     //初审提交时间

	IsPractice bool   `json:"isPractice" gorm:"comment:'练习分块'"` //练习分块
	Crypto     string `json:"crypto" gorm:"comment:'crypto'"`   //分块秘钥
}
