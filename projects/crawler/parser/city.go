package parser

import (
	"../engine"
	"regexp"
	)

type CityParser struct {
}

//const cityRe = `<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`

var (
	cityRe    = regexp.MustCompile(`<a href="(http://album.zhenai.com/u/([0-9]+))"[^>]*>([^<]+)</a>`)
	cityUrlRe = regexp.MustCompile(`<a href="(http://www.zhenai.com/zhenghun/[^"]+)"`)
)

func (p CityParser) Parse (contents []byte) engine.ParseResult {
	profiles := cityRe.FindAllSubmatch(contents, -1)

	result := engine.ParseResult{}
	for _, c := range profiles {
		name:=string(c[3])
		url :=string(c[1])
		//userID:=string(c[2])
		//result.Items = append(result.Items, engine.Item{
		//	Url:url ,
		//	Type:"zhenai:city",
		//	Id:userID,
		//	Payload:name,
		//}) //用户名字

		result.Requests = append(result.Requests, engine.Request{
			Url: url       ,
			Parser: ProfileParser{Url:url,Name:name},
		})
	}

	//爬取下一页
	cityUrls := cityUrlRe.FindAllSubmatch(contents, -1)
	for _,cityUrl :=range cityUrls{
		result.Requests = append(result.Requests,engine.Request{
			Url:string(cityUrl[1]),
			Parser:CityParser{},
		})
	}

	return result
}

func (p CityParser) Serialize() (name string, args interface{}){
	return "CityParser",nil
}