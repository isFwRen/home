package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"server/module/dingding/model"
	res "server/module/dingding/model/request"
	resp "server/module/dingding/model/response"
)

const (
	AppKey    = "dingcghid7icsasqc7n9"
	AppSecret = "xYoZTH-ufma_nlqaSqjftbI1rsC15NbmcBN0GRfMJpMvIoESA8wXkL0hyis_Iq7T"
	AgentID   = "1329034673"
)

// DdingTalk process
//1、Get the AccessTaken
//2、Get the UserID
//3、Send the message
func DdingTalk(UserPhone, Captcha string) (SuccessMsg string, FailMsg string) {
	AccessTakenS, err := GetAccessTaken()
	if err != "" {
		return "", "获取AccessTaken失败, 错误信息:" + err
	}
	if AccessTakenS.AccessToken != "" {
		DingUserIDResponse := GetUserID(AccessTakenS.AccessToken, UserPhone)
		if DingUserIDResponse.ErrMsg != "" && DingUserIDResponse.ErrMsg != "ok" {
			return "", "获取UserID失败," + DingUserIDResponse.ErrCode + ":" + DingUserIDResponse.ErrMsg
		}
		if DingUserIDResponse.Result.UserId != "" {
			fmt.Println(DingUserIDResponse.Result.UserId)
			DingMessage := res.SendDingMessage{
				AgentId:    AgentID,
				UserIdList: DingUserIDResponse.Result.UserId,
				Msg: res.SendDingMeg{
					MsgType: "text",
					Text:    res.CT{Content: "验证码:" + Captcha + "，1分钟内有效，仅用于汇流登录"},
				},
			}

			//buf, err := json.MarshalIndent(DingMessage, "", " ")
			buf, err := json.Marshal(DingMessage)
			if err != nil {
				panic(err)
			}
			//SendDingMessage := "{\"userid_list\":\"" + DingUserIDResponse.Result.UserId + "\",\"agent_id\":\"1329034673\",\"msg\":{\n\t\"msgtype\":\"text\",\"text\":{\n\t\"content\":\"验证码:" + Captcha + ",1分钟内有效,仅用于汇流登录\"\n}}}"
			DingTalk(AccessTakenS.AccessToken, buf)
		} else {
			return "", "UserId为空"
		}
	} else {
		return "", "AccessTaken为空"
	}
	return
}

// GetAccessTaken Get access to enterprise internal AccessTaken
//request method：GET
//request url：https://oapi.dingtalk.com/gettoken
func GetAccessTaken() (model.DingAccessTaken, string) {
	var accesstaken model.DingAccessTaken
	url := "https://oapi.dingtalk.com/gettoken?" + "appkey=" + AppKey + "&appsecret=" + AppSecret
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &accesstaken)
	if accesstaken.AccessToken != "" {
		return accesstaken, ""
	}
	return accesstaken, accesstaken.ErrMsg
}

// GetUserID Obtain the user ID according to the mobile phone number
//request method：POST
//request url：https://oapi.dingtalk.com/topapi/v2/user/getbymobile
func GetUserID(AccessTaken, UserPhone string) resp.DingUserIDResponse {
	var DingUserIDResponse resp.DingUserIDResponse
	url := "https://oapi.dingtalk.com/topapi/v2/user/getbymobile?access_token=" + AccessTaken
	client := &http.Client{}
	DingUserMessage := model.DingUserMessage{
		Mobile: UserPhone,
	}
	DM, _ := json.Marshal(DingUserMessage)
	res, err := http.NewRequest("POST", url, bytes.NewReader(DM))
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(res)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &DingUserIDResponse)
	return DingUserIDResponse
}

// DingTalk Send the message to DingUser
//request method：POST
//request url：https://oapi.dingtalk.com/topapi/message/corpconversation/asyncsend_v2
//DingMessage must be byte arr of SendDingMessage struct
func DingTalk(AccessTaken string, DingMessage []byte) resp.SendDingMessageResponse {
	var SendDingMessageResponse resp.SendDingMessageResponse
	url := "https://oapi.dingtalk.com/topapi/message/corpconversation/asyncsend_v2?access_token=" + AccessTaken
	client := &http.Client{}
	//res, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(SendDingMessage)))
	res, err := http.NewRequest("POST", url, bytes.NewBuffer(DingMessage))
	res.Header.Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(res)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err1 := json.Unmarshal(body, &SendDingMessageResponse)
	if err1 != nil {
		fmt.Println(err1)
	}
	fmt.Println("DingTakl1", SendDingMessageResponse)
	fmt.Println("DingTalk", SendDingMessageResponse.ErrMsg)
	return SendDingMessageResponse
}
