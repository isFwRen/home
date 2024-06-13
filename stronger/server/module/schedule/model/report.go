package model

import "server/module/sys_base/model"

type BusinessAnalysis struct {
	model.Model
	//ID                 string `json:"id"`
	CreateAt           string `json:"create_at"`
	AHalfPastEight     int32  `json:"a_half_past_eight"`
	AHalfPastNine      int32  `json:"a_half_past_nine"`
	AHalfPastTen       int32  `json:"a_half_past_ten"`
	AHalfPastEleven    int32  `json:"a_half_past_eleven"`
	AHalfPastTwelve    int32  `json:"a_half_past_twelve"`
	AHalfPastThirteen  int32  `json:"a_half_past_thirteen"`
	AHalfPastFourteen  int32  `json:"a_half_past_fourteen"`
	AHalfPastFifteen   int32  `json:"a_half_past_fifteen"`
	AHalfPastSixteen   int32  `json:"a_half_past_sixteen"`
	AHalfPastSeventeen int32  `json:"a_half_past_seventeen"`
	AHalfPastEighteen  int32  `json:"a_half_past_eighteen"`
	BeforeZeroHour     int32  `json:"before_zero_hour"`
	Types              string `json:"types"`
}

type BusinessUploadAnalysis struct {
	model.Model
	//ID                     string  `json:"id"`
	CreateAt               string  `json:"createAt" example:"创建时间"`
	VolumeOfBusiness       int32   `json:"volumeOfBusiness" example:"业务量"`
	TheAverageTime         float64 `json:"theAverageTime" example:"平均时长"`
	OneHours               int32   `json:"oneHours" example:"0-1小时内回传数量"`
	TwoHours               int32   `json:"twoHours" example:"1-2小时内回传数量"`
	ThreeHours             int32   `json:"threeHours" example:"2-3小时内回传数量"`
	MoreThanThreeHours     int32   `json:"moreThanThreeHours" example:"3小时以上回传数量"`
	OneHoursRate           float64 `json:"oneHoursRate" example:"0-1小时内回传比例"`
	TwoHoursRate           float64 `json:"twoHoursRate" example:"1-2小时内回传比例"`
	ThreeHoursRate         float64 `json:"threeHoursRate" example:"2-3小时内回传比例"`
	MoreThanThreeHoursRate float64 `json:"moreThanThreeHoursRate" example:"3小时以上回传比例"`
}
