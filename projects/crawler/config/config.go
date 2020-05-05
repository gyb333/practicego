package config

import (
	"../../utils"
	"github.com/spf13/viper"
)


func save()  {
	var maps =make( map[string]interface{})
	maps["Qps"]=20
	maps["Url"] = "https://www.zhenai.com/zhenghun"
	maps["CityUrl"] = "http://www.zhenai.com/zhenghun/aba"
	maps["ProfileUrl"] = "http://album.zhenai.com/u/1280064210"
	utils.SaveConfig("src/crawler/config.toml",maps)
}



var (

	WorkerPort0   = 5566 //Worker

	ItemSaverPort =":7788"


	CityParser = "CityParser"
	CityListParser = "CityListParser"
	ProfileParser = "ProfileParser"
	NilParser = "NilParser"

	ElasticIndex = "dating_profile"

	//RPC远程服务要调用的方法名称
	ItemSaverRpc    = "ItemSaveService.Save"
	CrawlServiceRpc = "CrawlService.Process"

	_= utils.ReadConfig("src/crawler/config.toml")
	//Rate limiting
	Qps =viper.GetInt("Qps")
	Url=viper.GetString("Url")
	CityUrl=viper.GetString("CityUrl")
	ProfileUrl=viper.GetString("ProfileUrl")

)