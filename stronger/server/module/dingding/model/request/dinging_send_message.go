package request

/*
SendDingMessage :
	AgentId    发送消息时使用的微应用的AgentID。
	UserIdList 接收者的userid列表，最大用户列表长度100。
	DeptIdList 接收者的部门id列表，最大列表长度20。接收者是部门ID时，包括子部门下的所有用户。
	ToAllUser  是否发送给企业全部用户。
	Msg        消息内容，最长不超过2048个字节，支持以下消息类型：
						文本消息
						图片消息
						语音消息
						文件消息
						链接消息
						OA消息 注意 OA消息支持通过status_bar参数设置消息的状态文案和颜色，消息发送后可调用更新工作通知状态栏接口更新消息状态和颜色。
						Markdown消息
						卡片消息
*/
type SendDingMessage struct {
	AgentId    string `json:"agent_id"`
	UserIdList string      `json:"userid_list"`
	//DeptIdList string      `json:"dept_id_list"`
	ToAllUser  bool        `json:"to_all_user"`
	Msg        SendDingMeg `json:"msg"`
}

type SendDingMeg struct {
	MsgType string `json:"msgtype"`
	Text    CT     `json:"text"`
}

type CT struct {
	Content string `json:"content"`
}
