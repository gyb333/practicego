package process

import (
	"../../common/message"
	"encoding/json"
	"log"
)

func outputGroupMes(mes *message.Message) { //这个地方mes一定SmsMes
	//显示即可
	//1. 反序列化mes.Data
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		log.Println("json.Unmarshal err=", err.Error())
		return
	}

	//显示信息
	log.Printf("用户id:\t%d 对大家说:\t%s\n",
		smsMes.UserId, smsMes.Content)

}
