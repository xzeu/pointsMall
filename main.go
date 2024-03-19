package main

import (
	"fmt"
	"pointMall/service"
	"pointMall/setting"
	"strings"
	"time"
)

func main() {
	//1. 加载配置文件
	setting.Viper()
	//fmt.Println(setting.Cnf)
	conf := setting.Cnf.Account
	for _, userInfo := range conf {
		fmt.Println(strings.Repeat("=", 60))
		fmt.Println("用户：", userInfo.UserName)
		service.Login(userInfo.UserName, userInfo.PassWd)
		token := service.GetAuthorization()
		fmt.Println("token:", token)
		//签到
		service.SignIn(token)
		err, respJson := service.DailyQA(token)
		if err != nil {
			fmt.Println(err)
			// panic(err)
			// return
		} else {
			//等待5-10秒再提交
			//timeNterval := rand.Intn(5) + 6
			fmt.Println(userInfo.TimeInterval)
			time.Sleep(time.Duration(userInfo.TimeInterval) * time.Second)
			// 提交答案
			answerUrl := "https://xiaoyou.dgtis.com/admin/answer/submit"
			answerResult := service.Postbin(answerUrl, token, respJson)
			fmt.Println(string(answerResult))
		}
	}
}
