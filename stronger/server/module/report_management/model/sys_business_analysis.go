package model

type BusinessAnalysisSearch struct {
	ProCode    string `json:"proCode" form:"proCode" example:"项目编码"` //all 或者 B0114    //form:"proCode"
	StartTime  string `json:"startTime" form:"startTime" example:"开始时间"`
	EndTime    string `json:"endTime" form:"endTime" example:"结束时间"`
	Types      string `json:"types" form:"types" example:"类型,2种:来量-download/回传-upload"`
	IsCheckAll bool   `json:"isCheckAll" form:"isCheckAll"` //全部true 或者 明细false
}

type AllBusinessUploadAnalysis struct {
	ProCode                string  `json:"proCode" example:"项目编码" excel:"项目编码"`
	CreateAt               string  `json:"createAt" example:"创建时间" excel:"时间"`
	VolumeOfBusiness       int32   `json:"volumeOfBusiness" example:"业务量" excel:"业务量"`
	TheAverageTime         float64 `json:"theAverageTime" example:"平均时长" excel:"平均时长"`
	OneHours               int32   `json:"oneHours" example:"0-1小时内回传数量" excel:"0-1小时内回传数量"`
	TwoHours               int32   `json:"twoHours" example:"1-2小时内回传数量" excel:"1-2小时内回传数量"`
	ThreeHours             int32   `json:"threeHours" example:"2-3小时内回传数量" excel:"2-3小时内回传数量"`
	MoreThanThreeHours     int32   `json:"moreThanThreeHours" example:"3小时以上回传数量" excel:"3小时以上回传数量"`
	OneHoursRate           string  `json:"oneHoursRate" example:"0-1小时内回传比例" excel:"0-1小时内回传比例"`
	TwoHoursRate           string  `json:"twoHoursRate" example:"1-2小时内回传比例" excel:"1-2小时内回传比例"`
	ThreeHoursRate         string  `json:"threeHoursRate" example:"2-3小时内回传比例" excel:"2-3小时内回传比例"`
	MoreThanThreeHoursRate string  `json:"moreThanThreeHoursRate" example:"3小时以上回传比例" excel:"3小时以上回传比例"`
}

type AllBusinessAnalysis struct {
	ProCode            string `json:"proCode" example:"项目编码" excel:"项目编码"`
	CreateAt           string `json:"createAt" example:"创建时间" excel:"时间"`
	AHalfPastEight     int32  `json:"aHalfPastEight" example:"0.00-8.30" excel:"0.00-8.30"`
	AHalfPastNine      int32  `json:"aHalfPastNine" example:"8.30-9.30" excel:"8.30-9.30"`
	AHalfPastTen       int32  `json:"aHalfPastTen" example:"9.30-10.30" excel:"9.30-10.30"`
	AHalfPastEleven    int32  `json:"aHalfPastEleven" example:"10.30-11.30" excel:"10.30-11.30"`
	AHalfPastTwelve    int32  `json:"aHalfPastTwelve" example:"11.30-12.30" excel:"11.30-12.30"`
	AHalfPastThirteen  int32  `json:"aHalfPastThirteen" example:"12.30-13.30" excel:"12.30-13.30"`
	AHalfPastFourteen  int32  `json:"aHalfPastFourteen" example:"13.30-14.30" excel:"13.30-14.30"`
	AHalfPastFifteen   int32  `json:"aHalfPastFifteen" example:"14.30-15.30" excel:"14.30-15.30"`
	AHalfPastSixteen   int32  `json:"aHalfPastSixteen" example:"15.30-16.30" excel:"15.30-16.30"`
	AHalfPastSeventeen int32  `json:"aHalfPastSeventeen" example:"16.30-17.30" excel:"16.30-17.30"`
	AHalfPastEighteen  int32  `json:"aHalfPastEighteen" example:"17.30-18.30" excel:"17.30-18.30"`
	BeforeZeroHour     int32  `json:"beforeZeroHour" example:"18.30-0.00" excel:"18.30-0.00"`
}

type BusinessAnalysis struct {
	CreateAt           string `json:"createAt" example:"创建时间" excel:"时间"`
	AHalfPastEight     int32  `json:"aHalfPastEight" example:"0.00-8.30" excel:"0.00-8.30"`
	AHalfPastNine      int32  `json:"aHalfPastNine" example:"8.30-9.30" excel:"8.30-9.30"`
	AHalfPastTen       int32  `json:"aHalfPastTen" example:"9.30-10.30" excel:"9.30-10.30"`
	AHalfPastEleven    int32  `json:"aHalfPastEleven" example:"10.30-11.30" excel:"10.30-11.30"`
	AHalfPastTwelve    int32  `json:"aHalfPastTwelve" example:"11.30-12.30" excel:"11.30-12.30"`
	AHalfPastThirteen  int32  `json:"aHalfPastThirteen" example:"12.30-13.30" excel:"12.30-13.30"`
	AHalfPastFourteen  int32  `json:"aHalfPastFourteen" example:"13.30-14.30" excel:"13.30-14.30"`
	AHalfPastFifteen   int32  `json:"aHalfPastFifteen" example:"14.30-15.30" excel:"14.30-15.30"`
	AHalfPastSixteen   int32  `json:"aHalfPastSixteen" example:"15.30-16.30" excel:"15.30-16.30"`
	AHalfPastSeventeen int32  `json:"aHalfPastSeventeen" example:"16.30-17.30" excel:"16.30-17.30"`
	AHalfPastEighteen  int32  `json:"aHalfPastEighteen" example:"17.30-18.30" excel:"17.30-18.30"`
	BeforeZeroHour     int32  `json:"beforeZeroHour" example:"18.30-0.00" excel:"18.30-0.00"`
	Types              string `json:"types" example:"类型,2种:来量-download/回传-upload"`
}

type BusinessUploadAnalysis struct {
	CreateAt               string  `json:"createAt" example:"创建时间" excel:"创建时间"`
	VolumeOfBusiness       int32   `json:"volumeOfBusiness" example:"业务量" excel:"业务量"`
	TheAverageTime         float64 `json:"theAverageTime" example:"平均时长" excel:"平均时长"`
	OneHours               int32   `json:"oneHours" example:"0-1小时内回传数量" excel:"0-1小时内回传数量"`
	TwoHours               int32   `json:"twoHours" example:"1-2小时内回传数量" excel:"1-2小时内回传数量"`
	ThreeHours             int32   `json:"threeHours" example:"2-3小时内回传数量" excel:"2-3小时内回传数量"`
	MoreThanThreeHours     int32   `json:"moreThanThreeHours" example:"3小时以上回传数量" excel:"3小时以上回传数量"`
	OneHoursRate           string  `json:"oneHoursRate" example:"0-1小时内回传比例" excel:"0-1小时内回传比例"`
	TwoHoursRate           string  `json:"twoHoursRate" example:"1-2小时内回传比例" excel:"1-2小时内回传比例"`
	ThreeHoursRate         string  `json:"threeHoursRate" example:"2-3小时内回传比例" excel:"2-3小时内回传比例"`
	MoreThanThreeHoursRate string  `json:"moreThanThreeHoursRate" example:"3小时以上回传比例" excel:"3小时以上回传比例"`
}
