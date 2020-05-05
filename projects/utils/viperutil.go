package utils

import (
	"github.com/spf13/viper"
	"fmt"
	)

func SaveConfig(pathname string,kv map[string]interface{})  {
	viper.SetConfigFile(pathname)//文件名
	for k,v :=range kv{
		//viper.Set("Address","0.0.0.0:9090")//统一把Key处理成小写 Address->address
		viper.Set(k,v)
	}

	err := viper.WriteConfig()//写入文件
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func ReadConfig(pathname string) error {
	viper.SetConfigFile(pathname)
	err := viper.ReadInConfig() // 会查找和读取配置文件
	if err != nil { // Handle errors reading the config file
		return fmt.Errorf("Fatal error config file: %s \n", err)
	}
	return nil
}