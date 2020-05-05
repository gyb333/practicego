package message

const (
	LoginMsgType            = "LoginMsgType"
	LoginRstMsgType         = "LoginRstMsgType"
	RegisterMsgType         = "RegisterMsgType"
	RegisterRstMsgType      = "RegisterRstMsgType"
	NotifyUserStatusMsgType = "NotifyUserStatusMsgType"
	SmsMsgType              = "SmsMsgType"
	LoginExistMsgType 		="LoginExistMsgType"
	LoginExistRstMsgType 		="LoginExistRstMsgType"
)
const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

type Message struct {
	Type string `json:"type"` //消息类型
	Data string `json:"data"` //消息类型数据
}

type LoginMsg struct {
	UserId   int    `json:"userId"`
	UserName string `json:"userName"`
	UserPwd  string `json:"userPwd"`
}

type LoginRstMsg struct {
	Code    int `json:"code"`
	UsersId []int
	Error   string `json:"error"`
}

type RegisterMsg struct {
	User User `json:"user"`
}

type RegisterRstMsg struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

//增加一个SmsMes //发送的消息
type SmsMes struct {
	Content string `json:"content"` //内容
	User           //匿名结构体，继承
}

//为了配合服务器端推送用户状态变化的消息
type NotifyUserStatusMes struct {
	UserId int `json:"userId"` //用户id
	Status int `json:"status"` //用户的状态
}
