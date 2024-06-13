package request

import "time"

type RegisterStruct struct {
	EntryTime        string `json:"entry_time"`
	CreatedAt        time.Time
	Username         string `json:"username"`
	IDCard           string `json:"IDCard"`
	Address          string `json:"address"`
	BankBlanch       string `json:"bankBlanch"`
	BankId           string `json:"bankId"`
	BankNickname     string `json:"bankNickname"`
	Duration         string `json:"duration"`
	Email            string `json:"email"`
	Nickname         string `json:"nickname"`
	Other            string `json:"other"`
	OtherWay         string `json:"otherWay"`
	Phone            string `json:"phone"`
	QQ               string `json:"qq"`
	RecomID          string `json:"recomId"`
	Sex              string `json:"sex"`
	Time             string `json:"time"`
	Times            string `json:"times"`
	Way              string `json:"way"`
	MarriageValue    string `json:"marriage_value"`
	Bank             string `json:"bank"`
	EducationalValue string `json:"educational_value"`
	ProfessionValue  string `json:"profession_value"`
	SelectValue      string `json:"select_value"`
	ImgA             string `json:"imga"`
	ImgB             string `json:"imgb"`
	Team             string `json:"team"`
	State            string `json:"state"`
	WorkTime         string `json:"work_time"`
	Reason           string `json:"reason"`
	StartTime        string `json:"start_time"`
	DepartureTime    string `json:"departure_time"`
	X1               string `json:"x_1"`
	X2               string `json:"x_2"`
	X3               string `json:"x_3"`
	X4               string `json:"x_4"`
	X5               string `json:"x_5"`
}


